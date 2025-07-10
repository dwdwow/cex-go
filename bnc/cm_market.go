package bnc

import (
	"errors"
	"strconv"
)

func GetCMRawOrderBook(params ParamsOrderBook) (orderBook RawOrderBook, err error) {
	if params.Symbol == "" {
		return RawOrderBook{}, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsOrderBook]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/depth",
		Params:  params,
	}
	resp, err := Request[ParamsOrderBook, RawOrderBook](req)
	return resp.Data, err
}

func GetCMOrderBook(params ParamsOrderBook) (orderBook OrderBook, err error) {
	rawOrderBook, err := GetCMRawOrderBook(params)
	if err != nil {
		return OrderBook{}, err
	}
	orderBook.LastUpdateID = rawOrderBook.LastUpdateID
	orderBook.Symbol = rawOrderBook.Symbol
	orderBook.Pair = rawOrderBook.Pair
	orderBook.EventTime = rawOrderBook.EventTime
	orderBook.TransactionTime = rawOrderBook.TransactionTime
	orderBook.Bids = make([][]float64, len(rawOrderBook.Bids))
	orderBook.Asks = make([][]float64, len(rawOrderBook.Asks))
	for i, bid := range rawOrderBook.Bids {
		orderBook.Bids[i] = make([]float64, len(bid))
		for j, s := range bid {
			orderBook.Bids[i][j], err = strconv.ParseFloat(s, 64)
			if err != nil {
				return
			}
		}
	}
	for i, ask := range rawOrderBook.Asks {
		orderBook.Asks[i] = make([]float64, len(ask))
		for j, s := range ask {
			orderBook.Asks[i][j], err = strconv.ParseFloat(s, 64)
			if err != nil {
				return
			}
		}
	}
	return
}

func GetCMTrades(params ParamsTrades) (trades []Trade, err error) {
	if params.Symbol == "" {
		return nil, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsTrades]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/trades",
		Params:  params,
	}
	resp, err := Request[ParamsTrades, []Trade](req)
	return resp.Data, err
}

// GetCMHistoricalTrades deprecated
// should use API KEY
func GetCMHistoricalTrades(params ParamsHistoricalTrades) (trades []Trade, err error) {
	if params.Symbol == "" {
		return nil, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsHistoricalTrades]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/historicalTrades",
		Params:  params,
	}
	resp, err := Request[ParamsHistoricalTrades, []Trade](req)
	return resp.Data, err
}

type CMPremiumIndexInfo struct {
	Symbol          string  `json:"symbol"`                      // Symbol name, e.g. "BTCUSD_PERP"
	Pair            string  `json:"pair"`                        // Trading pair, e.g. "BTCUSD"
	MarkPrice       float64 `json:"markPrice,string"`            // Mark price
	IndexPrice      float64 `json:"indexPrice,string"`           // Index price
	EstSettlePrice  float64 `json:"estimatedSettlePrice,string"` // Estimated Settle Price, only useful in the last hour before settlement
	LastFundingRate float64 `json:"lastFundingRate,string"`      // Latest funding rate (perpetual contracts only, empty for delivery)
	InterestRate    float64 `json:"interestRate,string"`         // Base asset interest rate (perpetual contracts only, empty for delivery)
	NextFundingTime int64   `json:"nextFundingTime"`             // Next funding time in ms (perpetual contracts only, 0 for delivery)
	Time            int64   `json:"time"`                        // Current timestamp in ms
}

type ParamsPremiumIndexInfo struct {
	Symbol string `json:"symbol"`
	Pair   string `json:"pair"`
}

func GetCMPremiumIndexInfo(params ParamsPremiumIndexInfo) (premiumIndexInfo []CMPremiumIndexInfo, err error) {
	req := Req[ParamsPremiumIndexInfo]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/premiumIndex",
		Params:  params,
	}
	resp, err := Request[ParamsPremiumIndexInfo, []CMPremiumIndexInfo](req)
	return resp.Data, err
}

