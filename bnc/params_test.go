package bnc

import (
	"fmt"
	"testing"
)

func TestStructToSignedQuery(t *testing.T) {
	p := struct{}{}
	v := structToSignedQuery(p, "")
	fmt.Println(v)
}
