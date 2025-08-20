package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/cex-go/bnc"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dataDir := filepath.Join(home, "cache", "cex", "bnc", "depth_update")
	pairs, err := bnc.GetSpotSymbols()
	if err != nil {
		panic(err)
	}
	var symbols []string
	for _, pair := range pairs {
		if !pair.Tradable {
			continue
		}
		symbols = append(symbols, pair.Symbol)
	}
	go bnc.CacheDepthUpdate(dataDir, cex.SYMBOL_TYPE_SPOT, symbols...)

	pairs, err = bnc.GetUMSymbols()
	if err != nil {
		panic(err)
	}
	for _, pair := range pairs {
		if !pair.Tradable {
			continue
		}
		if !pair.IsPerpetual {
			continue
		}
		symbols = append(symbols, pair.Symbol)
	}
	go bnc.CacheDepthUpdate(dataDir, cex.SYMBOL_TYPE_UM_FUTURES, symbols...)

	time.Sleep(time.Hour * 2)

	dataDir = filepath.Join(home, "cache", "cex", "bnc", "depth_update_redundancy")
	pairs, err = bnc.GetSpotSymbols()
	if err != nil {
		panic(err)
	}
	symbols = []string{}
	for _, pair := range pairs {
		if !pair.Tradable {
			continue
		}
		symbols = append(symbols, pair.Symbol)
	}
	go bnc.CacheDepthUpdate(dataDir, cex.SYMBOL_TYPE_SPOT, symbols...)

	pairs, err = bnc.GetUMSymbols()
	if err != nil {
		panic(err)
	}
	for _, pair := range pairs {
		if !pair.Tradable {
			continue
		}
		if !pair.IsPerpetual {
			continue
		}
		symbols = append(symbols, pair.Symbol)
	}
	go bnc.CacheDepthUpdate(dataDir, cex.SYMBOL_TYPE_UM_FUTURES, symbols...)

	select {}
}
