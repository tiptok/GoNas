package global

import (
	"log"

	"github.com/tiptok/gotransfer/comm"
)

//CacheManage 缓存管理
type CacheManage struct {
	TimerManage *comm.TimerWork
}

//初始化
func (cacheMana CacheManage) Init() (result bool) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("CacheManage Init panic recover! p: %v", p)
			//debug.PrintStack()
			result = false
		}
	}()
	cacheMana.TimerManage = comm.NewTimerWork()

	//企业信息缓存
	tmpPlatformInfoCahce := &CacheBase{}
	tmpPlatformInfoCahce.NewCache("PlatformInfoCahce", 60, MSPlatformInfoCacheLoader{}.Load)
	PInfoCahce = tmpPlatformInfoCahce
	cacheMana.TimerManage.RegistTask(tmpPlatformInfoCahce.TimerTask)

	//启动缓存定时更新
	cacheMana.TimerManage.Start()
	result = true
	return result
}

/*MSPlatformInfo  企业信息*/
type MSPlatformInfo struct {
	CompanyId   int
	CompanyName string //公司名
	AccessCode  string //接入码
	CompanyIP   string //接入公司Ip
	UserId      string //用户名
	Password    string //密码
}

/*MSPlatformInfoCacheLoader 企业信息缓存加载器*/
type MSPlatformInfoCacheLoader struct {
	CacheBase
}

func (cache MSPlatformInfoCacheLoader) Load(p interface{}) {
	defer func() {
		if p := recover(); p != nil {
			Error("GetPlatformInfoCahce Recover:%v", p)
		}
	}()
	sql := "select CompanyId,CompanyName,AccessCode,CompanyIP,UserId,[Password] from biz_809CompanysInfoManager" //,AssociateUser,PlatformId
	rows, err := DBInstance().Query(sql)
	if err != nil {
		Error("MSPlatformInfoCahce GetData Error:%v", err)
	}
	for rows.Next() {
		info := &MSPlatformInfo{}
		err = rows.Scan(&info.CompanyId, &info.CompanyName, &info.AccessCode, &info.CompanyIP, &info.UserId, &info.Password)
		if err != nil {
			Error("MSPlatformInfoCahce Scan Row Error:%v", err)
			continue
		}
		PInfoCahce.AddCache(info.AccessCode, info)
		//fmt.Println(info, string(info.CompanyName))
	}
	Info("企业信息缓存 Load Cache Size:%d", len(PInfoCahce.CacheValue.DataStore))
}
