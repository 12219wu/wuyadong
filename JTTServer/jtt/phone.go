package jtt

import (
	"bytes"
	"errors"
)

func init() {
	RegisterMarshals(&Marshal{
		Cmd: msgIDTELCallback,
		NewMarshaler: func() Marshaler {
			return telCallbackMarshal
		},
	}, &Marshal{
		Cmd: msgIDSetContacts,
		NewMarshaler: func() Marshaler {
			return setContactsMarshal
		},
	})
}

// 电话回拨
func telCallbackMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgTELCallback)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 设置电话本
func setContactsMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgContactsSettings)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}
