package SwitchIn808

import (
	"encoding/hex"
	"log"
	"runtime/debug"

	"github.com/tiptok/GoNas/global"
	"github.com/tiptok/gotransfer/conn"
)

type SvrHander808 struct {
}

//连接事件
func (trans *SvrHander808) OnConnect(c *conn.Connector) bool {
	defer func() {
		//conn.MyRecover()
	}()
	log.Println(c.RemoteAddress, "On Connect.")
	return true
}

//断开事件
func (trans *SvrHander808) OnClose(c *conn.Connector) {
	log.Println(c.RemoteAddress, "On Close.")
}

//接收事件
func (trans *SvrHander808) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	log.Printf("%v On Receive Data : %v", c.RemoteAddress, hex.EncodeToString(d.Bytes()))
	defer func() {
		if p := recover(); p != nil {
			log.Printf("panic recover! p: %v", p)
			debug.PrintStack()
		}
	}()
	/*解析分包数据*/
	if c.P == nil {
		//记录空日志
		return false
	}
	packdata, _, err := c.P.ParseMsg(d.Bytes(), c)
	if err != nil {
		log.Println(err.Error())
	}
	/*解析完整包*/
	if packdata != nil {
		for i:=0;i<len(packdata);i++{
			if len(packdata[i])<=0{
				continue
			}	
			log.Printf("%v On Receive Part Data : %v", c.RemoteAddress, hex.EncodeToString(packdata[i]))
			obj, err1 := c.P.Parse(packdata[i])
			if err1 != nil {
				log.Panicln(err1.Error())
			}
			if def, ok := obj.(conn.DefaultTcpData); ok {
				log.Printf("收到MsgTypeId：%v  Begin: %v  End: %v", def.MsgTypeId, def.BEGIN, def.END)
				/*添加上行*/
				global.UpHandler.UpData(def)
			} else {
				log.Println("Convert To Type Error.")
			}
		}
	}
	/*剩余bytes*/
	if c.Leftbuf.Len()>0{
		log.Printf("%v On Receive Part Data : %v", c.RemoteAddress, hex.EncodeToString(c.Leftbuf.Bytes()))
	}
	/*解析出实体*/
	if c.IsConneted {
		c.SendChan <- d
	}
	return true
}
