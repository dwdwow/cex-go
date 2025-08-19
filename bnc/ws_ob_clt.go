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

type OrderBookWs struct {
	spClt *orderBookMergedWs
	umClt *orderBookMergedWs
	cmClt *orderBookMergedWs
}

func NewOrderBookWs(ctx context.Context, logger *slog.Logger) *OrderBookWs {
	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("ws", "bnc_ob_clt")
	return &OrderBookWs{
		spClt: newOrderBookMergedWs(ctx, cex.SYMBOL_TYPE_SPOT, logger),
		umClt: newOrderBookMergedWs(ctx, cex.SYMBOL_TYPE_UM_FUTURES, logger),
		cmClt: newOrderBookMergedWs(ctx, cex.SYMBOL_TYPE_CM_FUTURES, logger),
	}
}

func (oc *OrderBookWs) Sub(symbolType cex.SymbolType, symbols ...string) (unsubed []string, err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return oc.spClt.subSymbols(symbols...)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return oc.umClt.subSymbols(symbols...)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return oc.cmClt.subSymbols(symbols...)
	}
	return nil, errors.New("bnc: unknown symbol type")
}

func (oc *OrderBookWs) Unsub(symbolType cex.SymbolType, symbols ...string) (err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return oc.spClt.unsubSymbols(symbols...)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return oc.umClt.unsubSymbols(symbols...)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return oc.cmClt.unsubSymbols(symbols...)
	}
	return errors.New("bnc: unknown symbol type")
}

func (oc *OrderBookWs) NewCh(symbolType cex.SymbolType, symbol string) (ch <-chan ob.Data, err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return oc.spClt.newCh(symbol)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return oc.umClt.newCh(symbol)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return oc.cmClt.newCh(symbol)
	}
	return nil, errors.New("bnc: unknown symbol type")
}

func (oc *OrderBookWs) RemoveCh(ch <-chan ob.Data) {
	oc.spClt.removeCh(ch)
	oc.umClt.removeCh(ch)
	oc.cmClt.removeCh(ch)
}

type orderBookMergedWs struct {
	mu sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc

	symbolType cex.SymbolType

	clts []*orderBookBaseWs
	sybs *props.SafeRWMap[string, *orderBookBaseWs]

	logger *slog.Logger
}

func newOrderBookMergedWs(ctx context.Context, symbolType cex.SymbolType, logger *slog.Logger) *orderBookMergedWs {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("bnc_ob_merged_clt", symbolType)
	return &orderBookMergedWs{
		ctx:        ctx,
		cancel:     cancel,
		symbolType: symbolType,
		sybs:       props.NewSafeRWMap[string, *orderBookBaseWs](),
		logger:     logger,
	}
}

func (oc *orderBookMergedWs) scanAllSymbols() {
	for _, clt := range oc.clts {
		for _, s := range clt.workers.Keys() {
			oc.sybs.SetIfNotExists(s, clt)
		}
	}
}

func (oc *orderBookMergedWs) subSymbols(symbols ...string) (unsubed []string, err error) {
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
		unsubed, err = clt.subSymbols(unsubed...)
		if err != nil && err != ErrNotAllStreamSubed {
			oc.logger.Error("bnc: sub symbols failed", "err", err)
			return
		}
	}
	for {
		if len(unsubed) == 0 {
			return
		}
		clt := newOrderBookBaseWs(oc.ctx, oc.symbolType, oc.logger)
		unsubed, err = clt.subSymbols(unsubed...)
		if err != nil && err != ErrNotAllStreamSubed {
			oc.logger.Error("bnc: sub symbols failed", "err", err)
			return
		}
		oc.clts = append(oc.clts, clt)
	}
}

func (oc *orderBookMergedWs) unsubSymbols(symbols ...string) (err error) {
	oc.mu.Lock()
	defer oc.mu.Unlock()
	defer oc.scanAllSymbols()
	for _, s := range symbols {
		clt, ok := oc.sybs.GetVWithOk(s)
		if ok {
			err = clt.unsubSymbols(s)
			if err != nil {
				return
			}
		}
	}
	return nil
}

// newCh will return a channel that will receive the ob data
// if symbol is not found, ok will be false,
// should subscribe the symbol first
func (oc *orderBookMergedWs) newCh(symbol string) (ch <-chan ob.Data, err error) {
	clt, ok := oc.sybs.GetVWithOk(symbol)
	if !ok {
		return nil, errors.New("bnc: symbol not found")
	}
	return clt.newCh(symbol), nil
}

