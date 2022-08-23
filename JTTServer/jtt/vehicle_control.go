package jtt

import (
	"bytes"
	"errors"
)

func init() {
	RegisterUnmarshals(&Unmarshal{
		Cmd: msgIDVehicleControlResp,
		NewUnmarshaler: func() Unmarshaler {
			return vehicleControlRespUnmarshal
		},
	})

	RegisterMarshals(&Marshal{
		Cmd: msgIDVehicleControl,
		NewMarshaler: func() Marshaler {
			return vehicleControlMarshal
		},
	})
}

// 车辆控制
func vehicleControlMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgVehicleControl)
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

// 车辆控制应答
func vehicleControlRespUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgVehicleControlResp

	if buf.Len() < 30 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}
