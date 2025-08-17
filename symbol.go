package cex

import "slices"

type SymbolType string

const (
	SYMBOL_TYPE_SPOT       SymbolType = "SPOT"
	SYMBOL_TYPE_UM_FUTURES SymbolType = "UM_FUTURES"
	SYMBOL_TYPE_CM_FUTURES SymbolType = "CM_FUTURES"
)

var SymbolTypes = []SymbolType{
	SYMBOL_TYPE_SPOT,
	SYMBOL_TYPE_UM_FUTURES,
	SYMBOL_TYPE_CM_FUTURES,
}

func NotSymbolType(t SymbolType) bool {
	return !slices.Contains(SymbolTypes, t)
}

type Symbol struct {
	// must be contained
	Cex        CexName    `json:"cex" bson:"cex"`
	Type       SymbolType `json:"type" bson:"type"`
	Asset      string     `json:"asset" bson:"asset"`
	Quote      string     `json:"quote" bson:"quote"`
	Symbol     string     `json:"symbol" bson:"symbol"`
	MidSymbol  string     `json:"midSymbol" bson:"midSymbol"`
	QPrecision int        `json:"qPrecision" bson:"qPrecision"`
	PPrecision int        `json:"pPrecision" bson:"pPrecision"`

	// may be omitted
	TakerFeeTier  float64 `json:"takerFeeTier" bson:"takerFeeTier"`
	MakerFeeTier  float64 `json:"makerFeeTier" bson:"makerFeeTier"`
	MinTradeQty   float64 `json:"minTradeQty" bson:"minTradeQty"`
	MinTradeQuote float64 `json:"minTradeQuote" bson:"minTradeQuote"`
	Tradable      bool    `json:"tradable" bson:"tradable"`
	CanMarket     bool    `json:"canMarket" bson:"canMarket"`
	CanMargin     bool    `json:"canMargin" bson:"canMargin"`
	IsCross       bool    `json:"isCross" bson:"isCross"`

	IsPerpetual bool `json:"isPerpetual" bson:"isPerpetual"`

	// just for cm futures
	ContractSize float64 `json:"contractSize" bson:"contractSize"`
	ContractType string  `json:"contractType" bson:"contractType"`

	// just for delivery contract
	DeliveryDate int64 `json:"deliveryDate" bson:"deliveryDate"`
	OnboardDate  int64 `json:"onboardDate" bson:"onboardDate"`
}
