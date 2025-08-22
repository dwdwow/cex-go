package bnc

import "errors"

var ErrUnknownSymbolType = errors.New("bnc: unknown symbol type")
var ErrSymbolNotFound = errors.New("bnc: symbol not found")
var ErrCachingObDepthUpdate = errors.New("bnc: caching ob depth update")
var ErrObQueryLimit = errors.New("bnc: ob query limit")
var ErrObDepthUpdateCacheTooLarge = errors.New("bnc: depth update cache too large")
var ErrWsStreamUnsubed = errors.New("bnc: ws stream unsubed")
var ErrNotAllStreamSubed = errors.New("bnc: not all stream subed")
var ErrWsTooFrequentRequest = errors.New("bnc: ws too frequent request")
var ErrWsTooFrequentRestart = errors.New("bnc: ws too frequent restart")
var ErrWsClientClosed = errors.New("bnc: ws client is closed")
