package bnc

import (
	"log/slog"
	"os"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/cex-go/ob"
)

type PublicWs struct {
	klClt *KlineWs
	agClt *PublicStream[WsAggTradeStream]
	dpClt *PublicStream[WsDepthStream]
	btClt *PublicStream[WsBookTickerStream]
	obClt *OrderBookWs
}

func NewPublicWs(logger *slog.Logger) *PublicWs {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	klClt := NewKlineWs(logger)
	agClt := NewAggTradeWs(logger)
	obClt := NewOrderBookWs(logger)
	dpClt := NewDepthUpdateWs(logger)
	btClt := NewBookTickerWs(logger)
	return &PublicWs{
		klClt: klClt,
		agClt: agClt,
		obClt: obClt,
		dpClt: dpClt,
		btClt: btClt,
	}
}

func (w *PublicWs) Close() {
	w.klClt.Close()
	w.agClt.Close()
	w.obClt.Close()
	w.dpClt.Close()
	w.btClt.Close()
}

func (w *PublicWs) SubKlines(sybType cex.SymbolType, interval KlineInterval, symbols ...string) (unsubed []string, err error) {
	return w.klClt.Sub(sybType, interval, symbols...)
}

func (w *PublicWs) UnsubKlines(sybType cex.SymbolType, interval KlineInterval, symbols ...string) (err error) {
	return w.klClt.Unsub(sybType, interval, symbols...)
}

func (w *PublicWs) NewKlineCh(sybType cex.SymbolType, interval KlineInterval, symbol string) (ch <-chan KlineMsg, err error) {
	return w.klClt.NewCh(sybType, interval, symbol)
}

func (w *PublicWs) RemoveKlineCh(ch <-chan KlineMsg) {
	w.klClt.RemoveCh(ch)
}

func (w *PublicWs) SubAggTrades(sybType cex.SymbolType, symbols ...string) (unsubed []string, err error) {
	return w.agClt.Sub(sybType, symbols...)
}

func (w *PublicWs) UnsubAggTrades(sybType cex.SymbolType, symbols ...string) (err error) {
	return w.agClt.Unsub(sybType, symbols...)
}

func (w *PublicWs) NewAggTradeCh(sybType cex.SymbolType, symbol string) (ch <-chan PublicStreamMsg[WsAggTradeStream], err error) {
	return w.agClt.NewCh(sybType, symbol)
}

func (w *PublicWs) RemoveAggTradeCh(ch <-chan PublicStreamMsg[WsAggTradeStream]) {
	w.agClt.RemoveCh(ch)
}

func (w *PublicWs) SubDepthUpdates(sybType cex.SymbolType, symbols ...string) (unsubed []string, err error) {
	return w.dpClt.Sub(sybType, symbols...)
}

func (w *PublicWs) UnsubDepthUpdates(sybType cex.SymbolType, symbols ...string) (err error) {
	return w.dpClt.Unsub(sybType, symbols...)
}

func (w *PublicWs) NewDepthUpdateCh(sybType cex.SymbolType, symbol string) (ch <-chan PublicStreamMsg[WsDepthStream], err error) {
	return w.dpClt.NewCh(sybType, symbol)
}

func (w *PublicWs) RemoveDepthUpdateCh(ch <-chan PublicStreamMsg[WsDepthStream]) {
	w.dpClt.RemoveCh(ch)
}

func (w *PublicWs) SubBookTickers(sybType cex.SymbolType, symbols ...string) (unsubed []string, err error) {
	return w.btClt.Sub(sybType, symbols...)
}

func (w *PublicWs) UnsubBookTickers(sybType cex.SymbolType, symbols ...string) (err error) {
	return w.btClt.Unsub(sybType, symbols...)
}

func (w *PublicWs) NewBookTickerCh(sybType cex.SymbolType, symbol string) (ch <-chan PublicStreamMsg[WsBookTickerStream], err error) {
	return w.btClt.NewCh(sybType, symbol)
}

func (w *PublicWs) RemoveBookTickerCh(ch <-chan PublicStreamMsg[WsBookTickerStream]) {
	w.btClt.RemoveCh(ch)
}

func (w *PublicWs) SubOrderBooks(sybType cex.SymbolType, symbols ...string) (unsubed []string, err error) {
	return w.obClt.Sub(sybType, symbols...)
}

func (w *PublicWs) UnsubOrderBooks(sybType cex.SymbolType, symbols ...string) (err error) {
	return w.obClt.Unsub(sybType, symbols...)
}

func (w *PublicWs) NewOrderBookCh(sybType cex.SymbolType, symbol string) (ch <-chan ob.Data[WsDepthStream], err error) {
	return w.obClt.NewCh(sybType, symbol)
}

func (w *PublicWs) RemoveOrderBookCh(ch <-chan ob.Data[WsDepthStream]) {
	w.obClt.RemoveCh(ch)
}
