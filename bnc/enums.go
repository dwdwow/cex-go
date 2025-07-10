package bnc

type SymbolStatus string

const (
	TRADING    SymbolStatus = "TRADING"
	END_OF_DAY SymbolStatus = "END_OF_DAY"
	HALT       SymbolStatus = "HALT"
	BREAK      SymbolStatus = "BREAK"
)

type AcctSybPermission string

const (
	SPOT        AcctSybPermission = "SPOT"
	MARGIN      AcctSybPermission = "MARGIN"
	LEVERAGED   AcctSybPermission = "LEVERAGED"
	TRD_GRP_002 AcctSybPermission = "TRD_GRP_002"
	TRD_GRP_003 AcctSybPermission = "TRD_GRP_003"
	TRD_GRP_004 AcctSybPermission = "TRD_GRP_004"
	TRD_GRP_005 AcctSybPermission = "TRD_GRP_005"
	TRD_GRP_006 AcctSybPermission = "TRD_GRP_006"
	TRD_GRP_007 AcctSybPermission = "TRD_GRP_007"
	TRD_GRP_008 AcctSybPermission = "TRD_GRP_008"
	TRD_GRP_009 AcctSybPermission = "TRD_GRP_009"
	TRD_GRP_010 AcctSybPermission = "TRD_GRP_010"
	TRD_GRP_011 AcctSybPermission = "TRD_GRP_011"
	TRD_GRP_012 AcctSybPermission = "TRD_GRP_012"
	TRD_GRP_013 AcctSybPermission = "TRD_GRP_013"
	TRD_GRP_014 AcctSybPermission = "TRD_GRP_014"
	TRD_GRP_015 AcctSybPermission = "TRD_GRP_015"
	TRD_GRP_016 AcctSybPermission = "TRD_GRP_016"
	TRD_GRP_017 AcctSybPermission = "TRD_GRP_017"
	TRD_GRP_018 AcctSybPermission = "TRD_GRP_018"
	TRD_GRP_019 AcctSybPermission = "TRD_GRP_019"
	TRD_GRP_020 AcctSybPermission = "TRD_GRP_020"
	TRD_GRP_021 AcctSybPermission = "TRD_GRP_021"
	TRD_GRP_022 AcctSybPermission = "TRD_GRP_022"
	TRD_GRP_023 AcctSybPermission = "TRD_GRP_023"
	TRD_GRP_024 AcctSybPermission = "TRD_GRP_024"
	TRD_GRP_025 AcctSybPermission = "TRD_GRP_025"
)

type OrderStatus string

const (
	NEW              OrderStatus = "NEW"
	PENDING_NEW      OrderStatus = "PENDING_NEW"
	PARTIALLY_FILLED OrderStatus = "PARTIALLY_FILLED"
	FILLED           OrderStatus = "FILLED"
	CANCELED         OrderStatus = "CANCELED"
	PENDING_CANCEL   OrderStatus = "PENDING_CANCEL"
	REJECTED         OrderStatus = "REJECTED"
	EXPIRED          OrderStatus = "EXPIRED"
	EXPIRED_IN_MATCH OrderStatus = "EXPIRED_IN_MATCH"
)

type ListStatusType string

const (
	RESPONSE     ListStatusType = "RESPONSE"
	EXEC_STARTED ListStatusType = "EXEC_STARTED"
	UPDATED      ListStatusType = "UPDATED"
	ALL_DONE     ListStatusType = "ALL_DONE"
)

type ListOrderStatus string

const (
	EXECUTING     ListOrderStatus = "EXECUTING"
	LIST_ALL_DONE ListOrderStatus = "ALL_DONE"
	REJECT        ListOrderStatus = "REJECT"
)

type ContingencyType string

const (
	OCO ContingencyType = "OCO"
	OTO ContingencyType = "OTO"
)

type AllocationType string

const (
	ALLOCATION_SOR AllocationType = "SOR"
)

type OrderType string

const (
	LIMIT             OrderType = "LIMIT"
	MARKET            OrderType = "MARKET"
	STOP_LOSS         OrderType = "STOP_LOSS"
	STOP_LOSS_LIMIT   OrderType = "STOP_LOSS_LIMIT"
	TAKE_PROFIT       OrderType = "TAKE_PROFIT"
	TAKE_PROFIT_LIMIT OrderType = "TAKE_PROFIT_LIMIT"
	LIMIT_MAKER       OrderType = "LIMIT_MAKER"
)

type NewOrderRespType string

const (
	ACK    NewOrderRespType = "ACK"
	RESULT NewOrderRespType = "RESULT"
	FULL   NewOrderRespType = "FULL"
)

type WorkingFloor string

const (
	EXCHANGE          WorkingFloor = "EXCHANGE"
	WORKING_FLOOR_SOR WorkingFloor = "SOR"
)

type OrderSide string

const (
	BUY  OrderSide = "BUY"
	SELL OrderSide = "SELL"
)

type TimeInForce string

const (
	GTC TimeInForce = "GTC"
	IOC TimeInForce = "IOC"
	FOK TimeInForce = "FOK"
)

type RateLimitType string

const (
	REQUEST_WEIGHT RateLimitType = "REQUEST_WEIGHT"
	ORDERS         RateLimitType = "ORDERS"
	RAW_REQUESTS   RateLimitType = "RAW_REQUESTS"
)

type RateLimiterInterval string

const (
	SECOND RateLimiterInterval = "SECOND"
	MINUTE RateLimiterInterval = "MINUTE"
	DAY    RateLimiterInterval = "DAY"
)

type STPMode string

const (
	NONE         STPMode = "NONE"
	EXPIRE_MAKER STPMode = "EXPIRE_MAKER"
	EXPIRE_TAKER STPMode = "EXPIRE_TAKER"
	EXPIRE_BOTH  STPMode = "EXPIRE_BOTH"
	DECREMENT    STPMode = "DECREMENT"
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
