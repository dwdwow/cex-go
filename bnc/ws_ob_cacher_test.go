package bnc

import (
	"testing"

	"github.com/dwdwow/cex-go"
)

func TestCacheOrderbook(t *testing.T) {
	CacheOrderbook("testdata", cex.SYMBOL_TYPE_SPOT, "BTCUSDT")
}
