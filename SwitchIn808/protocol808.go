package SwitchIn808

import "github.com/tiptok/gotransfer/conn"
import "github.com/tiptok/gotransfer/comm"
import "errors"
import "log"

type protocol808 struct {
}

func (p protocol808) PacketMsg(obj interface{}) (data []byte, err error) {
	return nil, nil
}

/*
	打包数据体
	obj 数据体
*/
func (p protocol808) Packet(obj interface{}) (packdata []byte, err error) {
	return nil, nil
}

/*
	分包处理
	packdata 解析出一个完整包
	leftdata 解析剩余报文的数据
	err 	 分包错误
*/
func (p protocol808) ParseMsg(data []byte, c *conn.Connector) (packdata []byte, leftdata []byte, err error) {
	defer func() {
		conn.MyRecover()
	}()
	if data == nil || len(data) == 0 {
		err = errors.New("未包含tcp数据")
		return packdata, leftdata, err
	}
	ibegin := -1
	iEnd := -1
	for i := 0; i < len(data); i++ {
		log.Printf("Index:%x  %x %t", i, data[i], data[i] == 0x7e)
		if data[i] == 0x7e {
			ibegin = i
		}
		if data[i] == 0x7d && ibegin >= 0 && ibegin != i {
			iEnd = i + 1
			log.Printf("Begin:%x End:%x", ibegin, iEnd)
		}
		if ibegin >= 0 && iEnd > 0 {
			packdata = data[ibegin:iEnd]
			if iEnd < len(data) {
				leftdata = data[iEnd:]
				_, err := c.WriteLeftData(leftdata)
				if err != nil {
					log.Println(err.Error())
				}
			}
			break
		}
	}
	/*未找到头标识 说明报文是非法数据*/
	if ibegin < 0 {
		err = errors.New("tcp数据格式不对")
	}
	return packdata, leftdata, err
}

/*
	解析数据
	obj 解析出对应得数据结构
	err 解析出错
	7e
	1001
	0001
	00000000
	00000000
	00
	7d
	7e100200010000000000000000007d
	7e100300010000000000000000007d
*/
func (p protocol808) Parse(packdata []byte) (obj interface{}, err error) {
	defer func() {
		conn.MyRecover()
	}()
	def := conn.DefaultTcpData{}
	def.BEGIN = packdata[0]
	def.MsgTypeId = comm.BinaryHelper.ToInt16(packdata, 1)
	def.Id = comm.BinaryHelper.ToInt16(packdata, 3)
	def.Length = comm.BinaryHelper.ToInt32(packdata, 5)
	def.PackagesProperty = comm.BinaryHelper.ToInt32(packdata, 9)
	def.Valid = packdata[14]
	def.END = packdata[len(packdata)-1]
	obj = def
	return obj, err
}