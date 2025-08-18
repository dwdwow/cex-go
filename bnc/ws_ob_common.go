package bnc

import (
	"time"

	"github.com/dwdwow/limiter-go"
)

// var lastObQueryFailTsMilli = props.SafeRWData[int64]{}

var publicRestLimitter = limiter.New(time.Minute, 1000)

const maxObWsStream = 20
const maxWsChCap = 10000
