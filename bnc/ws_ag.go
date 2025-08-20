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

type AggTradeWs struct {
	spClt *aggTradeMergedWs
	umClt *aggTradeMergedWs
	cmClt *aggTradeMergedWs
}

func NewAggTradeWs(logger *slog.Logger) *AggTradeWs {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("ws", "bnc_ag_ws")
	return &AggTradeWs{
		spClt: newAggTradeMergedWs(cex.SYMBOL_TYPE_SPOT, logger),
		umClt: newAggTradeMergedWs(cex.SYMBOL_TYPE_UM_FUTURES, logger),
		cmClt: newAggTradeMergedWs(cex.SYMBOL_TYPE_CM_FUTURES, logger),
	}
}

func (aw *AggTradeWs) Sub(symbolType cex.SymbolType, symbols ...string) (unsubed []string, err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return aw.spClt.sub(symbols...)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return aw.umClt.sub(symbols...)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return aw.cmClt.sub(symbols...)
	}
	return nil, errors.New("bnc: unknown symbol type")
}

func (aw *AggTradeWs) Unsub(symbolType cex.SymbolType, symbols ...string) (err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return aw.spClt.unsub(symbols...)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return aw.umClt.unsub(symbols...)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return aw.cmClt.unsub(symbols...)
	}
	return errors.New("bnc: unknown symbol type")
}

func (aw *AggTradeWs) NewCh(symbolType cex.SymbolType, symbol string) (ch <-chan AggTradeMsg, err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return aw.spClt.newCh(symbol)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return aw.umClt.newCh(symbol)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return aw.cmClt.newCh(symbol)
	}
	return nil, errors.New("bnc: unknown symbol type")
}

func (aw *AggTradeWs) RemoveCh(ch <-chan AggTradeMsg) {
	aw.spClt.removeCh(ch)
	aw.umClt.removeCh(ch)
	aw.cmClt.removeCh(ch)
}

func (aw *AggTradeWs) Close() {
	aw.spClt.close()
	aw.umClt.close()
	aw.cmClt.close()
}

type aggTradeMergedWs struct {
	mu sync.Mutex

	symbolType cex.SymbolType

	clts       []*aggTradeBaseWs
	wsByStream *props.SafeRWMap[string, *aggTradeBaseWs]

	logger *slog.Logger
}

func newAggTradeMergedWs(symbolType cex.SymbolType, logger *slog.Logger) *aggTradeMergedWs {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	mw := &aggTradeMergedWs{
		symbolType: symbolType,
		clts:       make([]*aggTradeBaseWs, 0),
		wsByStream: props.NewSafeRWMap[string, *aggTradeBaseWs](),
		logger:     logger,
	}
	return mw
}

func (mw *aggTradeMergedWs) scanAllStream() {
	for _, clt := range mw.clts {
		for _, s := range clt.rawWs.Streams() {
			mw.wsByStream.SetIfNotExists(s, clt)
		}
	}
}

func (mw *aggTradeMergedWs) sub(symbols ...string) (unsubed []string, err error) {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	defer mw.scanAllStream()
	for _, s := range symbols {
		if !mw.wsByStream.HasKey(aggTradeStream(s)) {
			unsubed = append(unsubed, s)
		}
	}
	for _, clt := range mw.clts {
		if len(unsubed) == 0 {
			return
		}
		unsubed, err = clt.sub(unsubed...)
		if err != nil && err != ErrNotAllStreamSubed {
			return
		}
	}
	for {
		if len(unsubed) == 0 {
			return
		}
		var clt *aggTradeBaseWs
		clt, err = startNewAggTradeBaseWs(mw.symbolType, mw.logger)
		if err != nil {
			return
		}
		unsubed, err = clt.sub(unsubed...)
		if err != nil && err != ErrNotAllStreamSubed {
			return
		}
		mw.clts = append(mw.clts, clt)
	}
}

func (mw *aggTradeMergedWs) unsub(symbols ...string) (err error) {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	defer mw.scanAllStream()
	for _, s := range symbols {
		clt, ok := mw.wsByStream.GetVWithOk(aggTradeStream(s))
		if ok {
			err = clt.unsub(s)
			if err != nil {
				return
			}
		}
	}
	return nil
}

func (mw *aggTradeMergedWs) newCh(symbol string) (ch <-chan AggTradeMsg, err error) {
	clt, ok := mw.wsByStream.GetVWithOk(aggTradeStream(symbol))
	if !ok {
		return nil, errors.New("bnc: symbol not found")
	}
	return clt.newCh(symbol), nil
}

func (mw *aggTradeMergedWs) removeCh(ch <-chan AggTradeMsg) {
	for _, clt := range mw.clts {
		clt.removeCh(ch)
	}
}

