package jtt

import (
	"JTTServer/util"
	"bytes"
	"encoding/binary"
	"time"

	"github.com/axgle/mahonia"
)

// JTT808协议版本号
const (
	// JTT808-2011版本号
	version2011 = byte(0)
	// JTT808-2019版本号
	version2019 = byte(1)
)

// Param 终端参数项
type Param struct {
	// 参数id
	ID uint32 `json:"id"`
	// 参数长度
	Len byte `json:"len"`
	// 参数值
	Value interface{} `json:"value"`
}

// 终端参数ID枚举
const (
	// 终端心跳发送间隔
	ParamIDHeartbeat = uint32(0x0001)
	// TCP消息应答超时时间，单位为秒
	ParamIDTCPTimeOut = uint32(0x0002)
	// TCP消息重传次数
	ParamIDTCPRetrans = uint32(0x0003)
	// UDP消息应答超时时间，单位为秒
	ParamIDUDPTimeOut = uint32(0x0004)
	// UDP消息重传次数
	ParamIDUDPRetrans = uint32(0x0005)
	// SMS消息应答超时时间，单位为秒
	ParamIDSMSTimeOut = uint32(0x0006)
	// SMS消息重传次数
	ParamIDSMSRetrans = uint32(0x0007)
	// 主服务器APN
	ParamIDMainAPN = uint32(0x0010)
	// 主服务器无线通信拨号用户名
	ParamIDMainDialUser = uint32(0x0011)
	// 主服务器无线通信拨号密码
	ParamIDMainDialPW = uint32(0x0012)
	// 主服务器IP端口号
	ParamIDMainServer = uint32(0x0013)
	// 备份服务器APN
	ParamIDSpareAPN = uint32(0x0014)
	// 备份服务器无线通信拨号用户名
	ParamIDSpareDialUser = uint32(0x0015)
	// 备份服务器无线通信拨号密码
	ParamIDSpareDialPW = uint32(0x0016)
	// 备用服务器IP端口号
	ParamIDSpareServer = uint32(0x0017)
	// 服务器TCP端口号 808-2011
	ParamIDServerTCPPort = uint32(0x0018)
	// 服务器UDP端口号 808-2011
	ParamIDServerUDPPort = uint32(0x0019)
	// 道路运输证IC卡认证主服务器IP
	ParamIDICServerIP = uint32(0x001A)
	// 道路运输证IC卡认证主服务器TCP端口
	ParamIDICServerTCPPort = uint32(0x001B)
	// 道路运输证IC卡认证主服务器UDP端口
	ParamIDICServerUDPPort = uint32(0x001C)
	// 道路运输证IC卡认证备用服务器IP
	ParamIDICSpareServerIP = uint32(0x001D)
	// 位置汇报策略
	ParamIDLocReptStrategy = uint32(0x0020)
	// 位置汇报方案
	ParamIDLocReptScheme = uint32(0x0021)
	// 驾驶员未登录时汇报时间间隔，单位为秒
	ParamIDLocTimeInLogOut = uint32(0x0022)
	// 从服务器APN
	ParamIDSlaveServerAPN = uint32(0x0023)
	// 从服务器无线通信拨号用户名
	ParamIDSlaveDialUser = uint32(0x0024)
	// 从服务器无线通信拨号密码
	ParamIDSlaveDialPW = uint32(0x0025)
	// 从服务器IP端口号
	ParamIDSlaveServer = uint32(0x0026)
	// 休眠时汇报时间间隔，单位为秒
	ParamIDLocTimeInSleep = uint32(0x0027)
	// 紧急报警时汇报时间间隔
	ParamIDLocTimeInEmAlarm = uint32(0x0028)
	// 默认汇报时间间隔
	ParamIDLocTimeDefault = uint32(0x0029)
	// 默认汇报距离间隔，单位为m
	ParamIDLocStepDefault = uint32(0x002C)
	// 驾驶员未登录时汇报距离间隔
	ParamIDLocStepInLogOut = uint32(0x002D)
	// 休眠时汇报距离间隔
	ParamIDLocStepInSleep = uint32(0x002E)
	// 紧急报警时汇报距离间隔
	ParamIDLocStepInEmAlarm = uint32(0x002F)
	// 拐点补传角度，值小于180
	ParamIDVertexAddForAngle = uint32(0x0030)
	// 电子围栏半径，单位为m
	ParamIDEleFenceRadius = uint32(0x0031)
	// 违规行驶时段范围
	ParamIDIllegalDriDuration = uint32(0x0032)
	// 监控平台电话号码
	ParamIDServerTEL = uint32(0x0040)
	// 复位电话号码
	ParamIDResetTEL = uint32(0x0041)
	// 恢复出厂设置电话号码
	ParamIDRestoreTEL = uint32(0x0042)
	// 监控平台SMS电话号码
	ParamIDServerSMSTEL = uint32(0x0043)
	// 接收终端SMS文本报警号码
	ParamIDTerSMSTextAlarmTEL = uint32(0x0044)
	// 终端电话接听策略
	ParamIDAnswerPhoneStrategy = uint32(0x0045)
	// 每次最长通话时间
	ParamIDMaxTalkTimeOnce = uint32(0x0046)
	// 当月最长通话时间
	ParamIDMaxTalkTimeMonth = uint32(0x0047)
	// 监听电话号码
	ParamIDMonitorTEL = uint32(0x0048)
	// 监管平台特权短信号码
	ParamIDServerOwnTextTEL = uint32(0x0049)
	// 报警屏蔽字
	ParamIDAlarmMaskWord = uint32(0x0050)
	// 报警发送文本SMS开关
	ParamIDAlarmTextSMSSwitch = uint32(0x0051)
	// 报警拍摄开关
	ParamIDAlarmShootSwitch = uint32(0x0052)
	// 报警拍摄存储标志
	ParamIDAlarmPhotoSaveFlag = uint32(0x0053)
	// 报警关键标志
	ParamIDAlarmKeyFlag = uint32(0x0054)
	// 最高速度，单位为千米每小时
	ParamIDMaxSpeed = uint32(0x0055)
	// 超速持续时间
	ParamIDDurOfOverspeed = uint32(0x0056)
	// 连续驾驶时间门限
	ParamIDDurOfDriving = uint32(0x0057)
	// 当天累计驾驶时间门限
	ParamIDTTimeOfDriToday = uint32(0x0058)
	// 最小休息时间
	ParamIDMinRestTime = uint32(0x0059)
	// 最长停车时间
	ParamIDMaxStopTime = uint32(0x005A)
	// 超速预警差值，单位为1/10km/h
	ParamIDOverSpeedWarnDiff = uint32(0x005B)
	// 疲劳驾驶预警差值
	ParamIDFatigueDriWarnValue = uint32(0x005C)
	// 碰撞报警参数设置
	ParamIDCrashAlarmParam = uint32(0x005D)
	// 侧翻报警参数设置
	ParamIDRolloverAlarmParam = uint32(0x005E)
	// 定时拍照控制
	ParamIDTimingShootParam = uint32(0x0064)
	// 定距拍照控制
	ParamIDIntervalShootParam = uint32(0x0065)
	// 图像/视频质量
	ParamIDDefinition = uint32(0x0070)
	// 亮度
	ParamIDBrightness = uint32(0x0071)
	// 对比度
	ParamIDContrast = uint32(0x0072)
	// 饱和度
	ParamIDSaturation = uint32(0x0073)
	// 色度
	ParamIDChroma = uint32(0x0074)

	// 音视频参数设置
	ParamIDAVParam = uint32(0x0075)
	// 音视频通道列表设置
	ParamIDAVChannelList = uint32(0x0076)
	// 单独视频通道设置
	ParamIDVideoSingleChannel = uint32(0x0077)
	// 特殊报警录像设置
	ParamIDVideoSpecialAlarm = uint32(0x0079)
	// 视频相关报警屏蔽字设置
	ParamIDVideoAlarmShieldWord = uint32(0x007A)
	// 图像分析报警设置
	ParamIDVideoAnalysisAlarm = uint32(0x007B)
	// 终端休眠唤醒模式设置
	ParamIDWakeUpMode = uint32(0x007C)

	// 车辆里程表读数，单位1/10km/h
	ParamIDOdometerReading = uint32(0x0080)
	// 车辆所在的省域ID
	ParamIDProvince = uint32(0x0081)
	// 车辆所在的市域ID
	ParamIDCity = uint32(0x0082)
	// 车牌号
	ParamIDLicencePlate = uint32(0x0083)
	// 车牌颜色
	ParamIDPlateColor = uint32(0x0084)
	// GNSS定位模式
	ParamIDGNSSMode = uint32(0x0090)
	// GNSS波特率
	ParamIDGNSSBaudRate = uint32(0x0091)
	// GNSS模块详细定位数据输出频率
	ParamIDGNSSStatFreq = uint32(0x0092)
	// GNSS模块详细定位数据采集频率
	ParamIDGNSSStatCollFreq = uint32(0x0093)
	// GNSS模块详细定位数据上传方式
	ParamIDGNSSStatUpSchema = uint32(0x0094)
	// GNSS模块详细定位数据上传设置
	ParamIDGNSSStatUpParam = uint32(0x0095)
	// CAN总线通道1采集时间间隔
	ParamIDCAN1CollInterval = uint32(0x0100)
	// CAN总线通道1上传时间间隔
	ParamIDCAN1UpInterval = uint32(0x0101)
	// CAN总线通道2采集时间间隔
	ParamIDCAN2CollInterval = uint32(0x0102)
	// CAN总线通道2上传时间间隔
	ParamIDCAN2UpInterval = uint32(0x0103)
	// CAN总线ID单独采集设置
	ParamIDCANCollParam = uint32(0x0110)
	// 硬控通信端口号
	ParamIDHWCtrlPort = uint32(0xF000)
	// 车辆VIN码
	ParamIDVIN = uint32(0xF001)
	// 0xF002未被使用
	// 车牌种类
	ParamIDLicencePlateKind = uint32(0xF003)
	// CMS平台地址
	ParamIDCMSServer = uint32(0xF004)
	//DMS告警分级速度阈值
	ParamSpeed = uint32(0xF101)
	//人脸朝向左阈值
	ParamYawLeftThreshold = uint32(0xF102)
	//人脸朝向右阈值
	ParamYawRightThreshold = uint32(0xF103)
	//人脸朝向上阈值
	ParamPitchUpThreshold = uint32(0xF104)
	//人脸朝向下阈值
	ParamPitchDownThreshold = uint32(0xF105)
	//警告文本1，默认：人脸丢失
	ParamWarningText1 = uint32(0xF106)
	//警告文本2，默认：请勿打电话
	ParamWarningText2 = uint32(0xF107)
	//警告文本3，默认：请勿喝水
	ParamWarningText3 = uint32(0xF108)
	//警告文本4，默认：请勿抽烟
	ParamWarningText4 = uint32(0xF109)
	//警告文本5，默认：请目视前方
	ParamWarningText5 = uint32(0xF10A)
	//警告文本6，默认：请勿过度抬头
	ParamWarningText6 = uint32(0xF10B)
	//警告文本7，默认：请勿过度低头
	ParamWarningText7 = uint32(0xF10C)
	//警告文本8，默认：请勿打瞌睡
	ParamWarningText8 = uint32(0xF10D)
	//警告文本9，默认：请勿疲劳驾驶
	ParamWarningText9 = uint32(0xF10E)
	//眼睛的高宽比阈值
	ParamEyeThreshold = uint32(0xF10F)
	//嘴巴的高宽比阈值
	ParamMouthThreshold = uint32(0xF110)
	//人脸质量分数阈值
	ParamFaceQualityThreshold = uint32(0xF111)
	//获取人脸质量分数超时时间阈值
	ParamGetFaceQualityTimeOutThreshold = uint32(0xF112)
	//ADAS告警分级速度阈值
	ParamSpeedThreshold = uint32(0xF113)
	//判定车辆是否属于当前车道线比值阈值
	ParamIsSameLaneThreshold = uint32(0xF114)
	//车距过近的车距阈值
	ParamVehicleWarningDist = uint32(0xF115)
	//碰撞预警时间阈值
	ParamCrashWarningTime = uint32(0xF116)
	// 脉冲系数
	ParamIDPulseFactor = uint32(0xF200)
	// 高级驾驶辅助系统参数
	ParamIDADAS = uint32(0xF364)
	// 驾驶员状态监测系统参数
	ParamIDDMS = uint32(0xF365)
	// 轮胎气压监测系统
	ParamIDTPMS = uint32(0xF366)
	// 盲点监测系统
	ParamIDBSD = uint32(0xF367)
)

