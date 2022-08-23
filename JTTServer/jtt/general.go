package jtt

import (
	"bytes"
	"errors"
)

func init() {
	RegisterUnmarshals(&Unmarshal{
		Cmd: msgIDTerminalResponse,
		NewUnmarshaler: func() Unmarshaler {
			return terminalResponseUnmarshal
		},
	}, &Unmarshal{
		Cmd: MsgIDTerminalHeartbeat,
		NewUnmarshaler: func() Unmarshaler {
			return terminalHeartbeatUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDServerTime,
		NewUnmarshaler: func() Unmarshaler {
			return getServerTimeUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDTerminalPackResend,
		NewUnmarshaler: func() Unmarshaler {
			return terminalGetSubpacketUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDTerminalRSAPublickey,
		NewUnmarshaler: func() Unmarshaler {
			return terminalRSAPublicKeyUnmarshal
		},
	})
	RegisterMarshals(&Marshal{
		Cmd: msgIDServerResponse,
		NewMarshaler: func() Marshaler {
			return serverResponseMarshal
		},
	}, &Marshal{
		Cmd: msgIDServerTimeResp,
		NewMarshaler: func() Marshaler {
			return getServerTimeRespMarshal
		},
	}, &Marshal{
		Cmd: msgIDServerPackResend,
		NewMarshaler: func() Marshaler {
			return serverGetSubpacketMarshal
		},
	}, &Marshal{
		Cmd: msgIDLinkCheck,
		NewMarshaler: func() Marshaler {
			return linkCheckMarshal
		},
	}, &Marshal{
		Cmd: msgIDPlatformRSAPublickey,
		NewMarshaler: func() Marshaler {
			return serverRSAPublicKeyMarshal
		},
	})
}

// 终端通用应答
func terminalResponseUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgTerminalResponse

	if buf.Len() != 5 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 平台通用应答
func serverResponseMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgServerResponse)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 终端心跳
func terminalHeartbeatUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgTerminalHeartbeat

	if buf.Len() > 0 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 查询服务器时间请求
func getServerTimeUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgGetServerTime

	if buf.Len() > 0 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 终端补传分包请求
func terminalGetSubpacketUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgTerGetSubpacket

	if buf.Len() < 4 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 终端RSA公钥
func terminalRSAPublicKeyUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgTerRSAPublicKey

	if buf.Len() != 132 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 查询服务器时间应答
func getServerTimeRespMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgServerTimeResp)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 服务器补传分包请求
func serverGetSubpacketMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgSerGetSubpacket)
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

// 链路检测
func linkCheckMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgLinkCheck)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 平台RSA公钥
func serverRSAPublicKeyMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgServerRSAPublicKey)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}
