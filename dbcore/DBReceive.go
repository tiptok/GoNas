package dbcore

import (
	"database/sql"
	"time"

	"bytes"

	"github.com/tiptok/GoNas/global"
)

var DBConn *sql.DB

type DBReceive struct {
	Rec           chan DBTask
	Err           chan string
	SleepInterval time.Duration
	DBSessionName string
	OnceTaskSize  int //一次入库任务大小 默认1000条
}

func NewDBReceive(sName string) *DBReceive {
	return &DBReceive{
		Rec:           make(chan DBTask, 10000),
		Err:           make(chan string, 1000),
		SleepInterval: 1000,
		OnceTaskSize:  1000,
		DBSessionName: sName,
	}
}
func (db *DBReceive) Start() {
	go db.OnDBTaskWork()
	//go db.OnDBErrTaskWork()
}
func (db *DBReceive) OnDBTaskWork() {
	for {
		select {
		case task, isClose := <-db.Rec:
			if !isClose {
				global.Debug("%v Close Rec Chan", db.DBSessionName)
				return
			}
			taskBuffer := bytes.NewBuffer(nil)
			taskBuffer.WriteString(task.Sql)
			taskBuffer.WriteString("\n")
			taskSize := 0 //一次大小
			if len(db.Rec) < db.OnceTaskSize {
				taskSize = len(db.Rec)
			} else {
				taskSize = db.OnceTaskSize
			}

			for i := 1; i <= taskSize; i++ {
				task, isClose = <-db.Rec
				taskBuffer.WriteString(task.Sql)
				taskBuffer.WriteString("\n")
			}
			//global.Debug("执行sql:%s",taskBuffer.String())
			db.ExecDBTask(taskBuffer) //执行sql
		default:
			time.Sleep(db.SleepInterval)
		}
	}
}

func (db *DBReceive) OnDBErrTaskWork() {
	for {
		select {
		case task, isClose := <-db.Rec:
			if !isClose {
				global.Debug("DB Err Task %v", task)
			}
		default:
			time.Sleep(db.SleepInterval)
		}
	}
}

func (db *DBReceive) ExecDBTask(buf *bytes.Buffer) {
	defer func() {
		if err := recover(); err != nil {
			global.Error("ExecDBTask Error:%v", err)
			db.Err <- buf.String()
		}
	}()

	_, err := global.DBInstance().Exec(buf.String())
	if err != nil {
		global.Error("ExecDBTask Error:%v%v", err, buf.String())
		db.Err <- buf.String()
	}
}

type DBTask struct {
	Sql   string
	MsgId int
}
