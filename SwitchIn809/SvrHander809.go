package SwitchIn809

import (
	"encoding/hex"
	"log"

	"strings"

	"fmt"

	"strconv"

	"github.com/tiptok/GoNas/global"
	"github.com/tiptok/GoNas/model"
	"github.com/tiptok/gotransfer/comm"
	"github.com/tiptok/gotransfer/conn"
)

// type SvrHander809 struct {
// 	conn.TcpServerBase
// }

//连接事件
func (trans *Tcp809Server) OnConnect(c *conn.Connector) bool {
	defer func() {
		//conn.MyRecover()
	}()
	log.Println(c.RemoteAddress, "On Connect.")
	return true
}

//断开事件
func (trans *Tcp809Server) OnClose(c *conn.Connector) {
	log.Println(c.RemoteAddress, "On Close.")
}

//接收事件
func (trans *Tcp809Server) OnReceive(c *conn.Connector, d conn.TcpData) bool {
	var bUpData bool = true
	global.Debug(global.F(global.TCP, global.SVR809, "%v On Receive Data : %v"), c.RemoteAddress, hex.EncodeToString(d.Bytes()))
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
		cmdcode := entity.MsgId.(uint16)
		if entity.SubMsgId != nil && entity.SubMsgId.(uint16) != 0 {
			cmdcode = entity.SubMsgId.(uint16)
		}
		global.Debug(global.F(global.TCP, global.SVR809, "MsgId:%X  MsgSN:%d AccessCode:%v"), cmdcode, entity.MsgSN, entity.AccessCode)
		switch cmdcode {
		case model.J主链路登录请求:
			login := obj.(*model.UP_CONNECT_REQ)
			result, errMsg := chkPlatInfo(login)
			if result {
				global.Info(global.F(global.TCP, global.SVR809, "主链路登录结果:%v %v"), result, login.String())
				rspEntity = &model.UP_CONNECT_RSP{EntityBase: model.EntityBase{MsgId: model.J主链路登录应答}, Result: 0, Verify_Code: int32(global.Param.VerifyCode)}
				//添加到从链路缓存
				if _, isExists := trans.SubList.GetOk(login.AccessCode); !isExists {
					subCli := NewTcpSubClient(login, trans)
					trans.SubList.Set(login.AccessCode, subCli)
				}
			} else {
				global.Info("主链路登录失败 %v 错误:%v", login.String(), errMsg)
			}
			// case model.主链路注销请求:
		case model.J主链路连接保持请求:
			global.Info(global.F(global.TCP, global.SVR809, "收到 %v %v 主链路连接保持请求"), entity.AccessCode, c.RemoteAddress)
		case model.J实时上传车辆定位信息:
			bUpData = false
			global.UpHandler.UpData((obj.(*model.UP_EXG_MSG_REAL_LOCATION)).GetConvEntity())
			//global.Debug("接收到实体%v", obj)
		case model.J车辆定位信息自动补报:
			bUpData = false
			hisLocation := obj.(*model.UP_EXG_MSG_HISTORY_LOCATION)
			for _, val := range hisLocation.GNSS_DATA_LIST {
				pos := &model.UP_EXG_MSG_REAL_LOCATION{
					UP_EXG_MSG: hisLocation.UP_EXG_MSG,
					GNSS_DATA:  val,
				}
				global.UpHandler.UpData(pos.GetConvEntity())
			}
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

//SendCmdAsync  异步发送指令
func SendCmdAsync(c *conn.Connector, e model.IEntity) {
	//IEntity
	data, err := conn.SendEntity(e, c)
	if err != nil {
		global.Error("SvrHander Send Entity Error:%v", err)
	} else {
		global.Debug("SvrHander Send Data:%s", comm.BinaryHelper.ToBCDString(data, 0, int32(len(data))))
	}
}

//检查主链路登录信息
func chkPlatInfo(req *model.UP_CONNECT_REQ) (result bool, errMsg string) {
	result = false
	obj := global.PInfoCahce.GetCache(req.AccessCode)
	if obj != nil {
		pCache := obj.(*global.MSPlatformInfo)
		if pCache.UserId != strconv.Itoa(int(req.USERID)) {
			errMsg = fmt.Sprintf("用户校验失败,正确用户:%v", pCache.UserId)
		} else if strings.Compare(pCache.Password, req.PASSWORD) != 0 {
			errMsg = fmt.Sprintf("密码校验失败,正确密码:%v", pCache.Password)
		} else if strings.Compare(pCache.CompanyIP, req.DOWN_LINK_IP) != 0 {
			errMsg = fmt.Sprintf("IP未认证,当前认证Ip:%v", pCache.CompanyIP)
		} else {
			result = true
		}
	} else {
		errMsg = fmt.Sprintf("未找到企业,接入码:%v", req.AccessCode)
	}
	return result, errMsg
}
