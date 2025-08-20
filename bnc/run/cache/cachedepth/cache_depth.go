package main

import (
	"strings"
	"time"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/cex-go/bnc"
)

func main() {
	mongoUri := "mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000"
	dbName := "bnc_depth"
	pairs, err := bnc.GetSpotSymbols()
	if err != nil {
		panic(err)
	}
	var symbols []string
	for _, pair := range pairs {
		if !pair.Tradable {
			continue
		}
		if !strings.HasSuffix(pair.Symbol, "USDT") {
			continue
		}
		symbols = append(symbols, pair.Symbol)
	}

	// Split symbols into groups of 30 units
	var symbolGroups [][]string
	for i := 0; i < len(symbols); i += 30 {
		end := min(i+30, len(symbols))
		symbolGroups = append(symbolGroups, symbols[i:end])
	}
	for _, group := range symbolGroups {
		go bnc.CacheDepthUpdate(mongoUri, dbName, cex.SYMBOL_TYPE_SPOT, group...)
	}

	pairs, err = bnc.GetUMSymbols()
	if err != nil {
		panic(err)
	}
	symbols = nil
	for _, pair := range pairs {
		if !pair.Tradable {
			continue
		}
		if !pair.IsPerpetual {
			continue
		}
		if !strings.HasSuffix(pair.Symbol, "USDT") {
			continue
		}
		symbols = append(symbols, pair.Symbol)
	}
	// Split symbols into groups of 30 units
	symbolGroups = nil
	for i := 0; i < len(symbols); i += 30 {
		end := min(i+30, len(symbols))
		symbolGroups = append(symbolGroups, symbols[i:end])
	}
	for _, group := range symbolGroups {
		go bnc.CacheDepthUpdate(mongoUri, dbName, cex.SYMBOL_TYPE_UM_FUTURES, group...)
	}

	time.Sleep(time.Hour * 2)

	dbName = "bnc_depth_redundancy"

	pairs, err = bnc.GetSpotSymbols()
	if err != nil {
		panic(err)
	}
	symbols = nil
	for _, pair := range pairs {
		if !pair.Tradable {
			continue
		}
		if !strings.HasSuffix(pair.Symbol, "USDT") {
			continue
		}
		symbols = append(symbols, pair.Symbol)
	}
	// Split symbols into groups of 30 units
	symbolGroups = nil
	for i := 0; i < len(symbols); i += 30 {
		end := min(i+30, len(symbols))
		symbolGroups = append(symbolGroups, symbols[i:end])
	}
	for _, group := range symbolGroups {
		go bnc.CacheDepthUpdate(mongoUri, dbName, cex.SYMBOL_TYPE_SPOT, group...)
	}

	pairs, err = bnc.GetUMSymbols()
	if err != nil {
		panic(err)
	}
	symbols = nil
	for _, pair := range pairs {
		if !pair.Tradable {
			continue
		}
		if !pair.IsPerpetual {
			continue
		}
		if !strings.HasSuffix(pair.Symbol, "USDT") {
			continue
		}
		symbols = append(symbols, pair.Symbol)
	}
	symbolGroups = nil
	for i := 0; i < len(symbols); i += 30 {
		end := min(i+30, len(symbols))
		symbolGroups = append(symbolGroups, symbols[i:end])
	}
	for _, group := range symbolGroups {
		go bnc.CacheDepthUpdate(mongoUri, dbName, cex.SYMBOL_TYPE_UM_FUTURES, group...)
	}

	select {}
}
