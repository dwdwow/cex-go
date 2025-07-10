package bnc

import (
	"fmt"
	"testing"
)

func TestGetUMServerTime(t *testing.T) {
	serverTime, err := GetUMServerTime()
	if err != nil {
		t.Fatalf("GetUMServerTime failed: %v", err)
	}
	t.Logf("Server time: %v", serverTime.ServerTime)
}

func TestGetUMExchangeInfo(t *testing.T) {
	exchangeInfo, err := GetUMExchangeInfo()
	if err != nil {
		t.Fatalf("GetUMExchangeInfo failed: %v", err)
	}
	fmt.Println(len(exchangeInfo.Assets), len(exchangeInfo.Symbols))
	fmt.Println(exchangeInfo.Assets[0])
	fmt.Printf("%+v\n", exchangeInfo.Symbols[0])
}
