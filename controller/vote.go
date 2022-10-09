package controller

import (
	"bluebell_app/logic"
	"bluebell_app/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func PostVoteController(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("err:", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}

	//获取当前用户id
	uid, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 处理请求
	if err := logic.VoteForPost(uid, p); err != nil {
		zap.L().Error("logic.VoteForPost failed,", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}
