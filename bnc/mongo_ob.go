package bnc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/cex-go/ob"
	"github.com/dwdwow/props"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoObClientCfg struct {
	SymbolType cex.SymbolType
	Symbol     string
	MongoUri   string
	DbName     string
	StartTime  time.Time
}

type MongoObClient struct {
	ctx    context.Context
	cancel context.CancelFunc

	cfg       MongoObClientCfg
	obUpdater obUpdater

	clt  *mongo.Client
	db   *mongo.Database
	coll *mongo.Collection
	cur  *mongo.Cursor

	_cache *props.SafeRWMap[string, []WsDepthStream]
	_exist *props.SafeRWMap[string, bool]
	_ods   *props.SafeRWMap[string, ob.Data[WsDepthStream]]

	logger *slog.Logger
}

func NewMongoObClient(ctx context.Context, cfg MongoObClientCfg, logger *slog.Logger) *MongoObClient {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("bnc_mongo_ob_client", cfg.Symbol)
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	return &MongoObClient{
		ctx:       ctx,
		cancel:    cancel,
		cfg:       cfg,
		obUpdater: defaultObUpdater(cfg.SymbolType),
		_cache:    props.NewSafeRWMap[string, []WsDepthStream](),
		_exist:    props.NewSafeRWMap[string, bool](),
		_ods:      props.NewSafeRWMap[string, ob.Data[WsDepthStream]](),
		logger:    logger,
	}
}

func (c *MongoObClient) symbolType() cex.SymbolType {
	return c.cfg.SymbolType
}

func (c *MongoObClient) cache() *props.SafeRWMap[string, []WsDepthStream] {
	return c._cache
}

func (c *MongoObClient) exist() *props.SafeRWMap[string, bool] {
	return c._exist
}

func (c *MongoObClient) ods() *props.SafeRWMap[string, ob.Data[WsDepthStream]] {
	return c._ods
}

func (c *MongoObClient) Run() (err error) {
	clt, err := mongo.Connect(options.Client().ApplyURI(c.cfg.MongoUri))
	if err != nil {
		return
	}
	c.clt = clt
	c.db = clt.Database(c.cfg.DbName)
	c.coll = c.db.Collection("depth_" + c.cfg.Symbol + "_" + string(c.cfg.SymbolType))
	obColl := c.db.Collection("ob_" + string(c.cfg.SymbolType))
	res := obColl.FindOne(c.ctx, bson.D{
		{Key: "symbol", Value: c.cfg.Symbol},
		{Key: "localTime", Value: c.cfg.StartTime.UnixNano()},
	})
	if res.Err() != nil {
		return
	}
	var o OrderBook
	err = res.Decode(&o)
	if err != nil {
		return
	}
	c._ods.SetKV(c.cfg.Symbol, ob.Data[WsDepthStream]{
		Cex:        cex.BINANCE,
		Type:       c.cfg.SymbolType,
		Symbol:     o.Symbol,
		Version:    strconv.FormatInt(o.LastUpdateID, 10),
		UpdateTime: time.Now().UnixNano(),
		Asks:       o.Asks,
		Bids:       o.Bids,
	})
	c.cur, err = c.coll.Find(c.ctx, bson.D{
		{Key: "s", Value: c.cfg.Symbol},
		{Key: "u", Value: bson.D{
			{Key: "$gte", Value: o.LastUpdateID},
		}},
	})
	if err != nil {
		return
	}
	if c.cur.Next(c.ctx) {
		fmt.Println("true")
	}
	c._exist.SetKV(c.cfg.Symbol, true)
	return
}

func (c *MongoObClient) Read() (obData ob.Data[WsDepthStream], err error) {
	if !c.cur.Next(c.ctx) {
		err = errors.New("bnc: no more depth update data in mongodb")
		return
	}
	var depthData WsDepthStream
	err = c.cur.Decode(&depthData)
	if err != nil {
		return
	}
	buffer := c._cache.GetV(c.cfg.Symbol)
	buffer = append(buffer, depthData)
	c._cache.SetKV(c.cfg.Symbol, buffer)
	obData = c.obUpdater(c, depthData)
	return
}

func (c *MongoObClient) Close() {
	if c.clt != nil {
		c.clt.Disconnect(c.ctx)
	}
	if c.cancel != nil {
		c.cancel()
	}
}
