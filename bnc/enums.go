package bnc

type SymbolStatus string

const (
	// Common
	SYMBOL_STATUS_TRADING SymbolStatus = "TRADING"

	// Spot
	SYMBOL_STATUS_END_OF_DAY SymbolStatus = "END_OF_DAY"
	SYMBOL_STATUS_HALT       SymbolStatus = "HALT"
	SYMBOL_STATUS_BREAK      SymbolStatus = "BREAK"

	// Futures
	SYMBOL_STATUS_PENDING_TRADING SymbolStatus = "PENDING_TRADING"
	SYMBOL_STATUS_PRE_DELIVERING  SymbolStatus = "PRE_DELIVERING"
	SYMBOL_STATUS_DELIVERING      SymbolStatus = "DELIVERING"
	SYMBOL_STATUS_DELIVERED       SymbolStatus = "DELIVERED"
	SYMBOL_STATUS_PRE_SETTLE      SymbolStatus = "PRE_SETTLE"
	SYMBOL_STATUS_SETTLING        SymbolStatus = "SETTLING"
	SYMBOL_STATUS_CLOSE           SymbolStatus = "CLOSE"
)

type AcctSybPermission string

const (
	PERMISSION_SPOT        AcctSybPermission = "SPOT"
	PERMISSION_MARGIN      AcctSybPermission = "MARGIN"
	PERMISSION_LEVERAGED   AcctSybPermission = "LEVERAGED"
	PERMISSION_TRD_GRP_002 AcctSybPermission = "TRD_GRP_002"
	PERMISSION_TRD_GRP_003 AcctSybPermission = "TRD_GRP_003"
	PERMISSION_TRD_GRP_004 AcctSybPermission = "TRD_GRP_004"
	PERMISSION_TRD_GRP_005 AcctSybPermission = "TRD_GRP_005"
	PERMISSION_TRD_GRP_006 AcctSybPermission = "TRD_GRP_006"
	PERMISSION_TRD_GRP_007 AcctSybPermission = "TRD_GRP_007"
	PERMISSION_TRD_GRP_008 AcctSybPermission = "TRD_GRP_008"
	PERMISSION_TRD_GRP_009 AcctSybPermission = "TRD_GRP_009"
	PERMISSION_TRD_GRP_010 AcctSybPermission = "TRD_GRP_010"
	PERMISSION_TRD_GRP_011 AcctSybPermission = "TRD_GRP_011"
	PERMISSION_TRD_GRP_012 AcctSybPermission = "TRD_GRP_012"
	PERMISSION_TRD_GRP_013 AcctSybPermission = "TRD_GRP_013"
	PERMISSION_TRD_GRP_014 AcctSybPermission = "TRD_GRP_014"
	PERMISSION_TRD_GRP_015 AcctSybPermission = "TRD_GRP_015"
	PERMISSION_TRD_GRP_016 AcctSybPermission = "TRD_GRP_016"
	PERMISSION_TRD_GRP_017 AcctSybPermission = "TRD_GRP_017"
	PERMISSION_TRD_GRP_018 AcctSybPermission = "TRD_GRP_018"
	PERMISSION_TRD_GRP_019 AcctSybPermission = "TRD_GRP_019"
	PERMISSION_TRD_GRP_020 AcctSybPermission = "TRD_GRP_020"
	PERMISSION_TRD_GRP_021 AcctSybPermission = "TRD_GRP_021"
	PERMISSION_TRD_GRP_022 AcctSybPermission = "TRD_GRP_022"
	PERMISSION_TRD_GRP_023 AcctSybPermission = "TRD_GRP_023"
)

type OrderStatus string

const (
	ORDER_STATUS_NEW              OrderStatus = "NEW"
	ORDER_STATUS_PENDING_NEW      OrderStatus = "PENDING_NEW"
	ORDER_STATUS_PARTIALLY_FILLED OrderStatus = "PARTIALLY_FILLED"
	ORDER_STATUS_FILLED           OrderStatus = "FILLED"
	ORDER_STATUS_CANCELED         OrderStatus = "CANCELED"
	ORDER_STATUS_PENDING_CANCEL   OrderStatus = "PENDING_CANCEL"
	ORDER_STATUS_REJECTED         OrderStatus = "REJECTED"
	ORDER_STATUS_EXPIRED          OrderStatus = "EXPIRED"
	ORDER_STATUS_EXPIRED_IN_MATCH OrderStatus = "EXPIRED_IN_MATCH"
)

