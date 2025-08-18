package bnc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/dwdwow/props"
	"github.com/gorilla/websocket"
)

type RawWsClientMsg struct {
	Data []byte `json:"data"`
	Err  error  `json:"err"`
}

type WsClientStatus int

const (
	WS_CLIENT_STATUS_CLOSED WsClientStatus = iota - 1
	WS_CLIENT_STATUS_NEW
	WS_CLIENT_STATUS_STARTING
	WS_CLIENT_STATUS_STARTED
)

type RawWsClient struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	gmu       sync.Mutex

	cfg WsCfg

	user *User

	listenKey string

	conn *websocket.Conn

	fan *props.Fanout[RawWsClientMsg]

	stream []string

	status WsClientStatus

	muxReqToken   sync.Mutex
	crrTokenIndex int
	latestTokens  []int64
	muxReqId      sync.Mutex
	reqId         int64

	logger *slog.Logger
}

func NewRawWsClient(ctx context.Context, cfg WsCfg, logger *slog.Logger) *RawWsClient {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("bnc_raw_ws", cfg.Url)
	var user *User
	if cfg.APIKey != "" && cfg.APISecretKey != "" {
		user = NewUser(UserCfg{
			APIKey:       cfg.APIKey,
			APISecretKey: cfg.APISecretKey,
		})
	}
	fan := props.NewFanout(
		props.WithFanoutDur[RawWsClientMsg](time.Second),
		props.WithFanoutLogger[RawWsClientMsg](logger),
		props.WithFanoutChCap[RawWsClientMsg](cfg.ChCap),
	)
	return &RawWsClient{
		ctx:          ctx,
		ctxCancel:    cancel,
		cfg:          cfg,
		user:         user,
		fan:          fan,
		latestTokens: make([]int64, cfg.MaxReqPerDur),
		reqId:        1000,
		logger:       logger,
	}
}

func (w *RawWsClient) Status() WsClientStatus {
	w.gmu.Lock()
	defer w.gmu.Unlock()
	return w.status
}

// Close will close the client and cancel the context
// client can not be restarted after close
func (w *RawWsClient) Close() {
	w.gmu.Lock()
	defer w.gmu.Unlock()
	defer func() {
		w.status = WS_CLIENT_STATUS_CLOSED
	}()
	if w.ctxCancel != nil {
		w.ctxCancel()
	}
	if w.conn != nil {
		_ = w.conn.Close()
	}
	w.conn = nil
}

func (w *RawWsClient) startWithLock() error {
	w.gmu.Lock()
	defer w.gmu.Unlock()
	return w.start()
}

