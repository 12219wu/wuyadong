package protocol

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec"
	"github.com/go-netty/go-netty/utils"
)

// EscapeCodec create Escape codec
func EscapeCodec(escapeChar byte, escapeMap map[byte]byte, unescapeChar byte, unescapeMap map[byte]byte) codec.Codec {
	utils.AssertIf(nil == escapeMap, "escapeMap must be not nil")
	utils.AssertIf(nil == unescapeMap, "unescapeMap must be not nil")
	return &escapeCodec{
		escapeChar:   escapeChar,
		unescapeChar: unescapeChar,
		escapeMap:    escapeMap,
		unescapeMap:  unescapeMap,
	}
}

type escapeCodec struct {
	escapeChar   byte          // 转义符
	escapeMap    map[byte]byte // 转义表
	unescapeChar byte          // 反转义符
	unescapeMap  map[byte]byte // 反转义表
	escapeBuf    bytes.Buffer  // 转义缓存
	unescapeBuf  bytes.Buffer  // 反转义缓存
}

func (*escapeCodec) CodecName() string {
	return "escape-codec"
}

func (e *escapeCodec) HandleRead(ctx netty.InboundContext, message netty.Message) {
	bts := message.([]byte)

	hasEscape := false
	for idx, bt := range bts {
		if e.unescapeChar == bt {
			hasEscape = true
			continue
		}

		if hasEscape {
			if bt, ok := e.unescapeMap[bt]; ok {
				e.unescapeBuf.WriteByte(bt)
			} else {
				e.unescapeBuf.Reset()
				panic(fmt.Errorf("接收的数据协议转义报错，数据：%s,位置：%d", hex.EncodeToString(bts), idx+1))
			}

			hasEscape = false
			continue
		}

		e.unescapeBuf.WriteByte(bt)
	}

	ctx.HandleRead(e.unescapeBuf.Bytes())
	e.unescapeBuf.Reset()
}

func (e *escapeCodec) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	bts := message.([]byte)

	for _, bt := range bts {
		if bt, ok := e.escapeMap[bt]; ok {
			e.escapeBuf.WriteByte(e.escapeChar)
			e.escapeBuf.WriteByte(bt)
			continue
		}

		e.escapeBuf.WriteByte(bt)
	}

	ctx.HandleWrite(e.escapeBuf.Bytes())
	e.escapeBuf.Reset()
}
