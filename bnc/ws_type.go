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
	WsEventBookTicker                    WsEvent = "bookTicker"
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
	Method WsMethod `json:"method,omitempty" bson:"method,omitempty"`
	Params P        `json:"params,omitempty" bson:"params,omitempty"`
	// id can be int64/string/null
	Id ID `json:"id,omitempty" bson:"id,omitempty"`
}

type WsResp[R any] struct {
	// id can be int64/string/null
	Id         any               `json:"id" bson:"id"`
	Status     int64             `json:"status" bson:"status"`
	Result     R                 `json:"result" bson:"result"`
	RateLimits []WsRespRateLimit `json:"rateLimits" bson:"rateLimits"`
	Error      *WsRespErr        `json:"error" bson:"error"`
}

type WsRespErr struct {
	Code int64  `json:"code" bson:"code"`
	Msg  string `json:"msg" bson:"msg"`
}

type WsRespRateLimit struct {
	RateLimitType string `json:"rateLimitType" bson:"rateLimitType"`
	Interval      string `json:"interval" bson:"interval"`
	IntervalNum   int64  `json:"intervalNum" bson:"intervalNum"`
	Limit         int64  `json:"limit" bson:"limit"`
	Count         int64  `json:"count" bson:"count"`
}

type WsDepthStream struct {
	EventType WsEvent    `json:"e" bson:"e"`
	EventTime int64      `json:"E" bson:"E"`
	Symbol    string     `json:"s" bson:"s"`
	FirstId   int64      `json:"U" bson:"U"`
	LastId    int64      `json:"u" bson:"u"`
	Bids      [][]string `json:"b" bson:"b"`
	Asks      [][]string `json:"a" bson:"a"`

	// just for future ob
	TxTime  int64 `json:"T" bson:"T"`
	PLastId int64 `json:"pu" bson:"pu"`
}

type WsBookTickerStream struct {
	EventType         WsEvent `json:"e" bson:"e"`
	OrderBookUpdateId int64   `json:"u" bson:"u"`
	EventTime         int64   `json:"E" bson:"E"`
	TransactionTime   int64   `json:"T" bson:"T"`
	Symbol            string  `json:"s" bson:"s"`
	BestBidPrice      string  `json:"b" bson:"b"`
	BestBidQty        string  `json:"B" bson:"B"`
	BestAskPrice      string  `json:"a" bson:"a"`
	BestAskQty        string  `json:"A" bson:"A"`
}

type WsFuAggTradeStream struct {
	EventType    WsEvent `json:"e" bson:"e"`
	EventTime    int64   `json:"E" bson:"E"`
	Symbol       string  `json:"s" bson:"s"`
	AggID        int64   `json:"a" bson:"a"`
	Price        float64 `json:"p,string" bson:"p,string"`
	Quantity     float64 `json:"q,string" bson:"q,string"`
	FirstTradeId int64   `json:"f" bson:"f"`
	LastTradeId  int64   `json:"l" bson:"l"`
	TradeTime    int64   `json:"T" bson:"T"`
	IsBuyerMaker bool    `json:"m" bson:"m"`
}

type WsSpotBalance struct {
	Asset  string  `json:"a" bson:"a"`
	Free   float64 `json:"f,string" bson:"f,string"`
	Locked float64 `json:"l,string" bson:"l,string"`
}

type WsSpotAccountUpdate struct {
	EventType  string          `json:"e" bson:"e"`
	EventTime  int64           `json:"E" bson:"E"`
	UpdateTime int64           `json:"u" bson:"u"`
	Balances   []WsSpotBalance `json:"B" bson:"B"`

	// just for margin user data stream
	TimeUpdateId int64 `json:"U" bson:"U"`
}

type WsSpotBalanceUpdate struct {
	EventType        string  `json:"e" bson:"e"`
	EventTime        int64   `json:"E" bson:"E"`
	Asset            string  `json:"a" bson:"a"`
	SignBalanceDelta float64 `json:"d,string" bson:"d,string"`
	ClearTime        int64   `json:"T" bson:"T"`
}

