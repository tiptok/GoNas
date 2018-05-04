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
	CacheMana global.CacheManage
	//上行处理
	//分析处理
	//下行处理
	//分发处理
	//入库处理
}

func (h *Host) Start(protocol string) {
	var init bool = false

	//加载缓存管理
	h.CacheMana = global.CacheManage{}
	init = h.CacheMana.Init()

	switch protocol {
	case "JTB808":
		h.NasServer = &SwitchIn808.Tcp808Server{}
	case "JTB809":
		h.NasServer = &SwitchIn809.Tcp809Server{}
	}
	uphandler := &Up808Data{}
	uphandler.BizDB = NewMSDBHandler()
	global.UpHandler = uphandler
	//global.d
	if h.NasServer != nil {
		init = h.NasServer.Start()
	}
	log.Printf("Host Start %s Result:%v", protocol, init)
}
