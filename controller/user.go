package controller

import (
	"bluebell_app/logic"
	"bluebell_app/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// 1 参数校验
// 2 业务处理
// 3 返回相应

func SignUpHandler(c *gin.Context) {
	// 获取参数校验
	p := new(models.ParamSignUp)

	// ShouldBindJSON 只能校验类型（json）、以及字段类型（字符串类型），其他的校验做不到
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误
		zap.L().Error("SignUp with invalid param", zap.Error(err)) //打印日志
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求参数有误",
		})
		return
	}
	// 手动对参数进行业务规则校验（不要相信前端）
	if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
		zap.L().Error("SignUp with invalid param : err") //打印日志
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求参数有误",
		})
		return
	}
	fmt.Println(p)
	// 业务处理
	logic.Signup(p)
	// 返回数据
	c.JSON(http.StatusOK, gin.H{

		"msg": "success",
	})
}