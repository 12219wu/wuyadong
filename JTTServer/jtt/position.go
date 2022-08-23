package jtt

import (
	"bytes"
	"errors"
)

func init() {
	RegisterUnmarshals(&Unmarshal{
		Cmd: MsgIDPositionReport,
		NewUnmarshaler: func() Unmarshaler {
			return positionReportUnmarshal
		},
	}, &Unmarshal{
		Cmd: MsgIDPositionBatchReport,
		NewUnmarshaler: func() Unmarshaler {
			return positionBatchReportUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDGetPositionResp,
		NewUnmarshaler: func() Unmarshaler {
			return positionRespUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDBDLocationCheck,
		NewUnmarshaler: func() Unmarshaler {
			return bdLocationCheckUnmarshal
		},
	})

	RegisterMarshals(&Marshal{
		Cmd: msgIDGetPosition,
		NewMarshaler: func() Marshaler {
			return getPositionMarshal
		},
	}, &Marshal{
		Cmd: msgIDTrackControl,
		NewMarshaler: func() Marshaler {
			return trackControlMarshal
		},
	}, &Marshal{
		Cmd: msgIDManualConfirmAlarm,
		NewMarshaler: func() Marshaler {
			return manualConfirmAlarmMarshal
		},
	})
}

// 位置上报
func positionReportUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgPositionReport

	if buf.Len() < 28 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 定位数据批量上传
func positionBatchReportUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgPosBatchReport

	if buf.Len() < 28 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 位置查询应答
func positionRespUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgPositionResp

	if buf.Len() < 30 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 北斗验真上报
func bdLocationCheckUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgBDLocCheck

	if buf.Len() < 32 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 位置信息查询
func getPositionMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgGetPosition)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 临时位置跟踪控制
func trackControlMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgTrackControl)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 人工确认报警
func manualConfirmAlarmMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgManualConfirmAlarm)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}
