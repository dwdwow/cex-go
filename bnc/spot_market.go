package bnc

import (
	"errors"
	"strconv"
)

type RawOrderBook struct {
	// Common
	LastUpdateID int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`

	// Futures
	EventTime       int64 `json:"E"`
	TransactionTime int64 `json:"T"`

	// CM Futures
	Symbol string `json:"symbol"`
	Pair   string `json:"pair"`
}

type OrderBook struct {
	// Common
	LastUpdateID int64       `json:"lastUpdateId"`
	Bids         [][]float64 `json:"bids"`
	Asks         [][]float64 `json:"asks"`

	// Futures
	EventTime       int64 `json:"E"`
	TransactionTime int64 `json:"T"`

	// CM Futures
	Symbol string `json:"symbol"`
	Pair   string `json:"pair"`
}

type ParamsOrderBook struct {
	Symbol string `json:"symbol"`
	Limit  int64  `json:"limit,omitempty"` // for spot, default 100, range [100, 5000]; for futures, default 500, range 5, 10, 20, 50, [100, 1000]
}

func GetSpotRawOrderBook(params ParamsOrderBook) (orderBook RawOrderBook, err error) {
	req := Req[ParamsOrderBook]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/depth",
		Params:  params,
	}
	resp, err := Request[ParamsOrderBook, RawOrderBook](req)
	return resp.Data, err
}

func GetSpotOrderBook(params ParamsOrderBook) (orderBook OrderBook, err error) {
	rawOrderBook, err := GetSpotRawOrderBook(params)
	if err != nil {
		return
	}
	orderBook.LastUpdateID = rawOrderBook.LastUpdateID
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

type Trade struct {
	ID           int64  `json:"id"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
	// Spot
	IsBestMatch bool `json:"isBestMatch"`
}

type ParamsTrades struct {
	Symbol string `json:"symbol"`
	Limit  int64  `json:"limit,omitempty"` // default 500, range [500, 1000]
}

func GetSpotTrades(params ParamsTrades) (trades []Trade, err error) {
	req := Req[ParamsTrades]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/trades",
		Params:  params,
	}
	resp, err := Request[ParamsTrades, []Trade](req)
	return resp.Data, err
}

type ParamsHistoricalTrades struct {
	Symbol string `json:"symbol"`
	Limit  int64  `json:"limit,omitempty"` // for spot, default 500, range [500, 1000]; for futures, default 100, range [100, 500]
	FromID int64  `json:"fromId,omitempty"`
}

func GetSpotHistoricalTrades(params ParamsHistoricalTrades) (trades []Trade, err error) {
	req := Req[ParamsHistoricalTrades]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/historicalTrades",
		Params:  params,
	}
	resp, err := Request[ParamsHistoricalTrades, []Trade](req)
	return resp.Data, err
}

type AggTrade struct {
	ID           int64  `json:"a"`
	Price        string `json:"p"`
	Qty          string `json:"q"`
	FirstTradeID int64  `json:"f"`
	LastTradeID  int64  `json:"l"`
	Time         int64  `json:"T"`
	IsBuyerMaker bool   `json:"m"`
	IsBestMatch  bool   `json:"M"`
}

type ParamsAggTrades struct {
	Symbol    string `json:"symbol"`              // Trading symbol, e.g. BTCUSDT
	FromID    int64  `json:"fromId,omitempty"`    // ID, aggTradeId, to get aggregate trades from INCLUSIVE
	StartTime int64  `json:"startTime,omitempty"` // Timestamp in ms to get aggregate trades from INCLUSIVE
	EndTime   int64  `json:"endTime,omitempty"`   // Timestamp in ms to get aggregate trades until INCLUSIVE
	Limit     int64  `json:"limit,omitempty"`     // Default: 500, Maximum: 1000
}

func GetSpotAggTrades(params ParamsAggTrades) (trades []AggTrade, err error) {
	req := Req[ParamsAggTrades]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/aggTrades",
		Params:  params,
	}
	resp, err := Request[ParamsAggTrades, []AggTrade](req)
	return resp.Data, err
}

// RawKline is the raw kline data from the API
// [
//
//	  1499040000000,      // Kline open time
//	  "0.01634790",       // Open price
//	  "0.80000000",       // High price
//	  "0.01575800",       // Low price
//	  "0.01577100",       // Close price
//	  "148976.11427815",  // Volume
//	  1499644799999,      // Kline Close time
//	  "2434.19055334",    // Quote asset volume
//	  308,                // Number of trades
//	  "1756.87402397",    // Taker buy base asset volume
//	  "28.46694368",      // Taker buy quote asset volume
//	  "0"                 // Unused field, ignore.
//	]
type RawKline [12]any

