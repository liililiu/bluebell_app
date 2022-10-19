package controller

import (
	"bluebell_app/logic"
	"bluebell_app/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// PostVoteController
// @Summary 投票接口
// @Description 投票接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 令牌"
// @Param models.ParamVoteData body models.ParamVoteData true "投票参数"
// @Success 200 {object}   ResponseData
// @Failure 400  {object}  ResponseData
// @Failure 401  {object}  ResponseData
// @Failure 500  {object}  ResponseData
// @Router /vote [post]
func PostVoteController(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("err:", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			Response400(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}

	//获取当前用户id
	uid, err := getCurrentUser(c)
	if err != nil {
		Response401(c, CodeNeedLogin)
		return
	}
	// 处理请求
	// 两个参数,哪个用户给哪个帖子投了什么票
	if err := logic.VoteForPost(uid, p); err != nil {
		zap.L().Error("logic.VoteForPost failed,", zap.Error(err))
		Response500(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}
