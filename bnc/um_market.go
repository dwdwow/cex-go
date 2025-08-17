package bnc

import (
	"errors"
	"strconv"

	"github.com/dwdwow/cex-go/ob"
)

func GetUMRawOrderBook(params ParamsOrderBook) (orderBook RawOrderBook, err error) {
	req := Req[ParamsOrderBook]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/depth",
		Params:  params,
	}
	resp, err := Request[ParamsOrderBook, RawOrderBook](req)
	return resp.Data, err
}

func GetUMOrderBook(params ParamsOrderBook) (orderBook OrderBook, err error) {
	rawOrderBook, err := GetUMRawOrderBook(params)
	if err != nil {
		return
	}
	orderBook.LastUpdateID = rawOrderBook.LastUpdateID
	orderBook.Bids = make(ob.Book, len(rawOrderBook.Bids))
	orderBook.Asks = make(ob.Book, len(rawOrderBook.Asks))
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

func GetUMTrades(params ParamsTrades) (trades []Trade, err error) {
	req := Req[ParamsTrades]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/trades",
		Params:  params,
	}
	resp, err := Request[ParamsTrades, []Trade](req)
	return resp.Data, err
}

// GetUMHistoricalTrades Deprecated
// must use api key, weird
func GetUMHistoricalTrades(params ParamsHistoricalTrades) (trades []Trade, err error) {
	req := Req[ParamsHistoricalTrades]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/historicalTrades",
		Params:  params,
	}
	resp, err := Request[ParamsHistoricalTrades, []Trade](req)
	return resp.Data, err
}

func GetUMAggTrades(params ParamsAggTrades) (trades []AggTrade, err error) {
	req := Req[ParamsAggTrades]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/aggTrades",
		Params:  params,
	}
	resp, err := Request[ParamsAggTrades, []AggTrade](req)
	return resp.Data, err
}

func GetUMRawKlines(params ParamsKlines) (klines []RawKline, err error) {
	req := Req[ParamsKlines]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/klines",
		Params:  params,
	}
	resp, err := Request[ParamsKlines, []RawKline](req)
	return resp.Data, err
}

func GetUMKlines(params ParamsKlines) (klines []Kline, err error) {
	rawKlines, err := GetUMRawKlines(params)
	if err != nil {
		return
	}
	klines = make([]Kline, len(rawKlines))
	for i, rawKline := range rawKlines {
		klines[i] = rawKline.ToKline()
	}
	return
}

func GetUMRawContinuousKlines(params ParamsKlines) (klines []RawKline, err error) {
	req := Req[ParamsKlines]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/continuousKlines",
		Params:  params,
	}
	resp, err := Request[ParamsKlines, []RawKline](req)
	return resp.Data, err
}

func GetUMContinousKlines(params ParamsKlines) (klines []Kline, err error) {
	rawKlines, err := GetUMRawContinuousKlines(params)
	if err != nil {
		return
	}
	klines = make([]Kline, len(rawKlines))
	for i, rawKline := range rawKlines {
		klines[i] = rawKline.ToKline()
	}
	return
}

// [
//     1591256400000,      	// Open time
//     "9653.69440000",    	// Open
//     "9653.69640000",     // High
//     "9651.38600000",     // Low
//     "9651.55200000",     // Close (or latest price)
//     "0", 				// Ignore
//     1591256459999,      	// Close time
//     "0",    				// Ignore,
//     60,                	// Ignore,
//     "0",    				// Ignore
//     "0",      			// Ignore
//     "0" 					// Ignore
// ]

func GetUMRawIndexPirceKlines(params ParamsKlines) (klines []RawKline, err error) {
	req := Req[ParamsKlines]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/indexPriceKlines",
		Params:  params,
	}
	resp, err := Request[ParamsKlines, []RawKline](req)
	return resp.Data, err
}

func GetUMIndexPirceKlines(params ParamsKlines) (klines []Kline, err error) {
	rawKlines, err := GetUMRawIndexPirceKlines(params)
	if err != nil {
		return
	}
	klines = make([]Kline, len(rawKlines))
	for i, rawKline := range rawKlines {
		klines[i] = rawKline.ToKline()
	}
	return
}

