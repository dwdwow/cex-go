package bnc

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/limiter-go"
)

func CacheDepthUpdate(dataDir string, symbolType cex.SymbolType, symbols ...string) {
	lim := limiter.New(time.Second, 1)
	slog.Info("bnc: cache orderbook", "dataDir", dataDir, "symbolType", symbolType, "symbols", symbols)
	err := os.MkdirAll(dataDir, 0777)
	if err != nil {
		panic(err)
	}
	slightWs, err := StartNewPublicSlightWsClient(symbolType)
	if err != nil {
		panic(err)
	}
	res, err := slightWs.SubDepthUpdateStream100ms(symbols...)
	if err != nil {
		fmt.Println(res)
		panic(err)
	}

	for _, symbol := range symbols {
		sub := slightWs.SubDepthUpdate100ms(symbol)
		ch := sub.Chan()
		depthDir := filepath.Join(dataDir, "depth", string(symbolType), symbol)
		os.MkdirAll(depthDir, 0777)
		obDir := filepath.Join(dataDir, "ob", string(symbolType), symbol)
		os.MkdirAll(obDir, 0777)
		go func() {
			go func() {
				for {
					lim.Wait(context.Background())
					o, err := GetOrderbook(symbolType, ParamsOrderBook{
						Symbol: symbol,
						Limit:  1000,
					})
					if err != nil {
						slog.Error("bnc: get orderbook failed", "err", err)
						continue
					}
					data, err := json.Marshal(o)
					if err != nil {
						slog.Error("bnc: marshal orderbook to json failed", "err", err, "ob", o)
						continue
					}
					os.WriteFile(filepath.Join(obDir, fmt.Sprintf("%d.json", time.Now().UnixMilli())), data, 0777)
					<-time.After(time.Minute)
				}
			}()
			for msg := range ch {
				if msg.Err != nil {
					slog.Error("bnc: cache orderbook failed", "err", msg.Err)
					continue
				}
				data, err := json.Marshal(msg.Data)
				if err != nil {
					slog.Error("bnc: marshal orderbook to json failed", "err", err, "ob", msg.Data)
					continue
				}
				os.WriteFile(filepath.Join(depthDir, fmt.Sprintf("%d.json", msg.Data.EventTime)), data, 0777)
			}
		}()
	}
	select {}
}
