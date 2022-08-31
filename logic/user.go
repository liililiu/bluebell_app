package logic

import (
	"bluebell_app/dao/mysql"
	"bluebell_app/models"
	sf "bluebell_app/pkg/snowflake"
)

// 存放业务逻辑的代码

// 登录业务处理

func Signup(p *models.ParamSignUp) (err error) {
	//判断用户是否存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	//生成uuid
	userID := sf.GenID()
	// 构造一个user实例
	user := &models.User{
		UserID:   userID,
		UserName: p.Username,
		PassWord: p.Password,
	}
	//保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (err error) {
	user := &models.User{
		UserName: p.Username,
		PassWord: p.Password,
	}
	return mysql.Login(user)
}
