package main

import "github.com/dwdwow/cex-go/bnc"

func main() {
	bnc.CacheAllSymbolsDepthAndBookTicker()
	select {}
}
