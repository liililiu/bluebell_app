package mysql

import (
	"bluebell_app/settings"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init(cfg *settings.MysqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/"+"%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
	)
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("sqlx.Connect failed...", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(cfg.MaxOpenConnects)
	db.SetMaxIdleConns(cfg.MaxIdleConnects)
	return
}

// Close main函数中初始化连接后，defer调用
func Close() {
	_ = db.Close()
}
