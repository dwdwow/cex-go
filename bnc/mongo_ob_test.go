package bnc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dwdwow/cex-go"
)

func TestMongoObClient(t *testing.T) {
	mongoUri := "mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000"
	cfg := MongoObClientCfg{
		SymbolType: cex.SYMBOL_TYPE_UM_FUTURES,
		Symbol:     "BTCUSDT",
		MongoUri:   mongoUri,
		DbName:     "bnc_realtime_cache",
		StartTime:  time.Now().Add(-time.Hour),
	}
	client := NewMongoObClient(context.Background(), cfg, nil)
	err := client.Run()
	if err != nil {
		t.Fatal(err)
	}
	for {
		obData, err := client.Read()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(obData.Symbol, obData.Asks[0], obData.Bids[0])
	}
}
