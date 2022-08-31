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

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// 返回错误信息；返回状态码

func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})

}

// 返回正确信息；返回数据

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}

// 返回自定义信息；返回msg

func ResponseSuccessWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