func (w WsSpotBalanceUpdate) AbsBalanceDelta() float64 {
	return math.Abs(w.SignBalanceDelta)
}

type WsSpotListStatusObject struct {
	Symbol        string `json:"s" bson:"s"`
	OrderId       int64  `json:"i" bson:"i"`
	ClientOrderId string `json:"c" bson:"c"`
}

type WsSpotListStatus struct {
	EventType         string                   `json:"e" bson:"e"`
	EventTime         int64                    `json:"E" bson:"E"`
	Symbol            string                   `json:"s" bson:"s"`
	OrderListId       int                      `json:"g" bson:"g"`
	ContingencyType   string                   `json:"c" bson:"c"`
	ListStatusType    string                   `json:"l" bson:"l"`
	ListOrderStatus   string                   `json:"L" bson:"L"`
	ListRejectReason  string                   `json:"r" bson:"r"`
	ListClientOrderId string                   `json:"C" bson:"C"`
	Time              int64                    `json:"T" bson:"T"`
	Objects           []WsSpotListStatusObject `json:"O" bson:"O"`
}

type WsAggTradeStream struct {
	EventType    WsEvent `json:"e" bson:"e"`
	EventTime    int64   `json:"E" bson:"E"`
	Symbol       string  `json:"s" bson:"s"`
	AggTradeId   int64   `json:"a" bson:"a"`
	Price        float64 `json:"p,string" bson:"p,string"`
	Qty          float64 `json:"q,string" bson:"q,string"`
	FirstTradeId int64   `json:"f" bson:"f"`
	LastTradeId  int64   `json:"l" bson:"l"`
	TradeTime    int64   `json:"T" bson:"T"`
	IsMaker      bool    `json:"m" bson:"m"`
	M            bool    `json:"M" bson:"M"` // ignore
}

type WsTradeStream struct {
	EventType    WsEvent `json:"e" bson:"e"`
	EventTime    int64   `json:"E" bson:"E"`
	Symbol       string  `json:"s" bson:"s"`
	TradeID      int64   `json:"t" bson:"t"`
	Price        float64 `json:"p,string" bson:"p,string"`
	Quantity     float64 `json:"q,string" bson:"q,string"`
	TradeTime    int64   `json:"T" bson:"T"`
	IsBuyerMaker bool    `json:"m" bson:"m"`
	M            bool    `json:"M" bson:"M"` // ignore

	//
	BuyerOrderID  int64 `json:"b" bson:"b"`
	SellerOrderID int64 `json:"a" bson:"a"`
}

type WsKline struct {
	StartTime                int64         `json:"t" bson:"t"`
	CloseTime                int64         `json:"T" bson:"T"`
	Symbol                   string        `json:"s" bson:"s"`
	Interval                 KlineInterval `json:"i" bson:"i"`
	FirstTradeId             int64         `json:"f" bson:"f"`
	LastTradeId              int64         `json:"L" bson:"L"`
	OpenPrice                float64       `json:"o,string" bson:"o,string"`
	ClosePrice               float64       `json:"c,string" bson:"c,string"`
	HighPrice                float64       `json:"h,string" bson:"h,string"`
	LowPrice                 float64       `json:"l,string" bson:"l,string"`
	BaseAssetVolume          float64       `json:"v,string" bson:"v,string"`
	TradesNum                int64         `json:"n" bson:"n"`
	IsKlineClosed            bool          `json:"x" bson:"x"`
	QuoteVolume              float64       `json:"q,string" bson:"q,string"`
	BaseAssetTakerBuyVolume  float64       `json:"V,string" bson:"V,string"`
	QuoteAssetTakerBuyVolume float64       `json:"Q,string" bson:"Q,string"`
	B                        string        `json:"B" bson:"B"` // ignore
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
	EventType WsEvent `json:"e" bson:"e"`
	EventTime int64   `json:"E" bson:"E"`
	Symbol    string  `json:"s" bson:"s"`
	Kline     WsKline `json:"k" bson:"k"`
}