func (rk RawKline) ToKline() (k Kline) {
	k[0] = rk[0].(float64)
	var err error
	k[1], err = strconv.ParseFloat(rk[1].(string), 64)
	if err != nil {
		panic(err)
	}
	k[2], err = strconv.ParseFloat(rk[2].(string), 64)
	if err != nil {
		panic(err)
	}
	k[3], err = strconv.ParseFloat(rk[3].(string), 64)
	if err != nil {
		panic(err)
	}
	k[4], err = strconv.ParseFloat(rk[4].(string), 64)

	if err != nil {
		panic(err)
	}
	k[5], err = strconv.ParseFloat(rk[5].(string), 64)
	if err != nil {
		panic(err)
	}
	k[6] = rk[6].(float64)
	k[7], err = strconv.ParseFloat(rk[7].(string), 64)
	if err != nil {
		panic(err)
	}
	k[8] = rk[8].(float64)
	k[9], err = strconv.ParseFloat(rk[9].(string), 64)
	if err != nil {
		panic(err)
	}
	k[10], err = strconv.ParseFloat(rk[10].(string), 64)
	if err != nil {
		panic(err)
	}
	return
}

type Kline [11]float64

func (k Kline) OpenTime() int64 {
	return int64(k[0])
}

func (k Kline) OpenPrice() float64 {
	return k[1]
}

func (k Kline) HighPrice() float64 {
	return k[2]
}

func (k Kline) LowPrice() float64 {
	return k[3]
}

func (k Kline) ClosePrice() float64 {
	return k[4]
}

func (k Kline) Volume() float64 {
	return k[5]
}

func (k Kline) CloseTime() int64 {
	return int64(k[6])
}

func (k Kline) QuoteAssetVolume() float64 {
	return k[7]
}

func (k Kline) NumberOfTrades() int64 {
	return int64(k[8])
}

func (k Kline) TakerBuyBaseAssetVolume() float64 {
	return k[9]
}

func (k Kline) TakerBuyQuoteAssetVolume() float64 {
	return k[10]
}

// ParamsKlines is for spot, um, cm
// For spot, symbol is required
// For um continuous kline, pair is required
// For cm continuous kline, symbol and pair are required
type ParamsKlines struct {
	Symbol    string        `json:"symbol,omitempty"`
	Interval  KlineInterval `json:"interval"`
	StartTime int64         `json:"startTime,omitempty"`
	EndTime   int64         `json:"endTime,omitempty"`
	TimeZone  string        `json:"timeZone,omitempty"`
	Limit     int64         `json:"limit,omitempty"` // default 500, range [500, 1000], for continuous contract kline, max is 1500

	// Continuous Contract, Index Price Kline
	Pair         string       `json:"pair,omitempty"`
	ContractType ContractType `json:"contractType,omitempty"` // PERPETUAL, CURRENT_QUARTER, NEXT_QUARTER
}

func GetSpotRawKlines(params ParamsKlines) (klines []RawKline, err error) {
	req := Req[ParamsKlines]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/klines",
		Params:  params,
	}
	resp, err := Request[ParamsKlines, []RawKline](req)
	return resp.Data, err
}

func GetSpotKlines(params ParamsKlines) (klines []Kline, err error) {
	rawKlines, err := GetSpotRawKlines(params)
	if err != nil {
		return
	}
	for _, rawKline := range rawKlines {
		klines = append(klines, rawKline.ToKline())
	}
	return
}

type AvgPrice struct {
	Mins      int64   `json:"mins"`
	Price     float64 `json:"price,string"`
	CloseTime int64   `json:"closeTime"`
}

type ParamsAvgPrice struct {
	Symbol string `json:"symbol"`
}

func GetSpotAvgPrice(params ParamsAvgPrice) (avgPrice AvgPrice, err error) {
	req := Req[ParamsAvgPrice]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/avgPrice",
		Params:  params,
	}
	res, err := Request[ParamsAvgPrice, AvgPrice](req)
	return res.Data, err
}

