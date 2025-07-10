package bnc

import (
	"fmt"
	"testing"
)

func TestGetCMRawOrderBook(t *testing.T) {
	orderBook, err := GetCMRawOrderBook(ParamsOrderBook{
		Symbol: "BTCUSD_PERP",
		Limit:  10,
	})
	if err != nil {
		t.Errorf("GetCMRawOrderBook() error = %v", err)
	}
	fmt.Println(orderBook)
}

func TestGetCMOrderBook(t *testing.T) {
	orderBook, err := GetCMOrderBook(ParamsOrderBook{
		Symbol: "BTCUSD_PERP",
		Limit:  10,
	})
	if err != nil {
		t.Errorf("GetCMOrderBook() error = %v", err)
	}
	fmt.Println(orderBook)
}

func TestGetCMTrades(t *testing.T) {
	trades, err := GetCMTrades(ParamsTrades{
		Symbol: "BTCUSD_PERP",
		Limit:  10,
	})
	if err != nil {
		t.Errorf("GetCMTrades() error = %v", err)
	}
	fmt.Println(trades)
}

func TestGetCMHistoricalTrades(t *testing.T) {
	trades, err := GetCMHistoricalTrades(ParamsHistoricalTrades{
		Symbol: "BTCUSD_PERP",
		Limit:  10,
	})
	if err != nil {
		t.Errorf("GetCMHistoricalTrades() error = %v", err)
	}
	fmt.Println(trades)
}

func TestGetCMPremiumIndexInfo(t *testing.T) {
	premiumIndexInfo, err := GetCMPremiumIndexInfo(ParamsPremiumIndexInfo{
		Symbol: "BTCUSD_PERP",
	})
	if err != nil {
		t.Errorf("GetCMPremiumIndexInfo() error = %v", err)
	}
	fmt.Println(premiumIndexInfo)
}

func TestGetCMFundingRateHistory(t *testing.T) {
	fundingRate, err := GetCMFundingRateHistory(ParamsFundingRate{
		Symbol: "BTCUSD_PERP",
	})
	if err != nil {
		t.Errorf("GetCMFundingRateHistory() error = %v", err)
	}
	fmt.Println(fundingRate)
}

func TestGetCMFundingInfoList(t *testing.T) {
	fundingInfo, err := GetCMFundingInfoList()
	if err != nil {
		t.Errorf("GetCMFundingInfoList() error = %v", err)
	}
	fmt.Println(fundingInfo)
}

func TestGetCMRawKlines(t *testing.T) {
	klines, err := GetCMRawKlines(ParamsKlines{
		Symbol:   "BTCUSD_PERP",
		Interval: KLINE_INTERVAL_1m,
		Limit:    10,
	})
	if err != nil {
		t.Errorf("GetCMRawKlines() error = %v", err)
	}
	fmt.Println(klines)
}

func TestGetCMKlines(t *testing.T) {
	klines, err := GetCMKlines(ParamsKlines{
		Symbol:   "BTCUSD_PERP",
		Interval: KLINE_INTERVAL_1m,
		Limit:    10,
	})
	if err != nil {
		t.Errorf("GetCMKlines() error = %v", err)
	}
	fmt.Println(klines)
}

func TestGetCMRawContinuousKlines(t *testing.T) {
	klines, err := GetCMRawContinuousKlines(ParamsKlines{
		Symbol:       "BTCUSD_PERP",
		Pair:         "BTCUSD",
		ContractType: CONTRACT_TYPE_PERPETUAL,
		Interval:     KLINE_INTERVAL_1m,
		Limit:        10,
	})
	if err != nil {
		t.Errorf("GetCMRawContinuousKlines() error = %v", err)
	}
	fmt.Println(klines)
}

