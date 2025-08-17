package bnc

import (
	"context"
	"fmt"
	"testing"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/cex-go/ob"
	"github.com/dwdwow/props"
	"github.com/dwdwow/spub"
)

func TestUMObWs(t *testing.T) {
	publisher := spub.NewSimplePublisher(ob.NewSimplePublisherChannelUtil(), spub.SimpleRcvCapOption[ob.Data](100))
	err := publisher.Start(context.TODO())
	props.PanicIfNotNil(err)
	wsCex := NewWsUMObMsgHandler(nil)
	obWs := ob.NewProducer(wsCex, publisher, nil)
	err = obWs.Start(context.TODO())
	props.PanicIfNotNil(err)
	id, err := ob.ID(cex.BINANCE, cex.SYMBOL_TYPE_UM_FUTURES, "ETHUSDT")
	props.PanicIfNotNil(err)
	sub, err := publisher.Subscribe(context.TODO(), id)
	props.PanicIfNotNil(err)
	for {
		o, closed, err := sub.Receive(context.TODO())
		props.PanicIfNotNil(err)
		if closed {
			panic("closed")
		}
		if o.Err != nil {
			fmt.Println(o.Err)
		} else {
			fmt.Println(o.Asks[0])
		}
	}
}
