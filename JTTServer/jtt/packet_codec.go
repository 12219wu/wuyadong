package jtt

import (
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec"
	"github.com/go-netty/go-netty/utils"
)

// PacketCodec create packet codec
func PacketCodec() codec.Codec {
	return &packetCodec{}
}

type packetCodec struct {
}

// CodecName 编码器名称
func (*packetCodec) CodecName() string {
	return "packet-codec"
}

func (p *packetCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {
	bts := message.([]byte)

	packet := &packet{}
	utils.Assert(packet.unpack(bts))
	ctx.HandleRead(packet)
}

func (p *packetCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	// TODO message类型为output接口，首先判断output是否为补传消息，如果是补传消息就从输出包表中找到对应的包，继续补传，
	// 如果找不到，就会报错。
	// 如果不是补传消息，将消息编码为packet。
	// 判断packet的body是否过大，不超过规定大小，直接由后续Context处理；过大则需要分包发送（循环发送），并将当前包存放
	// 到输出包表中。
	packet := message.(*packet)

	// 检查协议包是否需要分包
	if len(packet.body) > maxBodySize {
		packet.head.attr.subpackage()
		packet.head.pack.index = 1
		packet.head.pack.total = uint16((len(packet.body) + maxBodySize - 1) / maxBodySize)
	}

	// opts := ctx.Channel().Attachment().(Options)
	// packet.setNumber(opts.nextNumber())

	ctx.HandleWrite(packet.pack())
}
