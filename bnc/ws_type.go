package bnc

import "math"

const (
	WsBaseUrl         = "wss://stream.binance.com:9443/ws"
	FutureWsBaseUrl   = "wss://fstream.binance.com/ws"
	CMFutureWsBaseUrl = "wss://dstream.binance.com/ws"
)

type WsMethod string

const (
	WsMethodSub     WsMethod = "SUBSCRIBE"
	WsMethodUnsub   WsMethod = "UNSUBSCRIBE"
	WsMethodRequest WsMethod = "REQUEST"
	WsMethodListSub WsMethod = "LIST_SUBSCRIPTIONS"
)

type WsEvent string

const (
	WsEDepthUpdate                  WsEvent = "depthUpdate"
	WsTrade                         WsEvent = "trade"
	WsAggTrade                      WsEvent = "aggTrade"
	WsMarginCall                    WsEvent = "MARGIN_CALL"
	WsAccountUpdate                 WsEvent = "ACCOUNT_UPDATE"
	WsOrderTradeUpdate              WsEvent = "ORDER_TRADE_UPDATE"
	WsAccountConfigUpdate           WsEvent = "ACCOUNT_CONFIG_UPDATE"
	WsStrategyUpdate                WsEvent = "STRATEGY_UPDATE"
	WsGridUpdate                    WsEvent = "GRID_UPDATE"
	WsConditionalOrderTriggerReject WsEvent = "CONDITIONAL_ORDER_TRIGGER_REJECT"

	WsEventOutboundAccountPosition       WsEvent = "outboundAccountPosition"
	WsEventBalanceUpdate                 WsEvent = "balanceUpdate"
	WsEventExecutionReport               WsEvent = "executionReport"
	WsEventListStatus                    WsEvent = "listStatus"
	WsEventListenKeyExpired              WsEvent = "listenKeyExpired"
	WsEventTrade                         WsEvent = "trade"
	WsEventAggTrade                      WsEvent = "aggTrade"
	WsEventKline                         WsEvent = "kline"
	WsEventDepthUpdate                   WsEvent = "depthUpdate"
	WsEventMarkPriceUpdate               WsEvent = "markPriceUpdate"
	WsEventIndexPriceUpdate              WsEvent = "indexPriceUpdate"
	WsEventForceOrder                    WsEvent = "forceOrder"
	WsEventMarginCall                    WsEvent = "MARGIN_CALL"
	WsEventAccountUpdate                 WsEvent = "ACCOUNT_UPDATE"
	WsEventOrderTradeUpdate              WsEvent = "ORDER_TRADE_UPDATE"
	WsEventAccountConfigUpdate           WsEvent = "ACCOUNT_CONFIG_UPDATE"
	WsEventStrategyUpdate                WsEvent = "STRATEGY_UPDATE"
	WsEventGridUpdate                    WsEvent = "GRID_UPDATE"
	WsEventConditionalOrderTradeUpdate   WsEvent = "CONDITIONAL_ORDER_TRADE_UPDATE"
	WsEventConditionalOrderTriggerReject WsEvent = "CONDITIONAL_ORDER_TRIGGER_REJECT"
	WsEventOpenOrderLoss                 WsEvent = "openOrderLoss"
	WsEventLiabilityChange               WsEvent = "liabilityChange"
	WsEventRiskLevelChange               WsEvent = "riskLevelChange"
)

type WsMsgId interface {
	~string | ~int64
}

type WsReq[P any, ID WsMsgId] struct {
	Method WsMethod `json:"method,omitempty"`
	Params P        `json:"params,omitempty"`
	// id can be int64/string/null
	Id ID `json:"id,omitempty"`
}

type WsResp[R any] struct {
	// id can be int64/string/null
	Id         any               `json:"id"`
	Status     int64             `json:"status"`
	Result     R                 `json:"result"`
	RateLimits []WsRespRateLimit `json:"rateLimits"`
	Error      *WsRespErr        `json:"error"`
}

type WsRespErr struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

type WsRespRateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int64  `json:"intervalNum"`
	Limit         int64  `json:"limit"`
	Count         int64  `json:"count"`
}

