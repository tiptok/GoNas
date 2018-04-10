package model

//1001 主链路登录请求
type UP_CONNECT_REQ struct {
	EntityBase
	USERNAME       string //用户名 作为用户ID的补充（有些809协议吧UserId定义为串型）
	USERID         uint32 //用户ID
	PASSWORD       string //用户密码
	DOWN_LINK_IP   string //从链路IP
	DOWN_LINK_PORT uint32 //从链路端口
}

//1005 主链路连接保持请求
type UP_LINKTEST_REQ struct {
	EntityBase
}

//1008 下级平台主动关闭链路通知
type UP_CLOSELINK_INFORM struct {
	EntityBase
	Reason_Code byte //链路关闭原因
}

//IsValid valid UP_CLOSELINK_INFORM
func (e *UP_CLOSELINK_INFORM) IsValid(args ...interface{}) (errMsg string) {
	if e.Reason_Code == 0x00 {
		errMsg = "网关重启"
	} else if e.Reason_Code == 0x01 {
		errMsg = "其他原因"
	} else {
		errMsg = "其他原因"
	}
	return errMsg
}
