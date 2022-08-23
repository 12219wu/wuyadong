package jtt

import (
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/go-netty/go-netty"
	"github.com/spaolacci/murmur3"
)

// subscriber 终端事件订阅者
type subscriber interface {
	// 终端接入
	onClientConnected(Client)

	// 终端断开
	onClientDisconnected(Client)

	// 终端消息
	onMessage(Client, Input)
}

func sequenceCounter() Counter {
	number := int32(-1)
	return func() uint16 {
		return uint16(atomic.AddInt32(&number, 1))
	}
}

// Client JTT通信协议终端
type Client interface {
	// 终端ID
	ID() uint32

	// 发送消息
	Send(Output) error

	// 本地地址 0.0.0.0:0
	LocalAddr() string

	// 终端地址 0.0.0.0:0
	RemoteAddr() string
}

type client struct {
	id        uint32        // murmurhash值，根据服务地址计算得到
	channel   netty.Channel // 数据通道
	subsriber subscriber    // 终端（连接）事件订阅者
}

func (c *client) ID() uint32 {
	return c.id
}

func (c *client) Send(output Output) error {
	c.channel.Pipeline().FireChannelWrite(output)
	return nil
}

func (c *client) LocalAddr() string {
	return c.channel.LocalAddr()
}

func (c *client) RemoteAddr() string {
	return c.channel.RemoteAddr()
}

func (c *client) onReceive(input Input) {
	c.subsriber.onMessage(c, input)
	// switch msg := input.(type) {
	// case *MsgTerminalAuth:
	// 	resp := NewMsgServerResponse(msg.Number, msg.ID, 0)
	// 	c.Send(resp)
	// case *MsgTerminalHeartbeat:
	// 	resp := NewMsgServerResponse(msg.Number, msg.ID, 0)
	// 	c.Send(resp)
	// case *MsgPositionReport:
	// 	resp := NewMsgServerResponse(msg.Number, msg.ID, 0)
	// 	c.Send(resp)
	// case *MsgPosBatchReport:
	// 	resp := NewMsgServerResponse(msg.Number, msg.ID, 0)
	// 	c.Send(resp)
	// default:
	// 	log.Println("未知消息", input)
	// }
}

func (c *client) onEvent(event netty.Event) {
	switch e := event.(type) {
	case ActiveEvent:
		log.Printf("终端[%s]数据通信已开始", c.RemoteAddr())
		c.subsriber.onClientConnected(c)
	case InactiveEvent:
		log.Printf("终端[%s]数据通信已结束,原因：%s", c.RemoteAddr(), e.Error())
		c.subsriber.onClientDisconnected(c)
	case netty.Exception:
		if _, ok := e.Unwrap().(*net.OpError); ok || io.EOF == e.Unwrap() {
			log.Printf("终端[%s]数据通信已结束,原因：%s", c.RemoteAddr(), e.Error())
			c.subsriber.onClientDisconnected(c)
			c.channel.Close()
		} else {
			log.Println(e)
			log.Println(string(e.Stack()))
		}
	}
}

// 生成哈希
func generateHash(value string) uint32 {
	h32 := murmur3.New32()
	h32.Write([]byte(value))
	return h32.Sum32()
}
