package logic

import (
	"bluebell_app/dao/mysql"
	"bluebell_app/models"
	"bluebell_app/pkg/jwt"
	sf "bluebell_app/pkg/snowflake"
	"go.uber.org/zap"
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

func Login(p *models.ParamLogin) (u *models.User, err error) {
	user := &models.User{
		UserName: p.Username,
		PassWord: p.Password,
	}
	// 因jwt需要，获取user.userid
	uid, err := mysql.GetUserID(user.UserName)
	if err != nil {
		zap.L().Error("mysql.GetUserID failed ", zap.Error(err))
		return
	}
	user.UserID = uid
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	//登录成功，签发token

	token, err := jwt.GenToken(user.UserID, user.UserName)
	if err != nil {
		return
	}
	user.Token = token
	return
}
