package jtt

import (
	"JTTServer/util"
	"bytes"
	"encoding/binary"
	"log"
	"reflect"
	"time"

	"github.com/axgle/mahonia"
)

//=====================================接收的消息类=====================================//

// MsgTerminalResponse 终端通用应答消息
type MsgTerminalResponse struct {
	ResponseMark
	// 应答流水号
	Number uint16
	// 应答消息ID
	MsgID uint16
	// 结果
	Result byte
}

// 写入缓存中
func (m *MsgTerminalResponse) readBy(buf *bytes.Buffer) {
	// 应答流水号
	m.Number = binary.BigEndian.Uint16(buf.Next(2))
	// 应答消息ID
	m.MsgID = binary.BigEndian.Uint16(buf.Next(2))
	// 结果
	m.Result, _ = buf.ReadByte()
}

// MsgTerminalHeartbeat 终端心跳
type MsgTerminalHeartbeat struct {
	InputMark
}

func (m *MsgTerminalHeartbeat) readBy(buf *bytes.Buffer) {

}

// MsgGetServerTime 查询服务器时间请求
type MsgGetServerTime struct {
	InputMark
}

func (m *MsgGetServerTime) readBy(buf *bytes.Buffer) {

}

// MsgTerGetSubpacket 终端补传分包请求
type MsgTerGetSubpacket struct {
	InputMark
	// 原始流水号
	ReqNum uint16
	// 重传包ID列表
	IDs []uint16
}

func (m *MsgTerGetSubpacket) readBy(buf *bytes.Buffer) {
	// 原始流水号
	m.ReqNum = binary.BigEndian.Uint16(buf.Next(2))
	// 重传包列表大小
	total := binary.BigEndian.Uint16(buf.Next(2))
	// 重传包ID列表
	for i := uint16(0); i < total; i++ {
		m.IDs = append(m.IDs, binary.BigEndian.Uint16(buf.Next(2)))
	}
}

// MsgTerminalLogin 终端注册
type MsgTerminalLogin struct {
	MsgTerminalLogin2011
}

// 写入缓存中
func (m *MsgTerminalLogin) readBy(buf *bytes.Buffer) {
	// 省域id
	m.Province = binary.BigEndian.Uint16(buf.Next(2))
	// 市域id
	m.City = binary.BigEndian.Uint16(buf.Next(2))
	// 制造商id
	m.Vendor = string(buf.Next(11))
	// 终端型号
	m.Model = string(buf.Next(30))
	// 终端id
	m.TerID = string(buf.Next(30))
	// 车牌颜色
	m.Color, _ = buf.ReadByte()
	// 车牌
	decoder := mahonia.NewDecoder("gbk")
	m.LicencePlate = decoder.ConvertString(buf.String())
}

func (m *MsgTerminalLogin) base() Input {
	return &m.MsgTerminalLogin2011
}

// MsgTerminalLogin2011 终端注册
type MsgTerminalLogin2011 struct {
	InputMark
	// 省域id
	Province uint16
	// 市域id
	City uint16
	// 制造商id
	Vendor string
	// 终端型号
	Model string
	// 终端id
	TerID string
	// 车牌颜色
	Color byte
	// 车牌
	LicencePlate string
}

// 写入缓存中
func (m *MsgTerminalLogin2011) readBy(buf *bytes.Buffer) {
	// 省域id
	m.Province = binary.BigEndian.Uint16(buf.Next(2))
	// 市域id
	m.City = binary.BigEndian.Uint16(buf.Next(2))
	// 制造商id
	m.Vendor = string(buf.Next(5))
	// 终端型号 补充协议-对道路运输车辆卫星定位系统标准20个字节
	m.Model = string(buf.Next(20))
	// 终端id
	m.TerID = string(buf.Next(7))
	// 车牌颜色
	m.Color, _ = buf.ReadByte()
	// 车牌
	decoder := mahonia.NewDecoder("gbk")
	m.LicencePlate = decoder.ConvertString(buf.String())
}

// MsgTerminalLogout 终端注销
type MsgTerminalLogout struct {
	InputMark
}

// 写入缓存中
func (m *MsgTerminalLogout) readBy(buf *bytes.Buffer) {

}

// MsgTerminalAuth 终端鉴权
type MsgTerminalAuth struct {
	MsgTerminalAuth2011
	// 终端IMEI
	IMEI string
	// 软件版本号
	Version string
}

func (m *MsgTerminalAuth) readBy(buf *bytes.Buffer) {
	// 鉴权码，一般是纯字母+数字，可以不用转码
	len, _ := buf.ReadByte()
	m.Token = string(buf.Next(int(len)))
	// IMEI
	m.IMEI = string(buf.Next(15))
	// 软件版本号
	m.Version = string(buf.Next(20))
}

func (m *MsgTerminalAuth) base() Input {
	return &m.MsgTerminalAuth2011
}

// MsgTerminalAuth2011 终端鉴权
type MsgTerminalAuth2011 struct {
	InputMark
	// 鉴权码
	Token string
}

func (m *MsgTerminalAuth2011) readBy(buf *bytes.Buffer) {
	m.Token = buf.String()
}

// MsgGetTerminalParamsResp 查询终端参数应答
type MsgGetTerminalParamsResp struct {
	InputMark
	// 应答流水号
	ReqNum uint16
	// 参数项列表
	Params []Param
}

// 写入缓存中
func (m *MsgGetTerminalParamsResp) readBy(buf *bytes.Buffer) {
	// 应答流水号
	m.ReqNum = binary.BigEndian.Uint16(buf.Next(2))
	// 应答参数个数
	count, _ := buf.ReadByte()
	// 参数项列表
	var param Param
	for i := byte(0); i < count; i++ {
		param.ID = binary.BigEndian.Uint32(buf.Next(4))
		param.Len, _ = buf.ReadByte()
		param.Value = readParamValue(buf, param.Len, paramTypeMap[param.ID])
		if nil != param.Value {
			m.Params = append(m.Params, param)
		}
	}
}

// readParamValue 读取参数值
func readParamValue(buf *bytes.Buffer, len byte, tpName string) interface{} {
	switch tpName {
	case "uint8":
		bt, _ := buf.ReadByte()
		return bt
	case "byte":
		bt, _ := buf.ReadByte()
		return bt
	case "uint32":
		return binary.BigEndian.Uint32(buf.Next(4))
	case "uint16":
		return binary.BigEndian.Uint16(buf.Next(2))
	case "string":
		decoder := mahonia.NewDecoder("gbk")
		_, bts, _ := decoder.Translate(buf.Next(int(len)), true)
		return string(bts)
	case "[]uint8":
		//fallthrough
		bts := make([]byte, len)
		buf.Read(bts)
		return bts
	case "[]byte":
		bts := make([]byte, len)
		buf.Read(bts)
		return bts
	default:
		return nil
	}
}

// MsgGetTerminalAttrResp 查询终端属性应答
type MsgGetTerminalAttrResp struct {
	MsgGetTerminalAttrResp2011
}

func (m *MsgGetTerminalAttrResp) readBy(buf *bytes.Buffer) {
	m.TerminalAttr.readBy(buf, version2019)
}

func (m *MsgGetTerminalAttrResp) base() Input {
	return &m.MsgGetTerminalAttrResp2011
}

// MsgGetTerminalAttrResp2011 查询终端属性应答
type MsgGetTerminalAttrResp2011 struct {
	InputMark
	TerminalAttr
}

func (m *MsgGetTerminalAttrResp2011) readBy(buf *bytes.Buffer) {
	m.TerminalAttr.readBy(buf, version2011)
}

// MsgTerminalUpgradeResp 终端升级结果应答
type MsgTerminalUpgradeResp struct {
	InputMark
	// 升级类型
	Type byte
	// 升级结果
	Result byte
}

// 写入缓存中
func (m *MsgTerminalUpgradeResp) readBy(buf *bytes.Buffer) {
	// 升级类型
	m.Type, _ = buf.ReadByte()
	// 升级结果
	m.Result, _ = buf.ReadByte()
}

// MsgPositionReport 位置汇报消息
type MsgPositionReport struct {
	InputMark
	Position
}

func (m *MsgPositionReport) readBy(buf *bytes.Buffer) {
	m.Position.ReadBy(buf)
}

// MsgPosBatchReport 定位信息批量上传
type MsgPosBatchReport struct {
	InputMark
	// 位置数据类型，0-正常批量汇报，1-盲区补报
	Type byte
	// 位置信息项
	Positions []Position
}

func (m *MsgPosBatchReport) readBy(buf *bytes.Buffer) {
	// 定位数据项个数
	count := binary.BigEndian.Uint16(buf.Next(2))
	// 汇报类型，0-批量汇报，1-盲区补报
	buf.WriteByte(m.Type)
	m.Type, _ = buf.ReadByte()
	// 位置数据
	var len uint16
	var position Position
	for i := uint16(0); i < count; i++ {
		// 位置汇报数据体长度
		len = binary.BigEndian.Uint16(buf.Next(2))
		// 位置汇报数据体
		position.ReadBy(bytes.NewBuffer(buf.Next(int(len))))

		m.Positions = append(m.Positions, position)
	}
}

// MsgPositionResp 位置查询应答消息
type MsgPositionResp struct {
	InputMark
	// 应答消息流水号
	ReqNum uint16
	// 位置信息
	Position Position
}

func (m *MsgPositionResp) readBy(buf *bytes.Buffer) {
	m.ReqNum = binary.BigEndian.Uint16(buf.Next(2))
	m.Position.ReadBy(buf)
}

// MsgBDLocCheck 北斗验真上报
type MsgBDLocCheck struct {
	InputMark
	// 基础位置信息
	Position
	// 卫星状态信息
	GnnsStatus
}

func (m *MsgBDLocCheck) readBy(buf *bytes.Buffer) {
	m.Position.ReadBaseBy(buf)
	m.GnnsStatus.ReadBy(buf)
}

// MsgVehicleControlResp 车辆控制应答
type MsgVehicleControlResp struct {
	InputMark
	// 应答消息流水号
	ReqNum uint16
	// 位置信息
	Position Position
}

func (m *MsgVehicleControlResp) readBy(buf *bytes.Buffer) {
	m.ReqNum = binary.BigEndian.Uint16(buf.Next(2))
	m.Position.ReadBy(buf)
}

// MsgGetAreaOrPathResp 查询区域或路线数据应答
type MsgGetAreaOrPathResp struct {
	MsgGetAreaOrPathResp2011
}

func (m *MsgGetAreaOrPathResp) readBy(buf *bytes.Buffer) {
	// 查询类型
	m.Type, _ = buf.ReadByte()
	// 查询数量
	count := binary.BigEndian.Uint32(buf.Next(4))
	// 查询数据列表
	var area Area
	decoder := mahonia.NewDecoder("gbk")
	for i := uint32(0); i < count; i++ {
		switch m.Type {
		case ShapeTypeCircle:
			area = &RectArea{}
		case ShapeTypeRect:
			area = &RectArea{}
		case ShapeTypePolygon:
			area = &PolygonArea{}
		case ShapeTypeRoad:
			area = &Polyline{}
		}
		area.readBy(buf, decoder, version2019)

		m.Areas = append(m.Areas, area)
	}
}

func (m *MsgGetAreaOrPathResp) base() Input {
	return &m.MsgGetAreaOrPathResp2011
}