func (w *RawWsClient) start() error {
	switch w.status {
	case WS_CLIENT_STATUS_NEW:
	case WS_CLIENT_STATUS_STARTED:
		return nil
	case WS_CLIENT_STATUS_CLOSED:
		return fmt.Errorf("bnc: ws client is closed, cannot restart")
	default:
		// should not happen
		return fmt.Errorf("bnc: ws client status is %d, cannot start", w.status)
	}

	w.status = WS_CLIENT_STATUS_STARTING

	defer func() {
		if w.conn == nil {
			w.status = WS_CLIENT_STATUS_NEW
		}
	}()

	w.logger.Info("Starting")
	if w.ctxCancel != nil {
		w.ctxCancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.ctx = ctx
	w.ctxCancel = cancel

	if w.conn != nil {
		_ = w.conn.Close()
		w.conn = nil
	}

	lk, err := w.newAndKeepListenKey()

	if err != nil {
		return err
	}

	dialer := websocket.Dialer{}
	var path string
	if lk != "" {
		path = "/" + lk
	}
	conn, resp, err := dialer.DialContext(w.ctx, w.cfg.Url+path, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 101 {
		return fmt.Errorf("bnc_ws: response status code %d", resp.StatusCode)
	}

	w.conn = conn
	w.status = WS_CLIENT_STATUS_STARTED

	stream := w.stream
	w.stream = nil
	if len(stream) > 0 {
		_, err = w.subStreamNoLock(stream...)
		if err != nil {
			w.logger.Error("Cannot Sub Stream", "err", err)
			w.fan.Broadcast(RawWsClientMsg{nil, err})
			w.conn.Close()
			w.conn = nil
			w.ctxCancel()
			return err
		}
	}

	go w.connListener(conn)
	go w.ctxWaiter()

	return nil
}

func (w *RawWsClient) ctxWaiter() {
	<-w.ctx.Done()
	w.Close()
}

func (w *RawWsClient) connListener(conn *websocket.Conn) {
	if conn == nil {
		w.logger.Error("bnc: nil ws conn")
		return
	}

	w.logger.Info("Conn listener started")

	for {
		t, d, err := w.read()
		if err != nil {
			w.logger.Error("Read message", "err", err)
			w.fan.Broadcast(RawWsClientMsg{nil, err})
			break
		}
		switch t {
		case websocket.PingMessage:
			w.logger.Info("Server ping received", "msg", string(d))
			err = w.write(d)
			if err != nil {
				w.logger.Error("Write pong msg", "err", err)
			}
		case websocket.PongMessage:
			w.logger.Info("Server pong received", "msg", string(d))
		case websocket.TextMessage:
			w.fan.Broadcast(RawWsClientMsg{Data: d})
		case websocket.BinaryMessage:
			w.logger.Info("Server binary received", "msg", string(d), "binary", d)
		case websocket.CloseMessage:
			w.logger.Info("Server closed message", "msg", string(d))
		}
	}

	for {
		w.gmu.Lock()
		if w.conn == nil {
			_ = w.conn.Close()
			w.conn = nil
		}
		if w.status == WS_CLIENT_STATUS_CLOSED {
			w.gmu.Unlock()
			return
		}
		w.status = WS_CLIENT_STATUS_NEW
		w.gmu.Unlock()
		err := w.startWithLock()
		if err != nil {
			w.logger.Error("Cannot Restart", "err", err)
			w.fan.Broadcast(RawWsClientMsg{nil, err})
			time.Sleep(time.Second)
			continue
		}
		break
	}
}

func (w *RawWsClient) write(data []byte) error {
	// cannot read and write concurrently
	if w.conn == nil {
		return errors.New("bnc: nil conn")
	}
	if !w.canWriteMsg() {
		return errors.New("bnc: too frequent write")
	}
	return w.conn.WriteMessage(websocket.TextMessage, data)
}

func (w *RawWsClient) read() (msgType int, data []byte, err error) {
	// cannot read and write concurrently
	if w.conn == nil {
		return -1, nil, errors.New("bnc: nil conn")
	}
	return w.conn.ReadMessage()
}

func (w *RawWsClient) Sub() <-chan RawWsClientMsg {
	return w.fan.Sub()
}

func (w *RawWsClient) Unsub(c <-chan RawWsClientMsg) {
	w.fan.Unsub(c)
}

func rawWsNewStreamFilter(oldStream, params []string, maxStream int) (existedStream, newStream, remainingParams []string) {
	for _, s := range params {
		if slices.Contains(oldStream, s) {
			existedStream = append(existedStream, s)
		} else {
			newStream = append(newStream, s)
		}
	}
	if maxStream < 0 {
		return
	}
	extraNum := len(oldStream) + len(newStream) - maxStream
	if extraNum > 0 {
		ns := newStream
		i := len(ns) - extraNum
		newStream = ns[:i]
		remainingParams = ns[i:]
	}
	return
}

type RawWsSubStreamResult struct {
	ExistedParams  []string
	NewSubedParams []string
	UnsubedParams  []string
}

var ErrNotAllStreamSubed = errors.New("bnc: not all stream subed")

func (w *RawWsClient) subStreamNoLock(params ...string) (result RawWsSubStreamResult, err error) {
	if w.status == WS_CLIENT_STATUS_NEW {
		err = w.start()
		if err != nil {
			return
		}
	}
	existed, newSubed, unsubed := rawWsNewStreamFilter(w.stream, params, w.cfg.MaxStream)
	result.ExistedParams = existed
	result.UnsubedParams = append(result.UnsubedParams, newSubed...)
	result.UnsubedParams = append(result.UnsubedParams, unsubed...)
	if len(newSubed) == 0 {
		return
	}
	var data []byte
	if strings.Contains(w.cfg.Url, "dstream.binance.com") {
		data, err = json.Marshal(WsSubMsgInt64Id{
			Method: WsMethodSub,
			Params: newSubed,
			Id:     1,
		})
	} else {
		data, err = json.Marshal(WsSubMsg{
			Method: WsMethodSub,
			Params: newSubed,
			Id:     "1",
		})
	}
	if err != nil {
		return
	}
	err = w.write(data)
	if err != nil {
		return
	}
	w.stream = append(w.stream, newSubed...)
	result.NewSubedParams = newSubed
	result.UnsubedParams = unsubed
	if len(result.UnsubedParams) > 0 {
		err = ErrNotAllStreamSubed
	}
	return
}

// SubStream will return extra params that cannot be subscribed
// because of the max stream limit
func (w *RawWsClient) SubStream(params ...string) (result RawWsSubStreamResult, err error) {
	w.gmu.Lock()
	defer w.gmu.Unlock()
	return w.subStreamNoLock(params...)
}

func (w *RawWsClient) unsubStream(params ...string) error {
	var err error
	if strings.Contains(w.cfg.Url, "dstream.binance.com") {
		err = w.conn.WriteJSON(WsSubMsgInt64Id{
			Method: WsMethodUnsub,
			Params: params,
			Id:     1,
		})
	} else {
		err = w.conn.WriteJSON(WsSubMsg{
			Method: WsMethodUnsub,
			Params: params,
			Id:     "1",
		})
	}
	if err != nil {
		return err
	}
	for _, s := range params {
		i := slices.Index(w.stream, s)
		if i > -1 {
			w.stream = slices.Delete(w.stream, i, i+1)
		}
	}
	return nil
}

func (w *RawWsClient) UnsubStream(params ...string) error {
	w.gmu.Lock()
	defer w.gmu.Unlock()
	return w.unsubStream(params...)
}

func (w *RawWsClient) newAndKeepListenKey() (string, error) {
	if w.user == nil || w.cfg.ListenKeyUrl == "" {
		return "", nil
	}
	w.logger.Info("Getting new listen key")
	lk, err := w.newListenKey()
	if err != nil {
		return "", err
	}
	w.logger.Info("Listen key gotten")
	w.listenKey = lk
	go w.listenKeyKeeper(w.ctx)
	return lk, nil
}

func (w *RawWsClient) newListenKey() (string, error) {
	return w.user.NewListenKey(w.cfg.ListenKeyUrl)
}

func (w *RawWsClient) listenKeyKeeper(ctx context.Context) {
	w.logger.Info("Listen key keeper started")
	for {
		select {
		case <-ctx.Done():
		case <-time.After(time.Minute * 20):
			w.logger.Info("Keep listening key")
			_, err := w.user.KeepListenKey(w.cfg.ListenKeyUrl, w.listenKey)
			if err != nil {
				w.logger.Error("Cannot Keep Listen Key", "err", err)
			}
		}
	}
}

func (w *RawWsClient) canWriteMsg() bool {
	w.muxReqToken.Lock()
	defer w.muxReqToken.Unlock()
	t := time.Now().UnixMilli()
	withinDurNum := 0
	for _, v := range w.latestTokens {
		if t-v < w.cfg.ReqDur.Milliseconds() {
			withinDurNum++
		}
	}
	maxTokenNum := len(w.latestTokens)
	if withinDurNum >= maxTokenNum {
		return false
	}
	i := w.crrTokenIndex + 1
	if i >= maxTokenNum {
		i -= maxTokenNum
	}
	w.latestTokens[i] = t
	w.crrTokenIndex = i
	return true
}

func (w *RawWsClient) Request(method WsMethod, params []any) (id string, err error) {
	w.muxReqId.Lock()
	w.reqId++
	id = fmt.Sprintf("%d", w.reqId)
	w.muxReqId.Unlock()
	d, err := json.Marshal(WsReqMsg{
		Method: method,
		Params: params,
		Id:     id,
	})
	if err != nil {
		return
	}
	err = w.write(d)
	return
}

type MergedWsRawClient struct {
	ctx    context.Context
	mu     sync.Mutex
	cfg    WsCfg
	fan    *props.Fanout[RawWsClientMsg]
	clts   []*RawWsClient
	stream []string
	logger *slog.Logger
}

func NewMergedWsRawClient(ctx context.Context, cfg WsCfg, logger *slog.Logger) *MergedWsRawClient {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("ws", "bnc_merged_ws_raw")
	return &MergedWsRawClient{
		ctx:    ctx,
		fan:    props.NewFanout(props.WithFanoutDur[RawWsClientMsg](time.Second)),
		cfg:    cfg,
		logger: logger,
	}
}

func (m *MergedWsRawClient) SubStream(params ...string) (result RawWsSubStreamResult, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	existed, unsubed, unsubedParams := rawWsNewStreamFilter(m.stream, params, -1)
	result.ExistedParams = existed
	result.UnsubedParams = unsubedParams

	defer func() {
		var allStream []string
		for _, clt := range m.clts {
			allStream = append(allStream, clt.stream...)
		}
		m.stream = allStream
		result.UnsubedParams = append(result.UnsubedParams, unsubed...)
		if len(result.UnsubedParams) > 0 {
			err = ErrNotAllStreamSubed
		}
	}()

	for _, clt := range m.clts {
		if len(unsubed) == 0 {
			return
		}
		res, err := clt.SubStream(unsubed...)
		if err != nil {
			return result, err
		}
		unsubed = res.UnsubedParams
		result.NewSubedParams = append(result.NewSubedParams, res.NewSubedParams...)
	}
	for {
		if len(unsubed) == 0 {
			return
		}
		clt := NewRawWsClient(m.ctx, m.cfg, m.logger)
		clt.fan = m.fan
		res, err := clt.SubStream(unsubed...)
		if err != nil {
			return result, err
		}
		m.clts = append(m.clts, clt)
		result.NewSubedParams = append(result.NewSubedParams, res.NewSubedParams...)
		unsubed = res.UnsubedParams
	}
}

func (m *MergedWsRawClient) UnsubStream(params ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	defer func() {
		var allStream []string
		for _, clt := range m.clts {
			allStream = append(allStream, clt.stream...)
		}
		m.stream = allStream
	}()
	for _, clt := range m.clts {
		var unsubStream []string
		for _, p := range params {
			if slices.Contains(clt.stream, p) {
				unsubStream = append(unsubStream, p)
			}
		}
		err := clt.UnsubStream(unsubStream...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MergedWsRawClient) NewSuber() <-chan RawWsClientMsg {
	return m.fan.Sub()
}

type WsClientMsg struct {
	Event   WsEvent `json:"event"`
	Data    any     `json:"data"`
	IsArray bool    `json:"isArray"`
	Err     error   `json:"err"`
}

type WsClientSubscriptionMsg[D any] struct {
	Data D     `json:"data"`
	Err  error `json:"err"`
}

type WsClientSubscription[D any] struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	ws        *WsClient
	event     string
	rawCh     <-chan WsClientMsg
	ch        chan WsClientSubscriptionMsg[D]
}

func newWsClientSubscription[D any](ws *WsClient, event string) *WsClientSubscription[D] {
	ctx, ctxCancel := context.WithCancel(ws.ctx)
	suber := &WsClientSubscription[D]{
		ctx:       ctx,
		ctxCancel: ctxCancel,
		ws:        ws,
		event:     event,
		rawCh:     ws.newSuberCh(event),
		ch:        make(chan WsClientSubscriptionMsg[D], 1),
	}
	suber.start()
	return suber
}

func (w *WsClientSubscription[D]) start() {
	go func() {
		for {
			var rawMsg WsClientMsg
			select {
			case <-w.ctx.Done():
				return
			case rawMsg = <-w.rawCh:
			}
			var msg WsClientSubscriptionMsg[D]
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

func (w *WsClientSubscription[D]) Close() {
	w.ctxCancel()
	w.ws.Unsub(w.event, w.rawCh)
}

func (w *WsClientSubscription[D]) Chan() <-chan WsClientSubscriptionMsg[D] {
	return w.ch
}

// WsClient is a simple and lightweight websocket client for binance
type WsClient struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	wsCfg WsCfg
	rawWs *MergedWsRawClient

	user *User

	radio *props.Radio[WsClientMsg]

	logger *slog.Logger
}

func NewWsClient(ctx context.Context, cfg WsCfg, logger *slog.Logger) *WsClient {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	var user *User
	if cfg.APIKey != "" && cfg.APISecretKey != "" {
		user = NewUser(UserCfg{
			APIKey:       cfg.APIKey,
			APISecretKey: cfg.APISecretKey,
		})
	}
	radio := props.NewRadio(
		props.WithFanoutDur[WsClientMsg](time.Second),
		props.WithFanoutLogger[WsClientMsg](logger),
		props.WithFanoutChCap[WsClientMsg](cfg.ChCap),
	)
	ws := &WsClient{
		ctx:       ctx,
		ctxCancel: cancel,
		wsCfg:     cfg,
		user:      user,
		radio:     radio,
		logger:    logger,
	}
	ws.start()
	return ws
}

func (w *WsClient) start() {
	rawWs := NewMergedWsRawClient(w.ctx, w.wsCfg, w.logger)
	w.rawWs = rawWs
	ch := rawWs.NewSuber()
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				return
			case msg := <-ch:
				w.dataHandler(msg)
			}
		}
	}()
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

func (w *WsClient) dataHandler(msg RawWsClientMsg) {
	if msg.Err != nil {
		w.sendToAll(WsClientMsg{Data: msg.Data, Err: msg.Err})
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
	newMsg := WsClientMsg{Event: e, Data: data, IsArray: isArray}
	var err error
	if w.rawWs.cfg.DataUnmarshaler != nil {
		newMsg.Data, err = w.rawWs.cfg.DataUnmarshaler(e, isArray, data)
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

func (w *WsClient) specFans(isArray bool, newMsg WsClientMsg) (fans []*props.Fanout[WsClientMsg]) {
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

func (w *WsClient) sendToAll(msg WsClientMsg) {
	w.radio.Broadcast(mfanKeyAll, msg)
}

// newSuberCh
// Pass empty string if you want listen all events.
// Should SubStream firstly.
func (w *WsClient) newSuberCh(event string) <-chan WsClientMsg {
	return w.radio.Sub(event)
}

func (w *WsClient) Unsub(event string, ch <-chan WsClientMsg) {
	w.radio.Unsub(event, ch)
}

func (w *WsClient) subStream(events ...string) (result RawWsSubStreamResult, err error) {
	return w.rawWs.SubStream(events...)
}

// SubAggTradeStream real time
func (w *WsClient) SubAggTradeStream(symbols ...string) (result RawWsSubStreamResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@aggTrade")
	}
	return w.subStream(params...)
}

// SubAggTrade real time
// if symbol is empty, will listen all aggTrade events
func (w *WsClient) SubAggTrade(symbol string) *WsClientSubscription[WsAggTradeStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@aggTrade"
	} else {
		event = string(WsEventAggTrade)
	}
	return newWsClientSubscription[WsAggTradeStream](w, event)
}

// SubTradeStream real time
// just for spot ws
func (w *WsClient) SubTradeStream(symbols ...string) (result RawWsSubStreamResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@trade")
	}
	return w.subStream(params...)
}

// SubTrade real time
// if symbol is empty, will listen all trade events
func (w *WsClient) SubTrade(symbol string) *WsClientSubscription[WsTradeStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@trade"
	} else {
		event = string(WsEventTrade)
	}
	return newWsClientSubscription[WsTradeStream](w, event)
}

// SubKlineStream 1000ms for 1s, 2000ms for others
// 1s just for spot kline
func (w *WsClient) SubKlineStream(interval KlineInterval, symbols ...string) (result RawWsSubStreamResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, fmt.Sprintf("%s@kline_%v", strings.ToLower(symbol), interval))
	}
	return w.subStream(params...)
}

// SubKline
// if symbol is empty, will listen all kline events
func (w *WsClient) SubKline(symbol string, interval KlineInterval) *WsClientSubscription[WsKlineStream] {
	var event string
	if symbol != "" {
		event = fmt.Sprintf("%s@kline_%v", strings.ToLower(symbol), interval)
	} else {
		event = string(WsEventKline)
	}
	return newWsClientSubscription[WsKlineStream](w, event)
}

// SubDepthUpdateStream 1000ms for spot, 250ms for futures
func (w *WsClient) SubDepthUpdateStream(symbols ...string) (result RawWsSubStreamResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@depth")
	}
	return w.subStream(params...)
}