// 终端参数类型对照表
var paramTypeMap = map[uint32]string{
	uint32(0x0001): "uint32", // 终端心跳发送间隔，单位为（m）
	uint32(0x0002): "uint32", // TCP消息应答超时时间，单位为（m）
	// Someone is needed to complete these comments! I have no time to do this shit!
	uint32(0x0003): "uint32",
	uint32(0x0004): "uint32",
	uint32(0x0005): "uint32",
	uint32(0x0006): "uint32",
	uint32(0x0007): "uint32",
	uint32(0x0010): "string",
	uint32(0x0011): "string",
	uint32(0x0012): "string",
	uint32(0x0013): "string",
	uint32(0x0014): "string",
	uint32(0x0015): "string",
	uint32(0x0016): "string",
	uint32(0x0017): "string",
	uint32(0x001A): "string",
	uint32(0x001B): "uint32",
	uint32(0x001C): "uint32",
	uint32(0x001D): "string",
	uint32(0x0020): "uint32",
	uint32(0x0021): "uint32",
	uint32(0x0022): "uint32",
	uint32(0x0023): "string",
	uint32(0x0024): "string",
	uint32(0x0025): "string",
	uint32(0x0026): "string",
	uint32(0x0027): "uint32",
	uint32(0x0028): "uint32",
	uint32(0x0029): "uint32",
	uint32(0x002C): "uint32",
	uint32(0x002D): "uint32",
	uint32(0x002E): "uint32",
	uint32(0x002F): "uint32",
	uint32(0x0030): "uint32",
	uint32(0x0031): "uint16",
	uint32(0x0032): "[]uint8",
	uint32(0x0040): "string",
	uint32(0x0041): "string",
	uint32(0x0042): "string",
	uint32(0x0043): "string",
	uint32(0x0044): "string",
	uint32(0x0045): "uint32",
	uint32(0x0046): "uint32",
	uint32(0x0047): "uint32",
	uint32(0x0048): "string",
	uint32(0x0049): "string",
	uint32(0x0050): "uint32",
	uint32(0x0051): "uint32",
	uint32(0x0052): "uint32",
	uint32(0x0053): "uint32",
	uint32(0x0054): "uint32",
	uint32(0x0055): "uint32",
	uint32(0x0056): "uint32",
	uint32(0x0057): "uint32",
	uint32(0x0058): "uint32",
	uint32(0x0059): "uint32",
	uint32(0x005A): "uint32",
	uint32(0x005B): "uint16",
	uint32(0x005C): "uint16",
	uint32(0x005D): "uint16",
	uint32(0x005E): "uint16",
	uint32(0x0064): "uint32",
	uint32(0x0065): "uint32",
	uint32(0x0070): "uint32",
	uint32(0x0071): "uint32",
	uint32(0x0072): "uint32",
	uint32(0x0073): "uint32",
	uint32(0x0074): "uint32",
	uint32(0x0075): "[]uint8",
	uint32(0x0076): "[]uint8",
	uint32(0x0077): "[]uint8",
	uint32(0x0079): "[]uint8",
	uint32(0x007A): "uint32",
	uint32(0x007B): "[]uint8",
	uint32(0x007C): "[]uint8",
	uint32(0x0080): "uint32",
	uint32(0x0081): "uint16",
	uint32(0x0082): "uint16",
	uint32(0x0083): "string",
	uint32(0x0084): "uint8",
	uint32(0x0090): "uint8",
	uint32(0x0091): "uint8",
	uint32(0x0092): "uint8",
	uint32(0x0093): "uint32",
	uint32(0x0094): "uint8",
	uint32(0x0095): "uint32",  // GNSS模块详细定位数据上传设置
	uint32(0x0100): "uint32",  // CAN总线通道1采集时间间隔，单位为毫秒（s）
	uint32(0x0101): "uint16",  // CAN总线通道1上传时间间隔，单位为秒（s）
	uint32(0x0102): "uint32",  // CAN总线通道2采集时间间隔，单位为毫秒（s）
	uint32(0x0103): "uint16",  // CAN总线通道2上传时间间隔，单位为秒（s）
	uint32(0x0110): "[]uint8", // CAN总线ID单独采集设置
	// 0x0111~0x01FF 用于其他CAN总线ID单独采集设置，此处有待优化
	uint32(0x0111): "[]uint8", // 其他CAN总线ID单独采集设置
	uint32(0x0112): "[]uint8", // 其他CAN总线ID单独采集设置
	uint32(0x0113): "[]uint8", // 其他CAN总线ID单独采集设置
	uint32(0x0114): "[]uint8", // 其他CAN总线ID单独采集设置
	uint32(0x0115): "[]uint8", // 其他CAN总线ID单独采集设置
	uint32(0x0116): "[]uint8", // 其他CAN总线ID单独采集设置
	uint32(0x0117): "[]uint8", // 其他CAN总线ID单独采集设置
	uint32(0x0118): "[]uint8", // 其他CAN总线ID单独采集设置
	uint32(0x0119): "[]uint8", // 其他CAN总线ID单独采集设置
	uint32(0xF000): "uint16",  // 硬控通信端口号
	uint32(0xF001): "string",  // 车架号
	uint32(0xF003): "uint8",   // 车牌种类
	uint32(0xF004): "string",  // CMS平台地址
	uint32(0xF101): "uint32",  //DMS告警分级速度阈值
	uint32(0xF102): "uint32",  //人脸朝向左阈值
	uint32(0xF103): "uint32",  //人脸朝向右阈值
	uint32(0xF104): "uint32",  //人脸朝向上阈值
	uint32(0xF105): "uint32",  //人脸朝向下阈值
	uint32(0xF106): "string",  //警告文本1"人脸丢失"
	uint32(0xF107): "string",  //警告文本2"请勿打电话"
	uint32(0xF108): "string",  //警告文本3"请勿喝水"
	uint32(0xF109): "string",  //警告文本4"请勿抽烟"
	uint32(0xF10A): "string",  //警告文本5"请目视前方"
	uint32(0xF10B): "string",  //警告文本6"请勿过度抬头"
	uint32(0xF10C): "string",  //警告文本7"请勿过度低头"
	uint32(0xF10D): "string",  //警告文本8"请勿打瞌睡"
	uint32(0xF10E): "string",  //警告文本9"请勿疲劳驾驶"
	uint32(0xF10F): "uint32",  //眼睛的高宽比阈值
	uint32(0xF110): "uint32",  //嘴巴的高宽比阈值
	uint32(0xF111): "uint32",  //人脸质量分数阈值
	uint32(0xF112): "uint32",  //获取人脸质量分数超时时间阈值
	uint32(0xF113): "uint32",  //ADAS告警分级速度阈值
	uint32(0xF114): "uint32",  //判定车辆是否属于当前车道线比值阈值
	uint32(0xF115): "uint32",  //车距过近的车距阈值
	uint32(0xF116): "uint32",  //碰撞预警时间阈值
	uint32(0xF200): "uint32",  // 脉冲系数
	uint32(0xF364): "[]uint8", // 高级驾驶辅助系统参数
	uint32(0xF365): "[]uint8", // 驾驶员状态监测系统参数
	uint32(0xF366): "[]uint8", // 胎压监测系统参数
	uint32(0xF367): "[]uint8", // 盲区监测系统参数
}

// GetParamTypeMap 获取终端参数类型映射表
func GetParamTypeMap() map[uint32]string {
	return paramTypeMap
}

// GNSSAttr GNSS模块属性
type GNSSAttr byte

// COMMAttr 通信模块属性
type COMMAttr byte

// TerminalAttr 终端属性
type TerminalAttr struct {
	// 终端类型
	TerType uint16 `json:"ter_type"`
	// 制造商id
	VendorID string `json:"vendor_id"`
	// 终端型号
	Model string `json:"model"`
	// 终端id
	TerID string `json:"ter_id"`
	// 终端SIM卡ICCID号
	ICCID string `json:"iccid"`
	// 终端本机号码
	Phone string `json:"phone"`
	// 国际移动设备识别码
	IMEI string `json:"imei"`
	// 终端硬件版本号
	HWVersion string `json:"hw_version"`
	// 终端固件版本号
	FWVersion string `json:"fw_version"`
	// GNSS模块属性
	GnssAttr GNSSAttr `json:"gnss_attr"`
	// 通信模块属性
	CommAttr COMMAttr `json:"comm_attr"`
	// 初次安装时间
	SetupTimestamp int64 `json:"setup_timestamp"`
}

func (t *TerminalAttr) readBy(buf *bytes.Buffer, version byte) {
	// 终端类型
	t.TerType = binary.BigEndian.Uint16(buf.Next(2))
	if version2011 == version {
		// 制造商id
		t.VendorID = string(buf.Next(5))
		// 终端型号
		t.Model = string(buf.Next(20))
		// 终端id
		t.TerID = string(buf.Next(7))
	} else {
		t.VendorID = string(buf.Next(11))
		// 终端型号
		t.Model = string(buf.Next(30))
		// 终端id
		t.TerID = string(buf.Next(30))
	}
	// 终端SIM卡ICCID号
	t.ICCID = string(buf.Next(10))
	// 终端硬件版本号
	decoder := mahonia.NewDecoder("gbk")
	len, _ := buf.ReadByte()
	t.HWVersion = decoder.ConvertString(string(buf.Next(int(len))))
	// 终端固件版本号
	len, _ = buf.ReadByte()
	t.FWVersion = decoder.ConvertString(string(buf.Next(int(len))))
	// GNSS模块属性
	bt, _ := buf.ReadByte()
	t.GnssAttr = GNSSAttr(bt)
	// 通信模块属性
	bt, _ = buf.ReadByte()
	t.CommAttr = COMMAttr(bt)
}

type VehicleAttr struct {
	// 车辆载货状态
	CargoStatus uint8 `json:"cargo_status"`
}

// 位置附加项信息id
const (
	// 里程
	AuxIDMileage = byte(0x01)
	// 油量
	AuxIDOilMass = byte(0x02)
	// 行驶记录获取的速度
	AuxIDRecordSpeed = byte(0x03)
	// 需要人工确认报警事件的ID
	AuxIDAlarmEventID = byte(0x04)
	// 胎压
	AuxIDTirePressure = byte(0x05)
	// 车厢温度
	AuxIDCarriageTemp = byte(0x06)
	// 冷链货仓温度1
	AuxIDCoolerTemp1 = byte(0x07)
	// 冷链货仓温度2
	AuxIDCoolerTemp2 = byte(0x08)
	// 冷链货仓湿度1
	AuxIDCoolerHumi1 = byte(0x09)
	// 冷链货仓湿度2
	AuxIDCoolerHumi2 = byte(0x0A)
	// 附加报警项
	AuxIDAlarmPlus = byte(0x0B)
	// 门磁1状态
	AuxIDDoorMagnetic1 = byte(0x0C)
	// 门磁2状态
	AuxIDDoorMagnetic2 = byte(0x0D)
	// 超速报警信息
	AuxIDSpeedAlarm = byte(0x11)
	// 进出区域报警信息
	AuxIDLocalAlarm = byte(0x12)
	// 路线行驶时间过长/不足报警信息
	AuxIDPathTimeAlarm = byte(0x13)
	// 视频相关报警
	AuxIDVideoRelatedAlarm = byte(0x14)
	// 视频信号丢失
	AuxIDVideoSignalLoss = byte(0x15)
	// 视频信号遮挡
	AuxIDVideoSignalOcclude = byte(0x16)
	// 视频存储器故障
	AuxIDVideoMemoryFault = byte(0x17)
	// 视频异常驾驶行为
	AuxIDVideoAbnormalDriving = byte(0x18)
	// 车辆扩展状态位
	AuxIDVehicleStatus = byte(0x25)
	// IO状态位
	AuxIDIOStatus = byte(0x2A)
	// 模拟量
	AuxIDAnalogQuantity = byte(0x2B)
	// 无线通信网络信号强度
	AuxIDWlanSignal = byte(0x30)
	// GNSS定位卫星数
	AuxIDGnssStvCnt = byte(0x31)

	// AI报警
	AuxIDAIAlarm = byte(0xE2)
	// ADAS报警
	AuxIDAdasAlarm = byte(0x64)
	// DMS报警
	AuxIDDMSAlarm = byte(0x65)
	// 轮胎状态监测报警
	AuxIDTPMSAlarm = byte(0x66)
	// 盲点检测报警
	AuxIDBsdAlarm = byte(0x67)
)

// AlarmFlagMask 硬控终端故障报警掩码
const AlarmFlagMask = AlarmFlag(0x7F031FE1)

// AlarmFlag 报警标识
type AlarmFlag uint32

// IsEmergency 是否是紧急报警
func (a *AlarmFlag) IsEmergency() bool {
	return *a&1 > 0
}

// Emergency 紧急报警设置
func (a *AlarmFlag) Emergency(enabled bool) {
	if enabled {
		*a |= 1
	} else {
		*a &= 0xFFFFFFFE
	}
}

// IsOverSpeed 是否是超速报警
func (a *AlarmFlag) IsOverSpeed() bool {
	return *a&2 > 0
}

// OverSpeed 超速报警设置
func (a *AlarmFlag) OverSpeed(enabled bool) {
	if enabled {
		*a |= 2
	} else {
		*a &= 0xFFFFFFFD
	}
}

// IsTired 是否是疲劳驾驶报警
func (a *AlarmFlag) IsTired() bool {
	return *a&4 > 0
}

// Tired 疲劳驾驶报警设置
func (a *AlarmFlag) Tired(enabled bool) {
	if enabled {
		*a |= 4
	} else {
		*a &= 0xFFFFFFFB
	}
}

// IsDangerous 是否是危险驾驶行为报警
func (a *AlarmFlag) IsDangerous() bool {
	return *a&8 > 0
}

// Dangerous 危险驾驶行为报警设置
func (a *AlarmFlag) Dangerous(enabled bool) {
	if enabled {
		*a |= 8
	} else {
		*a &= 0xFFFFFFF7
	}
}

// IsGNSSBreakdown 是否是GNSS故障报警
func (a *AlarmFlag) IsGNSSBreakdown() bool {
	return *a&16 > 0
}

// GNSSBreakdown GNSS故障报警设置
func (a *AlarmFlag) GNSSBreakdown(enabled bool) {
	if enabled {
		*a |= 16
	} else {
		*a &= 0xFFFFFFEF
	}
}

// IsGNSSUnhook 是否是GNSS未接或剪断报警
func (a *AlarmFlag) IsGNSSUnhook() bool {
	return *a&32 > 0
}

// GNSSUnhook GNSS未接或剪断报警设置
func (a *AlarmFlag) GNSSUnhook(enabled bool) {
	if enabled {
		*a |= 32
	} else {
		*a &= 0xFFFFFFDF
	}
}

// IsGNSSShortOut 是否是GNSS天线短路报警
func (a *AlarmFlag) IsGNSSShortOut() bool {
	return *a&64 > 0
}

// GNSSShortOut GNSS天线短路报警设置
func (a *AlarmFlag) GNSSShortOut(enabled bool) {
	if enabled {
		*a |= 64
	} else {
		*a &= 0xFFFFFFBF
	}
}

// IsMainLowVoltage 是否是主电源低电压报警
func (a *AlarmFlag) IsMainLowVoltage() bool {
	return *a&128 > 0
}

// MainLowVoltage 主电源低电压报警设置
func (a *AlarmFlag) MainLowVoltage(enabled bool) {
	if enabled {
		*a |= 128
	} else {
		*a &= 0xFFFFFF7F
	}
}

// IsMainPowerDown 是否是主电源掉电报警
func (a *AlarmFlag) IsMainPowerDown() bool {
	return *a&256 > 0
}

// MainPowerDown 主电源掉电报警设置
func (a *AlarmFlag) MainPowerDown(enabled bool) {
	if enabled {
		*a |= 256
	} else {
		*a &= 0xFFFFFEFF
	}
}

// IsDisplayBreakdown 是否是LCD或显示器故障报警
func (a *AlarmFlag) IsDisplayBreakdown() bool {
	return *a&512 > 0
}

