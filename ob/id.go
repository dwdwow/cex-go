package ob

import (
	"fmt"
	"strings"

	"github.com/dwdwow/cex-go"
)

const obIDInterval = "_sP#l#It_"

func ID(cexName cex.CexName, obType cex.SymbolType, symbol string) (string, error) {
	if symbol == "USDUSD" {
		return "USDUSD", nil
	}
	if cex.NotCexName(cexName) {
		return "", fmt.Errorf("ob id: invalid cex name %v", cexName)
	}
	if cex.NotSymbolType(obType) {
		return "", fmt.Errorf("ob id: invalid ob type %v", obType)
	}
	if len(symbol) < 3 {
		return "", fmt.Errorf("ob id: invalid symbol %v", symbol)
	}
	return fmt.Sprintf("%v%v%v%v%v", cexName, obIDInterval, obType, obIDInterval, symbol), nil
}

func ParseID(id string) (cexName cex.CexName, obType cex.SymbolType, symbol string, err error) {
	if id == "USDUSD" {
		return "", "", "USDUSD", nil
	}
	l := strings.Split(id, obIDInterval)
	if len(l) != 3 {
		return
	}
	cexName = cex.CexName(l[0])
	obType = cex.SymbolType(l[1])
	symbol = l[2]
	if cex.NotCexName(cexName) {
		err = fmt.Errorf("ob parse id: invalid cex name %v", cexName)
		return
	}
	if cex.NotSymbolType(obType) {
		err = fmt.Errorf("ob parse id: invalid ob type %v", obType)
		return
	}
	if len(symbol) < 3 {
		err = fmt.Errorf("ob parse id: invalid symbol %v", symbol)
		return
	}
	return
}
