package client

import (
	"common/protocol"
	"fmt"
	"github.com/go-netty/go-netty"
)

func SetUpCodec() {
	setupCodec := func(channel netty.Channel) {
		channel.Pipeline().
			AddLast(protocol.DelimiterCodec(0x7E, true, 2048))
	}
	fmt.Println(&setupCodec)
}
