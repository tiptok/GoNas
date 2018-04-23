package SwitchIn809

import (
	"encoding/hex"
	"log"

	"strings"

	"github.com/tiptok/GoNas/global"
	"github.com/tiptok/GoNas/model"
	"github.com/tiptok/gotransfer/comm"
	"github.com/tiptok/gotransfer/conn"
)

type SvrHander809 struct {
	conn.TcpServerBase
}

//连接事件
func (trans *SvrHander809) OnConnect(c *conn.Connector) bool {
	defer func() {
		//conn.MyRecover()
	}()
	log.Println(c.RemoteAddress, "On Connect.")
	return true
}

//断开事件
func (trans *SvrHander809) OnClose(c *conn.Connector) {
	log.Println(c.RemoteAddress, "On Close.")
}

//接收事件
func (trans *SvrHander809) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	var bUpData bool = true
	global.Debug("%v On Receive Data : %v", c.RemoteAddress, hex.EncodeToString(d.Bytes()))
	defer func() {
		if p := recover(); p != nil {
			log.Printf("SvrHander809 OnReceive panic recover! p: %v", p)
			//debug.PrintStack()
		}
	}()
	obj, err := c.ParseToEntity(d.Bytes())
	if err != nil {
		global.Error(err.Error())
		return false
	}
	var rspEntity model.IEntity //应答实体
	if def, ok := obj.(model.IEntity); ok {
		entity := def.GetEntityBase()
		cmdcode := entity.MsgId.(int16)
		if entity.SubMsgId != nil && entity.SubMsgId.(int16) != 0 {
			cmdcode = entity.SubMsgId.(int16)
		}
		global.Debug("MsgId:%X  MsgSN:%d AccessCode:%v", cmdcode, entity.MsgSN, entity.AccessCode)
		switch cmdcode {
		case model.J主链路登录请求:
			login := obj.(*model.UP_CONNECT_REQ)
			//login.AccessCode global.Param.AccessCode && login.USERID="" && login.PASSWORD==""
			if strings.Compare(login.AccessCode, "12345678") == 0 {
				rspEntity = &model.UP_CONNECT_RSP{EntityBase: model.EntityBase{MsgId: model.J主链路登录应答}, Result: 0, Verify_Code: int32(global.Param.VerifyCode)}
			}
			// case model.主链路注销请求:
		case model.J主链路连接保持请求:
		case model.J实时上传车辆定位信息:
			bUpData = false
			global.UpHandler.UpData((obj.(*model.UP_EXG_MSG_REAL_LOCATION)).GetConvEntity())
			//global.Debug("接收到实体%v", obj)
		default:
		}
		if rspEntity != nil {
			base := rspEntity.GetEntityBase()
			base.AccessCode = entity.AccessCode
		}
		//上行
		if bUpData {
			global.UpHandler.UpData(def)
		}
	} else {
		global.Debug("接收到实体%v", obj)
	}
	//发送应答
	if rspEntity != nil {
		//IEntity
		data, err := conn.SendEntity(rspEntity, c)
		if err != nil {
			global.Error("SvrHander Send Entity Error:%v", err)
		} else {
			global.Debug("SvrHander Send Data:%s", comm.BinaryHelper.ToBCDString(data, 0, int32(len(data))))
		}
	}
	return true
}