func GetUMRawMarkPriceKlines(params ParamsKlines) (klines []RawKline, err error) {
	req := Req[ParamsKlines]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/markPriceKlines",
		Params:  params,
	}
	resp, err := Request[ParamsKlines, []RawKline](req)
	return resp.Data, err
}

func GetUMMarkPriceKlines(params ParamsKlines) (klines []Kline, err error) {
	rawKlines, err := GetUMRawMarkPriceKlines(params)
	if err != nil {
		return
	}
	klines = make([]Kline, len(rawKlines))
	for i, rawKline := range rawKlines {
		klines[i] = rawKline.ToKline()
	}
	return
}

func GetUMRawPremiumIndexKlines(params ParamsKlines) (klines []RawKline, err error) {
	req := Req[ParamsKlines]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/premiumIndexKlines",
		Params:  params,
	}
	resp, err := Request[ParamsKlines, []RawKline](req)
	return resp.Data, err
}

func GetUMPremiumIndexKlines(params ParamsKlines) (klines []Kline, err error) {
	rawKlines, err := GetUMRawPremiumIndexKlines(params)
	if err != nil {
		return
	}
	klines = make([]Kline, len(rawKlines))
	for i, rawKline := range rawKlines {
		klines[i] = rawKline.ToKline()
	}
	return
}

type MarkPriceInfo struct {
	Symbol          string  `json:"symbol"`
	MarkPrice       float64 `json:"markPrice,string"`
	IndexPrice      float64 `json:"indexPrice,string"`
	EstSettlePrice  float64 `json:"estimatedSettlePrice,string"`
	LastFundingRate float64 `json:"lastFundingRate,string"`
	InterestRate    float64 `json:"interestRate,string"`
	NextFundingTime int64   `json:"nextFundingTime"`
	Time            int64   `json:"time"`
}

type ParamsMarkPriceInfo struct {
	Symbol string `json:"symbol"`
}

func GetUMMarkPriceInfo(params ParamsMarkPriceInfo) (markPriceInfo MarkPriceInfo, err error) {
	if params.Symbol == "" {
		return MarkPriceInfo{}, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsMarkPriceInfo]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/premiumIndex",
		Params:  params,
	}
	resp, err := Request[ParamsMarkPriceInfo, MarkPriceInfo](req)
	return resp.Data, err
}

func GetUMAllMarkPriceInfo() (markPriceInfos []MarkPriceInfo, err error) {
	req := Req[EmptyStruct]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/premiumIndex",
	}
	resp, err := Request[EmptyStruct, []MarkPriceInfo](req)
	return resp.Data, err
}

type FundingRate struct {
	Symbol      string  `json:"symbol"`
	FundingRate float64 `json:"fundingRate,string"`
	FundingTime int64   `json:"fundingTime"`
	MarkPrice   float64 `json:"markPrice,string"`
}

type ParamsFundingRate struct {
	Symbol    string `json:"symbol,omitempty"` // for um, default all symbols; for cm, is required
	StartTime int64  `json:"startTime,omitempty"`
	EndTime   int64  `json:"endTime,omitempty"`
	Limit     int64  `json:"limit,omitempty"` // default 100, range [100, 1000]
}

func GetUMFundingRateHistory(params ParamsFundingRate) (fundingRate []FundingRate, err error) {
	req := Req[ParamsFundingRate]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/fundingRate",
		Params:  params,
	}
	resp, err := Request[ParamsFundingRate, []FundingRate](req)
	return resp.Data, err
}

type FundingInfo struct {
	Symbol                   string `json:"symbol"`
	AdjustedFundingRateCap   string `json:"adjustedFundingRateCap"`
	AdjustedFundingRateFloor string `json:"adjustedFundingRateFloor"`
	FundingIntervalHours     int64  `json:"fundingIntervalHours"`
	Disclaimer               bool   `json:"disclaimer"` // ignored
}

