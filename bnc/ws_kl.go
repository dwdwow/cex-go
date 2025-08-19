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

func NewKlineWs(logger *slog.Logger) *KlineWs {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("ws", "bnc_kl_clt")
	return &KlineWs{
		spClt: newKlineMergedWs(cex.SYMBOL_TYPE_SPOT, logger),
		umClt: newKlineMergedWs(cex.SYMBOL_TYPE_UM_FUTURES, logger),
		cmClt: newKlineMergedWs(cex.SYMBOL_TYPE_CM_FUTURES, logger),
	}
}

func (kw *KlineWs) Sub(symbolType cex.SymbolType, interval KlineInterval, symbols ...string) (unsubed []string, err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return kw.spClt.sub(interval, symbols...)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return kw.umClt.sub(interval, symbols...)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return kw.cmClt.sub(interval, symbols...)
	}
	return nil, errors.New("bnc: unknown symbol type")
}

func (kw *KlineWs) Unsub(symbolType cex.SymbolType, interval KlineInterval, symbols ...string) (err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return kw.spClt.unsub(interval, symbols...)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return kw.umClt.unsub(interval, symbols...)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return kw.cmClt.unsub(interval, symbols...)
	}
	return errors.New("bnc: unknown symbol type")
}

func (kw *KlineWs) NewCh(symbolType cex.SymbolType, interval KlineInterval, symbol string) (ch <-chan KlineMsg, err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return kw.spClt.newCh(interval, symbol)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return kw.umClt.newCh(interval, symbol)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return kw.cmClt.newCh(interval, symbol)
	}
	return nil, errors.New("bnc: unknown symbol type")
}

func (kw *KlineWs) RemoveCh(ch <-chan KlineMsg) {
	kw.spClt.removeCh(ch)
	kw.umClt.removeCh(ch)
	kw.cmClt.removeCh(ch)
}

func (kw *KlineWs) Close() {
	kw.spClt.close()
	kw.umClt.close()
	kw.cmClt.close()
}

type klineMergedWs struct {
	mu sync.Mutex

	symbolType cex.SymbolType

	clts       []*klineBaseWs
	wsByStream *props.SafeRWMap[string, *klineBaseWs]

	logger *slog.Logger
}

func newKlineMergedWs(symbolType cex.SymbolType, logger *slog.Logger) *klineMergedWs {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("bnc_kl_merged_clt", symbolType)
	return &klineMergedWs{
		symbolType: symbolType,
		wsByStream: props.NewSafeRWMap[string, *klineBaseWs](),
		logger:     logger,
	}
}

func (kw *klineMergedWs) scanAllStream() {
	for _, clt := range kw.clts {
		for _, s := range clt.rawWs.Streams() {
			kw.wsByStream.SetIfNotExists(s, clt)
		}
	}
}

func (kw *klineMergedWs) sub(interval KlineInterval, symbols ...string) (unsubed []string, err error) {
	kw.mu.Lock()
	defer kw.mu.Unlock()
	defer kw.scanAllStream()
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
			return
		}
	}
	for {
		if len(unsubed) == 0 {
			return
		}
		var clt *klineBaseWs
		clt, err = startNewKlineBaseWs(kw.symbolType, kw.logger)
		if err != nil {
			return
		}
		unsubed, err = clt.sub(interval, unsubed...)
		if err != nil && err != ErrNotAllStreamSubed {
			return
		}
		kw.clts = append(kw.clts, clt)
	}
}