// DisplayBreakdown LCD或显示器故障报警设置
func (a *AlarmFlag) DisplayBreakdown(enabled bool) {
	if enabled {
		*a |= 512
	} else {
		*a &= 0xFFFFFDFF
	}
}

// IsTTSBreakdown 是否是TSS模块故障报警
func (a *AlarmFlag) IsTTSBreakdown() bool {
	return *a&1024 > 0
}

// TTSBreakdown TSS模块故障报警设置
func (a *AlarmFlag) TTSBreakdown(enabled bool) {
	if enabled {
		*a |= 1024
	} else {
		*a &= 0xFFFFFBFF
	}
}

// IsCameraBreakdown 是否是摄像头故障报警
func (a *AlarmFlag) IsCameraBreakdown() bool {
	return *a&2048 > 0
}

// CameraBreakdown 摄像头故障报警设置
func (a *AlarmFlag) CameraBreakdown(enabled bool) {
	if enabled {
		*a |= 2048
	} else {
		*a &= 0xFFFFF7FF
	}
}

// IsICCardBreakdown 是否是IC卡模块故障报警
func (a *AlarmFlag) IsICCardBreakdown() bool {
	return *a&4096 > 0
}

// ICCardBreakdown IC卡模块故障报警设置
func (a *AlarmFlag) ICCardBreakdown(enabled bool) {
	if enabled {
		*a |= 4096
	} else {
		*a &= 0xFFFFEFFF
	}
}

// IsOverSpeedEarly 是否是超速预警
func (a *AlarmFlag) IsOverSpeedEarly() bool {
	return *a&8192 > 0
}

// OverSpeedEarly 超速预警设置
func (a *AlarmFlag) OverSpeedEarly(enabled bool) {
	if enabled {
		*a |= 8192
	} else {
		*a &= 0xFFFFDFFF
	}
}

// IsTiredEarly 是否是疲劳驾驶预警
func (a *AlarmFlag) IsTiredEarly() bool {
	return *a&16384 > 0
}

// TiredEarly 疲劳驾驶预警设置
func (a *AlarmFlag) TiredEarly(enabled bool) {
	if enabled {
		*a |= 16384
	} else {
		*a &= 0xFFFFBFFF
	}
}

// IsViolation 是否是违规行驶报警
func (a *AlarmFlag) IsViolation() bool {
	return *a&32768 > 0
}

// Violation 违规行驶报警设置
func (a *AlarmFlag) Violation(enabled bool) {
	if enabled {
		*a |= 32768
	} else {
		*a &= 0xFFFF7FFF
	}
}

// IsTirePressureEarly 是否是胎压预警
func (a *AlarmFlag) IsTirePressureEarly() bool {
	return *a&65536 > 0
}

// TirePressureEarly 胎压预警设置
func (a *AlarmFlag) TirePressureEarly(enabled bool) {
	if enabled {
		*a |= 65536
	} else {
		*a &= 0xFFFEFFFF
	}
}

// IsRightBlindAbnormal 是否是右转盲区异常报警
func (a *AlarmFlag) IsRightBlindAbnormal() bool {
	return *a&131072 > 0
}

// RightBlindAbnormal 右转盲区异常报警设置
func (a *AlarmFlag) RightBlindAbnormal(enabled bool) {
	if enabled {
		*a |= 131072
	} else {
		*a &= 0xFFFDFFFF
	}
}

// IsTimeout 是否是当天累计驾驶超时报警
func (a *AlarmFlag) IsTimeout() bool {
	return *a&262144 > 0
}

// Timeout 当天累计驾驶超时报警设置
func (a *AlarmFlag) Timeout(enabled bool) {
	if enabled {
		*a |= 262144
	} else {
		*a &= 0xFFFBFFFF
	}
}

// IsTimeoutParking 是否是超时停车报警
func (a *AlarmFlag) IsTimeoutParking() bool {
	return *a&524288 > 0
}

// TimeoutParking 超时停车报警设置
func (a *AlarmFlag) TimeoutParking(enabled bool) {
	if enabled {
		*a |= 524288
	} else {
		*a &= 0xFFF7FFFF
	}
}

// IsPassArea 是否是进出区域报警
func (a *AlarmFlag) IsPassArea() bool {
	return *a&1048576 > 0
}

// PassArea 进出区域报警设置
func (a *AlarmFlag) PassArea(enabled bool) {
	if enabled {
		*a |= 1048576
	} else {
		*a &= 0xFFEFFFFF
	}
}

// IsPassPath 是否是进出路线报警
func (a *AlarmFlag) IsPassPath() bool {
	return *a&2097152 > 0
}

// PassPath 进出路线报警设置
func (a *AlarmFlag) PassPath(enabled bool) {
	if enabled {
		*a |= 2097152
	} else {
		*a &= 0xFFDFFFFF
	}
}

// IsDurationErrInSegment 是否是路段行驶时间不足/过长报警
func (a *AlarmFlag) IsDurationErrInSegment() bool {
	return *a&4194304 > 0
}

// DurationErrInSegment 路段行驶时间不足/过长报警设置
func (a *AlarmFlag) DurationErrInSegment(enabled bool) {
	if enabled {
		*a |= 4194304
	} else {
		*a &= 0xFFBFFFFF
	}
}

// IsDeviatePath 是否是路线偏离报警
func (a *AlarmFlag) IsDeviatePath() bool {
	return *a&8388608 > 0
}

// DeviatePath 路线偏离报警设置
func (a *AlarmFlag) DeviatePath(enabled bool) {
	if enabled {
		*a |= 8388608
	} else {
		*a &= 0xFF7FFFFF
	}
}

// IsVSSBreakdown 是否是车辆VSS报警
func (a *AlarmFlag) IsVSSBreakdown() bool {
	return *a&16777216 > 0
}

// VSSBreakdown 车辆VSS报警设置
func (a *AlarmFlag) VSSBreakdown(enabled bool) {
	if enabled {
		*a |= 16777216
	} else {
		*a &= 0xFEFFFFFF
	}
}

// IsOilMassAbnormal 是否是油量异常报警
func (a *AlarmFlag) IsOilMassAbnormal() bool {
	return *a&33554432 > 0
}

// OilMassAbnormal 油量异常报警设置
func (a *AlarmFlag) OilMassAbnormal(enabled bool) {
	if enabled {
		*a |= 33554432
	} else {
		*a &= 0xFDFFFFFF
	}
}

// IsBeStolen 是否是车辆被盗报警
func (a *AlarmFlag) IsBeStolen() bool {
	return *a&67108864 > 0
}

// BeStolen 车辆被盗报警设置
func (a *AlarmFlag) BeStolen(enabled bool) {
	if enabled {
		*a |= 67108864
	} else {
		*a &= 0xFBFFFFFF
	}
}

// IsIllegalBoot 是否是车辆非法启动报警
func (a *AlarmFlag) IsIllegalBoot() bool {
	return *a&134217728 > 0
}

// IllegalBoot 车辆非法启动报警设置
func (a *AlarmFlag) IllegalBoot(enabled bool) {
	if enabled {
		*a |= 134217728
	} else {
		*a &= 0xF7FFFFFF
	}
}

// IsIllegalMove 是否是车辆非法移动报警
func (a *AlarmFlag) IsIllegalMove() bool {
	return *a&268435456 > 0
}

// IllegalMove 车辆非法移动报警设置
func (a *AlarmFlag) IllegalMove(enabled bool) {
	if enabled {
		*a |= 268435456
	} else {
		*a &= 0xEFFFFFFF
	}
}

// IsRollover 是否是车辆碰撞侧翻报警
func (a *AlarmFlag) IsRollover() bool {
	return *a&536870912 > 0
}

// Rollover 车辆碰撞侧翻报警设置
func (a *AlarmFlag) Rollover(enabled bool) {
	if enabled {
		*a |= 536870912
	} else {
		*a &= 0xDFFFFFFF
	}
}

// IsRolloverEarly 是否是车辆侧翻预警
func (a *AlarmFlag) IsRolloverEarly() bool {
	return *a&1073741824 > 0
}

// RolloverEarly 车辆侧翻预警设置
func (a *AlarmFlag) RolloverEarly(enabled bool) {
	if enabled {
		*a |= 1073741824
	} else {
		*a &= 0xBFFFFFFF
	}
}

// HasOtherAlarm 是否有其他报警
func (a *AlarmFlag) HasOtherAlarm() bool {
	return *a&2147483648 > 0
}

// OtherAlarm 存在其他报警设置
func (a *AlarmFlag) OtherAlarm(enabled bool) {
	if enabled {
		*a |= 2147483648
	} else {
		*a &= 0x7FFFFFFF
	}
}

// SpeedAlarm 超速报警附加信息项
type SpeedAlarm struct {
	// 位置类型
	Shape ShapeType
	// 区域或路段id
	TagID uint32
}

// ShapeType 超速报警区域的类型枚举
type ShapeType = uint8

const (
	// ShapeTypeNone 无类型
	ShapeTypeNone ShapeType = iota
	// ShapeTypeCircle 圆形
	ShapeTypeCircle
	// ShapeTypeRect 矩形
	ShapeTypeRect
	// ShapeTypePolygon 多边形
	ShapeTypePolygon
	// ShapeTypeRoad 路段
	ShapeTypeRoad
)

// readBy 从缓存中读
func (s *SpeedAlarm) readBy(buf *bytes.Buffer) {
	shape, _ := buf.ReadByte()
	s.Shape = ShapeType(shape)
	if ShapeTypeNone == s.Shape {
		return
	}

	value := make([]byte, 4)
	buf.Read(value)
	s.TagID = binary.BigEndian.Uint32(value)
}

// writeTo 写入缓存
func (s *SpeedAlarm) writeTo(buf *bytes.Buffer) {
	buf.WriteByte(byte(s.Shape))
	if ShapeTypeNone != s.Shape {
		value := make([]byte, 4)
		binary.BigEndian.PutUint32(value, s.TagID)
		buf.Write(value)
	}
}

// size 数据大小
func (s *SpeedAlarm) size() byte {
	if ShapeTypeNone == s.Shape {
		return byte(1)
	}

	return byte(5)
}

// LocalAlarm 进出区域/路线报警附加信息项
type LocalAlarm struct {
	// 位置类型
	Shape ShapeType
	// 区域或路线id
	TagID uint32
	// 方向，0-进，1-出
	Direction byte
}

// readBy 从缓存中读
func (l *LocalAlarm) readBy(buf *bytes.Buffer) {
	shape, _ := buf.ReadByte()
	l.Shape = ShapeType(shape)

	value := make([]byte, 4)
	buf.Read(value)
	l.TagID = binary.BigEndian.Uint32(value)

	l.Direction, _ = buf.ReadByte()
}

// writeTo 写入缓存
func (l *LocalAlarm) writeTo(buf *bytes.Buffer) {
	buf.WriteByte(byte(l.Shape))
	value := make([]byte, 4)
	binary.BigEndian.PutUint32(value, l.TagID)
	buf.Write(value)
	buf.WriteByte(l.Direction)
}

// size 数据大小
func (l *LocalAlarm) size() byte {
	return byte(6)
}

// RuntimeAlarm 路线行驶时间不足/过长报警附加信息项
type RuntimeAlarm struct {
	// 路段id
	TagID uint32
	// 路段行驶时间，单位（s）
	Duration uint16
	// 结果，0-不足，1-过长
	Verdict byte
}

// readBy 从缓存中读
func (r *RuntimeAlarm) readBy(buf *bytes.Buffer) {
	value := make([]byte, 4)

	buf.Read(value)
	r.TagID = binary.BigEndian.Uint32(value)
	buf.Read(value[:2])
	r.Duration = binary.BigEndian.Uint16(value)
	r.Verdict, _ = buf.ReadByte()
}

// writeTo 写入缓存中
func (r *RuntimeAlarm) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 4)

	binary.BigEndian.PutUint32(value, r.TagID)
	buf.Write(value)
	binary.BigEndian.PutUint16(value, r.Duration)
	buf.Write(value[:2])
	buf.WriteByte(r.Verdict)
}

// size 数据大小
func (r *RuntimeAlarm) size() byte {
	return byte(7)
}

// StatusFlagMask 硬控终端状态掩码
//const StatusFlagMask = StatusFlag(0x00000001)
const StatusFlagMask = StatusFlag(0x00000301)

// StatusFlag 状态标识
type StatusFlag uint32

// IsACCOpen ACC是否已打开
func (s *StatusFlag) IsACCOpen() bool {
	return *s&1 > 0
}

// ACCOpen 开关ACC
func (s *StatusFlag) ACCOpen(open bool) {
	if open {
		*s |= 1
	} else {
		*s &= 0xFFFFFFFE
	}
}

// IsFixed 是否已定位
func (s *StatusFlag) IsFixed() bool {
	return *s&2 > 0
}

// Fixed 定位与否设置
func (s *StatusFlag) Fixed(fixed bool) {
	if fixed {
		*s |= 2
	} else {
		*s &= 0xFFFFFFFD
	}
}

// IsSourth 是否是南北纬
func (s *StatusFlag) IsSourth() bool {
	return *s&4 > 0
}

// Sourth 南北纬设置
func (s *StatusFlag) Sourth(truth bool) {
	if truth {
		*s |= 4
	} else {
		*s &= 0xFFFFFFFB
	}
}

// IsWest 是否是东西经
func (s *StatusFlag) IsWest() bool {
	return *s&8 > 0
}

// West 东西经设置
func (s *StatusFlag) West(truth bool) {
	if truth {
		*s |= 8
	} else {
		*s &= 0xFFFFFFF7
	}
}

// IsOutage 是否是停运状态
func (s *StatusFlag) IsOutage() bool {
	return *s&16 > 0
}

// Outage 停运设置
func (s *StatusFlag) Outage(truth bool) {
	if truth {
		*s |= 16
	} else {
		*s &= 0xFFFFFFEF
	}
}

// IsEncrypted 是否是经纬度已加密
func (s *StatusFlag) IsEncrypted() bool {
	return *s&32 > 0
}

// Encrypted 经纬度已加密设置
func (s *StatusFlag) Encrypted(truth bool) {
	if truth {
		*s |= 32
	} else {
		*s &= 0xFFFFFFDF
	}
}

// IsCrashEarly 是否是前撞预警
func (s *StatusFlag) IsCrashEarly() bool {
	return *s&64 > 0
}

// CrashEarly 前撞预警设置
func (s *StatusFlag) CrashEarly(truth bool) {
	if truth {
		*s |= 64
	} else {
		*s &= 0xFFFFFFBF
	}
}

