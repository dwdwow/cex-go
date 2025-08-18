package bnc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/cex-go/ob"
	"github.com/dwdwow/props"
)

type UmObWs struct {
	ctx context.Context

	ws *RawWsClient

	cache *props.SafeRWMap[string, []WsDepthMsg]
	exist *props.SafeRWMap[string, bool]
	ods   *props.SafeRWMap[string, ob.Data]

	rawSuber <-chan RawWsClientMsg

	radio *props.Radio[ob.Data]

	logger *slog.Logger
}

func NewUmObWs(ctx context.Context, logger *slog.Logger) *UmObWs {
	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("ws", "bnc_um_ob")
	cfg := DefaultUmPublicWsCfg()
	cfg.MaxStream = maxObWsStream
	ws := NewRawWsClient(ctx, cfg, logger)
	clt := &UmObWs{
		ctx:      ctx,
		ws:       ws,
		rawSuber: ws.Sub(),
		cache:    props.NewSafeRWMap[string, []WsDepthMsg](),
		exist:    props.NewSafeRWMap[string, bool](),
		ods:      props.NewSafeRWMap[string, ob.Data](),
		radio:    props.NewRadio(props.WithFanoutDur[ob.Data](time.Second)),
		logger:   logger,
	}
	clt.run()
	return clt
}

// SubNewSymbols will subscribe to new symbols
// symbol is case insensitive
func (w *UmObWs) SubNewSymbols(symbols ...string) (unsubed []string, err error) {
	return w.sub(symbols...)
}

// NewCh will return a channel that will receive the ob data
// symbol is case insensitive
func (w *UmObWs) NewCh(symbol string) <-chan ob.Data {
	return w.radio.Sub(symbol)
}

func (w *UmObWs) RemoveCh(ch <-chan ob.Data) {
	w.radio.UnsubAll(ch)
}

func (w *UmObWs) run() {
	go w.listener()
}

func (w *UmObWs) listener() {
	for {
		select {
		case msg := <-w.rawSuber:
			obs := w.handle(msg)
			if len(obs) > 0 {
				for _, o := range obs {
					if o.Err == ErrCachingObDepthUpdate {
						continue
					}
					w.radio.Broadcast(o.Symbol, o)
				}
			}
		case <-w.ctx.Done():
			return
		}
	}
}

func (w *UmObWs) sub(symbols ...string) (unsubed []string, err error) {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@depth")
	}
	res, err := w.ws.SubStream(params...)
	if err == nil {
		return
	}
	for _, p := range res.UnsubedParams {
		for _, s := range symbols {
			if strings.Contains(p, strings.ToLower(s)) {
				unsubed = append(unsubed, s)
			}
		}
	}
	return unsubed, err
}

func (w *UmObWs) handle(msg RawWsClientMsg) []ob.Data {
	data := msg.Data
	err := msg.Err
	if err != nil {
		// set ob data to empty
		obs := w.makeAllEmptyObs(err)
		return obs
	}
	m := new(WsDepthMsg)
	err = json.Unmarshal(data, m)
	if err != nil {
		obs := w.makeAllEmptyObs(err)
		return obs
	}
	if m.EventType == WsEDepthUpdate {
		obData := w.update(*m)
		return []ob.Data{obData}
	}
	if string(data) == "{\"result\":null,\"id\":\"1\"}" ||
		string(data) == "{\"result\":null,\"id\":1}" {
		// this means subscribe success
		return nil
	}
	w.logger.Warn("bnc: unexpected ob ws msg, unknown event type", "msg", string(data))
	return nil
}

func (w *UmObWs) makeAllEmptyObs(err error) []ob.Data {
	symbols := w.exist.Keys()
	var obs []ob.Data
	for _, symbol := range symbols {
		empty := ob.NewData(cex.BINANCE, cex.SYMBOL_TYPE_UM_FUTURES, symbol)
		empty.SetErr(err)
		w.ods.SetKV(symbol, empty)
		obs = append(obs, empty)
	}
	return obs
}

func (w *UmObWs) update(depthData WsDepthMsg) ob.Data {
	symbol := depthData.Symbol
	err := w.cacheRawData(depthData)
	if err != nil {
		return w.setError(symbol, err)
	}
	if w.needQueryOb(symbol) {
		err = w.queryOb(symbol)
		if err != nil {
			return w.setError(symbol, err)
		}
	}
	o := w.updateOb(depthData)
	w.ods.SetKV(symbol, o)
	return o
}

func (w *UmObWs) setError(symbol string, err error) ob.Data {
	empty := ob.NewData(cex.BINANCE, cex.SYMBOL_TYPE_UM_FUTURES, symbol)
	empty.SetErr(err)
	w.ods.SetKV(symbol, empty)
	return empty
}

func (w *UmObWs) cacheRawData(depthData WsDepthMsg) error {
	symbol := depthData.Symbol
	oldCache := w.cache.GetV(symbol)
	if len(oldCache) > 100 {
		// clear cache
		w.cache.SetKV(symbol, nil)
		return errors.New("bnc: too many ob depth data cache")
	}
	newCache := append(oldCache, depthData)
	sort.Slice(newCache, func(i, j int) bool {
		iLastId := newCache[i].LastId
		jPu := newCache[j].PLastId
		return iLastId == jPu
	})
	cacheLen := len(newCache)
	for i := 0; i < cacheLen-1; i++ {
		if newCache[i].LastId != newCache[i+1].PLastId {
			w.cache.SetKV(symbol, nil)
			return errors.New("bnc: ob depth data cache is not continuous")
		}
	}
	w.cache.SetKV(symbol, newCache)
	return nil
}

