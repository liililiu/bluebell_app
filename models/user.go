package models

type User struct {
	UserID   int64  `db:"user_id"` // int64和uuid类型保持一致
	UserName string `db:"username"`
	PassWord string `db:"password"` // tag保持与数据库一致
}
