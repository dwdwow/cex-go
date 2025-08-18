package bnc

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestUmObWs(t *testing.T) {
	ws := NewUmObWs(context.TODO(), nil)
	unsubed, err := ws.SubNewSymbols(
		"BTCUSDT",
		"ETHUSDT",
		"SOLUSDT",
		"DOGEUSDT",
		"XRPUSDT",
		"TRXUSDT",
		// "LTCUSDT",
		// "ADAUSDT",
		// "BNBUSDT",
		// "DOTUSDT",
		// "LINKUSDT",
		// "MATICUSDT",
		// "AVAXUSDT",
		// "SHIBUSDT",
		// "ATOMUSDT",
		// "UNIUSDT",
		// "ETCUSDT",
		// "NEARUSDT",
		// "ALGOUSDT",
		// "FTMUSDT",
		// "SANDUSDT",
		// "MANAUSDT",
		// "GALAUSDT",
		// "APEUSDT",
		// "AXSUSDT",
		// "RUNEUSDT",
		// "GMTUSDT",
		// "EGLDUSDT",
		// "THETAUSDT",
		// "HBARUSDT",
		//
		// "XTZUSDT",
		// "ZILUSDT",
		// "VETUSDT",
		// "ICPUSDT",
		// "AAVEUSDT",
		// "EOSUSDT",
		// "MKRUSDT",
		// "CAKEUSDT",
		// "XMRUSDT",
		// "FLOWUSDT",
	)
	if err != nil {
		t.Fatal(err, unsubed)
	}
	ch := ws.NewCh("BTCUSDT")
	ch2 := ws.NewCh("ETHUSDT")
	ch3 := ws.NewCh("SOLUSDT")
	go func() {
		time.Sleep(time.Second * 30)
		ws.RemoveCh(ch)
	}()
	go func() {
		for {
			o, ok := <-ch
			if !ok {
				fmt.Println("ch closed")
				return
			}
			if o.Err != nil {
				fmt.Println(o.Err)
				continue
			}
			fmt.Println(o.Symbol, o.Asks[0])
		}
	}()
	go func() {
		for {
			o, ok := <-ch2
			if !ok {
				fmt.Println("ch2 closed")
				return
			}
			if o.Err != nil {
				fmt.Println(o.Err)
				continue
			}
			fmt.Println(o.Symbol, o.Asks[0])
		}
	}()
	go func() {
		for {
			o, ok := <-ch3
			if !ok {
				fmt.Println("ch3 closed")
				return
			}
			if o.Err != nil {
				fmt.Println(o.Err)
				continue
			}
			fmt.Println(o.Symbol, o.Asks[0])
		}
	}()
	select {}
}
