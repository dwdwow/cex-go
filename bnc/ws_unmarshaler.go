package bnc

import (
	"encoding/json"
	"fmt"
)

func unmarshal[T any](data []byte) (t T, err error) {
	err = json.Unmarshal(data, &t)
	return
}

func SpotWsPrivateMsgUnmarshaler(e WsEvent, _ bool, data []byte) (any, error) {
	switch e {
	case WsEventOutboundAccountPosition:
		return unmarshal[WsSpotAccountUpdate](data)
	case WsEventBalanceUpdate:
		return unmarshal[WsSpotBalanceUpdate](data)
	case WsEventExecutionReport:
		return unmarshal[WsOrderExecutionReport](data)
	case WsEventListStatus:
		return unmarshal[WsSpotListStatus](data)
	case WsEventListenKeyExpired:
		return unmarshal[WsListenKeyExpired](data)
	default:
		return nil, fmt.Errorf("bnc: unknown event %v", e)
	}
}

func SpotWsPublicMsgUnmarshaler(e WsEvent, _ bool, data []byte) (any, error) {
	switch e {
	case WsEventAggTrade:
		return unmarshal[WsAggTradeStream](data)
	case WsEventTrade:
		return unmarshal[WsTradeStream](data)
	case WsEventKline:
		return unmarshal[WsKlineStream](data)
	case WsEventDepthUpdate:
		return unmarshal[WsDepthStream](data)
	default:
		return nil, fmt.Errorf("bnc: unknown event %v", e)
	}
}

func UmFuturesWsPublicMsgUnmarshaler(e WsEvent, isArray bool, data []byte) (any, error) {
	switch e {
	case WsEventAggTrade:
		return unmarshal[WsAggTradeStream](data)
	case WsEventMarkPriceUpdate:
		if isArray {
			return unmarshal[[]WsMarkPriceStream](data)
		}
		return unmarshal[WsMarkPriceStream](data)
	case WsEventForceOrder:
		if isArray {
			return unmarshal[[]WsLiquidationOrderStream](data)
		}
		return unmarshal[WsLiquidationOrderStream](data)
	case WsEventKline:
		return unmarshal[WsKlineStream](data)
	case WsEventDepthUpdate:
		return unmarshal[WsDepthStream](data)
	default:
		return nil, fmt.Errorf("bnc: unknown event %v", e)
	}
}

func CmFuturesWsPublicMsgUnmarshaler(e WsEvent, isArray bool, data []byte) (any, error) {
	switch e {
	case WsEventAggTrade:
		return unmarshal[WsAggTradeStream](data)
	case WsEventIndexPriceUpdate:
		return unmarshal[WsCMIndexPriceStream](data)
	case WsEventMarkPriceUpdate:
		if isArray {
			return unmarshal[[]WsMarkPriceStream](data)
		}
		return unmarshal[WsMarkPriceStream](data)
	case WsEventForceOrder:
		if isArray {
			return unmarshal[[]WsLiquidationOrderStream](data)
		}
		return unmarshal[WsLiquidationOrderStream](data)
	case WsEventKline:
		return unmarshal[WsKlineStream](data)
	case WsEventDepthUpdate:
		return unmarshal[WsDepthStream](data)
	default:
		return nil, fmt.Errorf("bnc: unknown event %v", e)
	}
}
