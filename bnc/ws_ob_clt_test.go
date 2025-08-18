package bnc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dwdwow/cex-go"
)

func TestOrderbookClient(t *testing.T) {
	clt := newOrderBookBaseClient(context.Background(), cex.SYMBOL_TYPE_UM_FUTURES, nil)
	unsubed, err := clt.SubSymbols(
		"BTCUSDT",
		"ETHUSDT",
		"SOLUSDT",
		"DOGEUSDT",
		"XRPUSDT",
		"TRXUSDT",
		"LTCUSDT",
		"ADAUSDT",
		"BNBUSDT",
		"DOTUSDT",
		"LINKUSDT",
		"MATICUSDT",
		"AVAXUSDT",
		"SHIBUSDT",
		"ATOMUSDT",
		"UNIUSDT",
		"ETCUSDT",
		"NEARUSDT",
		"ALGOUSDT",
		"FTMUSDT",
		"AAVEUSDT",
		"EOSUSDT",
		"FILUSDT",
		"VETUSDT",
		"XLMUSDT",
		"ICPUSDT",
		"SANDUSDT",
		"MANAUSDT",
		"AXSUSDT",
		"RUNEUSDT",
		"APTUSDT",
		"LDOUSDT",
		"OPUSDT",
		"GMTUSDT",
		"GALAUSDT",
		"CHZUSDT",
		"THETAUSDT",
		"SNXUSDT",
		"APEUSDT",
		"GRTUSDT",
		"HBARUSDT",
		"XTZUSDT",
		"ZILUSDT",
		"FLOWUSDT",
		"QNTUSDT",
		"ENJUSDT",
		"MKRUSDT",
		"LRCUSDT",
		"ONEUSDT",
		"COMPUSDT",
	)

	if err != nil {
		t.Fatal(err, unsubed)
	}
	ch := clt.NewCh("BTCUSDT")
	ch2 := clt.NewCh("ETHUSDT")
	ch3 := clt.NewCh("SOLUSDT")
	go func() {
		time.Sleep(time.Second * 10)
		clt.UnsubSymbols("BTCUSDT")
		time.Sleep(time.Second * 10)
		clt.RemoveCh(ch2)
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
			fmt.Println(o.Type, o.Symbol, o.Asks[0])
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
			fmt.Println(o.Type, o.Symbol, o.Asks[0])
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
			fmt.Println(o.Type, o.Symbol, o.Asks[0])
		}
	}()
	select {}
}