func TestGetCMContinousKlines(t *testing.T) {
	klines, err := GetCMContinousKlines(ParamsKlines{
		Symbol:       "BTCUSD_PERP",
		Pair:         "BTCUSD",
		ContractType: CONTRACT_TYPE_PERPETUAL,
		Interval:     KLINE_INTERVAL_1m,
		Limit:        10,
	})
	if err != nil {
		t.Errorf("GetCMContinousKlines() error = %v", err)
	}
	fmt.Println(klines)
}

func TestGetCMRawIndexPriceKlines(t *testing.T) {
	klines, err := GetCMRawIndexPriceKlines(ParamsKlines{
		Symbol:   "BTCUSD_PERP",
		Pair:     "BTCUSD",
		Interval: KLINE_INTERVAL_1m,
		Limit:    10,
	})
	if err != nil {
		t.Errorf("GetCMRawIndexPriceKlines() error = %v", err)
	}
	fmt.Println(klines)
}

func TestGetCMIndexPriceKlines(t *testing.T) {
	klines, err := GetCMIndexPriceKlines(ParamsKlines{
		Symbol:   "BTCUSD_PERP",
		Pair:     "BTCUSD",
		Interval: KLINE_INTERVAL_1m,
		Limit:    10,
	})
	if err != nil {
		t.Errorf("GetCMIndexPriceKlines() error = %v", err)
	}
	fmt.Println(klines)
}

func TestGetCMRawMarkPriceKlines(t *testing.T) {
	klines, err := GetCMRawMarkPriceKlines(ParamsKlines{
		Symbol:   "BTCUSD_PERP",
		Interval: KLINE_INTERVAL_1m,
		Limit:    10,
	})
	if err != nil {
		t.Errorf("GetCMRawMarkPriceKlines() error = %v", err)
	}
	fmt.Println(klines)
}

func TestGetCMMarkPriceKlines(t *testing.T) {
	klines, err := GetCMMarkPriceKlines(ParamsKlines{
		Symbol:   "BTCUSD_PERP",
		Interval: KLINE_INTERVAL_1m,
		Limit:    10,
	})
	if err != nil {
		t.Errorf("GetCMMarkPriceKlines() error = %v", err)
	}
	fmt.Println(klines)
}

func TestGetCMRawPremiumIndexKlines(t *testing.T) {
	klines, err := GetCMRawPremiumIndexKlines(ParamsKlines{
		Symbol:   "BTCUSD_PERP",
		Interval: KLINE_INTERVAL_1m,
		Limit:    10,
	})
	if err != nil {
		t.Errorf("GetCMRawPremiumIndexKlines() error = %v", err)
	}
	fmt.Println(klines)
}

func TestGetCMPremiumIndexKlines(t *testing.T) {
	klines, err := GetCMPremiumIndexKlines(ParamsKlines{
		Symbol:   "BTCUSD_PERP",
		Interval: KLINE_INTERVAL_1m,
		Limit:    10,
	})
	if err != nil {
		t.Errorf("GetCMPremiumIndexKlines() error = %v", err)
	}
	fmt.Println(klines)
}

func TestGetCMTicker24HrStats(t *testing.T) {
	ticker, err := GetCMTicker24HrStats(ParamsTicker24HrStats{
		Symbol: "BTCUSD_PERP",
	})
	if err != nil {
		t.Errorf("GetCMTicker24HrStats() error = %v", err)
	}
	fmt.Println(ticker)
}

func TestGetCMTicker24HrStatsList(t *testing.T) {
	ticker, err := GetCMTicker24HrStatsList()
	if err != nil {
		t.Errorf("GetCMTicker24HrStatsList() error = %v", err)
	}
	fmt.Println(ticker)
}

func TestGetCMTickerPrice(t *testing.T) {
	ticker, err := GetCMTickerPrice(ParamsTickerPrice{
		Symbol: "BTCUSD_PERP",
	})
	if err != nil {
		t.Errorf("GetCMTickerPrice() error = %v", err)
	}
	fmt.Println(ticker)
}

func TestGetCMTickerPriceList(t *testing.T) {
	ticker, err := GetCMTickerPriceList()
	if err != nil {
		t.Errorf("GetCMTickerPriceList() error = %v", err)
	}
	fmt.Println(ticker)
}

