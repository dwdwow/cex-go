package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/cex-go/bnc"
)

func main() {

	mongoUri := "mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000"
	cfg := bnc.MongoObClientCfg{
		SymbolType: cex.SYMBOL_TYPE_UM_FUTURES,
		Symbol:     "BTCUSDT",
		MongoUri:   mongoUri,
		DbName:     "bnc_realtime_cache",
		StartTime:  time.Now().Add(-time.Hour),
	}
	client := bnc.NewMongoObClient(context.Background(), cfg, nil)
	err := client.Run()
	if err != nil {
		panic(err)
	}
	for {
		obData, err := client.Read()
		if err != nil {
			panic(err)
		}
		fmt.Println(obData.Symbol, obData.Asks[0], obData.Bids[0])
	}
}
