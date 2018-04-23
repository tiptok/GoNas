package model

import (
	"bytes"
	"fmt"
	"time"
)

type TermPosition struct {
	EntityBase

	/*基本信息*/
	SimNum       string    //终端手机号码
	GpsTime      time.Time //卫星时间:精确到秒
	AlarmFlag    int64     //报警标志
	StateFlag    int64     //状态标志
	ExStateFlag  int64     //扩展状态标志（扩展）
	Reserved     int
	Lon          float32 //经度
	Lat          float32 //纬度
	Speed        float32 //卫星速度: #公里/小时
	ADRSpeed     float32 //行驶记录仪速度：#公里/小时
	Mileage      float32 //里程：#公里
	Altitude     int     //海拔高度：#米
	Direction    int     //方向：0-359，正北为 0，顺时针
	ResidualFuel float32 //剩余油量：#升，对应车上油量表读数

	StarCount      int //定位卫星数
	WirelessSignal int //无线信号强度
	Analogue       int //模拟量

	/*超速报警信息*/
	SpeedingPosType int //“超速报警”位置类型( 0：无特定位置； 1：圆形区域； 2：矩形区域； 3：多边形区域；4：路段)
	SpeedingPosId   int //“超速报警”位置区域或路段Id（区域ID或路段ID）
	LimitSpeed      int //限制速度：#公里/小时（扩展）

	/*区域/路线报警信息*/
	InOutPosType int //“进出区域/路线报警” 位置类型 (1：圆形区域； 2：矩形区域； 3：多边形区域；4：路线)
	InOutPosId   int //“进出区域/路线报警”Id
	InOutFlag    int //“进出区域/路线报警”标志（0：进； 1：出 ）

	/*路段行驶过长/不足报警信息*/
	DrivingPathSectionId int //"行驶报警"路段Id
	DrivingDuration      int //“行驶报警”路段行驶时长：#秒
	DrivingFlag          int //“行驶报警”标识(0：不足；1：过长)

	/*多媒体信息*/
	MediaType       int    //多媒体类型 0：图像；1：音频；2：视频； -1：无多媒体
	MediaFormatCode int    //多媒体格式编码 0：JPEG；1：TIF；2：MP3；3：WAV；4：WMV；其他保留
	MediaEventCode  int    // 多媒体事件项编码 0：平台下发指令；1：定时动作；2：抢劫报警触发；3：碰撞侧翻报警触发；0xE1：驾驶员刷卡触发； 其他保留
	MediaChannelId  int    //多媒体通道Id; 多媒体类型为-1时此字段表示外设ID
	MediaData       []byte //多媒体数据; 多媒体类型为-1时此字段表示外设数据

	AdditionData string // 附加数据（16进制字符串）
}

func (e *TermPosition) GetMsgId() interface{} {
	return e.MsgId
}
func (e *TermPosition) GetEntityBase() *EntityBase {
	return e.EntityBase.GetEntityBase()
}

func (e *TermPosition) GetDBSql() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("EXEC spNas_AddTrack ")
	buf.WriteString(fmt.Sprintf("'%v',", e.SimNum))
	buf.WriteString(fmt.Sprintf("'%v',", e.GpsTime.Format("2006-01-02 03:04:05")))
	buf.WriteString(fmt.Sprintf("%v,", e.AlarmFlag))
	buf.WriteString(fmt.Sprintf("%v,", e.StateFlag))
	buf.WriteString(fmt.Sprintf("%v,", e.ExStateFlag))
	buf.WriteString(fmt.Sprintf("%v,", e.Lon))
	buf.WriteString(fmt.Sprintf("%v,", e.Lat))
	buf.WriteString(fmt.Sprintf("%v,", e.Speed))
	buf.WriteString(fmt.Sprintf("%v,", e.Altitude))
	buf.WriteString(fmt.Sprintf("%v,", e.Direction))
	buf.WriteString(fmt.Sprintf("%v,", e.Mileage))

	return buf.String()
}
