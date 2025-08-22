package bnc

import (
	"log/slog"
	"os"
	"sync"

	"github.com/dwdwow/cex-go"
	"github.com/dwdwow/props"
)

type PublicStreamMsg[S any] struct {
	SymbolType cex.SymbolType
	Stream     S
	Err        error
}

type PublicStreamHandler[S any] interface {
	Run() error
	Sub(symbols ...string) (unsubed []string, err error)
	Unsub(symbols ...string) (err error)
	NewCh(symbol string) <-chan PublicStreamMsg[S]
	RemoveCh(ch <-chan PublicStreamMsg[S])
	Streams() []string
	Close()
}

type PublicStreamHandlerKit[S any] interface {
	New(symbolType cex.SymbolType, logger *slog.Logger) PublicStreamHandler[S]
	Stream(symbolType cex.SymbolType, symbol string) string
}

type PublicStreamMergedHandler[S any] struct {
	mu sync.Mutex

	symbolType cex.SymbolType

	kit        PublicStreamHandlerKit[S]
	clts       []PublicStreamHandler[S]
	wsByStream *props.SafeRWMap[string, PublicStreamHandler[S]]

	logger *slog.Logger
}

func NewPublicStreamMergedHandler[S any](symbolType cex.SymbolType, kit PublicStreamHandlerKit[S], logger *slog.Logger) *PublicStreamMergedHandler[S] {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("ws", "bnc_public_stream_merged_handler")
	mw := &PublicStreamMergedHandler[S]{
		symbolType: symbolType,
		kit:        kit,
		wsByStream: props.NewSafeRWMap[string, PublicStreamHandler[S]](),
		logger:     logger,
	}
	return mw
}

func (mw *PublicStreamMergedHandler[S]) ScanAllStream() {
	for _, clt := range mw.clts {
		for _, s := range clt.Streams() {
			mw.wsByStream.SetIfNotExists(s, clt)
		}
	}
}

func (mw *PublicStreamMergedHandler[S]) Sub(symbols ...string) (unsubed []string, err error) {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	defer mw.ScanAllStream()
	for _, s := range symbols {
		if !mw.wsByStream.HasKey(mw.kit.Stream(mw.symbolType, s)) {
			unsubed = append(unsubed, s)
		}
	}
	for _, clt := range mw.clts {
		if len(unsubed) == 0 {
			return
		}
		unsubed, err = clt.Sub(unsubed...)
		if err != nil && err != ErrNotAllStreamSubed {
			return
		}
	}
	for {
		if len(unsubed) == 0 {
			return
		}
		clt := mw.kit.New(mw.symbolType, mw.logger)
		err = clt.Run()
		if err != nil {
			return
		}
		unsubed, err = clt.Sub(unsubed...)
		if err != nil && err != ErrNotAllStreamSubed {
			return
		}
		mw.clts = append(mw.clts, clt)
	}
}

func (mw *PublicStreamMergedHandler[S]) Unsub(symbols ...string) (err error) {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	defer mw.ScanAllStream()
	for _, s := range symbols {
		clt, ok := mw.wsByStream.GetVWithOk(mw.kit.Stream(mw.symbolType, s))
		if ok {
			err = clt.Unsub(s)
			if err != nil {
				return
			}
		}
	}
	return nil
}

func (mw *PublicStreamMergedHandler[S]) NewCh(symbol string) (ch <-chan PublicStreamMsg[S], err error) {
	clt, ok := mw.wsByStream.GetVWithOk(mw.kit.Stream(mw.symbolType, symbol))
	if !ok {
		return nil, ErrSymbolNotFound
	}
	return clt.NewCh(symbol), nil
}

func (mw *PublicStreamMergedHandler[S]) RemoveCh(ch <-chan PublicStreamMsg[S]) {
	for _, clt := range mw.clts {
		clt.RemoveCh(ch)
	}
}

func (mw *PublicStreamMergedHandler[S]) Close() {
	for _, clt := range mw.clts {
		clt.Close()
	}
}

type PublicStream[S any] struct {
	spClt *PublicStreamMergedHandler[S]
	umClt *PublicStreamMergedHandler[S]
	cmClt *PublicStreamMergedHandler[S]
}

func NewPublicStream[S any](kit PublicStreamHandlerKit[S], logger *slog.Logger) *PublicStream[S] {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	logger = logger.With("ws", "bnc_public_stream")
	return &PublicStream[S]{
		spClt: NewPublicStreamMergedHandler(cex.SYMBOL_TYPE_SPOT, kit, logger),
		umClt: NewPublicStreamMergedHandler(cex.SYMBOL_TYPE_UM_FUTURES, kit, logger),
		cmClt: NewPublicStreamMergedHandler(cex.SYMBOL_TYPE_CM_FUTURES, kit, logger),
	}
}

func (aw *PublicStream[S]) Sub(symbolType cex.SymbolType, symbols ...string) (unsubed []string, err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return aw.spClt.Sub(symbols...)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return aw.umClt.Sub(symbols...)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return aw.cmClt.Sub(symbols...)
	}
	return nil, ErrUnknownSymbolType
}

func (aw *PublicStream[S]) Unsub(symbolType cex.SymbolType, symbols ...string) (err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return aw.spClt.Unsub(symbols...)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return aw.umClt.Unsub(symbols...)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return aw.cmClt.Unsub(symbols...)
	}
	return ErrUnknownSymbolType
}

func (aw *PublicStream[S]) NewCh(symbolType cex.SymbolType, symbol string) (ch <-chan PublicStreamMsg[S], err error) {
	switch symbolType {
	case cex.SYMBOL_TYPE_SPOT:
		return aw.spClt.NewCh(symbol)
	case cex.SYMBOL_TYPE_UM_FUTURES:
		return aw.umClt.NewCh(symbol)
	case cex.SYMBOL_TYPE_CM_FUTURES:
		return aw.cmClt.NewCh(symbol)
	}
	return nil, ErrUnknownSymbolType
}

func (aw *PublicStream[S]) RemoveCh(ch <-chan PublicStreamMsg[S]) {
	aw.spClt.RemoveCh(ch)
	aw.umClt.RemoveCh(ch)
	aw.cmClt.RemoveCh(ch)
}

func (aw *PublicStream[S]) Close() {
	aw.spClt.Close()
	aw.umClt.Close()
	aw.cmClt.Close()
}
