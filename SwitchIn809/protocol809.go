package SwitchIn809

import (
	"log"

	"github.com/tiptok/gotransfer/comm"
	"github.com/tiptok/gotransfer/conn"
)

//var _JTB808ParseBase *JTB808ParseBase = &JTB808ParseBase{}
type protocol809 struct {
}

func (p protocol809) PacketMsg(obj interface{}) (data []byte, err error) {
	return nil, nil
}

/*
	打包数据体
	obj 数据体
*/
func (p protocol809) Packet(obj interface{}) (packdata []byte, err error) {
	return nil, nil
}

/*
	分包处理
	packdata 解析出一个完整包
	leftdata 解析剩余报文的数据
	err 	 分包错误
*/
func (p protocol809) ParseMsg(data []byte, c *conn.Connector) (packdata [][]byte, leftdata []byte, err error) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("protocol809.ParseMsg panic recover! p: %v", p)
		}
	}()
	return comm.ParseHelper.ParsePart(data, 0x5b, 0x5d)
}

/*
	解析数据
	obj 解析出对应得数据结构
	err 解析出错
*/
func (p protocol809) Parse(packdata []byte) (obj interface{}, err error) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("protocol809.Parse panic recover! p: %v", p)
		}
	}()
	return obj, err
}
