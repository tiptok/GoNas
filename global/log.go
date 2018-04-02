package global

import "github.com/astaxie/beego/logs"

//var _log *log.Logger

func init() {
	//_log = logs.GetLogger("809")
	logs.SetLogger(logs.AdapterFile, `{"filename":"809.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
}

//Debug log debug
func Debug(f interface{}, v ...interface{}) {
	logs.Debug(f, v)
}

//Info  log info
func Info(f interface{}, v ...interface{}) {
	logs.Info(f, v)
}

//Warning log warning
func Warning(f interface{}, v ...interface{}) {
	logs.Warning(f, v)
}

//Error log error
func Error(f interface{}, v ...interface{}) {
	logs.Error(f, v)
}
