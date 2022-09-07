package protocol

import "bytes"

// 打包协议表
var marshals = make(map[uint16]*Marshal)

// 解包协议表
var unmarshals = make(map[uint16]*Unmarshal)

// Marshaler 协议打包函数
type Marshaler func(param Output, version byte) ([]byte, error)

// Unmarshaler 协议解包函数
type Unmarshaler func(buf *bytes.Buffer, version byte) (Input, error)

// Marshal 打包协议
type Marshal struct {
	// 命令字
	Cmd uint16
	// 新建打包器
	NewMarshaler func() Marshaler
}

// Unmarshal 解包协议
type Unmarshal struct {
	// 命令字
	Cmd uint16
	// 新建解包器
	NewUnmarshaler func() Unmarshaler
}

// NewMarshaler 获取协议打包函数
func NewMarshaler(cmd uint16) Marshaler {
	marshal := marshals[cmd]
	if nil == marshal {
		return nil
	}
	return marshal.NewMarshaler()
}

// NewUnmarshaler 获取协议解析函数
func NewUnmarshaler(cmd uint16) Unmarshaler {
	unmarshal := unmarshals[cmd]
	if nil == unmarshal {
		return nil
	}
	return unmarshal.NewUnmarshaler()
}

// RegisterMarshals 注册协议
func RegisterMarshals(newMarshals ...*Marshal) {
	for _, marshal := range newMarshals {
		if nil == marshal {
			continue
		}
		marshals[marshal.Cmd] = marshal
	}
}

// RegisterUnmarshals 注册协议
func RegisterUnmarshals(newUnmarshals ...*Unmarshal) {
	for _, unmarshal := range newUnmarshals {
		if nil == unmarshal {
			continue
		}
		unmarshals[unmarshal.Cmd] = unmarshal
	}
}
