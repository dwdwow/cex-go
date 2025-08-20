package ob

import (
	"encoding/json"
	"time"

	"github.com/dwdwow/cex-go"
)

func NewData[N any](cexName cex.CexName, obType cex.SymbolType, symbol string) Data[N] {
	return Data[N]{
		Cex:        cexName,
		Type:       obType,
		Symbol:     symbol,
		UpdateTime: time.Now().UnixNano(),
	}
}

func RedisMsgUnmarshal[N any](payload string) (Data[N], error) {
	o := new(Data[N])
	err := json.Unmarshal([]byte(payload), o)
	return *o, err
}
