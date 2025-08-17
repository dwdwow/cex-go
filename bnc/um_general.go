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
	Symbol             string           `json:"symbol"` // e.g. SPOT/PERPETUAL:"BTCUSDT", Deliviry:"BTCUSDT_200925" CM:"BTCUSD_PERP"
	Pair               string           `json:"pair"`   // underlying symbol e.g. "BTCUSDT" "BTCUSD"
	ContractType       ContractType     `json:"contractType"`
	DeliveryDate       int64            `json:"deliveryDate"`
	OnboardDate        int64            `json:"onboardDate"`
	Status             SymbolStatus     `json:"status"`
	ContractStatus     ContractStatus   `json:"contractStatus"`
	ContractSize       float64          `json:"contractSize"`
	BaseAsset          string           `json:"baseAsset"`
	QuoteAsset         string           `json:"quoteAsset"`
	MarginAsset        string           `json:"marginAsset"`
	PricePrecision     int64            `json:"pricePrecision"`
	QuantityPrecision  int64            `json:"quantityPrecision"`
	BaseAssetPrecision int64            `json:"baseAssetPrecision"`
	QuotePrecision     int64            `json:"quotePrecision"`
	UnderlyingType     string           `json:"underlyingType"`
	UnderlyingSubType  []string         `json:"underlyingSubType"`
	SettlePlan         int64            `json:"settlePlan"`
	TriggerProtect     string           `json:"triggerProtect"`
	Filters            []map[string]any `json:"filters"`
	OrderTypes         []OrderType      `json:"orderTypes"`
	TimeInForce        []TimeInForce    `json:"timeInForce"`
	LiquidationFee     float64          `json:"liquidationFee,string"`
	MarketTakeBound    float64          `json:"marketTakeBound,string"`
	PermissionSets     []any            `json:"permissionSets"`
}

func (info FuturesExchangeSymbol) ToPair() (cex.Pair, error) {
	filtersInfo, err := AnalyzeExchangeSymbolFilters(info.Filters)
	if err != nil {
		return cex.Pair{}, err
	}
	pair := cex.Pair{
		Cex:         cex.BINANCE,
		Type:        cex.PairTypeFutures,
		Asset:       info.BaseAsset,
		Quote:       info.QuoteAsset,
		PairSymbol:  info.Symbol,
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

func GetUMPairs() (pairs []cex.Pair, err error) {
	exchangeInfo, err := GetUMExchangeInfo()
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
