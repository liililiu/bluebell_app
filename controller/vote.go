package controller

import (
	"bluebell_app/logic"
	"bluebell_app/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func PostVoteController(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	// 处理请求
	logic.PostVote()
	// 返回响应
	ResponseSuccess(c, nil)
}
