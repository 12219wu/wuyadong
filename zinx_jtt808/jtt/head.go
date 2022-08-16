package jtt

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// 消息体属性
type msgAttr uint16

// 设置消息体大小
func (m *msgAttr) setBodySize(size uint16) {
	*m |= msgAttr(size)
}

// 获取消息体大小
func (m *msgAttr) getBodySize() uint16 {
	return uint16(*m) & 0x03FF
}

// 设置加密方式
func (m *msgAttr) setEncryption(enc uint16) {
	*m |= msgAttr(enc << 10)
}

// subpackage 指定分包
func (m *msgAttr) subpackage() {
	*m |= msgAttr(0x2000)
}

// isSubpackage 是否分包
func (m *msgAttr) isSubpackage() bool {
	return (*m & 0x2000) > 0
}

// versionTag 添加版本标识
func (m *msgAttr) versionTag() {
	*m |= 0x4000
}

// hasVersionTag 是否有版本标识
func (m *msgAttr) hasVersionTag() bool {
	return (*m & 0x4000) > 0
}

// 消息包封装项
type packIndex struct {
	// 消息总包数
	total uint16
	// 包序号
	index uint16
}

// 消息头
type head struct {
	// 消息id
	id uint16
	// 消息体属性
	attr msgAttr
	// 协议版本号
	version byte
	// 终端手机号,10位BCD码，不足前面补0
	phone []byte
	// 消息流水号，从0开始循环累加
	number uint16
	// 仅消息体属性中明确有消息包处理时有效
	pack packIndex
}

// len 消息头大小
func (h *head) len() int {
	//2019
	if h.attr.hasVersionTag() {
		if h.attr.isSubpackage() {
			return 21
		}
		return 17
	}
	//2013
	if h.attr.isSubpackage() {
		return 16
	}
	return 12
}

func (h *head) bytes() []byte {
	var buffer bytes.Buffer
	value := make([]byte, 2)

	binary.BigEndian.PutUint16(value, h.id)
	buffer.Write(value)

	binary.BigEndian.PutUint16(value, uint16(h.attr))
	buffer.Write(value)

	if h.version != 0 {
		buffer.WriteByte(h.version)
	}
	buffer.Write(h.phone[:10])

	binary.BigEndian.PutUint16(value, h.number)
	buffer.Write(value)

	if h.attr.isSubpackage() {
		binary.BigEndian.PutUint16(value, h.pack.total)
		buffer.Write(value)
		binary.BigEndian.PutUint16(value, h.pack.index)
		buffer.Write(value)
	}

	return buffer.Bytes()
}

// readBy 从缓存中读取消息头内容
func (h *head) readBy(buf *bytes.Buffer) error {
	l := buf.Len()
	if l < h.len() {
		return errors.New("The bad protocol data: < head.length")
	}

	value := make([]byte, 2)

	buf.Read(value)
	h.id = binary.BigEndian.Uint16(value)

	buf.Read(value)
	h.attr = msgAttr(binary.BigEndian.Uint16(value))

	// 再次验证长度，分包数据，消息头长度大一些
	if l < h.len() {
		return errors.New("The bad protocol data: < head.length")
	}
	if h.attr.hasVersionTag() {
		h.version, _ = buf.ReadByte()
		h.phone = make([]byte, 10)
	} else {
		h.version = 0
		h.phone = make([]byte, 6)
	}

	buf.Read(h.phone)

	buf.Read(value)
	h.number = binary.BigEndian.Uint16(value)

	if h.attr.isSubpackage() {
		buf.Read(value)
		h.pack.total = binary.BigEndian.Uint16(value)
		buf.Read(value)
		h.pack.index = binary.BigEndian.Uint16(value)
	}

	return nil
}

// 写入缓存中
func (h *head) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 2)

	binary.BigEndian.PutUint16(value, h.id)
	buf.Write(value)

	binary.BigEndian.PutUint16(value, uint16(h.attr))
	buf.Write(value)

	if h.version != 0 {
		buf.WriteByte(h.version)
		buf.Write(h.phone[:10])
	} else {
		buf.Write(h.phone[:6])
	}

	binary.BigEndian.PutUint16(value, h.number)
	buf.Write(value)

	if h.attr.isSubpackage() {
		binary.BigEndian.PutUint16(value, h.pack.total)
		buf.Write(value)
		binary.BigEndian.PutUint16(value, h.pack.index)
		buf.Write(value)
	}
}
