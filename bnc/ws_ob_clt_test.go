package bnc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dwdwow/cex-go"
)

func TestOrderbookClient(t *testing.T) {
	clt := NewOrderBookClient(context.Background(), nil)
	clt.Sub(cex.SYMBOL_TYPE_SPOT, "BTCUSDT", "ETHUSDT", "SOLUSDT")
	clt.Sub(cex.SYMBOL_TYPE_UM_FUTURES, "BTCUSDT", "ETHUSDT", "SOLUSDT")
	ch1, err := clt.NewCh(cex.SYMBOL_TYPE_UM_FUTURES, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			o, ok := <-ch1
			if !ok {
				fmt.Println("ch1 closed")
				return
			}
			if o.Err != nil {
				fmt.Println(o.Err)
				continue
			}
			fmt.Println(o.Type, o.Symbol, o.Asks[0])
		}
	}()
	ch2, err := clt.NewCh(cex.SYMBOL_TYPE_UM_FUTURES, "ETHUSDT")
	if err != nil {
		t.Fatal(err)
	}
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
	ch3, err := clt.NewCh(cex.SYMBOL_TYPE_UM_FUTURES, "SOLUSDT")
	if err != nil {
		t.Fatal(err)
	}
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
	ch4, err := clt.NewCh(cex.SYMBOL_TYPE_SPOT, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			o, ok := <-ch4
			if !ok {
				fmt.Println("ch4 closed")
				return
			}
			if o.Err != nil {
				fmt.Println(o.Err)
				continue
			}
			fmt.Println(o.Type, o.Symbol, o.Asks[0])
		}
	}()
	ch5, err := clt.NewCh(cex.SYMBOL_TYPE_SPOT, "ETHUSDT")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			o, ok := <-ch5
			if !ok {
				fmt.Println("ch5 closed")
				return
			}
			if o.Err != nil {
				fmt.Println(o.Err)
				continue
			}
			fmt.Println(o.Type, o.Symbol, o.Asks[0])
		}
	}()
	ch6, err := clt.NewCh(cex.SYMBOL_TYPE_SPOT, "SOLUSDT")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			o, ok := <-ch6
			if !ok {
				fmt.Println("ch6 closed")
				return
			}
			if o.Err != nil {
				fmt.Println(o.Err)
				continue
			}
			fmt.Println(o.Type, o.Symbol, o.Asks[0])
		}
	}()
	time.Sleep(time.Second * 10)
	clt.Unsub(cex.SYMBOL_TYPE_UM_FUTURES, "BTCUSDT")
	time.Sleep(time.Second * 10)
	clt.RemoveCh(ch2)
	time.Sleep(time.Second * 10)
	clt.Unsub(cex.SYMBOL_TYPE_SPOT, "ETHUSDT", "SOLUSDT")
	time.Sleep(time.Second * 10)
	clt.RemoveCh(ch4)
	select {}
}

