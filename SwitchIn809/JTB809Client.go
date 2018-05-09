package SwitchIn809

import (
	"encoding/hex"
	"log"

	"fmt"

	"github.com/tiptok/GoNas/global"
	"github.com/tiptok/GoNas/model"
	"github.com/tiptok/gotransfer/comm"
	"github.com/tiptok/gotransfer/conn"
)

type TcpSubClient struct {
	//conn.TcpServerBase
	Client         *conn.TcpClient
	IsLogin        bool
	PlatInfo       *global.MSPlatformInfo
	HeartTimerWork *comm.TimerWork
	VerifyCode     int32
	AccessCode     string
	SVR            *Tcp809Server
}

//NewTcpSubClient  新建一个从链路
func NewTcpSubClient(login *model.UP_CONNECT_REQ, svr *Tcp809Server) *TcpSubClient {
	subCli := &TcpSubClient{}
	subCli.HeartTimerWork = comm.NewTimerWork()
	subCli.SVR = svr

	task := &comm.Task{
		Interval: 10,
		TaskId:   fmt.Sprintf("%v-HeartTimerWork", login.DOWN_LINK_IP),
		Work:     subCli.ChKSubHeartBeart,
		Param:    nil,
	}

	subCli.HeartTimerWork.RegistTask(task)

	//启动tcp服务
	subCli.Client = &conn.TcpClient{}
	subCli.Client.NewTcpClient(login.DOWN_LINK_IP, int(login.DOWN_LINK_PORT), 500, 500)
	subCli.AccessCode = login.AccessCode
	subCli.VerifyCode = int32(global.Param.VerifyCode)
	subCli.Client.P = &protocol809{}
	subCli.Client.Config.IsParsePartMsg = true //进行分包
	subCli.Client.Start(subCli)
	subCli.HeartTimerWork.Start() //启动定时
	return subCli
}

//ChKSubHeartBeart  心跳检查
func (subCli *TcpSubClient) ChKSubHeartBeart(obj interface{}) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("SubClient ChKSubHeartBeart panic recover! p: %v", p)
		}
	}()
	if subCli.Client != nil && subCli.Client.Conn != nil {
		if !subCli.Client.Conn.IsConneted {
			subCli.Client.ReStart()
		}
		if !subCli.IsLogin {
			//发送从链路登录
			subLogin := &model.DOWN_CONNECT_REQ{
				EntityBase: model.EntityBase{
					AccessCode: subCli.AccessCode,
					MsgId:      model.J从链路连接请求,
				},
				VERIFY_CODE: subCli.VerifyCode,
			}
			SendCmdAsync(subCli.Client.Conn, subLogin)
		} else {
			//发送心跳
			heart := &model.DOWN_LINKTEST_REQ{
				EntityBase: model.EntityBase{
					AccessCode: subCli.AccessCode,
					MsgId:      model.J从链路连接保持请求,
				},
			}
			SendCmdAsync(subCli.Client.Conn, heart)
		}
	}
}

//连接事件
func (trans *TcpSubClient) OnConnect(c *conn.Connector) bool {
	defer func() {
		//conn.MyRecover()
	}()
	log.Println(global.F(global.TCP, global.SUB809, ""), c.RemoteAddress, "On Connect.")
	return true
}

//断开事件
func (trans *TcpSubClient) OnClose(c *conn.Connector) {
	//trans.SVR.SubList.Delete(trans.AccessCode)
	trans.IsLogin = false
	log.Println(global.F(global.TCP, global.SUB809, ""), c.RemoteAddress, "On Close.")
}

//接收事件
func (trans *TcpSubClient) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	var bUpData bool = true
	global.Debug(global.F(global.TCP, global.SUB809, "%v On Receive Data : %v"), c.RemoteAddress, hex.EncodeToString(d.Bytes()))
	defer func() {
		if p := recover(); p != nil {
			log.Printf("SubClient OnReceive panic recover! p: %v", p)
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
		cmdcode := entity.MsgId.(uint16)
		if entity.SubMsgId != nil && entity.SubMsgId.(uint16) != 0 {
			cmdcode = entity.SubMsgId.(uint16)
		}
		global.Debug(global.F(global.TCP, global.SUB809, "MsgId:%X  MsgSN:%d AccessCode:%v"), cmdcode, entity.MsgSN, entity.AccessCode)
		switch cmdcode {
		case model.J从链路连接应答:
			connReply := obj.(*model.DOWN_CONNECT_RSP)
			if connReply.Result == 0 {
				trans.IsLogin = true
			}
			global.Info(global.F(global.TCP, global.SUB809, "收到 %v %v 从链路连接应答 应答结果:%v"), entity.AccessCode, c.RemoteAddress, connReply.Result)
		case model.J从链路连接保持请求应答:
			trans.SVR.SubList.Refresh(trans.AccessCode) //刷新活动时间
			global.Info(global.F(global.TCP, global.SUB809, "收到 %v %v 从链路连接保持请求应答 "), entity.AccessCode, c.RemoteAddress)
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
		SendCmdAsync(c, rspEntity)
	}
	return true
}

func (trans *TcpSubClient) Stop() {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("SubClient Stop panic recover! p: %v", p)
		}
	}()
	trans.HeartTimerWork.Stop()                //心跳检查 停止
	trans.SVR.SubList.Delete(trans.AccessCode) //从从链路列表中移除
	trans.Stop()                               //Client Stop
}
