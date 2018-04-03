package main

import (
	"runtime"

	"github.com/tiptok/GoNas/core"
	"github.com/tiptok/GoNas/global"
)

var host core.Host
var exit chan int

func main() {
	defer func() {
		exit <- 1 //异常退出
	}()
	runtime.GOMAXPROCS(runtime.NumCPU())
	host = core.Host{}
	host.Start(global.Param.Protocol)
	//等待退出
	<-exit
}
