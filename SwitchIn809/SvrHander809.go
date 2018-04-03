package SwitchIn809

import (
	"encoding/hex"
	"log"

	"github.com/tiptok/GoNas/global"
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
	log.Printf("%v On Receive Data : %v", c.RemoteAddress, hex.EncodeToString(d.Bytes()))
	defer func() {
		if p := recover(); p != nil {
			log.Printf("panic recover! p: %v", p)
			//debug.PrintStack()
		}
	}()
	obj, err := c.ParseToEntity(d.Bytes())
	if err != nil {
		log.Panicln(err.Error())
	}
	// if def, ok := obj.(conn.DefaultTcpData); ok {
	// 	log.Printf("收到MsgTypeId：%v  Begin: %v  End: %v", def.MsgTypeId, def.BEGIN, def.END)
	// 	/*添加上行*/
	// 	global.UpHandler.UpData(def)
	// } else {
	// 	log.Println("Convert To Type Error.")
	// }
	global.UpHandler.UpData(obj.(interface{}))
	return true
}