type ListStatusType string

const (
	LIST_STATUS_TYPE_RESPONSE     ListStatusType = "RESPONSE"
	LIST_STATUS_TYPE_EXEC_STARTED ListStatusType = "EXEC_STARTED"
	LIST_STATUS_TYPE_UPDATED      ListStatusType = "UPDATED"
	LIST_STATUS_TYPE_ALL_DONE     ListStatusType = "ALL_DONE"
)

type ListOrderStatus string

const (
	LIST_ORDER_STATUS_EXECUTING ListOrderStatus = "EXECUTING"
	LIST_ORDER_STATUS_ALL_DONE  ListOrderStatus = "ALL_DONE"
	LIST_ORDER_STATUS_REJECT    ListOrderStatus = "REJECT"
)

type ContingencyType string

const (
	CONINGENCY_TYPE_OCO ContingencyType = "OCO"
	CONINGENCY_TYPE_OTO ContingencyType = "OTO"
)

type AllocationType string

const (
	ALLOCATION_TYPE_SOR AllocationType = "SOR"
)

type OrderType string

const (
	// Common
	ORDER_TYPE_LIMIT       OrderType = "LIMIT"
	ORDER_TYPE_MARKET      OrderType = "MARKET"
	ORDER_TYPE_TAKE_PROFIT OrderType = "TAKE_PROFIT"

	// Spot
	ORDER_TYPE_STOP_LOSS         OrderType = "STOP_LOSS"
	ORDER_TYPE_STOP_LOSS_LIMIT   OrderType = "STOP_LOSS_LIMIT"
	ORDER_TYPE_TAKE_PROFIT_LIMIT OrderType = "TAKE_PROFIT_LIMIT"
	ORDER_TYPE_LIMIT_MAKER       OrderType = "LIMIT_MAKER"

	// Futures
	ORDER_TYPE_STOP                 OrderType = "STOP"
	ORDER_TYPE_STOP_MARKET          OrderType = "STOP_MARKET"
	ORDER_TYPE_TAKE_PROFIT_MARKET   OrderType = "TAKE_PROFIT_MARKET"
	ORDER_TYPE_TRAILING_STOP_MARKET OrderType = "TRAILING_STOP_MARKET"
)

type NewOrderRespType string

const (
	// Common
	NEW_ORDER_RESP_TYPE_ACK    NewOrderRespType = "ACK"
	NEW_ORDER_RESP_TYPE_RESULT NewOrderRespType = "RESULT"

	// Spot
	NEW_ORDER_RESP_TYPE_FULL NewOrderRespType = "FULL"
)

type WorkingFloor string

const (
	WORKING_FLOOR_EXCHANGE WorkingFloor = "EXCHANGE"
	WORKING_FLOOR_SOR      WorkingFloor = "SOR"
)

type OrderSide string

const (
	ORDER_SIDE_BUY  OrderSide = "BUY"
	ORDER_SIDE_SELL OrderSide = "SELL"
)

type TimeInForce string

const (
	TIME_IN_FORCE_GTC TimeInForce = "GTC"
	TIME_IN_FORCE_IOC TimeInForce = "IOC"
	TIME_IN_FORCE_FOK TimeInForce = "FOK"
)

type RateLimitType string

const (
	REQUEST_WEIGHT_TYPE_REQUEST_WEIGHT RateLimitType = "REQUEST_WEIGHT"
	REQUEST_WEIGHT_TYPE_ORDERS         RateLimitType = "ORDERS"
	REQUEST_WEIGHT_TYPE_RAW_REQUESTS   RateLimitType = "RAW_REQUESTS"
)

type RateLimiterInterval string

const (
	REQUEST_WEIGHT_INTERVAL_SECOND RateLimiterInterval = "SECOND"
	REQUEST_WEIGHT_INTERVAL_MINUTE RateLimiterInterval = "MINUTE"
	REQUEST_WEIGHT_INTERVAL_DAY    RateLimiterInterval = "DAY"
)

type STPMode string

