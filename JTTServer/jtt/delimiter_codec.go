package jtt

import (
	"encoding/hex"
	"log"

	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec"
	"github.com/go-netty/go-netty/utils"
)

// DelimiterCodec create delimiter codec
func DelimiterCodec(delimiter byte, stripDelimiter bool, bufSize int) codec.Codec {
	utils.AssertIf(bufSize <= 0, "bufSize must be a positive integer")
	return &delimiterCodec{
		delimiter:      delimiter,
		stripDelimiter: stripDelimiter,
		readBuf:        make([]byte, bufSize),
		// writeBuf:       make([]byte, bufSize),
	}
}

type delimiterCodec struct {
	delimiter      byte   // 定界符
	stripDelimiter bool   // 是否剥离定界符
	readBuf        []byte // 读缓存
	readIdx        int    // 读索引
	// writeBuf       []byte // 写缓存
	// writeIdx       int    // 写索引
}

func (*delimiterCodec) CodecName() string {
	return "delimiter-codec"
}

func (d *delimiterCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {
	// message 为tcptransport，即tcp socket，直接可从里面读数据
	reader := utils.MustToReader(message)
	read := utils.AssertLength(reader.Read(d.readBuf[d.readIdx:]))

	mark, start := 0, -1
	for idx, bt := range d.readBuf[:d.readIdx+read] {
		if d.delimiter == bt {
			if -1 == start {
				start = idx
				continue
			}

			// 取出完整协议包
			mark = idx + 1
			if d.stripDelimiter {
				log.Println("接收到数据包", hex.EncodeToString(d.readBuf[start:mark]))
				ctx.HandleRead(d.readBuf[start+1 : mark-1])
			} else {
				ctx.HandleRead(d.readBuf[start:mark])
			}
			start = -1
		}
	}

	if start < 0 {
		d.readIdx = 0
		return
	}

	// 移动剩余字节到读缓存头部
	d.readIdx = start
	copy(d.readBuf, d.readBuf[start:])
}

func (d *delimiterCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	bts := message.([]byte)
	ctx.HandleWrite([][]byte{{d.delimiter}, bts, {d.delimiter}})
}
