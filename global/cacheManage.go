package global

import (
	"log"

	"fmt"

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

	//车辆基础信息缓存
	tmpVehicleInfoCache := &CacheBase{}
	tmpVehicleInfoCache.NewCache("VehicleInfoCache", 60, MSVehiclesCacheLoader{}.Load)
	VehiclesCache = tmpVehicleInfoCache
	cacheMana.TimerManage.RegistTask(tmpVehicleInfoCache.TimerTask)

	//从链路clients 缓存列表
	// tmpSubClisCache := &CacheBase{}
	// tmpSubClisCache.NewCache("SubCliCache", 60, SubClisCacheLoader{}.Load)
	// SubCliCache = tmpSubClisCache
	// cacheMana.TimerManage.RegistTask(tmpSubClisCache.TimerTask)

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

/*车辆基本信息*/
type VehicleInfo struct {
	PlateNum        string //SIM卡号
	ColorCode       int    //车牌号码
	SimNum          string //车牌颜色
	VehicleTypeCode string //终端ID

	TerminalId string //车辆类型

	TerminalTypeId  string //终端类型ID
	ProtocolName    string //协议名称
	ProtocolVersion string //终端协议版本
	ADRVersion      string //行驶记录仪版本
	RegFlag         int    //注册码

	AuthCode string //鉴权码

	LimitSpeed int // 车辆配置的限制速度

	VehicleOperateState int //车辆营运状态 1：营运 2：派单 3.维修 4.停运 5.报停
	OperatorId          int //运营商编号 (809/自有协议接入使用)
}

func (e VehicleInfo) Key() string {
	return fmt.Sprintf("%v%v", e.PlateNum, e.ColorCode)
}

/*MSVehiclesCacheLoader 终端车辆加载器*/
type MSVehiclesCacheLoader struct {
	CacheBase
}

func (cache MSVehiclesCacheLoader) Load(p interface{}) {
	defer func() {
		if p := recover(); p != nil {
			Error("VehiclesCache Recover:%v", p)
		}
	}()
	sql := `SELECT PlateNum ,a.ColorCode,a.SimNum,a.VehicleTypeCode,ISNULL(a.TerminalId,'') TerminalId,a.TerminalTypeId,p.ProtocolName,p.ProtocolVersion,
(case when a.RecorderVersion IS null or a.RecorderVersion='' then '2012' else a.RecorderVersion end) as ADRVersion,RegFlag,ISNULL(AuthCode,'') AuthCode,ISNULL(c.LimitSpeed,0) LimitSpeed ,ISNULL(c.VehicleOperateState,0) VehicleOperateState,a.OperatorId 
FROM bas_Vehicle a with(NOLOCK) 
INNER JOIN bas_TerminalType t ON a.TerminalTypeId=t.TerminalTypeId
INNER JOIN bas_Protocol p ON t.ProtocolCode=p.ProtocolCode
LEFT OUTER JOIN bas_VehicleConfig c with(NOLOCK) ON a.VehicleId = c.VehicleId;` //,AssociateUser,PlatformId
	rows, err := DBInstance().Query(sql)
	if err != nil {
		Error("VehiclesCache GetData Error:%v", err)
	}
	for rows.Next() {
		info := &VehicleInfo{}
		err = rows.Scan(&info.PlateNum, &info.ColorCode, &info.SimNum, &info.VehicleTypeCode, &info.TerminalId, &info.TerminalTypeId, &info.ProtocolName, &info.ProtocolVersion,
			&info.ADRVersion, &info.RegFlag, &info.AuthCode, &info.LimitSpeed, &info.VehicleOperateState, &info.OperatorId)
		if err != nil {
			Error("VehiclesCache Scan Row Error:%v", err)
			continue
		}
		VehiclesCache.AddCache(info.Key(), info)
		//fmt.Println(info, string(info.CompanyName))
	}
	var sKey string
	VehiclesCache.CacheValue.Purge(sKey, 120) //清理超时
	Info("车辆基本信息缓存 Load Cache Size:%d", len(VehiclesCache.CacheValue.DataStore))
}

/*从链路缓存器*/
type SubClisCacheLoader struct {
}
