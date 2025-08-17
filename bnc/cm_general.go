package bnc

import "github.com/dwdwow/cex-go"

func CheckCMFuturesServerTime() (serverTime ServerTime, err error) {
	req := Req[EmptyStruct]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/time",
	}
	resp, err := Request[EmptyStruct, ServerTime](req)
	return resp.Data, err
}

func GetCMExchangeInfo() (exchangeInfo FuturesExchangeInfo, err error) {
	req := Req[EmptyStruct]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/exchangeInfo",
	}
	resp, err := Request[EmptyStruct, FuturesExchangeInfo](req)
	return resp.Data, err
}

func GetCMPairs() (pairs []cex.Pair, err error) {
	exchangeInfo, err := GetCMExchangeInfo()
	if err != nil {
		return nil, err
	}
	for _, symbol := range exchangeInfo.Symbols {
		pair, err := symbol.ToPair()
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, pair)
	}
	return pairs, nil
}
