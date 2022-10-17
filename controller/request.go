package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

const ContextUserIDKey = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

// getCurrentUser 获取当前用户ID
func getCurrentUser(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(ContextUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return userID, err
}

func getPageInfo(c *gin.Context) (int64, int64, error) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		zap.L().Error("PostList.pageStr failed ", zap.Error(err))
		Response400(c, CodeInvalidParam)
		return 0, 0, err
	}
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		zap.L().Error("PostList.pageStr failed ", zap.Error(err))
		Response400(c, CodeInvalidParam)
		return 0, 0, err

	}
	return page, size, nil
}
