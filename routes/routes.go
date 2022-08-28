package routes

import (
	"bluebell_app/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hi,baby!")

	})
	return r
}
