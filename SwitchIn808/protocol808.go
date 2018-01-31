package SwitchIn808

import "github.com/tiptok/gotransfer/conn"
import "github.com/tiptok/gotransfer/comm"
import "errors"
import "log"

var _JTB808ParseBase *JTB808ParseBase = &JTB808ParseBase{}
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
func (p protocol808) ParseMsg(data []byte, c *conn.Connector) (packdata [][]byte, leftdata []byte, err error) {

	defer func() {
		conn.MyRecover()
	}()
	if data == nil || len(data) == 0 {
		err = errors.New("未包含tcp数据")
		return packdata, leftdata, err
	}
	ibegin := -1
	iEnd := -1
	packdata = make([][]byte,1)
	for i := 0; i < len(data); i++ {
		log.Printf("Index:%x  %x %t", i, data[i], data[i] == 0x7e)
		if data[i] == 0x7e {
			ibegin = i
		}
		if data[i] == 0x7e && ibegin >= 0 && ibegin != i {
			iEnd = i + 1
			log.Printf("Begin:%x End:%x", ibegin, iEnd)
		}
		if ibegin >= 0 && iEnd > 0 {
			/*添加到data list */
			packdata = append(packdata, data[ibegin:iEnd])
			//
			/*重置下标*/
			ibegin, iEnd = -1, -1
			continue
		}
		/*退出分包 将剩余bytes写到leftbuffer 里面*/
		if ibegin >= 0 && i+1==len(data) {
			if iEnd < len(data) {
				leftdata = data[ibegin:]
				_, err := c.WriteLeftData(leftdata)
				if err != nil {
					log.Println(err.Error())
				}
			}
			break
		}
	}
	/*未找到头标识 说明报文是非法数据*/
	if ibegin < 0 && len(packdata) == 1 {
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
	7e100300010000000000000000007d7e1004   leftData :7e1004
	7e100200010000000000000000007d7e100300010000000000000000007d 多包
	def := conn.DefaultTcpData{}
	def.BEGIN = packdata[0]
	def.MsgTypeId = comm.BinaryHelper.ToInt16(packdata, 1)
	def.Id = comm.BinaryHelper.ToInt16(packdata, 3)
	def.Length = comm.BinaryHelper.ToInt32(packdata, 5)
	def.PackagesProperty = comm.BinaryHelper.ToInt32(packdata, 9)
	def.Valid = packdata[14]
	def.END = packdata[len(packdata)-1]
	obj = def
*/
func (p protocol808) Parse(packdata []byte) (obj interface{}, err error) {
	defer func() {
		conn.MyRecover()
	}()
	data,err :=comm.BinaryHelper.Byte808Descape(packdata,0,len(packdata))
	if err!=nil{
		return nil,err
	}
	iDataLength :=len(data)
	if(iDataLength==0){
		return nil,errors.New("终端上传数据长度错误,为空包")
	}
	iBodyLength:=(data[3] + ((data[2] & 0x03) << 8))//协议消息头的消息体属性中解析出的消息体长度
	bPartMsg := (data[2] & 0x20) != 0

	if bPartMsg{
		if iBodyLength!=byte(iDataLength - 17){
			return nil,errors.New("终端上传数据长度错误")
		}
	}else{
		if iBodyLength!=byte(iDataLength - 13){
			return nil,errors.New("终端上传数据长度错误")
		}
	}

	/*CRC Check*/
	if !(comm.BinaryHelper.CRCCheck(data)){
		return nil,errors.New("终端上传数据CRC错误")
	}
	

	MsgId := ((data[0] << 8) + data[1])
	SimNum:= comm.BinaryHelper.ToBCDString(data,4,6)
	MsgSN :=comm.BinaryHelper.ToInt16(data,10)
	var TotalMsgCount,CurrMsgOrder int16
	var msgBody []byte
	if bPartMsg{
		TotalMsgCount=comm.BinaryHelper.ToInt16(data,12)
		CurrMsgOrder=comm.BinaryHelper.ToInt16(data,14)
		msgBody = comm.BinaryHelper.CloneRange(data,16,(int32)(iBodyLength))
		log.Println("终端:",SimNum," MsgId:",MsgId,"MsgSN:",MsgSN,"TotalMsgCount:",TotalMsgCount,"CurrMsgOrder:",CurrMsgOrder,"Data:",comm.BinaryHelper.ToBCDString(packdata,0,int32(len(packdata))))
	}else{
		msgBody = comm.BinaryHelper.CloneRange(data,12,(int32)(iBodyLength))
		log.Println("终端:",SimNum," MsgId:",MsgId,"MsgSN:",MsgSN,comm.BinaryHelper.ToBCDString(packdata,0,int32(len(packdata))))
	}
	head :=NewMsgHead(int(MsgId),int(MsgSN),SimNum,int(TotalMsgCount),int(CurrMsgOrder))
	log.Println("Method Name:",head.GetMethodName())
	obj,err=_JTB808ParseBase.Parse(head,msgBody)
	// def := conn.DefaultTcpData{}
	// def.BEGIN = packdata[0]
	// def.MsgTypeId = comm.BinaryHelper.ToInt16(packdata, 1)
	// def.Id = comm.BinaryHelper.ToInt16(packdata, 3)
	// def.Length = comm.BinaryHelper.ToInt32(packdata, 5)
	// def.PackagesProperty = comm.BinaryHelper.ToInt32(packdata, 9)
	// def.Valid = packdata[14]
	// def.END = packdata[len(packdata)-1]
	// obj = def
	return obj, err
}
