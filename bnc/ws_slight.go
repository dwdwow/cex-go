package bnc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/dwdwow/props"
)

type SlightWsClientMsg struct {
	Event   WsEvent `json:"event"`
	Data    any     `json:"data"`
	IsArray bool    `json:"isArray"`
	Err     error   `json:"err"`
}

type SlightWsClientSubscriptionMsg[D any] struct {
	Data D     `json:"data"`
	Err  error `json:"err"`
}

type WsSlightClientSubscription[D any] struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	ws        *SlightWsClient
	event     string
	rawCh     <-chan SlightWsClientMsg
	ch        chan SlightWsClientSubscriptionMsg[D]
}

func newSlightWsClientSubscription[D any](ws *SlightWsClient, event string) *WsSlightClientSubscription[D] {
	ctx, ctxCancel := context.WithCancel(ws.ctx)
	suber := &WsSlightClientSubscription[D]{
		ctx:       ctx,
		ctxCancel: ctxCancel,
		ws:        ws,
		event:     event,
		rawCh:     ws.newSuberCh(event),
		ch:        make(chan SlightWsClientSubscriptionMsg[D], 1),
	}
	suber.start()
	return suber
}

func (w *WsSlightClientSubscription[D]) start() {
	go func() {
		for {
			var rawMsg SlightWsClientMsg
			select {
			case <-w.ctx.Done():
				return
			case rawMsg = <-w.rawCh:
			}
			var msg SlightWsClientSubscriptionMsg[D]
			if rawMsg.Err != nil {
				msg.Err = rawMsg.Err
			} else {
				var ok bool
				msg.Data, ok = rawMsg.Data.(D)
				if !ok {
					msg.Err = errors.New("invalid data type")
				}
			}
			w.ch <- msg
		}
	}()
}

func (w *WsSlightClientSubscription[D]) Close() {
	w.ctxCancel()
	w.ws.Unsub(w.event, w.rawCh)
}

func (w *WsSlightClientSubscription[D]) Chan() <-chan SlightWsClientSubscriptionMsg[D] {
	return w.ch
}

// SlightWsClient is a simple and lightweight websocket client for binance
type SlightWsClient struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	wsCfg       RawWsCfg
	rawWs       *RawWs
	unmarshaler WsDataUnmarshaler

	user *User

	radio *props.Radio[SlightWsClientMsg]

	logger *slog.Logger
}

func NewSlightWsClient(cfg RawWsCfg, unmarshaler WsDataUnmarshaler, logger *slog.Logger) *SlightWsClient {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	ctx, cancel := context.WithCancel(context.Background())
	var user *User
	if cfg.APIKey != "" && cfg.APISecretKey != "" {
		user = NewUser(UserCfg{
			APIKey:       cfg.APIKey,
			APISecretKey: cfg.APISecretKey,
		})
	}
	radio := props.NewRadio(
		props.WithFanoutDur[SlightWsClientMsg](time.Second),
		props.WithFanoutLogger[SlightWsClientMsg](logger),
		props.WithFanoutChCap[SlightWsClientMsg](1024),
	)
	ws := &SlightWsClient{
		ctx:         ctx,
		ctxCancel:   cancel,
		wsCfg:       cfg,
		unmarshaler: unmarshaler,
		user:        user,
		radio:       radio,
		logger:      logger,
	}
	ws.start()
	return ws
}

func (w *SlightWsClient) start() {
	rawWs, err := StartNewRawWs(w.wsCfg, w.logger)
	if err != nil {
		return
	}
	w.rawWs = rawWs
	go func() {
		for {
			msg := rawWs.Wait()
			if errors.Is(msg.Err, ErrWsClientClosed) {
				return
			}
			w.dataHandler(msg)
		}
	}()
}

func (w *SlightWsClient) Close() {
	w.ctxCancel()
	w.rawWs.Close()
}

const mfanKeyAll = "__all__"

func getWsEvent(data []byte) (event WsEvent, isArray, ok bool) {
	sd := string(data)
	isArray = strings.Index(sd, "[") == 0
	ss := strings.SplitSeq(sd, ",")
	for s := range ss {
		r := strings.Split(s, "\"e\":")
		if len(r) == 2 {
			return WsEvent(strings.ReplaceAll(r[1], "\"", "")), isArray, true
		}
	}
	return "", isArray, false
}

