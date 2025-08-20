package bnc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/dwdwow/cex-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CacheDepthUpdate(mongoUri, dbName string, symbolType cex.SymbolType, symbols ...string) {
	client, err := mongo.Connect(options.Client().ApplyURI(mongoUri))
	if err != nil {
		slog.Error("Can not connect mongo server", "err", err, "uri", mongoUri)
		panic(err)
	}
	db := client.Database(dbName)

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
		go func() {
			go func() {
				coll := db.Collection("ob_" + string(symbolType))
				for {
					publicRestLimitter.Wait(context.Background())
					o, err := GetOrderbook(symbolType, ParamsOrderBook{
						Symbol: symbol,
						Limit:  1000,
					})
					if err != nil {
						slog.Error("bnc: get orderbook failed", "err", err)
						continue
					}
					_, err = coll.InsertOne(context.Background(), o)
					if err != nil {
						slog.Error("bnc: insert orderbook to mongo failed", "err", err, "ob", o)
						continue
					}
					<-time.After(time.Minute)
				}
			}()
			coll := db.Collection("depth_" + symbol + "_" + string(symbolType))
			for msg := range ch {
				if msg.Err != nil {
					slog.Error("bnc: cache depth update failed", "err", msg.Err)
					continue
				}
				_, err = coll.InsertOne(context.Background(), msg.Data)
				if err != nil {
					slog.Error("bnc: insert depth update to mongo failed", "err", err, "depth", msg.Data)
				}
			}
		}()
	}
	select {}
}
