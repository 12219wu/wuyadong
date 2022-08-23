package jtt

import (
	"bytes"
	"errors"
)

func init() {
	RegisterUnmarshals(&Unmarshal{
		Cmd: msgIDDataUpPenetrate,
		NewUnmarshaler: func() Unmarshaler {
			return dataUplinkUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDDataCompressionReport,
		NewUnmarshaler: func() Unmarshaler {
			return dataCompressUnmarshal
		},
	})

	RegisterMarshals(&Marshal{
		Cmd: msgIDDataDownPenetrate,
		NewMarshaler: func() Marshaler {
			return dataDownlinkMarshal
		},
	})
}

// 数据上行
func dataUplinkUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgDataUplink

	if buf.Len() < 1 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 数据下行
func dataDownlinkMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgDataDownlink)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 数据压缩上报
func dataCompressUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgDataCompress

	if buf.Len() < 1 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}
