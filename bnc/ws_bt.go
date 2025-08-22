package bnc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/props"
)

func NewBookTickerStream(kit PublicStreamHandlerKit[WsBookTickerStream], logger *slog.Logger) *PublicStream[WsBookTickerStream] {
	return NewPublicStream(kit, logger)
}

func NewBookTickerWs(logger *slog.Logger) *PublicStream[WsBookTickerStream] {
	return NewBookTickerStream(BookTickerOneTopicWsKit{}, logger)
}

type BookTickerOneTopicWsKit struct {
}

func bookTickerStream(symbol string) string {
	return fmt.Sprintf("%s@bookTicker", strings.ToLower(symbol))
}

func (k BookTickerOneTopicWsKit) Stream(symbolType cex.SymbolType, symbol string) string {
	return bookTickerStream(symbol)
}

func (k BookTickerOneTopicWsKit) New(symbolType cex.SymbolType, logger *slog.Logger) PublicStreamHandler[WsBookTickerStream] {
	return NewBookTickerSingleWs(symbolType, logger)
}

type bookTickerBaseWs struct {
	ctx    context.Context
	cancel context.CancelFunc

	sybType cex.SymbolType

	producer StreamProducer

	radio *props.Radio[PublicStreamMsg[WsBookTickerStream]]

	logger *slog.Logger
}

func NewBookTickerHandler(sybType cex.SymbolType, producer StreamProducer, logger *slog.Logger) *bookTickerBaseWs {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("bnc_bt_base_clt", sybType)
	ctx, cancel := context.WithCancel(context.Background())
	radio := props.NewRadio(
		props.WithFanoutLogger[PublicStreamMsg[WsBookTickerStream]](logger),
		props.WithFanoutChCap[PublicStreamMsg[WsBookTickerStream]](10000),
	)
	return &bookTickerBaseWs{
		ctx:      ctx,
		cancel:   cancel,
		producer: producer,
		sybType:  sybType,
		radio:    radio,
		logger:   logger,
	}
}

func NewBookTickerSingleWs(sybType cex.SymbolType, logger *slog.Logger) *bookTickerBaseWs {
	return NewBookTickerHandler(sybType, NewRawWs(DefaultPublicWsCfg(sybType), logger), logger)
}

func StartNewBookTickerSingleWs(sybType cex.SymbolType, logger *slog.Logger) (clt *bookTickerBaseWs, err error) {
	clt = NewBookTickerSingleWs(sybType, logger)
	err = clt.Run()
	return
}

func (clt *bookTickerBaseWs) createSubParams(symbols ...string) (params []string) {
	for _, symbol := range symbols {
		params = append(params, bookTickerStream(symbol))
	}
	return params
}

func (clt *bookTickerBaseWs) Streams() []string {
	return clt.producer.Streams()
}

func (clt *bookTickerBaseWs) Sub(symbols ...string) (unsubed []string, err error) {
	params := clt.createSubParams(symbols...)
	res, err := clt.producer.Sub(params...)
	if err == nil {
		return
	}
	for _, p := range res.UnsubedStreams {
		for _, s := range symbols {
			if p == bookTickerStream(s) {
				unsubed = append(unsubed, s)
				break
			}
		}
	}
	return unsubed, err
}

func (clt *bookTickerBaseWs) Unsub(symbols ...string) (err error) {
	params := clt.createSubParams(symbols...)
	_, err = clt.producer.Unsub(params...)
	if err != nil {
		return err
	}
	for _, symbol := range symbols {
		clt.radio.Broadcast(bookTickerStream(symbol), PublicStreamMsg[WsBookTickerStream]{
			SymbolType: clt.sybType,
			Stream:     WsBookTickerStream{Symbol: symbol},
			Err:        ErrWsStreamUnsubed,
		})
	}
	return nil
}

func (clt *bookTickerBaseWs) NewCh(symbol string) <-chan PublicStreamMsg[WsBookTickerStream] {
	return clt.radio.Sub(bookTickerStream(symbol))
}

func (clt *bookTickerBaseWs) RemoveCh(ch <-chan PublicStreamMsg[WsBookTickerStream]) {
	clt.radio.UnsubAll(ch)
}

func (clt *bookTickerBaseWs) Run() error {
	err := clt.producer.Run()
	if err != nil {
		return err
	}
	go clt.listener()
	return nil
}

func (clt *bookTickerBaseWs) Close() {
	clt.cancel()
	clt.producer.Close()
}

func (clt *bookTickerBaseWs) listener() {
	for {
		msg, err := clt.producer.Wait()
		if errors.Is(err, ErrWsClientClosed) {
			return
		}
		bookTickerMsg := PublicStreamMsg[WsBookTickerStream]{
			SymbolType: clt.sybType,
		}
		if err != nil {
			bookTickerMsg.Err = err
			clt.radio.BroadcastAll(bookTickerMsg)
			continue
		}
		var stream WsBookTickerStream
		if err := json.Unmarshal(msg, &stream); err != nil {
			bookTickerMsg.Err = err
			clt.radio.BroadcastAll(bookTickerMsg)
			continue
		}
		bookTickerMsg.Stream = stream
		if stream.EventType == WsEventBookTicker {
			clt.radio.Broadcast(bookTickerStream(stream.Symbol), bookTickerMsg)
			continue
		}
		// spot stream has no event type, event time and transaction time
		if stream.OrderBookUpdateId > 0 {
			bookTickerMsg.Stream.EventType = WsEventBookTicker
			bookTickerMsg.Stream.EventTime = time.Now().UnixMilli()
			clt.radio.Broadcast(bookTickerStream(stream.Symbol), bookTickerMsg)
			continue
		}
		if stream.EventType != "" {
			// should not happen
			clt.logger.Error("unhandled book ticker ws msg", "msg", string(msg))
			bookTickerMsg.Err = fmt.Errorf("unknown event type msg: %s", string(msg))
			if stream.Symbol != "" {
				clt.radio.Broadcast(bookTickerStream(stream.Symbol), bookTickerMsg)
			} else {
				clt.radio.BroadcastAll(bookTickerMsg)
			}
			continue
		}
		var resp WsResp[any]
		if err := json.Unmarshal(msg, &resp); err == nil && resp.Error != nil {
			bookTickerMsg.Err = fmt.Errorf("bnc: %d %s", resp.Error.Code, resp.Error.Msg)
			clt.radio.BroadcastAll(bookTickerMsg)
			continue
		}
		clt.logger.Warn("unhandled book ticker ws msg", "msg", string(msg))
	}
}
