package bnc

import (
	"time"

	"github.com/dwdwow/cex-go"
)

const maxWsChCap = 10000

var spotPublicWsCfg = WsCfg{
	Url:             WsBaseUrl,
	MaxStream:       1024,
	ChCap:           maxWsChCap,
	ReqDur:          time.Second,
	MaxReqPerDur:    5,
	DataUnmarshaler: SpotWsPublicMsgUnmarshaler,
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
	DataUnmarshaler: SpotWsPrivateMsgUnmarshaler,
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
	DataUnmarshaler: UmFuturesWsPublicMsgUnmarshaler,
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
	DataUnmarshaler: CmFuturesWsPublicMsgUnmarshaler,
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
