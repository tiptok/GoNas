package SwitchIn809

import (
	"github.com/tiptok/GoNas/global"
	"github.com/tiptok/gotransfer/comm"
	"github.com/tiptok/gotransfer/conn"
)

type Tcp809Server struct {
	Server             conn.TcpServer
	SubList            *comm.DataContext
	SubClientExpireChk *comm.TimerWork
}

func (svr *Tcp809Server) Start() bool {
	//启动tcp服务
	go func() {
		svr.Server.NewTcpServer(global.Param.ServerPort, 500, 500)
		svr.Server.Config.IsParsePartMsg = true //进行分包
		svr.Server.P = &protocol809{}
		svr.SubList = &comm.DataContext{}

		svr.SubClientExpireChk = comm.NewTimerWork()
		taskChkExpire := &comm.Task{
			Interval: 60,
			TaskId:   "ClientExpireCheck",
			Work:     svr.CheckClientExpire,
		}
		svr.SubClientExpireChk.RegistTask(taskChkExpire)

		svr.Server.Start(svr)
		svr.SubClientExpireChk.Start()
	}()
	return true
}

func (svr *Tcp809Server) Stop() bool {
	svr.SubClientExpireChk.Stop()
	return true
}

//CheckClientExpire 检查从链路是否连接超时(无数据交互)
func (svr *Tcp809Server) CheckClientExpire(val interface{}) {
	defer func() {
		if p := recover(); p != nil {
			global.Error("SubClisCacheLoader Load Recover:%v", p)
		}
	}()
	var key string
	svr.SubList.PurgeWithFunc(key, 300, svr.OnRemvoe)
}

func (svr *Tcp809Server) OnRemvoe(k interface{}, val interface{}) {
	defer func() {
		if p := recover(); p != nil {
			global.Error("SubClisCacheLoader OnRemvoe Recover:%v", p)
		}
	}()
	subCli := val.(*TcpSubClient)
	if subCli != nil {
		global.Info("从链路 %v 执行超时处理", subCli.AccessCode)
		subCli.Stop() //超时关闭  做从链路关闭操作
	}
}
