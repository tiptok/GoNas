package SwitchIn809

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"encoding/hex"

	"encoding/json"

	_ "github.com/alexbrainman/odbc"
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
	//time.Parse
}

func Test809PackerBase(t *testing.T) {
	p := &protocol809{}
	rspEntity := &model.UP_CONNECT_RSP{EntityBase: model.EntityBase{MsgId: model.J主链路登录应答}, Result: 0, Verify_Code: 12345678}
	rspEntity.AccessCode = "12345678"
	data, err := p.PacketMsg(rspEntity)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Packer Data:%s", hex.EncodeToString(data))
}

func TestDB(t *testing.T) {
	var DBConn *sql.DB
	DBConn, err := sql.Open("odbc", "driver={SQL Server};server=192.168.3.87;uid=sa;pwd=top@db123;database=TopDB;Connect Timeout=120;")
	if err != nil {
		log.Println("Get DBInstance Err:%v", err)
		DBConn.Close()
	}
	if DBConn == nil {
		log.Println("Init Error")
	}
	log.Println("Init DB %v", DBConn)
}

//登录
//5b000000480000000e100100bc614e010001000000000000bc614e31323334353637383132372e302e302e3100000000000000000000000000000000000000000000004671cc385d
//2018/04/10 16:06:03 MsgId:4101  MsgSN:3 AccessCode:12345678
//心跳
//2018/04/10 16:06:10 MsgId:4101  MsgSN:3 AccessCode:12345678

// 5b0000005a0200000aa112000000350b0100010000000000b6f5414c45313536000000000000000000000000000212020000002400110407e20a3b1106ccbea501cc71c700000000000098f70144000700000002000000008ce65d
// 2018/04/19 18:24:46.369 [D] MsgId:1202  MsgSN:2721 AccessCode:13579
// 2018/04/19 18:24:46.369 [D] 接收到实体&{{{4608 2721 4610 13579} 鄂ALE156 2 } {0 2018-04-17 10:59:17 +0800 CST 114081445 30175687 0 0 39159 324 7 2 0 0}}