func TestGetCMOrderBookTicker(t *testing.T) {
	ticker, err := GetCMOrderBookTicker(ParamsOrderBookTicker{
		Symbol: "BTCUSD_PERP",
	})
	if err != nil {
		t.Errorf("GetCMTickerBookTicker() error = %v", err)
	}
	fmt.Println(ticker)
}

func TestGetCMOpenInterest(t *testing.T) {
	openInterest, err := GetCMOpenInterest(ParamsOpenInterest{
		Symbol: "BTCUSD_PERP",
	})
	if err != nil {
		t.Errorf("GetCMOpenInterest() error = %v", err)
	}
	fmt.Println(openInterest)
}

func TestGetCMOpenInterestStats(t *testing.T) {
	openInterestStats, err := GetCMOpenInterestStats(ParamsCMOpenInterestStats{
		Pair:         "BTCUSD",
		ContractType: CONTRACT_TYPE_PERPETUAL,
		Period:       "1d",
		Limit:        10,
	})
	if err != nil {
		t.Errorf("GetCMOpenInterestStats() error = %v", err)
	}
	fmt.Println(openInterestStats)
}

func TestGetCMTopTraderLongShortPositionRatio(t *testing.T) {
	openInterestStats, err := GetCMTopTraderLongShortPositionRatio(ParamsCMLongShortRatio{
		Pair:   "BTCUSD",
		Period: "1d",
		Limit:  10,
	})
	if err != nil {
		t.Errorf("GetCMTopTraderLongShortPositionRatio() error = %v", err)
	}
	fmt.Println(openInterestStats)
}

func TestGetCMTopTraderLongShortAccountRatio(t *testing.T) {
	openInterestStats, err := GetCMTopTraderLongShortAccountRatio(ParamsCMLongShortRatio{
		Pair:   "BTCUSD",
		Period: "1d",
		Limit:  10,
	})
	if err != nil {
		t.Errorf("GetCMTopTraderLongShortAccountRatio() error = %v", err)
	}
	fmt.Println(openInterestStats)
}

func TestGetCMGlobalLongShortAccountRatio(t *testing.T) {
	openInterestStats, err := GetCMGlobalLongShortAccountRatio(ParamsCMLongShortRatio{
		Pair:   "BTCUSD",
		Period: "1d",
		Limit:  10,
	})
	if err != nil {
		t.Errorf("GetCMGlobalLongShortAccountRatio() error = %v", err)
	}
	fmt.Println(openInterestStats)
}

func TestGetCMTakerBuySellVolume(t *testing.T) {
	takerBuySellVolume, err := GetCMTakerBuySellVolume(ParamsCMBuySellVolume{
		Pair:         "BTCUSD",
		ContractType: CONTRACT_TYPE_PERPETUAL,
		Period:       "1d",
		Limit:        10,
	})
	if err != nil {
		t.Errorf("GetCMTakerBuySellVolume() error = %v", err)
	}
	fmt.Println(takerBuySellVolume)
}

func TestGetCMBasisInfoList(t *testing.T) {
	basisInfo, err := GetCMBasisInfoList(ParamsCMBasisInfo{
		Pair:         "BTCUSD",
		ContractType: CONTRACT_TYPE_PERPETUAL,
		Period:       "1d",
		Limit:        10,
	})
	if err != nil {
		t.Errorf("GetCMBasisInfoList() error = %v", err)
	}
	fmt.Println(basisInfo)
}

func TestGetCMIndexPriceConstituentInfo(t *testing.T) {
	indexPriceConstituentInfo, err := GetCMIndexPriceConstituentInfo(ParamsIndexPriceConstituentInfo{
		Symbol: "BTCUSD",
	})
	if err != nil {
		t.Errorf("GetCMIndexPriceConstituentInfo() error = %v", err)
	}
	fmt.Println(indexPriceConstituentInfo)
}
