package controller

import (
	"bluebell_app/logic"
	"bluebell_app/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
)

// 1 参数校验
// 2 业务处理
// 3 返回相应

func SignUpHandler(c *gin.Context) {
	// 1.获取参数校验
	p := new(models.ParamSignUp)

	// 在web开发中一个不可避免的环节就是对请求参数进行校验，通常我们会在代码中定义与请求参数相对应的模型（结构体）
	// ShouldBindJSON 只能校验类型（json）、以及字段类型（字符串类型），其他的校验做不到
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误,绑定失败
		zap.L().Error("SignUp with invalid param", zap.Error(err)) //打印日志
		// 判断err是否是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				// 非validator.ValidationErrors类型错误直接返回
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			// 调用错误类型翻译器进行中文翻译
			// 去除返回前端信息中的，msg的结构体前缀；涉及反射，会影响效率
			"msg": removeTopStruct(errs.Translate(trans)),
		})
		return
	}
	// 手动对参数进行业务规则校验（不要相信前端）;被注释掉了，采用上面的validator去校验
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
	//	zap.L().Error("SignUp with invalid param : err") //打印日志
	//
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数有误",
	//	})
	//	return
	//}
	fmt.Println(p)
	// 2.业务处理
	if err := logic.Signup(p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// 3.返回数据
	c.JSON(http.StatusOK, gin.H{

		"msg": "success",
	})
}
