package bnc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/dwdwow/limiter-go"
	"github.com/dwdwow/props"
	"github.com/gorilla/websocket"
)

type StreamProducer interface {
	Streams() []string
	Run() error
	Sub(streams ...string) (result StreamsSubResult, err error)
	Unsub(streams ...string) (id int64, err error)
	Wait() ([]byte, error)
	Close()
}

type RawWsStatus int

const (
	RAW_WS_STATUS_CLOSED RawWsStatus = iota - 1
	RAW_WS_STATUS_NEW
	RAW_WS_STATUS_STARTED
	RAW_WS_STATUS_STARTING
)

type RawWs struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	gmu       sync.Mutex

	cfg RawWsCfg

	user *User

	listenKey string

	conn *websocket.Conn

	streams []string

	status RawWsStatus

	reqLimiter *limiter.Limiter

	// raw ws dose not use this channel
	// just for high level ws client to wait response
	respByReqId *props.SafeRWMap[int64, chan WsResp[any]]

	logger *slog.Logger
}

func NewRawWs(cfg RawWsCfg, logger *slog.Logger) *RawWs {
	ctx, cancel := context.WithCancel(context.Background())
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
	return &RawWs{
		ctx:         ctx,
		ctxCancel:   cancel,
		cfg:         cfg,
		user:        user,
		reqLimiter:  limiter.New(time.Second, cfg.MaxReqPerSecond),
		respByReqId: props.NewSafeRWMap[int64, chan WsResp[any]](),
		logger:      logger,
	}
}

func StartNewRawWs(cfg RawWsCfg, logger *slog.Logger) (ws *RawWs, err error) {
	ctx, cancel := context.WithCancel(context.Background())
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
	ws = &RawWs{
		ctx:         ctx,
		ctxCancel:   cancel,
		cfg:         cfg,
		user:        user,
		reqLimiter:  limiter.New(time.Second, cfg.MaxReqPerSecond),
		respByReqId: props.NewSafeRWMap[int64, chan WsResp[any]](),
		logger:      logger,
	}
	err = ws.start()
	return
}

func (w *RawWs) Status() RawWsStatus {
	w.gmu.Lock()
	defer w.gmu.Unlock()
	return w.status
}

func (w *RawWs) Streams() []string {
	w.gmu.Lock()
	defer w.gmu.Unlock()
	return w.streams
}

// Close will close the client and cancel the context
// client can not be restarted after close
func (w *RawWs) Close() {
	w.gmu.Lock()
	defer w.gmu.Unlock()
	w.status = RAW_WS_STATUS_CLOSED
	if w.ctxCancel != nil {
		w.ctxCancel()
	}
	if w.conn != nil {
		_ = w.conn.Close()
	}
}

func (w *RawWs) Run() error {
	return w.start()
}

func (w *RawWs) start() error {
	w.gmu.Lock()
	defer w.gmu.Unlock()
	return w.startNoLock()
}