func (w *SlightWsClient) dataHandler(msg RawWsMsg) {
	if msg.Err != nil {
		w.sendToAll(SlightWsClientMsg{Data: msg.Data, Err: msg.Err})
		return
	}
	data := msg.Data
	e, isArray, ok := getWsEvent(data)
	if !ok {
		if string(data) == "{\"result\":null,\"id\":\"1\"}" ||
			string(data) == "{\"result\":null,\"id\":1}" {
			return
		}
		w.logger.Error("Can not get event", "data", string(data))
		return
	}
	// var specFans []*props.Fanout[WsClientMsg]
	// singleFan := w.mfan[string(e)]
	// allFan := w.mfan[mfanKeyAll]
	newMsg := SlightWsClientMsg{Event: e, Data: data, IsArray: isArray}
	var err error
	if w.unmarshaler != nil {
		newMsg.Data, err = w.unmarshaler(e, isArray, data)
		if err != nil {
			w.logger.Error("Can not unmarshal msg", "err", err, "data", string(data))
			return
		}
		w.specFans(isArray, newMsg)
	}
	w.radio.Broadcast(string(e), newMsg)
	w.radio.Broadcast(mfanKeyAll, newMsg)
	// if singleFan != nil {
	// 	singleFan.Broadcast(newMsg)
	// }
	// if allFan != nil {
	// 	allFan.Broadcast(newMsg)
	// }
	// for _, fan := range specFans {
	// 	if fan != nil {
	// 		fan.Broadcast(newMsg)
	// 	}
	// }
}

func (w *SlightWsClient) specFans(isArray bool, newMsg SlightWsClientMsg) (fans []*props.Fanout[SlightWsClientMsg]) {
	d := newMsg.Data
	switch newMsg.Event {
	case WsEventAggTrade:
		s, ok := d.(WsAggTradeStream)
		if !ok {
			return
		}
		w.radio.Broadcast(strings.ToLower(s.Symbol)+"@"+string(WsEventAggTrade), newMsg)
	case WsEventTrade:
		s, ok := d.(WsTradeStream)
		if !ok {
			return
		}
		w.radio.Broadcast(strings.ToLower(s.Symbol)+"@"+string(WsEventTrade), newMsg)
	case WsEventKline:
		s, ok := d.(WsKlineStream)
		if !ok {
			return
		}
		w.radio.Broadcast(strings.ToLower(s.Symbol)+"@"+string(WsEventKline)+"_"+string(s.Kline.Interval), newMsg)
		w.radio.Broadcast(strings.ToLower(s.Symbol)+"@"+string(WsEventKline)+"_", newMsg)
	case WsEventDepthUpdate:
		s, ok := d.(WsDepthStream)
		if !ok {
			return
		}
		w.radio.Broadcast(strings.ToLower(s.Symbol)+"@depth", newMsg)
		w.radio.Broadcast(strings.ToLower(s.Symbol)+"@depth@100ms", newMsg)
		w.radio.Broadcast(strings.ToLower(s.Symbol)+"@depth@250ms", newMsg)
		w.radio.Broadcast(strings.ToLower(s.Symbol)+"@depth@500ms", newMsg)
	case WsEventMarkPriceUpdate:
		if isArray {
			w.radio.Broadcast("!markPrice@arr@1s", newMsg)
			w.radio.Broadcast("!markPrice@arr", newMsg)
			return
		}
		s, ok := d.(WsMarkPriceStream)
		if !ok {
			return
		}
		w.radio.Broadcast(strings.ToLower(s.Symbol)+"@markPrice", newMsg)
		w.radio.Broadcast(strings.ToLower(s.Symbol)+"@markPrice@1s", newMsg)
	case WsEventIndexPriceUpdate:
		s, ok := d.(WsCMIndexPriceStream)
		if !ok {
			return
		}
		w.radio.Broadcast(strings.ToLower(s.Pair)+"@indexPrice", newMsg)
		w.radio.Broadcast(strings.ToLower(s.Pair)+"@indexPrice@1s", newMsg)
	case WsEventForceOrder:
		s, ok := d.(WsLiquidationOrderStream)
		if !ok {
			return
		}
		w.radio.Broadcast(strings.ToLower(s.Order.Symbol)+"@"+string(WsEventForceOrder), newMsg)
		w.radio.Broadcast("!forceOrder@arr", newMsg)
	}
	return
}

func (w *SlightWsClient) sendToAll(msg SlightWsClientMsg) {
	w.radio.Broadcast(mfanKeyAll, msg)
}

// newSuberCh
// Pass empty string if you want listen all events.
// Should SubStream firstly.
func (w *SlightWsClient) newSuberCh(event string) <-chan SlightWsClientMsg {
	return w.radio.Sub(event)
}

func (w *SlightWsClient) Unsub(event string, ch <-chan SlightWsClientMsg) {
	w.radio.Unsub(event, ch)
}

