package client

import (
	"common/protocol"
	"github.com/go-netty/go-netty"
)

func ChannelInitializer() netty.ChannelInitializer {
	return func(channel netty.Channel) {
		channel.Pipeline().
			AddLast(protocol.DelimiterCodec(0x7E, true, 2048))
	}
}
