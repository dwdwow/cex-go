package bnc

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dwdwow/cex-go"
)

type WsDataUnmarshaler func(e WsEvent, isArray bool, data []byte) (any, error)

func unmarshal[T any](data []byte) (t T, err error) {
	err = json.Unmarshal(data, &t)
	return
}

func spotWsPrivateMsgUnmarshaler(e WsEvent, _ bool, data []byte) (any, error) {
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

func spotWsPublicMsgUnmarshaler(e WsEvent, _ bool, data []byte) (any, error) {
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

func umFuturesWsPublicMsgUnmarshaler(e WsEvent, isArray bool, data []byte) (any, error) {
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

func cmFuturesWsPublicMsgUnmarshaler(e WsEvent, isArray bool, data []byte) (any, error) {
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

type WsCfg struct {
	// ws url without streams and auth tokens
	Url          string
	ListenKeyUrl string

	// max stream per websocket client
	// ex. spot 1024, portfolio margin / futures 200
	// normally just for public websocket
	MaxStream int

	// channel capacity
	ChCap int

	// binance has incoming massage limitation
	// ex. spot 5/s, futures 10/s
	ReqDur       time.Duration
	MaxReqPerDur int

	APIKey       string
	APISecretKey string

	// just use in WsClient
	DataUnmarshaler WsDataUnmarshaler
}

var spotPublicWsCfg = WsCfg{
	Url:             WsBaseUrl,
	MaxStream:       1024,
	ChCap:           maxWsChCap,
	ReqDur:          time.Second,
	MaxReqPerDur:    5,
	DataUnmarshaler: spotWsPublicMsgUnmarshaler,
}

func DefaultSpotPublicWsCfg() WsCfg {
	return spotPublicWsCfg
}

var spotPrivateWsCfg = WsCfg{
	Url:             WsBaseUrl,
	ListenKeyUrl:    API_ENDPOINT + API_V3 + "/userDataStream",
	MaxStream:       1024,
	ChCap:           maxWsChCap,
	ReqDur:          time.Second,
	MaxReqPerDur:    10,
	DataUnmarshaler: spotWsPrivateMsgUnmarshaler,
}

func DefaultSpotPrivateWsCfg() WsCfg {
	return spotPrivateWsCfg
}

var umPublicWsCfg = WsCfg{
	Url:             FutureWsBaseUrl,
	MaxStream:       200,
	ChCap:           maxWsChCap,
	ReqDur:          time.Second,
	MaxReqPerDur:    10,
	DataUnmarshaler: umFuturesWsPublicMsgUnmarshaler,
}

func DefaultUmPublicWsCfg() WsCfg {
	return umPublicWsCfg
}

var cmPublicWsCfg = WsCfg{
	Url:             CMFutureWsBaseUrl,
	MaxStream:       200,
	ChCap:           maxWsChCap,
	ReqDur:          time.Second,
	MaxReqPerDur:    10,
	DataUnmarshaler: cmFuturesWsPublicMsgUnmarshaler,
}

func DefaultCmPublicWsCfg() WsCfg {
	return cmPublicWsCfg
}

func DefaultPublicWsCfg(symbolType cex.SymbolType) WsCfg {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return DefaultSpotPublicWsCfg()
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return DefaultUmPublicWsCfg()
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return DefaultCmPublicWsCfg()
	}
	panic("bnc: unknown symbol type")
}
