package bnc

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

type ExchangeSymbol struct {
	Symbol                          string                `json:"symbol"`
	Status                          string                `json:"status"`
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
	Filters                         []any                 `json:"filters"`
	Permissions                     []AcctSybPermission   `json:"permissions"`
	PermissionSets                  [][]AcctSybPermission `json:"permissionSets"`
	DefaultSelfTradePreventionMode  string                `json:"defaultSelfTradePreventionMode"`
	AllowedSelfTradePreventionModes []string              `json:"allowedSelfTradePreventionModes"`
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

type ExchangeInfo struct {
	Timezone        string           `json:"timezone"`
	ServerTime      int64            `json:"serverTime"`
	RateLimits      []RateLimier     `json:"rateLimits"`
	ExchangeFilters []any            `json:"exchangeFilters"`
	Symbols         []ExchangeSymbol `json:"symbols"`
	Sors            []SOR            `json:"sors"`
}

func GetSpotExchangeInfo(params ParamsSpotExchangeInfo) (exchangeInfo ExchangeInfo, err error) {
	req := Req[ParamsSpotExchangeInfo]{
		BaseURL: API_ENDPOINT,
		Path:    API_V3 + "/exchangeInfo",
		Params:  params,
	}
	resp, err := Request[ParamsSpotExchangeInfo, ExchangeInfo](req)
	if err != nil {
		return
	}
	exchangeInfo = resp.Data
	return
}