// MsgGetAreaOrPathResp2011 查询区域或路线数据应答
type MsgGetAreaOrPathResp2011 struct {
	InputMark
	// 查询类型
	Type ShapeType
	// 查询数据列表
	Areas []Area
}

func (m *MsgGetAreaOrPathResp2011) readBy(buf *bytes.Buffer) {
	// 查询类型
	m.Type, _ = buf.ReadByte()
	// 查询数量
	count := binary.BigEndian.Uint32(buf.Next(4))
	// 查询数据列表
	var area Area
	decoder := mahonia.NewDecoder("gbk")
	for i := uint32(0); i < count; i++ {
		switch m.Type {
		case ShapeTypeCircle:
			area = &RectArea{}
		case ShapeTypeRect:
			area = &RectArea{}
		case ShapeTypePolygon:
			area = &PolygonArea{}
		case ShapeTypeRoad:
			area = &Polyline{}
		}
		area.readBy(buf, decoder, version2011)

		m.Areas = append(m.Areas, area)
	}
}

// MsgDrivingRecordReport 行驶记录数据上传
type MsgDrivingRecordReport struct {
	InputMark
	// 应答流水号
	ReqNum uint16
	// 命令字
	CMD byte
	// 数据块
	Data []byte
}

func (m *MsgDrivingRecordReport) readBy(buf *bytes.Buffer) {
	// 应答流水号
	m.ReqNum = binary.BigEndian.Uint16(buf.Next(2))
	// 命令字
	m.CMD, _ = buf.ReadByte()
	// 数据块
	m.Data = make([]byte, buf.Len())
	buf.Read(m.Data)
}

// MsgWaybillReport 电子运单上报
type MsgWaybillReport struct {
	InputMark
	// 电子运单数据包
	Packet []byte
}

func (m *MsgWaybillReport) readBy(buf *bytes.Buffer) {
	// 电子运单长度
	len := binary.BigEndian.Uint32(buf.Next(4))
	// 电子运单内容
	m.Packet = make([]byte, len)
	buf.Read(m.Packet)
}

// MsgICCardReport 驾驶员身份信息上报
type MsgICCardReport struct {
	MsgICCardReport2011
	// 驾驶员身份证号
	IDCard string
}

// 写入缓存中
func (m *MsgICCardReport) readBy(buf *bytes.Buffer) {
	// 插卡状态
	m.Operation, _ = buf.ReadByte()
	// 时间
	m.Time, _ = time.ParseInLocation("060102150405", string(util.ParseBCD(buf.Next(6))), time.FixedZone("CST", 28800))
	// 插卡时
	if byte(1) == m.Operation {
		// IC卡读取结果
		m.Result, _ = buf.ReadByte()
		// 插卡读取成功
		if byte(0) == m.Result {
			decoder := mahonia.NewDecoder("gbk")
			// 驾驶员姓名
			len, _ := buf.ReadByte()
			m.Name = decoder.ConvertString(string(buf.Next(int(len))))
			// 从业资格证编码
			m.Credential = decoder.ConvertString(string(buf.Next(20)))
			// 发证机构名称
			len, _ = buf.ReadByte()
			m.Agency = decoder.ConvertString(string(buf.Next(int(len))))
			// 证件有效期
			m.Expire, _ = time.ParseInLocation("20060102", string(util.ParseBCD(buf.Next(4))), time.FixedZone("CST", 28800))
			// 驾驶员身份证号
			m.IDCard = decoder.ConvertString(string(buf.Next(20)))
		}
	}
}

func (m *MsgICCardReport) base() Input {
	return &m.MsgICCardReport2011
}

// MsgICCardReport2011 驾驶员身份信息上报
type MsgICCardReport2011 struct {
	InputMark
	// 操作，1-插卡，2-拔卡
	Operation byte
	// 时间
	Time time.Time
	// IC卡读取结果
	Result byte
	// 驾驶员姓名
	Name string
	// 从业资格证编码
	Credential string
	// 发证机构
	Agency string
	// 证件有效期
	Expire time.Time
}

// 写入缓存中
func (m *MsgICCardReport2011) readBy(buf *bytes.Buffer) {
	// 插卡状态
	m.Operation, _ = buf.ReadByte()
	// 时间
	m.Time, _ = time.ParseInLocation("060102150405", string(util.ParseBCD(buf.Next(6))), time.FixedZone("CST", 28800))
	// 插卡时
	if byte(1) == m.Operation {
		// IC卡读取结果
		m.Result, _ = buf.ReadByte()
		// 插卡读取成功
		if byte(0) == m.Result {
			decoder := mahonia.NewDecoder("gbk")
			// 驾驶员姓名
			len, _ := buf.ReadByte()
			m.Name = decoder.ConvertString(string(buf.Next(int(len))))
			// 从业资格证编码
			m.Credential = decoder.ConvertString(string(buf.Next(20)))
			// 发证机构名称
			len, _ = buf.ReadByte()
			m.Agency = decoder.ConvertString(string(buf.Next(int(len))))
			// 证件有效期
			m.Expire, _ = time.ParseInLocation("20060102", string(util.ParseBCD(buf.Next(4))), time.FixedZone("CST", 28800))
		}
	}
}

// MsgCANDataReport CAN总线数据上传
type MsgCANDataReport struct {
	InputMark
	// CAN总线数据接收时间
	Time time.Time
	// CAN数据项列表
	Datas []CanData
}

func (m *MsgCANDataReport) readBy(buf *bytes.Buffer) {
	// 数据项个数
	count := binary.BigEndian.Uint16(buf.Next(2))
	// 数据接收时间
	m.Time, _ = time.ParseInLocation("150405000", string(util.ParseBCD(buf.Next(5))), time.FixedZone("CST", 28800))
	// CAN总线数据项
	var data CanData
	for i := uint16(0); i < count; i++ {
		data.ID = binary.BigEndian.Uint32(buf.Next(4))
		buf.Read(data.Data[:])
		m.Datas = append(m.Datas, data)
	}
}

// MsgMultimediaEventReport 多媒体事件信息上传
type MsgMultimediaEventReport struct {
	InputMark
	// 多媒体数据ID
	DataID uint32
	// 多媒体类型,0-图像，1-音频，2-视频
	MimeType byte
	// 多媒体格式编码，0-JPEG，1-TIF，2-MP3，3-WAV，4-WMV
	MediaFmt byte
	// 事件项编码，0-平台下发指令，1-定时动作，2-抢劫报警触发，3-碰撞侧翻报警触发，
	// 4-门开拍照，5-门关拍照，6-车门由开变关，车速从小于20km/h到超过20km/h，7-定距拍照
	EventCode byte
	// 通道ID
	ChannelID byte
}

func (m *MsgMultimediaEventReport) readBy(buf *bytes.Buffer) {
	// 多媒体数据ID
	m.DataID = binary.BigEndian.Uint32(buf.Next(4))
	// 多媒体类型
	m.MimeType, _ = buf.ReadByte()
	// 多媒体格式编码
	m.MediaFmt, _ = buf.ReadByte()
	// 事件项编码
	m.EventCode, _ = buf.ReadByte()
	// 通道ID
	m.ChannelID, _ = buf.ReadByte()
}

// MsgMultimediaDataReport 多媒体数据上传
type MsgMultimediaDataReport struct {
	InputMark
	// 多媒体ID
	DataID uint32
	// 多媒体类型,0-图像，1-音频，2-视频
	MimeType byte
	// 多媒体格式编码，0-JPEG，1-TIF，2-MP3，3-WAV，4-WMV
	MediaFmt byte
	// 事件项编码，0-平台下发指令，1-定时动作，2-抢劫报警触发，3-碰撞侧翻报警触发，
	// 4-打开车门，5-关闭车门
	EventCode byte
	// 通道ID
	ChannelID byte
	// 位置信息
	Pos Position
	// 多媒体数据包
	Media []byte
}

func (m *MsgMultimediaDataReport) readBy(buf *bytes.Buffer) {
	// 多媒体ID
	m.DataID = binary.BigEndian.Uint32(buf.Next(4))
	// 多媒体类型
	m.MimeType, _ = buf.ReadByte()
	// 多媒体格式编码
	m.MediaFmt, _ = buf.ReadByte()
	// 事件项编码
	m.EventCode, _ = buf.ReadByte()
	// 通道ID
	m.ChannelID, _ = buf.ReadByte()
	// 位置信息
	m.Pos.ReadBaseBy(buf)
	// 多媒体数据包
	m.Media = make([]byte, buf.Len())
	buf.Read(m.Media)
}

// MsgSnapshootResp 摄像头立即拍摄命令应答
type MsgSnapshootResp struct {
	InputMark
	// 应答流水号
	ReqNum uint16
	// 结果，0-成功，1-失败，2-通道不支持
	Result byte
	// 拍摄成功的多媒体个数
	MediaIDs []uint32
}

func (m *MsgSnapshootResp) readBy(buf *bytes.Buffer) {
	// 应答流水号
	m.ReqNum = binary.BigEndian.Uint16(buf.Next(2))
	// 结果
	m.Result, _ = buf.ReadByte()
	// 多媒体ID列表
	if m.Result == 0 {
		// ID个数
		count := binary.BigEndian.Uint16(buf.Next(2))
		for i := uint16(0); i < count; i++ {
			m.MediaIDs = append(m.MediaIDs, binary.BigEndian.Uint32(buf.Next(4)))
		}
	}
}

// MsgSearchLocalMultimediaResp 存储多媒体数据检索应答
type MsgSearchLocalMultimediaResp struct {
	MsgSearchLocalMultimediaResp2011
}

func (m *MsgSearchLocalMultimediaResp) readBy(buf *bytes.Buffer) {
	// 应答流水号
	m.ReqNum = binary.BigEndian.Uint16(buf.Next(2))
	// 多媒体项总数
	count := binary.BigEndian.Uint16(buf.Next(2))
	// 多媒体项
	var multimedia Multimedia
	for i := uint16(0); i < count; i++ {
		multimedia.readBy(buf, version2019)
		m.Multimedias = append(m.Multimedias, multimedia)
	}
}

// MsgSearchLocalMultimediaResp2011 存储多媒体数据检索应答
type MsgSearchLocalMultimediaResp2011 struct {
	InputMark
	// 应答流水号
	ReqNum uint16
	// 多媒体项列表
	Multimedias []Multimedia
}

func (m *MsgSearchLocalMultimediaResp) base() Input {
	return &m.MsgSearchLocalMultimediaResp2011
}

// 写入缓存中
func (m *MsgSearchLocalMultimediaResp2011) readBy(buf *bytes.Buffer) {
	// 应答流水号
	m.ReqNum = binary.BigEndian.Uint16(buf.Next(2))
	// 多媒体项总数
	count := binary.BigEndian.Uint16(buf.Next(2))
	// 多媒体项
	var multimedia Multimedia
	for i := uint16(0); i < count; i++ {
		multimedia.readBy(buf, version2011)
		m.Multimedias = append(m.Multimedias, multimedia)
	}
}

// MsgDataUplink 数据上行透传
type MsgDataUplink struct {
	InputMark
	// 透传消息类型，0-GNSS详细数据，0x0B-IC卡信息，0x41-串口1数据，0x42-串口2数据，≥0xF0-自定义
	DataType byte
	// 透传消息内容
	Content []byte
}

func (m *MsgDataUplink) readBy(buf *bytes.Buffer) {
	m.DataType, _ = buf.ReadByte()
	m.Content = make([]byte, buf.Len())
	buf.Read(m.Content)
}