type WsOrderExecutionReport struct {
	EventType               string             `json:"e" bson:"e"`
	EventTime               int64              `json:"E" bson:"E"`
	Symbol                  string             `json:"s" bson:"s"`
	ClientOrderId           string             `json:"c" bson:"c"`
	Side                    OrderSide          `json:"S" bson:"S"`
	Type                    OrderType          `json:"o" bson:"o"`
	TimeInForce             TimeInForce        `json:"f" bson:"f"`
	Qty                     float64            `json:"q,string" bson:"q,string"`
	Price                   float64            `json:"p,string" bson:"p,string"`
	StopPrice               float64            `json:"P,string" bson:"P,string"`
	IcebergQty              float64            `json:"F,string" bson:"F,string"`
	OrderListId             int64              `json:"g" bson:"g"`
	OriginalClientId        string             `json:"C" bson:"C"`
	ExecutionType           OrderExecutionType `json:"x" bson:"x"`
	Status                  OrderStatus        `json:"X" bson:"X"`
	RejectReason            string             `json:"r" bson:"r"`
	OrderId                 int64              `json:"i" bson:"i"`
	LastExecutedQty         float64            `json:"l,string" bson:"l,string"`
	FilledQty               float64            `json:"z,string" bson:"z,string"`
	LastExecutedPrice       float64            `json:"L,string" bson:"L,string"`
	CommissionAmt           float64            `json:"n,string" bson:"n,string"`
	CommissionAsset         string             `json:"N" bson:"N"`
	Time                    int64              `json:"T" bson:"T"`
	TradeId                 int64              `json:"t" bson:"t"`
	PreventedMatchId        int64              `json:"v" bson:"v"`
	Ignore                  int64              `json:"I" bson:"I"`
	IsOrderOnTheBook        bool               `json:"w" bson:"w"`
	IsMaker                 bool               `json:"m" bson:"m"`
	Ignore1                 bool               `json:"M" bson:"M"`
	CreationTime            int64              `json:"O" bson:"O"`
	FilledQuote             float64            `json:"Z,string" bson:"Z,string"`
	LastExecutedQuote       float64            `json:"Y,string" bson:"Y,string"`
	QuoteOrderQty           float64            `json:"Q,string" bson:"Q,string"`
	WorkingTime             int64              `json:"W" bson:"W"`
	SelfTradePreventionMode STPMode            `json:"V" bson:"V"`

	// just for margin order
	TrailingDelta      float64 `json:"d" bson:"d"`       // Trailing Delta; This is only visible if the order was a trailing stop order.
	TrailingTime       int64   `json:"Data" bson:"Data"` // Trailing Time; This is only visible if the trailing stop order has been activated.
	MarginStrategyId   int64   `json:"j" bson:"j"`
	MarginStrategyType int64   `json:"J" bson:"J"`
	TradeGroupId       int64   `json:"u" bson:"u"`
	CounterOrderId     int64   `json:"U" bson:"U"`
	PreventedQty       float64 `json:"A,string" bson:"A,string"`
	LastPreventedQty   float64 `json:"B,string" bson:"B,string"`

	// just for futures order
	AvgPrice            string       `json:"ap" bson:"ap"`
	Sp                  string       `json:"sp" bson:"sp"` // ignore
	BidNotional         string       `json:"b" bson:"b"`
	AskNotional         string       `json:"a" bson:"a"`
	IsReduceOnly        bool         `json:"R" bson:"R"`
	PositionSide        PositionSide `json:"ps" bson:"ps"`
	RealizedProfit      string       `json:"rp" bson:"rp"`
	FuturesStrategyType string       `json:"st" bson:"st"`
	FuturesStrategyId   int64        `json:"si" bson:"si"`
	Gtd                 int64        `json:"gtd" bson:"gtd"`
}

type WsListenKeyExpired struct {
	EventType string `json:"e" bson:"e"`
	EventTime int64  `json:"E" bson:"E"`
	ListenKey string `json:"listenKey" bson:"listenKey"`
}

