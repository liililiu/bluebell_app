package main

import (
	"bluebell_app/logger"
	"bluebell_app/settings"
	"fmt"
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
	//3.初始化mysql
	//4.初始化redis
	//5.注册路由
	//6.启停服务（优雅关机）
}
