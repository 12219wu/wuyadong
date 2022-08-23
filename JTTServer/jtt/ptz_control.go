package jtt

import (
	"bytes"
	"errors"
)

func init() {
	RegisterMarshals(&Marshal{
		Cmd: msgIDPtzTurn,
		NewMarshaler: func() Marshaler {
			return ptzTurnMarshal
		},
	}, &Marshal{
		Cmd: msgIDPtzFocus,
		NewMarshaler: func() Marshaler {
			return ptzFocusMarshal
		},
	}, &Marshal{
		Cmd: msgIDPtzAperture,
		NewMarshaler: func() Marshaler {
			return ptzApertureMarshal
		},
	}, &Marshal{
		Cmd: msgIDPtzWiper,
		NewMarshaler: func() Marshaler {
			return ptzWiperMarshal
		},
	}, &Marshal{
		Cmd: msgIDPtzZoom,
		NewMarshaler: func() Marshaler {
			return ptzZoomMarshal
		},
	}, &Marshal{
		Cmd: msgIDPtzFilllight,
		NewMarshaler: func() Marshaler {
			return ptzFillLightMarshal
		},
	})
}

func ptzTurnMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgPtzTurn)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

func ptzFocusMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgPtzFocus)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

func ptzApertureMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgPtzAperture)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

func ptzWiperMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgPtzWiper)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

func ptzFillLightMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgPtzFillLight)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

func ptzZoomMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgPtzZoom)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}
