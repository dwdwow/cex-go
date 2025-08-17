package bnc

import (
	"fmt"
	"slices"
	"strings"

	"github.com/dwdwow/cex-go"
)

func PingSpotEndpoint() (err error) {
	req := Req[EmptyStruct]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/ping",
	}

	_, err = Request[EmptyStruct, EmptyStruct](req)

	return
}

type ServerTime struct {
	ServerTime int64 `json:"serverTime"`
}

func GetSpotServerTime() (serverTime ServerTime, err error) {
	req := Req[EmptyStruct]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/time",
	}
	resp, err := Request[EmptyStruct, ServerTime](req)
	if err != nil {
		return
	}
	serverTime = resp.Data
	return
}

type ParamsSpotExchangeInfo struct {
	Symbol             string              `json:"symbol,omitempty"`
	Symbols            []string            `json:"symbols,omitempty"`
	Permissions        []AcctSybPermission `json:"permissions,omitempty"`
	ShowPermissionSets *bool               `json:"showPermissionSets,omitempty"`
	SymbolStatus       SymbolStatus        `json:"symbolStatus,omitempty"`
}

type SpotExchangeSymbol struct {
	Symbol                          string                `json:"symbol"`
	Status                          SymbolStatus          `json:"status"`
	BaseAsset                       string                `json:"baseAsset"`
	BaseAssetPrecision              int64                 `json:"baseAssetPrecision"`
	QuoteAsset                      string                `json:"quoteAsset"`
	QuotePrecision                  int64                 `json:"quotePrecision"`
	QuoteAssetPrecision             int64                 `json:"quoteAssetPrecision"`
	BaseCommissionPrecision         int64                 `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision        int64                 `json:"quoteCommissionPrecision"`
	OrderTypes                      []OrderType           `json:"orderTypes"`
	IcebergAllowed                  bool                  `json:"icebergAllowed"`
	OcoAllowed                      bool                  `json:"ocoAllowed"`
	OtoAllowed                      bool                  `json:"otoAllowed"`
	QuoteOrderQtyMarketAllowed      bool                  `json:"quoteOrderQtyMarketAllowed"`
	AllowTrailingStop               bool                  `json:"allowTrailingStop"`
	CancelReplaceAllowed            bool                  `json:"cancelReplaceAllowed"`
	AmendAllowed                    bool                  `json:"amendAllowed"`
	IsSpotTradingAllowed            bool                  `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed          bool                  `json:"isMarginTradingAllowed"`
	Filters                         []map[string]any      `json:"filters"`
	Permissions                     []AcctSybPermission   `json:"permissions"`
	PermissionSets                  [][]AcctSybPermission `json:"permissionSets"`
	DefaultSelfTradePreventionMode  string                `json:"defaultSelfTradePreventionMode"`
	AllowedSelfTradePreventionModes []string              `json:"allowedSelfTradePreventionModes"`
}

func GetPrecJustForBinanceFilter(size string) (int, error) {
	ss := strings.Split(size, ".")
	s0 := ss[0]
	i1 := strings.Index(s0, "1")
	if i1 != -1 {
		return i1 - len(s0) + 1, nil
	}
	if len(ss) == 1 {
		return 0, fmt.Errorf("unknown size: %v", size)
	}
	s1 := ss[1]
	f1 := strings.Index(s1, "1")
	if f1 != -1 {
		return f1 + 1, nil
	}
	return 0, fmt.Errorf("unknown size: %v", size)
}

type ExchangeSymbolFiltersInfo struct {
	Pprec     int
	Qprec     int
	CanMarket bool
}

func AnalyzeExchangeSymbolFilters(filters []map[string]any) (info ExchangeSymbolFiltersInfo, err error) {
	var pPrec, qPrec int
	for _, filter := range filters {
		t := filter["filterType"]
		switch t := t.(type) {
		case string:
			switch t {
			case "PRICE_FILTER":
				tick := filter["tickSize"]
				switch tick := tick.(type) {
				case string:
					s, err := GetPrecJustForBinanceFilter(tick)
					if err != nil {
						return info, err
					}
					pPrec = s
				default:
					// should not get here
					return info, fmt.Errorf("exchange info tickSize type is not string, tick size %v", tick)
				}
			case "LOT_SIZE":
				step := filter["stepSize"]
				switch step := step.(type) {
				case string:
					s, err := GetPrecJustForBinanceFilter(step)
					if err != nil {
						return info, err
					}
					qPrec = s
				default:
					// should not get here
					return info, fmt.Errorf("exchange info stepSize type is not string, step size %v", step)
				}
			}
		default:
			// should not get here
			return info, fmt.Errorf("exchange info filter type is not string, type %v", t)
		}
	}
	info.Pprec = pPrec
	info.Qprec = qPrec
	return info, nil
}

func (info SpotExchangeSymbol) ToSymbol() (cex.Symbol, error) {
	filtersInfo, err := AnalyzeExchangeSymbolFilters(info.Filters)
	if err != nil {
		return cex.Symbol{}, err
	}
	pair := cex.Symbol{
		Cex:        cex.BINANCE,
		Type:       cex.SYMBOL_TYPE_SPOT,
		Asset:      info.BaseAsset,
		Quote:      info.QuoteAsset,
		Symbol:     info.Symbol,
		QPrecision: filtersInfo.Qprec,
		PPrecision: filtersInfo.Pprec,
		Tradable:   info.Status == SYMBOL_STATUS_TRADING,
		CanMarket:  slices.Contains(info.OrderTypes, ORDER_TYPE_MARKET),
		CanMargin:  info.IsMarginTradingAllowed,
	}
	return pair, nil
}

func GetSpotSymbols() (symbols []cex.Symbol, err error) {
	exchangeInfo, err := GetSpotExchangeInfo(ParamsSpotExchangeInfo{})
	if err != nil {
		return nil, err
	}
	for _, symbol := range exchangeInfo.Symbols {
		s, err := symbol.ToSymbol()
		if err != nil {
			return nil, err
		}
		symbols = append(symbols, s)
	}
	return
}

type RateLimier struct {
	RateLimitType RateLimitType       `json:"rateLimitType"`
	Interval      RateLimiterInterval `json:"interval"`
	IntervalNum   int64               `json:"intervalNum"`
	Limit         int64               `json:"limit"`
}

type SOR struct {
	BaseAsset string   `json:"baseAsset"`
	Symbols   []string `json:"symbols"`
}

type SpotExchangeInfo struct {
	Timezone        string               `json:"timezone"`
	ServerTime      int64                `json:"serverTime"`
	RateLimits      []RateLimier         `json:"rateLimits"`
	ExchangeFilters []any                `json:"exchangeFilters"`
	Symbols         []SpotExchangeSymbol `json:"symbols"`
	Sors            []SOR                `json:"sors"`
}

func GetSpotExchangeInfo(params ParamsSpotExchangeInfo) (exchangeInfo SpotExchangeInfo, err error) {
	req := Req[ParamsSpotExchangeInfo]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/exchangeInfo",
		Params:  params,
	}
	resp, err := Request[ParamsSpotExchangeInfo, SpotExchangeInfo](req)
	if err != nil {
		return
	}
	exchangeInfo = resp.Data
	return
}
