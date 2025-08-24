package bnc

import (
	"fmt"
	"testing"
)

func TestGetSpotRawOrderBook(t *testing.T) {
	orderBook, err := GetSpotRawOrderBook(ParamsOrderBook{
		Symbol: "BTCUSDT",
		Limit:  5000,
	})
	if err != nil {
		t.Errorf("GetSpotOrderBook() error = %v", err)
	}
	fmt.Println(orderBook.Symbol, len(orderBook.Bids))
}

func TestGetSpotOrderBook(t *testing.T) {
	orderBook, err := GetSpotOrderBook(ParamsOrderBook{
		Symbol: "BTCUSDT",
		Limit:  5000,
	})
	if err != nil {
		t.Errorf("GetSpotOrderBook() error = %v", err)
	}
	fmt.Println(orderBook.Symbol, orderBook.Bids[0])
}

func TestGetSpotTrades(t *testing.T) {
	trades, err := GetSpotTrades(ParamsTrades{
		Symbol: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetSpotTrades() error = %v", err)
	}
	fmt.Println(trades[0])
}

func TestGetSpotHistoricalTrades(t *testing.T) {
	trades, err := GetSpotHistoricalTrades(ParamsHistoricalTrades{
		Symbol: "BTCUSDT",
		Limit:  500,
		FromID: 5072786358,
	})
	if err != nil {
		t.Errorf("GetSpotHistoricalTrades() error = %v", err)
	}
	fmt.Println(len(trades), trades[0])
}

func TestGetSpotAggTrades(t *testing.T) {
	trades, err := GetSpotAggTrades(ParamsAggTrades{
		Symbol: "BTCUSDT",
		Limit:  500,
		FromID: 3621471982,
	})
	if err != nil {
		t.Errorf("GetSpotAggTrades() error = %v", err)
	}
	fmt.Println(len(trades), trades[0])
}

func TestGetSpotRawKlines(t *testing.T) {
	klines, err := GetSpotRawKlines(ParamsKlines{
		Symbol:   "BTCUSDT",
		Interval: KLINE_INTERVAL_1m,
	})
	if err != nil {
		t.Errorf("GetSpotKlines() error = %v", err)
	}
	fmt.Println(len(klines), klines[0])
}

func TestGetSpotKlines(t *testing.T) {
	klines, err := GetSpotKlines(ParamsKlines{
		Symbol:   "BTCUSDT",
		Interval: KLINE_INTERVAL_1m,
	})
	if err != nil {
		t.Errorf("GetSpotKlines() error = %v", err)
	}
	fmt.Println(len(klines), klines[0])
}

func TestGetSpotAvgPrice(t *testing.T) {
	avgPrice, err := GetSpotAvgPrice(ParamsAvgPrice{
		Symbol: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetSpotAvgPrice() error = %v", err)
	}
	fmt.Println(avgPrice)
}

func TestGetSpotTicker24Hr(t *testing.T) {
	ticker, err := GetSpotTicker24HrStats(ParamsTicker24HrStats{
		Symbol: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetSpotTicker24Hhr() error = %v", err)
	}
	fmt.Println(ticker)
}

func TestGetSpotTicker24Hrs(t *testing.T) {
	tickers, err := GetSpotTicker24HrStatsList(ParamsTicker24HrStatsList{
		// Symbols: []string{"BTCUSDT", "ETHUSDT"},
	})
	if err != nil {
		t.Errorf("GetSpotTicker24Hrs() error = %v", err)
	}
	fmt.Println(len(tickers), tickers[0])
}

func TestGetSpotTickerTradingDayStats(t *testing.T) {
	ticker, err := GetSpotTickerTradingDayStats(ParamsTickerTradingDayStats{
		Symbol: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetSpotTickerTradingDayStats() error = %v", err)
	}
	fmt.Println(ticker)
}

func TestGetSpotTickerTradingDayStatsList(t *testing.T) {
	tickers, err := GetSpotTickerTradingDayStatsList(ParamsTickerTradingDayStatsList{
		Symbols: []string{"BTCUSDT", "ETHUSDT"},
	})
	if err != nil {
		t.Errorf("GetSpotTickerTradingDayStatsList() error = %v", err)
	}
	fmt.Println(len(tickers), tickers[0])
}

func TestGetSpotTickerPrice(t *testing.T) {
	ticker, err := GetSpotTickerPrice(ParamsTickerPrice{
		Symbol: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetSpotTickerPrice() error = %v", err)
	}
	fmt.Println(ticker)
}

func TestGetSpotTickerPriceList(t *testing.T) {
	tickers, err := GetSpotTickerPriceList(ParamsTickerPriceList{
		Symbols: []string{"BTCUSDT", "ETHUSDT"},
	})
	if err != nil {
		t.Errorf("GetSpotTickerPriceList() error = %v", err)
	}
	fmt.Println(len(tickers), tickers[0])
}

func TestGetSpotOrderBookTicker(t *testing.T) {
	ticker, err := GetSpotOrderBookTicker(ParamsOrderBookTicker{
		Symbol: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetSpotTickerBookTicker() error = %v", err)
	}
	fmt.Println(ticker)
}

func TestGetSpotTickerBookTickerList(t *testing.T) {
	tickers, err := GetSpotTickerBookTickerList(ParamsTickerBookTickerList{
		Symbols: []string{"BTCUSDT", "ETHUSDT"},
	})
	if err != nil {
		t.Errorf("GetSpotTickerBookTickerList() error = %v", err)
	}
	fmt.Println(len(tickers), tickers[0])
}

func TestGetSpotTickerStats(t *testing.T) {
	ticker, err := GetSpotTickerStats(ParamsTickerStats{
		Symbol: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetSpotTickerStats() error = %v", err)
	}
	fmt.Println(ticker)
}

func TestGetSpotTickerStatsList(t *testing.T) {
	tickers, err := GetSpotTickerStatsList(ParamsTickerStatsList{
		Symbols: []string{"BTCUSDT", "ETHUSDT"},
	})
	if err != nil {
		t.Errorf("GetSpotTickerStatsList() error = %v", err)
	}
	fmt.Println(len(tickers), tickers[0])
}
