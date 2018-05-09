package core

import (
	"io/ioutil"
	"os"
	"time"

	"sync"

	"path/filepath"

	"github.com/tiptok/GoNas/dbcore"
	"github.com/tiptok/GoNas/global"
	"github.com/tiptok/GoNas/model"
	"github.com/tiptok/gotransfer/comm"
)

/*数据接收*/
type MSDBHandler struct {
	PosRec *dbcore.DBReceive
}

//UpData 接收数据
func (db MSDBHandler) UpData(rcv model.IEntity) {
	base := rcv.GetEntityBase()
	cmdcode := base.CmdCode().(uint16)
	posTask := dbcore.DBTask{Sql: rcv.GetDBSql()}
	switch cmdcode {
	case 0x1202: //添加任务
		//dbcore.DBTask{Sql: "insert into biz_DriverLogDetailed(SimNum,DriverId,DriverName,OrgName,[State],ReceTime,UpDataTime,WorkLicenseId)Values('1','2','3','4',5,GETDATE(),'2018-04-20 16:10:16','6');"}
		db.PosRec.Rec <- posTask
	case 0x1203:
		db.PosRec.Rec <- posTask
	default:
		//global.Debug("未识别上报数据类型:%x", cmdcode)
	}
}

//NewMSDBHandler new MSDBHandler
func NewMSDBHandler() MSDBHandler {
	wg := new(sync.WaitGroup)
	//启动位置线程
	posRec := dbcore.NewDBReceive("TaskPos")
	posRec.Start()
	wg.Add(1)
	func() {
		go OnDBErrTaskReWork(30) //错误文件重新入库 30s 执行一次
		wg.Done()
	}()
	wg.Wait()
	return MSDBHandler{
		PosRec: posRec,
	}
}

//OnDBErrTaskReWork 执行入库异常的文件 重新执行
//sec 执行间隔
func OnDBErrTaskReWork(sec int64) {
	for {
		dir, err := os.OpenFile(global.Param.CachePath, os.O_RDONLY, os.ModeDir)
		if err != nil {
			defer dir.Close()
			global.Error("OnDBErrTaskReWork:%v", err)
		} else {
			files, err := dir.Readdir(-1)
			if err != nil {
				global.Error("OnDBErrTaskReWork:%v", err)
			} else {
				for _, file := range files {
					if !file.IsDir() {
						sfilepath := filepath.Join(global.Param.CachePath, file.Name())
						data, err := ioutil.ReadFile(sfilepath)
						if err != nil {
							global.Error("OnDBErrTaskReWork :%v", err)
						} else {
							_, err = global.DBInstance().Exec(string(data)) //执行 sql
							if err != nil {
								global.Error("OnDBErrTaskReWork :%v", err)
							} else {
								os.Remove(sfilepath) //执行成功 移除文件
							}
						}
						if err != nil {
							/*执行错误写入错误文件*/
							efilepath := filepath.Join(global.Param.CachePath, "ErrorFile")

							if _, exist := comm.UnityToolHelper.FileExist(efilepath); !exist {
								comm.UnityToolHelper.MKdir(efilepath)
							}
							efilepath = filepath.Join(efilepath, file.Name())
							// ioutil.WriteFile(errfp, data, os.ModePerm)
							// os.Remove(filepath)
							os.Rename(sfilepath, efilepath) //移动文件夹 old -> new
						}
					}
				}
			}
		}
		time.Sleep(time.Duration(sec) * time.Second)
	}
}