// SubDepthUpdateStream500ms 500ms
// just for futures ws
func (w *WsClient) SubDepthUpdateStream500ms(symbols ...string) (result RawWsSubStreamResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@depth@500ms")
	}
	return w.subStream(params...)
}

// SubDepthUpdateStream100ms 100ms
func (w *WsClient) SubDepthUpdateStream100ms(symbols ...string) (result RawWsSubStreamResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@depth@100ms")
	}
	return w.subStream(params...)
}

// SubDepthUpdate
// if symbol is empty, will listen all depthUpdate events
func (w *WsClient) SubDepthUpdate(symbol string) *WsClientSubscription[WsDepthStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@depth"
	} else {
		event = string(WsEventDepthUpdate)
	}
	return newWsClientSubscription[WsDepthStream](w, event)
}

// SubDepthUpdate500ms
// if symbol is empty, will listen all depthUpdate 500ms events
func (w *WsClient) SubDepthUpdate500ms(symbol string) *WsClientSubscription[WsDepthStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@depth@500ms"
	} else {
		event = string(WsEventDepthUpdate)
	}
	return newWsClientSubscription[WsDepthStream](w, event)
}

// SubDepthUpdate100ms
// if symbol is empty, will listen all depthUpdate 100ms events
func (w *WsClient) SubDepthUpdate100ms(symbol string) *WsClientSubscription[WsDepthStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@depth@100ms"
	} else {
		event = string(WsEventDepthUpdate)
	}
	return newWsClientSubscription[WsDepthStream](w, event)
}

