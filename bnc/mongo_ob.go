package bnc

import (
	"context"
	"log/slog"
	"time"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/cex-go/ob"
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

type MongoDbClient struct {
	ctx    context.Context
	cancel context.CancelFunc

	cfg MongoObClientCfg

	clt    *mongo.Client
	db     *mongo.Database
	coll   *mongo.Collection
	obColl *mongo.Collection

	od ob.Data[WsDepthStream]

	logger *slog.Logger
}

func NewMongoDbClient(ctx context.Context, cfg MongoObClientCfg, logger *slog.Logger) *MongoDbClient {
	return &MongoDbClient{
		ctx:    ctx,
		cfg:    cfg,
		logger: logger,
	}
}

func (c *MongoDbClient) Run() (err error) {
	clt, err := mongo.Connect(options.Client().ApplyURI(c.cfg.MongoUri))
	if err != nil {
		return
	}
	c.clt = clt
	c.db = clt.Database(c.cfg.DbName)
	c.coll = c.db.Collection(c.cfg.Symbol)
	c.obColl = c.db.Collection("ob_" + string(c.cfg.SymbolType))

	// res := c.coll.FindOne(c.ctx, bson.D{{"symbol", c.cfg.Symbol}})

	return
}

func (c *MongoDbClient) Close() {
	if c.clt != nil {
		c.clt.Disconnect(c.ctx)
	}
	if c.cancel != nil {
		c.cancel()
	}
}