// IsLaneShiftEarly 是否是车道偏移预警
func (s *StatusFlag) IsLaneShiftEarly() bool {
	return *s&128 > 0
}

// LaneShiftEarly 车道偏移预警设置
func (s *StatusFlag) LaneShiftEarly(truth bool) {
	if truth {
		*s |= 128
	} else {
		*s &= 0xFFFFFF7F
	}
}

// GetCargoStatus 获取载货状态
func (s *StatusFlag) GetCargoStatus() uint32 {
	return uint32(*s&768) >> 8
}

// SetCargoStatus 设置载货状态
func (s *StatusFlag) SetCargoStatus(status uint32) {
	*s |= StatusFlag((status & 3) << 8)
}

// IsOffOil 是否是油路断开
func (s *StatusFlag) IsOffOil() bool {
	return *s&1024 > 0
}

// OffOil 油路断开设置
func (s *StatusFlag) OffOil(truth bool) {
	if truth {
		*s |= 1024
	} else {
		*s &= 0xFFFFFDFF
	}
}

// IsOffCircuit 是否是电路断开
func (s *StatusFlag) IsOffCircuit() bool {
	return *s&2048 > 0
}

// OffCircuit 电路断开设置
func (s *StatusFlag) OffCircuit(truth bool) {
	if truth {
		*s |= 2048
	} else {
		*s &= 0xFFFFFBFF
	}
}

// IsDoorLock 是否是车门锁定
func (s *StatusFlag) IsDoorLock() bool {
	return *s&4096 > 0
}

// DoorLock 车门锁定设置
func (s *StatusFlag) DoorLock(truth bool) {
	if truth {
		*s |= 4096
	} else {
		//*s &= 0xFFFFF7FF
		*s &= 0xFFFFEFFF
	}
}

// IsFrontDoorOpened 是否是前门已打开
func (s *StatusFlag) IsFrontDoorOpened() bool {
	return *s&8192 > 0
}

// FrontDoorOpened 前门已打开设置
func (s *StatusFlag) FrontDoorOpened(truth bool) {
	if truth {
		*s |= 8192
	} else {
		*s &= 0xFFFFEFFF
	}
}

// IsMiddleDoorOpened 是否是中门已打开
func (s *StatusFlag) IsMiddleDoorOpened() bool {
	return *s&16384 > 0
}

// MiddleDoorOpened 中门已打开设置
func (s *StatusFlag) MiddleDoorOpened(truth bool) {
	if truth {
		*s |= 16384
	} else {
		*s &= 0xFFFFDFFF
	}
}

// IsBackDoorOpened 是否是后门已打开
func (s *StatusFlag) IsBackDoorOpened() bool {
	return *s&32768 > 0
}

// BackDoorOpened 后门已打开设置
func (s *StatusFlag) BackDoorOpened(truth bool) {
	if truth {
		*s |= 32768
	} else {
		*s &= 0xFFFFBFFF
	}
}

// IsSideDoorOpened 是否是驾驶席门已打开
func (s *StatusFlag) IsSideDoorOpened() bool {
	return *s&65536 > 0
}

// SideDoorOpened 驾驶席门已打开设置
func (s *StatusFlag) SideDoorOpened(truth bool) {
	if truth {
		*s |= 65536
	} else {
		*s &= 0xFFFF7FFF
	}
}

// IsOtherDoorOpened 是否是自定义门已打开
func (s *StatusFlag) IsOtherDoorOpened() bool {
	return *s&131072 > 0
}

// OtherDoorOpened 自定义门已打开设置
func (s *StatusFlag) OtherDoorOpened(truth bool) {
	if truth {
		*s |= 131072
	} else {
		*s &= 0xFFFEFFFF
	}
}

// IsGPSOpened 是否有使用GPS定位
func (s *StatusFlag) IsGPSOpened() bool {
	return *s&262144 > 0
}

// GPSOpened 使用GPS定位设置
func (s *StatusFlag) GPSOpened(truth bool) {
	if truth {
		*s |= 262144
	} else {
		*s &= 0xFFFBFFFF //0xFFFDFFFF
	}
}

// IsBDOpened 是否有使用BD定位
func (s *StatusFlag) IsBDOpened() bool {
	return *s&524288 > 0
}

// BDOpened 使用BD定位设置
func (s *StatusFlag) BDOpened(truth bool) {
	if truth {
		*s |= 524288
	} else {
		*s &= 0xFFF7FFFF //0xFFFBFFFF
	}
}

// IsGLOpened 是否有使用GL定位
func (s *StatusFlag) IsGLOpened() bool {
	return *s&1048576 > 0
}

// GLOpened 使用GL定位设置
func (s *StatusFlag) GLOpened(truth bool) {
	if truth {
		*s |= 1048576
	} else {
		*s &= 0xFFF7FFFF
	}
}

// IsGAOpened 是否有使用伽利略定位
func (s *StatusFlag) IsGAOpened() bool {
	return *s&2097152 > 0
}

// GAOpened 使用伽利略定位设置
func (s *StatusFlag) GAOpened(truth bool) {
	if truth {
		*s |= 2097152
	} else {
		*s &= 0xFFEFFFFF
	}
}

// IsRunning 是否是行驶状态
func (s *StatusFlag) IsRunning() bool {
	return *s&4194304 > 0
}

// Running 行驶状态设置
func (s *StatusFlag) Running(truth bool) {
	if truth {
		*s |= 4194304
	} else {
		*s &= 0xFFDFFFFF
	}
}

// Signal 车辆信号扩展状态位
type Signal uint32

// GetHeadLamp 获取前大灯状态,1-近光灯，2-远光灯
func (s *Signal) GetHeadLamp() uint32 {
	return uint32(*s & 3)
}

// SetHeadLamp 设置前大灯状态
func (s *Signal) SetHeadLamp(lamp uint32) {
	*s |= Signal(lamp & 3)
}

// GetTurnSignal 获取转向灯状态,1-右转灯，2-左转灯
func (s *Signal) GetTurnSignal() uint32 {
	return (uint32(*s) >> 2) & 3
}

// SetTurnSignal 设置转向灯状态
func (s *Signal) SetTurnSignal(signal uint32) {
	*s |= Signal(signal&3) << 2
}

// IsBrake 是否制动
func (s *Signal) IsBrake() bool {
	return *s&16 > 0
}

// Brake 制动
func (s *Signal) Brake(truth bool) {
	if truth {
		*s |= 16
	} else {
		*s &= 0xFFFFFFEF
	}
}

// IsReverseGear 是否倒挡
func (s *Signal) IsReverseGear() bool {
	return *s&32 > 0
}

// ReverseGear 倒挡
func (s *Signal) ReverseGear(truth bool) {
	if truth {
		*s |= 32
	} else {
		*s &= 0xFFFFFFDF
	}
}

// IsFogLamp 是否雾灯
func (s *Signal) IsFogLamp() bool {
	return *s&64 > 0
}

// FogLamp 雾灯
func (s *Signal) FogLamp(truth bool) {
	if truth {
		*s |= 64
	} else {
		*s &= 0xFFFFFFBF
	}
}

// IsMarkerLamp 是否示廓灯
func (s *Signal) IsMarkerLamp() bool {
	return *s&128 > 0
}

// MarkerLamp 示廓灯
func (s *Signal) MarkerLamp(truth bool) {
	if truth {
		*s |= 128
	} else {
		*s &= 0xFFFFFF7F
	}
}

// IsBlow 是否鸣笛
func (s *Signal) IsBlow() bool {
	return *s&256 > 0
}

// Blow 鸣笛
func (s *Signal) Blow(truth bool) {
	if truth {
		*s |= 256
	} else {
		*s &= 0xFFFFFEFF
	}
}

// IsAirCondOpened 是否空调已打开
func (s *Signal) IsAirCondOpened() bool {
	return *s&512 > 0
}

// AirCondOpened 空调已打开
func (s *Signal) AirCondOpened(truth bool) {
	if truth {
		*s |= 512
	} else {
		*s &= 0xFFFFFDFF
	}
}

// IsNeutralGear 是否空挡
func (s *Signal) IsNeutralGear() bool {
	return *s&1024 > 0
}

// NeutralGear 空挡
func (s *Signal) NeutralGear(truth bool) {
	if truth {
		*s |= 1024
	} else {
		*s &= 0xFFFFFBFF
	}
}

// IsRetarderWorking 是否缓速器已工作
func (s *Signal) IsRetarderWorking() bool {
	return *s&2048 > 0
}

// RetarderWorking 缓速器已工作
func (s *Signal) RetarderWorking(truth bool) {
	if truth {
		*s |= 2048
	} else {
		*s &= 0xFFFFF7FF
	}
}

// IsABSWorking 是否ABS已工作
func (s *Signal) IsABSWorking() bool {
	return *s&4096 > 0
}

// ABSWorking ABS已工作
func (s *Signal) ABSWorking(truth bool) {
	if truth {
		*s |= 4096
	} else {
		*s &= 0xFFFFEFFF
	}
}

// IsHeaterWorking 是否加热器已工作
func (s *Signal) IsHeaterWorking() bool {
	return *s&8192 > 0
}

// HeaterWorking 加热器已工作
func (s *Signal) HeaterWorking(truth bool) {
	if truth {
		*s |= 8192
	} else {
		*s &= 0xFFFFDFFF
	}
}

// IsClutchReleased 是否离合器已松开
func (s *Signal) IsClutchReleased() bool {
	return *s&16384 > 0
}

// ClutchReleased 离合器已松开
func (s *Signal) ClutchReleased(truth bool) {
	if truth {
		*s |= 16384
	} else {
		*s &= 0xFFFFBFFF
	}
}

// PosAuxMap 位置附加信息对照表
var PosAuxMap = map[byte]string{
	byte(0x01): "uint32",       // 里程，单位为1/10km，对应车上的里程表读数
	byte(0x02): "uint16",       // 油量，单位为1/10L，对应车上油量表读数
	byte(0x03): "uint16",       // 行驶记录功能获取的速度，单位为1/10km/h
	byte(0x04): "uint16",       // 需要人工确认报警事件的ID，从1开始计数
	byte(0x05): "[]uint8",      // 胎压，单位为Pa，标定轮子的顺序为从车头开始从左到右顺序排列，多余的字节为0xFF，表示无效数据
	byte(0x06): "int16",        // 车厢温度，单位为摄氏度，取值范围为-32767~+32767，最高位为1表示负数
	byte(0x07): "int16",        // 冷链货仓温度（传感器1），单位为1/10摄氏度，取值范围为-32768~+32767，最高位为1表示负数
	byte(0x08): "int16",        // 冷链货仓温度（传感器2），单位为1/10摄氏度，取值范围为-32768~+32767，最高位为1表示负数
	byte(0x09): "uint8",        // 冷链货仓湿度1，单位为百分比，取值范围为0-100
	byte(0x0A): "uint8",        // 冷链货仓湿度2，单位为百分比，取值范围为0-100
	byte(0x0B): "uint8",        // 附加的自定义故障报警，包括门磁和温湿度传感器
	byte(0x0C): "uint8",        // 门磁1（未安装时不发送此ID）0：门磁1解锁；1：门磁1加锁（货仓后门）
	byte(0x0D): "uint8",        // 门磁2（未安装时不发送此ID）0：门磁2解锁；1：门磁2加锁（货仓侧门）
	byte(0x11): "SpeedAlarm",   // 超速报警附加信息
	byte(0x12): "LocalAlarm",   // 进出区域/路线报警附加信息
	byte(0x13): "RuntimeAlarm", // 路段行驶时间不足/过长报警附加信息
	byte(0x25): "uint32",       // 扩展车辆信号状态位
	byte(0x2A): "uint16",       // IO状态位
	byte(0x2B): "uint32",       // 模拟量，bit0-15，AD0；bit16-31，AD1
	byte(0x30): "uint8",        // 无线通信网络信号强度
	byte(0x31): "uint8",        // GNSS定位卫星数
	byte(0xE0): "unknown",      // 后续自定义信息长度
	byte(0xE2): "AIAlarm",      // AI报警
	byte(0x65): "DMSAlarm",     // DMS报警
	byte(0x66): "TPMSAlarm",    // 轮胎状态监测报警
}

// PositionAux 位置附加信息项
type PositionAux struct {
	// 附加项id
	ID byte
	// 附加项长度
	Len byte
	// 附加项内容
	Value interface{}
}

// ReadBy 从缓存中度
func (p *PositionAux) ReadBy(buf *bytes.Buffer) {
	// 附加项id
	p.ID, _ = buf.ReadByte()
	// 附加项长度
	p.Len, _ = buf.ReadByte()
	// 附加项内容
	switch PosAuxMap[p.ID] {
	case "uint8":
		p.Value, _ = buf.ReadByte()
	case "int16":
		bts := make([]byte, 2)
		buf.Read(bts)
		p.Value = int16(binary.BigEndian.Uint16(bts))
	case "uint16":
		bts := make([]byte, 2)
		buf.Read(bts)
		p.Value = binary.BigEndian.Uint16(bts)
	case "uint32":
		bts := make([]byte, 4)
		buf.Read(bts)
		p.Value = binary.BigEndian.Uint32(bts)
	case "[]uint8":
		bts := make([]byte, p.Len)
		buf.Read(bts)
		p.Value = bts
	case "string":
		decoder := mahonia.NewDecoder("gbk")
		bts := make([]byte, p.Len)
		buf.Read(bts)
		p.Value = decoder.ConvertString(string(bts))
	case "SpeedAlarm":
		var alarm SpeedAlarm
		alarm.readBy(buf)
		p.Value = &alarm
	case "LocalAlarm":
		var alarm LocalAlarm
		alarm.readBy(buf)
		p.Value = &alarm
	case "RuntimeAlarm":
		var alarm RuntimeAlarm
		alarm.readBy(buf)
		p.Value = &alarm
	case "DMSAlarm":
		var alarm DMSAlarm
		alarm.readBy(buf)
		p.Value = &alarm
	case "TPMSAlarm":
		// todo
	default:
	}
}

