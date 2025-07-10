package bnc

import (
	"fmt"
	"testing"
)

func TestGetUMRawOrderBook(t *testing.T) {
	orderBook, err := GetUMRawOrderBook(ParamsOrderBook{
		Symbol: "BTCUSDT",
		Limit:  1000,
	})
	if err != nil {
		t.Errorf("GetUMRawOrderBook() error = %v", err)
	}
	fmt.Println(len(orderBook.Bids))
}

func TestGetUMOrderBook(t *testing.T) {
	orderBook, err := GetUMOrderBook(ParamsOrderBook{
		Symbol: "BTCUSDT",
		Limit:  1000,
	})
	if err != nil {
		t.Errorf("GetUMOrderBook() error = %v", err)
	}
	fmt.Println(orderBook.Bids[0])
}

func TestGetUMTrades(t *testing.T) {
	trades, err := GetUMTrades(ParamsTrades{
		Symbol: "BTCUSDT",
		Limit:  1000,
	})
	if err != nil {
		t.Errorf("GetUMTrades() error = %v", err)
	}
	fmt.Println(trades[0])
}

func TestGetUMHistoricalTrades(t *testing.T) {
	trades, err := GetUMHistoricalTrades(ParamsHistoricalTrades{
		Symbol: "BTCUSDT",
		Limit:  500,
	})
	if err != nil {
		t.Errorf("GetUMHistoricalTrades() error = %v", err)
	}
	fmt.Println(trades[0])
}

func TestGetUMAggTrades(t *testing.T) {
	trades, err := GetUMAggTrades(ParamsAggTrades{
		Symbol: "BTCUSDT",
		Limit:  1000,
		FromID: 2781582246,
	})
	if err != nil {
		t.Errorf("GetUMAggTrades() error = %v", err)
	}
	fmt.Println(trades[0])
}

func TestGetUMRawKlines(t *testing.T) {
	trades, err := GetUMRawKlines(ParamsKlines{
		Symbol:   "BTCUSDT",
		Interval: KLINE_INTERVAL_1m,
		Limit:    1000,
	})
	if err != nil {
		t.Errorf("GetUMRawKlines() error = %v", err)
	}
	fmt.Println(trades[0])
}

func TestGetUMKlines(t *testing.T) {
	trades, err := GetUMKlines(ParamsKlines{
		Symbol:   "BTCUSDT",
		Interval: KLINE_INTERVAL_1m,
		Limit:    1000,
	})
	if err != nil {
		t.Errorf("GetUMKlines() error = %v", err)
	}
	fmt.Println(trades[0])
}

func TestGetUMRawContinuousKlines(t *testing.T) {
	trades, err := GetUMRawContinuousKlines(ParamsKlines{
		Pair:         "BTCUSDT",
		ContractType: CONTRACT_TYPE_PERPETUAL,
		Interval:     KLINE_INTERVAL_1m,
		Limit:        1500,
	})
	if err != nil {
		t.Errorf("GetUMRawContinuousKlines() error = %v", err)
	}
	fmt.Println(trades[0])
}

func TestGetUMContinousKlines(t *testing.T) {
	trades, err := GetUMContinousKlines(ParamsKlines{
		Pair:         "BTCUSDT",
		ContractType: CONTRACT_TYPE_PERPETUAL,
		Interval:     KLINE_INTERVAL_1m,
		Limit:        1500,
	})
	if err != nil {
		t.Errorf("GetUMContinousKlines() error = %v", err)
	}
	fmt.Println(trades[0])
}

func TestGetUMRawIndexPirceKlines(t *testing.T) {
	trades, err := GetUMRawIndexPirceKlines(ParamsKlines{
		Pair:     "BTCUSDT",
		Interval: KLINE_INTERVAL_1m,
		Limit:    1000,
	})
	if err != nil {
		t.Errorf("GetUMRawIndexPirceKlines() error = %v", err)
	}
	fmt.Println(trades[0])
}

func TestGetUMIndexPirceKlines(t *testing.T) {
	trades, err := GetUMIndexPirceKlines(ParamsKlines{
		Pair:     "BTCUSDT",
		Interval: KLINE_INTERVAL_1m,
		Limit:    1000,
	})
	if err != nil {
		t.Errorf("GetUMIndexPirceKlines() error = %v", err)
	}
	fmt.Println(trades[0])
}

func TestGetUMRawMarkPriceKlines(t *testing.T) {
	trades, err := GetUMRawMarkPriceKlines(ParamsKlines{
		Symbol:   "BTCUSDT",
		Interval: KLINE_INTERVAL_1m,
		Limit:    1000,
	})
	if err != nil {
		t.Errorf("GetUMRawMarkPriceKlines() error = %v", err)
	}
	fmt.Println(trades[0])
}

func TestGetUMMarkPriceKlines(t *testing.T) {
	trades, err := GetUMMarkPriceKlines(ParamsKlines{
		Symbol:   "BTCUSDT",
		Interval: KLINE_INTERVAL_1m,
		Limit:    1000,
	})
	if err != nil {
		t.Errorf("GetUMMarkPriceKlines() error = %v", err)
	}
	fmt.Println(trades[0])
}