func (w *UmObWs) needQueryOb(symbol string) bool {
	obData, ok := w.ods.GetVWithOk(symbol)
	return !ok || obData.Empty()
}

func (w *UmObWs) queryOb(symbol string) error {
	oldCache := w.cache.GetV(symbol)
	if len(oldCache) < 10 {
		return ErrCachingObDepthUpdate
	}
	if time.Now().UnixMilli()-lastObQueryFailTsMilli.Get() < 3000 {
		return errors.New("bnc: can not query orderbook within 3000 milliseconds")
	}
	if w.exist.SetKV(symbol, true) {
		return errors.New("bnc: lock to query orderbook")
	}
	defer w.exist.SetKV(symbol, false)
	// because ws orderbook default limit is 1000
	// so limit must be 1000
	rawOrderbook, err := GetUMOrderBook(ParamsOrderBook{
		Symbol: symbol,
		Limit:  1000,
	})
	if err != nil {
		lastObQueryFailTsMilli.Set(time.Now().UnixMilli())
		return err
	}
	obData := ob.Data{
		Cex:     cex.BINANCE,
		Type:    cex.SYMBOL_TYPE_UM_FUTURES,
		Symbol:  symbol,
		Version: strconv.FormatInt(rawOrderbook.LastUpdateID, 10),
		Time:    time.Now().UnixMilli(),
		Asks:    rawOrderbook.Asks,
		Bids:    rawOrderbook.Bids,
	}
	w.ods.SetKV(symbol, obData)
	return nil
}

func (w *UmObWs) updateOb(depthData WsDepthMsg) ob.Data {
	symbol := depthData.Symbol
	buffer := w.cache.GetV(symbol)
	empty := ob.NewData(cex.BINANCE, cex.SYMBOL_TYPE_UM_FUTURES, symbol)
	obData, ok := w.ods.GetVWithOk(symbol)
	if !ok || obData.Empty() {
		empty.SetErr(errors.New("bnc: unexpected error: binance update ob: if !ok || obData.Empty()"))
		return empty
	}
	currentVersion, err := strconv.ParseInt(obData.Version, 10, 64)
	if err != nil {
		empty.SetErr(fmt.Errorf("bnc: can not parse ob data version %s, err: %w", obData.Version, err))
		return empty
	}
	if buffer[0].PLastId > currentVersion {
		empty.SetErr(errors.New("bnc: current ob version is small"))
		return empty
	}
	lastIndex := 0
	_id := int64(0)
	for i, _depthData := range buffer {
		firstId := _depthData.FirstId
		lastId := _depthData.LastId
		pu := _depthData.PLastId
		if _id > 0 {
			if pu == _id {
				_id = lastId
			} else {
				w.cache.SetKV(symbol, buffer[i:])
				empty.SetErr(fmt.Errorf("bnc: ob depth update is not continuous, websocket message may be blocked"))
				return empty
			}
		} else {
			_id = lastId
		}
		if pu != currentVersion {
			if firstId <= currentVersion && lastId >= currentVersion {
				// first updating
			} else {
				if pu > currentVersion {
					lastIndex = i - 1
					break
				}
				if lastId < currentVersion {
					lastIndex = i
					continue
				}
			}
		}
		lastIndex = i
		asks := _depthData.Asks
		bids := _depthData.Bids
		currentVersion = lastId
		for _, ask := range asks {
			price, err := strconv.ParseFloat(ask[0], 64)
			if err != nil {
				empty.SetErr(fmt.Errorf("bnc: can not parse ask price %s, err: %w", ask[0], err))
				return empty
			}
			qty, err := strconv.ParseFloat(ask[1], 64)
			if err != nil {
				empty.SetErr(fmt.Errorf("bnc: can not parse ask qty %s, err: %w", ask[1], err))
				return empty
			}
			err = obData.UpdateAskDeltas(ob.Book{{price, qty}}, strconv.FormatInt(currentVersion, 10))
			if err != nil {
				empty.SetErr(fmt.Errorf("bnc: can not update ask deltas, err: %w", err))
				return empty
			}
		}
		for _, bid := range bids {
			price, err := strconv.ParseFloat(bid[0], 64)
			if err != nil {
				empty.SetErr(fmt.Errorf("bnc: can not parse bid price %s, err: %w", bid[0], err))
				return empty
			}
			qty, err := strconv.ParseFloat(bid[1], 64)
			if err != nil {
				empty.SetErr(fmt.Errorf("bnc: can not parse bid qty %s, err: %w", bid[1], err))
				return empty
			}
			err = obData.UpdateBidDeltas(ob.Book{{price, qty}}, strconv.FormatInt(currentVersion, 10))
			if err != nil {
				empty.SetErr(fmt.Errorf("bnc: can not update bid deltas, err: %w", err))
				return empty
			}
		}
	}
	if len(buffer) <= lastIndex+1 {
		w.cache.SetKV(symbol, []WsDepthMsg{})
	} else {
		w.cache.SetKV(symbol, buffer[lastIndex+1:])
	}
	return obData
}