type Ticker24HrStats struct {
	// MINI FULL
	Symbol      string  `json:"symbol"`
	LastPrice   float64 `json:"lastPrice,string"`
	OpenPrice   float64 `json:"openPrice,string"`
	HighPrice   float64 `json:"highPrice,string"`
	LowPrice    float64 `json:"lowPrice,string"`
	Volume      float64 `json:"volume,string"`
	QuoteVolume float64 `json:"quoteVolume,string"`
	OpenTime    int64   `json:"openTime"`
	CloseTime   int64   `json:"closeTime"`
	FirstID     int64   `json:"firstId"`
	LastID      int64   `json:"lastId"`
	Count       int64   `json:"count"`
	// FULL
	PriceChange        float64 `json:"priceChange,string"`
	PriceChangePercent float64 `json:"priceChangePercent,string"`
	WeightedAvgPrice   float64 `json:"weightedAvgPrice,string"`
	PrevClosePrice     float64 `json:"prevClosePrice,string"`
	LastQty            float64 `json:"lastQty,string"`
	// FULL, SPOT
	BidPrice float64 `json:"bidPrice,string"`
	BidQty   float64 `json:"bidQty,string"`
	AskPrice float64 `json:"askPrice,string"`
	AskQty   float64 `json:"askQty,string"`
}

type ParamsTicker24HrStats struct {
	Symbol string `json:"symbol"`
	// SPOT
	Type string `json:"type,omitempty"` // FULL or MINI, default is FULL
}

func GetSpotTicker24HrStats(params ParamsTicker24HrStats) (ticker Ticker24HrStats, err error) {
	if params.Symbol == "" {
		return Ticker24HrStats{}, errors.New("bnc: GetSpotTicker24Hhr, symbol is required")
	}
	req := Req[ParamsTicker24HrStats]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/ticker/24hr",
		Params:  params,
	}
	res, err := Request[ParamsTicker24HrStats, Ticker24HrStats](req)
	return res.Data, err
}

type ParamsTicker24HrStatsList struct {
	Symbols []string `json:"symbols,omitempty"`
	Type    string   `json:"type,omitempty"` // FULL or MINI, default is FULL
}

func GetSpotTicker24HrStatsList(params ParamsTicker24HrStatsList) (tickers []Ticker24HrStats, err error) {
	req := Req[ParamsTicker24HrStatsList]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/ticker/24hr",
		Params:  params,
	}
	res, err := Request[ParamsTicker24HrStatsList, []Ticker24HrStats](req)
	return res.Data, err
}

type TickerTradingDayStats struct {
	Symbol             string  `json:"symbol"`
	PriceChange        float64 `json:"priceChange,string"`
	PriceChangePercent float64 `json:"priceChangePercent,string"`
	WeightedAvgPrice   float64 `json:"weightedAvgPrice,string"`
	OpenPrice          float64 `json:"openPrice,string"`
	HighPrice          float64 `json:"highPrice,string"`
	LowPrice           float64 `json:"lowPrice,string"`
	LastPrice          float64 `json:"lastPrice,string"`
	Volume             float64 `json:"volume,string"`
	QuoteVolume        float64 `json:"quoteVolume,string"`
	OpenTime           int64   `json:"openTime"`
	CloseTime          int64   `json:"closeTime"`
	FirstID            int64   `json:"firstId"`
	LastID             int64   `json:"lastId"`
	Count              int64   `json:"count"`
}

type ParamsTickerTradingDayStats struct {
	Symbol   string `json:"symbol"`
	TimeZone string `json:"timeZone,omitempty"`
	Type     string `json:"type,omitempty"` // FULL or MINI, default is FULL
}

func GetSpotTickerTradingDayStats(params ParamsTickerTradingDayStats) (ticker TickerTradingDayStats, err error) {
	if params.Symbol == "" {
		return TickerTradingDayStats{}, errors.New("bnc: GetSpotTickerTradingDayStats, symbol is required")
	}
	req := Req[ParamsTickerTradingDayStats]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/ticker/tradingDay",
		Params:  params,
	}
	res, err := Request[ParamsTickerTradingDayStats, TickerTradingDayStats](req)
	return res.Data, err
}

type ParamsTickerTradingDayStatsList struct {
	Symbols  []string `json:"symbols"` // max 100 symbols
	TimeZone string   `json:"timeZone,omitempty"`
	Type     string   `json:"type,omitempty"` // FULL or MINI, default is FULL
}

func GetSpotTickerTradingDayStatsList(params ParamsTickerTradingDayStatsList) (tickers []TickerTradingDayStats, err error) {
	req := Req[ParamsTickerTradingDayStatsList]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/ticker/tradingDay",
		Params:  params,
	}
	res, err := Request[ParamsTickerTradingDayStatsList, []TickerTradingDayStats](req)
	return res.Data, err
}

type TickerPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
	// Futures
	Time int64 `json:"time"`
	// CM Futures
	Pair string `json:"pair"`
}

