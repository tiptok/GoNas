package SwitchIn809

import (
	"log"

	"errors"

	"fmt"

	"github.com/tiptok/GoNas/model"
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
	//转义
	data, err := Byte809Descape(packdata, 0, len(packdata))

	//CRC Check
	tmpCrc := comm.BinaryHelper.ToInt16(data, int32(len(data)-2))
	checkCrc := comm.BinaryHelper.CRC16Check(data[:len(data)-2])
	if checkCrc != tmpCrc {
		err = errors.New(fmt.Sprintf("CRC CHECK Error->%d != %d", tmpCrc, checkCrc))
		return nil, err
	}
	h := model.EntityBase{}
	//数据头
	msgBodyLenght := comm.BinaryHelper.ToInt32(data, 0) - 26
	h.MsgSN = int(comm.BinaryHelper.ToInt32(data, 4))
	h.MsgId = fmt.Sprintf("%d", comm.BinaryHelper.ToInt16(data, 8))
	h.AccessCode = fmt.Sprintf("%d", comm.BinaryHelper.ToInt32(data, 10))

	isEncrypt := data[17] == 0
	if isEncrypt {
		//解密
	}
	msgBody := data[22 : msgBodyLenght+22]

	sMethodName := fmt.Sprintf("J%s", comm.BinaryHelper.ToBCDString(data, 8, 2))
	//InvokeFunc
	//log.Println(msgBodyLenght, sMethodName, comm.BinaryHelper.ToBCDString(msgBody, 0, int32(len(msgBody))))
	value, err := comm.ParseHelper.InvokeFunc(&JTB809ParseBase{}, sMethodName, msgBody, h)
	if value != nil && len(value) > 0 {
		obj = value[0].Interface()
	}

	return obj, err
}