func GetUMFundingInfoList() (fundingInfo []FundingInfo, err error) {
	req := Req[EmptyStruct]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/fundingInfo",
	}
	resp, err := Request[EmptyStruct, []FundingInfo](req)
	return resp.Data, err
}

func GetUMTicker24HrStats(params ParamsTicker24HrStats) (ticker Ticker24HrStats, err error) {
	if params.Symbol == "" {
		return Ticker24HrStats{}, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsTicker24HrStats]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/ticker/24hr",
		Params:  params,
	}
	resp, err := Request[ParamsTicker24HrStats, Ticker24HrStats](req)
	return resp.Data, err
}

func GetUMTicker24HrStatsList() (ticker []Ticker24HrStats, err error) {
	req := Req[ParamsTicker24HrStatsList]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/ticker/24hr",
	}
	resp, err := Request[ParamsTicker24HrStatsList, []Ticker24HrStats](req)
	return resp.Data, err
}

func GetUMTickerPrice(params ParamsTickerPrice) (ticker TickerPrice, err error) {
	if params.Symbol == "" {
		return TickerPrice{}, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsTickerPrice]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/ticker/price",
		Params:  params,
	}
	resp, err := Request[ParamsTickerPrice, TickerPrice](req)
	return resp.Data, err
}

func GetUMTickerPriceList() (ticker []TickerPrice, err error) {
	req := Req[ParamsTickerPriceList]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/ticker/price",
	}
	resp, err := Request[ParamsTickerPriceList, []TickerPrice](req)
	return resp.Data, err
}

func GetUMOrderBookTicker(params ParamsOrderBookTicker) (ticker OrderBookTicker, err error) {
	if params.Symbol == "" {
		return OrderBookTicker{}, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsOrderBookTicker]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/ticker/bookTicker",
		Params:  params,
	}
	resp, err := Request[ParamsOrderBookTicker, OrderBookTicker](req)
	return resp.Data, err
}

func GetUMTickerBookTickerList() (ticker []OrderBookTicker, err error) {
	req := Req[ParamsTickerBookTickerList]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/ticker/bookTicker",
	}
	resp, err := Request[ParamsTickerBookTickerList, []OrderBookTicker](req)
	return resp.Data, err
}

type DeliveryPrice struct {
	DeliveryTime  int64   `json:"deliveryTime"`
	DeliveryPrice float64 `json:"deliveryPrice"`
}

type ParamsDeliveryPrice struct {
	Pair string `json:"pair"`
}

func GetUMDeliveryPriceList(params ParamsDeliveryPrice) (deliveryPrice []DeliveryPrice, err error) {
	if params.Pair == "" {
		return nil, errors.New("bnc: pair is required")
	}
	req := Req[ParamsDeliveryPrice]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    "/futures/data/delivery-price",
		Params:  params,
	}
	resp, err := Request[ParamsDeliveryPrice, []DeliveryPrice](req)
	return resp.Data, err
}

type OpenInterest struct {
	Symbol       string  `json:"symbol"`
	OpenInterest float64 `json:"openInterest,string"`
	Time         int64   `json:"time"`
	// CM Futures
	Pair         string       `json:"pair"`
	ContractType ContractType `json:"contractType"`
}

type ParamsOpenInterest struct {
	Symbol string `json:"symbol"`
}

func GetUMOpenInterest(params ParamsOpenInterest) (openInterest OpenInterest, err error) {
	if params.Symbol == "" {
		return OpenInterest{}, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsOpenInterest]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/openInterest",
		Params:  params,
	}
	resp, err := Request[ParamsOpenInterest, OpenInterest](req)
	return resp.Data, err
}

type UMOpenInterestStats struct {
	Symbol               string  `json:"symbol"`
	SumOpenInterest      float64 `json:"sumOpenInterest,string"`
	SumOpenInterestValue float64 `json:"sumOpenInterestValue,string"`
	CMCCirculatingSupply float64 `json:"CMCCirculatingSupply,string"`
	Timestamp            int64   `json:"timestamp"`
}