// WsCMIndexPriceStream
// url: wss://dstream.binance.com/ws
// event: WsEventIndexPriceUpdate
type WsCMIndexPriceStream struct {
	EventType WsEvent `json:"e" bson:"e"`
	EventTime int64   `json:"E" bson:"E"`
	Pair      string  `json:"i" bson:"i"`
	Price     float64 `json:"p,string" bson:"p,string"`
}

// WsMarkPriceStream
// url: wss://dstream.binance.com/ws
// event: WsEventMarkPriceUpdate
type WsMarkPriceStream struct {
	EventType            WsEvent `json:"e" bson:"e"`
	EventTime            int64   `json:"E" bson:"E"`
	Symbol               string  `json:"s" bson:"s"`
	MarkPrice            string  `json:"p" bson:"p"`
	EstimatedSettlePrice string  `json:"P" bson:"P"`
	IndexPrice           string  `json:"i" bson:"i"`
	FundingRate          string  `json:"r" bson:"r"`
	NextFundingTime      int64   `json:"T" bson:"T"`
}

type WsStrategyOrder struct {
	Symbol               string       `json:"s" bson:"s"`
	ClientOrderId        string       `json:"c" bson:"c"`
	Id                   int64        `json:"si" bson:"si"`
	Side                 OrderSide    `json:"S" bson:"S"`
	Type                 string       `json:"st" bson:"st"`
	TimeInForce          TimeInForce  `json:"f" bson:"f"`
	Qty                  float64      `json:"q,string" bson:"q,string"`
	Price                float64      `json:"p,string" bson:"p,string"`
	StopPrice            float64      `json:"sp,string" bson:"sp,string"`
	Status               OrderStatus  `json:"os" bson:"os"`
	OrderBookTime        int64        `json:"T" bson:"T"`
	UpdateTime           int64        `json:"ut" bson:"ut"`
	IsThisReduceOnly     bool         `json:"R" bson:"R"`
	StopPriceWorkingType string       `json:"wt" bson:"wt"`
	PositionSide         PositionSide `json:"ps" bson:"ps"`
	Cp                   bool         `json:"cp" bson:"cp"` // If Close-All, pushed with conditional order
	ActivationPrice      float64      `json:"AP,string" bson:"AP,string"`
	CallbackRate         float64      `json:"cr,string" bson:"cr,string"`
	OrderId              int64        `json:"i" bson:"i"`
	STPMode              string       `json:"V" bson:"V"`
	Gtd                  int64        `json:"gtd" bson:"gtd"`
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
	EventType       WsEvent         `json:"e" bson:"e"`
	EventTime       int64           `json:"E" bson:"E"`
	TransactionTime int64           `json:"T" bson:"T"`
	BusinessUnit    WsBusinessUnit  `json:"fs" bson:"fs"`
	StrategyOrder   WsStrategyOrder `json:"so" bson:"so"`
}

type WsOpenOrderLossUpdate struct {
	Asset      string  `json:"a" bson:"a"`
	SignAmount float64 `json:"o" bson:"o"`
}

// WsOpenOrderLoss
// url: https://papi.binance.com/ws
// event: WsEventOpenOrderLoss
type WsOpenOrderLoss struct {
	EventType WsEvent                 `json:"e" bson:"e"`
	EventTime int64                   `json:"E" bson:"E"`
	OrderLoss []WsOpenOrderLossUpdate `json:"O" bson:"O"`
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
	EventType      WsEvent       `json:"e" bson:"e"`
	EventTime      int64         `json:"E" bson:"E"`
	Asset          string        `json:"a" bson:"a"`
	Type           LiabilityType `json:"t" bson:"t"`
	TxId           int64         `json:"tx" bson:"tx"`
	Principal      float64       `json:"p,string" bson:"p,string"`
	Interest       float64       `json:"i,string" bson:"i,string"`
	TotalLiability float64       `json:"l,string" bson:"l,string"`
}

