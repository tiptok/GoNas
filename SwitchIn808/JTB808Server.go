package SwitchIn808

import (
	"github.com/tiptok/gotransfer/conn"
)

/*统一Server接口*/
type TcpServer interface {
	Start() bool
	Stop() bool
}

//func ()xx

type TcpServerBase struct {
	Server conn.TcpServer
}

type Tcp808Server struct {
	TcpServerBase
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
