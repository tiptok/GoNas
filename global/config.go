package global

import "github.com/astaxie/beego/config"

var Param *Params

//Params config item
type Params struct {
	/*******基本参数*******/
	Protocol   string //协议类型 :JTB808 JTB809
	ServerPort int    //服务监听端口

	/*******809配置项******/
	IsEncrypt  bool
	Key        int
	M1         int
	IA1        int
	IC1        int
	AccessCode []string
}

func init() {
	Param = &Params{}
	Param.LoadConfig("ini", "param.cnf")
}

//LoadConfig load config
func (p *Params) LoadConfig(pType string, fName string) *Params {
	defer func() {
		if p := recover(); p != nil {
			Error("LoadConfig 读取配置异常r! p: %v", p)
			//debug.PrintStack()
		}
	}()
	con, err := config.NewConfig(pType, fName)
	if err != nil {
		Error("LoadConfig 加载配置异常", err)
	}
	p.Protocol = con.String("Protocol")
	p.ServerPort, _ = con.Int("ServerPort")
	p.IsEncrypt, _ = con.Bool("IsEncrypt")
	p.Key, _ = con.Int("Key")
	p.M1, _ = con.Int("M1")
	p.IA1, _ = con.Int("IA1")
	p.IC1, _ = con.Int("IC1")
	p.AccessCode = con.Strings("AccessCode")
	return p
}
