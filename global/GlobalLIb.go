package global

/*上行*/
var UpHandler IUpData

type IUpData interface {
	UpData(rcv interface{})
}

/*下行*/
var DownHandler IDownData

type IDownData interface {
	DownData(rcv interface{})
}
