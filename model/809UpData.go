package model

import (
	"time"

	"github.com/tiptok/gotransfer/comm"
)

//1001 主链路登录请求
type UP_CONNECT_REQ struct {
	EntityBase
	USERNAME       string //用户名 作为用户ID的补充（有些809协议吧UserId定义为串型）
	USERID         uint32 //用户ID
	PASSWORD       string //用户密码
	DOWN_LINK_IP   string //从链路IP
	DOWN_LINK_PORT uint32 //从链路端口
}

//1005 主链路连接保持请求
type UP_LINKTEST_REQ struct {
	EntityBase
}

//1008 下级平台主动关闭链路通知
type UP_CLOSELINK_INFORM struct {
	EntityBase
	Reason_Code byte //链路关闭原因
}

//IsValid valid UP_CLOSELINK_INFORM
func (e *UP_CLOSELINK_INFORM) IsValid(args ...interface{}) (errMsg string) {
	if e.Reason_Code == 0x00 {
		errMsg = "网关重启"
	} else if e.Reason_Code == 0x01 {
		errMsg = "其他原因"
	} else {
		errMsg = "其他原因"
	}
	return errMsg
}

//UP_EXG_MSG 0x1200 主链路动态信息交换
type UP_EXG_MSG struct {
	EntityBase
	Vehicle_No    string //车牌号
	Vehicle_Color byte   //车牌颜色
	SimNum        string //手机号
}

func (e *UP_EXG_MSG) New_UP_EXG_MSG(vn string, vc byte, subid interface{}) *UP_EXG_MSG {
	e.Vehicle_No = vn
	e.Vehicle_Color = vc
	//e.SimNum = s
	e.SubMsgId = subid
	return e
}

type LocationInfoItem struct {
	ENCRYPT    byte      // 加密标识 1 已加密 0未加密
	GPSTIME    time.Time // 时间
	LON        int       // 经度，单位1*10 -6度
	LAT        int       // 纬度，单位1*10 -6度
	VEC1       int16     // 速度 行车速度信息 单位 千米每小时（km/h）
	VEC2       int16     // 行驶记录仪速度单位 千米每小时（km/h）
	VEC3       int       // 车辆当前总里程数 单位千米（km）
	DIRECTION  int16     // 方向 0~359 单位为度 ，正北为0，顺时针
	ALTITUDE   int16     // 海拔高度 单位 米
	STATE      int       // 车辆状态
	ALARM      int       // 报警状态
	SPEEDLIMIT int       // 限制速度 （四川过检添加，用于替换 海拔高度字段)
}

//UP_EXG_MSG_REAL_LOCATION 实时上传车辆定位信息 子业务命令码（0x1202）
type UP_EXG_MSG_REAL_LOCATION struct {
	UP_EXG_MSG
	GNSS_DATA LocationInfoItem
}

func (e *UP_EXG_MSG_REAL_LOCATION) GetConvEntity() IEntity {
	termPos := &TermPosition{
		Id:         comm.BinaryHelper.UniqueId(), //NewObjectId(),
		SimNum:     "18860183011",
		GpsTime:    e.GNSS_DATA.GPSTIME,
		AlarmFlag:  int64(e.GNSS_DATA.ALARM),
		StateFlag:  int64(e.GNSS_DATA.STATE),
		Lon:        float32(e.GNSS_DATA.LON) / 1000000.0,
		Lat:        float32(e.GNSS_DATA.LAT) / 1000000.0,
		Speed:      float32(e.GNSS_DATA.VEC1),
		ADRSpeed:   float32(e.GNSS_DATA.VEC2),
		Mileage:    float32(e.GNSS_DATA.VEC3),
		Altitude:   int(e.GNSS_DATA.ALTITUDE),
		Direction:  int(e.GNSS_DATA.DIRECTION),
		EntityBase: *(e.GetEntityBase()),
	}
	return termPos
}