// WriteTo 写入缓存中
func (p *PositionAux) WriteTo(buf *bytes.Buffer) {
	// 附加项id
	buf.WriteByte(p.ID)
	// 附加项长度和内容
	switch p.Value.(type) {
	case int8:
		buf.WriteByte(byte(1))
		buf.WriteByte(byte(p.Value.(int8)))
	case uint8:
		buf.WriteByte(byte(1))
		buf.WriteByte(p.Value.(byte))
	case int16:
		buf.WriteByte(byte(2))
		bts := make([]byte, 2)
		binary.BigEndian.PutUint16(bts, uint16(p.Value.(int16)))
		buf.Write(bts)
	case uint16:
		buf.WriteByte(byte(2))
		bts := make([]byte, 2)
		binary.BigEndian.PutUint16(bts, p.Value.(uint16))
		buf.Write(bts)
	case uint32:
		buf.WriteByte(byte(4))
		bts := make([]byte, 4)
		binary.BigEndian.PutUint32(bts, p.Value.(uint32))
		buf.Write(bts)
	case []uint8:
		bts := p.Value.([]byte)
		buf.WriteByte(byte(len(bts)))
		buf.Write(bts)
	case string:
		encoder := mahonia.NewEncoder("gbk")
		bts := []byte(encoder.ConvertString(p.Value.(string)))
		buf.WriteByte(byte(len(bts)))
		buf.Write(bts)
	case SpeedAlarm:
		alarm := p.Value.(SpeedAlarm)
		buf.WriteByte(alarm.size())
		alarm.writeTo(buf)
	case LocalAlarm:
		alarm := p.Value.(LocalAlarm)
		buf.WriteByte(alarm.size())
		alarm.writeTo(buf)
	case RuntimeAlarm:
		alarm := p.Value.(RuntimeAlarm)
		buf.WriteByte(alarm.size())
		alarm.writeTo(buf)
	case DMSAlarm:
		alarm := p.Value.(DMSAlarm)
		buf.WriteByte(alarm.size())
		alarm.writeTo(buf)
	case TPMSAlarm:
		// todo
	default:
	}
}

// 载货状态
const (
	CargoStatusEmpty uint32 = 0 // 空载
	CargoStatusHalf  uint32 = 1 // 半载
	CargoStatusFull  uint32 = 3 // 满载
)

// Position 位置信息
type Position struct {
	// 数据主键
	PK int64 `orm:"column(id);pk;auto"`
	// 报警标识
	Alarm AlarmFlag
	// 状态标识
	Status StatusFlag
	// 纬度，单位0.000001°
	Latitude uint32
	// 经度，单位0.000001°
	Longitude uint32
	// 高程，单位m
	Altitude uint16
	// 速度，单位0.1km/h
	Speed uint16
	// 方向，单位°
	Bearing uint16
	// 时间，单位s，GMT+8
	Time time.Time `orm:"type(datetime)"`
	// 附加信息项列表
	Auxs map[byte]PositionAux `orm:"-"`
	// 附加信息项字节数组
	AuxsBts string

	//sync.Mutex // add by zhangxinyi0811,otherwise it will cause [pos.Auxs] rw failed
}

// ReadBy 从缓存中读
func (p *Position) ReadBy(buf *bytes.Buffer) {
	value := make([]byte, 6)

	// 报警标志
	buf.Read(value[:4])
	p.Alarm = AlarmFlag(binary.BigEndian.Uint32(value))
	// 状态
	buf.Read(value[:4])
	p.Status = StatusFlag(binary.BigEndian.Uint32(value))
	// 纬度
	buf.Read(value[:4])
	p.Latitude = binary.BigEndian.Uint32(value)
	// 经度
	buf.Read(value[:4])
	p.Longitude = binary.BigEndian.Uint32(value)
	// 高程
	buf.Read(value[:2])
	p.Altitude = binary.BigEndian.Uint16(value)
	// 速度
	buf.Read(value[:2])
	p.Speed = binary.BigEndian.Uint16(value)
	// 方向
	buf.Read(value[:2])
	p.Bearing = binary.BigEndian.Uint16(value)
	// 时间
	buf.Read(value)
	p.Time, _ = time.ParseInLocation("060102150405", string(util.ParseBCD(value)), time.FixedZone("CST", 28800))
	// 附加项
	if buf.Len() > 0 && nil == p.Auxs {
		p.Auxs = make(map[byte]PositionAux)
	}
	var aux PositionAux
	for buf.Len() > 0 {
		aux.ReadBy(buf)
		p.Auxs[aux.ID] = aux
		// p.Auxs = append(p.Auxs, aux)
	}
}

// ReadBaseBy 从缓存中读位置基本信息
func (p *Position) ReadBaseBy(buf *bytes.Buffer) {
	// 报警标志
	value := binary.BigEndian.Uint32(buf.Next(4))
	p.Alarm = AlarmFlag(value)
	// 状态
	value = binary.BigEndian.Uint32(buf.Next(4))
	p.Status = StatusFlag(value)
	// 纬度
	p.Latitude = binary.BigEndian.Uint32(buf.Next(4))
	// 经度
	p.Longitude = binary.BigEndian.Uint32(buf.Next(4))
	// 高程
	p.Altitude = binary.BigEndian.Uint16(buf.Next(2))
	// 速度
	p.Speed = binary.BigEndian.Uint16(buf.Next(2))
	// 方向
	p.Bearing = binary.BigEndian.Uint16(buf.Next(2))
	// 时间
	p.Time, _ = time.ParseInLocation("060102150405", string(util.ParseBCD(buf.Next(6))), time.FixedZone("CST", 28800))
}

// WriteTo 写入缓存中
func (p *Position) WriteTo(buf *bytes.Buffer) {
	value := make([]byte, 4)

	// 报警标志
	binary.BigEndian.PutUint32(value, uint32(p.Alarm))
	buf.Write(value)
	// 状态
	binary.BigEndian.PutUint32(value, uint32(p.Status))
	buf.Write(value)
	// 纬度
	binary.BigEndian.PutUint32(value, p.Latitude)
	buf.Write(value)
	// 经度
	binary.BigEndian.PutUint32(value, p.Longitude)
	buf.Write(value)
	// 高程
	binary.BigEndian.PutUint16(value, p.Altitude)
	buf.Write(value[:2])
	// 速度
	binary.BigEndian.PutUint16(value, p.Speed)
	buf.Write(value[:2])
	// 方向
	binary.BigEndian.PutUint16(value, p.Bearing)
	buf.Write(value[:2])
	// 时间
	buf.Write(util.ToBCD([]byte(p.Time.Format("060102150405"))))
	// 位置附加项
	for _, aux := range p.Auxs {
		aux.WriteTo(buf)
	}
}

// WriteBaseTo 将位置基本信息写入缓存中
func (p *Position) WriteBaseTo(buf *bytes.Buffer) {
	value := make([]byte, 4)

	// 报警标志
	binary.BigEndian.PutUint32(value, uint32(p.Alarm))
	buf.Write(value)
	// 状态
	binary.BigEndian.PutUint32(value, uint32(p.Status))
	buf.Write(value)
	// 纬度
	binary.BigEndian.PutUint32(value, p.Latitude)
	buf.Write(value)
	// 经度
	binary.BigEndian.PutUint32(value, p.Longitude)
	buf.Write(value)
	// 高程
	binary.BigEndian.PutUint16(value, p.Altitude)
	buf.Write(value[:2])
	// 速度
	binary.BigEndian.PutUint16(value, p.Speed)
	buf.Write(value[:2])
	// 方向
	binary.BigEndian.PutUint16(value, p.Bearing)
	buf.Write(value[:2])
	// 时间
	buf.Write(util.ToBCD([]byte(p.Time.Format("060102150405"))))
}

// WriteLocTo 将位置定位信息写入缓存中
func (p *Position) WriteLocTo(buf *bytes.Buffer) {
	value := make([]byte, 4)

	// 纬度
	binary.BigEndian.PutUint32(value, p.Latitude)
	buf.Write(value)
	// 经度
	binary.BigEndian.PutUint32(value, p.Longitude)
	buf.Write(value)
	// 高程
	binary.BigEndian.PutUint16(value, p.Altitude)
	buf.Write(value[:2])
	// 速度
	binary.BigEndian.PutUint16(value, p.Speed)
	buf.Write(value[:2])
	// 方向
	binary.BigEndian.PutUint16(value, p.Bearing)
	buf.Write(value[:2])
	// 时间
	buf.Write(util.ToBCD([]byte(p.Time.Format("060102150405"))))
}

// Bytes 序列化
func (p *Position) Bytes() []byte {
	var buf bytes.Buffer
	p.WriteTo(&buf)
	return buf.Bytes()
}

// BaseBytes 基础信息序列化
func (p *Position) BaseBytes() []byte {
	var buf bytes.Buffer
	p.WriteBaseTo(&buf)
	return buf.Bytes()
}

//var mtxPos sync.Mutex
// Duplicate 复制定位信息
func (p *Position) Duplicate() *Position {
	//p.Lock()
	//defer p.Unlock() //为什么在这里上锁仍然要报Auxs读写错误

	pos := new(Position)
	pos.Alarm = p.Alarm
	pos.Status = p.Status
	pos.Latitude = p.Latitude
	pos.Longitude = p.Longitude
	pos.Altitude = p.Altitude
	pos.Speed = p.Speed
	pos.Bearing = p.Bearing
	pos.Time = p.Time
	pos.Auxs = make(map[byte]PositionAux)
	if len(p.Auxs) > 0 {
		//pos.Auxs = map[byte]PositionAux{}
		for key, val := range p.Auxs {
			pos.Auxs[key] = val
		}
	}
	return pos
}

// Satellite 卫星状态
type Satellite struct {
	// 卫星id
	ID byte `json:"id"`
	// 俯仰角
	Elevation byte `json:"elevation"`
	// 方位角
	Azimuth uint16 `json:"azimuth"`
}

// 写入缓存中
func (s *Satellite) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 2)

	// 卫星id
	buf.WriteByte(s.ID)
	// 俯仰角
	buf.WriteByte(s.Elevation)
	// 方位角
	binary.BigEndian.PutUint16(value, s.Azimuth)
	buf.Write(value)
}

// 从缓存中读
func (s *Satellite) readBy(buf *bytes.Buffer) {
	value := make([]byte, 2)
	// 卫星id
	s.ID, _ = buf.ReadByte()
	// 俯仰角
	s.Elevation, _ = buf.ReadByte()
	// 方位角
	buf.Read(value)
	s.Azimuth = binary.BigEndian.Uint16(value)
}

// GnnsStatus 定位状态信息
type GnnsStatus struct {
	// 北斗卫星状态
	BDSates []Satellite `json:"bd_status"`
	// GPS卫星状态
	GPSSates []Satellite `json:"gps_status"`
	// 格勒纳斯卫星状态
	GLSates []Satellite
	// 伽利略卫星状态
	GOSates []Satellite
}

// WriteTo 写入缓存中
func (g *GnnsStatus) WriteTo(buf *bytes.Buffer) {
	// 北斗卫星个数
	buf.WriteByte(byte(len(g.BDSates)))
	// 北斗卫星状态
	for _, sate := range g.BDSates {
		sate.writeTo(buf)
	}
	// GPS卫星个数
	buf.WriteByte(byte(len(g.GPSSates)))
	// GPS卫星状态
	for _, sate := range g.GPSSates {
		sate.writeTo(buf)
	}
	// 格勒纳斯卫星个数
	buf.WriteByte(byte(len(g.GLSates)))
	// 格勒纳斯卫星状态
	for _, sate := range g.GLSates {
		sate.writeTo(buf)
	}
	// 伽利略卫星个数
	buf.WriteByte(byte(len(g.GOSates)))
	// 伽利略卫星状态
	for _, sate := range g.GOSates {
		sate.writeTo(buf)
	}
}

// ReadBy 从缓存中读
func (g *GnnsStatus) ReadBy(buf *bytes.Buffer) {
	// 北斗卫星个数
	bdCount, _ := buf.ReadByte()
	// 北斗卫星状态
	var bdStatus []Satellite
	for i := 0; i < int(bdCount); i++ {
		var sate Satellite
		sate.readBy(buf)
		bdStatus = append(bdStatus, sate)
	}
	// GPS卫星个数
	gpsCount, _ := buf.ReadByte()
	// GPS卫星状态
	var gpsStatus []Satellite
	for i := 0; i < int(gpsCount); i++ {
		var sate Satellite
		sate.readBy(buf)
		gpsStatus = append(gpsStatus, sate)
	}
	// 格勒纳斯卫星个数
	glCount, _ := buf.ReadByte()
	// 格勒纳斯卫星状态
	var glStatus []Satellite
	for i := 0; i < int(glCount); i++ {
		var sate Satellite
		sate.readBy(buf)
		glStatus = append(glStatus, sate)
	}
	// 伽利略卫星个数
	goCount, _ := buf.ReadByte()
	// 伽利略卫星状态
	var goStatus []Satellite
	for i := 0; i < int(goCount); i++ {
		var sate Satellite
		sate.readBy(buf)
		goStatus = append(goStatus, sate)
	}
	g.BDSates = bdStatus
	g.GPSSates = gpsStatus
	g.GLSates = glStatus
	g.GOSates = goStatus
}

// Duplicate 新建副本
func (g *GnnsStatus) Duplicate() *GnnsStatus {
	gnss := new(GnnsStatus)
	gnss.BDSates = g.BDSates
	gnss.GPSSates = g.GPSSates
	return gnss
}

// AlarmConfirm 报警确认信息
type AlarmConfirm uint32

// ConfirmEmergency 是否确认紧急报警
func (a *AlarmConfirm) ConfirmEmergency() bool {
	return *a&1 > 0
}

// ConfirmDanger 确认危险报警
func (a *AlarmConfirm) ConfirmDanger() bool {
	return *a&8 > 0
}

// ConfirmInOutArea 确认进出区域报警
func (a *AlarmConfirm) ConfirmInOutArea() bool {
	return *a&1048576 > 0
}

// ConfirmInOutPath 确认进出路线报警
func (a *AlarmConfirm) ConfirmInOutPath() bool {
	return *a&2097152 > 0
}

// ConfirmRoadTime 确认路段时间不符报警
func (a *AlarmConfirm) ConfirmRoadTime() bool {
	return *a&4194304 > 0
}

// ConfirmIgnition 确认非法点火报警
func (a *AlarmConfirm) ConfirmIgnition() bool {
	return *a&134217728 > 0
}

// ConfirmMove 确认非法位移报警
func (a *AlarmConfirm) ConfirmMove() bool {
	return *a&268435456 > 0
}

// TextFlag 文本标志信息
type TextFlag byte

