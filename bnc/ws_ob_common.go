package bnc

import (
	"time"

	"github.com/dwdwow/limiter-go"
)

// var lastObQueryFailTsMilli = props.SafeRWData[int64]{}

var publicRestLimitter = limiter.New(time.Second, 2)

const maxObWsStream = 30
const maxWsChCap = 10000
const maxSpWsReqPerSec = 5
const maxFuWsReqPerSec = 10