// SubMarkPriceStream1s 1s
func (w *WsClient) SubMarkPriceStream1s(symbols ...string) (result RawWsSubStreamResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@markPrice@1s")
	}
	return w.subStream(params...)
}

// SubMarkPriceStream3s 3s
func (w *WsClient) SubMarkPriceStream3s(symbols ...string) (result RawWsSubStreamResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@markPrice")
	}
	return w.subStream(params...)
}

func (w *WsClient) SubMarkPrice1s(symbol string) *WsClientSubscription[WsMarkPriceStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@markPrice@1s"
	} else {
		event = string(WsEventMarkPriceUpdate)
	}
	return newWsClientSubscription[WsMarkPriceStream](w, event)
}

func (w *WsClient) SubMarkPrice3s(symbol string) *WsClientSubscription[WsMarkPriceStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@markPrice"
	} else {
		event = string(WsEventMarkPriceUpdate)
	}
	return newWsClientSubscription[WsMarkPriceStream](w, event)
}

// SubAllMarkPriceStream1s 1s
// just for um futures
func (w *WsClient) SubAllMarkPriceStream1s() (result RawWsSubStreamResult, err error) {
	return w.subStream("!markPrice@arr@1s")
}

// SubAllMarkPriceStream3s 3s
// just for um futures
func (w *WsClient) SubAllMarkPriceStream3s() (result RawWsSubStreamResult, err error) {
	return w.subStream("!markPrice@arr")
}

