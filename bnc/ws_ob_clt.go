package bnc

import (
	"context"
	"encoding/json"
	"errors"
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

type obGetter func(ParamsOrderBook) (OrderBook, error)

func defaultObGetter(symbolType cex.SymbolType) obGetter {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return GetSpotOrderBook
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return GetUMOrderBook
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return GetCMOrderBook
	}
	panic("bnc: unknown symbol type")
}

type OrderBookClient struct {
	ctx    context.Context
	cancel context.CancelFunc

	sybType   cex.SymbolType
	obGetter  obGetter
	obUpdater obUpdater

	ws *RawWsClient

	_cache *props.SafeRWMap[string, []WsDepthMsg]
	_exist *props.SafeRWMap[string, bool]
	_ods   *props.SafeRWMap[string, ob.Data]

	rawMsgCh <-chan RawWsClientMsg

	radio *props.Radio[ob.Data]

	logger *slog.Logger
}

func NewOrderBookClient(ctx context.Context, symbolType cex.SymbolType, logger *slog.Logger) *OrderBookClient {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("bnc_ob_clt", symbolType)
	wsCfg := DefaultPublicWsCfg(symbolType)
	wsCfg.MaxStream = maxObWsStream
	ws := NewRawWsClient(ctx, wsCfg, logger)
	clt := &OrderBookClient{
		ctx:       ctx,
		cancel:    cancel,
		sybType:   symbolType,
		obGetter:  defaultObGetter(symbolType),
		obUpdater: defaultObUpdater(symbolType),
		ws:        ws,
		rawMsgCh:  ws.Sub(),
		_cache:    props.NewSafeRWMap[string, []WsDepthMsg](),
		_exist:    props.NewSafeRWMap[string, bool](),
		_ods:      props.NewSafeRWMap[string, ob.Data](),
		radio:     props.NewRadio(props.WithFanoutDur[ob.Data](time.Second)),
		logger:    logger,
	}
	clt.run()
	return clt
}

func (oc *OrderBookClient) SymbolType() cex.SymbolType {
	return oc.sybType
}

func (oc *OrderBookClient) symbolType() cex.SymbolType {
	return oc.sybType
}

func (oc *OrderBookClient) cache() *props.SafeRWMap[string, []WsDepthMsg] {
	return oc._cache
}

func (oc *OrderBookClient) exist() *props.SafeRWMap[string, bool] {
	return oc._exist
}

func (oc *OrderBookClient) ods() *props.SafeRWMap[string, ob.Data] {
	return oc._ods
}

// SubSymbols will subscribe to new symbols
// symbol is case insensitive
func (oc *OrderBookClient) SubSymbols(symbols ...string) (unsubed []string, err error) {
	return oc.sub(symbols...)
}

func (oc *OrderBookClient) UnsubSymbols(symbols ...string) (err error) {
	return oc.unsub(symbols...)
}

// NewCh will return a channel that will receive the ob data
// symbol is case insensitive
func (oc *OrderBookClient) NewCh(symbol string) <-chan ob.Data {
	return oc.radio.Sub(symbol)
}

func (oc *OrderBookClient) RemoveCh(ch <-chan ob.Data) {
	oc.radio.UnsubAll(ch)
}

func (oc *OrderBookClient) run() {
	go oc.listener()
}

func (oc *OrderBookClient) listener() {
	for {
		select {
		case msg := <-oc.rawMsgCh:
			obs := oc.handle(msg)
			if len(obs) > 0 {
				for _, o := range obs {
					if o.Err == ErrCachingObDepthUpdate {
						continue
					}
					oc.radio.Broadcast(o.Symbol, o)
				}
			}
		case <-oc.ctx.Done():
			return
		}
	}
}

func (oc *OrderBookClient) createSubParams(symbols ...string) []string {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@depth@100ms")
	}
	return params
}

func (oc *OrderBookClient) sub(symbols ...string) (unsubed []string, err error) {
	params := oc.createSubParams(symbols...)
	res, err := oc.ws.SubStream(params...)
	if err == nil {
		return
	}
	for _, p := range res.UnsubedParams {
		for _, s := range symbols {
			if strings.Contains(p, strings.ToLower(s)) {
				unsubed = append(unsubed, s)
				break
			}
		}
	}
	return unsubed, err
}

func (oc *OrderBookClient) unsub(symbols ...string) (err error) {
	params := oc.createSubParams(symbols...)
	err = oc.ws.UnsubStream(params...)
	if err != nil {
		return err
	}
	// if delete cache data, may create some problems
	// for _, symbol := range symbols {
	// 	oc._exist.Delete(symbol)
	// 	oc._cache.Delete(symbol)
	// 	oc._ods.Delete(symbol)
	// }
	return nil
}