func TestGetUMRawPremiumIndexKlines(t *testing.T) {
	trades, err := GetUMRawPremiumIndexKlines(ParamsKlines{
		Symbol:   "BTCUSDT",
		Interval: KLINE_INTERVAL_1m,
		Limit:    1000,
	})
	if err != nil {
		t.Errorf("GetUMRawPremiumIndexKlines() error = %v", err)
	}
	fmt.Println(trades[0])
}

func TestGetUMPremiumIndexKlines(t *testing.T) {
	trades, err := GetUMPremiumIndexKlines(ParamsKlines{
		Symbol:   "BTCUSDT",
		Interval: KLINE_INTERVAL_1m,
		Limit:    1000,
	})
	if err != nil {
		t.Errorf("GetUMPremiumIndexKlines() error = %v", err)
	}
	fmt.Println(trades[0])
}

func TestGetUMMarkPriceInfo(t *testing.T) {
	markPriceInfo, err := GetUMMarkPriceInfo(ParamsMarkPriceInfo{
		Symbol: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetUMMarkPriceInfo() error = %v", err)
	}
	fmt.Println(markPriceInfo)
}

func TestGetUMAllMarkPriceInfo(t *testing.T) {
	markPriceInfos, err := GetUMAllMarkPriceInfo()
	if err != nil {
		t.Errorf("GetUMAllMarkPriceInfo() error = %v", err)
	}
	fmt.Println(len(markPriceInfos))
	fmt.Println(markPriceInfos[0])
}

func TestGetUMFundingRateHistory(t *testing.T) {
	fundingRate, err := GetUMFundingRateHistory(ParamsFundingRate{
		Symbol: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetUMFundingRate() error = %v", err)
	}
	fmt.Println(fundingRate)
}

func TestGetUMFundingInfoList(t *testing.T) {
	fundingInfo, err := GetUMFundingInfoList()
	if err != nil {
		t.Errorf("GetUMFundingInfoList() error = %v", err)
	}
	fmt.Println(len(fundingInfo))
	fmt.Println(fundingInfo[0])
}

func TestGetUMTicker24HrStats(t *testing.T) {
	ticker, err := GetUMTicker24HrStats(ParamsTicker24HrStats{
		Symbol: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetUMTicker24HrStats() error = %v", err)
	}
	fmt.Println(ticker)
}

func TestGetUMTicker24HrStatsList(t *testing.T) {
	ticker, err := GetUMTicker24HrStatsList()
	if err != nil {
		t.Errorf("GetUMTicker24HrStatsList() error = %v", err)
	}
	fmt.Println(len(ticker))
	fmt.Println(ticker[0])
}

func TestGetUMTickerPrice(t *testing.T) {
	ticker, err := GetUMTickerPrice(ParamsTickerPrice{
		Symbol: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetUMTickerPrice() error = %v", err)
	}
	fmt.Println(ticker)
}

func TestGetUMTickerPriceList(t *testing.T) {
	ticker, err := GetUMTickerPriceList()
	if err != nil {
		t.Errorf("GetUMTickerPriceList() error = %v", err)
	}
	fmt.Println(len(ticker))
	fmt.Println(ticker[0])
}

func TestGetUMTickerBookTicker(t *testing.T) {
	ticker, err := GetUMTickerBookTicker(ParamsTickerBookTicker{
		Symbol: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetUMTickerBookTicker() error = %v", err)
	}
	fmt.Println(ticker)
}

func TestGetUMTickerBookTickerList(t *testing.T) {
	ticker, err := GetUMTickerBookTickerList()
	if err != nil {
		t.Errorf("GetUMTickerBookTickerList() error = %v", err)
	}
	fmt.Println(len(ticker))
	fmt.Println(ticker[0])
}

func TestGetUMDeliveryPriceList(t *testing.T) {
	deliveryPrice, err := GetUMDeliveryPriceList(ParamsDeliveryPrice{
		Pair: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetUMDeliveryPrice() error = %v", err)
	}
	fmt.Println(deliveryPrice)
}

func TestGetUMOpenInterest(t *testing.T) {
	openInterest, err := GetUMOpenInterest(ParamsOpenInterest{
		Symbol: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetUMOpenInterest() error = %v", err)
	}
	fmt.Println(openInterest)
}

func TestGetUMOpenInterestStats(t *testing.T) {
	openInterestStats, err := GetUMOpenInterestStats(ParamsOpenInterestStats{
		Symbol: "BTCUSDT",
		Period: "5m",
	})
	if err != nil {
		t.Errorf("GetUMOpenInterestStats() error = %v", err)
	}
	fmt.Println(openInterestStats)
}

func TestGetUMTopTraderLongShortPositionRatio(t *testing.T) {
	longShortRatio, err := GetUMTopTraderLongShortPositionRatio(ParamsLongShortRatio{
		Symbol: "BTCUSDT",
		Period: "5m",
	})
	if err != nil {
		t.Errorf("GetUMTopTraderLongShortPositionRatio() error = %v", err)
	}
	fmt.Println(longShortRatio)
}

func TestGetUMTopTraderLongShortAccountRatio(t *testing.T) {
	longShortRatio, err := GetUMTopTraderLongShortAccountRatio(ParamsLongShortRatio{
		Symbol: "BTCUSDT",
		Period: "5m",
	})
	if err != nil {
		t.Errorf("GetUMTopTraderLongShortAccountRatio() error = %v", err)
	}
	fmt.Println(longShortRatio)
}

func TestGetUMGlobalLongShortAccountRatio(t *testing.T) {
	longShortRatio, err := GetUMGlobalLongShortAccountRatio(ParamsLongShortRatio{
		Symbol: "BTCUSDT",
		Period: "5m",
	})
	if err != nil {
		t.Errorf("GetUMTopTraderLongShortRatio() error = %v", err)
	}
	fmt.Println(longShortRatio)
}

func TestGetUMTakerBuySellRatio(t *testing.T) {
	takerBuySellRatio, err := GetUMTakerBuySellRatio(ParamsTakerBuySellRatio{
		Symbol: "BTCUSDT",
		Period: "5m",
	})
	if err != nil {
		t.Errorf("GetUMTakerBuySellRatio() error = %v", err)
	}
	fmt.Println(takerBuySellRatio)
}

func TestGetUMFuturesBasis(t *testing.T) {
	futuresBasis, err := GetUMFuturesBasis(ParamsFuturesBasis{
		Pair:         "BTCUSDT",
		ContractType: CONTRACT_TYPE_PERPETUAL,
		Period:       "5m",
		Limit:        30,
	})
	if err != nil {
		t.Errorf("GetUMFuturesBasis() error = %v", err)
	}
	fmt.Println(futuresBasis)
}

func TestGetUMCompositeIndexSymbolInfo(t *testing.T) {
	compositeIndexSymbolInfo, err := GetUMCompositeIndexSymbolInfo(ParamsCompositeIndexSymbolInfo{
		Symbol: "DEFIUSDT",
	})
	if err != nil {
		t.Errorf("GetUMCompositeIndexSymbolInfo() error = %v", err)
	}
	fmt.Println(compositeIndexSymbolInfo)
}

func TestGetUMCompositeIndexSymbolInfoList(t *testing.T) {
	compositeIndexSymbolInfo, err := GetUMCompositeIndexSymbolInfoList()
	if err != nil {
		t.Errorf("GetUMCompositeIndexSymbolInfoList() error = %v", err)
	}
	fmt.Println(len(compositeIndexSymbolInfo))
	fmt.Println(compositeIndexSymbolInfo)
}

func TestGetUMMultiAssetsModeAssetIndex(t *testing.T) {
	multiAssetsModeAssetIndex, err := GetUMMultiAssetsModeAssetIndex(ParamsMultiAssetsModeAssetIndex{
		Symbol: "BTCUSD",
	})
	if err != nil {
		t.Errorf("GetUMMultiAssetsModeAssetIndex() error = %v", err)
	}
	fmt.Println(multiAssetsModeAssetIndex)
}

func TestGetUMMultiAssetsModeAssetIndexList(t *testing.T) {
	multiAssetsModeAssetIndex, err := GetUMMultiAssetsModeAssetIndexList()
	if err != nil {
		t.Errorf("GetUMMultiAssetsModeAssetIndexList() error = %v", err)
	}
	fmt.Println(len(multiAssetsModeAssetIndex))
	fmt.Println(multiAssetsModeAssetIndex)
}

func TestGetUMIndexPriceConstituentInfo(t *testing.T) {
	indexPriceConstituentInfo, err := GetUMIndexPriceConstituentInfo(ParamsIndexPriceConstituentInfo{
		Symbol: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetUMIndexPriceConstituentInfo() error = %v", err)
	}
	fmt.Println(indexPriceConstituentInfo)
}

func TestGetUMInsuranceFundBalanceSnapshot(t *testing.T) {
	insuranceFundBalanceSnapshot, err := GetUMInsuranceFundBalanceSnapshot(ParamsInsuranceFundBalanceSnapshot{
		Symbol: "BTCUSDT",
	})
	if err != nil {
		t.Errorf("GetUMInsuranceFundBalanceSnapshot() error = %v", err)
	}
	fmt.Println(insuranceFundBalanceSnapshot)
}

func TestGetUMInsuranceFundBalanceSnapshotList(t *testing.T) {
	insuranceFundBalanceSnapshot, err := GetUMInsuranceFundBalanceSnapshotList()
	if err != nil {
		t.Errorf("GetUMInsuranceFundBalanceSnapshotList() error = %v", err)
	}
	fmt.Println(len(insuranceFundBalanceSnapshot))
	fmt.Println(insuranceFundBalanceSnapshot)
}