type ParamsUMOpenInterestStats struct {
	Symbol    string `json:"symbol"`
	Period    string `json:"period"` // "5m","15m","30m","1h","2h","4h","6h","12h","1d"
	Limit     int64  `json:"limit,omitempty"`
	EndTime   int64  `json:"endTime,omitempty"`
	StartTime int64  `json:"startTime,omitempty"`
}

func GetUMOpenInterestStats(params ParamsUMOpenInterestStats) (openInterestStats []UMOpenInterestStats, err error) {
	if params.Symbol == "" {
		return nil, errors.New("bnc: symbol is required")
	}
	if params.Period == "" {
		return nil, errors.New("bnc: period is required")
	}
	req := Req[ParamsUMOpenInterestStats]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FUTURES_DATA + "/openInterestHist",
		Params:  params,
	}
	resp, err := Request[ParamsUMOpenInterestStats, []UMOpenInterestStats](req)
	return resp.Data, err
}

type UMLongShortRatio struct {
	Symbol         string  `json:"symbol"`                // e.g. "BTCUSDT"
	LongShortRatio float64 `json:"longShortRatio,string"` // Long/short position/account ratio of top traders
	LongAccount    float64 `json:"longAccount,string"`    // Long positions ratio of top traders
	ShortAccount   float64 `json:"shortAccount,string"`   // Short positions ratio of top traders
	Timestamp      int64   `json:"timestamp"`             // Unix timestamp in milliseconds
}

type ParamsUMLongShortRatio struct {
	Symbol    string `json:"symbol"`
	Period    string `json:"period"`          // "5m","15m","30m","1h","2h","4h","6h","12h","1d"
	Limit     int64  `json:"limit,omitempty"` // default 30, max 500
	StartTime int64  `json:"startTime,omitempty"`
	EndTime   int64  `json:"endTime,omitempty"`
}

func GetUMTopTraderLongShortPositionRatio(params ParamsUMLongShortRatio) (longShortRatio []UMLongShortRatio, err error) {
	if params.Symbol == "" {
		return nil, errors.New("bnc: symbol is required")
	}
	if params.Period == "" {
		return nil, errors.New("bnc: period is required")
	}
	req := Req[ParamsUMLongShortRatio]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FUTURES_DATA + "/topLongShortPositionRatio",
		Params:  params,
	}
	resp, err := Request[ParamsUMLongShortRatio, []UMLongShortRatio](req)
	return resp.Data, err
}

func GetUMTopTraderLongShortAccountRatio(params ParamsUMLongShortRatio) (longShortRatio []UMLongShortRatio, err error) {
	if params.Symbol == "" {
		return nil, errors.New("bnc: symbol is required")
	}
	if params.Period == "" {
		return nil, errors.New("bnc: period is required")
	}
	req := Req[ParamsUMLongShortRatio]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FUTURES_DATA + "/topLongShortAccountRatio",
		Params:  params,
	}
	resp, err := Request[ParamsUMLongShortRatio, []UMLongShortRatio](req)
	return resp.Data, err
}

func GetUMGlobalLongShortAccountRatio(params ParamsUMLongShortRatio) (longShortRatio []UMLongShortRatio, err error) {
	if params.Symbol == "" {
		return nil, errors.New("bnc: symbol is required")
	}
	if params.Period == "" {
		return nil, errors.New("bnc: period is required")
	}
	req := Req[ParamsUMLongShortRatio]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FUTURES_DATA + "/globalLongShortAccountRatio",
		Params:  params,
	}
	resp, err := Request[ParamsUMLongShortRatio, []UMLongShortRatio](req)
	return resp.Data, err
}

type TakerBuySellRatio struct {
	BuySellRatio float64 `json:"buySellRatio,string"`
	BuyVol       float64 `json:"buyVol,string"`
	SellVol      float64 `json:"sellVol,string"`
	Timestamp    int64   `json:"timestamp"`
}

