package bnc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/props"
)

type KlineWs struct {
	spClt *klineMergedWs
	umClt *klineMergedWs
	cmClt *klineMergedWs
}

func NewKlineWs(ctx context.Context, logger *slog.Logger) *KlineWs {
	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("ws", "bnc_kl_clt")
	return &KlineWs{
		spClt: newKlineMergedWs(ctx, cex.SYMBOL_TYPE_SPOT, logger),
		umClt: newKlineMergedWs(ctx, cex.SYMBOL_TYPE_UM_FUTURES, logger),
		cmClt: newKlineMergedWs(ctx, cex.SYMBOL_TYPE_CM_FUTURES, logger),
	}
}

func (kw *KlineWs) Sub(symbolType cex.SymbolType, interval KlineInterval, symbols ...string) (unsubed []string, err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return kw.spClt.subKlines(interval, symbols...)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return kw.umClt.subKlines(interval, symbols...)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return kw.cmClt.subKlines(interval, symbols...)
	}
	return nil, errors.New("bnc: unknown symbol type")
}

func (kw *KlineWs) Unsub(symbolType cex.SymbolType, interval KlineInterval, symbols ...string) (err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return kw.spClt.unsubKlines(interval, symbols...)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return kw.umClt.unsubKlines(interval, symbols...)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return kw.cmClt.unsubKlines(interval, symbols...)
	}
	return errors.New("bnc: unknown symbol type")
}

func (kw *KlineWs) NewCh(symbolType cex.SymbolType, symbol string, interval KlineInterval) (ch <-chan KlineMsg, err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return kw.spClt.newCh(symbol, interval)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return kw.umClt.newCh(symbol, interval)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return kw.cmClt.newCh(symbol, interval)
	}
	return nil, errors.New("bnc: unknown symbol type")
}

func (kw *KlineWs) RemoveCh(ch <-chan KlineMsg) {
	kw.spClt.removeCh(ch)
	kw.umClt.removeCh(ch)
	kw.cmClt.removeCh(ch)
}

type klineMergedWs struct {
	mu sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc

	symbolType cex.SymbolType

	clts       []*klineBaseWs
	wsByStream *props.SafeRWMap[string, *klineBaseWs]

	logger *slog.Logger
}

func newKlineMergedWs(ctx context.Context, symbolType cex.SymbolType, logger *slog.Logger) *klineMergedWs {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("bnc_kl_merged_clt", symbolType)
	return &klineMergedWs{
		ctx:        ctx,
		cancel:     cancel,
		symbolType: symbolType,
		wsByStream: props.NewSafeRWMap[string, *klineBaseWs](),
		logger:     logger,
	}
}

func (kw *klineMergedWs) scanAllSymbols() {
	for _, clt := range kw.clts {
		for _, s := range clt.ws.Stream() {
			kw.wsByStream.SetIfNotExists(s, clt)
		}
	}
}

func (kw *klineMergedWs) subKlines(interval KlineInterval, symbols ...string) (unsubed []string, err error) {
	kw.mu.Lock()
	defer kw.mu.Unlock()
	defer kw.scanAllSymbols()
	for _, s := range symbols {
		if !kw.wsByStream.HasKey(klineStream(s, interval)) {
			unsubed = append(unsubed, s)
		}
	}
	for _, clt := range kw.clts {
		if len(unsubed) == 0 {
			return
		}
		unsubed, err = clt.sub(interval, unsubed...)
		if err != nil && err != ErrNotAllStreamSubed {
			kw.logger.Error("bnc: sub klines failed", "err", err)
			return
		}
	}
	for {
		if len(unsubed) == 0 {
			return
		}
		clt := newKlineBaseWs(kw.ctx, kw.symbolType, kw.logger)
		unsubed, err = clt.sub(interval, unsubed...)
		if err != nil && err != ErrNotAllStreamSubed {
			kw.logger.Error("bnc: sub klines failed", "err", err)
			return
		}
		kw.clts = append(kw.clts, clt)
	}
}

func (kw *klineMergedWs) unsubKlines(interval KlineInterval, symbols ...string) (err error) {
	kw.mu.Lock()
	defer kw.mu.Unlock()
	defer kw.scanAllSymbols()
	for _, s := range symbols {
		clt, ok := kw.wsByStream.GetVWithOk(klineStream(s, interval))
		if ok {
			err = clt.unsub(interval, s)
			if err != nil {
				return
			}
		}
	}
	return nil
}

