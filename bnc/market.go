package bnc

import (
	"errors"

	"github.com/dwdwow/cex-go"
)

func GetOrderbook(symbolType cex.SymbolType, params ParamsOrderBook) (orderbook OrderBook, err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return GetSpotOrderBook(params)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return GetUMOrderBook(params)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return GetCMOrderBook(params)
	default:
		return orderbook, errors.New("bnc: invalid symbol type")
	}
}
