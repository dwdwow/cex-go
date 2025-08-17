package bnc

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/cex-go/ob"
	"github.com/dwdwow/props"
	"github.com/dwdwow/ws/wsclt"
	"github.com/gorilla/websocket"
)

type WsUMObMsgHandler struct {
	mgClt            *wsclt.MergedClient
	svDataCacheBySyb props.SafeRWMap[string, []WsDepthMsg]
	gettingObSybs    props.SafeRWMap[string, bool]
	obCacheBySyb     props.SafeRWMap[string, ob.Data]
}

func NewWsUMObMsgHandler(logger *slog.Logger) *WsUMObMsgHandler {
	mgClt := wsclt.
		NewMergedClient(FutureWsBaseUrl, true, maxTopicNumPerWs, logger).
		SetTopicSuber(topicSuber).
		SetTopicUnsuber(topicUnsuber).
		SetPong(pong)
	return &WsUMObMsgHandler{
		mgClt: mgClt,
	}
}

func (w *WsUMObMsgHandler) Name() cex.CexName {
	return cex.BINANCE
}

func (w *WsUMObMsgHandler) Type() cex.SymbolType {
	return cex.SYMBOL_TYPE_UM_FUTURES
}

func (w *WsUMObMsgHandler) Client() *wsclt.MergedClient {
	return w.mgClt
}

func (w *WsUMObMsgHandler) Topics(symbols ...string) []string {
	var topics []string
	for _, s := range symbols {
		topics = append(topics, CreateObTopic(s))
	}
	return topics
}

func (w *WsUMObMsgHandler) Handle(msg wsclt.MergedClientMsg) ([]ob.Data, error) {
	return w.handle(msg)
}

func (w *WsUMObMsgHandler) handle(msg wsclt.MergedClientMsg) ([]ob.Data, error) {
	if msg.Err != nil {
		// set ob data to empty
		var obs []ob.Data
		topics := msg.Client.Topics()
		for _, topic := range topics {
			topicSplit := strings.Split(topic, "@depth")
			if len(topicSplit) != 2 {
				// should not get here
				fmt.Println("bnc: unexpected error: binance future ob ws msg handle: can not parse topic", topic)
				continue
			}
			symbol := topicSplit[0]
			empty := ob.NewData(cex.BINANCE, cex.SYMBOL_TYPE_UM_FUTURES, symbol)
			empty.SetErr(msg.Err)
			w.obCacheBySyb.SetKV(symbol, empty)
			obs = append(obs, empty)
		}
		return obs, nil
	}
	if msg.MsgType != websocket.TextMessage {
		return nil, fmt.Errorf("bnc: ws receive unknown msg type %v", msg.MsgType)
	}
	msgData := msg.Data
	data := new(WsDepthMsg)
	err := json.Unmarshal(msgData, data)
	if err != nil {
		return nil, fmt.Errorf("bnc: ws msg unmarshal, msg: %v, %w", string(msgData), err)
	}
	if data.EventType == WsEDepthUpdate {
		obData := w.update(*data)
		return []ob.Data{obData}, nil
	}
	return nil, nil
}

func (w *WsUMObMsgHandler) update(depthData WsDepthMsg) ob.Data {
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
	w.obCacheBySyb.SetKV(symbol, o)
	return o
}

func (w *WsUMObMsgHandler) setError(symbol string, err error) ob.Data {
	empty := ob.NewData(cex.BINANCE, cex.SYMBOL_TYPE_UM_FUTURES, symbol)
	empty.SetErr(err)
	w.obCacheBySyb.SetKV(symbol, empty)
	return empty
}

func (w *WsUMObMsgHandler) cacheRawData(depthData WsDepthMsg) error {
	symbol := depthData.Symbol
	oldCache := w.svDataCacheBySyb.GetV(symbol)
	if len(oldCache) > 100 {
		// clear cache
		w.svDataCacheBySyb.SetKV(symbol, nil)
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
			w.svDataCacheBySyb.SetKV(symbol, nil)
			return errors.New("bnc: ob depth data cache is not continuous")
		}
	}
	w.svDataCacheBySyb.SetKV(symbol, newCache)
	return nil
}

func (w *WsUMObMsgHandler) needQueryOb(symbol string) bool {
	obData, ok := w.obCacheBySyb.GetVWithOk(symbol)
	return !ok || obData.Empty()
}

func (w *WsUMObMsgHandler) queryOb(symbol string) error {
	oldCache := w.svDataCacheBySyb.GetV(symbol)
	if len(oldCache) < 10 {
		return ErrCachingObDepthUpdate
	}
	if time.Now().UnixMilli()-lastObQueryFailTsMilli.Get() < 3000 {
		return errors.New("bnc: can not query orderbook within 3000 milliseconds")
	}
	if w.gettingObSybs.SetKV(symbol, true) {
		return errors.New("bnc: lock to query orderbook")
	}
	defer w.gettingObSybs.SetKV(symbol, false)
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
	w.obCacheBySyb.SetKV(symbol, obData)
	return nil
}

func (w *WsUMObMsgHandler) updateOb(depthData WsDepthMsg) ob.Data {
	symbol := depthData.Symbol
	buffer := w.svDataCacheBySyb.GetV(symbol)
	empty := ob.NewData(cex.BINANCE, cex.SYMBOL_TYPE_UM_FUTURES, symbol)
	obData, ok := w.obCacheBySyb.GetVWithOk(symbol)
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
				safeMapObDataBuffer.SetKV(symbol, buffer[i:])
				empty.SetErr(fmt.Errorf("pu != lastId, %v %v", _id, firstId))
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
		w.svDataCacheBySyb.SetKV(symbol, []WsDepthMsg{})
	} else {
		w.svDataCacheBySyb.SetKV(symbol, buffer[lastIndex+1:])
	}
	return obData
}
