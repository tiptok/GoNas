package SwitchIn809

import (
	"encoding/hex"
	"log"

	"github.com/tiptok/GoNas/global"
	"github.com/tiptok/GoNas/model"
	"github.com/tiptok/gotransfer/conn"
)

type SvrHander809 struct {
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
	global.Debug("%v On Receive Data : %v", c.RemoteAddress, hex.EncodeToString(d.Bytes()))
	defer func() {
		if p := recover(); p != nil {
			log.Printf("panic recover! p: %v", p)
			//debug.PrintStack()
		}
	}()
	obj, err := c.ParseToEntity(d.Bytes())
	if err != nil {
		global.Error(err.Error())
		return false
	}
	// if obj != nil {
	// 	log.Println("Receive Entity:", obj)
	// }
	if def, ok := obj.(model.IEntity); ok {
		entity := def.GetEntityBase()
		//log.Printf("MsgId:%v  MsgSN:%v AccessCode:%v", entity.MsgId, entity.MsgSN, entity.AccessCode)
		global.Debug("MsgId:%v  MsgSN:%v AccessCode:%v", entity.MsgId, entity.MsgSN, entity.AccessCode)
	}
	global.UpHandler.UpData(obj.(interface{}))
	return true
}
