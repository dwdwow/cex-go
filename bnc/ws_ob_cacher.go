package bnc

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/dwdwow/cex-go"
)

func CacheDepthUpdate(dataDir string, symbolType cex.SymbolType, symbols ...string) {
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
		dir := filepath.Join(dataDir, string(symbolType), symbol)
		os.MkdirAll(dir, 0777)
		go func() {
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
				os.WriteFile(filepath.Join(dir, fmt.Sprintf("%d.json", msg.Data.EventTime)), data, 0777)
			}
		}()
	}
	select {}
}
