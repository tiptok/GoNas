package core

import (
	"log"

	"github.com/tiptok/GoNas/SwitchIn808"
	"github.com/tiptok/GoNas/global"
)

type Host struct {
	NasServer SwitchIn808.TcpServer
	//上行处理
	//分析处理
	//下行处理
	//分发处理
	//入库处理
}

func (h *Host) Start(protocol string) {
	var init bool = false
	switch protocol {
	case "JTB808":
		h.NasServer = &SwitchIn808.Tcp808Server{}
	}
	global.UpHandler = &Up808Data{}
	if h.NasServer != nil {
		init = h.NasServer.Start()
	}
	log.Printf("Host Start %s Result:%v", protocol, init)
}
