package bnc

import (
	"log/slog"
	"os"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/cex-go/ob"
)

type PublicWs struct {
	klClt *KlineWs
	agClt *AggTradeWs
	obClt *OrderBookWs
}

func NewPublicWs(logger *slog.Logger) *PublicWs {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	klClt := NewKlineWs(logger)
	agClt := NewAggTradeWs(logger)
	obClt := NewOrderBookWs(logger)
	return &PublicWs{
		klClt: klClt,
		agClt: agClt,
		obClt: obClt,
	}
}

func (w *PublicWs) SubKlines(sybType cex.SymbolType, interval KlineInterval, symbols ...string) (unsubed []string, err error) {
	return w.klClt.Sub(sybType, interval, symbols...)
}

func (w *PublicWs) UnsubKlines(sybType cex.SymbolType, interval KlineInterval, symbols ...string) (err error) {
	return w.klClt.Unsub(sybType, interval, symbols...)
}

func (w *PublicWs) SubAggTrades(sybType cex.SymbolType, symbols ...string) (unsubed []string, err error) {
	return w.agClt.Sub(sybType, symbols...)
}

func (w *PublicWs) UnsubAggTrades(sybType cex.SymbolType, symbols ...string) (err error) {
	return w.agClt.Unsub(sybType, symbols...)
}

func (w *PublicWs) SubOrderBooks(sybType cex.SymbolType, symbols ...string) (unsubed []string, err error) {
	return w.obClt.Sub(sybType, symbols...)
}

func (w *PublicWs) UnsubOrderBooks(sybType cex.SymbolType, symbols ...string) (err error) {
	return w.obClt.Unsub(sybType, symbols...)
}

func (w *PublicWs) NewKlineCh(sybType cex.SymbolType, interval KlineInterval, symbol string) (ch <-chan KlineMsg, err error) {
	return w.klClt.NewCh(sybType, interval, symbol)
}

func (w *PublicWs) NewAggTradeCh(sybType cex.SymbolType, symbol string) (ch <-chan AggTradeMsg, err error) {
	return w.agClt.NewCh(sybType, symbol)
}

func (w *PublicWs) NewOrderBookCh(sybType cex.SymbolType, symbol string) (ch <-chan ob.Data, err error) {
	return w.obClt.NewCh(sybType, symbol)
}

func (w *PublicWs) RemoveKlineCh(ch <-chan KlineMsg) {
	w.klClt.RemoveCh(ch)
}

func (w *PublicWs) RemoveAggTradeCh(ch <-chan AggTradeMsg) {
	w.agClt.RemoveCh(ch)
}

func (w *PublicWs) RemoveOrderBookCh(ch <-chan ob.Data) {
	w.obClt.RemoveCh(ch)
}

func (w *PublicWs) Close() {
	w.klClt.Close()
	w.agClt.Close()
	w.obClt.Close()
}
