package SwitchIn809

import (
	"fmt"
	"testing"

	"encoding/hex"

	"encoding/json"

	"github.com/tiptok/GoNas/model"
)

func TestByte809Descape(t *testing.T) {
	data := []byte{0x5b, 0x5a, 0x01, 0x5e, 0x02, 0x03, 0x5e, 0x01, 0x5a, 0x02, 0x5b}
	fmt.Println("srcData:", hex.EncodeToString(data))
	desData, err := Byte809Descape(data, 0, len(data))
	if err != nil {
		t.Log(err)
	}
	fmt.Println("desData:", hex.EncodeToString(desData))
}

func TestByte809Enscape(t *testing.T) {
	data := []byte{0x5b, 0x5e, 0x03, 0x5d, 0x5a}
	fmt.Println("srcData:", hex.EncodeToString(data))
	desData := Byte809Enscape(data, 0, len(data))
	fmt.Println("desData:", hex.EncodeToString(desData))
	outdata := []byte{0x5b, 0x5a, 0x01, 0x5e, 0x02, 0x03, 0x5e, 0x01, 0x5a, 0x02, 0x5b}
	fmt.Println("outdata:", hex.EncodeToString(outdata))
}

func TestInherit(t *testing.T) {
	a_a := &model.Entity809_a{}
	a_a.MsgId = "1001"
	a_a.SimNum = "18860183050"
	a_a.Result = "0"
	a_a.Message = "Hello 2018"
	fmt.Println(a_a)

	var i interface{} = a_a

	if a, ok := i.(model.IEntity); ok {
		fmt.Println(a)
		jd, err := json.Marshal(a)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Json Data:", string(jd))
		}
	}

	// a := *model.EntityBase(a_a)
	// fmt.Println("EntityBase", a)

	// a := model.Entity809_a(a_a)
	// fmt.Println(a)
}

//登录
//5b000000480000000e100100bc614e010001000000000000bc614e31323334353637383132372e302e302e3100000000000000000000000000000000000000000000004671cc385d
2018/04/10 16:06:03 MsgId:4101  MsgSN:3 AccessCode:12345678
//心跳
//2018/04/10 16:06:10 MsgId:4101  MsgSN:3 AccessCode:12345678