const (
	STP_MODE_NONE         STPMode = "NONE"
	STP_MODE_EXPIRE_MAKER STPMode = "EXPIRE_MAKER"
	STP_MODE_EXPIRE_TAKER STPMode = "EXPIRE_TAKER"
	STP_MODE_EXPIRE_BOTH  STPMode = "EXPIRE_BOTH"
	STP_MODE_DECREMENT    STPMode = "DECREMENT"
)

type KlineInterval string

const (
	KLINE_INTERVAL_1s  KlineInterval = "1s"
	KLINE_INTERVAL_1m  KlineInterval = "1m"
	KLINE_INTERVAL_3m  KlineInterval = "3m"
	KLINE_INTERVAL_5m  KlineInterval = "5m"
	KLINE_INTERVAL_15m KlineInterval = "15m"
	KLINE_INTERVAL_30m KlineInterval = "30m"
	KLINE_INTERVAL_1h  KlineInterval = "1h"
	KLINE_INTERVAL_2h  KlineInterval = "2h"
	KLINE_INTERVAL_4h  KlineInterval = "4h"
	KLINE_INTERVAL_6h  KlineInterval = "6h"
	KLINE_INTERVAL_8h  KlineInterval = "8h"
	KLINE_INTERVAL_12h KlineInterval = "12h"
	KLINE_INTERVAL_1d  KlineInterval = "1d"
	KLINE_INTERVAL_3d  KlineInterval = "3d"
	KLINE_INTERVAL_1w  KlineInterval = "1w"
	KLINE_INTERVAL_1M  KlineInterval = "1M"
)

type Symbol string

const (
	SYMBOL_TYPE_FUTURE Symbol = "FUTURE"
)

type ContractType string

const (
	CONTRACT_TYPE_PERPETUAL            ContractType = "PERPETUAL"
	CONTRACT_TYPE_CURRENT_MONTH        ContractType = "CURRENT_MONTH"
	CONTRACT_TYPE_NEXT_MONTH           ContractType = "NEXT_MONTH"
	CONTRACT_TYPE_CURRENT_QUARTER      ContractType = "CURRENT_QUARTER"
	CONTRACT_TYPE_NEXT_QUARTER         ContractType = "NEXT_QUARTER"
	CONTRACT_TYPE_PERPETUAL_DELIVERING ContractType = "PERPETUAL_DELIVERING"
)

type PositionSide string

const (
	POSITION_SIDE_BOTH  PositionSide = "BOTH"
	POSITION_SIDE_LONG  PositionSide = "LONG"
	POSITION_SIDE_SHORT PositionSide = "SHORT"
)

type WorkingType string

const (
	WORKING_TYPE_MARK_PRICE     WorkingType = "MARK_PRICE"
	WORKING_TYPE_CONTRACT_PRICE WorkingType = "CONTRACT_PRICE"
)

type PriceMatch string

const (
	PRICE_MATCH_NONE        PriceMatch = "NONE"
	PRICE_MATCH_OPPONENT    PriceMatch = "OPPONENT"
	PRICE_MATCH_OPPONENT_5  PriceMatch = "OPPONENT_5"
	PRICE_MATCH_OPPONENT_10 PriceMatch = "OPPONENT_10"
	PRICE_MATCH_OPPONENT_20 PriceMatch = "OPPONENT_20"
	PRICE_MATCH_QUEUE       PriceMatch = "QUEUE"
	PRICE_MATCH_QUEUE_5     PriceMatch = "QUEUE_5"
	PRICE_MATCH_QUEUE_10    PriceMatch = "QUEUE_10"
	PRICE_MATCH_QUEUE_20    PriceMatch = "QUEUE_20"
)

type OrderExecutionType string

const (
	OrderExecutionTypeNew         OrderExecutionType = "NEW"
	OrderExecutionTypeCanceled    OrderExecutionType = "CANCELED"
	OrderExecutionTypeReplaced    OrderExecutionType = "REPLACED"
	OrderExecutionRejected        OrderExecutionType = "REJECTED"
	OrderExecutionTrade           OrderExecutionType = "TRADE"
	OrderExecutionExpired         OrderExecutionType = "EXPIRED"
	OrderExecutionTradePrevention OrderExecutionType = "TRADE_PREVENTION"
)
