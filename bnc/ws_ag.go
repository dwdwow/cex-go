package bnc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/props"
)

func NewAggTradeStream(kit PublicStreamHandlerKit[WsAggTradeStream], logger *slog.Logger) *PublicStream[WsAggTradeStream] {
	return NewPublicStream(kit, logger)
}

func NewAggTradeWs(logger *slog.Logger) *PublicStream[WsAggTradeStream] {
	return NewAggTradeStream(AggTradeOneTopicWsKit{}, logger)
}

type AggTradeOneTopicWsKit struct {
}

func aggTradeStreamer(symbol string) string {
	return fmt.Sprintf("%s@aggTrade", strings.ToLower(symbol))
}

func (k AggTradeOneTopicWsKit) Stream(symbolType cex.SymbolType, symbol string) string {
	return aggTradeStreamer(symbol)
}

func (k AggTradeOneTopicWsKit) New(symbolType cex.SymbolType, logger *slog.Logger) PublicStreamHandler[WsAggTradeStream] {
	return NewAggTradeSingleWs(symbolType, logger)
}

type AggTradeHandler struct {
	ctx    context.Context
	cancel context.CancelFunc

	sybType cex.SymbolType

	producer StreamProducer

	radio *props.Radio[PublicStreamMsg[WsAggTradeStream]]

	logger *slog.Logger
}

func NewAggTradeHandler(sybType cex.SymbolType, producer StreamProducer, logger *slog.Logger) *AggTradeHandler {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("bnc_ag_base_clt", sybType)
	ctx, cancel := context.WithCancel(context.Background())
	radio := props.NewRadio(
		props.WithFanoutLogger[PublicStreamMsg[WsAggTradeStream]](logger),
		props.WithFanoutChCap[PublicStreamMsg[WsAggTradeStream]](10000),
	)
	return &AggTradeHandler{
		ctx:      ctx,
		cancel:   cancel,
		producer: producer,
		sybType:  sybType,
		radio:    radio,
		logger:   logger,
	}
}

func NewAggTradeSingleWs(sybType cex.SymbolType, logger *slog.Logger) *AggTradeHandler {
	return NewAggTradeHandler(sybType, NewRawWs(DefaultPublicWsCfg(sybType), logger), logger)
}

func StartNewAggTradeSingleWs(sybType cex.SymbolType, logger *slog.Logger) (clt *AggTradeHandler, err error) {
	clt = NewAggTradeSingleWs(sybType, logger)
	err = clt.Run()
	return
}

func (clt *AggTradeHandler) createSubParams(symbols ...string) (params []string) {
	for _, symbol := range symbols {
		params = append(params, aggTradeStreamer(symbol))
	}
	return params
}

func (clt *AggTradeHandler) Streams() []string {
	return clt.producer.Streams()
}

func (clt *AggTradeHandler) Sub(symbols ...string) (unsubed []string, err error) {
	params := clt.createSubParams(symbols...)
	res, err := clt.producer.Sub(params...)
	if err == nil {
		return
	}
	for _, p := range res.UnsubedStreams {
		for _, s := range symbols {
			if p == aggTradeStreamer(s) {
				unsubed = append(unsubed, s)
				break
			}
		}
	}
	return unsubed, err
}

func (clt *AggTradeHandler) Unsub(symbols ...string) (err error) {
	params := clt.createSubParams(symbols...)
	_, err = clt.producer.Unsub(params...)
	if err != nil {
		return err
	}
	for _, symbol := range symbols {
		clt.radio.Broadcast(aggTradeStreamer(symbol), PublicStreamMsg[WsAggTradeStream]{
			SymbolType: clt.sybType,
			Stream:     WsAggTradeStream{Symbol: symbol},
			Err:        ErrWsStreamUnsubed,
		})
	}
	return nil
}

func (clt *AggTradeHandler) NewCh(symbol string) <-chan PublicStreamMsg[WsAggTradeStream] {
	return clt.radio.Sub(aggTradeStreamer(symbol))
}

func (clt *AggTradeHandler) RemoveCh(ch <-chan PublicStreamMsg[WsAggTradeStream]) {
	clt.radio.UnsubAll(ch)
}

func (clt *AggTradeHandler) Run() error {
	err := clt.producer.Run()
	if err != nil {
		return err
	}
	go clt.listener()
	return nil
}

func (clt *AggTradeHandler) Close() {
	clt.cancel()
	clt.producer.Close()
}

func (clt *AggTradeHandler) listener() {
	for {
		msg, err := clt.producer.Wait()
		if errors.Is(err, ErrWsClientClosed) {
			return
		}
		aggTradeMsg := PublicStreamMsg[WsAggTradeStream]{
			SymbolType: clt.sybType,
		}
		if err != nil {
			aggTradeMsg.Err = err
			clt.radio.BroadcastAll(aggTradeMsg)
			continue
		}
		var stream WsAggTradeStream
		if err := json.Unmarshal(msg, &stream); err != nil {
			aggTradeMsg.Err = err
			clt.radio.BroadcastAll(aggTradeMsg)
			continue
		}
		aggTradeMsg.Stream = stream
		if stream.EventType == WsEventAggTrade {
			clt.radio.Broadcast(aggTradeStreamer(stream.Symbol), aggTradeMsg)
			continue
		}
		if stream.EventType != "" {
			// should not happen
			clt.logger.Error("unhandled ag trade ws msg", "msg", string(msg))
			aggTradeMsg.Err = fmt.Errorf("unknown event type msg: %s", string(msg))
			if stream.Symbol != "" {
				clt.radio.Broadcast(aggTradeStreamer(stream.Symbol), aggTradeMsg)
			} else {
				clt.radio.BroadcastAll(aggTradeMsg)
			}
			continue
		}
		var resp WsResp[any]
		if err := json.Unmarshal(msg, &resp); err == nil && resp.Error != nil {
			aggTradeMsg.Err = fmt.Errorf("bnc: %d %s", resp.Error.Code, resp.Error.Msg)
			clt.radio.BroadcastAll(aggTradeMsg)
			continue
		}
		clt.logger.Warn("unhandled ag trade ws msg", "msg", string(msg))
	}
}
