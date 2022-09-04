package controller

// 社区相关

import (
	"bluebell_app/dao/mysql"
	"bluebell_app/logic"
	"bluebell_app/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// CommunityHandler获取社区名称列表

func CommunityHandler(c *gin.Context) {
	//获取所有的社区相关的id以及name
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed, ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)

}

// CommunityDetailHandler 获取社区详情

func CommunityDetailHandler(c *gin.Context) {
	//根据id来实现；解析参数，拿到id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("用户请求id参数错误：strconv.ParseInt failed,", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//调用logic层进行业务处理
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		if err == mysql.ErrorNoRow {
			zap.L().Warn("用户所查找的社区id不存在", zap.Int64("id值：", id))
			ResponseError(c, CodeInvalidRow)
			return
		}
		zap.L().Error("logic.GetCommunityDetail failed,", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)

}

// CreatePost 创建帖子
func CreatePost(c *gin.Context) {
	//获取参数及校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreatePost.ShouldBindJSON failed ", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//业务处理，调用logic层

	// == 需要获取用户uuid
	id, err := getCurrentUser(c)
	if err != nil {
		return
	}
	p.AuthorID = id
	fmt.Println("##########")
	fmt.Println(p.AuthorID)
	fmt.Println("##########")

	// == 传参 进行业务处理
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回请求响应
	ResponseSuccess(c, nil)
}

// GetPostDetail 获取帖子细节
func GetPostDetail(c *gin.Context) {
	//参数处理，从url中获取
	communityID := c.Param("id")
	id, err := strconv.ParseInt(communityID, 10, 64)
	if err != nil {

		zap.L().Error("GetPostDetail.communityID.ParseInt failed ", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//业务处理
	data, err := logic.GetPostDetail(id)
	if err != nil {
		if err == mysql.ErrorNoRow {
			ResponseError(c, CodeInvalidRow)
			return
		} else {
			zap.L().Error(" logic.GetPostDetail failed ", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}

	}
	//返回响应
	ResponseSuccess(c, data)
}