// MsgDataCompress 数据压缩上报
type MsgDataCompress struct {
	InputMark
	// 压缩后的消息体
	Body []byte
}

func (m *MsgDataCompress) readBy(buf *bytes.Buffer) {
	// 压缩消息长度
	len := binary.BigEndian.Uint32(buf.Next(4))
	// 压缩后的消息体
	m.Body = make([]byte, len)
	buf.Read(m.Body)
}

// MsgTerRSAPublicKey 终端RSA公钥
type MsgTerRSAPublicKey struct {
	InputMark
	// 密钥中的e
	E uint32
	// 密钥中的n
	N [128]byte
}

func (m *MsgTerRSAPublicKey) readBy(buf *bytes.Buffer) {
	m.E = binary.BigEndian.Uint32(buf.Next(4))
	buf.Read(m.N[:])
}

// MsgDriverFaceReport 驾驶员人脸信息采集上报
type MsgDriverFaceReport struct {
	InputMark
	// 人脸采集时间
	Time time.Time
	// 驾驶员身份识别码
	IDUID string
	// 位置基础信息
	Pos Position
	// 多媒体类型，0-图像，1-音频，2-视频
	MediaType byte
	// 多媒体格式编码，0-JPEG,1-TIF,2-MP3,3-WAV,4-WMV
	MediaFmt byte
	// 多媒体数据包
	Media []byte
}

func (m *MsgDriverFaceReport) readBy(buf *bytes.Buffer) {
	// 人脸采集时间
	m.Time, _ = time.ParseInLocation("060102150405", string(util.ParseBCD(buf.Next(6))), time.FixedZone("CST", 28800))
	// 驾驶员身份识别码
	len, _ := buf.ReadByte()
	m.IDUID = string(buf.Next(int(len)))
	// 位置基础信息
	m.Pos.ReadBaseBy(buf)
	// 多媒体类型
	m.MediaType, _ = buf.ReadByte()
	// 多媒体格式
	m.MediaFmt, _ = buf.ReadByte()
	// 多媒体数据包
	m.Media = make([]byte, buf.Len())
	buf.Read(m.Media)
}

// MsgMediaResourceList 终端上传音视频资源列表
type MsgMediaResourceList struct {
	InputMark
	// 应答流水号
	ReqNum uint16
	// 音视频资源列表
	List []MediaResource
}

func (m *MsgMediaResourceList) readBy(buf *bytes.Buffer) {
	// 应答流水号
	m.ReqNum = binary.BigEndian.Uint16(buf.Next(2))
	// 总资源数
	count := binary.BigEndian.Uint32(buf.Next(4))
	// 资源列表
	var src MediaResource
	for i := uint32(0); i < count; i++ {
		src.readBy(buf)
		m.List = append(m.List, src)
	}
}

// MsgFileUploadFinish 文件上传完成通知
type MsgFileUploadFinish struct {
	InputMark
	// 应答流水号
	ReqNum uint16
	// 结果
	Result byte
}

func (m *MsgFileUploadFinish) readBy(buf *bytes.Buffer) {
	// 应答流水号
	m.ReqNum = binary.BigEndian.Uint16(buf.Next(2))
	// 结果
	m.Result, _ = buf.ReadByte()
}

// MsgMediaPropertyReply 终端上传音视频属性
type MsgMediaPropertyReply struct {
	InputMark
	// 音频编码方式
	AudioEncodeMode byte
	// 输入音频声道数
	Channels byte
	// 输入音频采样率
	SamplingRate byte
	// 输入音频采样率位数
	SamplingRateBits byte
	//音频帧长度
	AudioFrameLen uint16
	//是否支持音频输出
	SupportAudioOutput byte
	// 视频编码方式
	VideoEncodeMode byte
	//终端支持的最大音频物理通道数
	MaxAudioPhysicalChannels byte
	//终端支持的最大视频物理通道数
	MaxVideoPhysicalChannels byte
}

func (m *MsgMediaPropertyReply) readBy(buf *bytes.Buffer) {
	m.AudioEncodeMode, _ = buf.ReadByte()
	m.Channels, _ = buf.ReadByte()
	m.SamplingRate, _ = buf.ReadByte()
	m.SamplingRateBits, _ = buf.ReadByte()
	m.AudioFrameLen = binary.BigEndian.Uint16(buf.Next(2))
	m.SupportAudioOutput, _ = buf.ReadByte()
	m.VideoEncodeMode, _ = buf.ReadByte()
	m.MaxAudioPhysicalChannels, _ = buf.ReadByte()
	m.MaxVideoPhysicalChannels, _ = buf.ReadByte()
}

//=====================================发送的消息类=====================================//

// MsgServerResponse 平台通用应答
type MsgServerResponse struct {
	OutputMark
	// 应答流水号
	ReqNum uint16
	// 应答消息ID
	ReqID uint16
	// 结果，0-成功/确认，1-失败，2-消息有误，3-不支持，4-报警处理确认
	Result byte
}

func (m *MsgServerResponse) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 2)

	// 应答流水号
	binary.BigEndian.PutUint16(value, m.ReqNum)
	buf.Write(value)
	// 应答消息ID
	binary.BigEndian.PutUint16(value, m.ReqID)
	buf.Write(value)
	// 结果
	buf.WriteByte(m.Result)
}

// NewMsgServerResponse 新建平台通用应答
func NewMsgServerResponse(reqNum, reqId uint16, result byte) *MsgServerResponse {
	return &MsgServerResponse{
		OutputMark: OutputMark{
			ID: msgIDServerResponse,
		},
		ReqNum: reqNum,
		ReqID:  reqId,
		Result: result,
	}
}

// MsgServerTimeResp 查询服务器时间应答
type MsgServerTimeResp struct {
	OutputMark
	// 服务器时间
	Time time.Time
}

func (m *MsgServerTimeResp) writeTo(buf *bytes.Buffer) {
	buf.Write(util.ToBCD([]byte(m.Time.Format("060102150405"))))
}

// NewMsgServerTimeResp 新建查询服务器时间应答
func NewMsgServerTimeResp() *MsgServerTimeResp {
	return &MsgServerTimeResp{
		OutputMark: OutputMark{
			ID: msgIDServerTimeResp,
		},
	}
}

// MsgSerGetSubpacket 服务器补传分包请求
type MsgSerGetSubpacket struct {
	MsgSerGetSubpacket2011
}

func (m *MsgSerGetSubpacket) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 2)

	// 原始消息流水号
	binary.BigEndian.PutUint16(value, m.ReqNum)
	buf.Write(value)
	// 重传包总数
	binary.BigEndian.PutUint16(value, uint16(len(m.IDs)))
	buf.Write(value)
	// 重传包ID列表
	for _, id := range m.IDs {
		binary.BigEndian.PutUint16(value, id)
		buf.Write(value)
	}
}

func (m *MsgSerGetSubpacket) base() Output {
	return &m.MsgSerGetSubpacket2011
}

// NewMsgSerGetSubpacket 新建服务器补传分包请求
func NewMsgSerGetSubpacket() *MsgSerGetSubpacket {
	return &MsgSerGetSubpacket{
		MsgSerGetSubpacket2011: MsgSerGetSubpacket2011{
			OutputMark: OutputMark{
				ID: msgIDServerPackResend,
			},
		},
	}
}

// MsgSerGetSubpacket2011 服务器补传分包请求
type MsgSerGetSubpacket2011 struct {
	OutputMark
	// 原始消息流水号
	ReqNum uint16
	// 重传包ID列表
	IDs []uint16
}

func (m *MsgSerGetSubpacket2011) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 2)

	// 原始消息流水号
	binary.BigEndian.PutUint16(value, m.ReqNum)
	buf.Write(value)
	// 重传包总数
	buf.WriteByte(byte(len(m.IDs)))
	// 重传包ID列表
	for _, id := range m.IDs {
		binary.BigEndian.PutUint16(value, id)
		buf.Write(value)
	}
}

// // NewMsgSerGetSubpacket2011 新建服务器补传分包请求（2011版）
// func NewMsgSerGetSubpacket2011() *MsgSerGetSubpacket2011 {
// 	return &MsgSerGetSubpacket2011{
// 		OutputMark: OutputMark{
// 			ID: msgIDServerPackResend,
// 		},
// 	}
// }

// MsgTerminalLoginResp 终端注册应答
type MsgTerminalLoginResp struct {
	OutputMark
	// 应答流水号
	ReqNum uint16
	// 应答结果
	Result byte
	// 鉴权码
	Token string
}

func (m *MsgTerminalLoginResp) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 2)

	// 应答流水号
	binary.BigEndian.PutUint16(value, m.ReqNum)
	buf.Write(value)
	// 应答结果
	buf.WriteByte(m.Result)
	// 鉴权码
	if byte(0) == m.Result {
		buf.Write([]byte(m.Token))
	}
}

// NewMsgTerminalLoginResp 新建终端注册应答
func NewMsgTerminalLoginResp() *MsgTerminalLoginResp {
	return &MsgTerminalLoginResp{
		OutputMark: OutputMark{
			ID: msgIDServerPackResend,
		},
	}
}

// MsgTerParamsSettings 设置终端参数
type MsgTerParamsSettings struct {
	OutputMark
	// 参数项列表
	Params []Param
}

func (m *MsgTerParamsSettings) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 4)

	// 参数总数
	buf.WriteByte(byte(len(m.Params)))
	// 参数项列表
	for _, param := range m.Params {
		binary.BigEndian.PutUint32(value, param.ID)
		buf.Write(value)
		writeParamValue(buf, param.Value)
	}
}

// writeParamValue 写入参数值
func writeParamValue(buf *bytes.Buffer, value interface{}) {
	switch reflect.ValueOf(value).Kind() {
	case reflect.Uint8:
		buf.WriteByte(byte(1))
		buf.WriteByte(value.(byte))
	case reflect.Uint16:
		buf.WriteByte(byte(2))
		bts := make([]byte, 2)
		binary.BigEndian.PutUint16(bts, value.(uint16))
		buf.Write(bts)
	case reflect.Uint32:
		buf.WriteByte(byte(4))
		bts := make([]byte, 4)
		binary.BigEndian.PutUint32(bts, value.(uint32))
		buf.Write(bts)
	case reflect.Array:
		log.Println("终端参数类型不可为数组")
	case reflect.Slice:
		if reflect.TypeOf(value).Elem().Kind() != reflect.Uint8 {
			log.Println("未知切片类型的终端参数")
			break
		}
		bts := value.([]byte)
		buf.WriteByte(byte(len(bts)))
		buf.Write(bts)
	case reflect.String:
		encoder := mahonia.NewEncoder("gbk")
		bts := []byte(encoder.ConvertString(value.(string)))
		buf.WriteByte(byte(len(bts)))
		buf.Write(bts)
	default:
		log.Println("类型未知的终端参数", reflect.ValueOf(value).Kind())
	}
}

// NewMsgTerParamsSettings 新建终端参数设置消息
func NewMsgTerParamsSettings() *MsgTerParamsSettings {
	return &MsgTerParamsSettings{
		OutputMark: OutputMark{
			ID: msgIDSetTerminalParams,
		},
	}
}

// MsgGetTerminalParams 查询终端参数
type MsgGetTerminalParams struct {
	OutputMark
}

// 从缓存中读
func (m *MsgGetTerminalParams) writeTo(buf *bytes.Buffer) {

}

