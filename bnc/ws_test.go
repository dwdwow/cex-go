package bnc

import (
	"testing"

	"github.com/dwdwow/props"
)

func TestSpotPublicWsClient(t *testing.T) {
	ws := NewWsClient(SpotPublicWsCfg, nil, nil)
	err := ws.Start()
	props.PanicIfNotNil(err)
	err = ws.SubStream("btcusdt@depth", "ethusdt@depth@100ms")
	props.PanicIfNotNil(err)
	chAll, err := ws.Sub(string(WsEventDepthUpdate))
	props.PanicIfNotNil(err)
	for {
		msg := <-chAll
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}

func TestSpotPublicWsClientAggTrade(t *testing.T) {
	ws := NewWsClient(SpotPublicWsCfg, nil, nil)
	err := ws.Start()
	props.PanicIfNotNil(err)
	err = ws.SubAggTradeStream("BTCUSDT")
	props.PanicIfNotNil(err)
	sub, err := ws.SubAggTrade("BTCUSDT")
	props.PanicIfNotNil(err)
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
	ws := NewWsClient(UmFuturesWsCfg, nil, nil)
	err := ws.start()
	props.PanicIfNotNil(err)
	err = ws.SubStream("btcusdt@depth", "ethusdt@depth@100ms")
	props.PanicIfNotNil(err)
	chAll, err := ws.Sub(string(WsEventDepthUpdate))
	props.PanicIfNotNil(err)
	for {
		msg := <-chAll
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}

func TestCmFuturesPublicWsClient(t *testing.T) {
	ws := NewWsClient(CmFuturesWsCfg, nil, nil)
	err := ws.start()
	props.PanicIfNotNil(err)
	err = ws.SubStream("btcusd_perp@kline_1m", "ethusd_perp@kline_1m")
	props.PanicIfNotNil(err)
	chAll, err := ws.Sub(string(WsEventKline))
	props.PanicIfNotNil(err)
	for {
		msg := <-chAll
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}

func TestWsClient_SubTrade(t *testing.T) {
	ws := NewWsClient(SpotPublicWsCfg, nil, nil)
	err := ws.Start()
	props.PanicIfNotNil(err)
	err = ws.SubTradeStream("ETHUSDT", "BTCUSDT")
	props.PanicIfNotNil(err)
	sub, err := ws.SubTrade("ETHUSDT")
	props.PanicIfNotNil(err)
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
	ws := NewWsClient(CmFuturesWsCfg, nil, nil)
	err := ws.start()
	props.PanicIfNotNil(err)
	err = ws.SubAggTradeStream("BTCUSD_PERP", "ETHUSD_PERP")
	props.PanicIfNotNil(err)
	sub, err := ws.SubAggTrade("ETHUSD_PERP")
	props.PanicIfNotNil(err)
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
	ws := NewWsClient(SpotPublicWsCfg, nil, nil)
	err := ws.start()
	props.PanicIfNotNil(err)
	err = ws.SubKlineStream(KLINE_INTERVAL_1s, "ETHUSDT", "BTCUSDT")
	props.PanicIfNotNil(err)
	sub, err := ws.SubKline("ETHUSDT", KLINE_INTERVAL_1s)
	props.PanicIfNotNil(err)
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
	ws := NewWsClient(SpotPublicWsCfg, nil, nil)
	err := ws.start()
	props.PanicIfNotNil(err)
	err = ws.SubDepthUpdateStream100ms("ETHUSDT", "BTCUSDT")
	props.PanicIfNotNil(err)
	sub, err := ws.SubDepthUpdate100ms("BTCUSDT")
	props.PanicIfNotNil(err)
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
	ws := NewWsClient(CmFuturesWsCfg, nil, nil)
	err := ws.start()
	props.PanicIfNotNil(err)
	err = ws.SubMarkPriceStream3s("ETHUSD_PERP", "BTCUSD_PERP")
	props.PanicIfNotNil(err)
	sub, err := ws.SubMarkPrice3s("ETHUSD_PERP")
	props.PanicIfNotNil(err)
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
	ws := NewWsClient(CmFuturesWsCfg, nil, nil)
	err := ws.start()
	props.PanicIfNotNil(err)
	err = ws.SubAllMarkPriceStream1s()
	props.PanicIfNotNil(err)
	sub, err := ws.SubAllMarkPrice1s()
	props.PanicIfNotNil(err)
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
	ws := NewWsClient(CmFuturesWsCfg, nil, nil)
	err := ws.start()
	props.PanicIfNotNil(err)
	err = ws.SubCMIndexPriceStream3s("ETHUSD", "BTCUSD")
	props.PanicIfNotNil(err)
	sub, err := ws.SubCMIndexPrice3s("ETHUSD")
	props.PanicIfNotNil(err)
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
	ws := NewWsClient(CmFuturesWsCfg, nil, nil)
	err := ws.start()
	props.PanicIfNotNil(err)
	err = ws.SubLiquidationOrderStream("ETHUSD_PERP", "BTCUSD_PERP")
	props.PanicIfNotNil(err)
	sub, err := ws.SubLiquidationOrder("")
	props.PanicIfNotNil(err)
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
	ws := NewWsClient(CmFuturesWsCfg, nil, nil)
	err := ws.start()
	props.PanicIfNotNil(err)
	err = ws.SubAllMarketLiquidationOrderStream()
	props.PanicIfNotNil(err)
	sub, err := ws.SubAllMarketLiquidationOrder()
	props.PanicIfNotNil(err)
	for {
		msg := <-sub.Chan()
		if msg.Err != nil {
			t.Error(msg.Err)
			break
		}
		t.Logf("%+v", msg.Data)
	}
}
