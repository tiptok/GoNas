package SwitchIn808

import(
	"testing"
	"github.com/tiptok/gotransfer/comm"
	"log"
)

func TestParse0001(t *testing.T){
	//7E0001000506495944723706DD97A28300007D027E
	//7E0001000506495944723706DE97A3880100767E
	/*
	7E
	0001
	0005
	064959447237
	06DE
	97A3
	8801
	00
	76
	7E*/
	data,_ := comm.BinaryHelper.GetBCDString("7E0001000506495944723706DE97A3880100767E")
	protocol :=protocol808{}
	entity,err:=protocol.Parse(data)
	if err!=nil{
		//log.Println(err.Error())
		t.Log(err.Error())
	}
	if entity!=nil{
		t.Log(entity)
		//log.Println()
		log.Println(entity)
	}
	t.Log("end")
}