func (w *WsClient) SubAllMarkPrice1s() *WsClientSubscription[[]WsMarkPriceStream] {
	event := "!markPrice@arr@1s"
	return newWsClientSubscription[[]WsMarkPriceStream](w, event)
}

func (w *WsClient) SubAllMarkPrice3s() *WsClientSubscription[[]WsMarkPriceStream] {
	event := "!markPrice@arr"
	return newWsClientSubscription[[]WsMarkPriceStream](w, event)
}

func (w *WsClient) SubAllMarkPriceEvents() *WsClientSubscription[[]WsMarkPriceStream] {
	event := string(WsEventMarkPriceUpdate)
	return newWsClientSubscription[[]WsMarkPriceStream](w, event)
}

// SubCMIndexPriceStream3s 3s
// just for cm futures
func (w *WsClient) SubCMIndexPriceStream3s(pairs ...string) (result RawWsSubStreamResult, err error) {
	var params []string
	for _, pair := range pairs {
		params = append(params, strings.ToLower(pair)+"@indexPrice")
	}
	return w.subStream(params...)
}

// SubCMIndexPriceStream1s 1s
// just for cm futures
func (w *WsClient) SubCMIndexPriceStream1s(pairs ...string) (result RawWsSubStreamResult, err error) {
	var params []string
	for _, pair := range pairs {
		params = append(params, strings.ToLower(pair)+"@indexPrice@1s")
	}
	return w.subStream(params...)
}