func (w *RawWs) startNoLock() error {
	switch w.status {
	case RAW_WS_STATUS_STARTED:
		return nil
	case RAW_WS_STATUS_CLOSED:
		return ErrWsClientClosed
	case RAW_WS_STATUS_NEW:
	default:
		// should not happen
		return fmt.Errorf("bnc: ws client status is %d, cannot start", w.status)
	}

	w.status = RAW_WS_STATUS_STARTING

	defer func() {
		if w.status == RAW_WS_STATUS_STARTING {
			w.status = RAW_WS_STATUS_NEW
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
	if w.cfg.Microsecond {
		path += "?timeUnit=microsecond"
	}
	conn, resp, err := dialer.DialContext(w.ctx, w.cfg.Url+path, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 101 {
		return fmt.Errorf("bnc_ws: response status code %d", resp.StatusCode)
	}

	w.conn = conn

	w.logger.Info("Started")

	streams := w.streams
	w.streams = nil
	if len(streams) > 0 {
		w.logger.Info("Restarted, subbing existed streams")
		_, err = w.subStreamsNoLock(streams...)
		if err != nil {
			w.logger.Error("Cannot Sub Stream", "err", err)
			w.conn.Close()
			w.ctxCancel()
			return err
		}
	}

	w.status = RAW_WS_STATUS_STARTED

	return nil
}

var rawWsRestartAfterErrLimiter = limiter.New(time.Minute, 10)

func (w *RawWs) restartAfterErr() error {
	if !rawWsRestartAfterErrLimiter.TryWait() {
		return ErrWsTooFrequentRestart
	}
	w.logger.Info("Restarting")
	w.gmu.Lock()
	defer w.gmu.Unlock()
	if w.status == RAW_WS_STATUS_CLOSED {
		return ErrWsClientClosed
	}
	if w.status != RAW_WS_STATUS_STARTED {
		return fmt.Errorf("bnc: ws client status is %d, cannot restart", w.status)
	}
	w.status = RAW_WS_STATUS_NEW
	return w.startNoLock()
}

func (w *RawWs) Wait() (d []byte, err error) {
	for {
		t, d, err := w.conn.ReadMessage()
		if err != nil {
			e := w.restartAfterErr()
			if e != nil {
				return nil, fmt.Errorf("bnc: ws has err: %w, cannot restart: %w", err, e)
			}
			return nil, err
		}
		switch t {
		case websocket.PingMessage:
			w.logger.Info("Server ping received", "msg", string(d))
			if w.conn != nil {
				err := w.conn.WriteMessage(websocket.PongMessage, d)
				if err != nil {
					e := w.restartAfterErr()
					if e != nil {
						return nil, fmt.Errorf("bnc: ws has err: %w, cannot restart: %w", err, e)
					}
					return nil, err
				}
			} else {
				e := w.restartAfterErr()
				if e != nil {
					return nil, fmt.Errorf("bnc: ws has err: %w, cannot restart: %w", err, e)
				}
				return nil, errors.New("bnc: cannot write pong msg, conn is nil")
			}
		case websocket.PongMessage:
			w.logger.Info("Server pong received", "msg", string(d))
		case websocket.TextMessage:
			return d, nil
		case websocket.BinaryMessage:
			w.logger.Info("Server binary received", "msg", string(d), "binary", d)
		case websocket.CloseMessage:
			w.logger.Info("Server closed message", "msg", string(d))
		}
	}
}

var muWsReqId sync.Mutex
var latestWsReqId int64

func (w *RawWs) reqId() int64 {
	muWsReqId.Lock()
	defer muWsReqId.Unlock()
	ns := time.Now().UnixNano()
	if ns > latestWsReqId {
		latestWsReqId = ns
	} else {
		latestWsReqId++
	}
	return latestWsReqId
}

func (w *RawWs) SendMsg(method WsMethod, params any) (id int64, err error) {
	// cannot send msg concurrently
	// so remeber to lock the gmu before sending msg
	w.gmu.Lock()
	defer w.gmu.Unlock()
	return w.sendMsgNoLock(method, params)
}

func (w *RawWs) sendMsgNoLock(method WsMethod, params any) (id int64, err error) {
	// cannot send msg concurrently
	// so remeber to lock the gmu before sending msg
	if w.conn == nil {
		return 0, errors.New("bnc: nil conn")
	}
	id = w.reqId()
	d, err := json.Marshal(WsReq[any, int64]{
		Method: method,
		Params: params,
		Id:     id,
	})
	if err != nil {
		return
	}
	if !w.reqLimiter.TryWait() {
		return 0, ErrWsTooFrequentRequest
	}
	err = w.conn.WriteMessage(websocket.TextMessage, d)
	if err != nil {
		return
	}
	w.respByReqId.SetKV(id, make(chan WsResp[any], 1))
	return
}

func (w *RawWs) RespCh(id int64) (ch <-chan WsResp[any], ok bool) {
	return w.respByReqId.GetVWithOk(id)
}

func rawWsNewStreamFilter(oldStreams, subingStreams []string, maxStreams int) (existedStreams, newStreams, remainingParams []string) {
	for _, s := range subingStreams {
		if slices.Contains(oldStreams, s) {
			existedStreams = append(existedStreams, s)
		} else {
			newStreams = append(newStreams, s)
		}
	}
	if maxStreams < 0 {
		return
	}
	extraNum := len(oldStreams) + len(newStreams) - maxStreams
	if extraNum > 0 {
		ns := newStreams
		i := len(ns) - extraNum
		newStreams = ns[:i]
		remainingParams = ns[i:]
	}
	return
}

type StreamsSubResult struct {
	ExistedStreams  []string
	NewSubedStreams []string
	UnsubedStreams  []string
	Id              int64
}

func (w *RawWs) subStreamsNoLock(streams ...string) (result StreamsSubResult, err error) {
	defer func() {
		if err == nil && len(result.UnsubedStreams) > 0 {
			err = ErrNotAllStreamSubed
		}
	}()
	existed, newSubed, unsubed := rawWsNewStreamFilter(w.streams, streams, w.cfg.MaxStream)
	result.ExistedStreams = existed
	result.UnsubedStreams = append(result.UnsubedStreams, newSubed...)
	result.UnsubedStreams = append(result.UnsubedStreams, unsubed...)
	if len(newSubed) == 0 {
		return
	}

	id, err := w.sendMsgNoLock(WsMethodSub, newSubed)
	if err != nil {
		return
	}
	result.Id = id
	w.streams = append(w.streams, newSubed...)
	result.NewSubedStreams = newSubed
	result.UnsubedStreams = unsubed
	return
}

// Sub sub streams, binance dose not check stream validity,
// even if the stream is not valid, it will return success
// so we need to check the stream validity by ourselves
func (w *RawWs) Sub(streams ...string) (result StreamsSubResult, err error) {
	w.gmu.Lock()
	defer w.gmu.Unlock()
	return w.subStreamsNoLock(streams...)
}

func (w *RawWs) unsubStreamsNoLock(streams ...string) (id int64, err error) {
	id, err = w.sendMsgNoLock(WsMethodUnsub, streams)
	if err != nil {
		return
	}
	for _, s := range streams {
		i := slices.Index(w.streams, s)
		if i > -1 {
			w.streams = slices.Delete(w.streams, i, i+1)
		}
	}
	return
}

func (w *RawWs) Unsub(streams ...string) (id int64, err error) {
	w.gmu.Lock()
	defer w.gmu.Unlock()
	return w.unsubStreamsNoLock(streams...)
}

func (w *RawWs) newAndKeepListenKey() (string, error) {
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

func (w *RawWs) newListenKey() (string, error) {
	return w.user.NewListenKey(w.cfg.ListenKeyUrl)
}

func (w *RawWs) listenKeyKeeper(ctx context.Context) {
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