// NewMsgGetTerminalParams 新建查询终端参数消息
func NewMsgGetTerminalParams() *MsgGetTerminalParams {
	return &MsgGetTerminalParams{
		OutputMark: OutputMark{
			ID: msgIDGetTerminalParams,
		},
	}
}

// MsgGetTerSpecParams 查询指定终端参数
type MsgGetTerSpecParams struct {
	OutputMark
	// 参数id列表
	ParamIDs []uint32
}

// 从缓存中读
func (m *MsgGetTerSpecParams) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 4)

	// 参数总数
	buf.WriteByte(byte(len(m.ParamIDs)))
	// 参数id列表
	for _, id := range m.ParamIDs {
		binary.BigEndian.PutUint32(value, id)
		buf.Write(value)
	}
}

// NewMsgGetTerSpecParams 新建查询指定终端参数消息
func NewMsgGetTerSpecParams() *MsgGetTerSpecParams {
	return &MsgGetTerSpecParams{
		OutputMark: OutputMark{
			ID: msgIDGetTerminalSpecParams,
		},
	}
}

// MsgTerminalControl 终端控制
type MsgTerminalControl struct {
	OutputMark
	// 命令字
	Cmd byte
	// 命令参数
	Value string
}

func (m *MsgTerminalControl) writeTo(buf *bytes.Buffer) {
	// 命令字
	buf.WriteByte(m.Cmd)
	// 命令参数
	if byte(2) == m.Cmd {
		encoder := mahonia.NewEncoder("gbk")
		buf.Write([]byte(encoder.ConvertString(m.Value)))
	}
}

// NewMsgTerminalControl 新建终端控制消息
func NewMsgTerminalControl() *MsgTerminalControl {
	return &MsgTerminalControl{
		OutputMark: OutputMark{
			ID: msgIDTerminalControl,
		},
	}
}

// MsgGetTerminalAttr 查询终端属性
type MsgGetTerminalAttr struct {
	OutputMark
}

func (m *MsgGetTerminalAttr) writeTo(buf *bytes.Buffer) {

}

// NewMsgGetTerminalAttr 新建查询终端属性消息
func NewMsgGetTerminalAttr() *MsgGetTerminalAttr {
	return &MsgGetTerminalAttr{
		OutputMark: OutputMark{
			ID: msgIDGetTerminalAttr,
		},
	}
}

// MsgTerminalUpgrade 终端升级
type MsgTerminalUpgrade struct {
	OutputMark
	// 升级类型
	Type byte
	// 制造商ID
	VendorID string
	// 终端固件版本号
	Version string
	// 升级数据包
	Pkg []byte
}

func (m *MsgTerminalUpgrade) writeTo(buf *bytes.Buffer) {
	// 升级类型
	buf.WriteByte(m.Type)
	// 制造商ID
	bts := []byte(m.VendorID)
	if l := len(bts); l < 11 {
		buf.Write(bts)
		buf.Write(make([]byte, 11-l))
	} else {
		buf.Write(bts[:11])
	}
	// 终端固件版本号长度
	bts = []byte(m.Version)
	buf.WriteByte(byte(len(bts)))
	// 终端固件版本号
	buf.Write(bts)
	// 升级数据包长度
	value := make([]byte, 4)
	binary.BigEndian.PutUint32(value, uint32(len(m.Pkg)))
	// 升级数据包
	buf.Write(m.Pkg)
}

// NewMsgTerminalUpgrade 新建终端升级消息
func NewMsgTerminalUpgrade() *MsgTerminalUpgrade {
	return &MsgTerminalUpgrade{
		OutputMark: OutputMark{
			ID: msgIDTerminalUpgrade,
		},
	}
}

// MsgGetPosition 位置信息查询
type MsgGetPosition struct {
	OutputMark
}

func (m *MsgGetPosition) writeTo(buf *bytes.Buffer) {

}

// NewMsgGetPosition 新建位置信息查询消息
func NewMsgGetPosition() *MsgGetPosition {
	return &MsgGetPosition{
		OutputMark: OutputMark{
			ID: msgIDGetPosition,
		},
	}
}

// MsgTrackControl 临时位置跟踪控制
type MsgTrackControl struct {
	OutputMark
	// 时间间隔
	Duration uint16
	// 位置跟踪有效期
	Expire uint32
}

func (m *MsgTrackControl) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 4)

	// 时间间隔
	binary.BigEndian.PutUint16(value, m.Duration)
	buf.Write(value[:2])
	// 位置跟踪有效期
	if m.Duration > 0 {
		binary.BigEndian.PutUint32(value, m.Expire)
		buf.Write(value)
	}
}

// NewMsgTrackControl 新建临时位置跟踪控制消息
func NewMsgTrackControl() *MsgTrackControl {
	return &MsgTrackControl{
		OutputMark: OutputMark{
			ID: msgIDTrackControl,
		},
	}
}

// MsgManualConfirmAlarm 人工确认报警
type MsgManualConfirmAlarm struct {
	OutputMark
	// 报警消息流水号
	AlarmNum uint16
	// 人工确认报警类型
	AlarmType AlarmConfirm
}

func (m *MsgManualConfirmAlarm) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 4)

	// 报警消息流水号
	binary.BigEndian.PutUint16(value, m.AlarmNum)
	buf.Write(value[:2])
	// 人工确认报警类型
	binary.BigEndian.PutUint32(value, uint32(m.AlarmType))
	buf.Write(value)
}

// NewMsgManualConfirmAlarm 新建人工确认报警消息
func NewMsgManualConfirmAlarm() *MsgManualConfirmAlarm {
	return &MsgManualConfirmAlarm{
		OutputMark: OutputMark{
			ID: msgIDManualConfirmAlarm,
		},
	}
}

// MsgLinkCheck 链路检测
type MsgLinkCheck struct {
	OutputMark
}

func (m *MsgLinkCheck) writeTo(buf *bytes.Buffer) {

}

// NewMsgLinkCheck 新建链路检测消息
func NewMsgLinkCheck() *MsgLinkCheck {
	return &MsgLinkCheck{
		OutputMark: OutputMark{
			ID: msgIDLinkCheck,
		},
	}
}

// MsgTextIssued 文本信息下发
type MsgTextIssued struct {
	MsgTextIssued2011
}

// 从缓存中读
func (m *MsgTextIssued) writeTo(buf *bytes.Buffer) {
	// 标志
	buf.WriteByte(byte(m.Flag))
	// 类型
	buf.WriteByte(m.Type)
	// 文本
	encoder := mahonia.NewEncoder("gbk")
	buf.Write([]byte(encoder.ConvertString(m.Text)))
}

func (m *MsgTextIssued) base() Output {
	return &m.MsgTextIssued2011
}

// NewMsgTextIssued 新建文本信息下发消息
func NewMsgTextIssued() *MsgTextIssued {
	return &MsgTextIssued{
		MsgTextIssued2011: MsgTextIssued2011{
			OutputMark: OutputMark{
				ID: msgIDTextIssued,
			},
		},
	}
}

// MsgTextIssued2011 文本信息下发
type MsgTextIssued2011 struct {
	OutputMark
	// 标志
	Flag TextFlag
	// 类型，1-通知，2-服务
	Type byte
	// 文本
	Text string
}

// 从缓存中读
func (m *MsgTextIssued2011) writeTo(buf *bytes.Buffer) {
	// 标志
	buf.WriteByte(byte(m.Flag))
	// 文本
	encoder := mahonia.NewEncoder("gbk")
	buf.Write([]byte(encoder.ConvertString(m.Text)))
}

// // NewMsgTextIssued2011 新建文本信息下发消息
// func NewMsgTextIssued2011() *MsgTextIssued2011 {
// 	return &MsgTextIssued2011{
// 		OutputMark: OutputMark{
// 			ID: msgIDTextIssued,
// 		},
// 	}
// }

// MsgTELCallback 电话回拨
type MsgTELCallback struct {
	OutputMark
	// 标志，0-普通通话，1-监听
	Flag byte
	// 电话号码
	Phone string
}

// 从缓存中读
func (m *MsgTELCallback) writeTo(buf *bytes.Buffer) {
	// 标志
	buf.WriteByte(m.Flag)
	// 电话号码
	buf.Write([]byte(m.Phone))
}

// NewMsgTELCallback 新建电话回拨消息
func NewMsgTELCallback() *MsgTELCallback {
	return &MsgTELCallback{
		OutputMark: OutputMark{
			ID: msgIDTELCallback,
		},
	}
}

// MsgContactsSettings 设置通讯录
type MsgContactsSettings struct {
	OutputMark
	// 操作类型，0-删除所有，1-删除所有并追加，2-追加，3-修改
	Operation byte
	// 联系人列表
	Contacts []Contact
}

// 从缓存中读
func (m *MsgContactsSettings) writeTo(buf *bytes.Buffer) {
	// 标志
	buf.WriteByte(m.Operation)
	// 联系人总数
	if m.Operation != 0 {
		buf.WriteByte(byte(len(m.Contacts)))

		encoder := mahonia.NewEncoder("gbk")

		// 联系人项
		for _, contact := range m.Contacts {
			// 标志
			buf.WriteByte(contact.Flag)
			// 电话号码长度
			buf.WriteByte(byte(len(contact.Phone)))
			// 电话号码
			buf.Write([]byte(contact.Phone))
			// 联系人长度
			buf.WriteByte(byte(len(contact.Name)))
			// 联系人
			buf.Write([]byte(encoder.ConvertString(contact.Name)))
		}
	}
}

// NewMsgContactsSettings 新建设置通讯录消息
func NewMsgContactsSettings() *MsgContactsSettings {
	return &MsgContactsSettings{
		OutputMark: OutputMark{
			ID: msgIDSetContacts,
		},
	}
}

// MsgVehicleControl 车辆控制
type MsgVehicleControl struct {
	MsgVehicleControl2011
}

func (m *MsgVehicleControl) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 2)
	// 控制类型数量
	binary.BigEndian.PutUint16(value, uint16(len(m.Params)))
	buf.Write(value)
	// 控制参数项列表
	for _, param := range m.Params {
		// 控制类型ID
		binary.BigEndian.PutUint16(value, param.ID)
		// 控制值
		writeVehCtrlParamValue(buf, param.Value)
	}
}

func (m *MsgVehicleControl) base() Output {
	return &m.MsgVehicleControl2011
}

// NewMsgVehicleControl 新建车辆控制消息
func NewMsgVehicleControl() *MsgVehicleControl {
	return &MsgVehicleControl{
		MsgVehicleControl2011: MsgVehicleControl2011{
			OutputMark: OutputMark{
				ID: msgIDVehicleControl,
			},
		},
	}
}

// MsgVehicleControl2011 车辆控制
type MsgVehicleControl2011 struct {
	OutputMark
	// 控制类型参数
	Params []VehCtrlParam
}

func (m *MsgVehicleControl2011) writeTo(buf *bytes.Buffer) {
	for _, param := range m.Params {
		if param.ID == 1 {
			buf.WriteByte(param.Value.(byte))
			break
		}
	}
}

// writeVehCtrlParamValue 读取参数值
func writeVehCtrlParamValue(buf *bytes.Buffer, value interface{}) {
	switch v := value.(type) {
	case byte:
		buf.WriteByte(v)
	case uint32:
		bts := make([]byte, 4)
		binary.BigEndian.PutUint32(bts, v)
		buf.Write(bts)
	case uint16:
		bts := make([]byte, 2)
		binary.BigEndian.PutUint16(bts, v)
		buf.Write(bts)
	case string:
		buf.WriteByte(byte(len(v)))
		encoder := mahonia.NewEncoder("gbk")
		buf.Write([]byte(encoder.ConvertString(v)))
	case []byte:
		buf.WriteByte(byte(len(v)))
		buf.Write(v)
	default:
		log.Println("存在无效的车辆控制参数")
	}
}

