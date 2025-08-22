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

func NewDepthUpdateStream(kit PublicStreamHandlerKit[WsDepthStream], logger *slog.Logger) *PublicStream[WsDepthStream] {
	return NewPublicStream(kit, logger)
}

func NewDepthUpdateWs(logger *slog.Logger) *PublicStream[WsDepthStream] {
	return NewDepthUpdateStream(DepthUpdateOneTopicWsKit{}, logger)
}

type DepthUpdateOneTopicWsKit struct {
}

func depthUpdate100msStream(symbol string) string {
	return fmt.Sprintf("%s@depth@100ms", strings.ToLower(symbol))
}

func (k DepthUpdateOneTopicWsKit) Stream(symbolType cex.SymbolType, symbol string) string {
	return depthUpdate100msStream(symbol)
}

func (k DepthUpdateOneTopicWsKit) New(symbolType cex.SymbolType, logger *slog.Logger) PublicStreamHandler[WsDepthStream] {
	return NewDepthUpdateSingleWs(symbolType, logger)
}

type depthUpdateBaseWs struct {
	ctx    context.Context
	cancel context.CancelFunc

	sybType cex.SymbolType

	producer StreamProducer

	radio *props.Radio[PublicStreamMsg[WsDepthStream]]

	logger *slog.Logger
}

func NewDepthUpdateHandler(sybType cex.SymbolType, producer StreamProducer, logger *slog.Logger) *depthUpdateBaseWs {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("bnc_dp_base_clt", sybType)
	ctx, cancel := context.WithCancel(context.Background())
	radio := props.NewRadio(
		props.WithFanoutLogger[PublicStreamMsg[WsDepthStream]](logger),
		props.WithFanoutChCap[PublicStreamMsg[WsDepthStream]](10000),
	)
	return &depthUpdateBaseWs{
		ctx:      ctx,
		cancel:   cancel,
		producer: producer,
		sybType:  sybType,
		radio:    radio,
		logger:   logger,
	}
}

func NewDepthUpdateSingleWs(sybType cex.SymbolType, logger *slog.Logger) *depthUpdateBaseWs {
	return NewDepthUpdateHandler(sybType, NewRawWs(DefaultPublicWsCfg(sybType), logger), logger)
}

func StartNewDepthUpdateSingleWs(sybType cex.SymbolType, logger *slog.Logger) (clt *depthUpdateBaseWs, err error) {
	clt = NewDepthUpdateSingleWs(sybType, logger)
	err = clt.Run()
	return
}

func (clt *depthUpdateBaseWs) createSubParams(symbols ...string) (params []string) {
	for _, symbol := range symbols {
		params = append(params, depthUpdate100msStream(symbol))
	}
	return params
}

func (clt *depthUpdateBaseWs) Streams() []string {
	return clt.producer.Streams()
}

func (clt *depthUpdateBaseWs) Sub(symbols ...string) (unsubed []string, err error) {
	params := clt.createSubParams(symbols...)
	res, err := clt.producer.Sub(params...)
	if err == nil {
		return
	}
	for _, p := range res.UnsubedStreams {
		for _, s := range symbols {
			if p == depthUpdate100msStream(s) {
				unsubed = append(unsubed, s)
				break
			}
		}
	}
	return unsubed, err
}

func (clt *depthUpdateBaseWs) Unsub(symbols ...string) (err error) {
	params := clt.createSubParams(symbols...)
	_, err = clt.producer.Unsub(params...)
	if err != nil {
		return err
	}
	for _, symbol := range symbols {
		clt.radio.Broadcast(depthUpdate100msStream(symbol), PublicStreamMsg[WsDepthStream]{
			SymbolType: clt.sybType,
			Stream:     WsDepthStream{Symbol: symbol},
			Err:        ErrWsStreamUnsubed,
		})
	}
	return nil
}

func (clt *depthUpdateBaseWs) NewCh(symbol string) <-chan PublicStreamMsg[WsDepthStream] {
	return clt.radio.Sub(depthUpdate100msStream(symbol))
}

func (clt *depthUpdateBaseWs) RemoveCh(ch <-chan PublicStreamMsg[WsDepthStream]) {
	clt.radio.UnsubAll(ch)
}

func (clt *depthUpdateBaseWs) Run() error {
	err := clt.producer.Run()
	if err != nil {
		return err
	}
	go clt.listener()
	return nil
}

func (clt *depthUpdateBaseWs) Close() {
	clt.cancel()
	clt.producer.Close()
}

func (clt *depthUpdateBaseWs) listener() {
	for {
		msg, err := clt.producer.Wait()
		if errors.Is(err, ErrWsClientClosed) {
			return
		}
		depthUpdateMsg := PublicStreamMsg[WsDepthStream]{
			SymbolType: clt.sybType,
		}
		if err != nil {
			depthUpdateMsg.Err = err
			clt.radio.BroadcastAll(depthUpdateMsg)
			continue
		}
		var stream WsDepthStream
		if err := json.Unmarshal(msg, &stream); err != nil {
			depthUpdateMsg.Err = err
			clt.radio.BroadcastAll(depthUpdateMsg)
			continue
		}
		depthUpdateMsg.Stream = stream
		if stream.EventType == WsEventDepthUpdate {
			clt.radio.Broadcast(depthUpdate100msStream(stream.Symbol), depthUpdateMsg)
			continue
		}
		if stream.EventType != "" {
			// should not happen
			clt.logger.Error("unhandled depth update ws msg", "msg", string(msg))
			depthUpdateMsg.Err = fmt.Errorf("unknown event type msg: %s", string(msg))
			if stream.Symbol != "" {
				clt.radio.Broadcast(depthUpdate100msStream(stream.Symbol), depthUpdateMsg)
			} else {
				clt.radio.BroadcastAll(depthUpdateMsg)
			}
			continue
		}
		var resp WsResp[any]
		if err := json.Unmarshal(msg, &resp); err == nil && resp.Error != nil {
			depthUpdateMsg.Err = fmt.Errorf("bnc: %d %s", resp.Error.Code, resp.Error.Msg)
			clt.radio.BroadcastAll(depthUpdateMsg)
			continue
		}
		clt.logger.Warn("unhandled depth update ws msg", "msg", string(msg))
	}
}
