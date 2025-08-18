package bnc

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/cex-go/ob"
	"github.com/dwdwow/props"
)

type OrderBookClient struct {
	spClt *orderBookMergedClient
	umClt *orderBookMergedClient
	cmClt *orderBookMergedClient
}

func NewOrderBookClient(ctx context.Context, logger *slog.Logger) *OrderBookClient {
	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("ws", "bnc_ob_clt")
	return &OrderBookClient{
		spClt: newOrderBookMergedClient(ctx, cex.SYMBOL_TYPE_SPOT, logger),
		umClt: newOrderBookMergedClient(ctx, cex.SYMBOL_TYPE_UM_FUTURES, logger),
		cmClt: newOrderBookMergedClient(ctx, cex.SYMBOL_TYPE_CM_FUTURES, logger),
	}
}

func (oc *OrderBookClient) SubSymbols(symbolType cex.SymbolType, symbols ...string) (unsubed []string, err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return oc.spClt.SubSymbols(symbols...)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return oc.umClt.SubSymbols(symbols...)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return oc.cmClt.SubSymbols(symbols...)
	}
	return nil, errors.New("bnc: unknown symbol type")
}

func (oc *OrderBookClient) UnsubSymbols(symbolType cex.SymbolType, symbols ...string) (err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return oc.spClt.UnsubSymbols(symbols...)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return oc.umClt.UnsubSymbols(symbols...)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return oc.cmClt.UnsubSymbols(symbols...)
	}
	return errors.New("bnc: unknown symbol type")
}

func (oc *OrderBookClient) NewCh(symbolType cex.SymbolType, symbol string) (ch <-chan ob.Data, err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return oc.spClt.NewCh(symbol)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return oc.umClt.NewCh(symbol)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return oc.cmClt.NewCh(symbol)
	}
	return nil, errors.New("bnc: unknown symbol type")
}

func (oc *OrderBookClient) RemoveCh(symbolType cex.SymbolType, ch <-chan ob.Data) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		oc.spClt.RemoveCh(ch)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		oc.umClt.RemoveCh(ch)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		oc.cmClt.RemoveCh(ch)
	}
}

type orderBookMergedClient struct {
	mu sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc

	symbolType cex.SymbolType

	clts []*orderBookBaseClient
	sybs *props.SafeRWMap[string, *orderBookBaseClient]

	logger *slog.Logger
}

func newOrderBookMergedClient(ctx context.Context, symbolType cex.SymbolType, logger *slog.Logger) *orderBookMergedClient {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("bnc_ob_merged_clt", symbolType)
	return &orderBookMergedClient{
		ctx:        ctx,
		cancel:     cancel,
		symbolType: symbolType,
		sybs:       props.NewSafeRWMap[string, *orderBookBaseClient](),
		logger:     logger,
	}
}

func (oc *orderBookMergedClient) scanAllSymbols() {
	for _, clt := range oc.clts {
		for _, s := range clt.workers.Keys() {
			oc.sybs.SetIfNotExists(s, clt)
		}
	}
}

func (oc *orderBookMergedClient) SubSymbols(symbols ...string) (unsubed []string, err error) {
	oc.mu.Lock()
	defer oc.mu.Unlock()
	defer oc.scanAllSymbols()
	for _, s := range symbols {
		if !oc.sybs.HasKey(s) {
			unsubed = append(unsubed, s)
		}
	}
	for _, clt := range oc.clts {
		if len(unsubed) == 0 {
			return
		}
		unsubed, err = clt.SubSymbols(unsubed...)
		if err != nil && err != ErrNotAllStreamSubed {
			oc.logger.Error("bnc: sub symbols failed", "err", err)
			return
		}
	}
	for {
		if len(unsubed) == 0 {
			return
		}
		clt := newOrderBookBaseClient(oc.ctx, oc.symbolType, oc.logger)
		unsubed, err = clt.SubSymbols(unsubed...)
		if err != nil && err != ErrNotAllStreamSubed {
			oc.logger.Error("bnc: sub symbols failed", "err", err)
			return
		}
		oc.clts = append(oc.clts, clt)
	}
}

func (oc *orderBookMergedClient) UnsubSymbols(symbols ...string) (err error) {
	oc.mu.Lock()
	defer oc.mu.Unlock()
	defer oc.scanAllSymbols()
	for _, s := range symbols {
		clt, ok := oc.sybs.GetVWithOk(s)
		if ok {
			err = clt.UnsubSymbols(s)
			if err != nil {
				return
			}
		}
	}
	return nil
}

// NewCh will return a channel that will receive the ob data
// if symbol is not found, ok will be false,
// should subscribe the symbol first
func (oc *orderBookMergedClient) NewCh(symbol string) (ch <-chan ob.Data, err error) {
	clt, ok := oc.sybs.GetVWithOk(symbol)
	if !ok {
		return nil, errors.New("bnc: symbol not found")
	}
	return clt.NewCh(symbol), nil
}

func (oc *orderBookMergedClient) RemoveCh(ch <-chan ob.Data) {
	for _, clt := range oc.clts {
		clt.RemoveCh(ch)
	}
}

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

type orderBookBaseClient struct {
	ctx    context.Context
	cancel context.CancelFunc

	sybType   cex.SymbolType
	obGetter  obGetter
	obUpdater obUpdater

	ws *RawWsClient

	_cache  *props.SafeRWMap[string, []WsDepthMsg]
	_exist  *props.SafeRWMap[string, bool]
	_ods    *props.SafeRWMap[string, ob.Data]
	workers *props.SafeRWMap[string, chan WsDepthMsg]

	rawMsgCh <-chan RawWsClientMsg

	radio *props.Radio[ob.Data]

	logger *slog.Logger
}

