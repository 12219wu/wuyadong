package jtt

import (
	"bytes"
	"errors"
)

func init() {
	RegisterUnmarshals(&Unmarshal{
		Cmd: msgIDGetAreaResp,
		NewUnmarshaler: func() Unmarshaler {
			return getAreaOrPathRespUnmarshal
		},
	})

	RegisterMarshals(&Marshal{
		Cmd: msgIDSetRoundArea,
		NewMarshaler: func() Marshaler {
			return setRoundAreaMarshal
		},
	}, &Marshal{
		Cmd: msgIDDeleteRoundArea,
		NewMarshaler: func() Marshaler {
			return deleteRoundAreaMarshal
		},
	}, &Marshal{
		Cmd: msgIDSetRectArea,
		NewMarshaler: func() Marshaler {
			return setRectAreaMarshal
		},
	}, &Marshal{
		Cmd: msgIDDeleteRectArea,
		NewMarshaler: func() Marshaler {
			return deleteRectAreaMarshal
		},
	}, &Marshal{
		Cmd: msgIDSetPolygonArea,
		NewMarshaler: func() Marshaler {
			return setPolygonAreaMarshal
		},
	}, &Marshal{
		Cmd: msgIDDeletePolygonArea,
		NewMarshaler: func() Marshaler {
			return deletePolygonAreaMarshal
		},
	}, &Marshal{
		Cmd: msgIDSetPolyline,
		NewMarshaler: func() Marshaler {
			return setPathMarshal
		},
	}, &Marshal{
		Cmd: msgIDDeletePolyline,
		NewMarshaler: func() Marshaler {
			return deletePathMarshal
		},
	}, &Marshal{
		Cmd: msgIDGetArea,
		NewMarshaler: func() Marshaler {
			return getAreaOrPathMarshal
		},
	})
}

// 设置圆形区域
func setRoundAreaMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgRoundAreaSettings)
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

// 删除圆形区域
func deleteRoundAreaMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgRoundAreaDelete)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 设置矩形区域
func setRectAreaMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgRectAreaSettings)
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

// 删除矩形区域
func deleteRectAreaMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgRectAreaDelete)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 设置多边形区域
func setPolygonAreaMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgPolygonAreaSettings)
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

// 删除多边形区域
func deletePolygonAreaMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgPolygonAreaDelete)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 设置路线区域
func setPathMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgPathSettings)
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

// 删除路线区域
func deletePathMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgPathDelete)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 查询区域或线路数据
func getAreaOrPathMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgGetAreaOrPath)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 查询区域或线路数据应答
func getAreaOrPathRespUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgGetAreaOrPathResp

	if buf.Len() < 4 {
		return nil, errors.New("the bad protocol data:body error")
	}

	if version2011 == version {
		msg.base().readBy(buf)
	} else {
		msg.readBy(buf)
	}

	return &msg, nil
}