func (kw *klineMergedWs) newCh(symbol string, interval KlineInterval) (ch <-chan KlineMsg, err error) {
	clt, ok := kw.wsByStream.GetVWithOk(klineStream(symbol, interval))
	if !ok {
		return nil, errors.New("bnc: symbol not found")
	}
	return clt.newCh(symbol, interval), nil
}

func (kw *klineMergedWs) removeCh(ch <-chan KlineMsg) {
	for _, clt := range kw.clts {
		clt.removeCh(ch)
	}
}

type KlineMsg struct {
	SymbolType cex.SymbolType
	Kline      WsKline
	Err        error
}

type klineBaseWs struct {
	ctx    context.Context
	cancel context.CancelFunc

	sybType cex.SymbolType

	ws       *RawWsClient
	rawMsgCh <-chan RawWsClientMsg

	radio *props.Radio[KlineMsg]

	logger *slog.Logger
}

func newKlineBaseWs(ctx context.Context, sybType cex.SymbolType, logger *slog.Logger) *klineBaseWs {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("bnc_kl_base_clt", sybType)
	wsCfg := DefaultPublicWsCfg(sybType)
	wsCfg.MaxStream = maxObWsStream
	wsCfg.FanoutTimerDur = 0
	ws := NewRawWsClient(ctx, wsCfg, logger)
	radio := props.NewRadio(
		props.WithFanoutLogger[KlineMsg](logger),
		props.WithFanoutChCap[KlineMsg](10000),
	)
	clt := &klineBaseWs{
		ctx:      ctx,
		cancel:   cancel,
		ws:       ws,
		rawMsgCh: ws.Sub(),
		sybType:  sybType,
		radio:    radio,
		logger:   logger,
	}
	clt.run()
	return clt
}

func (clt *klineBaseWs) createSubParams(interval KlineInterval, symbols ...string) (params []string) {
	for _, symbol := range symbols {
		params = append(params, klineStream(symbol, interval))
	}
	return params
}

func (clt *klineBaseWs) sub(interval KlineInterval, symbols ...string) (unsubed []string, err error) {
	params := clt.createSubParams(interval, symbols...)
	res, err := clt.ws.SubStream(params...)
	if err == nil {
		return
	}
	for _, p := range res.UnsubedParams {
		for _, s := range symbols {
			if strings.Contains(p, strings.ToLower(s)) &&
				strings.Contains(p, "@kline_"+string(interval)) {
				unsubed = append(unsubed, s)
				break
			}
		}
	}
	return unsubed, err
}

func (clt *klineBaseWs) unsub(interval KlineInterval, symbols ...string) (err error) {
	params := clt.createSubParams(interval, symbols...)
	err = clt.ws.UnsubStream(params...)
	if err != nil {
		return err
	}
	for _, symbol := range symbols {
		clt.radio.Broadcast(klineStream(symbol, interval), KlineMsg{
			SymbolType: clt.sybType,
			Kline:      WsKline{Symbol: symbol, Interval: interval},
			Err:        ErrWsStreamUnsubed,
		})
	}
	return nil
}

func (clt *klineBaseWs) newCh(symbol string, interval KlineInterval) <-chan KlineMsg {
	return clt.radio.Sub(klineStream(symbol, interval))
}

func (clt *klineBaseWs) removeCh(ch <-chan KlineMsg) {
	clt.radio.UnsubAll(ch)
}

func (clt *klineBaseWs) run() {
	go clt.listener()
}

func (clt *klineBaseWs) listener() {
	for {
		select {
		case <-clt.ctx.Done():
			return
		case msg := <-clt.rawMsgCh:
			klineMsg := KlineMsg{
				SymbolType: clt.sybType,
			}
			if msg.Err != nil {
				klineMsg.Err = msg.Err
				clt.radio.BroadcastAll(klineMsg)
				continue
			}
			var stream WsKlineStream
			if err := json.Unmarshal(msg.Data, &stream); err != nil {
				klineMsg.Err = err
				clt.radio.BroadcastAll(klineMsg)
				continue
			}
			kline := stream.Kline
			klineMsg.Kline = kline
			clt.radio.Broadcast(klineStream(kline.Symbol, kline.Interval), klineMsg)
		}
	}
}

func klineStream(symbol string, interval KlineInterval) string {
	return fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval)
}
