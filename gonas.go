package main

import (
	"runtime"

	"github.com/tiptok/GoNas/core"
)

var host core.Host
var exit chan int

func main() {
	defer func() {
		exit <- 1 //异常退出
	}()
	runtime.GOMAXPROCS(runtime.NumCPU())
	host = core.Host{}
	host.Start("JTB808")
	//等待退出
	<-exit
}
