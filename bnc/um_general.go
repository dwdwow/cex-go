package bnc

import (
	"slices"

	"github.com/dwdwow/cex-go"
)

func GetUMServerTime() (serverTime ServerTime, err error) {
	req := Req[EmptyStruct]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/time",
	}
	resp, err := Request[EmptyStruct, ServerTime](req)
	return resp.Data, err
}

type FuturesExchangeSymbol struct {
	Symbol             string           `json:"symbol" bson:"symbol"` // e.g. SPOT/PERPETUAL:"BTCUSDT", Deliviry:"BTCUSDT_200925" CM:"BTCUSD_PERP"
	Pair               string           `json:"pair" bson:"pair"`   // underlying symbol e.g. "BTCUSDT" "BTCUSD"
	ContractType       ContractType     `json:"contractType" bson:"contractType"`
	DeliveryDate       int64            `json:"deliveryDate" bson:"deliveryDate"`
	OnboardDate        int64            `json:"onboardDate" bson:"onboardDate"`
	Status             SymbolStatus     `json:"status" bson:"status"`
	ContractStatus     ContractStatus   `json:"contractStatus" bson:"contractStatus"`
	ContractSize       float64          `json:"contractSize" bson:"contractSize"`
	BaseAsset          string           `json:"baseAsset" bson:"baseAsset"`
	QuoteAsset         string           `json:"quoteAsset" bson:"quoteAsset"`
	MarginAsset        string           `json:"marginAsset" bson:"marginAsset"`
	PricePrecision     int64            `json:"pricePrecision" bson:"pricePrecision"`
	QuantityPrecision  int64            `json:"quantityPrecision" bson:"quantityPrecision"`
	BaseAssetPrecision int64            `json:"baseAssetPrecision" bson:"baseAssetPrecision"`
	QuotePrecision     int64            `json:"quotePrecision" bson:"quotePrecision"`
	UnderlyingType     string           `json:"underlyingType" bson:"underlyingType"`
	UnderlyingSubType  []string         `json:"underlyingSubType" bson:"underlyingSubType"`
	SettlePlan         int64            `json:"settlePlan" bson:"settlePlan"`
	TriggerProtect     string           `json:"triggerProtect" bson:"triggerProtect"`
	Filters            []map[string]any `json:"filters" bson:"filters"`
	OrderTypes         []OrderType      `json:"orderTypes" bson:"orderTypes"`
	TimeInForce        []TimeInForce    `json:"timeInForce" bson:"timeInForce"`
	LiquidationFee     float64          `json:"liquidationFee,string" bson:"liquidationFee,string"`
	MarketTakeBound    float64          `json:"marketTakeBound,string" bson:"marketTakeBound,string"`
	PermissionSets     []any            `json:"permissionSets" bson:"permissionSets"`
}

func (info FuturesExchangeSymbol) ToSymbol() (cex.Symbol, error) {
	filtersInfo, err := AnalyzeExchangeSymbolFilters(info.Filters)
	if err != nil {
		return cex.Symbol{}, err
	}
	pair := cex.Symbol{
		Cex:         cex.BINANCE,
		Type:        cex.SYMBOL_TYPE_UM_FUTURES,
		Asset:       info.BaseAsset,
		Quote:       info.QuoteAsset,
		Symbol:      info.Symbol,
		QPrecision:  filtersInfo.Qprec,
		PPrecision:  filtersInfo.Pprec,
		Tradable:    info.Status == SYMBOL_STATUS_TRADING || info.ContractStatus == CONTRACT_STATUS_TRADING,
		CanMarket:   slices.Contains(info.OrderTypes, ORDER_TYPE_MARKET),
		IsPerpetual: info.ContractType == CONTRACT_TYPE_PERPETUAL,

		ContractSize: info.ContractSize,
		ContractType: string(info.ContractType),

		DeliveryDate: info.DeliveryDate,
		OnboardDate:  info.OnboardDate,
	}
	return pair, nil
}

type FuturesMarginAsset struct {
	Asset             string  `json:"asset" bson:"asset"`                    // The asset code (e.g. "BTC")
	MarginAvailable   bool    `json:"marginAvailable" bson:"marginAvailable"`          // Whether the asset can be used as margin in Multi-Assets mode
	AutoAssetExchange float64 `json:"autoAssetExchange,string" bson:"autoAssetExchange,string"` // Auto-exchange threshold in Multi-Assets margin mode
}

type FuturesExchangeInfo struct {
	Timezone        string                  `json:"timezone" bson:"timezone"`
	ServerTime      int64                   `json:"serverTime" bson:"serverTime"`
	RateLimits      []RateLimier            `json:"rateLimits" bson:"rateLimits"`
	ExchangeFilters []any                   `json:"exchangeFilters" bson:"exchangeFilters"`
	Assets          []FuturesMarginAsset    `json:"assets" bson:"assets"`
	Symbols         []FuturesExchangeSymbol `json:"symbols" bson:"symbols"`
}

func GetUMExchangeInfo() (exchangeInfo FuturesExchangeInfo, err error) {
	req := Req[EmptyStruct]{
		BaseURL: API_UM_FUTURES_ENDPOINT,
		Path:    FAPI_V1 + "/exchangeInfo",
	}
	resp, err := Request[EmptyStruct, FuturesExchangeInfo](req)
	return resp.Data, err
}

func GetUMSymbols() (symbols []cex.Symbol, err error) {
	exchangeInfo, err := GetUMExchangeInfo()
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
