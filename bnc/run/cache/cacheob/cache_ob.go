package main

import (
	"os"
	"path/filepath"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/cex-go/bnc"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dataDir := filepath.Join(home, "cache", "cex", "bnc", "ob")
	go bnc.CacheOrderbook(dataDir, cex.SYMBOL_TYPE_SPOT,
		"BTCUSDT", "ETHUSDT", "SOLUSDT", "PEPEUSDT", "DOGEUSDT",
	)
	go bnc.CacheOrderbook(dataDir, cex.SYMBOL_TYPE_UM_FUTURES,
		"BTCUSDT", "ETHUSDT", "SOLUSDT", "1000PEPEUSDT", "DOGEUSDT", "HYPEUSDT",
	)
	select {}
}
