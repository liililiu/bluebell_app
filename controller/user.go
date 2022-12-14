package controller

import (
	"bluebell_app/dao/mysql"
	"bluebell_app/logic"
	"bluebell_app/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 1 参数校验
// 2 业务处理
// 3 返回相应

// SignUpHandler
// @Summary 注册接口
// @Description 用户注册
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param models.ParamSignUp body models.ParamSignUp true "注册参数"
// @Success 200 {object}   ResponseData
// @Failure 400  {object}  ResponseData
// @Failure 500  {object}  ResponseData
// @Router /signup [post]
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
			Response400(c, CodeInvalidParam)
		}
		ResponseSuccessWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2.业务处理
	if err := logic.Signup(p); err != nil {
		zap.L().Error("logic.Signup failed", zap.Error(err)) //打印日志
		if errors.Is(err, mysql.ErrorUserExist) {
			Response500(c, CodeUserExist)
			return
		}
		Response500(c, CodeServerBusy)
		return
	}
	// 3.返回数据
	ResponseSuccess(c, nil)
}

// LoginHandler
// @Summary 登录接口
// @Description 用户登录
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param models.ParamLogin body models.ParamLogin true "登录参数"
// @Success 200 {object}   ResponseData
// @Failure 400  {object}  ResponseData
// @Failure 500  {object}  ResponseData
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	//1.参数校验(这一块可直接复用)
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误,绑定失败
		zap.L().Error("Login with invalid param", zap.Error(err)) //打印日志
		// 判断err是否是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			Response400(c, CodeInvalidParam)
			return
		}
		ResponseSuccessWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	//2.业务处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err)) //打印日志
		if err == mysql.ErrorUserExist {
			Response500(c, CodeUserExist)
			return
		}
		Response400(c, CodeInvalidPassword)
		return
	}
	//3.返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   user.UserID, //id值为int64可能大于前端number类型所能展示的最大值
		"user_name": user.UserName,
		"token":     user.Token,
	})
}
