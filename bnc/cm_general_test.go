package bnc

import (
	"fmt"
	"testing"
)

func TestCheckCMFuturesServerTime(t *testing.T) {
	serverTime, err := CheckCMFuturesServerTime()
	if err != nil {
		t.Errorf("CheckCMFuturesServerTime() error = %v", err)
	}
	fmt.Println(serverTime)
}

func TestGetCMExchangeInfo(t *testing.T) {
	exchangeInfo, err := GetCMExchangeInfo()
	if err != nil {
		t.Errorf("GetCMExchangeInfo() error = %v", err)
	}
	fmt.Println(len(exchangeInfo.Symbols))
	fmt.Printf("%+v\n", exchangeInfo.Symbols[0])
}
