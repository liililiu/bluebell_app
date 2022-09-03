package mysql

import (
	"bluebell_app/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

// 把每一步数据库操作封装成函数，
// 等待logic层根据业务需要调用

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("密码错误")
	Errorsql             = errors.New("数据库查询异常")
)

func CheckUserExist(uname string) error {
	sqlStr := `select count(user_id) from user where username=?`
	var count int
	if err := db.Get(&count, sqlStr, uname); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

func InsertUser(u *models.User) (err error) {
	//对密码进行加密
	newPassword := encryptPassword(u.PassWord)
	//执行sql入库语句
	sqlStr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err = db.Exec(sqlStr, u.UserID, u.UserName, newPassword)
	return err
}

func encryptPassword(oPassword string) string {
	secret := "This is secret"
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(p *models.User) (err error) {
	var user models.User
	//1.查询用户是否存在
	sqlStr := `select user_id,username,password from user where username =?`
	err = db.Get(&user, sqlStr, p.UserName)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err // 查询数据库异常
	}
	//2.比对密码
	oPassword := p.PassWord
	if encryptPassword(oPassword) != user.PassWord {
		return ErrorInvalidPassword
	}
	return nil
}
