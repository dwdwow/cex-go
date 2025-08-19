package bnc

import (
	"fmt"
	"testing"
	"time"

	"github.com/dwdwow/props"
)

func TestRawWs(t *testing.T) {
	ws, err := StartNewRawWs(cmPublicWsCfg, nil)
	props.PanicIfNotNil(err)
	go func() {
		for {
			msg := ws.Wait()
			if msg.Err != nil {
				t.Error(msg.Err)
			} else {
				t.Logf("%s", string(msg.Data))
			}
		}
	}()
	var streams []string
	for range 100000 {
		streams = append(streams, "b")
	}
	fmt.Println("send")
	_, err = ws.SendMsg(WsMethodSub, streams)
	props.PanicIfNotNil(err)
	time.Sleep(time.Second * 20)
	fmt.Println("send")
	_, err = ws.SendMsg(WsMethodSub, streams)
	props.PanicIfNotNil(err)
	time.Sleep(time.Second * 20)
	fmt.Println("send")
	_, err = ws.SendMsg(WsMethodSub, streams)
	props.PanicIfNotNil(err)
	time.Sleep(time.Second * 10)
}
