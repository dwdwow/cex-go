package ob

import (
	"errors"
	"strings"
)

type ChannelMarshal[N any] func(Data[N]) (string, error)
type ChannelUnmarshal[N any] func(string) (Data[N], error)

type ChannelUtil struct {
	marshal   ChannelMarshal[any]
	unmarshal ChannelUnmarshal[any]
}

func NewChannelUtil(marshal ChannelMarshal[any], unmarshal ChannelUnmarshal[any]) *ChannelUtil {
	return &ChannelUtil{
		marshal:   marshal,
		unmarshal: unmarshal,
	}
}

func (u *ChannelUtil) Marshal(data Data[any]) (string, error) {
	if u.marshal == nil {
		return "", errors.New("marshal is nil")
	}
	return u.marshal(data)
}

func (u *ChannelUtil) Unmarshal(s string) (Data[any], error) {
	if u.unmarshal == nil {
		return Data[any]{}, errors.New("unmarshal is nil")
	}
	return u.unmarshal(s)
}

func SimplePublisherChannelMarshal(d Data[any]) (string, error) {
	return ID(d.Cex, d.Type, d.Symbol)
}

func SimplePublisherChannelUnmarshal(c string) (Data[any], error) {
	cn, ot, syb, err := ParseID(c)
	od := Data[any]{
		Cex:    cn,
		Type:   ot,
		Symbol: syb,
	}
	return od, err
}

func NewSimplePublisherChannelUtil() *ChannelUtil {
	return NewChannelUtil(SimplePublisherChannelMarshal, SimplePublisherChannelUnmarshal)
}

const RedisPublisherChannelPrefix = "redis_ob_channel:"

func RedisPublisherChannelMarshal(d Data[any]) (string, error) {
	id, err := ID(d.Cex, d.Type, d.Symbol)
	if err != nil {
		return "", err
	}
	return RedisPublisherChannelPrefix + id, nil
}

func RedisPublisherChannelUnmarshal(c string) (Data[any], error) {
	split := strings.Split(c, RedisPublisherChannelPrefix)
	if len(split) != 2 {
		return Data[any]{}, errors.New("ob channel util: invalid channel " + c)
	}
	id := split[1]
	cn, ot, syb, err := ParseID(id)
	od := Data[any]{
		Cex:    cn,
		Type:   ot,
		Symbol: syb,
	}
	return od, err
}

func NewRedisPublisherChannelUtil() *ChannelUtil {
	return NewChannelUtil(RedisPublisherChannelMarshal, RedisPublisherChannelUnmarshal)
}
