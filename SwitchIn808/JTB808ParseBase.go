package SwitchIn808

import (
	//"log"
	"reflect"
	"errors"
	"github.com/tiptok/gotransfer/comm"
)

var (
	PARSEPACKAGE_ERROR = errors.New("解包异常")
)
type JTB808ParseBase struct{

}
func(p *JTB808ParseBase)Parse(h MsgHead,msgBody []byte)(interface{},error){
	sMethodName := h.GetMethodName()
	aRefV :=[]reflect.Value{reflect.ValueOf(msgBody),reflect.ValueOf(h)}
	method :=reflect.ValueOf(p).MethodByName(sMethodName)
	// if method=nil{
	// 	return nil,errors.New("NOT EXiSTS Method :"+sMethodName)
	// }
	rsp:=method.Call(aRefV)
	return rsp[0].Interface(),nil
}
/*
	解析终端通用应答 0x0001
*/
func(p *JTB808ParseBase) J0001(msgBody []byte,h MsgHead)(interface{}){
	outEntity := TermCommonReply{}
	outEntity.RspMsgSN =uint16(comm.BinaryHelper.ToInt16(msgBody,0)) //应答流水号 
	outEntity.RspMsgId = uint16(comm.BinaryHelper.ToInt16(msgBody,2))//应答ID
	outEntity.RspResult = int(msgBody[4]);//结果 0：成功/确认；1：失败；2：消息有误；3：不支持
	
	/*是否需要应答*/
	outEntity.IsNeedRsp = false
	outEntity.DownRspEntity = nil
	return outEntity
}
