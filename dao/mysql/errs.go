package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("密码错误")
	ErrorSql             = errors.New("数据库查询异常")
	ErrorNoRow           = errors.New("未查到相关记录")
)
