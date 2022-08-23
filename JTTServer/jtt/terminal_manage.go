package jtt

import (
	"bytes"
	"errors"
)

func init() {
	RegisterUnmarshals(&Unmarshal{
		Cmd: msgIDTerminalLogin,
		NewUnmarshaler: func() Unmarshaler {
			return terminalLoginUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDTerminalLogout,
		NewUnmarshaler: func() Unmarshaler {
			return terminalLogoutUnmarshal
		},
	}, &Unmarshal{
		Cmd: MsgIDTerminalAuth,
		NewUnmarshaler: func() Unmarshaler {
			return terminalAuthUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDGetTerminalParamsResp,
		NewUnmarshaler: func() Unmarshaler {
			return getTerminalParamsRespUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDGetTerminalAttrResp,
		NewUnmarshaler: func() Unmarshaler {
			return getTerminalAttrRespUnmarshal
		},
	}, &Unmarshal{
		Cmd: msgIDTerminalUpgradeResp,
		NewUnmarshaler: func() Unmarshaler {
			return terminalUpgradeRespUnmarshal
		},
	})

	RegisterMarshals(&Marshal{
		Cmd: msgIDTerminalLoginResp,
		NewMarshaler: func() Marshaler {
			return terminalLoginRespMarshal
		},
	}, &Marshal{
		Cmd: msgIDSetTerminalParams,
		NewMarshaler: func() Marshaler {
			return setTerminalParamsMarshal
		},
	}, &Marshal{
		Cmd: msgIDGetTerminalParams,
		NewMarshaler: func() Marshaler {
			return getTerminalParamsMarshal
		},
	}, &Marshal{
		Cmd: msgIDGetTerminalSpecParams,
		NewMarshaler: func() Marshaler {
			return getTerSpecParamsMarshal
		},
	}, &Marshal{
		Cmd: msgIDTerminalControl,
		NewMarshaler: func() Marshaler {
			return terminalControlMarshal
		},
	}, &Marshal{
		Cmd: msgIDGetTerminalAttr,
		NewMarshaler: func() Marshaler {
			return getTerminalAttrMarshal
		},
	}, &Marshal{
		Cmd: msgIDTerminalUpgrade,
		NewMarshaler: func() Marshaler {
			return terminalUpgradeMarshal
		},
	})
}

// 终端注册
func terminalLoginUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgTerminalLogin

	if version2011 == version {
		if buf.Len() < 37 {
			return nil, errors.New("the bad protocol data:body error")
		}

		msg.base().readBy(buf)
	} else {
		if buf.Len() < 76 {
			return nil, errors.New("the bad protocol data:body error")
		}

		msg.readBy(buf)
	}

	return &msg, nil
}

// 终端注销
func terminalLogoutUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgTerminalLogout

	if buf.Len() > 0 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 终端鉴权
func terminalAuthUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgTerminalAuth

	if version2011 == version {
		if buf.Len() < 1 {
			return nil, errors.New("the bad protocol data:body error")
		}

		msg.base().readBy(buf)
	} else {
		if buf.Len() < 36 {
			return nil, errors.New("the bad protocol data:body error")
		}

		msg.readBy(buf)
	}

	return &msg, nil
}

// 查询终端参数应答
func getTerminalParamsRespUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgGetTerminalParamsResp

	if buf.Len() < 3 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 查询终端属性应答
func getTerminalAttrRespUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgGetTerminalAttrResp

	if version2011 == version {
		if buf.Len() < 48 {
			return nil, errors.New("the bad protocol data:body error")
		}

		msg.base().readBy(buf)
	} else {
		if buf.Len() < 87 {
			return nil, errors.New("the bad protocol data:body error")
		}

		msg.readBy(buf)
	}

	return &msg, nil
}

// 终端升级结果应答
func terminalUpgradeRespUnmarshal(buf *bytes.Buffer, version byte) (Input, error) {
	var msg MsgTerminalUpgradeResp

	if buf.Len() != 2 {
		return nil, errors.New("the bad protocol data:body error")
	}

	msg.readBy(buf)

	return &msg, nil
}

// 终端注册应答
func terminalLoginRespMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgTerminalLoginResp)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 设置终端参数
func setTerminalParamsMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgTerParamsSettings)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 查询终端参数
func getTerminalParamsMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgGetTerminalParams)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 查询指定终端参数
func getTerSpecParamsMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgGetTerSpecParams)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 终端控制
func terminalControlMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgTerminalControl)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 查询终端属性
func getTerminalAttrMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgGetTerminalAttr)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}

// 终端升级
func terminalUpgradeMarshal(output Output, version byte) ([]byte, error) {
	var buf bytes.Buffer

	_, ok := output.(*MsgTerminalUpgrade)
	if !ok {
		return nil, errors.New("消息体数据与消息ID不符")
	}
	output.writeTo(&buf)

	return buf.Bytes(), nil
}
