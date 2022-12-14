package main

import (
	"bluebell_app/controller"
	"bluebell_app/dao/mysql"
	"bluebell_app/dao/redis"
	"bluebell_app/logger"
	sf "bluebell_app/pkg/snowflake"
	"bluebell_app/routes"
	"bluebell_app/settings"
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title GoWeb项目
// @version 1.0
// @description Gin+Mysql+Redis快速入门GoWeb开发
// @contact.url https://llbowen.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	//1.加载配置
	if err := settings.Init(); err != nil {
		zap.L().Error("settings.Init failed , err: %v", zap.Error(err))
		return
	}
	//2.初始化日志
	if err := logger.Init(settings.GlobalConfig.LogConfig, settings.GlobalConfig.Mode); err != nil {
		zap.L().Error("logger.Init failed , err: %v", zap.Error(err))
		return
	}
	zap.L().Info("logger.Init success...")

	defer zap.L().Sync()

	//3.初始化mysql
	if err := mysql.Init(settings.GlobalConfig.MysqlConfig); err != nil {
		zap.L().Error("mysql.Init failed , err: %v", zap.Error(err))
		return
	}
	zap.L().Info("mysql.Init success...")
	defer mysql.Close()
	//4.初始化redis
	if err := redis.Init(settings.GlobalConfig.RedisConfig); err != nil {
		zap.L().Error("redis.Init failed , err: %v", zap.Error(err))
		return
	}
	zap.L().Info("redis.Init success...")
	defer redis.Close()

	// 初始化gin框架内置的validator的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		zap.L().Error("validator.InitTrans init failed , err: %v", zap.Error(err))
		return
	}
	zap.L().Info("validator.InitTrans init success...")

	//uuid 初始化
	if err := sf.Init(settings.GlobalConfig.SnowflakeConfig.StartTime, settings.GlobalConfig.SnowflakeConfig.MachineID); err != nil {
		zap.L().Error("sf.Init failed , err: %v", zap.Error(err))
		return
	}

	zap.L().Info("sf.uuid.Init init success...")
	zap.L().Info("all init success, start to SetupRouter...")

	//5.注册路由
	r := routes.SetupRouter()

	//6.启停服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.GlobalConfig.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// ，为关闭服务器操作设置一个5等待中断信号来优雅地关闭服务器秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