// NewMsgVehicleControl2011 新建车辆控制消息（2011版）
func NewMsgVehicleControl2011() *MsgVehicleControl2011 {
	return &MsgVehicleControl2011{
		OutputMark: OutputMark{
			ID: msgIDVehicleControl,
		},
	}
}

// MsgRoundAreaSettings 设置圆形区域消息
type MsgRoundAreaSettings struct {
	MsgRoundAreaSettings2011
}

// 从缓存中读
func (m *MsgRoundAreaSettings) writeTo(buf *bytes.Buffer) {
	encoder := mahonia.NewEncoder("gbk")

	// 设置属性
	buf.WriteByte(m.Operation)
	// 区域总数
	buf.WriteByte(byte(len(m.Areas)))
	// 区域项列表
	for _, area := range m.Areas {
		area.writeTo(buf, encoder, version2019)
	}
}

func (m *MsgRoundAreaSettings) base() Output {
	return &m.MsgRoundAreaSettings2011
}

// NewMsgRoundAreaSettings 新建设置圆形区域消息
func NewMsgRoundAreaSettings() *MsgRoundAreaSettings {
	return &MsgRoundAreaSettings{
		MsgRoundAreaSettings2011: MsgRoundAreaSettings2011{
			OutputMark: OutputMark{
				ID: msgIDSetRoundArea,
			},
		},
	}
}

// MsgRoundAreaSettings2011 设置圆形区域消息
type MsgRoundAreaSettings2011 struct {
	OutputMark
	// 操作
	Operation byte
	// 区域列表
	Areas []RoundArea
}

// 从缓存中读
func (m *MsgRoundAreaSettings2011) writeTo(buf *bytes.Buffer) {
	encoder := mahonia.NewEncoder("gbk")

	// 设置属性
	buf.WriteByte(m.Operation)
	// 区域总数
	buf.WriteByte(byte(len(m.Areas)))
	// 区域项列表
	for _, area := range m.Areas {
		area.writeTo(buf, encoder, version2011)
	}
}

// // NewMsgRoundAreaSettings2011 新建设置圆形区域消息（2011版）
// func NewMsgRoundAreaSettings2011() *MsgRoundAreaSettings2011 {
// 	return &MsgRoundAreaSettings2011{
// 		OutputMark: OutputMark{
// 			ID: msgIDSetRoundArea,
// 		},
// 	}
// }

// MsgRoundAreaDelete 删除圆形区域消息
type MsgRoundAreaDelete struct {
	OutputMark
	// 区域ID列表
	IDs []uint32
}

func (m *MsgRoundAreaDelete) writeTo(buf *bytes.Buffer) {
	// 区域数
	buf.WriteByte(byte(len(m.IDs)))

	if len(m.IDs) == 0 {
		return
	}

	value := make([]byte, 4)

	// 区域项id列表
	for _, id := range m.IDs {
		binary.BigEndian.PutUint32(value, id)
		buf.Write(value)
	}
}

// NewMsgRoundAreaDelete 新建删除圆形区域消息
func NewMsgRoundAreaDelete() *MsgRoundAreaDelete {
	return &MsgRoundAreaDelete{
		OutputMark: OutputMark{
			ID: msgIDDeleteRoundArea,
		},
	}
}

// MsgRectAreaSettings 设置矩形区域消息
type MsgRectAreaSettings struct {
	MsgRectAreaSettings2011
}

// 从缓存中读
func (m *MsgRectAreaSettings) writeTo(buf *bytes.Buffer) {
	encoder := mahonia.NewEncoder("gbk")

	// 设置属性
	buf.WriteByte(m.Operation)
	// 区域总数
	buf.WriteByte(byte(len(m.Areas)))
	// 区域项列表
	for _, area := range m.Areas {
		area.writeTo(buf, encoder, version2019)
	}
}

func (m *MsgRectAreaSettings) base() Output {
	return &m.MsgRectAreaSettings2011
}

// NewMsgRectAreaSettings 新建设置矩形区域消息
func NewMsgRectAreaSettings() *MsgRectAreaSettings {
	return &MsgRectAreaSettings{
		MsgRectAreaSettings2011: MsgRectAreaSettings2011{
			OutputMark: OutputMark{
				ID: msgIDSetRectArea,
			},
		},
	}
}

// MsgRectAreaSettings2011 设置矩形区域消息
type MsgRectAreaSettings2011 struct {
	OutputMark
	// 操作
	Operation byte
	// 区域列表
	Areas []RectArea
}

// 从缓存中读
func (m *MsgRectAreaSettings2011) writeTo(buf *bytes.Buffer) {
	encoder := mahonia.NewEncoder("gbk")

	// 设置属性
	buf.WriteByte(m.Operation)
	// 区域总数
	buf.WriteByte(byte(len(m.Areas)))
	// 区域项列表
	for _, area := range m.Areas {
		area.writeTo(buf, encoder, version2011)
	}
}

// // NewMsgRectAreaSettings2011 新建设置矩形区域消息（2011版）
// func NewMsgRectAreaSettings2011() *MsgRectAreaSettings2011 {
// 	return &MsgRectAreaSettings2011{
// 		OutputMark: OutputMark{
// 			ID: msgIDSetRectArea,
// 		},
// 	}
// }

// MsgRectAreaDelete 删除矩形区域消息
type MsgRectAreaDelete struct {
	OutputMark
	// 区域ID列表
	IDs []uint32
}

func (m *MsgRectAreaDelete) writeTo(buf *bytes.Buffer) {
	// 区域数
	buf.WriteByte(byte(len(m.IDs)))
	// 删除所有
	if len(m.IDs) == 0 {
		return
	}

	value := make([]byte, 4)

	// 区域项id列表
	for _, id := range m.IDs {
		binary.BigEndian.PutUint32(value, id)
		buf.Write(value)
	}
}

// NewMsgRectAreaDelete 新建删除矩形区域消息
func NewMsgRectAreaDelete() *MsgRectAreaDelete {
	return &MsgRectAreaDelete{
		OutputMark: OutputMark{
			ID: msgIDDeleteRectArea,
		},
	}
}

// MsgPolygonAreaSettings 设置多边形区域消息
type MsgPolygonAreaSettings struct {
	MsgPolygonAreaSettings2011
}

// 从缓存中读
func (m *MsgPolygonAreaSettings) writeTo(buf *bytes.Buffer) {
	encoder := mahonia.NewEncoder("gbk")
	m.Area.writeTo(buf, encoder, version2019)
}

func (m *MsgPolygonAreaSettings) base() Output {
	return &m.MsgPolygonAreaSettings2011
}

// NewMsgPolygonAreaSettings 新建设置多边形区域消息
func NewMsgPolygonAreaSettings() *MsgPolygonAreaSettings {
	return &MsgPolygonAreaSettings{
		MsgPolygonAreaSettings2011: MsgPolygonAreaSettings2011{
			OutputMark: OutputMark{
				ID: msgIDSetPolygonArea,
			},
		},
	}
}

// MsgPolygonAreaSettings2011 设置多边形区域消息
type MsgPolygonAreaSettings2011 struct {
	OutputMark
	// 区域
	Area PolygonArea
}

// 从缓存中读
func (m *MsgPolygonAreaSettings2011) writeTo(buf *bytes.Buffer) {
	encoder := mahonia.NewEncoder("gbk")
	m.Area.writeTo(buf, encoder, version2011)
}

// // NewMsgPolygonAreaSettings2011 新建设置多边形区域消息
// func NewMsgPolygonAreaSettings2011() *MsgPolygonAreaSettings2011 {
// 	return &MsgPolygonAreaSettings2011{
// 		OutputMark: OutputMark{
// 			ID: msgIDSetPolygonArea,
// 		},
// 	}
// }

// MsgPolygonAreaDelete 删除多边形区域消息
type MsgPolygonAreaDelete struct {
	OutputMark
	// 区域ID列表
	IDs []uint32
}

func (m *MsgPolygonAreaDelete) writeTo(buf *bytes.Buffer) {
	// 区域数
	buf.WriteByte(byte(len(m.IDs)))
	// 删除所有
	if len(m.IDs) == 0 {
		return
	}

	value := make([]byte, 4)

	// 区域项id列表
	for _, id := range m.IDs {
		binary.BigEndian.PutUint32(value, id)
		buf.Write(value)
	}
}

// NewMsgPolygonAreaDelete 新建删除多边形区域消息
func NewMsgPolygonAreaDelete() *MsgPolygonAreaDelete {
	return &MsgPolygonAreaDelete{
		OutputMark: OutputMark{
			ID: msgIDDeletePolygonArea,
		},
	}
}

// MsgPathSettings 设置路线消息
type MsgPathSettings struct {
	MsgPathSettings2011
}

// 从缓存中读
func (m *MsgPathSettings) writeTo(buf *bytes.Buffer) {
	encoder := mahonia.NewEncoder("gbk")
	m.Path.writeTo(buf, encoder, version2019)
}

func (m *MsgPathSettings) base() Output {
	return &m.MsgPathSettings2011
}

// NewMsgPathSettings 新建设置路线消息
func NewMsgPathSettings() *MsgPathSettings {
	return &MsgPathSettings{
		MsgPathSettings2011: MsgPathSettings2011{
			OutputMark: OutputMark{
				ID: msgIDSetPolyline,
			},
		},
	}
}

// MsgPathSettings2011 设置路线消息
type MsgPathSettings2011 struct {
	OutputMark
	Path Polyline
}

// 从缓存中读
func (m *MsgPathSettings2011) writeTo(buf *bytes.Buffer) {
	encoder := mahonia.NewEncoder("gbk")
	m.Path.writeTo(buf, encoder, version2019)
}

// // NewMsgPathSettings2011 新建设置路线消息
// func NewMsgPathSettings2011() *MsgPathSettings2011 {
// 	return &MsgPathSettings2011{
// 		OutputMark: OutputMark{
// 			ID: msgIDSetPolyline,
// 		},
// 	}
// }

// MsgPathDelete 删除路线消息
type MsgPathDelete struct {
	OutputMark
	// 路线id列表
	IDs []uint32
}

// 从缓存中读
func (m *MsgPathDelete) writeTo(buf *bytes.Buffer) {
	// 区域数
	buf.WriteByte(byte(len(m.IDs)))
	// 删除所有
	if len(m.IDs) == 0 {
		return
	}

	value := make([]byte, 4)

	// 区域项id列表
	for _, id := range m.IDs {
		binary.BigEndian.PutUint32(value, id)
		buf.Write(value)
	}
}

// NewMsgPathDelete 新建删除路线消息
func NewMsgPathDelete() *MsgPathDelete {
	return &MsgPathDelete{
		OutputMark: OutputMark{
			ID: msgIDDeletePolyline,
		},
	}
}

// MsgGetAreaOrPath 查询区域或线路数据
type MsgGetAreaOrPath struct {
	OutputMark
	// 查询类型
	Type ShapeType
	// 查询项id列表，为空时表示查所有
	IDs []uint32
}

