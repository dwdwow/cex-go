package bnc

import (
	"testing"

	"github.com/dwdwow/cex-go"
)

func TestCacheDepthUpdate(t *testing.T) {
	CacheDepthUpdate("testdata", cex.SYMBOL_TYPE_SPOT, "BTCUSDT")
}
