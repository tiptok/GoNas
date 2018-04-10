package model

type IEntity interface {
	GetMsgId() string
	GetEntityBase() *EntityBase
}
type EntityBase struct {
	MsgId string
	//SimNum string
	MsgSN      int
	SubMsgId   string
	AccessCode string
}

func (e *EntityBase) GetMsgId() string {
	return e.MsgId
}
func (e *EntityBase) GetEntityBase() *EntityBase {
	return e
}

func (e *EntityBase) SetEntity(in EntityBase) {
	e.MsgId = in.MsgId
	e.MsgSN = in.MsgSN
	e.AccessCode = in.AccessCode
}

type Entity809_a struct {
	EntityBase
	Message string
}