func (m *MsgGetAreaOrPath) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 4)

	// 查询类型
	buf.WriteByte(byte(m.Type))
	// 查询数量
	binary.BigEndian.PutUint32(value, uint32(len(m.IDs)))
	buf.Write(value)
	// 查询项id列表
	for _, id := range m.IDs {
		binary.BigEndian.PutUint32(value, id)
		buf.Write(value)
	}
}

// NewMsgGetAreaOrPath 新建查询区域或线路数据消息
func NewMsgGetAreaOrPath() *MsgGetAreaOrPath {
	return &MsgGetAreaOrPath{
		OutputMark: OutputMark{
			ID: msgIDGetArea,
		},
	}
}

// MsgGatherDrivingRecord 行驶记录数据采集命令
type MsgGatherDrivingRecord struct {
	OutputMark
	// 命令字
	CMD byte
	// 数据块
	Data []byte
}

func (m *MsgGatherDrivingRecord) writeTo(buf *bytes.Buffer) {
	// 命令字
	buf.WriteByte(m.CMD)
	// 数据块
	buf.Write(m.Data)
}

// NewMsgGatherDrivingRecord 新建行驶记录数据采集命令消息
func NewMsgGatherDrivingRecord() *MsgGatherDrivingRecord {
	return &MsgGatherDrivingRecord{
		OutputMark: OutputMark{
			ID: msgIDGatherDrivingRecord,
		},
	}
}

// MsgDrivingRecordParamsIssued 行驶记录参数下传命令
type MsgDrivingRecordParamsIssued struct {
	OutputMark
	// 命令字
	CMD byte
	// 数据块
	Data []byte
}

func (m *MsgDrivingRecordParamsIssued) writeTo(buf *bytes.Buffer) {
	// 命令字
	buf.WriteByte(m.CMD)
	// 数据块
	buf.Write(m.Data)
}

// NewMsgDrivingRecordParamsIssued 新建行驶记录参数下传命令消息
func NewMsgDrivingRecordParamsIssued() *MsgDrivingRecordParamsIssued {
	return &MsgDrivingRecordParamsIssued{
		OutputMark: OutputMark{
			ID: msgIDDriRecordParamsIssued,
		},
	}
}

// MsgMultimediaReportResp 多媒体数据上传应答
type MsgMultimediaReportResp struct {
	OutputMark
	// 多媒体ID
	MediaID uint32
	// 重传包ID列表
	IDs []uint16
}

func (m *MsgMultimediaReportResp) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 4)

	// 多媒体ID
	binary.BigEndian.PutUint32(value, m.MediaID)
	buf.Write(value)
	// 重传包ID列表
	if len(m.IDs) > 0 {
		// 总数
		buf.WriteByte(byte(len(m.IDs)))
		// 列表
		for _, id := range m.IDs {
			binary.BigEndian.PutUint16(value, id)
			buf.Write(value[:2])
		}
	}
}

// NewMsgMultimediaReportResp 新建多媒体数据上传应答消息
func NewMsgMultimediaReportResp() *MsgMultimediaReportResp {
	return &MsgMultimediaReportResp{
		OutputMark: OutputMark{
			ID: msgIDMultimediaDataReportResp,
		},
	}
}

// MsgSnapshoot 摄像头立即拍摄命令
type MsgSnapshoot struct {
	OutputMark
	// 通道ID
	ChannelID byte
	// 拍摄命令，0-停止拍摄，0xFFFF-录像，其他-拍照张数
	Order uint16
	// 拍照间隔/录像时间，单位为s，0-表示按最小间隔拍照或一直录像
	Duration uint16
	// 保存标识，1-保存，0-实时上传
	SaveTag byte
	// 分辨率，0-最低分辨率，1-320×240，2-640×480，3-800×600，4-1024×768，
	// 5-176×144[Qcif]，6-352×288[Cif]，7-704×288[HALF D1]，8-704×576[D1]，0xFF-最高分辨率
	Resolution byte
	// 图片/视频质量，取值范围1~10，1-代表质量损失最小，10-表示压缩比最大
	Quality byte
	// 亮度，0~255
	Brightness byte
	// 对比度，0~127
	Contrast byte
	// 饱和度，0~127
	Saturation byte
	// 色度，0~255
	Chroma byte
}

func (m *MsgSnapshoot) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 2)

	// 通道ID
	buf.WriteByte(m.ChannelID)
	// 拍摄命令
	binary.BigEndian.PutUint16(value, m.Order)
	buf.Write(value)
	// 拍照间隔/录像时长
	binary.BigEndian.PutUint16(value, m.Duration)
	buf.Write(value)
	// 保存标识
	buf.WriteByte(m.SaveTag)
	// 分辨率
	buf.WriteByte(m.Resolution)
	// 图片/视频质量
	buf.WriteByte(m.Quality)
	// 亮度
	buf.WriteByte(m.Brightness)
	// 对比度
	buf.WriteByte(m.Contrast)
	// 饱和度
	buf.WriteByte(m.Saturation)
	// 色度
	buf.WriteByte(m.Chroma)
}

// NewMsgSnapshoot 新建摄像头立即拍摄命令消息
func NewMsgSnapshoot() *MsgSnapshoot {
	return &MsgSnapshoot{
		OutputMark: OutputMark{
			ID: msgIDSnapshot,
		},
	}
}

// MsgSearchLocalMultimedia 存储多媒体数据检索
type MsgSearchLocalMultimedia struct {
	OutputMark
	// 多媒体类型，0-图像，1-音频，2-视频
	MediaType byte
	// 通道ID，0表示检索该媒体类型的所有通道
	ChannelID byte
	// 事件项编码，0-平台下发指令，1-定时动作，2-抢劫报警触发，3-碰撞侧翻报警触发，其他保留
	EventCode byte
	// 起始时间
	STime time.Time
	// 结束时间
	ETime time.Time
}

func (m *MsgSearchLocalMultimedia) writeTo(buf *bytes.Buffer) {
	// 多媒体类型
	buf.WriteByte(m.MediaType)
	// 通道ID
	buf.WriteByte(m.ChannelID)
	// 事件项编码
	buf.WriteByte(m.EventCode)
	// 起始时间
	if m.STime.IsZero() {
		buf.Write(make([]byte, 6))
	} else {
		buf.Write(util.ToBCD([]byte(m.STime.Format("060102150405"))))
	}
	// 结束时间
	if m.ETime.IsZero() {
		buf.Write(make([]byte, 6))
	} else {
		buf.Write(util.ToBCD([]byte(m.ETime.Format("060102150405"))))
	}
}

// NewMsgSearchLocalMultimedia 新建存储多媒体数据检索消息
func NewMsgSearchLocalMultimedia() *MsgSearchLocalMultimedia {
	return &MsgSearchLocalMultimedia{
		OutputMark: OutputMark{
			ID: msgIDGetMultimediaSaveInfo,
		},
	}
}

// MsgPullLocalMultimedia 存储多媒体数据上传命令
type MsgPullLocalMultimedia struct {
	OutputMark
	// 多媒体类型，0-图像，1-音频，2-视频
	MediaType byte
	// 通道ID
	ChannelID byte
	// 事件项编码，0-平台下发指令，1-定时动作，2-抢劫报警触发，3-碰撞侧翻报警触发，其他保留
	EventCode byte
	// 起始时间
	STime time.Time
	// 结束时间
	ETime time.Time
	// 删除标识，0-保留，1-删除
	Deleted byte
}

func (m *MsgPullLocalMultimedia) writeTo(buf *bytes.Buffer) {
	// 多媒体类型
	buf.WriteByte(m.MediaType)
	// 通道ID
	buf.WriteByte(m.ChannelID)
	// 事件项编码
	buf.WriteByte(m.EventCode)
	// 起始时间
	if m.STime.IsZero() {
		buf.Write(make([]byte, 6))
	} else {
		buf.Write(util.ToBCD([]byte(m.STime.Format("060102150405"))))
	}
	// 结束时间
	if m.ETime.IsZero() {
		buf.Write(make([]byte, 6))
	} else {
		buf.Write(util.ToBCD([]byte(m.ETime.Format("060102150405"))))
	}
	// 删除标志
	buf.WriteByte(m.Deleted)
}

// NewMsgPullLocalMultimedia 新建存储多媒体数据上传命令消息
func NewMsgPullLocalMultimedia() *MsgPullLocalMultimedia {
	return &MsgPullLocalMultimedia{
		OutputMark: OutputMark{
			ID: msgIDGetMultimediaSaveUp,
		},
	}
}

// MsgRecording 录音开始命令
type MsgRecording struct {
	OutputMark
	// 录音命令，0-停止录音，1-开始录音
	Operation byte
	// 录音时间，0表示一直录着
	Duration uint16
	// 保存标识，0-实时上传，1-保存
	SaveFlag byte
	// 音频采样率，0-8K，1-11K，2-23K，3-32K
	SampleRate byte
}

func (m *MsgRecording) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 2)

	// 录音命令
	buf.WriteByte(m.Operation)
	// 录音时间
	binary.BigEndian.PutUint16(value, m.Duration)
	buf.Write(value)
	// 保存标识
	buf.WriteByte(m.SaveFlag)
	// 音频采样率
	buf.WriteByte(m.SampleRate)
}

// NewMsgRecording 新建录音开始命令消息
func NewMsgRecording() *MsgRecording {
	return &MsgRecording{
		OutputMark: OutputMark{
			ID: msgIDGetAudioStartRecord,
		},
	}
}

// MsgGetLocalMultimedia 单条存储多媒体数据检索上传命令
type MsgGetLocalMultimedia struct {
	OutputMark
	// 多媒体ID
	MediaID uint32
	// 删除标识，0-保留，1-删除
	Deleted byte
}

func (m *MsgGetLocalMultimedia) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 4)

	// 多媒体ID
	binary.BigEndian.PutUint32(value, m.MediaID)
	buf.Write(value)
	// 删除标志
	buf.WriteByte(m.Deleted)
}

// NewMsgGetLocalMultimedia 新建单条存储多媒体数据检索上传命令消息
func NewMsgGetLocalMultimedia() *MsgGetLocalMultimedia {
	return &MsgGetLocalMultimedia{
		OutputMark: OutputMark{
			ID: msgIDGetSingleMediaSaveUp,
		},
	}
}

// MsgDataDownlink 数据下行透传
type MsgDataDownlink struct {
	OutputMark
	// 透传消息类型，0-GNSS详细数据，0x0B-IC卡信息，0x41-串口1数据，0x42-串口2数据，≥0xF0-自定义
	DataType byte
	// 透传消息内容
	Content []byte
}

func (m *MsgDataDownlink) writeTo(buf *bytes.Buffer) {
	// 透传消息类型
	buf.WriteByte(m.DataType)
	// 透传消息内容
	buf.Write(m.Content)
}

// NewMsgDataDownlink 新建数据下行透传消息
func NewMsgDataDownlink() *MsgDataDownlink {
	return &MsgDataDownlink{
		OutputMark: OutputMark{
			ID: msgIDDataDownPenetrate,
		},
	}
}

// MsgServerRSAPublicKey 平台RSA公钥
type MsgServerRSAPublicKey struct {
	OutputMark
	// 密钥中的e
	E uint32
	// 密钥中的n
	N [128]byte
}

func (m *MsgServerRSAPublicKey) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 4)

	// 密钥中的e
	binary.BigEndian.PutUint32(value, m.E)
	buf.Write(value)
	// 密钥中的n
	buf.Write(m.N[:])
}

