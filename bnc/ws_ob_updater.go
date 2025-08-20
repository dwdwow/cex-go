package bnc

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/cex-go/ob"
	"github.com/dwdwow/props"
)

type obWsCacher interface {
	symbolType() cex.SymbolType
	cache() *props.SafeRWMap[string, []WsDepthMsg]
	exist() *props.SafeRWMap[string, bool]
	ods() *props.SafeRWMap[string, ob.Data[WsDepthMsg]]
}

type obUpdater func(w obWsCacher, depthData WsDepthMsg) ob.Data[WsDepthMsg]

func defaultObUpdater(symbolType cex.SymbolType) obUpdater {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return spObUpdater
	case cex.SYMBOL_TYPE_UM_FUTURES, cex.SYMBOL_TYPE_CM_FUTURES:
		return futuresObUpdater
	}
	panic("bnc: unknown symbol type")
}

func spObUpdater(w obWsCacher, depthData WsDepthMsg) ob.Data[WsDepthMsg] {
	symbol := depthData.Symbol
	buffer := w.cache().GetV(symbol)
	empty := ob.NewData[WsDepthMsg](cex.BINANCE, w.symbolType(), symbol)
	obData, ok := w.ods().GetVWithOk(symbol)
	if !ok || obData.Empty() {
		empty.SetErr(errors.New("bnc: unexpected error: binance update ob: if !ok || obData.Empty()"))
		return empty
	}
	currentVersion, err := strconv.ParseInt(obData.Version, 10, 64)
	if err != nil {
		empty.SetErr(fmt.Errorf("bnc: can not parse ob data version %s, err: %w", obData.Version, err))
		return empty
	}
	if buffer[0].FirstId > currentVersion+1 {
		empty.SetErr(errors.New("bnc: current ob version is small"))
		return empty
	}
	lastIndex := 0
	_id := int64(0)
	for i, _depthData := range buffer {
		firstId := _depthData.FirstId
		lastId := _depthData.LastId
		if _id > 0 {
			if firstId == _id+1 {
				_id = lastId
			} else {
				w.cache().SetKV(symbol, buffer[i:])
				empty.SetErr(errors.New("bnc: ob depth update is not continuous, websocket message may be blocked"))
				return empty
			}
		} else {
			_id = lastId
		}
		if firstId != currentVersion+1 {
			if firstId <= currentVersion+1 && lastId >= currentVersion+1 {
				// update firstly
			} else {
				if firstId > currentVersion+1 {
					lastIndex = i - 1
					break
				}
				if firstId < currentVersion+1 {
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
		w.cache().SetKV(symbol, []WsDepthMsg{})
	} else {
		w.cache().SetKV(symbol, buffer[lastIndex+1:])
	}
	obData.Time = depthData.EventTime
	obData.Note = depthData
	return obData
}

func futuresObUpdater(w obWsCacher, depthData WsDepthMsg) ob.Data[WsDepthMsg] {
	symbol := depthData.Symbol
	buffer := w.cache().GetV(symbol)
	empty := ob.NewData[WsDepthMsg](cex.BINANCE, cex.SYMBOL_TYPE_UM_FUTURES, symbol)
	obData, ok := w.ods().GetVWithOk(symbol)
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
				w.cache().SetKV(symbol, buffer[i:])
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
		w.cache().SetKV(symbol, []WsDepthMsg{})
	} else {
		w.cache().SetKV(symbol, buffer[lastIndex+1:])
	}
	obData.Time = depthData.EventTime
	obData.Note = depthData
	return obData
}
