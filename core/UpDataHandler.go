package core

import "github.com/tiptok/GoNas/model"

type Up808Data struct {
	BizDB MSDBHandler
}

func (u Up808Data) UpData(rcv model.IEntity) {
	//global.Debug("Up808Data UpData:%v", rcv)
	u.BizDB.UpData(rcv)
}

// func NewUp808Data(){
// 	up:=Up808Data{}
// 	up.BizDB =
// }