type ParamsTakerBuySellRatio struct {
	Symbol    string `json:"symbol"`
	Period    string `json:"period"`          // "5m","15m","30m","1h","2h","4h","6h","12h","1d"
	Limit     int64  `json:"limit,omitempty"` // default 30, max 500
	StartTime int64  `json:"startTime,omitempty"`
	EndTime   int64  `json:"endTime,omitempty"`
}

func GetUMTakerBuySellRatio(params ParamsTakerBuySellRatio) (takerBuySellRatio []TakerBuySellRatio, err error) {
	if params.Symbol == "" {
		return nil, errors.New("bnc: symbol is required")
	}
	if params.Period == "" {
		return nil, errors.New("bnc: period is required")
	}
	req := Req[ParamsTakerBuySellRatio]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FUTURES_DATA + "/takerlongshortRatio",
		Params:  params,
	}
	resp, err := Request[ParamsTakerBuySellRatio, []TakerBuySellRatio](req)
	return resp.Data, err
}

type FuturesBasis struct {
	IndexPrice          float64      `json:"indexPrice,string"`
	ContractType        ContractType `json:"contractType"`
	BasisRate           float64      `json:"basisRate,string"`
	FuturesPrice        float64      `json:"futuresPrice,string"`
	AnnualizedBasisRate string       `json:"annualizedBasisRate"`
	Basis               float64      `json:"basis,string"`
	Pair                string       `json:"pair"`
	Timestamp           int64        `json:"timestamp"`
}

type ParamsFuturesBasis struct {
	Pair         string       `json:"pair"`
	ContractType ContractType `json:"contractType"`
	Period       string       `json:"period"` // "5m","15m","30m","1h","2h","4h","6h","12h","1d"
	Limit        int64        `json:"limit"`  // [30, 500]
	StartTime    int64        `json:"startTime,omitempty"`
	EndTime      int64        `json:"endTime,omitempty"`
}

func GetUMFuturesBasis(params ParamsFuturesBasis) (futuresBasis []FuturesBasis, err error) {
	if params.Pair == "" {
		return nil, errors.New("bnc: pair is required")
	}
	if params.ContractType == "" {
		return nil, errors.New("bnc: contractType is required")
	}
	if params.Period == "" {
		return nil, errors.New("bnc: period is required")
	}
	req := Req[ParamsFuturesBasis]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FUTURES_DATA + "/basis",
		Params:  params,
	}
	resp, err := Request[ParamsFuturesBasis, []FuturesBasis](req)
	return resp.Data, err
}

type FuturesBaseAsset struct {
	BaseAsset          string  `json:"baseAsset"`
	QuoteAsset         string  `json:"quoteAsset"`
	WeightInQuantity   float64 `json:"weightInQuantity,string"`
	WeightInPercentage float64 `json:"weightInPercentage,string"`
}

type CompositeIndexSymbolInfo struct {
	Symbol        string             `json:"symbol"`
	Time          int64              `json:"time"`
	Component     string             `json:"component"`
	BaseAssetList []FuturesBaseAsset `json:"baseAssetList"`
}

type ParamsCompositeIndexSymbolInfo struct {
	Symbol string `json:"symbol"`
}

func GetUMCompositeIndexSymbolInfo(params ParamsCompositeIndexSymbolInfo) (compositeIndexSymbolInfo CompositeIndexSymbolInfo, err error) {
	if params.Symbol == "" {
		return CompositeIndexSymbolInfo{}, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsCompositeIndexSymbolInfo]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/indexInfo",
		Params:  params,
	}
	resp, err := Request[ParamsCompositeIndexSymbolInfo, CompositeIndexSymbolInfo](req)
	return resp.Data, err
}

func GetUMCompositeIndexSymbolInfoList() (compositeIndexSymbolInfo []CompositeIndexSymbolInfo, err error) {
	req := Req[EmptyStruct]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/indexInfo",
	}
	resp, err := Request[EmptyStruct, []CompositeIndexSymbolInfo](req)
	return resp.Data, err
}

