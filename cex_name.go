package cex

import "slices"

type CexName string

const (
	BINANCE CexName = "BINANCE"
)

var cexNames = []CexName{BINANCE}

func NotCexName(name CexName) bool {
	return !slices.Contains(cexNames, name)
}
