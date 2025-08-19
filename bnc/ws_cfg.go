package bnc

import (
	"github.com/dwdwow/cex-go"
)

type RawWsCfg struct {
	// ws url without streams and auth tokens
	Url string

	APIKey       string
	APISecretKey string
	ListenKeyUrl string

	// default is millisecond
	// if true, use microsecond
	// just spot stream has microsecond
	Microsecond bool

	// max stream per websocket client
	// ex. spot 1024, portfolio margin / futures 200
	// normally just for public websocket
	MaxStream int

	// max request per second
	MaxReqPerSecond int
}

var spotPublicWsCfg = RawWsCfg{
	Url:             WsBaseUrl,
	MaxStream:       1024,
	Microsecond:     true,
	MaxReqPerSecond: maxSpWsReqPerSec,
}

func DefaultSpotPublicWsCfg() RawWsCfg {
	return spotPublicWsCfg
}

var spotPrivateWsCfg = RawWsCfg{
	Url:             WsBaseUrl,
	ListenKeyUrl:    API_ENDPOINT + API_V3 + "/userDataStream",
	MaxStream:       1024,
	Microsecond:     true,
	MaxReqPerSecond: maxSpWsReqPerSec,
}

func DefaultSpotPrivateWsCfg() RawWsCfg {
	return spotPrivateWsCfg
}

var umPublicWsCfg = RawWsCfg{
	Url:             FutureWsBaseUrl,
	MaxStream:       200,
	MaxReqPerSecond: maxFuWsReqPerSec,
}

func DefaultUmPublicWsCfg() RawWsCfg {
	return umPublicWsCfg
}

var cmPublicWsCfg = RawWsCfg{
	Url:             CMFutureWsBaseUrl,
	MaxStream:       200,
	MaxReqPerSecond: maxFuWsReqPerSec,
}

func DefaultCmPublicWsCfg() RawWsCfg {
	return cmPublicWsCfg
}

func DefaultPublicWsCfg(symbolType cex.SymbolType) RawWsCfg {
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

type WebsocketCfg struct {
}

var defaultWebsocketCfg = WebsocketCfg{}

func DefaultWebsocketCfg() WebsocketCfg {
	return defaultWebsocketCfg
}
