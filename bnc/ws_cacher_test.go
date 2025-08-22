package bnc

import "testing"

func TestCacheAllSymbolsDepthAndBookTicker(t *testing.T) {
	CacheAllSymbolsDepthAndBookTicker()
	select {}
}
