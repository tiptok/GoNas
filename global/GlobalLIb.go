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

//企业信息缓存
var PInfoCahce *CacheBase

//终端车辆信息缓存
var VehiclesCache *CacheBase

var SubCliCache *CacheBase
