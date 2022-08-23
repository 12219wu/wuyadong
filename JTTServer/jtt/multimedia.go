package jtt

import (
	"bytes"
	"errors"
)

func init() {
	RegisterUnmarshals(&Unmarshal{
		Cmd: msgIDDriverFaceReport,
		NewUnmarshaler: func() Unmarshaler {
			return driverFaceReportUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDMultimediaEventReport,
		NewUnmarshaler: func() Unmarshaler {
			return multimediaEventReportUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDMultimediaDataReport,
		NewUnmarshaler: func() Unmarshaler {
			return multimediaDataReportUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDSnapshotResp,
		NewUnmarshaler: func() Unmarshaler {
			return snapshootRespUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDGetMultimediaSaveInfoResp,
		NewUnmarshaler: func() Unmarshaler {
			return searchLocalMultimediaRespUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDMediaResourceListReport,
		NewUnmarshaler: func() Unmarshaler {
			return mediaResourceListReportUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDMediaPropertyReport,
		NewUnmarshaler: func() Unmarshaler {
			return mediaPropertyReportUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDFileUploadFinish,
		NewUnmarshaler: func() Unmarshaler {
			return fileUploadFinishUnmarshal
		},
	})

	RegisterMarshals(&Marshal{
		Cmd: msgIDMultimediaDataReportResp,
		NewMarshaler: func() Marshaler {
			return multimediaReportRespMarshal
		},
	}, &Marshal{
		Cmd: msgIDSnapshot,
		NewMarshaler: func() Marshaler {
			return snapshootMarshal
		},
	}, &Marshal{
		Cmd: msgIDGetMultimediaSaveInfo,
		NewMarshaler: func() Marshaler {
			return searchLocalMultimediaMarshal
		},
	}, &Marshal{
		Cmd: msgIDGetMultimediaSaveUp,
		NewMarshaler: func() Marshaler {
			return pullLocalMultimediaMarshal
		},
	}, &Marshal{
		Cmd: msgIDGetAudioStartRecord,
		NewMarshaler: func() Marshaler {
			return recordingMarshal
		},
	}, &Marshal{
		Cmd: msgIDGetSingleMediaSaveUp,
		NewMarshaler: func() Marshaler {
			return getLocalMultimediaMarshal
		},
	}, &Marshal{
		Cmd: msgIDRemoteVideoReplay,
		NewMarshaler: func() Marshaler {
			return remoteMediaReplayMarshal
		},
	}, &Marshal{
		Cmd: msgIDRemoteReplayControl,
		NewMarshaler: func() Marshaler {
			return remoteMediaReplayControlMarshal
		},
	}, &Marshal{
		Cmd: msgIDMediaResourceSelect,
		NewMarshaler: func() Marshaler {
			return mediaResourceSelectMarshal
		},
	}, &Marshal{
		Cmd: msgIDQueryMediaProperty,
		NewMarshaler: func() Marshaler {
			return mediaPropertySelectMarshal
		},
	}, &Marshal{
		Cmd: msgIDRealMediaRequest,
		NewMarshaler: func() Marshaler {
			return realMediaRequestMarshal
		},
	}, &Marshal{
		Cmd: msgIDRealMediaControl,
		NewMarshaler: func() Marshaler {
			return realMediaControlMarshal
		},
	}, &Marshal{
		Cmd: msgIDRealMediaNotice,
		NewMarshaler: func() Marshaler {
			return realMediaNoticeMarshal
		},
	}, &Marshal{
		Cmd: msgIDFileUploadCmd,
		NewMarshaler: func() Marshaler {
			return fileUploadCmdMarshal
		},
	}, &Marshal{
		Cmd: msgIDFileUploadCtl,
		NewMarshaler: func() Marshaler {
			return fileUploadCtlMarshal
		},
	})
}

// 多媒体事件信息上传
func multimediaEventReportUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgMultimediaEventReport

	if buf.Len() != 8 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 多媒体数据上传
func multimediaDataReportUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgMultimediaDataReport

	if buf.Len() < 36 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 多媒体数据上传应答
func multimediaReportRespMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgMultimediaReportResp)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 摄像头立即拍摄命令
func snapshootMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgSnapshoot)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 摄像头立即拍摄命令应答
func snapshootRespUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgSnapshootResp

	if buf.Len() < 3 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 存储多媒体数据检索
func searchLocalMultimediaMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgSearchLocalMultimedia)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 存储多媒体数据检索应答
func searchLocalMultimediaRespUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgSearchLocalMultimediaResp

	if buf.Len() < 4 {
		return nil, errors.New("the bad protocol data:body error")
	}

	if version2011 == version && nil != msg.base() {
		msg.base().readBy(buf)
	} else {
		msg.readBy(buf)
	}

	return &msg, nil
}

// 存储多媒体数据上传命令
func pullLocalMultimediaMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgPullLocalMultimedia)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 文件上传完成通知
func fileUploadFinishUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgFileUploadFinish

	if buf.Len() != 3 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 录音命令
func recordingMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgRecording)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 单条存储多媒体数据检索上传命令
func getLocalMultimediaMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgGetLocalMultimedia)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 驾驶员人脸信息采集上报
func driverFaceReportUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgDriverFaceReport

	if buf.Len() < 37 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 平台下发远程录像回放请求
func remoteMediaReplayMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgRemoteVideoReplay)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 平台下发远程录像回放控制
func remoteMediaReplayControlMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgRemoteReplayControl)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 查询音视频资源
func mediaResourceSelectMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgMediaResourceSelect)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 终端上传音视频资源列表
func mediaResourceListReportUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgMediaResourceList

	if buf.Len() < 6 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

//平台下发查询音视频属性
func mediaPropertySelectMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgMediaProperty)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

//终端上传音视频属性
func mediaPropertyReportUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgMediaPropertyReply

	if buf.Len() != 10 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

//平台下发实时音视频传输请求
func realMediaRequestMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgRealMediaRequest)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

//平台下发实时音视频传输控制
func realMediaControlMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgRealMediaControl)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

//平台下发实时音视频传输通知
func realMediaNoticeMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgRealMediaControl)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 文件上传指令
func fileUploadCmdMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgFileUploadCmd)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 文件上传控制
func fileUploadCtlMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgFileUploadCtl)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}