type WsDepthStream struct {
	EventType WsEvent    `json:"e"`
	EventTime int64      `json:"E"`
	Symbol    string     `json:"s"`
	FirstId   int64      `json:"U"`
	LastId    int64      `json:"u"`
	Bids      [][]string `json:"b"`
	Asks      [][]string `json:"a"`

	// just for future ob
	TxTime  int64 `json:"T"`
	PLastId int64 `json:"pu"`
}

type WsFuAggTradeStream struct {
	EventType    WsEvent `json:"e"`
	EventTime    int64   `json:"E"`
	Symbol       string  `json:"s"`
	AggID        int64   `json:"a"`
	Price        float64 `json:"p,string"`
	Quantity     float64 `json:"q,string"`
	FirstTradeId int64   `json:"f"`
	LastTradeId  int64   `json:"l"`
	TradeTime    int64   `json:"T"`
	IsBuyerMaker bool    `json:"m"`
}

type WsSpotBalance struct {
	Asset  string  `json:"a"`
	Free   float64 `json:"f,string"`
	Locked float64 `json:"l,string"`
}

type WsSpotAccountUpdate struct {
	EventType  string          `json:"e"`
	EventTime  int64           `json:"E"`
	UpdateTime int64           `json:"u"`
	Balances   []WsSpotBalance `json:"B"`

	// just for margin user data stream
	TimeUpdateId int64 `json:"U"`
}

type WsSpotBalanceUpdate struct {
	EventType        string  `json:"e"`
	EventTime        int64   `json:"E"`
	Asset            string  `json:"a"`
	SignBalanceDelta float64 `json:"d,string"`
	ClearTime        int64   `json:"T"`
}

func (w WsSpotBalanceUpdate) AbsBalanceDelta() float64 {
	return math.Abs(w.SignBalanceDelta)
}

type WsSpotListStatusObject struct {
	Symbol        string `json:"s"`
	OrderId       int64  `json:"i"`
	ClientOrderId string `json:"c"`
}

type WsSpotListStatus struct {
	EventType         string                   `json:"e"`
	EventTime         int64                    `json:"E"`
	Symbol            string                   `json:"s"`
	OrderListId       int                      `json:"g"`
	ContingencyType   string                   `json:"c"`
	ListStatusType    string                   `json:"l"`
	ListOrderStatus   string                   `json:"L"`
	ListRejectReason  string                   `json:"r"`
	ListClientOrderId string                   `json:"C"`
	Time              int64                    `json:"T"`
	Objects           []WsSpotListStatusObject `json:"O"`
}

type WsAggTradeStream struct {
	EventType    WsEvent `json:"e"`
	EventTime    int64   `json:"E"`
	Symbol       string  `json:"s"`
	AggTradeId   int64   `json:"a"`
	Price        float64 `json:"p,string"`
	Qty          float64 `json:"q,string"`
	FirstTradeId int64   `json:"f"`
	LastTradeId  int64   `json:"l"`
	TradeTime    int64   `json:"T"`
	IsMaker      bool    `json:"m"`
	M            bool    `json:"M"` // ignore
}

type WsTradeStream struct {
	EventType    WsEvent `json:"e"`
	EventTime    int64   `json:"E"`
	Symbol       string  `json:"s"`
	TradeID      int64   `json:"t"`
	Price        float64 `json:"p,string"`
	Quantity     float64 `json:"q,string"`
	TradeTime    int64   `json:"T"`
	IsBuyerMaker bool    `json:"m"`
	M            bool    `json:"M"` // ignore

	//
	BuyerOrderID  int64 `json:"b"`
	SellerOrderID int64 `json:"a"`
}

type WsKline struct {
	StartTime                int64         `json:"t"`
	CloseTime                int64         `json:"T"`
	Symbol                   string        `json:"s"`
	Interval                 KlineInterval `json:"i"`
	FirstTradeId             int64         `json:"f"`
	LastTradeId              int64         `json:"L"`
	OpenPrice                float64       `json:"o,string"`
	ClosePrice               float64       `json:"c,string"`
	HighPrice                float64       `json:"h,string"`
	LowPrice                 float64       `json:"l,string"`
	BaseAssetVolume          float64       `json:"v,string"`
	TradesNum                int64         `json:"n"`
	IsKlineClosed            bool          `json:"x"`
	QuoteVolume              float64       `json:"q,string"`
	BaseAssetTakerBuyVolume  float64       `json:"V,string"`
	QuoteAssetTakerBuyVolume float64       `json:"Q,string"`
	B                        string        `json:"B"` // ignore
}