func newOrderBookBaseClient(ctx context.Context, symbolType cex.SymbolType, logger *slog.Logger) *orderBookBaseClient {
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
	radio := props.NewRadio(
		props.WithFanoutDur[ob.Data](time.Second),
		props.WithFanoutLogger[ob.Data](logger),
		props.WithFanoutChCap[ob.Data](10000),
	)
	clt := &orderBookBaseClient{
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
		workers:   props.NewSafeRWMap[string, chan WsDepthMsg](),
		radio:     radio,
		logger:    logger,
	}
	clt.run()
	return clt
}

func (oc *orderBookBaseClient) SymbolType() cex.SymbolType {
	return oc.sybType
}

func (oc *orderBookBaseClient) symbolType() cex.SymbolType {
	return oc.sybType
}

func (oc *orderBookBaseClient) cache() *props.SafeRWMap[string, []WsDepthMsg] {
	return oc._cache
}

func (oc *orderBookBaseClient) exist() *props.SafeRWMap[string, bool] {
	return oc._exist
}

func (oc *orderBookBaseClient) ods() *props.SafeRWMap[string, ob.Data] {
	return oc._ods
}

// SubSymbols will subscribe to new symbols
// symbol is case insensitive
func (oc *orderBookBaseClient) SubSymbols(symbols ...string) (unsubed []string, err error) {
	return oc.sub(symbols...)
}

func (oc *orderBookBaseClient) UnsubSymbols(symbols ...string) (err error) {
	return oc.unsub(symbols...)
}

// NewCh will return a channel that will receive the ob data
// symbol is case insensitive
func (oc *orderBookBaseClient) NewCh(symbol string) <-chan ob.Data {
	return oc.radio.Sub(symbol)
}

func (oc *orderBookBaseClient) RemoveCh(ch <-chan ob.Data) {
	oc.radio.UnsubAll(ch)
}

func (oc *orderBookBaseClient) createSubParams(symbols ...string) []string {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@depth@100ms")
	}
	return params
}

func (oc *orderBookBaseClient) sub(symbols ...string) (unsubed []string, err error) {
	defer func() {
		for _, s := range symbols {
			if !slices.Contains(unsubed, s) {
				oc.newWorker(s)
			}
		}
	}()
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

func (oc *orderBookBaseClient) unsub(symbols ...string) (err error) {
	params := oc.createSubParams(symbols...)
	err = oc.ws.UnsubStream(params...)
	if err != nil {
		return err
	}
	oc.workers.DeleteMany(symbols)
	// if delete cache data, may create some problems
	// for _, symbol := range symbols {
	// 	oc._exist.Delete(symbol)
	// 	oc._cache.Delete(symbol)
	// 	oc._ods.Delete(symbol)
	// }
	// remove workers
	return nil
}

func (oc *orderBookBaseClient) run() {
	go oc.listener()
}

func (oc *orderBookBaseClient) listener() {
	for {
		select {
		case msg := <-oc.rawMsgCh:
			oc.handle(msg)
		case <-oc.ctx.Done():
			return
		}
	}
}

func (oc *orderBookBaseClient) handle(msg RawWsClientMsg) {
	data := msg.Data
	err := msg.Err
	if err != nil {
		go oc.broadcastObs(oc.makeAllEmptyObs(err))
		return
	}
	m := WsDepthMsg{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		go oc.broadcastObs(oc.makeAllEmptyObs(err))
		return
	}
	if m.EventType == WsEDepthUpdate {
		defer func() {
			if err := recover(); err != nil {
				oc.logger.Error("bnc: panic in ob ws msg handler", "err", err)
			}
		}()
		w, ok := oc.workers.GetVWithOk(m.Symbol)
		if ok {
			w <- m
		}
		return
	}
	if string(data) == "{\"result\":null,\"id\":\"1\"}" ||
		string(data) == "{\"result\":null,\"id\":1}" {
		// this means subscribe success
		return
	}
	oc.logger.Warn("bnc: unexpected ob ws msg, unknown event type", "msg", string(data))
}

func (oc *orderBookBaseClient) worker(ch chan WsDepthMsg) {
	for {
		select {
		case data, ok := <-ch:
			if !ok {
				return
			}
			ob := oc.update(data)
			if ob.Err == ErrCachingObDepthUpdate {
				continue
			}
			oc.radio.Broadcast(data.Symbol, ob)
		case <-oc.ctx.Done():
			return
		}
	}
}

func (oc *orderBookBaseClient) newWorker(symbol string) {
	oc.workers.SetIfNotExists(symbol, func() chan WsDepthMsg {
		ch := make(chan WsDepthMsg, 10000)
		go oc.worker(ch)
		return ch
	}())
}

func (oc *orderBookBaseClient) makeAllEmptyObs(err error) []ob.Data {
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

func (oc *orderBookBaseClient) broadcastObs(obs []ob.Data) {
	for _, o := range obs {
		if o.Err == ErrCachingObDepthUpdate {
			continue
		}
		oc.radio.Broadcast(o.Symbol, o)
	}
}

func (oc *orderBookBaseClient) update(depthData WsDepthMsg) ob.Data {
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

func (oc *orderBookBaseClient) setError(symbol string, err error) ob.Data {
	empty := ob.NewData(cex.BINANCE, oc.sybType, symbol)
	empty.SetErr(err)
	oc._ods.SetKV(symbol, empty)
	return empty
}

func (oc *orderBookBaseClient) cacheRawData(depthData WsDepthMsg) error {
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

func (oc *orderBookBaseClient) needQueryOb(symbol string) bool {
	obData, ok := oc._ods.GetVWithOk(symbol)
	return !ok || obData.Empty()
}

var ErrCachingObDepthUpdate = errors.New("bnc: caching ob depth update")

func (oc *orderBookBaseClient) queryOb(symbol string) error {
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
