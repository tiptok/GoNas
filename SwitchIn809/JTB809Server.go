package SwitchIn809

import (
	"github.com/tiptok/GoNas/global"
	"github.com/tiptok/gotransfer/conn"
)

type Tcp809Server struct {
	conn.TcpServerBase
}

func (svr *Tcp809Server) Start() bool {
	//启动tcp服务
	go func() {
		svr.Server.NewTcpServer(global.Param.ServerPort, 500, 500)
		svr.Server.Config.IsParsePartMsg = true //进行分包
		svr.Server.P = &protocol809{}
		svr.Server.Start(&SvrHander809{})
	}()
	return true
}

func (svr *Tcp809Server) Stop() bool {
	return true
}
