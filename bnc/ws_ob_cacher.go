package bnc

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/dwdwow/cex-go"
)

func CacheOrderbook(dataDir string, symbolType cex.SymbolType, symbols ...string) {
	slog.Info("bnc: cache orderbook", "dataDir", dataDir, "symbolType", symbolType, "symbols", symbols)
	err := os.MkdirAll(dataDir, 0777)
	if err != nil {
		panic(err)
	}
	unsubed, err := DefaultWebsocket().Public().SubOrderBooks(symbolType, symbols...)
	if err != nil {
		slog.Error("bnc: cache orderbook failed", "err", err, "unsubed", unsubed)
		panic(err)
	}
	for _, symbol := range symbols {
		ch, err := DefaultWebsocket().Public().NewOrderBookCh(symbolType, symbol)
		if err != nil {
			panic(err)
		}
		dir := filepath.Join(dataDir, string(symbolType), symbol)
		os.MkdirAll(dir, 0777)
		go func() {
			for ob := range ch {
				data, err := json.Marshal(ob)
				if err != nil {
					slog.Error("bnc: marshal orderbook to json failed", "err", err, "ob", ob)
				}
				os.WriteFile(filepath.Join(dir, fmt.Sprintf("%d.json", ob.UpdateTime)), data, 0777)
			}
		}()
	}
	select {}
}
