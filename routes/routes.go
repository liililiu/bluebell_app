package routes

import (
	"bluebell_app/controller"
	"bluebell_app/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务路由

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong!")

	})

	r.POST("/signup", controller.SignUpHandler)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": 404,
		})
	})

	return r
}
