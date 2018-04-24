package core

import (
	"github.com/tiptok/GoNas/dbcore"
	"github.com/tiptok/GoNas/global"
	"github.com/tiptok/GoNas/model"
)

/*数据接收*/
type MSDBHandler struct {
	PosRec *dbcore.DBReceive
}

func (db MSDBHandler) UpData(rcv model.IEntity) {
	base := rcv.GetEntityBase()
	cmdcode := base.CmdCode().(int16)
	posTask := dbcore.DBTask{Sql: rcv.GetDBSql()}
	switch cmdcode {
	case 0x1202: //添加任务
		//dbcore.DBTask{Sql: "insert into biz_DriverLogDetailed(SimNum,DriverId,DriverName,OrgName,[State],ReceTime,UpDataTime,WorkLicenseId)Values('1','2','3','4',5,GETDATE(),'2018-04-20 16:10:16','6');"}
		db.PosRec.Rec <- posTask
	default:
		global.Debug("未识别上报数据类型:%x", cmdcode)
	}
}

//NewMSDBHandler new MSDBHandler
func NewMSDBHandler() MSDBHandler {
	//启动位置线程
	posRec := dbcore.NewDBReceive("TaskPos")
	posRec.Start()
	return MSDBHandler{
		PosRec: posRec,
	}
}