type MultiAssetsModeAssetIndex struct {
	Symbol                string `json:"symbol"`
	Time                  int64  `json:"time"`
	Index                 string `json:"index"`
	BidBuffer             string `json:"bidBuffer"`
	AskBuffer             string `json:"askBuffer"`
	BidRate               string `json:"bidRate"`
	AskRate               string `json:"askRate"`
	AutoExchangeBidBuffer string `json:"autoExchangeBidBuffer"`
	AutoExchangeAskBuffer string `json:"autoExchangeAskBuffer"`
	AutoExchangeBidRate   string `json:"autoExchangeBidRate"`
	AutoExchangeAskRate   string `json:"autoExchangeAskRate"`
}

type ParamsMultiAssetsModeAssetIndex struct {
	Symbol string `json:"symbol"`
}

func GetUMMultiAssetsModeAssetIndex(params ParamsMultiAssetsModeAssetIndex) (multiAssetsModeAssetIndex MultiAssetsModeAssetIndex, err error) {
	if params.Symbol == "" {
		return MultiAssetsModeAssetIndex{}, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsMultiAssetsModeAssetIndex]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/assetIndex",
		Params:  params,
	}
	resp, err := Request[ParamsMultiAssetsModeAssetIndex, MultiAssetsModeAssetIndex](req)
	return resp.Data, err
}

func GetUMMultiAssetsModeAssetIndexList() (multiAssetsModeAssetIndex []MultiAssetsModeAssetIndex, err error) {
	req := Req[EmptyStruct]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/assetIndex",
	}
	resp, err := Request[EmptyStruct, []MultiAssetsModeAssetIndex](req)
	return resp.Data, err
}

type IndexPriceConstituentInfo struct {
	Exchange string `json:"exchange"`
	Symbol   string `json:"symbol"`
	// UM Futures
	Price  float64 `json:"price,string"`
	Weight float64 `json:"weight,string"`
}

type IndexPriceConstituents struct {
	Symbol       string                      `json:"symbol"`
	Time         int64                       `json:"time"`
	Constituents []IndexPriceConstituentInfo `json:"constituents"`
}

type ParamsIndexPriceConstituentInfo struct {
	Symbol string `json:"symbol"`
}

func GetUMIndexPriceConstituentInfo(params ParamsIndexPriceConstituentInfo) (indexPriceConstituentInfo IndexPriceConstituents, err error) {
	if params.Symbol == "" {
		return IndexPriceConstituents{}, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsIndexPriceConstituentInfo]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/constituents",
		Params:  params,
	}
	resp, err := Request[ParamsIndexPriceConstituentInfo, IndexPriceConstituents](req)
	return resp.Data, err
}

type InsuranceFundAsset struct {
	Asset         string  `json:"asset"`
	MarginBalance float64 `json:"marginBalance,string"`
	UpdateTime    int64   `json:"updateTime"`
}

type InsuranceFundBalanceSnapshot struct {
	Symbols []string             `json:"symbols"`
	Assets  []InsuranceFundAsset `json:"assets"`
}

type ParamsInsuranceFundBalanceSnapshot struct {
	Symbol string `json:"symbol"`
}

func GetUMInsuranceFundBalanceSnapshot(params ParamsInsuranceFundBalanceSnapshot) (insuranceFundBalanceSnapshot InsuranceFundBalanceSnapshot, err error) {
	if params.Symbol == "" {
		return InsuranceFundBalanceSnapshot{}, errors.New("bnc: symbol is required")
	}
	req := Req[ParamsInsuranceFundBalanceSnapshot]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/insuranceBalance",
		Params:  params,
	}
	resp, err := Request[ParamsInsuranceFundBalanceSnapshot, InsuranceFundBalanceSnapshot](req)
	return resp.Data, err
}

func GetUMInsuranceFundBalanceSnapshotList() (insuranceFundBalanceSnapshot []InsuranceFundBalanceSnapshot, err error) {
	req := Req[EmptyStruct]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/insuranceBalance",
	}
	resp, err := Request[EmptyStruct, []InsuranceFundBalanceSnapshot](req)
	return resp.Data, err
}
