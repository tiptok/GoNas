package global

import "github.com/astaxie/beego/config"
import "log"

var Param *Params

//Params config item
type Params struct {
	/*******基本参数*******/
	Protocol     string //协议类型 :JTB808 JTB809
	ServerPort   int    //服务监听端口
	DBConnString string //数据库连接串
	TrackDBName  string //轨迹库名称
	CachePath    string //文件缓存路径

	/*******809配置项******/
	IsEncrypt  bool
	Key        int
	M1         int
	IA1        int
	IC1        int
	AccessCode []string
	VerifyCode int
}

func init() {
	Param = &Params{}
	Param.LoadConfig("ini", "param.conf")
}

//LoadConfig load config
func (p *Params) LoadConfig(pType string, fName string) *Params {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("LoadConfig 读取配置异常r! p: %v", p)
			//debug.PrintStack()
		}
	}()
	con, err := config.NewConfig(pType, fName)
	if err != nil {
		log.Printf("LoadConfig 加载配置异常 e: %v", err)
	}
	p.Protocol = con.String("Protocol")
	p.ServerPort, _ = con.Int("ServerPort")
	p.DBConnString = con.String("DBConnString")
	p.TrackDBName = con.String("TrackDBName")
	p.CachePath = con.String("CachePath")
	p.IsEncrypt, _ = con.Bool("IsEncrypt")
	p.Key, _ = con.Int("Key")
	p.M1, _ = con.Int("M1")
	p.IA1, _ = con.Int("IA1")
	p.IC1, _ = con.Int("IC1")
	p.AccessCode = con.Strings("AccessCode")
	p.VerifyCode, _ = con.Int("VerifyCode")
	log.Printf("Load Config:%v\n", *p)
	return p
}