func GetCMFundingRateHistory(params ParamsFundingRate) (fundingRate []FundingRate, err error) {
	if params.Symbol == "" {
		return nil, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsFundingRate]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/fundingRate",
		Params:  params,
	}
	resp, err := Request[ParamsFundingRate, []FundingRate](req)
	return resp.Data, err
}

func GetCMFundingInfoList() (fundingInfo []FundingInfo, err error) {
	req := Req[EmptyStruct]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/fundingInfo",
	}
	resp, err := Request[EmptyStruct, []FundingInfo](req)
	return resp.Data, err
}

func GetCMRawKlines(params ParamsKlines) (klines []RawKline, err error) {
	if params.Symbol == "" {
		return nil, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsKlines]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/klines",
		Params:  params,
	}
	resp, err := Request[ParamsKlines, []RawKline](req)
	return resp.Data, err
}

func GetCMKlines(params ParamsKlines) (klines []Kline, err error) {
	rawKlines, err := GetCMRawKlines(params)
	if err != nil {
		return nil, err
	}
	klines = make([]Kline, len(rawKlines))
	for i, r := range rawKlines {
		klines[i] = r.ToKline()
	}
	return
}

func GetCMRawContinuousKlines(params ParamsKlines) (klines []RawKline, err error) {
	if params.Symbol == "" {
		return nil, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsKlines]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/continuousKlines",
		Params:  params,
	}
	resp, err := Request[ParamsKlines, []RawKline](req)
	return resp.Data, err
}

func GetCMContinousKlines(params ParamsKlines) (klines []Kline, err error) {
	rawKlines, err := GetCMRawContinuousKlines(params)
	if err != nil {
		return nil, err
	}
	klines = make([]Kline, len(rawKlines))
	for i, r := range rawKlines {
		klines[i] = r.ToKline()
	}
	return
}

func GetCMRawIndexPriceKlines(params ParamsKlines) (klines []RawKline, err error) {
	if params.Symbol == "" {
		return nil, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsKlines]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/indexPriceKlines",
		Params:  params,
	}
	resp, err := Request[ParamsKlines, []RawKline](req)
	return resp.Data, err
}

func GetCMIndexPriceKlines(params ParamsKlines) (klines []Kline, err error) {
	rawKlines, err := GetCMRawIndexPriceKlines(params)
	if err != nil {
		return nil, err
	}
	klines = make([]Kline, len(rawKlines))
	for i, r := range rawKlines {
		klines[i] = r.ToKline()
	}
	return
}

func GetCMRawMarkPriceKlines(params ParamsKlines) (klines []RawKline, err error) {
	if params.Symbol == "" {
		return nil, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsKlines]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/markPriceKlines",
		Params:  params,
	}
	resp, err := Request[ParamsKlines, []RawKline](req)
	return resp.Data, err
}

func GetCMMarkPriceKlines(params ParamsKlines) (klines []Kline, err error) {
	rawKlines, err := GetCMRawMarkPriceKlines(params)
	if err != nil {
		return nil, err
	}
	klines = make([]Kline, len(rawKlines))
	for i, r := range rawKlines {
		klines[i] = r.ToKline()
	}
	return
}

func GetCMRawPremiumIndexKlines(params ParamsKlines) (klines []RawKline, err error) {
	if params.Symbol == "" {
		return nil, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsKlines]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/premiumIndexKlines",
		Params:  params,
	}
	resp, err := Request[ParamsKlines, []RawKline](req)
	return resp.Data, err
}

func GetCMPremiumIndexKlines(params ParamsKlines) (klines []Kline, err error) {
	rawKlines, err := GetCMRawPremiumIndexKlines(params)
	if err != nil {
		return nil, err
	}
	klines = make([]Kline, len(rawKlines))
	for i, r := range rawKlines {
		klines[i] = r.ToKline()
	}
	return
}

func GetCMTicker24HrStats(params ParamsTicker24HrStats) (ticker []Ticker24HrStats, err error) {
	req := Req[ParamsTicker24HrStats]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/ticker/24hr",
		Params:  params,
	}
	resp, err := Request[ParamsTicker24HrStats, []Ticker24HrStats](req)
	return resp.Data, err
}

