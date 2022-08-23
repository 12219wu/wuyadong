// Package jtt 信息采集类协议
package jtt

import (
	"bytes"
	"errors"
)

func init() {
	RegisterUnmarshals(&Unmarshal{
		Cmd: msgIDWaybillReport,
		NewUnmarshaler: func() Unmarshaler {
			return waybillReportUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDDriverIdentityReport,
		NewUnmarshaler: func() Unmarshaler {
			return driverIdentityReportUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDDrivingRecordReport,
		NewUnmarshaler: func() Unmarshaler {
			return drivingRecordReportUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDCANDataReport,
		NewUnmarshaler: func() Unmarshaler {
			return canDataReportUnmarshal
		},
	})

	RegisterMarshals(&Marshal{
		Cmd: msgIDGetDriverIdentity,
		NewMarshaler: func() Marshaler {
			return getDriverIdentityMarshal
		},
	}, &Marshal{
		Cmd: msgIDGatherDrivingRecord,
		NewMarshaler: func() Marshaler {
			return gatherDrivingRecordMarshal
		},
	}, &Marshal{
		Cmd: msgIDDriRecordParamsIssued,
		NewMarshaler: func() Marshaler {
			return drivingRecordParamsIssuedMarshal
		},
	})
}

// 电子运单上报
func waybillReportUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgWaybillReport

	if buf.Len() < 1 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 上报驾驶员身份信息请求
func getDriverIdentityMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgGetDriverIdentity)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 驾驶员身份信息上报
func driverIdentityReportUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgICCardReport

	if buf.Len() < 7 {
		return nil, errors.New("the bad protocol data:body error")
	}

	if version2011 == version && nil != msg.base() {
		msg.base().readBy(buf)
	} else {
		msg.readBy(buf)
	}

	return &msg, nil
}

// CAN总线数据上传
func canDataReportUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgCANDataReport

	if buf.Len() < 7 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 行驶记录数据采集命令
func gatherDrivingRecordMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgGatherDrivingRecord)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 行驶记录数据上传
func drivingRecordReportUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgDrivingRecordReport

	if buf.Len() < 3 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 行驶记录参数下传命令
func drivingRecordParamsIssuedMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgDrivingRecordParamsIssued)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}
