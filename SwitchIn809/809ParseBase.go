package SwitchIn809

import (
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