// SubCMIndexPrice3s
// just for cm futures
// if pair is empty, will listen all WsEventIndexPriceUpdate events
func (w *WsClient) SubCMIndexPrice3s(pair string) *WsClientSubscription[WsCMIndexPriceStream] {
	var event string
	if pair != "" {
		event = strings.ToLower(pair) + "@indexPrice"
	} else {
		event = string(WsEventIndexPriceUpdate)
	}
	return newWsClientSubscription[WsCMIndexPriceStream](w, event)
}

// SubCMIndexPrice1s
// just for cm futures
// if pair is empty, will listen all WsEventIndexPriceUpdate events
func (w *WsClient) SubCMIndexPrice1s(pair string) *WsClientSubscription[WsCMIndexPriceStream] {
	var event string
	if pair != "" {
		event = strings.ToLower(pair) + "@indexPrice@1s"
	} else {
		event = string(WsEventIndexPriceUpdate)
	}
	return newWsClientSubscription[WsCMIndexPriceStream](w, event)
}

// SubLiquidationOrderStream 1s
// just for futures
func (w *WsClient) SubLiquidationOrderStream(symbols ...string) (result RawWsSubStreamResult, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@forceOrder")
	}
	return w.subStream(params...)
}

// SubLiquidationOrder 1s
// just for futures
// if symbol is empty, will listen all WsEventForceOrder events
func (w *WsClient) SubLiquidationOrder(symbol string) *WsClientSubscription[WsLiquidationOrderStream] {
	var event string
	if symbol != "" {
		event = strings.ToLower(symbol) + "@forceOrder"
	} else {
		event = string(WsEventForceOrder)
	}
	return newWsClientSubscription[WsLiquidationOrderStream](w, event)
}

// SubAllMarketLiquidationOrderStream 1s
func (w *WsClient) SubAllMarketLiquidationOrderStream() (result RawWsSubStreamResult, err error) {
	return w.subStream("!forceOrder@arr")
}

func (w *WsClient) SubAllMarketLiquidationOrder() *WsClientSubscription[WsLiquidationOrderStream] {
	event := "!forceOrder@arr"
	return newWsClientSubscription[WsLiquidationOrderStream](w, event)
}
