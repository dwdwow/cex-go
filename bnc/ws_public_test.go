package bnc

import (
	"fmt"
	"testing"
	"time"

	"github.com/dwdwow/cex-go"
)

func TestPublicWs(t *testing.T) {
	clt := NewPublicWs(nil)
	unsubed, err := clt.SubKlines(cex.SYMBOL_TYPE_SPOT, KLINE_INTERVAL_1m, spotSymbols200...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	unsubed, err = clt.SubKlines(cex.SYMBOL_TYPE_UM_FUTURES, KLINE_INTERVAL_1m, spotSymbols200...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	unsubed, err = clt.SubAggTrades(cex.SYMBOL_TYPE_SPOT, spotSymbols200...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	unsubed, err = clt.SubAggTrades(cex.SYMBOL_TYPE_UM_FUTURES, spotSymbols200...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	unsubed, err = clt.SubOrderBooks(cex.SYMBOL_TYPE_SPOT, spotSymbols200...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	unsubed, err = clt.SubOrderBooks(cex.SYMBOL_TYPE_UM_FUTURES, spotSymbols200...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	unsubed, err = clt.SubDepthUpdates(cex.SYMBOL_TYPE_SPOT, spotSymbols200...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	unsubed, err = clt.SubDepthUpdates(cex.SYMBOL_TYPE_UM_FUTURES, spotSymbols200...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	unsubed, err = clt.SubBookTickers(cex.SYMBOL_TYPE_SPOT, spotSymbols200...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	unsubed, err = clt.SubBookTickers(cex.SYMBOL_TYPE_UM_FUTURES, spotSymbols200...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	ch1, err := clt.NewKlineCh(cex.SYMBOL_TYPE_SPOT, KLINE_INTERVAL_1m, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch2, err := clt.NewKlineCh(cex.SYMBOL_TYPE_UM_FUTURES, KLINE_INTERVAL_1m, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch3, err := clt.NewAggTradeCh(cex.SYMBOL_TYPE_SPOT, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch4, err := clt.NewAggTradeCh(cex.SYMBOL_TYPE_UM_FUTURES, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch5, err := clt.NewOrderBookCh(cex.SYMBOL_TYPE_SPOT, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch6, err := clt.NewOrderBookCh(cex.SYMBOL_TYPE_UM_FUTURES, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch7, err := clt.NewDepthUpdateCh(cex.SYMBOL_TYPE_SPOT, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch8, err := clt.NewDepthUpdateCh(cex.SYMBOL_TYPE_UM_FUTURES, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch9, err := clt.NewBookTickerCh(cex.SYMBOL_TYPE_SPOT, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch10, err := clt.NewBookTickerCh(cex.SYMBOL_TYPE_UM_FUTURES, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for msg := range ch1 {
			if msg.Err != nil {
				panic(msg.Err)
			}
			fmt.Println(msg.SymbolType, msg.Kline.Symbol, msg.Kline.Interval, msg.Kline.ClosePrice)
		}
	}()
	go func() {
		for msg := range ch2 {
			if msg.Err != nil {
				panic(msg.Err)
			}
			fmt.Println(msg.SymbolType, msg.Kline.Symbol, msg.Kline.Interval, msg.Kline.ClosePrice)
		}
	}()
	go func() {
		for msg := range ch3 {
			if msg.Err != nil {
				panic(msg.Err)
			}
			fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Stream.Price, msg.Stream.Qty)
		}
	}()
	go func() {
		for msg := range ch4 {
			if msg.Err != nil {
				panic(msg.Err)
			}
			fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Stream.Price, msg.Stream.Qty)
		}
	}()
	go func() {
		for msg := range ch5 {
			if msg.Err != nil {
				panic(msg.Err)
			}
			fmt.Println(msg.Type, msg.Symbol, msg.Asks[0], msg.Asks[0])
		}
	}()
	go func() {
		for msg := range ch6 {
			if msg.Err != nil {
				panic(msg.Err)
			}
			fmt.Println(msg.Type, msg.Symbol, msg.Asks[0], msg.Asks[0])
		}
	}()
	go func() {
		for msg := range ch7 {
			if msg.Err != nil {
				panic(msg.Err)
			}
			if len(msg.Stream.Asks) > 0 && len(msg.Stream.Bids) > 0 {
				fmt.Println(msg.Stream.Symbol, msg.Stream.Asks[0], msg.Stream.Bids[0])
			}
		}
	}()
	go func() {
		for msg := range ch8 {
			if msg.Err != nil {
				panic(msg.Err)
			}
			if len(msg.Stream.Asks) > 0 && len(msg.Stream.Bids) > 0 {
				fmt.Println(msg.Stream.Symbol, msg.Stream.Asks[0], msg.Stream.Bids[0])
			}
		}
	}()
	go func() {
		for msg := range ch9 {
			if msg.Err != nil {
				panic(msg.Err)
			}
			fmt.Println(msg.Stream.Symbol, msg.Stream.BestBidPrice, msg.Stream.BestAskPrice)
		}
	}()
	go func() {
		for msg := range ch10 {
			if msg.Err != nil {
				panic(msg.Err)
			}
			fmt.Println(msg.Stream.Symbol, msg.Stream.BestBidPrice, msg.Stream.BestAskPrice)
		}
	}()
	time.Sleep(time.Second * 10)
}
