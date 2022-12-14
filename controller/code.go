package controller

// 定义code以及编写获取code对应msg的方法

type ResCode int64

const (
	CodeSuccess ResCode = 10000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeNeedLogin
	CodeInvalidToken

	CodeInvalidRow
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名密码错误",
	CodeServerBusy:      "服务繁忙",
	CodeNeedLogin:       "需要登录",
	CodeInvalidToken:    "无效token",
	CodeInvalidRow:      "没有相关记录",
}

// Msg 获取code对应的msg的方法
func (code ResCode) Msg() string {
	msg, ok := codeMsgMap[code]
	if !ok {
		return codeMsgMap[CodeServerBusy]
	}
	return msg
}
