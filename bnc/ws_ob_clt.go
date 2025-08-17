package bnc

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/cex-go/ob"
	"github.com/dwdwow/spub"
)

type OrderBookClient struct {
	ctx        context.Context
	spConsumer spub.ConsumerService[ob.Data]
	umConsumer spub.ConsumerService[ob.Data]
	logger     *slog.Logger
}

func NewOrderBookClient(ctx context.Context, logger *slog.Logger) *OrderBookClient {
	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	return &OrderBookClient{
		ctx:    ctx,
		logger: logger,
	}
}

func (o *OrderBookClient) consumer(pairType cex.SymbolType) (consumer spub.ConsumerService[ob.Data], err error) {
	var handler ob.CexWsMsgHandler
	switch pairType {
	case cex.SYMBOL_TYPE_SPOT:
		if o.spConsumer != nil {
			return o.spConsumer, nil
		}
		handler = NewWsSpObMsgHandler(o.logger)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		if o.umConsumer != nil {
			return o.umConsumer, nil
		}
		handler = NewWsUMObMsgHandler(o.logger)
	default:
		return nil, fmt.Errorf("bnc: invalid symbol type: %s", pairType)
	}
	publisher := ob.NewSimplePublisher(handler, o.logger)
	err = publisher.Start(o.ctx)
	if err != nil {
		return nil, err
	}
	consumer = spub.ConsumerService[ob.Data](publisher)
	return
}

func (o *OrderBookClient) Sub(pairType cex.SymbolType, symbols ...string) (subscription spub.Subscription[ob.Data], err error) {
	consumer, err := o.consumer(pairType)
	if err != nil {
		return nil, err
	}
	var cs []string
	for _, symbol := range symbols {
		c, err := consumer.ChannelUtil().Marshal(ob.Data{Cex: cex.BINANCE, Type: pairType, Symbol: symbol})
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	subscription, err = consumer.Subscribe(o.ctx, cs...)
	return
}
