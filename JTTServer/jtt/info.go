package jtt

import (
	"bytes"
	"errors"
)

func init() {
	RegisterMarshals(&Marshal{
		Cmd: msgIDTextIssued,
		NewMarshaler: func() Marshaler {
			return textIssuedMarshal
		},
	})
}

// 文本信息下发
func textIssuedMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgTextIssued)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}

	if version2011 == version && nil != output.base() {
		output.base().writeTo(&buf)
	} else {
		output.writeTo(&buf)
	}

	return buf.Bytes(), nil
}
