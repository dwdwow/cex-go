package bnc

import (
	"testing"

	"github.com/dwdwow/cex-go"
)

func TestCacheDepthUpdate(t *testing.T) {
	CacheDepthUpdate("mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000", "testdata", cex.SYMBOL_TYPE_SPOT, "BTCUSDT")
}
