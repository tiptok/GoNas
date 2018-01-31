package SwitchIn808

import (
	"github.com/tiptok/gotransfer/comm"
	"fmt"
)

type IMethodName interface{
	GetMethodName() string  
}

/*
	MsgHead 存放
	报文头部信息
*/
type MsgHead struct{
	/*消息标识 1001*/
	Id int
	/*消息序号*/
	SN int
	/*标识 作SimNum*/
	Identify string
	/*总包数*/
	TotalPkgCount int
	/*当前包数*/
	CurOrder   int
}
func(h MsgHead)GetMethodName() string{
	data :=comm.BinaryHelper.Int16ToBytes(int16(h.Id))
	return fmt.Sprintf("J%s",comm.BinaryHelper.ToBCDString(data,int32(0),int32(len(data))))
}
/*
	New MsgHead
*/
func NewMsgHead(Id int,SN int,Identify string,TotalPkgCount int,CurOrder int) MsgHead{
	return MsgHead{
		Id:Id,
		SN:SN,
		Identify:Identify,
		TotalPkgCount:TotalPkgCount,
		CurOrder:CurOrder,
	}
}


/****************808实体定义*******************/
const(
	终端通用应答 = 0x0001
	终端心跳 = 2
    终端注销 = 3
    终端注册 = 256
    终端鉴权 = 258
)
type EntityBase struct{
	MsgHead 
}

type UpDataEntityBase struct{
	EntityBase
	CmdCode int
	/*是否需要应答*/
	IsNeedRsp bool
	DownRspEntity interface{}
}
type DownDataEntityBase struct{
	EntityBase
	CmdCode int
	/*是否需要应答*/
}
/*
	终端通用应答 0x0001
*/
type TermCommonReply struct{
	UpDataEntityBase
	RspMsgSN uint16
	RspMsgId uint16
	RspResult int
}