func (w *SlightWsClient) subStream(events ...string) (result RawWsSubStreamsResult, err error) {
	return w.rawWs.Sub(events...)
}

// SubAggTradeStream real time
func (w *SlightWsClient) SubAggTradeStream(symbols ...string) (result RawWsSubStreamsResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@aggTrade")
	}
	return w.subStream(params...)
}

// SubAggTrade real time
// if symbol is empty, will listen all aggTrade events
func (w *SlightWsClient) SubAggTrade(symbol string) *WsSlightClientSubscription[WsAggTradeStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@aggTrade"
	} else {
		event = string(WsEventAggTrade)
	}
	return newSlightWsClientSubscription[WsAggTradeStream](w, event)
}

// SubTradeStream real time
// just for spot ws
func (w *SlightWsClient) SubTradeStream(symbols ...string) (result RawWsSubStreamsResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@trade")
	}
	return w.subStream(params...)
}

// SubTrade real time
// if symbol is empty, will listen all trade events
func (w *SlightWsClient) SubTrade(symbol string) *WsSlightClientSubscription[WsTradeStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@trade"
	} else {
		event = string(WsEventTrade)
	}
	return newSlightWsClientSubscription[WsTradeStream](w, event)
}

// SubKlineStream 1000ms for 1s, 2000ms for others
// 1s just for spot kline
func (w *SlightWsClient) SubKlineStream(interval KlineInterval, symbols ...string) (result RawWsSubStreamsResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, fmt.Sprintf("%s@kline_%v", strings.ToLower(symbol), interval))
	}
	return w.subStream(params...)
}

// SubKline
// if symbol is empty, will listen all kline events
func (w *SlightWsClient) SubKline(symbol string, interval KlineInterval) *WsSlightClientSubscription[WsKlineStream] {
	var event string
	if symbol != "" {
		event = fmt.Sprintf("%s@kline_%v", strings.ToLower(symbol), interval)
	} else {
		event = string(WsEventKline)
	}
	return newSlightWsClientSubscription[WsKlineStream](w, event)
}

// SubDepthUpdateStream 1000ms for spot, 250ms for futures
func (w *SlightWsClient) SubDepthUpdateStream(symbols ...string) (result RawWsSubStreamsResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@depth")
	}
	return w.subStream(params...)
}

// SubDepthUpdateStream500ms 500ms
// just for futures ws
func (w *SlightWsClient) SubDepthUpdateStream500ms(symbols ...string) (result RawWsSubStreamsResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@depth@500ms")
	}
	return w.subStream(params...)
}

// SubDepthUpdateStream100ms 100ms
func (w *SlightWsClient) SubDepthUpdateStream100ms(symbols ...string) (result RawWsSubStreamsResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@depth@100ms")
	}
	return w.subStream(params...)
}

// SubDepthUpdate
// if symbol is empty, will listen all depthUpdate events
func (w *SlightWsClient) SubDepthUpdate(symbol string) *WsSlightClientSubscription[WsDepthStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@depth"
	} else {
		event = string(WsEventDepthUpdate)
	}
	return newSlightWsClientSubscription[WsDepthStream](w, event)
}

// SubDepthUpdate500ms
// if symbol is empty, will listen all depthUpdate 500ms events
func (w *SlightWsClient) SubDepthUpdate500ms(symbol string) *WsSlightClientSubscription[WsDepthStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@depth@500ms"
	} else {
		event = string(WsEventDepthUpdate)
	}
	return newSlightWsClientSubscription[WsDepthStream](w, event)
}

// SubDepthUpdate100ms
// if symbol is empty, will listen all depthUpdate 100ms events
func (w *SlightWsClient) SubDepthUpdate100ms(symbol string) *WsSlightClientSubscription[WsDepthStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@depth@100ms"
	} else {
		event = string(WsEventDepthUpdate)
	}
	return newSlightWsClientSubscription[WsDepthStream](w, event)
}

// SubMarkPriceStream1s 1s
func (w *SlightWsClient) SubMarkPriceStream1s(symbols ...string) (result RawWsSubStreamsResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@markPrice@1s")
	}
	return w.subStream(params...)
}

// SubMarkPriceStream3s 3s
func (w *SlightWsClient) SubMarkPriceStream3s(symbols ...string) (result RawWsSubStreamsResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@markPrice")
	}
	return w.subStream(params...)
}

func (w *SlightWsClient) SubMarkPrice1s(symbol string) *WsSlightClientSubscription[WsMarkPriceStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@markPrice@1s"
	} else {
		event = string(WsEventMarkPriceUpdate)
	}
	return newSlightWsClientSubscription[WsMarkPriceStream](w, event)
}