func (kw *klineMergedWs) unsub(interval KlineInterval, symbols ...string) (err error) {
	kw.mu.Lock()
	defer kw.mu.Unlock()
	defer kw.scanAllStream()
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

func (kw *klineMergedWs) newCh(interval KlineInterval, symbol string) (ch <-chan KlineMsg, err error) {
	clt, ok := kw.wsByStream.GetVWithOk(klineStream(symbol, interval))
	if !ok {
		return nil, errors.New("bnc: symbol not found")
	}
	return clt.newCh(interval, symbol), nil
}

func (kw *klineMergedWs) removeCh(ch <-chan KlineMsg) {
	for _, clt := range kw.clts {
		clt.removeCh(ch)
	}
}

func (kw *klineMergedWs) close() {
	for _, clt := range kw.clts {
		clt.close()
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

	rawWs *RawWs

	radio *props.Radio[KlineMsg]

	logger *slog.Logger
}

func startNewKlineBaseWs(sybType cex.SymbolType, logger *slog.Logger) (ws *klineBaseWs, err error) {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("bnc_kl_base_clt", sybType)
	wsCfg := DefaultPublicWsCfg(sybType)
	wsCfg.MaxStream = maxObWsStream
	rawWs, err := StartNewRawWs(wsCfg, logger)
	if err != nil {
		return
	}
	radio := props.NewRadio(
		props.WithFanoutLogger[KlineMsg](logger),
		props.WithFanoutChCap[KlineMsg](10000),
	)
	ctx, cancel := context.WithCancel(context.Background())
	ws = &klineBaseWs{
		ctx:     ctx,
		cancel:  cancel,
		rawWs:   rawWs,
		sybType: sybType,
		radio:   radio,
		logger:  logger,
	}
	ws.run()
	return
}

func (clt *klineBaseWs) createSubParams(interval KlineInterval, symbols ...string) (params []string) {
	for _, symbol := range symbols {
		params = append(params, klineStream(symbol, interval))
	}
	return params
}

func (clt *klineBaseWs) sub(interval KlineInterval, symbols ...string) (unsubed []string, err error) {
	params := clt.createSubParams(interval, symbols...)
	res, err := clt.rawWs.Sub(params...)
	if err == nil {
		return
	}
	for _, p := range res.UnsubedStreams {
		for _, s := range symbols {
			if p == klineStream(s, interval) {
				unsubed = append(unsubed, s)
				break
			}
		}
	}
	return unsubed, err
}

func (clt *klineBaseWs) unsub(interval KlineInterval, symbols ...string) (err error) {
	params := clt.createSubParams(interval, symbols...)
	_, err = clt.rawWs.Unsub(params...)
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

func (clt *klineBaseWs) newCh(interval KlineInterval, symbol string) <-chan KlineMsg {
	return clt.radio.Sub(klineStream(symbol, interval))
}

func (clt *klineBaseWs) removeCh(ch <-chan KlineMsg) {
	clt.radio.UnsubAll(ch)
}

func (clt *klineBaseWs) run() {
	go clt.listener()
}

func (clt *klineBaseWs) close() {
	clt.cancel()
	clt.rawWs.Close()
}

func (clt *klineBaseWs) listener() {
	for {
		msg := clt.rawWs.Wait()
		if errors.Is(msg.Err, ErrWsClientClosed) {
			return
		}
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
		if stream.EventType == WsEventKline {
			clt.radio.Broadcast(klineStream(kline.Symbol, kline.Interval), klineMsg)
			continue
		}
		if stream.EventType != "" {
			// should not happen
			klineMsg.Err = fmt.Errorf("unknown event type msg: %s", string(msg.Data))
			if kline.Symbol != "" && kline.Interval != "" {
				clt.radio.Broadcast(klineStream(kline.Symbol, kline.Interval), klineMsg)
			} else {
				clt.radio.BroadcastAll(klineMsg)
			}
			continue
		}
		var resp WsResp[any]
		if err := json.Unmarshal(msg.Data, &resp); err == nil && resp.Error != nil {
			klineMsg.Err = fmt.Errorf("bnc: %d %s", resp.Error.Code, resp.Error.Msg)
			clt.radio.BroadcastAll(klineMsg)
			continue
		}
		clt.logger.Warn("unhandled kline ws msg", "msg", string(msg.Data))
	}
}

func klineStream(symbol string, interval KlineInterval) string {
	return fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval)
}
