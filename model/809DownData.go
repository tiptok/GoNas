package model

//UP_CONNECT_RSP 1002 主链路登录应答
type UP_CONNECT_RSP struct {
	EntityBase
	Result      byte  //结果
	Verify_Code int32 //校验主链路登录请求
}

func (e *UP_CONNECT_RSP) GetMsgId() interface{} {
	return J主链路登录应答
}

//IsValid valid UP_CONNECT_RSP
func (e *UP_CONNECT_RSP) IsValid(args ...interface{}) (isLogin bool, errMsg string) {
	isLogin = false
	if e.Result == 0x00 {
		errMsg = "登录成功"
		isLogin = true
	} else if e.Result == 0x01 {
		errMsg = "登录失败：IP地址不正确"
	} else if e.Result == 0x02 {
		errMsg = "登录失败：接入码不正确"
	} else if e.Result == 0x03 {
		errMsg = "登录失败：用户没有注册"
	} else if e.Result == 0x04 {
		errMsg = "登录失败：密码错误"
	} else if e.Result == 0x05 {
		errMsg = "登录失败：资源紧张，稍后再连接"
	} else if e.Result == 0x06 {
		errMsg = "登录失败：其他原因"
	}
	return isLogin, errMsg
}

//UP_DISCONNECT_RSP  0x1004 主链路注销应答
type UP_DISCONNECT_RSP struct {
	EntityBase
}

//UP_LINKTEST_RSP 0x1006 主链路连接保持应答
type UP_LINKTEST_RSP struct {
	EntityBase
}

//DOWN_CONNECT_REQ 0x1009 从链路连接请求
type DOWN_CONNECT_REQ struct {
	EntityBase
	VERIFY_CODE int32 //校验码
}

//IsValid valid DOWN_CONNECT_REQ
func (e *DOWN_CONNECT_REQ) IsValid(args ...interface{}) (Result bool) {
	Result = args[0] == e.VERIFY_CODE
	return Result
}

//DOWN_CONNECT_RSP 0x9001 从链路连接应答
type DOWN_CONNECT_RSP struct {
	EntityBase
	Result byte
}

//DOWN_DISCONNECT_REQ 0x9003 从链路注销请求
type DOWN_DISCONNECT_REQ struct {
	EntityBase
	Verify_Code int32
}

//DOWN_DISCONNECT_RSP 0x9004 从链路注销应答
type DOWN_DISCONNECT_RSP struct {
	EntityBase
}

//DOWN_LINKTEST_RSP 0x9006 从链路连接保持请求应答
type DOWN_LINKTEST_RSP struct {
	EntityBase
}

//DOWN_DISCONNET_INFORM 0x9007 从链路断开通知
type DOWN_DISCONNET_INFORM struct {
	EntityBase
	Error_Code byte
}

//IsValid valid DOWN_CLOSELINK_INFORM
func (e *DOWN_DISCONNET_INFORM) IsValid(args ...interface{}) (errMsg string) {
	if e.Error_Code == 0x00 {
		errMsg = "无法连接下级平台指定服务的"
	} else if e.Error_Code == 0x01 {
		errMsg = "上级平台客户端与下级平台服务端断开"
	} else {
		errMsg = "其他原因"
	}
	return errMsg
}

//DOWN_CLOSELINK_INFORM 0x9008 上级平台主动关闭链路通知
type DOWN_CLOSELINK_INFORM struct {
	EntityBase
	Reason_Code byte
}

//IsValid valid DOWN_CLOSELINK_INFORM
func (e *DOWN_CLOSELINK_INFORM) IsValid(args ...interface{}) (errMsg string) {
	if e.Reason_Code == 0x00 {
		errMsg = "网关重启"
	} else if e.Reason_Code == 0x01 {
		errMsg = "其他原因"
	} else {
		errMsg = "其他原因"
	}
	return errMsg
}