func (oc *orderBookMergedWs) removeCh(ch <-chan ob.Data) {
	for _, clt := range oc.clts {
		clt.removeCh(ch)
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

type orderBookBaseWs struct {
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

func newOrderBookBaseWs(ctx context.Context, symbolType cex.SymbolType, logger *slog.Logger) *orderBookBaseWs {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("bnc_ob_base_clt", symbolType)
	wsCfg := DefaultPublicWsCfg(symbolType)
	wsCfg.MaxStream = maxObWsStream
	wsCfg.FanoutTimerDur = 0
	ws := NewRawWsClient(ctx, wsCfg, logger)
	radio := props.NewRadio(
		// props.WithFanoutDur[ob.Data](time.Second),
		props.WithFanoutLogger[ob.Data](logger),
		props.WithFanoutChCap[ob.Data](10000),
	)
	clt := &orderBookBaseWs{
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

func (oc *orderBookBaseWs) symbolType() cex.SymbolType {
	return oc.sybType
}

func (oc *orderBookBaseWs) cache() *props.SafeRWMap[string, []WsDepthMsg] {
	return oc._cache
}

func (oc *orderBookBaseWs) exist() *props.SafeRWMap[string, bool] {
	return oc._exist
}

func (oc *orderBookBaseWs) ods() *props.SafeRWMap[string, ob.Data] {
	return oc._ods
}

// subSymbols will subscribe to new symbols
// symbol is case insensitive
func (oc *orderBookBaseWs) subSymbols(symbols ...string) (unsubed []string, err error) {
	return oc.sub(symbols...)
}

func (oc *orderBookBaseWs) unsubSymbols(symbols ...string) (err error) {
	return oc.unsub(symbols...)
}

// newCh will return a channel that will receive the ob data
// symbol is case insensitive
func (oc *orderBookBaseWs) newCh(symbol string) <-chan ob.Data {
	return oc.radio.Sub(symbol)
}

func (oc *orderBookBaseWs) removeCh(ch <-chan ob.Data) {
	oc.radio.UnsubAll(ch)
}

func (oc *orderBookBaseWs) createSubParams(symbols ...string) []string {
	var params []string
	for _, symbol := range symbols {
		params = append(params, strings.ToLower(symbol)+"@depth@100ms")
	}
	return params
}

func (oc *orderBookBaseWs) sub(symbols ...string) (unsubed []string, err error) {
	defer func() {
		// create new worker for new symbols
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

func (oc *orderBookBaseWs) unsub(symbols ...string) (err error) {
	params := oc.createSubParams(symbols...)
	err = oc.ws.UnsubStream(params...)
	if err != nil {
		return err
	}
	closedWorkers := oc.workers.DeleteMany(symbols)
	for _, ch := range closedWorkers {
		close(ch)
	}
	for _, symbol := range symbols {
		od := oc.newObWithErr(symbol, ErrWsStreamUnsubed)
		oc.radio.Broadcast(symbol, od)
	}
	// if delete cache data, may create some problems
	for _, symbol := range symbols {
		oc._exist.Delete(symbol)
		oc._cache.Delete(symbol)
		oc._ods.Delete(symbol)
	}
	return nil
}

func (oc *orderBookBaseWs) run() {
	go oc.listener()
}

func (oc *orderBookBaseWs) listener() {
	for {
		select {
		case msg := <-oc.rawMsgCh:
			oc.handle(msg)
		case <-oc.ctx.Done():
			return
		}
	}
}

func (oc *orderBookBaseWs) handle(msg RawWsClientMsg) {
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

func (oc *orderBookBaseWs) worker(ch chan WsDepthMsg) {
	for {
		select {
		case data, ok := <-ch:
			if !ok {
				return
			}
			ob, updated := oc.update(data)
			if updated {
				oc.radio.Broadcast(data.Symbol, ob)
			} else {
				if ob.Err == ErrCachingObDepthUpdate {
					continue
				}
			}
		case <-oc.ctx.Done():
			return
		}
	}
}

func (oc *orderBookBaseWs) newWorker(symbol string) {
	oc.workers.SetIfNotExists(symbol, func() chan WsDepthMsg {
		ch := make(chan WsDepthMsg, 10000)
		go oc.worker(ch)
		return ch
	}())
}

func (oc *orderBookBaseWs) makeAllEmptyObs(err error) []ob.Data {
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

func (oc *orderBookBaseWs) broadcastObs(obs []ob.Data) {
	for _, o := range obs {
		if o.Err == ErrCachingObDepthUpdate {
			continue
		}
		oc.radio.Broadcast(o.Symbol, o)
	}
}

func (oc *orderBookBaseWs) update(depthData WsDepthMsg) (od ob.Data, updated bool) {
	symbol := depthData.Symbol
	err := oc.cacheRawData(depthData)
	if err != nil {
		return oc.newObWithErr(symbol, err), true
	}
	if oc.needQueryOb(symbol) {
		// check cache
		oldCache := oc._cache.GetV(symbol)
		if len(oldCache) < 10 {
			return oc.newObWithErr(symbol, ErrCachingObDepthUpdate), false
		}
		if !publicRestLimitter.TryWait() {
			return oc.newObWithErr(symbol, ErrObQueryLimit), true
		}
		err = oc.queryOb(symbol)
		if err != nil {
			return oc.newObWithErr(symbol, err), true
		}
	}
	od = oc.obUpdater(oc, depthData)
	oc._ods.SetKV(symbol, od)
	return od, true
}

func (oc *orderBookBaseWs) newObWithErr(symbol string, err error) ob.Data {
	empty := ob.NewData(cex.BINANCE, oc.sybType, symbol)
	empty.SetErr(err)
	oc._ods.SetKV(symbol, empty)
	return empty
}

func (oc *orderBookBaseWs) cacheRawData(depthData WsDepthMsg) error {
	symbol := depthData.Symbol
	oldCache := oc._cache.GetV(symbol)
	if len(oldCache) > 1000 {
		oc._cache.SetKV(symbol, nil)
		return ErrObDepthUpdateCacheTooLarge
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

func (oc *orderBookBaseWs) needQueryOb(symbol string) bool {
	obData, ok := oc._ods.GetVWithOk(symbol)
	return !ok || obData.Empty()
}

func (oc *orderBookBaseWs) queryOb(symbol string) error {
	// oldCache := oc._cache.GetV(symbol)
	// if len(oldCache) < 10 {
	// 	return ErrCachingObDepthUpdate
	// }
	// if time.Now().UnixMilli()-lastObQueryFailTsMilli.Get() < 3000 {
	// 	return errors.New("bnc: can not query orderbook within 3000 milliseconds")
	// }
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
		// lastObQueryFailTsMilli.Set(time.Now().UnixMilli())
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