// Type 文本信息类型，0-服务，2-紧急，3-通知
func (t *TextFlag) Type() byte {
	return byte(*t & 0x03)
}

// IsShow 是否终端显示器显示
func (t *TextFlag) IsShow() bool {
	return *t&4 > 0
}

// IsTTS 是否TTS语音播报
func (t *TextFlag) IsTTS() bool {
	return *t&8 > 0
}

// Kind 0-中心导航信息，1-CAN故障码信息
func (t *TextFlag) Kind() byte {
	return byte(*t & 32)
}

// Contact 联系人项
type Contact struct {
	// 标识，1-呼入，2-呼出，3-呼入/呼出
	Flag byte
	// 电话号码
	Phone string
	// 联系人姓名
	Name string
}

// VehCtrlParam 车辆控制参数
type VehCtrlParam struct {
	// ID
	ID uint16
	// 参数值
	Value interface{}
}

// Area 区域接口
type Area interface {
	readBy(buf *bytes.Buffer, decoder mahonia.Decoder, version byte)
	writeTo(buf *bytes.Buffer, encoder mahonia.Encoder, version byte)
}

// AreaAttr 区域属性
type AreaAttr uint16

// HasTime 是否有起始/结束时间判断
func (a *AreaAttr) HasTime() bool {
	return *a&AreaAttr(1) > 0
}

// SetTimeEnable 设置时间有效性
func (a *AreaAttr) SetTimeEnable(enable bool) {
	if enable {
		*a |= 1
	} else {
		*a &= 0xFFFE
	}
}

// HasSpeed 是否有速度值
func (a *AreaAttr) HasSpeed() bool {
	return *a&2 > 0
}

// SetSpeedEnable 设置速度有效性
func (a *AreaAttr) SetSpeedEnable(enable bool) {
	if enable {
		*a |= 2
	} else {
		*a &= 0xFFFD
	}
}

// HasEnterAlarmToDriver 进区域是否通知驾驶员
func (a *AreaAttr) HasEnterAlarmToDriver() bool {
	return *a&4 > 0
}

// SetEnterAlarmToDriver 设置进区域是否通知驾驶员
func (a *AreaAttr) SetEnterAlarmToDriver(enable bool) {
	if enable {
		*a |= 4
	} else {
		*a &= 0xFFFB
	}
}

// HasEnterAlarmToServer 进区域是否通知平台
func (a *AreaAttr) HasEnterAlarmToServer() bool {
	return *a&8 > 0
}

// SetEnterAlarmToServer 设置进区域是否通知平台
func (a *AreaAttr) SetEnterAlarmToServer(enable bool) {
	if enable {
		*a |= 8
	} else {
		*a &= 0xFFF7
	}
}

// HasExitAlarmToDriver 出区域是否通知驾驶员
func (a *AreaAttr) HasExitAlarmToDriver() bool {
	return *a&16 > 0
}

// SetExitAlarmToDriver 设置出区域是否通知驾驶员
func (a *AreaAttr) SetExitAlarmToDriver(enable bool) {
	if enable {
		*a |= 16
	} else {
		*a &= 0xFFEF
	}
}

// HasExitAlarmToServer 出区域是否通知平台
func (a *AreaAttr) HasExitAlarmToServer() bool {
	return *a&32 > 0
}

// SetExitAlarmToServer 设置出区域是否通知平台
func (a *AreaAttr) SetExitAlarmToServer(enable bool) {
	if enable {
		*a |= 32
	} else {
		*a &= 0xFFDF
	}
}

// IsSouth 南纬或北纬
func (a *AreaAttr) IsSouth() bool {
	return *a&64 > 0
}

// SetSouth 设置南纬或北纬
func (a *AreaAttr) SetSouth(enable bool) {
	if enable {
		*a |= 64
	} else {
		*a &= 0xFFBF
	}
}

// IsWest 西经或东经
func (a *AreaAttr) IsWest() bool {
	return *a&128 > 0
}

// SetWest 设置西经或东经
func (a *AreaAttr) SetWest(enable bool) {
	if enable {
		*a |= 128
	} else {
		*a &= 0xFF7F
	}
}

// IsDoorOpenable 是否允许开门
func (a *AreaAttr) IsDoorOpenable() bool {
	return *a&256 == 0
}

// SetDoorOpenable 设置车门可开
func (a *AreaAttr) SetDoorOpenable(enable bool) {
	if enable {
		*a &= 0xFEFF
	} else {
		*a |= 256
	}
}

// IsCloseCommWithEnter 进区域是否关闭通信模块
func (a *AreaAttr) IsCloseCommWithEnter() bool {
	return *a&16384 > 0
}

// SetCloseCommWithEnter 设置进区域是否关闭通信模块
func (a *AreaAttr) SetCloseCommWithEnter(enable bool) {
	if enable {
		*a |= 16384
	} else {
		*a &= 0xEFFF
	}
}

// IsGatherGNSSWithEnter 进区域是否采集GNSS详细定位数据
func (a *AreaAttr) IsGatherGNSSWithEnter() bool {
	return *a&32768 > 0
}

// SetGatherGNSSWithEnter 设置进区域是否采集GNSS详细定位数据
func (a *AreaAttr) SetGatherGNSSWithEnter(enable bool) {
	if enable {
		*a |= 32768
	} else {
		*a &= 0x7FFF
	}
}

// RoundArea 圆形区域项
type RoundArea struct {
	// 区域ID
	ID uint32
	// 区域属性
	Attr AreaAttr
	// 中心点纬度
	CenterY uint32
	// 中心点经度
	CenterX uint32
	// 半径
	Radius uint32
	// 起始时间
	STime string //time.Time
	// 结束时间
	ETime string //time.Time
	// 最高速度，单位km/h
	MaxSpeed uint16
	// 超速持续时间，单位s
	Duration byte
	// 夜间最高速度，单位km/h
	MaxSpeedInNight uint16
	// 区域名称
	Name string
}

// readBy 从缓存中读
func (r *RoundArea) readBy(buf *bytes.Buffer, decoder mahonia.Decoder, version byte) {
	value := make([]byte, 6)

	buf.Read(value[:4])
	r.ID = binary.BigEndian.Uint32(value)
	buf.Read(value[:2])
	r.Attr = AreaAttr(binary.BigEndian.Uint16(value))
	buf.Read(value[:4])
	r.CenterY = binary.BigEndian.Uint32(value)
	buf.Read(value[:4])
	r.CenterX = binary.BigEndian.Uint32(value)
	buf.Read(value[:4])
	r.Radius = binary.BigEndian.Uint32(value)

	if r.Attr.HasTime() {
		buf.Read(value)
		r.STime = string(util.ParseBCD(value))
		// r.STime, _ = time.ParseInLocation("060102150405", string(util.ParseBCD(value)), time.FixedZone("CST", 28800))
		buf.Read(value)
		r.ETime = string(util.ParseBCD(value))
		// r.ETime, _ = time.ParseInLocation("060102150405", string(util.ParseBCD(value)), time.FixedZone("CST", 28800))
	}

	if version2011 == version {
		if r.Attr.HasSpeed() {
			buf.Read(value[:2])
			r.MaxSpeed = binary.BigEndian.Uint16(value)
			r.Duration, _ = buf.ReadByte()
			r.MaxSpeedInNight = 0xFFFF
		}
		r.Name = ""
	} else {
		if r.Attr.HasSpeed() {
			buf.Read(value[:2])
			r.MaxSpeed = binary.BigEndian.Uint16(value)
			r.Duration, _ = buf.ReadByte()
			buf.Read(value[:2])
			r.MaxSpeedInNight = binary.BigEndian.Uint16(value)
		}

		buf.Read(value[:2])
		len := binary.BigEndian.Uint16(value)
		bts := make([]byte, len)
		buf.Read(bts)
		r.Name = decoder.ConvertString(string(bts))
	}
}

func (r *RoundArea) writeTo(buf *bytes.Buffer, encoder mahonia.Encoder, version byte) {
	value := make([]byte, 6)

	binary.BigEndian.PutUint32(value, r.ID)
	buf.Write(value[:4])
	binary.BigEndian.PutUint16(value, uint16(r.Attr))
	buf.Write(value[:2])
	binary.BigEndian.PutUint32(value, r.CenterY)
	buf.Write(value[:4])
	binary.BigEndian.PutUint32(value, r.CenterX)
	buf.Write(value[:4])
	binary.BigEndian.PutUint32(value, r.Radius)
	buf.Write(value[:4])

	if r.Attr.HasTime() {
		// bts := util.ToBCD([]byte(r.STime.Format("060102150405")))
		bts := util.ToBCD([]byte(r.STime))
		buf.Write(bts)
		// bts = util.ToBCD([]byte(r.ETime.Format("060102150405")))
		bts = util.ToBCD([]byte(r.ETime))
		buf.Write(bts)
	}

	if r.Attr.HasSpeed() {
		binary.BigEndian.PutUint16(value, r.MaxSpeed)
		buf.Write(value[:2])
		buf.WriteByte(r.Duration)
		binary.BigEndian.PutUint16(value, r.MaxSpeedInNight)
		buf.Write(value[:2])
	}

	bts := []byte(encoder.ConvertString(r.Name))
	len := uint16(len(bts))
	binary.BigEndian.PutUint16(value, len)
	buf.Write(value[:2])
	buf.Write(bts)
}

// RectArea 矩形区域项
type RectArea struct {
	// 区域ID
	ID uint32
	// 区域属性
	Attr AreaAttr
	// 左上纬度（西北）
	NWY uint32
	// 左上经度（西北）
	NWX uint32
	// 右下纬度（东南）
	SEY uint32
	// 右下经度（东南）
	SEX uint32
	// 起始时间
	STime string //time.Time
	// 结束时间
	ETime string //time.Time
	// 最高速度，单位km/h
	MaxSpeed uint16
	// 超速持续时间，单位s
	Duration byte
	// 夜间最高速度，单位km/h
	MaxSpeedInNight uint16
	// 区域名称
	Name string
}

// readBy 从缓存中读
func (r *RectArea) readBy(buf *bytes.Buffer, decoder mahonia.Decoder, version byte) {
	value := make([]byte, 6)

	buf.Read(value[:4])
	r.ID = binary.BigEndian.Uint32(value)
	buf.Read(value[:2])
	r.Attr = AreaAttr(binary.BigEndian.Uint16(value))
	buf.Read(value[:4])
	r.NWY = binary.BigEndian.Uint32(value)
	buf.Read(value[:4])
	r.NWX = binary.BigEndian.Uint32(value)
	buf.Read(value[:4])
	r.SEY = binary.BigEndian.Uint32(value)
	buf.Read(value[:4])
	r.SEX = binary.BigEndian.Uint32(value)

	if r.Attr.HasTime() {
		buf.Read(value)
		r.STime = string(util.ParseBCD(value))
		// r.STime, _ = time.ParseInLocation("060102150405", string(util.ParseBCD(value)), time.FixedZone("CST", 28800))
		buf.Read(value)
		r.ETime = string(util.ParseBCD(value))
		// r.ETime, _ = time.ParseInLocation("060102150405", string(util.ParseBCD(value)), time.FixedZone("CST", 28800))
	}

	if version2011 == version {
		if r.Attr.HasSpeed() {
			buf.Read(value[:2])
			r.MaxSpeed = binary.BigEndian.Uint16(value)
			r.Duration, _ = buf.ReadByte()
			r.MaxSpeedInNight = 0xFFFF
		}
		r.Name = ""
	} else {
		if r.Attr.HasSpeed() {
			buf.Read(value[:2])
			r.MaxSpeed = binary.BigEndian.Uint16(value)
			r.Duration, _ = buf.ReadByte()
			buf.Read(value[:2])
			r.MaxSpeedInNight = binary.BigEndian.Uint16(value)
		}

		buf.Read(value[:2])
		len := binary.BigEndian.Uint16(value)
		bts := make([]byte, len)
		buf.Read(bts)
		r.Name = decoder.ConvertString(string(bts))
	}
}

// 写入缓存中
func (r *RectArea) writeTo(buf *bytes.Buffer, encoder mahonia.Encoder, version byte) {
	value := make([]byte, 6)

	binary.BigEndian.PutUint32(value, r.ID)
	buf.Write(value[:4])
	binary.BigEndian.PutUint16(value, uint16(r.Attr))
	buf.Write(value[:2])
	binary.BigEndian.PutUint32(value, r.NWY)
	buf.Write(value[:4])
	binary.BigEndian.PutUint32(value, r.NWX)
	buf.Write(value[:4])
	binary.BigEndian.PutUint32(value, r.SEY)
	buf.Write(value[:4])
	binary.BigEndian.PutUint32(value, r.SEX)
	buf.Write(value[:4])

	if r.Attr.HasTime() {
		bts := util.ToBCD([]byte(r.STime))
		// bts := util.ToBCD([]byte(r.STime.Format("060102150405")))
		buf.Write(bts)
		bts = util.ToBCD([]byte(r.ETime))
		// bts = util.ToBCD([]byte(r.ETime.Format("060102150405")))
		buf.Write(bts)
	}

	if r.Attr.HasSpeed() {
		binary.BigEndian.PutUint16(value, r.MaxSpeed)
		buf.Write(value[:2])
		buf.WriteByte(r.Duration)
		binary.BigEndian.PutUint16(value, r.MaxSpeedInNight)
		buf.Write(value[:2])
	}

	bts := []byte(encoder.ConvertString(r.Name))
	len := uint16(len(bts))
	binary.BigEndian.PutUint16(value, len)
	buf.Write(value[:2])
	buf.Write(bts)
}

// PolygonArea 多边形区域项
type PolygonArea struct {
	// 区域ID
	ID uint32
	// 区域属性
	Attr AreaAttr
	// 顶点列表
	Vertexs []uint32
	// 起始时间
	STime string //time.Time
	// 结束时间
	ETime string //time.Time
	// 最高速度，单位km/h
	MaxSpeed uint16
	// 超速持续时间，单位s
	Duration byte
	// 夜间最高速度，单位km/h
	MaxSpeedInNight uint16
	// 区域名称
	Name string
}

