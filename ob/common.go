package ob

import (
	"encoding/json"
	"time"

	"github.com/dwdwow/cex-go"
)

func NewData(cexName cex.CexName, obType cex.SymbolType, symbol string) Data {
	return Data{
		Cex:    cexName,
		Type:   obType,
		Symbol: symbol,
		Time:   time.Now().UnixMilli(),
	}
}

func RedisMsgUnmarshal(payload string) (Data, error) {
	o := new(Data)
	err := json.Unmarshal([]byte(payload), o)
	return *o, err
}
