package bnc

func GetUMServerTime() (serverTime ServerTime, err error) {
	req := Req[EmptyStruct]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/time",
	}
	resp, err := Request[EmptyStruct, ServerTime](req)
	return resp.Data, err
}

type FuturesExchangeSymbol struct {
	Symbol             string       `json:"symbol"`
	Pair               string       `json:"pair"`
	ContractType       ContractType `json:"contractType"`
	DeliveryDate       int64        `json:"deliveryDate"`
	OnboardDate        int64        `json:"onboardDate"`
	Status             string       `json:"status"`
	BaseAsset          string       `json:"baseAsset"`
	QuoteAsset         string       `json:"quoteAsset"`
	MarginAsset        string       `json:"marginAsset"`
	PricePrecision     int64        `json:"pricePrecision"`
	QuantityPrecision  int64        `json:"quantityPrecision"`
	BaseAssetPrecision int64        `json:"baseAssetPrecision"`
	QuotePrecision     int64        `json:"quotePrecision"`
	UnderlyingType     string       `json:"underlyingType"`
	UnderlyingSubType  []string     `json:"underlyingSubType"`
	SettlePlan         int64        `json:"settlePlan"`
	TriggerProtect     string       `json:"triggerProtect"`
	Filters            []any        `json:"filters"`
	OrderTypes         []string     `json:"OrderType"`
	TimeInForce        []string     `json:"timeInForce"`
	LiquidationFee     float64      `json:"liquidationFee,string"`
	MarketTakeBound    float64      `json:"marketTakeBound,string"`
	PermissionSets     []any        `json:"permissionSets"`
}

type FuturesMarginAsset struct {
	Asset             string  `json:"asset"`                    // The asset code (e.g. "BTC")
	MarginAvailable   bool    `json:"marginAvailable"`          // Whether the asset can be used as margin in Multi-Assets mode
	AutoAssetExchange float64 `json:"autoAssetExchange,string"` // Auto-exchange threshold in Multi-Assets margin mode
}

type FuturesExchangeInfo struct {
	Timezone        string                  `json:"timezone"`
	ServerTime      int64                   `json:"serverTime"`
	RateLimits      []RateLimier            `json:"rateLimits"`
	ExchangeFilters []any                   `json:"exchangeFilters"`
	Assets          []FuturesMarginAsset    `json:"assets"`
	Symbols         []FuturesExchangeSymbol `json:"symbols"`
}

func GetUMExchangeInfo() (exchangeInfo FuturesExchangeInfo, err error) {
	req := Req[EmptyStruct]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/exchangeInfo",
	}
	resp, err := Request[EmptyStruct, FuturesExchangeInfo](req)
	return resp.Data, err
}
