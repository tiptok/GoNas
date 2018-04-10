package global

import "github.com/astaxie/beego/logs"

var _log *logs.BeeLogger

func init() {
	_log = logs.NewLogger()
	_log.SetLogger(logs.AdapterFile, `{"filename":"809.log","level":7,"maxlines":0,"maxsize":2097152,"daily":true,"maxdays":10}`)
}

//Debug log debug
func Debug(f string, v ...interface{}) {
	_log.Debug(f, v...)
}

//Info  log info
func Info(f string, v ...interface{}) {
	_log.Info(f, v...)
}

//Warning log warning
func Warning(f string, v ...interface{}) {
	_log.Warning(f, v...)
}

//Error log error
func Error(f string, v ...interface{}) {
	_log.Error(f, v...)
}
