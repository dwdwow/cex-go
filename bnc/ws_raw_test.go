package bnc

import (
	"context"
	"testing"

	"github.com/dwdwow/props"
)

func TestSpotPublicWsClient(t *testing.T) {
	ws := NewSlightWsClient(context.Background(), spotPublicWsCfg, nil)
	_, err := ws.SubDepthUpdateStream("btcusdt", "ethusdt")
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
	ws := NewSlightWsClient(context.Background(), spotPublicWsCfg, nil)
	_, err := ws.SubAggTradeStream("BTCUSDT")
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

// func TestSpotPublicWsClientKline(t *testing.T) {
// 	pairs, _, err := QuerySpotPairs()
// 	props.PanicIfNotNil(err)
// 	symbols := []string{}
// 	for _, p := range pairs {
// 		if p.Tradable && p.Quote == "USDT" {
// 			symbols = append(symbols, p.PairSymbol)
// 		}
// 	}
// 	fmt.Println("symbols", symbols)
// 	ws := NewWsClient(SpotPublicWsCfg, nil, nil)
// 	err = ws.Start()
// 	props.PanicIfNotNil(err)
// 	err = ws.SubKlineStream(KlineInterval1s, symbols[:10]...)
// 	props.PanicIfNotNil(err)
// 	sub, err := ws.SubKline(symbols[5], KlineInterval1s)
// 	props.PanicIfNotNil(err)
// 	for {
// 		msg := <-sub.Chan()
// 		if msg.Err != nil {
// 			t.Error(msg.Err)
// 			break
// 		}
// 		t.Logf("%+v", msg.Data)
// 	}
// }

func TestUmFuturesPublicWsClient(t *testing.T) {
	ws := NewSlightWsClient(context.Background(), umPublicWsCfg, nil)
	_, err := ws.SubDepthUpdateStream("btcusdt", "ethusdt")
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

func TestCmFuturesPublicWsClient(t *testing.T) {
	ws := NewSlightWsClient(context.Background(), cmPublicWsCfg, nil)
	_, err := ws.SubKlineStream("btcusd_perp", "ethusd_perp")
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
	ws := NewSlightWsClient(context.Background(), spotPublicWsCfg, nil)
	_, err := ws.SubTradeStream("ETHUSDT", "BTCUSDT")
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
	ws := NewSlightWsClient(context.Background(), cmPublicWsCfg, nil)
	_, err := ws.SubAggTradeStream("BTCUSD_PERP", "ETHUSD_PERP")
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
	ws := NewSlightWsClient(context.Background(), spotPublicWsCfg, nil)
	_, err := ws.SubKlineStream(KLINE_INTERVAL_1s, "ETHUSDT", "BTCUSDT")
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
	ws := NewSlightWsClient(context.Background(), spotPublicWsCfg, nil)
	_, err := ws.SubDepthUpdateStream100ms("ETHUSDT", "BTCUSDT")
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
	ws := NewSlightWsClient(context.Background(), cmPublicWsCfg, nil)
	_, err := ws.SubMarkPriceStream3s("ETHUSD_PERP", "BTCUSD_PERP")
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
	ws := NewSlightWsClient(context.Background(), cmPublicWsCfg, nil)
	_, err := ws.SubAllMarkPriceStream1s()
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
	ws := NewSlightWsClient(context.Background(), cmPublicWsCfg, nil)
	_, err := ws.SubCMIndexPriceStream3s("ETHUSD", "BTCUSD")
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
	ws := NewSlightWsClient(context.Background(), cmPublicWsCfg, nil)
	_, err := ws.SubLiquidationOrderStream("ETHUSD_PERP", "BTCUSD_PERP")
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
	ws := NewSlightWsClient(context.Background(), cmPublicWsCfg, nil)
	_, err := ws.SubAllMarketLiquidationOrderStream()
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
