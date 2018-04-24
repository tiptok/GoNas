package model

import "time"

type IEntity interface {
	GetMsgId() interface{}
	GetEntityBase() *EntityBase
	GetDBSql() string
}
type EntityBase struct {
	MsgId interface{}
	//SimNum string
	MsgSN      int
	SubMsgId   interface{}
	AccessCode string
	ReceTime   time.Time
}

func (e *EntityBase) CmdCode() interface{} {
	if e.SubMsgId != nil {
		return e.SubMsgId
	}
	return e.MsgId
}
func (e *EntityBase) GetMsgId() interface{} {
	return e.MsgId
}
func (e *EntityBase) GetEntityBase() *EntityBase {
	return e
}

func (e *EntityBase) GetDBSql() string {
	return ""
}

func (e *EntityBase) SetEntity(in EntityBase) {
	e.MsgId = in.MsgId
	e.MsgSN = in.MsgSN
	e.AccessCode = in.AccessCode
	e.ReceTime = time.Now()
}