func (oc *OrderBookClient) handle(msg RawWsClientMsg) []ob.Data {
	data := msg.Data
	err := msg.Err
	if err != nil {
		// set ob data to empty
		obs := oc.makeAllEmptyObs(err)
		return obs
	}
	m := new(WsDepthMsg)
	err = json.Unmarshal(data, m)
	if err != nil {
		obs := oc.makeAllEmptyObs(err)
		return obs
	}
	if m.EventType == WsEDepthUpdate {
		obData := oc.update(*m)
		return []ob.Data{obData}
	}
	if string(data) == "{\"result\":null,\"id\":\"1\"}" ||
		string(data) == "{\"result\":null,\"id\":1}" {
		// this means subscribe success
		return nil
	}
	oc.logger.Warn("bnc: unexpected ob ws msg, unknown event type", "msg", string(data))
	return nil
}

func (oc *OrderBookClient) makeAllEmptyObs(err error) []ob.Data {
	symbols := oc._exist.Keys()
	var obs []ob.Data
	for _, symbol := range symbols {
		empty := ob.NewData(cex.BINANCE, oc.sybType, symbol)
		empty.SetErr(err)
		oc._ods.SetKV(symbol, empty)
		obs = append(obs, empty)
	}
	return obs
}

func (oc *OrderBookClient) update(depthData WsDepthMsg) ob.Data {
	symbol := depthData.Symbol
	err := oc.cacheRawData(depthData)
	if err != nil {
		return oc.setError(symbol, err)
	}
	if oc.needQueryOb(symbol) {
		err = oc.queryOb(symbol)
		if err != nil {
			return oc.setError(symbol, err)
		}
	}
	o := oc.obUpdater(oc, depthData)
	oc._ods.SetKV(symbol, o)
	return o
}

func (oc *OrderBookClient) setError(symbol string, err error) ob.Data {
	empty := ob.NewData(cex.BINANCE, oc.sybType, symbol)
	empty.SetErr(err)
	oc._ods.SetKV(symbol, empty)
	return empty
}

func (oc *OrderBookClient) cacheRawData(depthData WsDepthMsg) error {
	symbol := depthData.Symbol
	oldCache := oc._cache.GetV(symbol)
	if len(oldCache) > 100 {
		oc._cache.SetKV(symbol, nil)
		return errors.New("bnc: too many ob depth data cache")
	}
	newCache := append(oldCache, depthData)
	// spot update data ids are continuous
	// um and cm update data ids are not continuous, must check pLastId
	switch oc.sybType {
	case cex.SYMBOL_TYPE_SPOT:
		sort.Slice(newCache, func(i, j int) bool {
			iLastId := newCache[i].LastId
			jFirstId := newCache[j].FirstId
			return iLastId+1 == jFirstId
		})
		for i := range len(newCache) - 1 {
			if newCache[i].LastId+1 != newCache[i+1].FirstId {
				oc._cache.SetKV(symbol, nil)
				return errors.New("bnc: ob depth data cache is not continuous")
			}
		}
	case cex.SYMBOL_TYPE_UM_FUTURES, cex.SYMBOL_TYPE_CM_FUTURES:
		sort.Slice(newCache, func(i, j int) bool {
			iLastId := newCache[i].LastId
			jPu := newCache[j].PLastId
			return iLastId == jPu
		})
		cacheLen := len(newCache)
		for i := 0; i < cacheLen-1; i++ {
			if newCache[i].LastId != newCache[i+1].PLastId {
				oc._cache.SetKV(symbol, nil)
				return errors.New("bnc: ob depth data cache is not continuous")
			}
		}
	}
	oc._cache.SetKV(symbol, newCache)
	return nil
}

func (oc *OrderBookClient) needQueryOb(symbol string) bool {
	obData, ok := oc._ods.GetVWithOk(symbol)
	return !ok || obData.Empty()
}

func (oc *OrderBookClient) queryOb(symbol string) error {
	oldCache := oc._cache.GetV(symbol)
	if len(oldCache) < 10 {
		return ErrCachingObDepthUpdate
	}
	if time.Now().UnixMilli()-lastObQueryFailTsMilli.Get() < 3000 {
		return errors.New("bnc: can not query orderbook within 3000 milliseconds")
	}
	if oc._exist.SetKV(symbol, true) {
		return errors.New("bnc: lock to query orderbook")
	}
	defer oc._exist.SetKV(symbol, false)
	// because ws orderbook default limit is 1000
	// so limit must be 1000
	rawOrderbook, err := oc.obGetter(ParamsOrderBook{
		Symbol: symbol,
		Limit:  1000,
	})
	if err != nil {
		lastObQueryFailTsMilli.Set(time.Now().UnixMilli())
		return err
	}
	obData := ob.Data{
		Cex:     cex.BINANCE,
		Type:    oc.sybType,
		Symbol:  symbol,
		Version: strconv.FormatInt(rawOrderbook.LastUpdateID, 10),
		Time:    time.Now().UnixMilli(),
		Asks:    rawOrderbook.Asks,
		Bids:    rawOrderbook.Bids,
	}
	oc._ods.SetKV(symbol, obData)
	return nil
}
