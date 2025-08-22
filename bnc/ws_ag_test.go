package bnc

import (
	"fmt"
	"testing"
	"time"

	"github.com/dwdwow/cex-go"
)

func TestAggTradeSingleWs(t *testing.T) {
	clt, err := StartNewAggTradeSingleWs(cex.SYMBOL_TYPE_SPOT, nil)
	if err != nil {
		t.Fatal(err)
	}
	unsubed, err := clt.Sub(spotSymbols200[:20]...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	ch1 := clt.NewCh("BTCUSDT")
	ch2 := clt.NewCh("BTCUSDT")
	ch3 := clt.NewCh("ETHUSDT")
	ch4 := clt.NewCh("ETHUSDT")
	go func() {
		for msg := range ch1 {
			if msg.Err != nil {
				fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Err)
				continue
			}
			fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Stream.Price, msg.Stream.Qty)
		}
		fmt.Println("ch1 done", cex.SYMBOL_TYPE_SPOT, "BTCUSDT")
	}()
	go func() {
		for msg := range ch2 {
			if msg.Err != nil {
				fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Err)
				continue
			}
			fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Stream.Price, msg.Stream.Qty)
		}
		fmt.Println("ch2 done", cex.SYMBOL_TYPE_UM_FUTURES, "BTCUSDT")
	}()
	go func() {
		for msg := range ch3 {
			if msg.Err != nil {
				fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Err)
				continue
			}
			fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Stream.Price, msg.Stream.Qty)
		}
		fmt.Println("ch3 done", cex.SYMBOL_TYPE_SPOT, "ETHUSDT")
	}()
	go func() {
		for msg := range ch4 {
			if msg.Err != nil {
				fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Err)
				continue
			}
			fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Stream.Price, msg.Stream.Qty)
		}
		fmt.Println("ch4 done", cex.SYMBOL_TYPE_UM_FUTURES, "ETHUSDT")
	}()
	time.Sleep(time.Second * 10)
	err = clt.Unsub("BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 10)
	clt.RemoveCh(ch2)
	time.Sleep(time.Second * 10)
	err = clt.Unsub("ETHUSDT")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 10)
	clt.RemoveCh(ch3)
	time.Sleep(time.Second * 10)
}

func TestAggTradeWs(t *testing.T) {
	clt := NewAggTradeWs(nil)
	unsubed, err := clt.Sub(cex.SYMBOL_TYPE_SPOT, spotSymbols200...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	unsubed, err = clt.Sub(cex.SYMBOL_TYPE_UM_FUTURES, spotSymbols200...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	ch1, err := clt.NewCh(cex.SYMBOL_TYPE_SPOT, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch2, err := clt.NewCh(cex.SYMBOL_TYPE_SPOT, "ETHUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch3, err := clt.NewCh(cex.SYMBOL_TYPE_UM_FUTURES, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch4, err := clt.NewCh(cex.SYMBOL_TYPE_UM_FUTURES, "ETHUSDT")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for msg := range ch1 {
			if msg.Err != nil {
				fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Err)
				continue
			}
			fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Stream.Price, msg.Stream.Qty)
		}
		fmt.Println("ch1 done", cex.SYMBOL_TYPE_SPOT, "BTCUSDT")
	}()
	go func() {
		for msg := range ch2 {
			if msg.Err != nil {
				fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Err)
				continue
			}
			fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Stream.Price, msg.Stream.Qty)
		}
		fmt.Println("ch2 done", cex.SYMBOL_TYPE_SPOT, "ETHUSDT")
	}()
	go func() {
		for msg := range ch3 {
			if msg.Err != nil {
				fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Err)
				continue
			}
			fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Stream.Price, msg.Stream.Qty)
		}
		fmt.Println("ch3 done", cex.SYMBOL_TYPE_UM_FUTURES, "BTCUSDT")
	}()
	go func() {
		for msg := range ch4 {
			if msg.Err != nil {
				fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Err)
				continue
			}
			fmt.Println(msg.SymbolType, msg.Stream.Symbol, msg.Stream.Price, msg.Stream.Qty)
		}
		fmt.Println("ch4 done", cex.SYMBOL_TYPE_UM_FUTURES, "ETHUSDT")
	}()
	time.Sleep(time.Second * 10)
	err = clt.Unsub(cex.SYMBOL_TYPE_SPOT, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 10)
	clt.RemoveCh(ch2)
	time.Sleep(time.Second * 10)
	err = clt.Unsub(cex.SYMBOL_TYPE_UM_FUTURES, "ETHUSDT")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 10)
	clt.RemoveCh(ch3)
	time.Sleep(time.Second * 10)
}
