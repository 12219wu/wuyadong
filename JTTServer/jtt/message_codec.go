package jtt

import (
	"bytes"
	"fmt"

	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec"
	"github.com/go-netty/go-netty/utils"
)

type (
	// 网络连接事件
	ActiveEvent struct{}

	// 网络断开事件
	InactiveEvent struct {
		netty.Exception
	}
)

// Counter 计数器，用来消息计数
type Counter func() uint16

// presenter 消息接收
type presenter interface {
	// 接收到消息
	onReceive(Input)

	// 事件通知
	onEvent(netty.Event)
}

// MessageCodec create packet codec
func MessageCodec(version byte, phone []byte, counter Counter) codec.Codec {
	utils.AssertIf(nil == counter, "参数[counter]不能为空")
	return &messageCodec{
		version: version,
		phone:   phone,
		count:   counter,
	}
}

type messageCodec struct {
	version byte    // 协议版本号
	phone   []byte  // 电话号码BCD码
	count   Counter // 计数器，计数消息序号
}

// CodecName 编码器名称
func (*messageCodec) CodecName() string {
	return "message-codec"
}

func (m *messageCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {
	packet := message.(*packet)

	// 接收到终端第一条协议消息
	if nil == m.phone {
		m.version = packet.head.version
		if version2011 == m.version {
			m.phone = make([]byte, 10)
			copy(m.phone[4:], packet.head.phone)
		} else {
			m.phone = packet.head.phone
		}
	}

	// 解析协议消息
	var err error
	unmarshaler := NewUnmarshaler(packet.head.id)
	if nil == unmarshaler {
		utils.Assert(fmt.Errorf("协议[%#x]解码器不存在！", packet.head.id))
	}

	buf := bytes.NewBuffer(packet.body)
	input, err := unmarshaler(buf, packet.head.version)
	utils.Assert(err)

	// 将消息id和流水号带上
	input.setIDAndNumber(packet.head.id, packet.head.number)

	ctx.Channel().Attachment().(presenter).onReceive(input)
}

func (m *messageCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	output := message.(Output)

	// 获取对应消息体打包函数
	reqID := output.msgID()
	marshal := NewMarshaler(reqID)
	if nil == marshal {
		utils.Assert(fmt.Errorf("协议[%#x]编码器不存在！", reqID))
	}

	// 打包消息体
	body, err := marshal(output, m.version)
	utils.Assert(err)

	packet := &packet{
		head: head{
			id:      reqID,
			number:  m.count(),
			phone:   m.phone,
			version: m.version,
		},
		body: body,
	}

	// 指定版本号
	if packet.head.version != 0 {
		packet.head.attr.versionTag()
	}

	ctx.HandleWrite(packet)
}

func (m *messageCodec) HandleEvent(ctx netty.EventContext, event netty.Event) {
	ctx.Attachment().(presenter).onEvent(event)
}

func (m *messageCodec) HandleActive(ctx netty.ActiveContext) {
	ctx.Attachment().(presenter).onEvent(ActiveEvent{})
}

func (m *messageCodec) HandleInActive(ctx netty.InactiveContext, ex netty.Exception) {
	ctx.Attachment().(presenter).onEvent(InactiveEvent{Exception: ex})
	ctx.HandleInactive(ex)
}

func (m *messageCodec) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	ctx.Attachment().(presenter).onEvent(ex)
}
