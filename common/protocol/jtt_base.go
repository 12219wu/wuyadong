package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// 协议包体最大字节数
const maxBodySize = int(1<<10 - 1)

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
func (m *msgAttr) setEncryption(enc byte) {
	*m |= msgAttr((enc & 0x07)) << 10
}

// 获取加密方式
func (m *msgAttr) getEncryption() byte {
	return byte(*m & 0x1C00)
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
	if h.attr.isSubpackage() {
		return 21
	}

	if h.attr.hasVersionTag() {
		return 17
	}

	return 12
}

// readBy 从缓存中读取消息头内容
func (h *head) readBy(buf *bytes.Buffer) error {
	l := buf.Len()
	if l < h.len() {
		return errors.New("the bad protocol data: < head.length")
	}

	value := make([]byte, 2)

	buf.Read(value)
	h.id = binary.BigEndian.Uint16(value)

	buf.Read(value)
	h.attr = msgAttr(binary.BigEndian.Uint16(value))

	// 再次验证长度，分包数据，消息头长度大一些
	if l < h.len() {
		return errors.New("the bad protocol data: < head.length")
	}

	if h.attr.hasVersionTag() {
		h.version, _ = buf.ReadByte()
		h.phone = make([]byte, 10)
	} else {
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

	if h.attr.hasVersionTag() {
		buf.WriteByte(h.version)
		buf.Write(h.phone[:10])
	} else {
		buf.Write(h.phone[4:10])
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

// packet 协议单包
type packet struct {
	// 协议头
	head head
	// 协议体
	body []byte
}

// 打包协议包，将协议包打包为字节序
func (p *packet) pack() []byte {
	var buf bytes.Buffer

	// 写入协议头
	p.head.writeTo(&buf)
	// 写入协议体
	buf.Write(p.body)
	// 添加校验码
	addChecksum(&buf)

	return buf.Bytes()
}

// 解包协议包，将字节序解包为协议包
func (p *packet) unpack(bts []byte) error {
	// 校验码
	l := len(bts)
	if !checksum(bts[0:l-1], bts[l-1]) {
		return errors.New("the bad protocol data:checksum error")
	}

	// 解析协议包
	buf := bytes.NewBuffer(bts[:l-1])
	if err := p.head.readBy(buf); nil != err {
		return err
	}
	p.body = buf.Bytes()

	return nil
}

func (p *packet) setNumber(number uint16) {
	p.head.number = number
}

func (p *packet) hashCode() uint32 {
	return uint32(p.head.id)<<16 | uint32(p.head.number)
}

// 获取子包（多协议包时有效）
func (p *packet) subpacket(idx uint16) (*packet, error) {
	size, off := len(p.body), int(idx-1)*maxBodySize
	if off < 0 || off > size {
		return nil, fmt.Errorf("包序号溢出[当前请求：%d，总包个数：%d]", idx, p.head.pack.total)
	}

	// 创建子包
	var packet packet
	packet.head = p.head
	packet.head.pack.index = idx
	if (size - off) < maxBodySize {
		packet.body = make([]byte, size-off)
	} else {
		packet.body = make([]byte, maxBodySize)
	}
	packet.head.attr.setBodySize(uint16(len(packet.body)))

	copy(packet.body, p.body[off:])

	return &packet, nil
}

// addChecksum 添加校验码
func addChecksum(buf *bytes.Buffer) {
	bts := buf.Bytes()
	digit := bts[0]
	for _, bt := range bts[1:] {
		digit ^= bt
	}

	buf.WriteByte(digit)
}

// checksum 检验校验码
func checksum(bts []byte, checksum byte) bool {
	digit := bts[0]
	for _, bt := range bts[1:] {
		digit ^= bt
	}
	return digit == checksum
}