// WsPmFuturesOrderTradeUpdate
// url: https://papi.binance.com/ws
// event: WsEventOrderTradeUpdate
type WsPmFuturesOrderTradeUpdate struct {
	EventType       WsEvent                `json:"e" bson:"e"`
	EventTime       int64                  `json:"E" bson:"E"`
	TransactionTime int64                  `json:"T" bson:"T"`
	BusinessUnit    WsBusinessUnit         `json:"fs" bson:"fs"`
	AccountAlias    string                 `json:"i" bson:"i"` // Account Alias,ignore for UM
	Order           WsOrderExecutionReport `json:"o" bson:"o"`
}

type WsFuturesWalletBalance struct {
	Asset              string  `json:"a" bson:"a"`
	WalletBalance      float64 `json:"wb,string" bson:"wb,string"`
	CrossWalletBalance float64 `json:"cw,string" bson:"cw,string"`
	BalanceChange      float64 `json:"bc,string" bson:"bc,string"`
}

type WsFuturesPosition struct {
	Symbol                    string       `json:"s" bson:"s"`
	SignPositionQty           float64      `json:"pa" bson:"pa"`
	EntryPrice                float64      `json:"ep" bson:"ep"`
	PreFeeAccumulatedRealized float64      `json:"cr" bson:"cr"`
	UnrealizedPnl             float64      `json:"up" bson:"up"`
	PositionSide              PositionSide `json:"ps" bson:"ps"`
	BreakevenPrice            string       `json:"bep" bson:"bep"`
}

type WsPmFuturesAcctUpdateData struct {
	Reason    string                   `json:"m" bson:"m"`
	Balances  []WsFuturesWalletBalance `json:"B" bson:"B"`
	Positions []WsFuturesPosition      `json:"P" bson:"P"`
}

// WsPmFuturesAcctUpdateStream
// url: https://papi.binance.com/ws
// event: WsEventAccountUpdate
type WsPmFuturesAcctUpdateStream struct {
	EventType       WsEvent                   `json:"e" bson:"e"`
	EventTime       int64                     `json:"E" bson:"E"`
	TransactionTime int64                     `json:"T" bson:"T"`
	BusinessUnit    WsBusinessUnit            `json:"fs" bson:"fs"`
	AccountAlias    string                    `json:"i" bson:"i"` // Account Alias,ignore for UM
	Data            WsPmFuturesAcctUpdateData `json:"a" bson:"a"`
}

type WsFuturesCfg struct {
	Symbol   string  `json:"s" bson:"s"`
	Leverage float64 `json:"l" bson:"l"`
}

// WsFuturesAccountCfgUpdate
// url: https://papi.binance.com/ws
// event: WsEventAccountConfigUpdate
type WsFuturesAccountCfgUpdate struct {
	EventType       WsEvent        `json:"e" bson:"e"`
	EventTime       int64          `json:"E" bson:"E"`
	TransactionTime int64          `json:"T" bson:"T"`
	BusinessUnit    WsBusinessUnit `json:"fs" bson:"fs"`
	Cfg             WsFuturesCfg   `json:"ac" bson:"ac"`
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
	Asset        string  `json:"a" bson:"a"`
	BalanceDelta float64 `json:"d,string" bson:"d,string"`
	UpdateId     int64   `json:"U" bson:"U"`
	ClearTime    int64   `json:"T" bson:"T"`
}

type WsLiquidationOrder struct {
	Symbol        string      `json:"s" bson:"s"`
	OrderSide     OrderSide   `json:"S" bson:"S"`
	OrderType     OrderType   `json:"o" bson:"o"`
	TimeInForce   TimeInForce `json:"f" bson:"f"`
	OriQty        float64     `json:"q,string" bson:"q,string"`
	Price         float64     `json:"p,string" bson:"p,string"`
	AvgPrice      float64     `json:"ap,string" bson:"ap,string"`
	OrderStatus   OrderStatus `json:"X" bson:"X"`
	LastFilledQty float64     `json:"l,string" bson:"l,string"`
	FilledQty     float64     `json:"z,string" bson:"z,string"`
	TradeTime     int64       `json:"T" bson:"T"`
}

type WsLiquidationOrderStream struct {
	EventType WsEvent            `json:"e" bson:"e"`
	EventTime int64              `json:"E" bson:"E"`
	Order     WsLiquidationOrder `json:"o" bson:"o"`
}
