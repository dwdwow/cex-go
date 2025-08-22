package bnc

import (
	"testing"
	"time"

	"github.com/dwdwow/cex-go"
)

func TestBookTickerSingleWs(t *testing.T) {
	ws, err := StartNewBookTickerSingleWs(cex.SYMBOL_TYPE_UM_FUTURES, nil)
	if err != nil {
		t.Fatal(err)
	}
	unsubed, err := ws.Sub(spotSymbols200[:20]...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	ch1 := ws.NewCh("BTCUSDT")
	ch2 := ws.NewCh("ETHUSDT")
	ch3 := ws.NewCh("BNBUSDT")
	ch4 := ws.NewCh("ADAUSDT")
	go func() {
		for msg := range ch1 {
			if msg.Err != nil {
				t.Log(msg.Err)
				continue
			}
			t.Log(msg)
		}
		t.Log("ch1 done")
	}()
	go func() {
		for msg := range ch2 {
			if msg.Err != nil {
				t.Log(msg.Err)
				continue
			}
			t.Log(msg)
		}
		t.Log("ch2 done")
	}()
	go func() {
		for msg := range ch3 {
			if msg.Err != nil {
				t.Log(msg.Err)
				continue
			}
			t.Log(msg)
		}
		t.Log("ch3 done")
	}()
	go func() {
		for msg := range ch4 {
			if msg.Err != nil {
				t.Log(msg.Err)
				continue
			}
			t.Log(msg)
		}
		t.Log("ch4 done")
	}()
	time.Sleep(time.Second * 5)
	err = ws.Unsub("BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 5)
	ws.RemoveCh(ch2)
	time.Sleep(time.Second * 5)
}

func TestBookTickerWs(t *testing.T) {
	ws := NewBookTickerWs(nil)
	unsubed, err := ws.Sub(cex.SYMBOL_TYPE_SPOT, spotSymbols200...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	unsubed, err = ws.Sub(cex.SYMBOL_TYPE_UM_FUTURES, spotSymbols200...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	ch1, err := ws.NewCh(cex.SYMBOL_TYPE_SPOT, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch2, err := ws.NewCh(cex.SYMBOL_TYPE_UM_FUTURES, "ETHUSDT")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for msg := range ch1 {
			if msg.Err != nil {
				t.Log(msg.Err)
				continue
			}
			// t.Log(msg)
		}
		t.Log("ch1 done")
	}()
	go func() {
		for msg := range ch2 {
			if msg.Err != nil {
				t.Log(msg.Err)
				continue
			}
			// t.Log(msg)
		}
		t.Log("ch2 done")
	}()
	time.Sleep(time.Second * 5)
	err = ws.Unsub(cex.SYMBOL_TYPE_UM_FUTURES, "ETHUSDT")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 5)
	ws.RemoveCh(ch1)
	time.Sleep(time.Second * 5)
}

func TestBookTickerWs2(t *testing.T) {
	ws := NewBookTickerWs(nil)
	unsubed, err := ws.Sub(cex.SYMBOL_TYPE_SPOT, spotSymbols200...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	ch1, err := ws.NewCh(cex.SYMBOL_TYPE_SPOT, "ETHBTC")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for msg := range ch1 {
			if msg.Err != nil {
				t.Log(msg.Err)
				continue
			}
			t.Log(msg)
		}
		t.Log("ch1 done")
	}()
	time.Sleep(time.Second * 5)
	ws.RemoveCh(ch1)
	time.Sleep(time.Second * 5)
}