func (mw *aggTradeMergedWs) close() {
	for _, clt := range mw.clts {
		clt.close()
	}
}

type AggTradeMsg struct {
	SymbolType cex.SymbolType
	AggTrade   WsAggTradeStream
	Err        error
}

type aggTradeBaseWs struct {
	ctx    context.Context
	cancel context.CancelFunc

	sybType cex.SymbolType

	rawWs *RawWs

	radio *props.Radio[AggTradeMsg]

	logger *slog.Logger
}

func startNewAggTradeBaseWs(sybType cex.SymbolType, logger *slog.Logger) (ws *aggTradeBaseWs, err error) {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("bnc_ag_base_clt", sybType)
	wsCfg := DefaultPublicWsCfg(sybType)
	wsCfg.MaxStream = maxObWsStream
	rawWs, err := StartNewRawWs(wsCfg, logger)
	if err != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	radio := props.NewRadio(
		props.WithFanoutLogger[AggTradeMsg](logger),
		props.WithFanoutChCap[AggTradeMsg](10000),
	)
	ws = &aggTradeBaseWs{
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

func (clt *aggTradeBaseWs) createSubParams(symbols ...string) (params []string) {
	for _, symbol := range symbols {
		params = append(params, aggTradeStream(symbol))
	}
	return params
}

func (clt *aggTradeBaseWs) sub(symbols ...string) (unsubed []string, err error) {
	params := clt.createSubParams(symbols...)
	res, err := clt.rawWs.Sub(params...)
	if err == nil {
		return
	}
	for _, p := range res.UnsubedStreams {
		for _, s := range symbols {
			if p == aggTradeStream(s) {
				unsubed = append(unsubed, s)
				break
			}
		}
	}
	return unsubed, err
}

func (clt *aggTradeBaseWs) unsub(symbols ...string) (err error) {
	params := clt.createSubParams(symbols...)
	_, err = clt.rawWs.Unsub(params...)
	if err != nil {
		return err
	}
	for _, symbol := range symbols {
		clt.radio.Broadcast(aggTradeStream(symbol), AggTradeMsg{
			SymbolType: clt.sybType,
			AggTrade:   WsAggTradeStream{Symbol: symbol},
			Err:        ErrWsStreamUnsubed,
		})
	}
	return nil
}

func (clt *aggTradeBaseWs) newCh(symbol string) <-chan AggTradeMsg {
	return clt.radio.Sub(aggTradeStream(symbol))
}

func (clt *aggTradeBaseWs) removeCh(ch <-chan AggTradeMsg) {
	clt.radio.UnsubAll(ch)
}

func (clt *aggTradeBaseWs) run() {
	go clt.listener()
}

func (clt *aggTradeBaseWs) close() {
	clt.cancel()
	clt.rawWs.Close()
}

func (clt *aggTradeBaseWs) listener() {
	for {
		msg := clt.rawWs.Wait()
		if errors.Is(msg.Err, ErrWsClientClosed) {
			return
		}
		aggTradeMsg := AggTradeMsg{
			SymbolType: clt.sybType,
		}
		if msg.Err != nil {
			aggTradeMsg.Err = msg.Err
			clt.radio.BroadcastAll(aggTradeMsg)
			continue
		}
		var stream WsAggTradeStream
		if err := json.Unmarshal(msg.Data, &stream); err != nil {
			aggTradeMsg.Err = err
			clt.radio.BroadcastAll(aggTradeMsg)
			continue
		}
		aggTradeMsg.AggTrade = stream
		if stream.EventType == WsEventAggTrade {
			clt.radio.Broadcast(aggTradeStream(stream.Symbol), aggTradeMsg)
			continue
		}
		if stream.EventType != "" {
			// should not happen
			clt.logger.Error("unhandled ag trade ws msg", "msg", string(msg.Data))
			aggTradeMsg.Err = fmt.Errorf("unknown event type msg: %s", string(msg.Data))
			if stream.Symbol != "" {
				clt.radio.Broadcast(aggTradeStream(stream.Symbol), aggTradeMsg)
			} else {
				clt.radio.BroadcastAll(aggTradeMsg)
			}
			continue
		}
		var resp WsResp[any]
		if err := json.Unmarshal(msg.Data, &resp); err == nil && resp.Error != nil {
			aggTradeMsg.Err = fmt.Errorf("bnc: %d %s", resp.Error.Code, resp.Error.Msg)
			clt.radio.BroadcastAll(aggTradeMsg)
			continue
		}
		clt.logger.Warn("unhandled ag trade ws msg", "msg", string(msg.Data))
	}
}

func aggTradeStream(symbol string) string {
	return fmt.Sprintf("%s@aggTrade", strings.ToLower(symbol))
}