func TestOrderbookClient2(t *testing.T) {
	clt := NewOrderBookClient(context.Background(), nil)
	symbols := []string{
		"BTCUSDT", "ETHUSDT", "BNBUSDT", "SOLUSDT", "XRPUSDT",
		"ADAUSDT", "DOGEUSDT", "MATICUSDT", "DOTUSDT", "LTCUSDT",
		"AVAXUSDT", "TRXUSDT", "LINKUSDT", "ATOMUSDT", "UNIUSDT",
		"ETCUSDT", "NEARUSDT", "ALGOUSDT", "FTMUSDT", "AAVEUSDT",
		"EOSUSDT", "FILUSDT", "VETUSDT", "XLMUSDT", "ICPUSDT",
		"SANDUSDT", "MANAUSDT", "AXSUSDT", "RUNEUSDT", "APTUSDT",
		"LDOUSDT", "OPUSDT", "GMTUSDT", "GALAUSDT", "CHZUSDT",
		"THETAUSDT", "SNXUSDT", "APEUSDT", "GRTUSDT", "HBARUSDT",
		"XTZUSDT", "ZILUSDT", "FLOWUSDT", "QNTUSDT", "ENJUSDT",
		"MKRUSDT", "LRCUSDT", "ONEUSDT", "COMPUSDT", "DASHUSDT",
		"WAVESUSDT", "ZECUSDT", "KSMUSDT", "BATUSDT", "CRVUSDT",
		"SUSHIUSDT", "YFIUSDT", "1INCHUSDT", "CELOUSDT", "ZENUSDT",
		"ANKRUSDT", "BCHUSDT", "RENUSDT", "SRMUSDT", "COTIUSDT",
		"IOSTUSDT", "ONTUSDT", "OMGUSDT", "ICXUSDT", "DGBUSDT",
		"ZRXUSDT", "RVNUSDT", "SKLUSDT", "NEOUSDT", "AUDIOUSDT",
		"BANDUSDT", "STORJUSDT", "KNCUSDT", "RSRUSDT", "OCEANUSDT",
		"CVCUSDT", "BELUSDT", "IOTXUSDT", "SCUSDT", "DENTUSDT",
		"MTLUSDT", "HNTUSDT", "CHRUSDT", "STXUSDT", "ARPAUSDT",
		"CELRUSDT", "REEFUSDT", "C98USDT", "XEMUSDT", "OGNUSDT",
		"KAVAUSDT", "BTSUSDT", "SFPUSDT", "CTKUSDT", "BAKEUSDT",
		"ALPHAUSDT", "GTCUSDT", "TORNUSDT", "MDXUSDT", "PERPUSDT",
		"ROSEUSDT", "FLMUSDT", "KEEPUSDT", "LINAUSDT", "RAYUSDT",
		"MASKUSDT", "XVGUSDT", "KLAYUSDT", "BTGUSDT", "TFUELUSDT",
		"RLCUSDT", "TRBUSDT", "POWRUSDT", "LITUSDT", "DODOUSDT",
		"FORTHUSDT", "QUICKUSDT", "UFTUSDT", "PUNDIXUSDT", "WINGUSDT",
		"HARDUSDT", "WNXMUSDT", "SCRTUSDT", "VTHOUSDT", "STRAXUSDT",
		"FORUSDT", "UNFIUSDT", "FRONTUSDT", "FIROUSDT", "PERLUSDT",
		"RAMPUSDT", "SUPERUSDT", "CFXUSDT", "TLMUSDT", "BARUSDT",
		"SYSUSDT", "ACMUSDT", "IRISUSDT", "BTCSTUSDT", "TRUUSDT",
		"CKBUSDT", "TWTUSDT", "LITUSDT", "DARUSDT", "BIFIUSDT",
		"DREPUSDT", "PNTUSDT", "DIAUSDT", "GBPUSDT", "EURUSDT",
		"XVGUSDT", "STMXUSDT", "LSKUSDT", "BNTUSDT", "LTOUSDT",
		"MBLUSDT", "NKNUSDT", "ZENUSDT", "VIDTUSDT", "WRXUSDT",
		"CTSIUSDT", "HIVEUSDT", "CHRUSDT", "MDTUSDT", "STPTUSDT",
		"REPUSDT", "IOTAUSDT", "DATAUSDT", "CTXCUSDT", "BCHABCUSDT",
		"DUSKUSDT", "ARDRUSDT", "LENDUSDT", "MITHUSDT", "ATOMUSDT",
		"FETUSDT", "CELRUSDT", "XZCUSDT", "RENBTCUSDT", "NPXSUSDT",
		"TCTUSDT", "PHBUSDT", "TOMOUSDT", "KEYUSDT", "DOCKUSDT",
		"FUNUSDT", "CELOUSDT", "HOTUSDT", "MANAUSDT", "BNBUSDT",
		"XRPUSDT", "ETHUSDT", "BTCUSDT", "LINKUSDT", "ADAUSDT",
		"DOGEUSDT", "SOLUSDT", "DOTUSDT", "UNIUSDT", "LTCUSDT",
	}
	unsubed, err := clt.Sub(cex.SYMBOL_TYPE_SPOT, symbols...)
	if err != nil {
		t.Fatal(err, unsubed)
	}
	ch1, err := clt.NewCh(cex.SYMBOL_TYPE_SPOT, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch2, err := clt.NewCh(cex.SYMBOL_TYPE_SPOT, "C98USDT")
	if err != nil {
		t.Fatal(err)
	}
	ch3, err := clt.NewCh(cex.SYMBOL_TYPE_SPOT, "DOTUSDT")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			o, ok := <-ch1
			if !ok {
				fmt.Println("ch1 closed")
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
	time.Sleep(time.Second * 10)
	clt.Unsub(cex.SYMBOL_TYPE_SPOT, "BTCUSDT")
	time.Sleep(time.Second * 10)
	clt.RemoveCh(ch2)
	time.Sleep(time.Second * 10)
	clt.Unsub(cex.SYMBOL_TYPE_SPOT, "C98USDT")
	time.Sleep(time.Second * 10)
	clt.RemoveCh(ch3)
	select {}
}

func TestOrderbookBaseClient(t *testing.T) {
	clt := newOrderBookBaseClient(context.Background(), cex.SYMBOL_TYPE_UM_FUTURES, nil)
	unsubed, err := clt.subSymbols(
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
	ch := clt.newCh("BTCUSDT")
	ch2 := clt.newCh("ETHUSDT")
	ch3 := clt.newCh("SOLUSDT")
	go func() {
		time.Sleep(time.Second * 10)
		clt.unsubSymbols("BTCUSDT")
		time.Sleep(time.Second * 10)
		clt.removeCh(ch2)
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
}
