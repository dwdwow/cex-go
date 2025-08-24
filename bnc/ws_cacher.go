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

func CacheSymbolDepthUpdateAndBookTicker(symbolType cex.SymbolType, symbol, mongoUri, dbName string, depthCh, redunDepthCh <-chan PublicStreamMsg[WsDepthStream], bookTickerCh, redunBookTickerCh <-chan PublicStreamMsg[WsBookTickerStream]) {
	client, err := mongo.Connect(options.Client().ApplyURI(mongoUri))
	if err != nil {
		slog.Error("Can not connect mongo server", "err", err, "uri", mongoUri)
		panic(err)
	}
	db := client.Database(dbName)
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

	go func() {
		coll := db.Collection("depth_" + symbol + "_" + string(symbolType))
		var latestEventTime int64
		var msgs []any
		for {
			var msg PublicStreamMsg[WsDepthStream]
			select {
			case msg = <-depthCh:
			case msg = <-redunDepthCh:
			}
			if msg.Err != nil {
				slog.Error("bnc: cache depth update failed", "err", msg.Err)
				continue
			}
			if msg.Stream.EventTime <= latestEventTime {
				continue
			}
			latestEventTime = msg.Stream.EventTime
			msgs = append(msgs, msg.Stream)
			if len(msgs) > 1000 {
				_, err = coll.InsertMany(context.Background(), msgs)
				if err != nil {
					slog.Error("bnc: insert depth update to mongo failed", "err", err, "depth", msgs)
					panic(err)
				}
				msgs = nil
			}
		}
	}()

	go func() {
		coll := db.Collection("book_ticker_" + symbol + "_" + string(symbolType))
		var latestOrderUpdateId int64
		var msgs []any
		for {
			var msg PublicStreamMsg[WsBookTickerStream]
			select {
			case msg = <-bookTickerCh:
			case msg = <-redunBookTickerCh:
			}
			if msg.Err != nil {
				slog.Error("bnc: cache book ticker failed", "err", msg.Err)
				continue
			}
			if msg.Stream.OrderBookUpdateId <= latestOrderUpdateId {
				continue
			}
			latestOrderUpdateId = msg.Stream.OrderBookUpdateId
			msgs = append(msgs, msg.Stream)
			if len(msgs) > 1000 {
				_, err = coll.InsertMany(context.Background(), msgs)
				if err != nil {
					slog.Error("bnc: insert book ticker to mongo failed", "err", err, "book_ticker", msg.Stream)
					panic(err)
				}
				msgs = nil
			}
		}
	}()
}

func CacheOneTypeAllSymbolsDepthAndBookTicker(symbolType cex.SymbolType) {
	mongoUri := "mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000"
	dbName := "bnc_realtime_cache"

	wsDepth := NewDepthUpdateWs(nil)
	redunWsDepth := NewDepthUpdateWs(nil)
	wsBookTicker := NewBookTickerWs(nil)
	redunWsBookTicker := NewBookTickerWs(nil)

	var pairs []cex.Symbol
	var err error

	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		pairs, err = GetSpotSymbols()
	case cex.SYMBOL_TYPE_UM_FUTURES:
		pairs, err = GetUMSymbols()
	case cex.SYMBOL_TYPE_CM_FUTURES:
		pairs, err = GetCMSymbols()
	default:
		panic("bnc: unknown symbol type")
	}
	if err != nil {
		panic(err)
	}

	var symbols []string
	for _, pair := range pairs {
		if !pair.Tradable {
			continue
		}
		symbols = append(symbols, pair.Symbol)
	}

	unsubed, err := wsDepth.Sub(symbolType, symbols...)
	if err != nil {
		fmt.Println(unsubed)
		panic(err)
	}
	// var mu sync.Mutex
	depthRedunChs := map[string]chan PublicStreamMsg[WsDepthStream]{}
	for _, symbol := range symbols {
		depthRedunChs[symbol] = make(chan PublicStreamMsg[WsDepthStream], 1000)
	}
	go func() {
		time.Sleep(time.Hour)
		unsubed, err := redunWsDepth.Sub(symbolType, symbols...)
		if err != nil {
			fmt.Println(unsubed)
			panic(err)
		}
		for symbol, ch := range depthRedunChs {
			nc, err := redunWsDepth.NewCh(symbolType, symbol)
			if err != nil {
				fmt.Println(symbolType, symbol, "redun depth", err)
				panic(err)
			}
			go func() {
				for msg := range nc {
					ch <- msg
				}
			}()
		}
	}()

	unsubed, err = wsBookTicker.Sub(symbolType, symbols...)
	if err != nil {
		fmt.Println(unsubed)
		panic(err)
	}
	bookTickerRedunChs := map[string]chan PublicStreamMsg[WsBookTickerStream]{}
	for _, symbol := range symbols {
		bookTickerRedunChs[symbol] = make(chan PublicStreamMsg[WsBookTickerStream], 1000)
	}
	go func() {
		time.Sleep(time.Hour)
		unsubed, err := redunWsBookTicker.Sub(symbolType, symbols...)
		if err != nil {
			fmt.Println(unsubed)
			panic(err)
		}
		for symbol, ch := range bookTickerRedunChs {
			nc, err := redunWsBookTicker.NewCh(symbolType, symbol)
			if err != nil {
				fmt.Println(symbolType, symbol, "redun book ticker", err)
				panic(err)
			}
			go func() {
				for msg := range nc {
					ch <- msg
				}
			}()
		}
	}()

	for _, symbol := range symbols {
		chDepth, err := wsDepth.NewCh(symbolType, symbol)
		if err != nil {
			fmt.Println(symbolType, symbol, "depth", err)
			panic(err)
		}
		redunChDepth, ok := depthRedunChs[symbol]
		if !ok {
			fmt.Println(symbolType, symbol, "redun depth not found")
			panic(fmt.Errorf("redun depth not found"))
		}
		chBookTicker, err := wsBookTicker.NewCh(symbolType, symbol)
		if err != nil {
			fmt.Println(symbolType, symbol, "book ticker", err)
			panic(err)
		}
		redunChBookTicker, ok := bookTickerRedunChs[symbol]
		if !ok {
			fmt.Println(symbolType, symbol, "redun book ticker not found")
			panic(fmt.Errorf("redun book ticker not found"))
		}
		CacheSymbolDepthUpdateAndBookTicker(
			symbolType, symbol, mongoUri, dbName,
			chDepth, redunChDepth, chBookTicker, redunChBookTicker,
		)
	}

}

