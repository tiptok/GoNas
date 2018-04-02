package global

import (
	"log"

	"github.com/astaxie/beego/config"
)

//Params config item
type Params struct {

	/*******809配置项******/
	IsEncrypt  bool
	Key        int
	M1         int
	IA1        int
	IC1        int
	AccessCode []string
}

//LoadConfig load config
func (p *Params) LoadConfig(pType string, fName string) *Params {
	con, err := config.NewConfig(pType, fName)
	if err != nil {
		log.Println("读取配置异常", err)
	}
	p.IsEncrypt, _ = con.Bool("IsEncrypt")
	p.Key, _ = con.Int("Key")
	p.M1, _ = con.Int("M1")
	p.IA1, _ = con.Int("IA1")
	p.IC1, _ = con.Int("IC1")
	p.AccessCode = con.Strings("AccessCode")
	return p
}
