package bnc

import (
	"fmt"
	"testing"
)

func TestPingSpotEndpoint(t *testing.T) {
	err := PingSpotEndpoint()
	if err != nil {
		t.Errorf("SpotEndpointPing() error = %v", err)
	}
}

func TestGetServerTime(t *testing.T) {
	serverTime, err := GetSpotServerTime()
	if err != nil {
		t.Errorf("GetServerTime() error = %v", err)
	}
	fmt.Println(serverTime)
}

func TestGetSpotExchangeInfo(t *testing.T) {
	exchangeInfo, err := GetSpotExchangeInfo(ParamsSpotExchangeInfo{
		Symbols: []string{"BTCUSDT", "ETHUSDT"},
	})
	if err != nil {
		t.Errorf("GetSpotExchangeInfo() error = %v", err)
	}
	fmt.Println(exchangeInfo)
}