func (w *SlightWsClient) SubMarkPrice3s(symbol string) *WsSlightClientSubscription[WsMarkPriceStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@markPrice"
	} else {
		event = string(WsEventMarkPriceUpdate)
	}
	return newSlightWsClientSubscription[WsMarkPriceStream](w, event)
}

// SubAllMarkPriceStream1s 1s
// just for um futures
func (w *SlightWsClient) SubAllMarkPriceStream1s() (result RawWsSubStreamsResult, err error) {
	return w.subStream("!markPrice@arr@1s")
}

// SubAllMarkPriceStream3s 3s
// just for um futures
func (w *SlightWsClient) SubAllMarkPriceStream3s() (result RawWsSubStreamsResult, err error) {
	return w.subStream("!markPrice@arr")
}

func (w *SlightWsClient) SubAllMarkPrice1s() *WsSlightClientSubscription[[]WsMarkPriceStream] {
	event := "!markPrice@arr@1s"
	return newSlightWsClientSubscription[[]WsMarkPriceStream](w, event)
}

func (w *SlightWsClient) SubAllMarkPrice3s() *WsSlightClientSubscription[[]WsMarkPriceStream] {
	event := "!markPrice@arr"
	return newSlightWsClientSubscription[[]WsMarkPriceStream](w, event)
}

func (w *SlightWsClient) SubAllMarkPriceEvents() *WsSlightClientSubscription[[]WsMarkPriceStream] {
	event := string(WsEventMarkPriceUpdate)
	return newSlightWsClientSubscription[[]WsMarkPriceStream](w, event)
}

// SubCMIndexPriceStream3s 3s
// just for cm futures
func (w *SlightWsClient) SubCMIndexPriceStream3s(pairs ...string) (result RawWsSubStreamsResult, err error) {
	var params []string
	for _, pair := range pairs {
		params = append(params, strings.ToLower(pair)+"@indexPrice")
	}
	return w.subStream(params...)
}

// SubCMIndexPriceStream1s 1s
// just for cm futures
func (w *SlightWsClient) SubCMIndexPriceStream1s(pairs ...string) (result RawWsSubStreamsResult, err error) {
	var params []string
	for _, pair := range pairs {
		params = append(params, strings.ToLower(pair)+"@indexPrice@1s")
	}
	return w.subStream(params...)
}

// SubCMIndexPrice3s
// just for cm futures
// if pair is empty, will listen all WsEventIndexPriceUpdate events
func (w *SlightWsClient) SubCMIndexPrice3s(pair string) *WsSlightClientSubscription[WsCMIndexPriceStream] {
	var event string
	if pair != "" {
		event = strings.ToLower(pair) + "@indexPrice"
	} else {
		event = string(WsEventIndexPriceUpdate)
	}
	return newSlightWsClientSubscription[WsCMIndexPriceStream](w, event)
}

// SubCMIndexPrice1s
// just for cm futures
// if pair is empty, will listen all WsEventIndexPriceUpdate events
func (w *SlightWsClient) SubCMIndexPrice1s(pair string) *WsSlightClientSubscription[WsCMIndexPriceStream] {
	var event string
	if pair != "" {
		event = strings.ToLower(pair) + "@indexPrice@1s"
	} else {
		event = string(WsEventIndexPriceUpdate)
	}
	return newSlightWsClientSubscription[WsCMIndexPriceStream](w, event)
}

// SubLiquidationOrderStream 1s
// just for futures
func (w *SlightWsClient) SubLiquidationOrderStream(symbols ...string) (result RawWsSubStreamsResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@forceOrder")
	}
	return w.subStream(params...)
}

// SubLiquidationOrder 1s
// just for futures
// if symbol is empty, will listen all WsEventForceOrder events
func (w *SlightWsClient) SubLiquidationOrder(symbol string) *WsSlightClientSubscription[WsLiquidationOrderStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@forceOrder"
	} else {
		event = string(WsEventForceOrder)
	}
	return newSlightWsClientSubscription[WsLiquidationOrderStream](w, event)
}

// SubAllMarketLiquidationOrderStream 1s
func (w *SlightWsClient) SubAllMarketLiquidationOrderStream() (result RawWsSubStreamsResult, err error) {
	return w.subStream("!forceOrder@arr")
}

func (w *SlightWsClient) SubAllMarketLiquidationOrder() *WsSlightClientSubscription[WsLiquidationOrderStream] {
	event := "!forceOrder@arr"
	return newSlightWsClientSubscription[WsLiquidationOrderStream](w, event)
}