func GetCMTicker24HrStatsList() (ticker []Ticker24HrStats, err error) {
	req := Req[ParamsTicker24HrStatsList]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/ticker/24hr",
	}
	resp, err := Request[ParamsTicker24HrStatsList, []Ticker24HrStats](req)
	return resp.Data, err
}

func GetCMTickerPrice(params ParamsTickerPrice) (ticker []TickerPrice, err error) {
	req := Req[ParamsTickerPrice]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/ticker/price",
		Params:  params,
	}
	resp, err := Request[ParamsTickerPrice, []TickerPrice](req)
	return resp.Data, err
}

func GetCMTickerPriceList() (ticker []TickerPrice, err error) {
	req := Req[ParamsTickerPriceList]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/ticker/price",
	}
	resp, err := Request[ParamsTickerPriceList, []TickerPrice](req)
	return resp.Data, err
}

func GetCMOrderBookTicker(params ParamsOrderBookTicker) (ticker []OrderBookTicker, err error) {
	req := Req[ParamsOrderBookTicker]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/ticker/bookTicker",
		Params:  params,
	}
	resp, err := Request[ParamsOrderBookTicker, []OrderBookTicker](req)
	return resp.Data, err
}

func GetCMOpenInterest(params ParamsOpenInterest) (openInterest OpenInterest, err error) {
	if params.Symbol == "" {
		return OpenInterest{}, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsOpenInterest]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/openInterest",
		Params:  params,
	}
	resp, err := Request[ParamsOpenInterest, OpenInterest](req)
	return resp.Data, err
}

type CMOpenInterestStats struct {
	Pair               string       `json:"pair"`
	ContractType       ContractType `json:"contractType"`
	SumOpenInterest    float64      `json:"sumOpenInterest,string"`
	SumOpenInterestVal float64      `json:"sumOpenInterestValue,string"`
	Timestamp          int64        `json:"timestamp"`
}

type ParamsCMOpenInterestStats struct {
	Pair         string       `json:"pair"`
	ContractType ContractType `json:"contractType"`
	Period       string       `json:"period"` // "5m","15m","30m","1h","2h","4h","6h","12h","1d"
	Limit        int64        `json:"limit,omitempty"`
	EndTime      int64        `json:"endTime,omitempty"`
	StartTime    int64        `json:"startTime,omitempty"`
}

func GetCMOpenInterestStats(params ParamsCMOpenInterestStats) (openInterestStats []CMOpenInterestStats, err error) {
	req := Req[ParamsCMOpenInterestStats]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    FUTURES_DATA + "/openInterestHist",
		Params:  params,
	}
	resp, err := Request[ParamsCMOpenInterestStats, []CMOpenInterestStats](req)
	return resp.Data, err
}

type CMLongShortRatio struct {
	Pair           string  `json:"pair"`                  // e.g. "BTCUSD"
	LongShortRatio float64 `json:"longShortRatio,string"` // e.g. 0.7869
	LongPosition   float64 `json:"longPosition,string"`   // Long account ratio e.g. 0.6442 (64.42%)
	ShortPosition  float64 `json:"shortPosition,string"`  // Short account ratio e.g. 0.4404 (44.04%)
	Timestamp      int64   `json:"timestamp"`             // Unix timestamp in milliseconds
}

type ParamsCMLongShortRatio struct {
	Pair      string `json:"pair"`
	Period    string `json:"period"` // "5m","15m","30m","1h","2h","4h","6h","12h","1d"
	Limit     int64  `json:"limit,omitempty"`
	EndTime   int64  `json:"endTime,omitempty"`
	StartTime int64  `json:"startTime,omitempty"`
}

