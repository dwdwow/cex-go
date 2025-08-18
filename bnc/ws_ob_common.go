package bnc

import (
	"time"

	"github.com/dwdwow/limiter-go"
)

// var lastObQueryFailTsMilli = props.SafeRWData[int64]{}

var publicRestLimitter = limiter.New(time.Second, 5)

const maxObWsStream = 20
const maxWsChCap = 10000
