package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
{
	"code": 10001, //程序中的错误码
	"msg":xx,    //提示信息
	"data":{},  //数据
}
*/

// ResponseData 响应结构体
type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	//omitempty  当字段为空时就不返回
	Data interface{} `json:"data,omitempty"`
	//Data interface{} `json:"data"`
}

// ResponseSuccess 返回正确信息；返回数据
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}

func Response500(c *gin.Context, code ResCode) {
	c.JSON(http.StatusInternalServerError, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

func Response400(c *gin.Context, code ResCode) {
	c.JSON(http.StatusBadRequest, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

func Response401(c *gin.Context, code ResCode) {
	c.JSON(http.StatusUnauthorized, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

func Response403(c *gin.Context, code ResCode) {
	c.JSON(http.StatusForbidden, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

// ResponseSuccessWithMsg 上述返回消息不能满足时，返回自定义信息
func ResponseSuccessWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// ResponseErrorWithMsg 上述返回消息不能满足时，返回自定义信息
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusInternalServerError, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