func GetCMTopTraderLongShortPositionRatio(params ParamsCMLongShortRatio) (longShortRatio []CMLongShortRatio, err error) {
	req := Req[ParamsCMLongShortRatio]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    FUTURES_DATA + "/topLongShortPositionRatio",
		Params:  params,
	}
	resp, err := Request[ParamsCMLongShortRatio, []CMLongShortRatio](req)
	return resp.Data, err
}

func GetCMTopTraderLongShortAccountRatio(params ParamsCMLongShortRatio) (longShortRatio []CMLongShortRatio, err error) {
	req := Req[ParamsCMLongShortRatio]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    FUTURES_DATA + "/topLongShortAccountRatio",
		Params:  params,
	}
	resp, err := Request[ParamsCMLongShortRatio, []CMLongShortRatio](req)
	return resp.Data, err
}

func GetCMGlobalLongShortAccountRatio(params ParamsCMLongShortRatio) (longShortRatio []CMLongShortRatio, err error) {
	req := Req[ParamsCMLongShortRatio]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    FUTURES_DATA + "/globalLongShortAccountRatio",
		Params:  params,
	}
	resp, err := Request[ParamsCMLongShortRatio, []CMLongShortRatio](req)
	return resp.Data, err
}

type CMTakerBuySellVolume struct {
	Pair              string       `json:"pair"`
	ContractType      ContractType `json:"contractType"`
	TakerBuyVol       int64        `json:"takerBuyVol,string"`
	TakerSellVol      int64        `json:"takerSellVol,string"`
	TakerBuyVolValue  float64      `json:"takerBuyVolValue,string"`
	TakerSellVolValue float64      `json:"takerSellVolValue,string"`
	Timestamp         int64        `json:"timestamp"`
}

type ParamsCMBuySellVolume struct {
	Pair         string       `json:"pair"`            // e.g. BTCUSD
	ContractType ContractType `json:"contractType"`    // ALL, CURRENT_QUARTER, NEXT_QUARTER, PERPETUAL
	Period       string       `json:"period"`          // "5m","15m","30m","1h","2h","4h","6h","12h","1d"
	Limit        int64        `json:"limit,omitempty"` // Default 30, Max 500
	EndTime      int64        `json:"endTime,omitempty"`
	StartTime    int64        `json:"startTime,omitempty"`
}

func GetCMTakerBuySellVolume(params ParamsCMBuySellVolume) (longShortRatio []CMTakerBuySellVolume, err error) {
	req := Req[ParamsCMBuySellVolume]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    FUTURES_DATA + "/takerBuySellVol",
		Params:  params,
	}
	resp, err := Request[ParamsCMBuySellVolume, []CMTakerBuySellVolume](req)
	return resp.Data, err
}

type CMBasisInfo struct {
	Pair                string       `json:"pair"`
	ContractType        ContractType `json:"contractType"`
	FuturesPrice        float64      `json:"futuresPrice,string"`
	IndexPrice          float64      `json:"indexPrice,string"`
	Basis               float64      `json:"basis,string"`
	BasisRate           float64      `json:"basisRate,string"`
	AnnualizedBasisRate string       `json:"annualizedBasisRate"`
	Timestamp           int64        `json:"timestamp"`
}

type ParamsCMBasisInfo ParamsCMBuySellVolume

func GetCMBasisInfoList(params ParamsCMBasisInfo) (longShortRatio []CMBasisInfo, err error) {
	req := Req[ParamsCMBasisInfo]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    FUTURES_DATA + "/basis",
		Params:  params,
	}
	resp, err := Request[ParamsCMBasisInfo, []CMBasisInfo](req)
	return resp.Data, err
}

func GetCMIndexPriceConstituentInfo(params ParamsIndexPriceConstituentInfo) (indexPriceConstituentInfo IndexPriceConstituents, err error) {
	req := Req[ParamsIndexPriceConstituentInfo]{
		BaseURL: API_CM_FUTURES_ENDPOINT,
		Path:    DAPI_V1 + "/constituents",
		Params:  params,
	}
	resp, err := Request[ParamsIndexPriceConstituentInfo, IndexPriceConstituents](req)
	return resp.Data, err
}
