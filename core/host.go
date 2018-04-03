package core

import (
	"log"

	"github.com/tiptok/GoNas/SwitchIn808"
	"github.com/tiptok/GoNas/SwitchIn809"
	"github.com/tiptok/GoNas/global"
	"github.com/tiptok/gotransfer/conn"
)

type Host struct {
	NasServer conn.ITcpServer
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
	case "JTB809":
		h.NasServer = &SwitchIn809.Tcp809Server{}
	}
	global.UpHandler = &Up808Data{}
	if h.NasServer != nil {
		init = h.NasServer.Start()
	}
	log.Printf("Host Start %s Result:%v", protocol, init)
}