type ParamsTickerPrice struct {
	Symbol string `json:"symbol"`
}

func GetSpotTickerPrice(params ParamsTickerPrice) (ticker TickerPrice, err error) {
	if params.Symbol == "" {
		return TickerPrice{}, errors.New("bnc: GetSpotTickerPrice, symbol is required")
	}
	req := Req[ParamsTickerPrice]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/ticker/price",
		Params:  params,
	}
	resp, err := Request[ParamsTickerPrice, TickerPrice](req)
	return resp.Data, err
}

type ParamsTickerPriceList struct {
	Symbols []string `json:"symbols,omitempty"`
}

func GetSpotTickerPriceList(params ParamsTickerPriceList) (tickers []TickerPrice, err error) {
	req := Req[ParamsTickerPriceList]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/ticker/price",
		Params:  params,
	}
	resp, err := Request[ParamsTickerPriceList, []TickerPrice](req)
	return resp.Data, err
}

type OrderBookTicker struct {
	Symbol   string  `json:"symbol"`
	BidPrice float64 `json:"bidPrice,string"`
	BidQty   float64 `json:"bidQty,string"`
	AskPrice float64 `json:"askPrice,string"`
	AskQty   float64 `json:"askQty,string"`
	// Futures
	Time int64 `json:"time"`
}

type ParamsOrderBookTicker struct {
	Symbol string `json:"symbol"`
}

func GetSpotOrderBookTicker(params ParamsOrderBookTicker) (ticker OrderBookTicker, err error) {
	if params.Symbol == "" {
		return OrderBookTicker{}, errors.New("bnc: GetSpotTickerBookTicker, symbol is required")
	}
	req := Req[ParamsOrderBookTicker]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/ticker/bookTicker",
		Params:  params,
	}
	resp, err := Request[ParamsOrderBookTicker, OrderBookTicker](req)
	return resp.Data, err
}

type ParamsTickerBookTickerList struct {
	Symbols []string `json:"symbols,omitempty"`
}

func GetSpotTickerBookTickerList(params ParamsTickerBookTickerList) (tickers []OrderBookTicker, err error) {
	req := Req[ParamsTickerBookTickerList]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/ticker/bookTicker",
		Params:  params,
	}
	resp, err := Request[ParamsTickerBookTickerList, []OrderBookTicker](req)
	return resp.Data, err
}

type TickerStats struct {
	Symbol             string  `json:"symbol"`
	PriceChange        float64 `json:"priceChange,string"`
	PriceChangePercent float64 `json:"priceChangePercent,string"`
	WeightedAvgPrice   float64 `json:"weightedAvgPrice,string"`
	OpenPrice          float64 `json:"openPrice,string"`
	HighPrice          float64 `json:"highPrice,string"`
	LowPrice           float64 `json:"lowPrice,string"`
	LastPrice          float64 `json:"lastPrice,string"`
	Volume             float64 `json:"volume,string"`
	QuoteVolume        float64 `json:"quoteVolume,string"`
	OpenTime           int64   `json:"openTime"`
	CloseTime          int64   `json:"closeTime"`
	FirstID            int64   `json:"firstId"`
	LastID             int64   `json:"lastId"`
	Count              int64   `json:"count"`
}

type ParamsTickerStats struct {
	Symbol     string `json:"symbol"`
	WindowSize string `json:"windowSize,omitempty"` // 1m-59m, 1h-23h, 1d-7d, default 1d
	Type       string `json:"type,omitempty"`       // FULL or MINI, default FULL
}

func GetSpotTickerStats(params ParamsTickerStats) (ticker TickerStats, err error) {
	if params.Symbol == "" {
		return TickerStats{}, errors.New("bnc: GetSpotTickerStats, symbol is required")
	}
	req := Req[ParamsTickerStats]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/ticker",
		Params:  params,
	}
	resp, err := Request[ParamsTickerStats, TickerStats](req)
	return resp.Data, err
}

type ParamsTickerStatsList struct {
	Symbols    []string `json:"symbols"`              // max 100 symbols
	WindowSize string   `json:"windowSize,omitempty"` // 1m-59m, 1h-23h, 1d-7d, default 1d
	Type       string   `json:"type,omitempty"`       // FULL or MINI, default FULL
}

func GetSpotTickerStatsList(params ParamsTickerStatsList) (tickers []TickerStats, err error) {
	req := Req[ParamsTickerStatsList]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/ticker",
		Params:  params,
	}
	resp, err := Request[ParamsTickerStatsList, []TickerStats](req)
	return resp.Data, err
}
