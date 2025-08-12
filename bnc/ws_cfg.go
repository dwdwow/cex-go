package bnc

import "time"

var SpotPublicWsCfg = WsCfg{
	Url:             WsBaseUrl,
	MaxStream:       1024,
	ReqDur:          time.Second,
	MaxReqPerDur:    5,
	DataUnmarshaler: SpotWsPublicMsgUnmarshaler,
}

var SpotPrivateWsCfg = WsCfg{
	Url:             WsBaseUrl,
	ListenKeyUrl:    API_ENDPOINT + API_V3 + "/userDataStream",
	MaxStream:       1024,
	ReqDur:          time.Second,
	MaxReqPerDur:    10,
	DataUnmarshaler: SpotWsPrivateMsgUnmarshaler,
}

var UmFuturesWsCfg = WsCfg{
	Url:             FutureWsBaseUrl,
	MaxStream:       200,
	ReqDur:          time.Second,
	MaxReqPerDur:    10,
	DataUnmarshaler: UmFuturesWsPublicMsgUnmarshaler,
}

var CmFuturesWsCfg = WsCfg{
	Url:             CMFutureWsBaseUrl,
	MaxStream:       200,
	ReqDur:          time.Second,
	MaxReqPerDur:    10,
	DataUnmarshaler: CmFuturesWsPublicMsgUnmarshaler,
}
