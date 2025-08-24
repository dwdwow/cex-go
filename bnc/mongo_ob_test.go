package bnc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dwdwow/cex-go"
)

func TestMongoObClient(t *testing.T) {
	cfg := MongoObClientCfg{
		SymbolType: cex.SYMBOL_TYPE_UM_FUTURES,
		Symbol:     "BTCUSDT",
		DbName:     "bnc_realtime_cache",
		StartTime:  time.Now().Add(-time.Hour),
	}
	clt := NewMongoObClient(context.Background(), cfg, nil)
	err := clt.Run()
	if err != nil {
		t.Fatal(err)
	}
	for {
		obData, err := clt.Read()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(time.UnixMilli(obData.Note.TxTime), obData.Symbol, obData.Asks[0], obData.Bids[0])
	}
}