func (wk WsKline) ToKline() Kline {
	return Kline{
		float64(wk.StartTime),
		float64(wk.OpenPrice),
		float64(wk.HighPrice),
		float64(wk.LowPrice),
		float64(wk.ClosePrice),
		float64(wk.BaseAssetVolume),
		float64(wk.CloseTime),
		float64(wk.QuoteVolume),
		float64(wk.TradesNum),
		float64(wk.BaseAssetTakerBuyVolume),
		float64(wk.QuoteAssetTakerBuyVolume),
	}
}

type WsKlineStream struct {
	EventType WsEvent `json:"e"`
	EventTime int64   `json:"E"`
	Symbol    string  `json:"s"`
	Kline     WsKline `json:"k"`
}

type WsOrderExecutionReport struct {
	EventType               string             `json:"e"`
	EventTime               int64              `json:"E"`
	Symbol                  string             `json:"s"`
	ClientOrderId           string             `json:"c"`
	Side                    OrderSide          `json:"S"`
	Type                    OrderType          `json:"o"`
	TimeInForce             TimeInForce        `json:"f"`
	Qty                     float64            `json:"q,string"`
	Price                   float64            `json:"p,string"`
	StopPrice               float64            `json:"P,string"`
	IcebergQty              float64            `json:"F,string"`
	OrderListId             int64              `json:"g"`
	OriginalClientId        string             `json:"C"`
	ExecutionType           OrderExecutionType `json:"x"`
	Status                  OrderStatus        `json:"X"`
	RejectReason            string             `json:"r"`
	OrderId                 int64              `json:"i"`
	LastExecutedQty         float64            `json:"l,string"`
	FilledQty               float64            `json:"z,string"`
	LastExecutedPrice       float64            `json:"L,string"`
	CommissionAmt           float64            `json:"n,string"`
	CommissionAsset         string             `json:"N"`
	Time                    int64              `json:"T"`
	TradeId                 int64              `json:"t"`
	PreventedMatchId        int64              `json:"v"`
	Ignore                  int64              `json:"I"`
	IsOrderOnTheBook        bool               `json:"w"`
	IsMaker                 bool               `json:"m"`
	Ignore1                 bool               `json:"M"`
	CreationTime            int64              `json:"O"`
	FilledQuote             float64            `json:"Z,string"`
	LastExecutedQuote       float64            `json:"Y,string"`
	QuoteOrderQty           float64            `json:"Q,string"`
	WorkingTime             int64              `json:"W"`
	SelfTradePreventionMode STPMode            `json:"V"`

	// just for margin order
	TrailingDelta      float64 `json:"d"`    // Trailing Delta; This is only visible if the order was a trailing stop order.
	TrailingTime       int64   `json:"Data"` // Trailing Time; This is only visible if the trailing stop order has been activated.
	MarginStrategyId   int64   `json:"j"`
	MarginStrategyType int64   `json:"J"`
	TradeGroupId       int64   `json:"u"`
	CounterOrderId     int64   `json:"U"`
	PreventedQty       float64 `json:"A,string"`
	LastPreventedQty   float64 `json:"B,string"`

	// just for futures order
	AvgPrice            string       `json:"ap"`
	Sp                  string       `json:"sp"` // ignore
	BidNotional         string       `json:"b"`
	AskNotional         string       `json:"a"`
	IsReduceOnly        bool         `json:"R"`
	PositionSide        PositionSide `json:"ps"`
	RealizedProfit      string       `json:"rp"`
	FuturesStrategyType string       `json:"st"`
	FuturesStrategyId   int64        `json:"si"`
	Gtd                 int64        `json:"gtd"`
}

type WsListenKeyExpired struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	ListenKey string `json:"listenKey"`
}

// WsCMIndexPriceStream
// url: wss://dstream.binance.com/ws
// event: WsEventIndexPriceUpdate
type WsCMIndexPriceStream struct {
	EventType WsEvent `json:"e"`
	EventTime int64   `json:"E"`
	Pair      string  `json:"i"`
	Price     float64 `json:"p,string"`
}

