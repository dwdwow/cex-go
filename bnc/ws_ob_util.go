package bnc

import "github.com/dwdwow/props"

var lastObQueryFailTsMilli = props.SafeRWData[int64]{}

const maxObWsStream = 50
