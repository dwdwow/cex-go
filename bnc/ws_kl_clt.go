package bnc

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/props"
)

type KlineMsg struct {
	SymbolType cex.SymbolType
	Kline      WsKline
	Err        error
}

type klineBaseWsClient struct {
	ctx    context.Context
	cancel context.CancelFunc

	sybType cex.SymbolType

	ws       *RawWsClient
	rawMsgCh <-chan RawWsClientMsg

	radio *props.Radio[KlineMsg]

	logger *slog.Logger
}

func newKlineBaseWsClient(ctx context.Context, sybType cex.SymbolType, logger *slog.Logger) *klineBaseWsClient {
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
	clt := &klineBaseWsClient{
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

func (clt *klineBaseWsClient) createSubParams(interval KlineInterval, symbols ...string) (params []string) {
	for _, symbol := range symbols {
		params = append(params, fmt.Sprintf("%s@kline_%s", strings.ToLower(symbol), interval))
	}
	return params
}

func (clt *klineBaseWsClient) sub(interval KlineInterval, symbols ...string) (unsubed []string, err error) {
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

func (clt *klineBaseWsClient) unsub(interval KlineInterval, symbols ...string) (err error) {
	params := clt.createSubParams(interval, symbols...)
	err = clt.ws.UnsubStream(params...)
	if err != nil {
		return err
	}
	for _, symbol := range symbols {
		clt.radio.Broadcast(klineRadioChannel(symbol, interval), KlineMsg{
			SymbolType: clt.sybType,
			Kline:      WsKline{Symbol: symbol, Interval: interval},
			Err:        ErrWsStreamUnsubed,
		})
	}
	return nil
}

func (clt *klineBaseWsClient) newCh(symbol string, interval KlineInterval) <-chan KlineMsg {
	return clt.radio.Sub(klineRadioChannel(symbol, interval))
}

func (clt *klineBaseWsClient) removeCh(ch <-chan KlineMsg) {
	clt.radio.UnsubAll(ch)
}

func (clt *klineBaseWsClient) run() {
	go clt.listener()
}

func (clt *klineBaseWsClient) listener() {
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
			clt.radio.Broadcast(klineRadioChannel(kline.Symbol, kline.Interval), klineMsg)
		}
	}
}

func klineRadioChannel(symbol string, interval KlineInterval) string {
	return fmt.Sprintf("%s-rAdIo-%s", symbol, interval)
}