// NewMsgServerRSAPublicKey 新建平台RSA公钥消息
func NewMsgServerRSAPublicKey() *MsgServerRSAPublicKey {
	return &MsgServerRSAPublicKey{
		OutputMark: OutputMark{
			ID: msgIDPlatformRSAPublickey,
		},
	}
}

// MsgGetDriverIdentity 上报驾驶员身份信息请求
type MsgGetDriverIdentity struct {
	OutputMark
}

func (m *MsgGetDriverIdentity) writeTo(buf *bytes.Buffer) {

}

// NewMsgGetDriverIdentity 新建上报驾驶员身份信息请求消息
func NewMsgGetDriverIdentity() *MsgGetDriverIdentity {
	return &MsgGetDriverIdentity{
		OutputMark: OutputMark{
			ID: msgIDGetDriverIdentity,
		},
	}
}

// MsgRemoteVideoReplay 平台下发远程录像回放请求
type MsgRemoteVideoReplay struct {
	OutputMark
	// 服务器IP地址
	IP string
	// TCP端口号
	TCPPort uint16
	// UDP端口号
	UDPPort uint16
	// 逻辑通道号
	Channel byte
	// 音视频类型：0-音视频，1-音频，2-视频，3-视频或音视频
	MediaType byte
	// 码流类型：0-主码流或子码流，1-主码流，2-子码流
	StreamType byte
	// 存储器类型：0-主存储器或灾备存储器，1-主存储器，2-灾备存储器
	StorageType byte
	// 回放方式：0-正常回放，1-快进回放，2-关键帧快退回放，3-关键帧播放，4-单帧上传
	ReplayMode byte
	// 快进/快退倍数：0-无效，1-1倍，2-2倍，3-4倍，4-8倍，5-16倍
	Multiple byte
	// 开始时间
	STime time.Time
	// 结束时间
	ETime time.Time
}

func (m *MsgRemoteVideoReplay) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 2)

	// 服务器IP地址
	buf.WriteByte(byte(len(m.IP)))
	buf.Write([]byte(m.IP))
	// TCP端口号
	binary.BigEndian.PutUint16(value, m.TCPPort)
	buf.Write(value)
	// UDP端口号
	binary.BigEndian.PutUint16(value, m.UDPPort)
	buf.Write(value)
	// 逻辑通道号
	buf.WriteByte(m.Channel)
	// 音视频类型
	buf.WriteByte(m.MediaType)
	// 码流类型
	buf.WriteByte(m.StreamType)
	// 存储器类型
	buf.WriteByte(m.StorageType)
	// 回放方式
	buf.WriteByte(m.ReplayMode)
	// 快进/快退倍数
	buf.WriteByte(m.Multiple)
	// 开始时间
	buf.Write(util.ToBCD([]byte(m.STime.Format("060102150405"))))
	// 结束时间，回访方式为4时，该字段无效，可以不读
	if m.ReplayMode != 4 {
		if m.ETime.IsZero() {
			buf.Write(make([]byte, 6))
		} else {
			buf.Write(util.ToBCD([]byte(m.ETime.Format("060102150405"))))
		}
	}
}

// NewMsgRemoteVideoReplay 新建平台下发远程录像回放请求消息
func NewMsgRemoteVideoReplay() *MsgRemoteVideoReplay {
	return &MsgRemoteVideoReplay{
		OutputMark: OutputMark{
			ID: msgIDRemoteVideoReplay,
		},
	}
}

// MsgRemoteReplayControl 平台下发远程录像回放控制
type MsgRemoteReplayControl struct {
	OutputMark
	// 音视频通道号
	Channel byte
	// 回放控制，0-开始回放，1-暂停回放，2-结束回放，3-快进回放，4-关键帧快退回放，5-拖动回放，6-关键帧播放
	Operation byte
	// 快进/快退倍数，0-无效，1-1倍，2-2倍，3-4倍，4-8倍，5-16倍
	Multiple byte
	// 拖动位置
	Seek time.Time
}

func (m *MsgRemoteReplayControl) writeTo(buf *bytes.Buffer) {
	// 音视频通道号
	buf.WriteByte(m.Channel)
	// 回放控制
	buf.WriteByte(m.Operation)
	// 快进/快退倍数
	buf.WriteByte(m.Multiple)
	// 拖动位置
	if m.Operation == 5 {
		buf.Write(util.ToBCD([]byte(m.Seek.Format("060102150405"))))
	}
}

// NewMsgRemoteReplayControl 新建平台下发远程录像回放控制消息
func NewMsgRemoteReplayControl() *MsgRemoteReplayControl {
	return &MsgRemoteReplayControl{
		OutputMark: OutputMark{
			ID: msgIDRemoteReplayControl,
		},
	}
}

// MsgMediaResourceSelect 查询音视频资源
type MsgMediaResourceSelect struct {
	OutputMark
	// 逻辑通道号
	Channel byte
	// 开始时间
	STime time.Time
	// 结束时间
	ETime time.Time
	// 报警标志
	Alarm uint64
	// 音视频类型：0-音视频，1-音频，2-视频,3-视频或音视频
	MediaType byte
	// 码流类型，0-所有码流，1-主码流，2-子码流
	StreamType byte
	// 存储器类型，0-所有存储器，1-主存储器，2-灾备存储器
	StorageType byte
}

func (m *MsgMediaResourceSelect) writeTo(buf *bytes.Buffer) {
	// 逻辑通道号
	buf.WriteByte(m.Channel)
	// 开始时间
	if m.STime.IsZero() {
		buf.Write(make([]byte, 6))
	} else {
		buf.Write(util.ToBCD([]byte(m.STime.Format("060102150405"))))
	}
	// 结束时间
	if m.ETime.IsZero() {
		buf.Write(make([]byte, 6))
	} else {
		buf.Write(util.ToBCD([]byte(m.ETime.Format("060102150405"))))
	}
	// 报警标志
	value := make([]byte, 8)
	binary.BigEndian.PutUint64(value, m.Alarm)
	buf.Write(value)
	// 音视频类型
	buf.WriteByte(m.MediaType)
	// 码流类型
	buf.WriteByte(m.StreamType)
	// 存储器类型
	buf.WriteByte(m.StorageType)
}

// NewMsgMediaResourceSelect 新建查询音视频资源消息
func NewMsgMediaResourceSelect() *MsgMediaResourceSelect {
	return &MsgMediaResourceSelect{
		OutputMark: OutputMark{
			ID: msgIDMediaResourceSelect,
		},
	}
}

// // MsgAttachUploadSb 苏标附件上传消息
// type MsgAttachUploadSb struct {
// 	OutputMark
// 	ServerAddr string       // 服务器地址
// 	TcpPort    uint16       // 服务器TCP端口号
// 	UdpPort    uint16       // 服务器UDP端口号
// 	AlarmId    AlarmMarking // 报警标识
// 	AlarmCode  string       // 报警编号32byte
// }

// func (m *MsgAttachUploadSb) readBy(buf *bytes.Buffer) {
// 	//服务器地址长度
// 	buf.WriteByte(byte(len(m.ServerAddr)))
// 	//服务器地址
// 	buf.Write([]byte(m.ServerAddr))
// 	//服务器TCP端口
// 	port := make([]byte, 2)
// 	binary.BigEndian.PutUint16(port, m.TcpPort)
// 	buf.Write(port)
// 	//服务器UDP端口
// 	binary.BigEndian.PutUint16(port, m.UdpPort)
// 	buf.Write(port)
// 	// 报警标识
// 	m.AlarmId.writeTo(buf)
// 	// 报警编号
// 	if l := len(m.AlarmCode); l >= 32 {
// 		buf.Write([]byte(m.AlarmCode)[:32])
// 	} else {
// 		buf.Write([]byte(m.AlarmCode))
// 		buf.Write(make([]byte, 32-l))
// 	}
// }

// // MsgAttachUpload 附件上传消息
// type MsgAttachUpload struct {
// 	OutputMark
// 	beans.Attach
// }

// func (m *MsgAttachUpload) readBy(buf *bytes.Buffer) {
// 	//服务器地址长度
// 	serverAddrLen, _ := buf.ReadByte()
// 	//服务器地址
// 	serverAddr := make([]byte, serverAddrLen)
// 	buf.Read(serverAddr)
// 	decoder := mahonia.NewDecoder("gbk")
// 	_, cdata, _ := decoder.Translate(serverAddr, true)
// 	m.ServerAddr = string(cdata)

// 	//服务器端口
// 	port := make([]byte, 2)
// 	buf.Read(port)
// 	m.ServerPort = binary.BigEndian.Uint16(port)

// 	//用户名长度
// 	userNameLen, _ := buf.ReadByte()
// 	//用户名
// 	userName := make([]byte, userNameLen)
// 	buf.Read(userName)
// 	_, c, _ := decoder.Translate(userName, true)
// 	m.UserName = string(c)

// 	//密码长度
// 	pwdLen, _ := buf.ReadByte()
// 	//密码
// 	pwd := make([]byte, pwdLen)
// 	buf.Read(pwd)
// 	_, cdata2, _ := decoder.Translate(pwd, true)
// 	m.Password = string(cdata2)

// 	//文件上传路径长度
// 	uploadPathLen, _ := buf.ReadByte()
// 	//文件上传路径
// 	uploadPath := make([]byte, uploadPathLen)
// 	buf.Read(uploadPath)
// 	_, cdata3, _ := decoder.Translate(uploadPath, true)
// 	m.UploadPath = string(cdata3)

// 	//报警附件总数
// 	m.AttachCount, _ = buf.ReadByte()

// 	//报警附件标识
// 	AttachIdentifier := make([]byte, 34)
// 	buf.Read(AttachIdentifier)
// 	fmt.Println(string(AttachIdentifier))
// 	m.AttachIdentifier = AttachIdentifier
// 	fmt.Println(string(m.AttachIdentifier))
// 	fmt.Printf("调试，附件上传：%+v\n", m.Attach)
// }

// //MsgAttachUpload 附件上传消息
// type MsgAttachUpload struct {
// 	InputMark
// 	//beans.Attach
// 	beans.AlarmAttachUploadData
// }

// func (m *MsgAttachUpload) readBy(buf *bytes.Buffer) {
// 	// 附件服务器IP地址长度
// 	attachServerIPLen, _ := buf.ReadByte()
// 	// 附件服务器IP地址
// 	attachServerIPAddr := make([]byte, attachServerIPLen)
// 	buf.Read(attachServerIPAddr)
// 	m.AttachServerAddr = string(attachServerIPAddr)
// 	// 附件服务器端口（TCP）
// 	port := make([]byte, 2)
// 	buf.Read(port)
// 	m.AttachServerTcpPort = binary.BigEndian.Uint16(port)
// 	// 附件服务器端口（UDP）
// 	buf.Read(port)
// 	m.AttachServerUdpPort = binary.BigEndian.Uint16(port)
// 	// 终端ID
// 	terminalID := make([]byte, 7)
// 	buf.Read(terminalID)
// 	m.TerminalID = string(terminalID)
// 	// 时间
// 	timeBcd := make([]byte, 6)
// 	buf.Read(timeBcd)
// 	m.Time, _ = time.ParseInLocation(
// 		"060102150405",
// 		string(util.ParseBCD(timeBcd)),
// 		time.FixedZone("CST", 28800))
// 	// 序号
// 	m.SerialNumber, _ = buf.ReadByte()
// 	// 附件数量
// 	m.AttachCount, _ = buf.ReadByte()
// 	// 预留
// 	buf.ReadByte()
// 	// 报警编号
// 	alaramNumber := make([]byte, 32)
// 	buf.Read(alaramNumber)
// 	m.AlarmNumber = string(alaramNumber)
// 	// 预留，暂不需要解析
// }

