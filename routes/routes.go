package routes

import (
	"bluebell_app/controller"
	"bluebell_app/logger"
	"bluebell_app/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")

	// 注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)
	v1.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong!")

	})
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		v1.POST("/post", controller.CreatePost)
		v1.GET("/post/:id", controller.GetPostDetail)
		//单纯的按照帖子从新到旧返回
		v1.GET("/post", controller.PostList)
		//可以按照时间顺序或者投票顺序返回
		v1.GET("/post2", controller.PostList2)
		//投票
		v1.POST("/vote", controller.PostVoteController)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": 404,
		})
	})

	return r
}
