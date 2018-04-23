package SwitchIn809

import (
	"fmt"
	"time"

	"strings"

	"github.com/axgle/mahonia"
	"github.com/tiptok/GoNas/model"
	"github.com/tiptok/gotransfer/comm"
)

type JTB809ParseBase struct {
}

// func(p *JTB809ParseBase)Parse(h MsgHead,msgBody []byte)(interface{},error){

// }

//J1001 0x1001 主链路登录请求
func (p *JTB809ParseBase) J1001(msgBody []byte, h model.EntityBase) interface{} {
	outEntity := &model.UP_CONNECT_REQ{}
	outEntity.SetEntity(h)
	outEntity.USERID = uint32(comm.BinaryHelper.ToInt32(msgBody, 0)) //应答ID
	outEntity.PASSWORD = comm.BinaryHelper.ToASCIIString(msgBody, 4, 8)
	outEntity.DOWN_LINK_IP = comm.BinaryHelper.ToASCIIString(msgBody, 12, 32) //结果
	outEntity.DOWN_LINK_PORT = uint32(comm.BinaryHelper.ToInt16(msgBody, 44)) //应答ID

	/*是否需要应答*/
	// outEntity.IsNeedRsp = false
	// outEntity.DownRspEntity = nil
	return outEntity
}

//J1002 主链路登录应答 0x1002
func (p *JTB809ParseBase) J1002(msgBody []byte, h model.EntityBase) interface{} {
	outEntity := &model.UP_CONNECT_RSP{}
	outEntity.SetEntity(h)
	outEntity.Result = msgBody[0]
	outEntity.Verify_Code = comm.BinaryHelper.ToInt32(msgBody, 1)
	/*是否需要应答*/
	// outEntity.IsNeedRsp = false
	// outEntity.DownRspEntity = nil
	return outEntity
}

//J1005 0x1005  主链路保持连接请求
func (p *JTB809ParseBase) J1005(msgBody []byte, h model.EntityBase) interface{} {
	outEntity := &model.UP_LINKTEST_REQ{}
	outEntity.SetEntity(h)
	// outEntity.Result = msgBody[0]
	// outEntity.Verify_Code = comm.BinaryHelper.ToInt32(msgBody, 1)
	/*是否需要应答*/
	// outEntity.IsNeedRsp = false
	// outEntity.DownRspEntity = nil
	return outEntity
}

//J1200 主链路动态信息交换
func (p *JTB809ParseBase) J1200(msgBody []byte, h model.EntityBase) interface{} {
	outEntity := model.UP_EXG_MSG{}
	outEntity.SetEntity(h)
	enc := mahonia.NewDecoder("gbk")
	outEntity.Vehicle_No = strings.Trim(comm.BinaryHelper.ToASCIIString(msgBody, 0, 21), string([]byte{0x00})) // strings.Trim(comm.BinaryHelper.ToASCIIString(msgBody, 0, 21), " ")
	outEntity.Vehicle_No = enc.ConvertString(outEntity.Vehicle_No)
	outEntity.Vehicle_Color = msgBody[21]
	outEntity.SubMsgId = comm.BinaryHelper.ToInt16(msgBody, 22)
	switch outEntity.SubMsgId.(int16) {
	case 0x1202:
		return J1202(msgBody, outEntity)
	default:
		panic(fmt.Sprintf("未找到对应方法:%v", outEntity.MsgId))
	}
	return nil
}
func J1202(msgBody []byte, h model.UP_EXG_MSG) interface{} {
	outEntity := &model.UP_EXG_MSG_REAL_LOCATION{
		UP_EXG_MSG: h,
	}
	outEntity.GNSS_DATA = GetGetLoactionInfo(msgBody, 28)
	return outEntity
}

func GetGetLoactionInfo(msgBody []byte, iIndex int32) (location model.LocationInfoItem) {
	location.ENCRYPT = msgBody[iIndex]
	iIndex += 1
	iDay := int(msgBody[iIndex])
	iMonth := int(msgBody[iIndex+1])
	iYear := int(comm.BinaryHelper.ToInt16(msgBody, int32(iIndex+2)))
	iHour := int(msgBody[iIndex+4])
	iMin := int(msgBody[iIndex+5])
	iSec := int(msgBody[iIndex+6])
	var err error
	//location.GPSTIME, err = time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%v-%v-%v %v:%v:%v", iYear, iMonth, iDay, iHour, iMin, iSec))
	location.GPSTIME = time.Date(iYear, time.Month(iMonth), iDay, iHour, iMin, iSec, 0, time.Local)
	if err != nil {
		panic(err)
	}
	location.LON = int(comm.BinaryHelper.ToInt32(msgBody, iIndex+7))
	location.LAT = int(comm.BinaryHelper.ToInt32(msgBody, iIndex+11))
	location.VEC1 = comm.BinaryHelper.ToInt16(msgBody, iIndex+15)
	location.VEC2 = comm.BinaryHelper.ToInt16(msgBody, iIndex+17)
	location.VEC3 = int(comm.BinaryHelper.ToInt32(msgBody, iIndex+19))
	location.DIRECTION = comm.BinaryHelper.ToInt16(msgBody, iIndex+23)
	location.ALTITUDE = comm.BinaryHelper.ToInt16(msgBody, iIndex+25)
	location.STATE = int(comm.BinaryHelper.ToInt32(msgBody, iIndex+27))
	location.ALARM = int(comm.BinaryHelper.ToInt32(msgBody, iIndex+31))
	return location
}