// MsgMediaProperty 查询音视频属性
type MsgMediaProperty struct {
	OutputMark
}

func (m *MsgMediaProperty) writeTo(buf *bytes.Buffer) {
	//不需要解析
}

// NewMsgMediaProperty 新建查询音视频属性消息
func NewMsgMediaProperty() *MsgMediaProperty {
	return &MsgMediaProperty{
		OutputMark: OutputMark{
			ID: msgIDQueryMediaProperty,
		},
	}
}

// MsgRealMediaRequest 实时音视频传输请求
type MsgRealMediaRequest struct {
	OutputMark
	// 服务器IP地址
	IP string
	// TCP端口号
	TCPPort uint16
	// UDP端口号
	UDPPort uint16
	// 逻辑通道号
	Channel byte
	// 音视频类型：0-音视频，1-视频，2-双向对讲，3-监听，4-中心广播，5-透传
	MediaType byte
	// 码流类型：0-主码流，1-子码流
	StreamType byte
}

func (m *MsgRealMediaRequest) writeTo(buf *bytes.Buffer) {
	// 服务器IP地址
	buf.WriteByte(byte(len(m.IP)))
	buf.Write([]byte(m.IP))
	// TCP端口号
	value := make([]byte, 2)
	binary.BigEndian.PutUint16(value, m.TCPPort)
	buf.Write(value)
	// UDP端口号
	binary.BigEndian.PutUint16(value, m.UDPPort)
	buf.Write(value)
	// 逻辑通道号
	buf.WriteByte(m.Channel)
	// 音视频类型
	buf.WriteByte(m.MediaType)
	// 码流类型
	buf.WriteByte(m.StreamType)
}

// NewMsgRealMediaRequest 新建实时音视频传输请求消息
func NewMsgRealMediaRequest() *MsgRealMediaRequest {
	return &MsgRealMediaRequest{
		OutputMark: OutputMark{
			ID: msgIDRealMediaRequest,
		},
	}
}

// MsgRealMediaControl 音视频实时传输控制
type MsgRealMediaControl struct {
	OutputMark
	// 逻辑通道号
	Channel byte
	// 控制指令，0-关闭音视频传输，1-切换码流，2-暂停，3-恢复，4关闭双向对讲
	Command byte
	// 操作类型，0-关闭该通道，1-只关闭音频，2-只关闭视频
	Operation byte
	// 切换码流类型，0-主码流，1-子码流
	StreamType byte
}

func (m *MsgRealMediaControl) writeTo(buf *bytes.Buffer) {
	// 逻辑通道号
	buf.WriteByte(m.Channel)
	// 控制指令
	buf.WriteByte(m.Command)
	// 操作类型
	buf.WriteByte(m.Operation)
	// 切换码流类型
	buf.WriteByte(m.StreamType)
}

// NewMsgRealMediaControl 新建音视频实时传输控制消息
func NewMsgRealMediaControl() *MsgRealMediaControl {
	return &MsgRealMediaControl{
		OutputMark: OutputMark{
			ID: msgIDRealMediaControl,
		},
	}
}

// MsgRealMediaStatusNotice 音视频实时传输状态通知
type MsgRealMediaStatusNotice struct {
	OutputMark
	//逻辑通道号
	Channel byte
	//丢包率
	PacketLossRate byte
}

func (m *MsgRealMediaStatusNotice) writeTo(buf *bytes.Buffer) {
	// 逻辑通道号
	buf.WriteByte(m.Channel)
	// 丢包率
	buf.WriteByte(m.PacketLossRate)
}

// NewMsgRealMediaStatusNotice 新建音视频实时传输状态通知消息
func NewMsgRealMediaStatusNotice() *MsgRealMediaStatusNotice {
	return &MsgRealMediaStatusNotice{
		OutputMark: OutputMark{
			ID: msgIDRealMediaNotice,
		},
	}
}

// MsgFileUploadCmd 文件上传指令
type MsgFileUploadCmd struct {
	OutputMark
	// 服务器地址
	ServerAddr string
	// 端口号
	ServerPort uint16
	// 用户名
	UserName string
	// 密码
	Password string
	// 文件上传路径
	Path string
	// 逻辑通道号
	Channel byte
	// 开始时间
	STime time.Time
	// 结束时间
	ETime time.Time
	// 报警标志
	Alarm uint64
	// 音视频资源类型，0-音视频，1-音频，2-视频，3-视频或音频
	MediaType byte
	// 码流类型，0-主码流或子码流，1-主码流，2-子码流
	StreamType byte
	// 存储位置，0-主存储器或灾备存储器，1-主存储器，2-灾备存储器
	StorageType byte
	// 任务执行条件
	ExeCondition byte
}

func (m *MsgFileUploadCmd) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 8)

	// 服务器地址长度
	buf.WriteByte(byte(len(m.ServerAddr)))
	// 服务器地址
	buf.Write([]byte(m.ServerAddr))
	// 端口
	binary.BigEndian.PutUint16(value, m.ServerPort)
	buf.Write(value[:2])
	// 用户名长度
	buf.WriteByte(byte(len(m.UserName)))
	// 用户名
	buf.Write([]byte(m.UserName))
	// 密码长度
	buf.WriteByte(byte(len(m.Password)))
	// 密码
	buf.Write([]byte(m.Password))
	// 文件上传路径长度
	buf.WriteByte(byte(len(m.Path)))
	// 文件上传路径
	buf.Write([]byte(m.Path))
	// 逻辑通道号
	buf.WriteByte(m.Channel)
	// 开始时间
	if m.STime.IsZero() {
		buf.Write(make([]byte, 6))
	} else {
		buf.Write(util.ToBCD([]byte(m.STime.Format("060102150405"))))
	}
	// 结束时间
	if m.ETime.IsZero() {
		buf.Write(make([]byte, 6))
	} else {
		buf.Write(util.ToBCD([]byte(m.ETime.Format("060102150405"))))
	}
	// 报警标志
	binary.BigEndian.PutUint64(value, m.Alarm)
	buf.Write(value)
	// 音视频资源类型
	buf.WriteByte(m.MediaType)
	// 码流类型
	buf.WriteByte(m.StreamType)
	// 存储位置
	buf.WriteByte(m.StorageType)
	// 任务执行条件
	buf.WriteByte(m.ExeCondition)
}

// NewMsgFileUploadCmd 新建文件上传指令消息
func NewMsgFileUploadCmd() *MsgFileUploadCmd {
	return &MsgFileUploadCmd{
		OutputMark: OutputMark{
			ID: msgIDFileUploadCmd,
		},
	}
}

// MsgFileUploadCtl 文件上传控制
type MsgFileUploadCtl struct {
	OutputMark
	// 应答流水号
	ReqNum uint16
	// 上传控制
	Ctl byte
}

func (m *MsgFileUploadCtl) writeTo(buf *bytes.Buffer) {
	// 应答流水号
	value := make([]byte, 2)
	binary.BigEndian.PutUint16(value, m.ReqNum)
	buf.Write(value)
	// 上传控制
	buf.WriteByte(m.Ctl)
}

// NewMsgFileUploadCtl 新建文件上传控制消息
func NewMsgFileUploadCtl() *MsgFileUploadCtl {
	return &MsgFileUploadCtl{
		OutputMark: OutputMark{
			ID: msgIDFileUploadCtl,
		},
	}
}

// MsgPtzTurn 云台转动控制
type MsgPtzTurn struct {
	OutputMark
	// 逻辑通道号
	Channel byte
	// 方向
	Direction byte
	//速度
	Speed byte
}

func (m *MsgPtzTurn) writeTo(buf *bytes.Buffer) {
	// 逻辑通道号
	buf.WriteByte(m.Channel)
	// 方向
	buf.WriteByte(m.Direction)
	// 速度
	buf.WriteByte(m.Speed)
}

// NewMsgPtzTurn 新建云台转动控制消息
func NewMsgPtzTurn() *MsgPtzTurn {
	return &MsgPtzTurn{
		OutputMark: OutputMark{
			ID: msgIDPtzTurn,
		},
	}
}

// MsgPtzFocus 云台焦距控制
type MsgPtzFocus struct {
	OutputMark
	// 逻辑通道号
	Channel byte
	// 操作
	Operate byte
}

func (m *MsgPtzFocus) writeTo(buf *bytes.Buffer) {
	// 逻辑通道号
	buf.WriteByte(m.Channel)
	// 操作
	buf.WriteByte(m.Operate)
}

// NewMsgPtzFocus 新建云台焦距控制消息
func NewMsgPtzFocus() *MsgPtzFocus {
	return &MsgPtzFocus{
		OutputMark: OutputMark{
			ID: msgIDPtzFocus,
		},
	}
}

// MsgPtzAperture 云台光圈控制
type MsgPtzAperture struct {
	OutputMark
	// 逻辑通道号
	Channel byte
	// 操作
	Operate byte
}

func (m *MsgPtzAperture) writeTo(buf *bytes.Buffer) {
	// 逻辑通道号
	buf.WriteByte(m.Channel)
	// 操作
	buf.WriteByte(m.Operate)
}

// NewMsgPtzAperture 新建云台光圈控制消息
func NewMsgPtzAperture() *MsgPtzAperture {
	return &MsgPtzAperture{
		OutputMark: OutputMark{
			ID: msgIDPtzAperture,
		},
	}
}

// MsgPtzWiper 云台雨刷控制
type MsgPtzWiper struct {
	OutputMark
	// 逻辑通道号
	Channel byte
	// 操作
	Operate byte
}

func (m *MsgPtzWiper) writeTo(buf *bytes.Buffer) {
	// 逻辑通道号
	buf.WriteByte(m.Channel)
	// 操作
	buf.WriteByte(m.Operate)
}

// NewMsgPtzWiper 新建云台雨刷控制消息
func NewMsgPtzWiper() *MsgPtzWiper {
	return &MsgPtzWiper{
		OutputMark: OutputMark{
			ID: msgIDPtzWiper,
		},
	}
}

// MsgPtzFillLight 云台红外补光
type MsgPtzFillLight struct {
	OutputMark
	// 逻辑通道号
	Channel byte
	// 操作
	Operate byte
}

func (m *MsgPtzFillLight) writeTo(buf *bytes.Buffer) {
	// 逻辑通道号
	buf.WriteByte(m.Channel)
	// 操作
	buf.WriteByte(m.Operate)
}

// NewMsgPtzFillLight 新建云台红外补光消息
func NewMsgPtzFillLight() *MsgPtzFillLight {
	return &MsgPtzFillLight{
		OutputMark: OutputMark{
			ID: msgIDPtzFilllight,
		},
	}
}

// MsgPtzZoom 云台变倍控制
type MsgPtzZoom struct {
	OutputMark
	// 逻辑通道号
	Channel byte
	// 操作
	Operate byte
}

func (m *MsgPtzZoom) writeTo(buf *bytes.Buffer) {
	// 逻辑通道号
	buf.WriteByte(m.Channel)
	// 操作
	buf.WriteByte(m.Operate)
}

// NewMsgPtzZoom 新建云台变倍控制消息
func NewMsgPtzZoom() *MsgPtzZoom {
	return &MsgPtzZoom{
		OutputMark: OutputMark{
			ID: msgIDPtzZoom,
		},
	}
}