// WsMarkPriceStream
// url: wss://dstream.binance.com/ws
// event: WsEventMarkPriceUpdate
type WsMarkPriceStream struct {
	EventType            WsEvent `json:"e"`
	EventTime            int64   `json:"E"`
	Symbol               string  `json:"s"`
	MarkPrice            string  `json:"p"`
	EstimatedSettlePrice string  `json:"P"`
	IndexPrice           string  `json:"i"`
	FundingRate          string  `json:"r"`
	NextFundingTime      int64   `json:"T"`
}

type WsStrategyOrder struct {
	Symbol               string       `json:"s"`
	ClientOrderId        string       `json:"c"`
	Id                   int64        `json:"si"`
	Side                 OrderSide    `json:"S"`
	Type                 string       `json:"st"`
	TimeInForce          TimeInForce  `json:"f"`
	Qty                  float64      `json:"q,string"`
	Price                float64      `json:"p,string"`
	StopPrice            float64      `json:"sp,string"`
	Status               OrderStatus  `json:"os"`
	OrderBookTime        int64        `json:"T"`
	UpdateTime           int64        `json:"ut"`
	IsThisReduceOnly     bool         `json:"R"`
	StopPriceWorkingType string       `json:"wt"`
	PositionSide         PositionSide `json:"ps"`
	Cp                   bool         `json:"cp"` // If Close-All, pushed with conditional order
	ActivationPrice      float64      `json:"AP,string"`
	CallbackRate         float64      `json:"cr,string"`
	OrderId              int64        `json:"i"`
	STPMode              string       `json:"V"`
	Gtd                  int64        `json:"gtd"`
}

type WsBusinessUnit string

const (
	WsBusinessUnitUM WsBusinessUnit = "UM"
	WsBusinessUnitCM WsBusinessUnit = "CM"
)

// WsConditionalOrderTradeUpdate
// url: https://papi.binance.com/ws
// event: WsEventConditionalOrderTradeUpdate
type WsConditionalOrderTradeUpdate struct {
	EventType       WsEvent         `json:"e"`
	EventTime       int64           `json:"E"`
	TransactionTime int64           `json:"T"`
	BusinessUnit    WsBusinessUnit  `json:"fs"`
	StrategyOrder   WsStrategyOrder `json:"so"`
}

type WsOpenOrderLossUpdate struct {
	Asset      string  `json:"a"`
	SignAmount float64 `json:"o"`
}

// WsOpenOrderLoss
// url: https://papi.binance.com/ws
// event: WsEventOpenOrderLoss
type WsOpenOrderLoss struct {
	EventType WsEvent                 `json:"e"`
	EventTime int64                   `json:"E"`
	OrderLoss []WsOpenOrderLossUpdate `json:"O"`
}

// WsMarginAccountUpdate
// url: https://papi.binance.com/ws
// event: WsEventOutboundAccountPosition
type WsMarginAccountUpdate WsSpotAccountUpdate

type LiabilityType string

const (
	LiabilityTypeBorrow LiabilityType = "BORROW"
)

// WsLiabilityChange
// url: https://papi.binance.com/ws
// event: WsEventLiabilityChange
type WsLiabilityChange struct {
	EventType      WsEvent       `json:"e"`
	EventTime      int64         `json:"E"`
	Asset          string        `json:"a"`
	Type           LiabilityType `json:"t"`
	TxId           int64         `json:"tx"`
	Principal      float64       `json:"p,string"`
	Interest       float64       `json:"i,string"`
	TotalLiability float64       `json:"l,string"`
}

// WsPmFuturesOrderTradeUpdate
// url: https://papi.binance.com/ws
// event: WsEventOrderTradeUpdate
type WsPmFuturesOrderTradeUpdate struct {
	EventType       WsEvent                `json:"e"`
	EventTime       int64                  `json:"E"`
	TransactionTime int64                  `json:"T"`
	BusinessUnit    WsBusinessUnit         `json:"fs"`
	AccountAlias    string                 `json:"i"` // Account Alias,ignore for UM
	Order           WsOrderExecutionReport `json:"o"`
}

