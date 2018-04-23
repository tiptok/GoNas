package SwitchIn809

import (
	"bytes"

	"github.com/tiptok/GoNas/model"
	"github.com/tiptok/gotransfer/comm"
)

type JTB809PackerBase struct {
}

/*
   J1002 主链路登录应答
*/
func (p *JTB809PackerBase) J1002(obj interface{}) (packdata []byte, err error) {
	buf := bytes.NewBuffer(nil)
	inEntity := obj.(*model.UP_CONNECT_RSP)
	buf.WriteByte(inEntity.Result)
	buf.Write(comm.BinaryHelper.Int32ToBytes(int(inEntity.Verify_Code)))
	return buf.Bytes(), nil
}