// readBy 从缓存中读
func (p *PolygonArea) readBy(buf *bytes.Buffer, decoder mahonia.Decoder, version byte) {
	value := make([]byte, 6)

	buf.Read(value[:4])
	p.ID = binary.BigEndian.Uint32(value)
	buf.Read(value[:2])
	p.Attr = AreaAttr(binary.BigEndian.Uint16(value))

	if p.Attr.HasTime() {
		buf.Read(value)
		p.STime = string(util.ParseBCD(value))
		// p.STime, _ = time.ParseInLocation("060102150405", string(util.ParseBCD(value)), time.FixedZone("CST", 28800))
		buf.Read(value)
		p.ETime = string(util.ParseBCD(value))
		// p.ETime, _ = time.ParseInLocation("060102150405", string(util.ParseBCD(value)), time.FixedZone("CST", 28800))
	}

	if p.Attr.HasSpeed() {
		buf.Read(value[:2])
		p.MaxSpeed = binary.BigEndian.Uint16(value)
		p.Duration, _ = buf.ReadByte()
	}

	buf.Read(value[:2])
	count := binary.BigEndian.Uint16(value) << 1
	p.Vertexs = make([]uint32, count)
	for i := uint16(0); i < count; i++ {
		buf.Read(value[:4])
		p.Vertexs[i] = binary.BigEndian.Uint32(value)
	}

	if version2011 == version {
		if p.Attr.HasSpeed() {
			p.MaxSpeedInNight = 0xFFFF
		}
		p.Name = ""
	} else {
		if p.Attr.HasSpeed() {
			buf.Read(value[:2])
			p.MaxSpeedInNight = binary.BigEndian.Uint16(value)
		}

		buf.Read(value[:2])
		len := binary.BigEndian.Uint16(value)
		bts := make([]byte, len)
		buf.Read(bts)
		p.Name = decoder.ConvertString(string(bts))
	}

}

// 写入缓存中
func (p *PolygonArea) writeTo(buf *bytes.Buffer, encoder mahonia.Encoder, version byte) {
	value := make([]byte, 6)

	binary.BigEndian.PutUint32(value, p.ID)
	buf.Write(value[:4])
	binary.BigEndian.PutUint16(value, uint16(p.Attr))
	buf.Write(value[:2])

	if p.Attr.HasTime() {
		bts := util.ToBCD([]byte(p.STime))
		// bts := util.ToBCD([]byte(p.STime.Format("060102150405")))
		buf.Write(bts)
		bts = util.ToBCD([]byte(p.ETime))
		// bts = util.ToBCD([]byte(p.ETime.Format("060102150405")))
		buf.Write(bts)
	}

	if p.Attr.HasSpeed() {
		binary.BigEndian.PutUint16(value, p.MaxSpeed)
		buf.Write(value[:2])
		buf.WriteByte(p.Duration)
	}

	count := uint16(len(p.Vertexs)) >> 1
	binary.BigEndian.PutUint16(value, count)
	buf.Write(value[:2])
	for _, ver := range p.Vertexs {
		binary.BigEndian.PutUint32(value, ver)
		buf.Write(value[:4])
	}

	if p.Attr.HasSpeed() {
		binary.BigEndian.PutUint16(value, p.MaxSpeedInNight)
		buf.Write(value[:2])
	}

	bts := []byte(encoder.ConvertString(p.Name))
	len := uint16(len(bts))
	binary.BigEndian.PutUint16(value, len)
	buf.Write(value[:2])
	buf.Write(bts)
}

// PolylineAttr 路线属性
type PolylineAttr uint16

// HasTime 是否有起始/结束时间判断
func (a *PolylineAttr) HasTime() bool {
	return *a&PolylineAttr(1) > 0
}

// SetTimeEnable 设置时间有效性
func (a *PolylineAttr) SetTimeEnable(enable bool) {
	if enable {
		*a |= 1
	} else {
		*a &= 0xFFFE
	}
}

// HasEnterAlarmToDriver 进路线是否通知驾驶员
func (a *PolylineAttr) HasEnterAlarmToDriver() bool {
	return *a&4 > 0
}

// SetEnterAlarmToDriver 设置进路线是否通知驾驶员
func (a *PolylineAttr) SetEnterAlarmToDriver(enable bool) {
	if enable {
		*a |= 4
	} else {
		*a &= 0xFFFB
	}
}

// HasEnterAlarmToServer 进路线是否通知平台
func (a *PolylineAttr) HasEnterAlarmToServer() bool {
	return *a&8 > 0
}

// SetEnterAlarmToServer 设置进路线是否通知平台
func (a *PolylineAttr) SetEnterAlarmToServer(enable bool) {
	if enable {
		*a |= 8
	} else {
		*a &= 0xFFF7
	}
}

// HasExitAlarmToDriver 出路线是否通知驾驶员
func (a *PolylineAttr) HasExitAlarmToDriver() bool {
	return *a&16 > 0
}

// SetExitAlarmToDriver 设置出路线是否通知驾驶员
func (a *PolylineAttr) SetExitAlarmToDriver(enable bool) {
	if enable {
		*a |= 16
	} else {
		*a &= 0xFFEF
	}
}

// HasExitAlarmToServer 出路线是否通知平台
func (a *PolylineAttr) HasExitAlarmToServer() bool {
	return *a&32 > 0
}

// SetExitAlarmToServer 设置出路线是否通知平台
func (a *PolylineAttr) SetExitAlarmToServer(enable bool) {
	if enable {
		*a |= 32
	} else {
		*a &= 0xFFDF
	}
}

// Vertex 路线拐点项
type Vertex struct {
	// 拐点id
	ID uint32
	// 纬度
	Latitude uint32
	// 经度
	Longitude uint32
	// 路段
	Tag Segment
}

// 从缓存中读
func (v *Vertex) readBy(buf *bytes.Buffer, version byte) {
	value := make([]byte, 4)

	buf.Read(value)
	v.ID = binary.BigEndian.Uint32(value)
	buf.Read(value)
	v.Tag.ID = binary.BigEndian.Uint32(value)
	buf.Read(value)
	v.Latitude = binary.BigEndian.Uint32(value)
	buf.Read(value)
	v.Longitude = binary.BigEndian.Uint32(value)
	v.Tag.Width, _ = buf.ReadByte()
	bt, _ := buf.ReadByte()
	v.Tag.Attr = SegmentAttr(bt)

	if v.Tag.Attr.HasTime() {
		buf.Read(value[:2])
		v.Tag.MaxDuration = binary.BigEndian.Uint16(value)
		buf.Read(value[:2])
		v.Tag.MinDuration = binary.BigEndian.Uint16(value)
	}

	if v.Tag.Attr.HasSpeedLimit() {
		buf.Read(value[:2])
		v.Tag.MaxSpeed = binary.BigEndian.Uint16(value)
		bt, _ = buf.ReadByte()
		v.Tag.Duration = bt

		if version2011 == version {
			v.Tag.MaxSpeedInNight = 0xFFFF
		} else {
			buf.Read(value[:2])
			v.Tag.MaxSpeedInNight = binary.BigEndian.Uint16(value)
		}
	}
}

// 写入缓存中
func (v *Vertex) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 4)

	binary.BigEndian.PutUint32(value, v.ID)
	buf.Write(value)
	binary.BigEndian.PutUint32(value, v.Tag.ID)
	buf.Write(value)
	binary.BigEndian.PutUint32(value, v.Latitude)
	buf.Write(value)
	binary.BigEndian.PutUint32(value, v.Longitude)
	buf.Write(value)
	buf.WriteByte(v.Tag.Width)
	buf.WriteByte(byte(v.Tag.Attr))

	if v.Tag.Attr.HasTime() {
		binary.BigEndian.PutUint16(value, v.Tag.MaxDuration)
		buf.Write(value[:2])
		binary.BigEndian.PutUint16(value, v.Tag.MinDuration)
		buf.Write(value[:2])
	}

	if v.Tag.Attr.HasSpeedLimit() {
		binary.BigEndian.PutUint16(value, v.Tag.MaxSpeed)
		buf.Write(value[:2])
		buf.WriteByte(v.Tag.Duration)
		binary.BigEndian.PutUint16(value, v.Tag.MaxSpeedInNight)
		buf.Write(value[:2])
	}
}

// SegmentAttr 路段属性
type SegmentAttr uint8

// HasTime 是否有行驶时间限制
func (s *SegmentAttr) HasTime() bool {
	return *s&1 > 0
}

// SetTimeEnable 设置行驶时间限制
func (s *SegmentAttr) SetTimeEnable(enable bool) {
	if enable {
		*s |= 1
	} else {
		*s &= 0xFE
	}
}

// HasSpeedLimit 是否有限速设置
func (s *SegmentAttr) HasSpeedLimit() bool {
	return *s&2 > 0
}

// SetSpeedLimitEnable 设置限速
func (s *SegmentAttr) SetSpeedLimitEnable(enable bool) {
	if enable {
		*s |= 2
	} else {
		*s &= 0xFD
	}
}

// IsSouth 南纬或北纬
func (s *SegmentAttr) IsSouth() bool {
	return *s&4 > 0
}

// SetSouth 设置南纬或北纬
func (s *SegmentAttr) SetSouth(enable bool) {
	if enable {
		*s |= 4
	} else {
		*s &= 0xFB
	}
}

// IsWest 西经或东经
func (s *SegmentAttr) IsWest() bool {
	return *s&8 > 0
}

// SetWest 设置西经或东经
func (s *SegmentAttr) SetWest(enable bool) {
	if enable {
		*s |= 8
	} else {
		*s &= 0xF7
	}
}

// Segment 路段项
type Segment struct {
	// id
	ID uint32
	// 路段宽度
	Width byte
	// 路段属性
	Attr SegmentAttr
	// 路段行驶过长阈值
	MaxDuration uint16
	// 路段行驶不足阈值
	MinDuration uint16
	// 路段最高速度
	MaxSpeed uint16
	// 路段超速持续时间
	Duration byte
	// 路段夜间最高速度
	MaxSpeedInNight uint16
}

// Polyline 路线项
type Polyline struct {
	// 区域ID
	ID uint32
	// 区域属性
	Attr PolylineAttr
	// 拐点列表
	Vertexs []Vertex
	// 起始时间
	STime string //time.Time
	// 结束时间
	ETime string //time.Time
	// 区域名称
	Name string
}

// readBy 从缓存中读
func (p *Polyline) readBy(buf *bytes.Buffer, decoder mahonia.Decoder, version byte) {
	value := make([]byte, 6)

	buf.Read(value[:4])
	p.ID = binary.BigEndian.Uint32(value)
	buf.Read(value[:2])
	p.Attr = PolylineAttr(binary.BigEndian.Uint16(value))

	if p.Attr.HasTime() {
		buf.Read(value)
		p.STime = string(util.ParseBCD(value))
		// p.STime, _ = time.ParseInLocation("060102150405", string(util.ParseBCD(value)), time.FixedZone("CST", 28800))
		buf.Read(value)
		p.ETime = string(util.ParseBCD(value))
		// p.ETime, _ = time.ParseInLocation("060102150405", string(util.ParseBCD(value)), time.FixedZone("CST", 28800))
	}

	buf.Read(value[:2])
	count := binary.BigEndian.Uint16(value)
	p.Vertexs = make([]Vertex, count)
	for i := uint16(0); i < count; i++ {
		p.Vertexs[i].readBy(buf, version)
	}

	if version2011 == version {
		p.Name = ""
		return
	}
	buf.Read(value[:2])
	len := binary.BigEndian.Uint16(value)
	bts := make([]byte, len)
	buf.Read(bts)
	p.Name = decoder.ConvertString(string(bts))
}

// 写入缓存中
func (p *Polyline) writeTo(buf *bytes.Buffer, encoder mahonia.Encoder, version byte) {
	value := make([]byte, 6)

	binary.BigEndian.PutUint32(value, p.ID)
	buf.Write(value[:4])
	binary.BigEndian.PutUint16(value, uint16(p.Attr))
	buf.Write(value[:2])

	if p.Attr.HasTime() {
		bts := util.ToBCD([]byte(p.STime))
		// bts := util.ToBCD([]byte(p.STime.Format("060102150405")))
		buf.Write(bts)
		bts = util.ToBCD([]byte(p.ETime))
		// bts = util.ToBCD([]byte(p.ETime.Format("060102150405")))
		buf.Write(bts)
	}

	count := uint16(len(p.Vertexs))
	binary.BigEndian.PutUint16(value, count)
	buf.Write(value[:2])
	for _, v := range p.Vertexs {
		v.writeTo(buf)
	}

	bts := []byte(encoder.ConvertString(p.Name))
	len := uint16(len(bts))
	binary.BigEndian.PutUint16(value, len)
	buf.Write(value[:2])
	buf.Write(bts)
}

// CanData CAN总线数据
type CanData struct {
	// CAN ID
	ID uint32
	// CAN DATA
	Data [8]byte
}

// Multimedia 多媒体信息项
type Multimedia struct {
	// 多媒体ID
	ID uint32
	// 多媒体类型，0-图像，1-音频，2-视频
	MimeType byte
	// 通道ID
	ChannelID byte
	// 事件项编码，0-平台下发指令，1-定时动作，2-抢劫报警触发，3-碰撞侧翻报警
	EventCode byte
	// 位置信息项
	Pos Position
}

func (m *Multimedia) writeTo(buf *bytes.Buffer, version byte) {
	value := make([]byte, 4)

	if version == version2019 {
		// 多媒体ID
		binary.BigEndian.PutUint32(value, m.ID)
		buf.Write(value)
	}
	// 多媒体类型
	buf.WriteByte(m.MimeType)
	// 通道ID
	buf.WriteByte(m.ChannelID)
	// 事件项编码
	buf.WriteByte(m.EventCode)
	// 位置信息项
	m.Pos.WriteTo(buf)
}

func (m *Multimedia) readBy(buf *bytes.Buffer, version byte) {
	if version == version2019 {
		// 多媒体ID
		m.ID = binary.BigEndian.Uint32(buf.Next(4))
	}
	// 多媒体类型
	m.MimeType, _ = buf.ReadByte()
	// 通道ID
	m.ChannelID, _ = buf.ReadByte()
	// 事件项编码
	m.EventCode, _ = buf.ReadByte()
	// 位置信息项
	m.Pos.ReadBy(buf)
}

// MediaAlarmFlag  多媒体报警标识
type MediaAlarmFlag uint32

