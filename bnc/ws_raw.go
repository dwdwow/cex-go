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
		props.WithFanoutDur[RawWsClientMsg](cfg.FanoutTimerDur),
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
