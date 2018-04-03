package SwitchIn808

import "github.com/tiptok/gotransfer/conn"

type Tcp808Server struct {
	conn.TcpServerBase
}

func (svr *Tcp808Server) Start() bool {
	//启动tcp服务
	go func() {
		svr.Server.NewTcpServer(9927, 500, 500)
		svr.Server.P = &protocol808{}
		svr.Server.Start(&SvrHander808{})
	}()
	return true
}

func (svr *Tcp808Server) Stop() bool {
	return true
}