// MediaResource 多媒体资源项
type MediaResource struct {
	// 逻辑通道号
	Channel byte `json:"channel"`
	// 开始时间
	STime uint32 `json:"start_time"`
	// 结束时间
	ETime uint32 `json:"end_time"`

	//// 多媒体报警标识
	//MediaAlarm MediaAlarmFlag `json:"media_alarm"`
	//// 普通报警标识
	//Alarm beans.AlarmFlag `json:"alarm"`

	Alarm uint64 `json:"alarm"`

	// 音视频类型：0-音视频，1-音频，2-视频
	MediaType byte `json:"media_type"`
	// 码流类型，1-主码流，2-子码流
	StreamType byte `json:"stream_type"`
	// 存储器类型，1-主存储器，2-灾备存储器
	StorageType byte `json:"storage_type"`
	// 文件大小（字节）
	FileSize uint32 `json:"file_size"`
}

func (m *MediaResource) writeTo(buf *bytes.Buffer) {
	value := make([]byte, 4)
	// 逻辑通道号
	buf.WriteByte(m.Channel)
	// 开始时间
	buf.Write(util.ToBCD([]byte(time.Unix(int64(m.STime), 0).In(time.FixedZone("CST", 28800)).Format("060102150405"))))
	// 结束时间
	buf.Write(util.ToBCD([]byte(time.Unix(int64(m.ETime), 0).In(time.FixedZone("CST", 28800)).Format("060102150405"))))

	//// 多媒体报警标识
	//binary.BigEndian.PutUint32(value, uint32(m.MediaAlarm))
	//buf.Write(value)
	//// 普通报警标识
	//binary.BigEndian.PutUint32(value, uint32(m.Alarm))
	//buf.Write(value)

	// 报警标识
	a := make([]byte, 8)
	binary.BigEndian.PutUint64(a, m.Alarm)
	buf.Write(a)

	// 音视频类型
	buf.WriteByte(m.MediaType)
	// 码流类型
	buf.WriteByte(m.StreamType)
	// 存储器类型
	buf.WriteByte(m.StorageType)
	// 文件大小
	binary.BigEndian.PutUint32(value, m.FileSize)
	buf.Write(value)
}

func (m *MediaResource) readBy(buf *bytes.Buffer) {
	// 逻辑通道号
	m.Channel, _ = buf.ReadByte()
	// 开始时间
	t, _ := time.ParseInLocation("060102150405", string(util.ParseBCD(buf.Next(6))), time.FixedZone("CST", 28800))
	m.STime = uint32(t.Unix())
	// 结束时间
	t, _ = time.ParseInLocation("060102150405", string(util.ParseBCD(buf.Next(6))), time.FixedZone("CST", 28800))
	m.ETime = uint32(t.Unix())

	//// 多媒体报警标识
	//binary.BigEndian.PutUint32(value, uint32(m.MediaAlarm))
	//buf.Write(value)
	//// 普通报警标识
	//binary.BigEndian.PutUint32(value, uint32(m.Alarm))
	//buf.Write(value)

	// 报警标识
	m.Alarm = binary.BigEndian.Uint64(buf.Next(8))

	// 音视频类型
	m.MediaType, _ = buf.ReadByte()
	// 码流类型
	m.StreamType, _ = buf.ReadByte()
	// 存储器类型
	m.StorageType, _ = buf.ReadByte()
	// 文件大小
	m.FileSize = binary.BigEndian.Uint32(buf.Next(4))
}

type DMSAlarm struct {
	// 报警ID
	AlarmID uint32
	// 标志状态
	SignState byte
	// 报警/事件类型
	AlarmType byte
	// 报警级别
	AlarmLevel byte
	// 疲劳程度
	Fatigue byte
	// 车速
	Speed byte
	// 高程
	Altitude uint16
	// 纬度
	Latitude uint32
	// 经度
	Longitude uint32
	// 日期时间
	Time time.Time
	// 车辆状态
	Status VehicleStatus
	// 报警标识号
	Marking AlarmMarking
}

// 从缓存中读
func (da *DMSAlarm) readBy(buf *bytes.Buffer) {
	// 报警ID
	value := make([]byte, 4)
	buf.Read(value)
	da.AlarmID = binary.BigEndian.Uint32(value)

	// 标志状态
	da.SignState, _ = buf.ReadByte()

	// 报警/事件类型
	da.AlarmType, _ = buf.ReadByte()

	// 报警级别
	da.AlarmLevel, _ = buf.ReadByte()

	// 疲劳程度
	da.Fatigue, _ = buf.ReadByte()

	// 预留
	buf.Read(value)

	// 车速
	da.Speed, _ = buf.ReadByte()

	// 高程
	buf.Read(value[:2])
	da.Altitude = binary.BigEndian.Uint16(value)

	// 纬度
	buf.Read(value)
	da.Latitude = binary.BigEndian.Uint32(value)

	// 经度
	buf.Read(value)
	da.Longitude = binary.BigEndian.Uint32(value)

	// 日期时间
	dateTimeValue := make([]byte, 6)
	buf.Read(dateTimeValue)
	da.Time, _ = time.ParseInLocation(
		"060102150405",
		string(util.ParseBCD(dateTimeValue)),
		time.FixedZone("CST", 28800))

	// 车辆状态
	buf.Read(value[:2])
	da.Status = VehicleStatus(binary.BigEndian.Uint16(value))

	// 报警标识号
	da.Marking.readBy(buf)
}

// 写入缓存中
func (da *DMSAlarm) writeTo(buf *bytes.Buffer) {
	// 报警ID
	value := make([]byte, 4)
	binary.BigEndian.PutUint32(value, da.AlarmID)
	buf.Write(value)

	// 标志状态
	buf.WriteByte(da.SignState)

	// 报警/事件类型
	buf.WriteByte(da.AlarmType)

	// 报警级别
	buf.WriteByte(da.AlarmLevel)

	// 疲劳程度
	buf.WriteByte(da.Fatigue)

	// 预留，暂时填充0
	binary.BigEndian.PutUint32(value, 0)
	buf.Write(value)

	// 车速
	buf.WriteByte(da.Speed)

	// 高程
	binary.BigEndian.PutUint16(value, da.Altitude)
	buf.Write(value[:2])

	// 纬度
	binary.BigEndian.PutUint32(value, da.Latitude)
	buf.Write(value)

	// 经度
	binary.BigEndian.PutUint32(value, da.Longitude)
	buf.Write(value)

	// 日期时间
	dateTimeValue := make([]byte, 6)
	dateTimeValue = util.ToBCD([]byte(da.Time.Format("060102150405")))
	buf.Write(dateTimeValue)

	// 车辆状态
	binary.BigEndian.PutUint16(value, uint16(da.Status))
	buf.Write(value[:2])

	// 报警标识号
	da.Marking.writeTo(buf)
}

func (da *DMSAlarm) size() byte {
	return 0
}

// 车辆状态类型
type VehicleStatus uint16

// 设置ACC状态标志
func (s *VehicleStatus) SetAccOpen(open bool) {
	if open {
		*s |= 1
	} else {
		*s &= 0xFFFE
	}
}

// 左转向状态标志
func (s *VehicleStatus) SetLeftTurn(left bool) {
	if left {
		*s |= 2
	} else {
		*s &= 0xFFFD
	}
}

// 右转向状态标志
func (s *VehicleStatus) SetRightTurn(right bool) {
	if right {
		*s |= 4
	} else {
		*s &= 0xFFFB
	}
}

// 雨刮器状态标志
func (s *VehicleStatus) SetWiperOpen(open bool) {
	if open {
		*s |= 8
	} else {
		*s &= 0xFFF7
	}
}

// 制动状态标志
func (s *VehicleStatus) SetBrake(brake bool) {
	if brake {
		*s |= 16
	} else {
		*s &= 0xFFEF
	}
}

// 插卡状态标志
func (s *VehicleStatus) SetCardIn(cardIn bool) {
	if cardIn {
		*s |= 32
	} else {
		*s &= 0xFFDF
	}
}

// 定位状态标志
func (s *VehicleStatus) SetLoc(loc bool) {
	if loc {
		*s |= 1024
	} else {
		*s &= 0xFBFF
	}
}

// 设置状态
func (s *VehicleStatus) SetStatus(signal Signal, flag StatusFlag, login bool) {
	// ACC
	s.SetAccOpen(flag.IsACCOpen())
	// 左转
	if signal.GetTurnSignal() == 2 || signal.GetTurnSignal() == 3 {
		s.SetLeftTurn(true)
	}
	// 右转
	if signal.GetTurnSignal() == 1 || signal.GetTurnSignal() == 3 {
		s.SetRightTurn(true)
	}
	// 雨刮器
	s.SetWiperOpen(false)
	// 制动
	s.SetBrake(signal.IsBrake())
	// 插拔卡
	s.SetCardIn(login)
	// 定位状态
	s.SetLoc(flag.IsFixed())
}

// AlarmMarking 报警标识（苏标）
type AlarmMarking struct {
	TerminalId string    // 终端id，7个字节
	Time       time.Time // 时间
	Number     byte      // 同一时间点报警的序号，从0开始循环累加
	Count      byte      // 表示该报警对应的附件数量
	reserved   byte      // 预留字段
}

func (a *AlarmMarking) readBy(buf *bytes.Buffer) {}

func (a *AlarmMarking) writeTo(buf *bytes.Buffer) {
	// 终端id
	if l := len(a.TerminalId); l >= 7 {
		buf.Write([]byte(a.TerminalId)[:7])
	} else {
		buf.Write([]byte(a.TerminalId))
		buf.Write(make([]byte, 7-l))
	}
	// 时间
	buf.Write(util.ToBCD([]byte(a.Time.Format("060102150405"))))
	// 序号
	buf.WriteByte(a.Number)
	// 附件数量
	buf.WriteByte(a.Count)
	// 保留
	buf.WriteByte(a.reserved)
}

type TPMSAlarm struct {
	// 报警ID
	AlarmID uint32
	// 标志状态
	SignState byte
	// 车速
	Speed byte
	// 高程
	Altitude uint16
	// 纬度
	Latitude uint32
	// 经度
	Longitude uint32
	// 日期时间
	Time time.Time
	// 车辆状态
	Status VehicleStatus
	// 报警标识号
	Marking AlarmMarking
	// 报警/事件列表总数
	ListCount byte
	// 报警/事件信息列表
	List []TPMSAlarmInfo
}

type TPMSAlarmInfo struct {
	// 胎压报警位置
	TpAlaramPos byte
	// 报警/事件类型
	AlarmType uint16
	// 胎压
	TirePressure uint16
	// 胎温
	TireTemperature uint16
	// 电池电量
	Battery uint16
}

// 从缓存中读
func (alarm *TPMSAlarm) readBy(buf *bytes.Buffer) {
	// 报警ID
	value := make([]byte, 4)
	buf.Read(value)
	alarm.AlarmID = binary.BigEndian.Uint32(value)

	// 标志状态
	alarm.SignState, _ = buf.ReadByte()

	// 车速
	alarm.Speed, _ = buf.ReadByte()

	// 高程
	buf.Read(value[:2])
	alarm.Altitude = binary.BigEndian.Uint16(value)

	// 纬度
	buf.Read(value)
	alarm.Latitude = binary.BigEndian.Uint32(value)

	// 经度
	buf.Read(value)
	alarm.Longitude = binary.BigEndian.Uint32(value)

	// 日期时间
	dateTimeValue := make([]byte, 6)
	buf.Read(dateTimeValue)
	alarm.Time, _ = time.ParseInLocation(
		"060102150405",
		string(util.ParseBCD(dateTimeValue)),
		time.FixedZone("CST", 28800))

	// 车辆状态
	buf.Read(value[:2])
	alarm.Status = VehicleStatus(binary.BigEndian.Uint16(value))

	// 报警标识号
	alarm.Marking.readBy(buf)

	// 报警/事件列表总数
	alarm.ListCount, _ = buf.ReadByte()

	// 报警/事件信息列表
	for i := 0; i < int(alarm.ListCount); i++ {
		// 胎压报警位置
		alarm.List[i].TpAlaramPos, _ = buf.ReadByte()
		// 报警/事件类型
		buf.Read(value[:2])
		alarm.List[i].AlarmType = binary.BigEndian.Uint16(value)
		// 胎压
		buf.Read(value[:2])
		alarm.List[i].TirePressure = binary.BigEndian.Uint16(value)
		// 胎温
		buf.Read(value[:2])
		alarm.List[i].TireTemperature = binary.BigEndian.Uint16(value)
		// 电池电量
		buf.Read(value[:2])
		alarm.List[i].Battery = binary.BigEndian.Uint16(value)
	}
}

// 写入缓存中
func (alarm *TPMSAlarm) writeTo(buf *bytes.Buffer) {
	// 报警ID
	value := make([]byte, 4)
	binary.BigEndian.PutUint32(value, alarm.AlarmID)
	buf.Write(value)

	// 标志状态
	buf.WriteByte(alarm.SignState)

	// 车速
	buf.WriteByte(alarm.Speed)

	// 高程
	binary.BigEndian.PutUint16(value[:2], alarm.Altitude)
	buf.Write(value[:2])

	// 纬度
	binary.BigEndian.PutUint32(value, alarm.Latitude)
	buf.Write(value)

	// 经度
	binary.BigEndian.PutUint32(value, alarm.Longitude)
	buf.Write(value)

	// 日期时间
	dateTimeValue := make([]byte, 6)
	dateTimeValue = util.ToBCD([]byte(alarm.Time.Format("060102150405")))
	buf.Write(dateTimeValue)

	// 车辆状态
	binary.BigEndian.PutUint16(value, uint16(alarm.Status))
	buf.Write(value[:2])

	// 报警标识号
	alarm.Marking.writeTo(buf)

	// 报警/事件列表总数
	buf.WriteByte(alarm.ListCount)

	// 报警/事件信息列表
	for i := 0; i < int(alarm.ListCount); i++ {
		// 胎压报警位置
		buf.WriteByte(alarm.List[i].TpAlaramPos)
		// 报警/事件类型
		binary.BigEndian.PutUint16(value[:2], alarm.List[i].AlarmType)
		buf.Write(value[:2])
		// 胎压
		binary.BigEndian.PutUint16(value[:2], alarm.List[i].TirePressure)
		buf.Write(value[:2])
		// 胎温
		binary.BigEndian.PutUint16(value[:2], alarm.List[i].TireTemperature)
		buf.Write(value[:2])
		// 电池电量
		binary.BigEndian.PutUint16(value[:2], alarm.List[i].Battery)
		buf.Write(value[:2])
	}
}
