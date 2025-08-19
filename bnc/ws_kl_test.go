package bnc

import (
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/dwdwow/cex-go"
)

var spotSymbols200 = []string{
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
}

func TestKlineBaseWs(t *testing.T) {
	clt, err := startNewKlineBaseWs(cex.SYMBOL_TYPE_UM_FUTURES, slog.Default())
	if err != nil {
		t.Fatal(err)
	}
	unsubed, err := clt.sub(KLINE_INTERVAL_1m, spotSymbols200[:20]...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	unsubed, err = clt.sub(KLINE_INTERVAL_1m, spotSymbols200[10:20]...)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	ch1 := clt.newCh(KLINE_INTERVAL_1m, "BTCUSDT")
	ch2 := clt.newCh(KLINE_INTERVAL_1m, "ETHUSDT")
	ch3 := clt.newCh(KLINE_INTERVAL_1m, "ADAUSDT")
	go func() {
		for msg := range ch1 {
			if msg.Err != nil {
				fmt.Println(msg.Err)
				continue
			}
			fmt.Println(msg.SymbolType, msg.Kline.Symbol, msg.Kline.Interval, msg.Kline.ClosePrice)
		}
	}()
	go func() {
		for msg := range ch2 {
			if msg.Err != nil {
				fmt.Println(msg.Err)
				continue
			}
			fmt.Println(msg.SymbolType, msg.Kline.Symbol, msg.Kline.Interval, msg.Kline.ClosePrice)
		}
	}()
	go func() {
		for msg := range ch3 {
			if msg.Err != nil {
				fmt.Println(msg.Err)
				continue
			}
			fmt.Println(msg.SymbolType, msg.Kline.Symbol, msg.Kline.Interval, msg.Kline.ClosePrice)
		}
	}()
	time.Sleep(time.Second * 10)
	err = clt.unsub(KLINE_INTERVAL_1m, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 10)
	clt.removeCh(ch2)
}

func TestKlineWs(t *testing.T) {
	clt := NewKlineWs(slog.Default())
	unsubed, err := clt.Sub(
		cex.SYMBOL_TYPE_SPOT,
		KLINE_INTERVAL_1m,
		spotSymbols200...,
	)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	unsubed, err = clt.Sub(
		cex.SYMBOL_TYPE_UM_FUTURES,
		KLINE_INTERVAL_1m,
		spotSymbols200...,
	)
	if err != nil {
		t.Fatal(err, len(unsubed), unsubed)
	}
	ch1, err := clt.NewCh(cex.SYMBOL_TYPE_SPOT, KLINE_INTERVAL_1m, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch2, err := clt.NewCh(cex.SYMBOL_TYPE_UM_FUTURES, KLINE_INTERVAL_1m, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch3, err := clt.NewCh(cex.SYMBOL_TYPE_SPOT, KLINE_INTERVAL_1m, "ETHUSDT")
	if err != nil {
		t.Fatal(err)
	}
	ch4, err := clt.NewCh(cex.SYMBOL_TYPE_UM_FUTURES, KLINE_INTERVAL_1m, "ETHUSDT")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for msg := range ch1 {
			if msg.Err != nil {
				fmt.Println(msg.SymbolType, msg.Kline.Symbol, msg.Kline.Interval, msg.Err)
				continue
			}
			fmt.Println(msg.SymbolType, msg.Kline.Symbol, msg.Kline.Interval, msg.Kline.ClosePrice)
		}
		fmt.Println("ch1 done", cex.SYMBOL_TYPE_SPOT, "BTCUSDT", KLINE_INTERVAL_1m)
	}()
	go func() {
		for msg := range ch2 {
			if msg.Err != nil {
				fmt.Println(msg.SymbolType, msg.Kline.Symbol, msg.Kline.Interval, msg.Err)
				continue
			}
			fmt.Println(msg.SymbolType, msg.Kline.Symbol, msg.Kline.Interval, msg.Kline.ClosePrice)
		}
		fmt.Println("ch2 done", cex.SYMBOL_TYPE_UM_FUTURES, "BTCUSDT", KLINE_INTERVAL_1m)
	}()
	go func() {
		for msg := range ch3 {
			if msg.Err != nil {
				fmt.Println(msg.SymbolType, msg.Kline.Symbol, msg.Kline.Interval, msg.Err)
				continue
			}
			fmt.Println(msg.SymbolType, msg.Kline.Symbol, msg.Kline.Interval, msg.Kline.ClosePrice)
		}
		fmt.Println("ch3 done", cex.SYMBOL_TYPE_SPOT, "ETHUSDT", KLINE_INTERVAL_1m)
	}()
	go func() {
		for msg := range ch4 {
			if msg.Err != nil {
				fmt.Println(msg.SymbolType, msg.Kline.Symbol, msg.Kline.Interval, msg.Err)
				continue
			}
			fmt.Println(msg.SymbolType, msg.Kline.Symbol, msg.Kline.Interval, msg.Kline.ClosePrice)
		}
		fmt.Println("ch4 done", cex.SYMBOL_TYPE_UM_FUTURES, "ETHUSDT", KLINE_INTERVAL_1m)
	}()
	time.Sleep(time.Second * 10)
	err = clt.Unsub(cex.SYMBOL_TYPE_SPOT, KLINE_INTERVAL_1m, "BTCUSDT")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 10)
	clt.RemoveCh(ch2)
	time.Sleep(time.Second * 10)
	err = clt.Unsub(cex.SYMBOL_TYPE_UM_FUTURES, KLINE_INTERVAL_1m, "ETHUSDT")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 10)
	clt.RemoveCh(ch3)
}
