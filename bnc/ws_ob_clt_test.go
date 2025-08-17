package bnc

import (
	"context"
	"fmt"
	"testing"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/props"
)

func TestWsObClt(t *testing.T) {
	clt := NewOrderBookClient(context.Background(), nil)
	spSuber, err := clt.Sub(cex.SYMBOL_TYPE_SPOT,
		"BTCUSDT",
		"ETHUSDT",
	)
	props.PanicIfNotNil(err)
	go func() {
		for {
			o, closed, err := spSuber.Receive(context.TODO())
			props.PanicIfNotNil(err)
			if closed {
				panic("closed")
			}
			if o.Err != nil {
				if o.Err == ErrCachingObDepthUpdate {
					continue
				}
				fmt.Println(o.Err)
			} else {
				fmt.Println(o.Type, o.Symbol, o.Asks[0])
			}
		}
	}()
	umSuber, err := clt.Sub(cex.SYMBOL_TYPE_UM_FUTURES,
		"ETHUSDT",
		"BTCUSDT",
	)
	props.PanicIfNotNil(err)
	for {
		o, closed, err := umSuber.Receive(context.TODO())
		props.PanicIfNotNil(err)
		if closed {
			panic("closed")
		}
		if o.Err != nil {
			if o.Err == ErrCachingObDepthUpdate {
				continue
			}
			fmt.Println(o.Err)
		} else {
			fmt.Println(o.Type, o.Symbol, o.Asks[0])
		}
	}
}