func CacheAllSymbolsDepthAndBookTicker() {
	CacheOneTypeAllSymbolsDepthAndBookTicker(cex.SYMBOL_TYPE_SPOT)
	CacheOneTypeAllSymbolsDepthAndBookTicker(cex.SYMBOL_TYPE_UM_FUTURES)
}

// spot data use microsecond, futures use millisecond, so should check
var ms20000101 = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
var us20000101 = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).UnixMicro()

func CacheOneSymbolAggTrades(symbolType cex.SymbolType, symbol, mongoUri, dbName string, ch, redunCh <-chan PublicStreamMsg[WsAggTradeStream]) {
	go func() {
		client, err := mongo.Connect(options.Client().ApplyURI(mongoUri))
		if err != nil {
			slog.Error("Can not connect mongo server", "err", err, "uri", mongoUri)
			panic(err)
		}
		db := client.Database(dbName)
		coll := db.Collection("agg_trades_" + symbol + "_" + string(symbolType))
		// agg trades are not ordered, so use map to check repetition
		exist := map[int64]int64{}
		var msgs []any
		for {
			var msg PublicStreamMsg[WsAggTradeStream]
			select {
			case msg = <-ch:
			case msg = <-redunCh:
			}
			if msg.Err != nil {
				slog.Error("bnc: cache agg trades failed", "err", msg.Err)
				continue
			}
			if exist[msg.Stream.AggTradeId] > 0 {
				continue
			}
			exist[msg.Stream.AggTradeId] = msg.Stream.EventTime
			msgs = append(msgs, msg.Stream)
			if len(msgs) > 1000 {
				_, err = coll.InsertMany(context.Background(), msgs)
				if err != nil {
					slog.Error("bnc: insert agg trades to mongo failed", "err", err, "agg_trades_len", len(msgs))
					panic(err)
				}
				msgs = nil
			}
			if len(exist) > 10000 {
				now := time.Now()
				newExist := map[int64]int64{}
				for k, v := range exist {
					var ti time.Time
					if v < us20000101 {
						ti = time.UnixMilli(v)
					} else {
						ti = time.UnixMicro(v)
					}
					if now.Sub(ti) < time.Second*10 {
						newExist[k] = v
					}
				}
				exist = newExist
			}
		}
	}()
}

func CacheAggTrades(symbolType cex.SymbolType, symbols []string) {
	mongoUri := "mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000"
	dbName := "bnc_realtime_cache"

	wsAggTrade := NewAggTradeWs(nil)
	redunWsAggTrade := NewAggTradeWs(nil)

	unsubed, err := wsAggTrade.Sub(symbolType, symbols...)
	if err != nil {
		fmt.Println(unsubed)
		panic(err)
	}
	aggTradeRedunChs := map[string]chan PublicStreamMsg[WsAggTradeStream]{}
	for _, symbol := range symbols {
		aggTradeRedunChs[symbol] = make(chan PublicStreamMsg[WsAggTradeStream], 1000)
	}
	go func() {
		time.Sleep(time.Hour)
		unsubed, err := redunWsAggTrade.Sub(symbolType, symbols...)
		if err != nil {
			fmt.Println(unsubed)
			panic(err)
		}
		for symbol, ch := range aggTradeRedunChs {
			nc, err := redunWsAggTrade.NewCh(symbolType, symbol)
			if err != nil {
				fmt.Println(symbolType, symbol, "redun agg trade", err)
				panic(err)
			}
			go func() {
				for msg := range nc {
					ch <- msg
				}
			}()
		}
	}()
	for _, symbol := range symbols {
		ch, err := wsAggTrade.NewCh(symbolType, symbol)
		if err != nil {
			fmt.Println(symbolType, symbol, "agg trade", err)
			panic(err)
		}
		redunCh, ok := aggTradeRedunChs[symbol]
		if !ok {
			fmt.Println(symbolType, symbol, "redun agg trade not found")
			panic(fmt.Errorf("redun agg trade not found"))
		}
		CacheOneSymbolAggTrades(symbolType, symbol, mongoUri, dbName, ch, redunCh)
	}
}
