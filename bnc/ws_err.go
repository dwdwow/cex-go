package bnc

import "errors"

var ErrCachingObDepthUpdate = errors.New("bnc: caching ob depth update")
var ErrObQueryLimit = errors.New("bnc: ob query limit")
var ErrObDepthUpdateCacheTooLarge = errors.New("bnc: depth update cache too large")
var ErrWsStreamUnsubed = errors.New("bnc: ws stream unsubed")
