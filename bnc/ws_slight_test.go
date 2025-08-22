package bnc

import (
	"fmt"
	"testing"
	"time"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/props"
)

func TestSpotPublicWsClient(t *testing.T) {
	ws, err := StartNewSlightWsClient(spotPublicWsCfg, spotWsPublicMsgUnmarshaler, nil)
	props.PanicIfNotNil(err)
	_, err = ws.SubDepthUpdateStream("btcusdt", "ethusdt")
	props.PanicIfNotNil(err)
	sub := ws.SubDepthUpdate("BTCUSDT")
	for {
		msg := <-sub.Chan()
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}

func TestSpotPublicWsClientAggTrade(t *testing.T) {
	ws, err := StartNewSlightWsClient(spotPublicWsCfg, spotWsPublicMsgUnmarshaler, nil)
	props.PanicIfNotNil(err)
	_, err = ws.SubAggTradeStream("BTCUSDT")
	props.PanicIfNotNil(err)
	sub := ws.SubAggTrade("BTCUSDT")
	for {
		msg := <-sub.Chan()
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}

func TestUmFuturesPublicWsClient(t *testing.T) {
	ws, err := StartNewSlightWsClient(umPublicWsCfg, umFuturesWsPublicMsgUnmarshaler, nil)
	props.PanicIfNotNil(err)
	_, err = ws.SubDepthUpdateStream("btcusdt", "ethusdt")
	props.PanicIfNotNil(err)
	ch := ws.SubDepthUpdate("BTCUSDT").Chan()
	ch3 := ws.SubDepthUpdate("ETHUSDT").Chan()
	ws2 := NewPublicSlightWsClient(cex.SYMBOL_TYPE_UM_FUTURES)
	ch2 := ws2.SubDepthUpdate("BTCUSDT").Chan()
	go func() {
		time.Sleep(time.Second*10 + time.Millisecond*77)
		err := ws2.Start()
		if err != nil {
			panic(err)
		}
		_, err = ws2.SubDepthUpdateStream("BTCUSDT")
		if err != nil {
			panic(err)
		}
	}()

	for {
		var msg SlightWsClientSubscriptionMsg[WsDepthStream]
		select {
		case msg = <-ch:
		case msg = <-ch2:
		case msg = <-ch3:
		}
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		data := msg.Data
		fmt.Println(data.Symbol, data.EventTime, data.FirstId, data.LastId)
	}
}

func TestCmFuturesPublicWsClient(t *testing.T) {
	ws, err := StartNewSlightWsClient(cmPublicWsCfg, cmFuturesWsPublicMsgUnmarshaler, nil)
	props.PanicIfNotNil(err)
	_, err = ws.SubKlineStream("btcusd_perp", "ethusd_perp")
	props.PanicIfNotNil(err)
	sub := ws.SubKline("BTCUSD_PERP", KLINE_INTERVAL_1m)
	for {
		msg := <-sub.Chan()
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}

func TestWsClient_SubTrade(t *testing.T) {
	ws, err := StartNewSlightWsClient(spotPublicWsCfg, spotWsPublicMsgUnmarshaler, nil)
	props.PanicIfNotNil(err)
	_, err = ws.SubTradeStream("ETHUSDT", "BTCUSDT")
	props.PanicIfNotNil(err)
	sub := ws.SubTrade("ETHUSDT")
	for {
		msg := <-sub.Chan()
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}

func TestWsClient_SubAggTrade(t *testing.T) {
	ws, err := StartNewSlightWsClient(cmPublicWsCfg, cmFuturesWsPublicMsgUnmarshaler, nil)
	props.PanicIfNotNil(err)
	_, err = ws.SubAggTradeStream("BTCUSD_PERP", "ETHUSD_PERP")
	props.PanicIfNotNil(err)
	sub := ws.SubAggTrade("ETHUSD_PERP")
	for {
		msg := <-sub.Chan()
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}

func TestWsClient_SubKline(t *testing.T) {
	ws, err := StartNewSlightWsClient(spotPublicWsCfg, spotWsPublicMsgUnmarshaler, nil)
	props.PanicIfNotNil(err)
	_, err = ws.SubKlineStream(KLINE_INTERVAL_1s, "ETHUSDT", "BTCUSDT")
	props.PanicIfNotNil(err)
	sub := ws.SubKline("ETHUSDT", KLINE_INTERVAL_1s)
	for {
		msg := <-sub.Chan()
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}

func TestWsClient_SubDepthUpdate(t *testing.T) {
	ws, err := StartNewSlightWsClient(spotPublicWsCfg, spotWsPublicMsgUnmarshaler, nil)
	props.PanicIfNotNil(err)
	_, err = ws.SubDepthUpdateStream100ms("ETHUSDT", "BTCUSDT")
	props.PanicIfNotNil(err)
	sub := ws.SubDepthUpdate100ms("BTCUSDT")
	for {
		msg := <-sub.Chan()
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}

func TestWsClient_SubMarkPrice1s(t *testing.T) {
	ws, err := StartNewSlightWsClient(cmPublicWsCfg, cmFuturesWsPublicMsgUnmarshaler, nil)
	props.PanicIfNotNil(err)
	_, err = ws.SubMarkPriceStream3s("ETHUSD_PERP", "BTCUSD_PERP")
	props.PanicIfNotNil(err)
	sub := ws.SubMarkPrice3s("ETHUSD_PERP")
	for {
		msg := <-sub.Chan()
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}

func TestWsClient_SubAllMarkPrice1s(t *testing.T) {
	ws, err := StartNewSlightWsClient(cmPublicWsCfg, cmFuturesWsPublicMsgUnmarshaler, nil)
	props.PanicIfNotNil(err)
	_, err = ws.SubAllMarkPriceStream1s()
	props.PanicIfNotNil(err)
	sub := ws.SubAllMarkPrice1s()
	for {
		msg := <-sub.Chan()
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}

func TestWsClient_SubCMIndexPrice1s(t *testing.T) {
	ws, err := StartNewSlightWsClient(cmPublicWsCfg, cmFuturesWsPublicMsgUnmarshaler, nil)
	props.PanicIfNotNil(err)
	_, err = ws.SubCMIndexPriceStream3s("ETHUSD", "BTCUSD")
	props.PanicIfNotNil(err)
	sub := ws.SubCMIndexPrice3s("ETHUSD")
	for {
		msg := <-sub.Chan()
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}

func TestWsClient_SubLiquidationOrder(t *testing.T) {
	ws, err := StartNewSlightWsClient(cmPublicWsCfg, cmFuturesWsPublicMsgUnmarshaler, nil)
	props.PanicIfNotNil(err)
	_, err = ws.SubLiquidationOrderStream("ETHUSD_PERP", "BTCUSD_PERP")
	props.PanicIfNotNil(err)
	sub := ws.SubLiquidationOrder("")
	for {
		msg := <-sub.Chan()
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}

func TestWsClient_SubAllMarketLiquidationOrder(t *testing.T) {
	ws, err := StartNewSlightWsClient(cmPublicWsCfg, cmFuturesWsPublicMsgUnmarshaler, nil)
	props.PanicIfNotNil(err)
	_, err = ws.SubAllMarketLiquidationOrderStream()
	props.PanicIfNotNil(err)
	sub := ws.SubAllMarketLiquidationOrder()
	for {
		msg := <-sub.Chan()
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}
