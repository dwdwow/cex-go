package bnc

import (
	"log/slog"
	"os"
	"sync"
)

type Websocket struct {
	cfg    WebsocketCfg
	pubClt *PublicWs
	logger *slog.Logger
}

func NewWebsocket(cfg WebsocketCfg, logger *slog.Logger) *Websocket {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	pubClt := NewPublicWs(logger)
	return &Websocket{
		cfg:    cfg,
		pubClt: pubClt,
		logger: logger,
	}
}

var muDefaultWebsocketOnce sync.Once
var defaultWebsocket *Websocket

func DefaultWebsocket() *Websocket {
	muDefaultWebsocketOnce.Do(func() {
		defaultWebsocket = NewWebsocket(DefaultWebsocketCfg(), nil)
	})
	return defaultWebsocket
}

func (w *Websocket) Public() *PublicWs {
	return w.pubClt
}

func (w *Websocket) Close() {
	w.pubClt.Close()
}
