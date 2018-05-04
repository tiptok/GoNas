package dbcore

import (
	"database/sql"
	"time"

	"bytes"

	"path/filepath"

	"fmt"

	"io/ioutil"
	"os"

	"github.com/tiptok/GoNas/global"
	"github.com/tiptok/gotransfer/comm"
)

var DBConn *sql.DB

type DBReceive struct {
	Rec             chan DBTask
	Err             chan string
	SleepInterval   time.Duration
	DBSessionName   string
	OnceTaskSize    int //一次入库任务大小 默认1000条
	OnceErrTaskSize int
}

func NewDBReceive(sName string) *DBReceive {
	return &DBReceive{
		Rec:             make(chan DBTask, 10000),
		Err:             make(chan string, 1000),
		SleepInterval:   1000,
		OnceTaskSize:    10, //修改长度 10
		OnceErrTaskSize: 10,
		DBSessionName:   sName,
	}
}
func (db *DBReceive) Start() {
	go db.OnDBTaskWork()
	go db.OnDBErrTaskWork()
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
			global.Debug("%s 执行sql数目:%d", db.DBSessionName, taskSize+1)
			db.ExecDBTask(taskBuffer) //执行sql
		default:
			time.Sleep(db.SleepInterval)
		}
	}
}

func (db *DBReceive) OnDBErrTaskWork() {
	for {
		select {
		case task, isClose := <-db.Err: //ErrSizeOnSave
			buf := bytes.NewBuffer(nil)
			buf.WriteString(task)
			for i := 1; i <= db.OnceErrTaskSize; i++ {
				if task, isClose = <-db.Err; !isClose {
					global.Debug("DB Err Task Close:%v", task)
					return
				}
				buf.WriteString(task)
			}
			result, sfile := getfilePath(db)
			/*执行错误的脚本 存储到文件中*/
			if result {
				err := savefile(sfile, buf.Bytes())
				if err != nil {
					global.Error("OnDBErrTaskWork Work Error:%v", err)
				}
			}
			buf.Reset()    //清空
			time.Sleep(50) //睡眠50ms
		// case task, isClose := <-db.Err:
		// 	if !isClose {
		// 		global.Debug("DB Err Task Close:%v", task)
		// 		return
		// 	}
		// 	result, sfile := getfilePath(db)
		// 	/*执行错误的脚本 存储到文件中*/
		// 	if result {
		// 		err := savefile(sfile, []byte(task))
		// 		if err != nil {
		// 			global.Error("OnDBErrTaskWork Work Error:%v", err)
		// 		}
		// 	}
		// 	time.Sleep(50) //睡眠50ms
		default:
			time.Sleep(db.SleepInterval)
		}
	}
}

//getfilePath 获取文件路径
func getfilePath(db *DBReceive) (bool, string) {
	result := false
	if f, exist := comm.UnityToolHelper.FileExist(global.Param.CachePath); exist {
		if !f.IsDir() {
			return result, ""
		}
		pathN := filepath.Join(global.Param.CachePath, db.DBSessionName)
		if _, exist = comm.UnityToolHelper.FileExist(pathN); !exist {
			err := comm.UnityToolHelper.MKdir(pathN) //新建下级路径 ./PosTask
			if err != nil {
				global.Error("DBReceive getfilePath:%v", err)
			}
		}
		pathN = filepath.Join(pathN, time.Now().Format("2006-01-02"))
		if _, exist = comm.UnityToolHelper.FileExist(pathN); !exist {
			err := comm.UnityToolHelper.MKdir(pathN) //新建下级路径 ./PosTask/2018-04-25
			if err != nil {
				global.Error("DBReceive getfilePath:%v", err)
			}
		}
		result = true
		sfilePath := filepath.Join(pathN, fmt.Sprintf("%d.dat", time.Now().UnixNano()))
		return result, sfilePath
	}
	return result, ""
}

func savefile(sfilePath string, data []byte) error {
	// result, sfile := getfilePath(db)
	// /*执行错误的脚本 存储到文件中*/
	// if result {
	// 	err := savefile(sfile, []byte(task))
	// 	if err != nil {
	// 		global.Error("OnDBErrTaskWork Work Error:%v", err)
	// 	}
	// }
	return ioutil.WriteFile(sfilePath, data, os.ModePerm)
}

func (db *DBReceive) ExecDBTask(buf *bytes.Buffer) {
	defer func() {
		if err := recover(); err != nil {
			global.Error("ExecDBTask Error:%v", err)
			db.Err <- buf.String()
		}
	}()

	reuslt, err := global.DBInstance().Exec(buf.String())
	global.Error("ExecDBTask:%v%v", buf.String())
	if err != nil {
		global.Error("ExecDBTask Error:%v%v", err, buf.String())
		db.Err <- buf.String()
	} else {
		count, err := reuslt.RowsAffected()
		if err != nil {
			global.Error("ExecDBTask Error:%v", err)
		}
		global.Debug("%s 执行sql影响数目:%d", db.DBSessionName, count)
	}
}

//DBTask  数据库任务
type DBTask struct {
	Sql   string
	MsgId int
}
