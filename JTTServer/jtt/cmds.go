package jtt

const (
	msgIDTerminalResponse          = uint16(1)      // 终端通用应答
	MsgIDTerminalHeartbeat         = uint16(2)      // 终端心跳
	msgIDTerminalLogout            = uint16(3)      // 终端注销
	msgIDServerTime                = uint16(4)      // 查询服务器时间
	msgIDTerminalPackResend        = uint16(5)      // 终端补传分包请求
	msgIDTerminalLogin             = uint16(0x0100) // 终端注册
	MsgIDTerminalAuth              = uint16(0x0102) // 终端鉴权
	msgIDGetTerminalParamsResp     = uint16(0x0104) // 查询终端参数应答
	msgIDGetTerminalAttrResp       = uint16(0x0107) // 查询终端属性应答
	msgIDTerminalUpgradeResp       = uint16(0x0108) // 终端升级应答
	MsgIDPositionReport            = uint16(0x0200) // 位置信息汇报
	msgIDGetPositionResp           = uint16(0x0201) // 位置信息查询应答
	msgIDBDLocationCheck           = uint16(0x0205) // 北斗验真上报
	msgIDVehicleControlResp        = uint16(0x0500) // 车辆控制应答
	msgIDGetAreaResp               = uint16(0x0608) // 查询区域或路线数据应答
	msgIDDrivingRecordReport       = uint16(0x0700) // 行驶记录数据上传
	msgIDWaybillReport             = uint16(0x0701) // 电子运单上报
	msgIDDriverIdentityReport      = uint16(0x0702) // 驾驶员身份信息上报
	MsgIDPositionBatchReport       = uint16(0x0704) // 定位数据批量上传
	msgIDCANDataReport             = uint16(0x0705) // CAN总线数据上传
	msgIDCanDataUpload             = uint16(0x0705) // CAN总线数据上传
	msgIDMultimediaEventReport     = uint16(0x0800) // 多媒体事件信息上传
	msgIDMultimediaDataReport      = uint16(0x0801) // 多媒体数据上传
	msgIDGetMultimediaSaveInfoResp = uint16(0x0802) // 存储多媒体数据检索应答
	msgIDDriverFaceReport          = uint16(0x0803) // 驾驶员人脸信息采集上报
	msgIDSnapshotResp              = uint16(0x0805) // 摄像头立即拍摄命令应答
	msgIDDataUpPenetrate           = uint16(0x0900) // 数据上行透传
	msgIDDataCompressionReport     = uint16(0x0901) // 数据压缩上报
	msgIDTerminalRSAPublickey      = uint16(0x0A00) // 终端RSA公钥
	msgIDMediaPropertyReport       = uint16(0x1003) // 终端上传音视频属性
	msgIDMediaResourceListReport   = uint16(0x1205) // 终端上传音视频资源列表
	msgIDFileUploadFinish          = uint16(0x1206) // 文件上传完成通知

	msgIDServerResponse           = uint16(0x8001) // 平台通用应答
	msgIDServerPackResend         = uint16(0x8003) // 服务器补传分包请求
	msgIDServerTimeResp           = uint16(0x8004) // 查询服务器时间应答
	msgIDTerminalLoginResp        = uint16(0x8100) // 终端注册应答
	msgIDSetTerminalParams        = uint16(0x8103) // 设置终端参数
	msgIDGetTerminalParams        = uint16(0x8104) // 查询终端参数
	msgIDTerminalControl          = uint16(0x8105) // 终端控制
	msgIDGetTerminalSpecParams    = uint16(0x8106) // 查询指定终端参数
	msgIDGetTerminalAttr          = uint16(0x8107) // 查询终端属性
	msgIDTerminalUpgrade          = uint16(0x8108) // 下发终端升级包
	msgIDGetPosition              = uint16(0x8201) // 位置信息查询
	msgIDTrackControl             = uint16(0x8202) // 临时位置跟踪控制
	msgIDManualConfirmAlarm       = uint16(0x8203) // 人工确认报警信息
	msgIDLinkCheck                = uint16(0x8204) // 链路检测
	msgIDTextIssued               = uint16(0x8300) // 文本信息下发
	msgIDTELCallback              = uint16(0x8400) // 电话回拨
	msgIDSetContacts              = uint16(0x8401) // 设置电话本
	msgIDVehicleControl           = uint16(0x8500) // 车辆控制
	msgIDSetRoundArea             = uint16(0x8600) // 设置圆形区域
	msgIDDeleteRoundArea          = uint16(0x8601) // 删除圆形区域
	msgIDSetRectArea              = uint16(0x8602) // 设置矩形区域
	msgIDDeleteRectArea           = uint16(0x8603) // 删除矩形区域
	msgIDSetPolygonArea           = uint16(0x8604) // 设置多边形区域
	msgIDDeletePolygonArea        = uint16(0x8605) // 删除多边形区域
	msgIDSetPolyline              = uint16(0x8606) // 设置路线
	msgIDDeletePolyline           = uint16(0x8607) // 删除路线
	msgIDGetArea                  = uint16(0x8608) // 查询区域或路线数据
	msgIDGatherDrivingRecord      = uint16(0x8700) // 行车记录数据采集
	msgIDDriRecordParamsIssued    = uint16(0x8701) // 行驶记录参数下发
	msgIDGetDriverIdentity        = uint16(0x8702) // 驾驶员身份信息查询
	msgIDMultimediaDataReportResp = uint16(0x8800) // 多媒体数据上传应答
	msgIDSnapshot                 = uint16(0x8801) // 摄像头立即拍摄命令
	msgIDGetMultimediaSaveInfo    = uint16(0x8802) // 存储多媒体数据检索
	msgIDGetMultimediaSaveUp      = uint16(0x8803) // 存储多媒体数据上传
	msgIDGetAudioStartRecord      = uint16(0x8804) // 开始录音
	msgIDGetSingleMediaSaveUp     = uint16(0x8805) // 单条存储多媒体数据检索上传
	msgIDDataDownPenetrate        = uint16(0x8900) // 数据下行透传
	msgIDPlatformRSAPublickey     = uint16(0x8A00) // 平台RSA公钥
	msgIDQueryMediaProperty       = uint16(0x9003) // 查询音视频属性
	msgIDRealMediaRequest         = uint16(0x9101) // 实时音频传输请求
	msgIDRealMediaControl         = uint16(0x9102) // 实时音频传输控制
	msgIDRealMediaNotice          = uint16(0x9105) // 实时音频传输状态通知
	msgIDRemoteVideoReplay        = uint16(0x9201) // 远程录像回放
	msgIDRemoteReplayControl      = uint16(0x9202) // 远程回放控制
	msgIDMediaResourceSelect      = uint16(0x9205) // 查询音视频资源
	msgIDFileUploadCmd            = uint16(0x9206) // 文件上传指令
	msgIDFileUploadCtl            = uint16(0x9207) // 文件上传控制
	msgIDAttachUpload             = uint16(0x9208) // 附件上传
	msgIDPtzTurn                  = uint16(0x9301) // 云台转动控制
	msgIDPtzFocus                 = uint16(0x9302) // 云台焦距控制
	msgIDPtzAperture              = uint16(0x9303) // 云台光圈控制
	msgIDPtzWiper                 = uint16(0x9304) // 云台雨刷控制
	msgIDPtzFilllight             = uint16(0x9305) // 云台红外补光
	msgIDPtzZoom                  = uint16(0x9306) // 云台变倍控制
)