type WsFuturesWalletBalance struct {
	Asset              string  `json:"a"`
	WalletBalance      float64 `json:"wb,string"`
	CrossWalletBalance float64 `json:"cw,string"`
	BalanceChange      float64 `json:"bc,string"`
}

type WsFuturesPosition struct {
	Symbol                    string       `json:"s"`
	SignPositionQty           float64      `json:"pa"`
	EntryPrice                float64      `json:"ep"`
	PreFeeAccumulatedRealized float64      `json:"cr"`
	UnrealizedPnl             float64      `json:"up"`
	PositionSide              PositionSide `json:"ps"`
	BreakevenPrice            string       `json:"bep"`
}

type WsPmFuturesAcctUpdateData struct {
	Reason    string                   `json:"m"`
	Balances  []WsFuturesWalletBalance `json:"B"`
	Positions []WsFuturesPosition      `json:"P"`
}

// WsPmFuturesAcctUpdateStream
// url: https://papi.binance.com/ws
// event: WsEventAccountUpdate
type WsPmFuturesAcctUpdateStream struct {
	EventType       WsEvent                   `json:"e"`
	EventTime       int64                     `json:"E"`
	TransactionTime int64                     `json:"T"`
	BusinessUnit    WsBusinessUnit            `json:"fs"`
	AccountAlias    string                    `json:"i"` // Account Alias,ignore for UM
	Data            WsPmFuturesAcctUpdateData `json:"a"`
}

type WsFuturesCfg struct {
	Symbol   string  `json:"s"`
	Leverage float64 `json:"l"`
}

// WsFuturesAccountCfgUpdate
// url: https://papi.binance.com/ws
// event: WsEventAccountConfigUpdate
type WsFuturesAccountCfgUpdate struct {
	EventType       WsEvent        `json:"e"`
	EventTime       int64          `json:"E"`
	TransactionTime int64          `json:"T"`
	BusinessUnit    WsBusinessUnit `json:"fs"`
	Cfg             WsFuturesCfg   `json:"ac"`
}

type WsPmRiskLevelChangeType string

const (
	WsPmRiskLevelChangeTypeMarginCall       WsPmRiskLevelChangeType = "MARGIN_CALL"
	WsPmRiskLevelChangeTypeSupplyMargin     WsPmRiskLevelChangeType = "SUPPLY_MARGIN"
	WsPmRiskLevelChangeTypeReduceOnly       WsPmRiskLevelChangeType = "REDUCE_ONLY"
	WsPmRiskLevelChangeTypeForceLiquidation WsPmRiskLevelChangeType = "FORCE_LIQUIDATION"
)

// WsPmRiskLevelChange
// url: https://papi.binance.com/ws
// event: WsEventRiskLevelChange
type WsPmRiskLevelChange struct {
	EventType                             WsEvent                 `json:"e"`
	EventTime                             int64                   `json:"E"`
	UniMMR                                string                  `json:"u"`
	ChangeType                            WsPmRiskLevelChangeType `json:"s"`
	AccountEquityUSD                      float64                 `json:"eq,string"`
	AccountEquityUSDWithoutCollateralRate float64                 `json:"ae,string"`
	MaintenanceMarginUSD                  float64                 `json:"m,string"`
}

// WsMarginBalanceUpdate
// url: https://papi.binance.com/ws
// event: balanceUpdate
type WsMarginBalanceUpdate struct {
	Asset        string  `json:"a"`
	BalanceDelta float64 `json:"d,string"`
	UpdateId     int64   `json:"U"`
	ClearTime    int64   `json:"T"`
}

type WsLiquidationOrder struct {
	Symbol        string      `json:"s"`
	OrderSide     OrderSide   `json:"S"`
	OrderType     OrderType   `json:"o"`
	TimeInForce   TimeInForce `json:"f"`
	OriQty        float64     `json:"q,string"`
	Price         float64     `json:"p,string"`
	AvgPrice      float64     `json:"ap,string"`
	OrderStatus   OrderStatus `json:"X"`
	LastFilledQty float64     `json:"l,string"`
	FilledQty     float64     `json:"z,string"`
	TradeTime     int64       `json:"T"`
}

type WsLiquidationOrderStream struct {
	EventType WsEvent            `json:"e"`
	EventTime int64              `json:"E"`
	Order     WsLiquidationOrder `json:"o"`
}
