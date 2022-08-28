package main

import (
	"bluebell_app/dao/mysql"
	"bluebell_app/dao/redis"
	"bluebell_app/logger"
	"bluebell_app/settings"
	"fmt"

	"go.uber.org/zap"
)

//GoWeb开发通用脚手架

func main() {
	//1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed , err: %v", err)
		return
	}
	//2.初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("init logger failed , err: %v", err)
		return
	}
	zap.L().Info("logger init success...")

	//3.初始化mysql
	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failed , err: %v", err)
		return
	}
	//4.初始化redis
	if err := redis.Init(); err != nil {
		fmt.Printf("init redis failed , err: %v", err)
		return
	}
	//5.注册路由
	//6.启停服务（优雅关机）
}
