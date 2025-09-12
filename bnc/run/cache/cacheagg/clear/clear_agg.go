package main

import (
	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/cex-go/bnc"
)

func main() {
	go bnc.ClearMongoCachedAggTrades(cex.SYMBOL_TYPE_UM_FUTURES, []string{
		"BTCUSDT",
		"ETHUSDT",
		"SOLUSDT",
		"1000PEPEUSDT",
		"HYPEUSDT",
	})
	go bnc.ClearMongoCachedAggTrades(cex.SYMBOL_TYPE_SPOT, []string{
		"BTCUSDT",
		"ETHUSDT",
		"SOLUSDT",
		"PEPEUSDT",
	})
	select {}
}
