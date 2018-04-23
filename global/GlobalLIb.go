package global

import "github.com/tiptok/GoNas/model"

/*上行*/
var UpHandler IUpData

type IUpData interface {
	UpData(rcv model.IEntity)
}

/*下行*/
var DownHandler IDownData

type IDownData interface {
	DownData(rcv model.IEntity)
}
