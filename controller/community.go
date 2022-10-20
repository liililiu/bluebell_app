package controller

// 社区相关
import (
	"bluebell_app/dao/mysql"
	"bluebell_app/logic"
	"bluebell_app/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

const (
	orderTime  = "time"
	orderScore = "score"
)

// CommunityHandler 获取社区名称列表
// @Summary 获取社区名称
// @Description 获取社区名称
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 令牌"
// @Security ApiKeyAuth
// @Success 200 {object} []models.Community
// @Failure 401  {object}  ResponseData
// @Failure 500  {object}  ResponseData
// @Router /community [get]
func CommunityHandler(c *gin.Context) {
	//获取所有的社区相关的id以及name
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed, ", zap.Error(err))
		Response500(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)

}

// CommunityDetailHandler 获取社区详情
// @Summary 获取社区详情
// @Description 获取指定id社区详情内容
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 令牌"
// @Param id path string true "社区id"
// @Security ApiKeyAuth
// @Success 200 {object}   models.CommunityDetail
// @Failure 401  {object}  ResponseData
// @Failure 500  {object}  ResponseData
// @Router /community/{id} [get]
func CommunityDetailHandler(c *gin.Context) {
	//根据id来实现；解析参数，拿到id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("用户请求id参数错误：strconv.ParseInt failed,", zap.Error(err))
		Response500(c, CodeServerBusy)

		return
	}

	//调用logic层进行业务处理
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		if err == mysql.ErrorNoRow {
			zap.L().Warn("用户所查找的社区id不存在", zap.Int64("id值：", id))
			Response400(c, CodeInvalidRow)
			return
		}
		zap.L().Error("logic.GetCommunityDetail failed,", zap.Error(err))
		Response500(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)

}

// CreatePost 创建帖子
// @Summary 创建帖子
// @Description 创建帖子
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 令牌"
// @Param models.Post body models.Post true "创建帖子参数"
// @Security ApiKeyAuth
// @Success 200 {object}   nil
// @Failure 401  {object}  ResponseData
// @Failure 500  {object}  ResponseData
// @Router /post [post]
func CreatePost(c *gin.Context) {
	//获取参数及校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreatePost.ShouldBindJSON failed ", zap.Error(err))
		Response400(c, CodeInvalidParam)
		return
	}
	//业务处理，调用logic层
	// == 需要获取用户uuid
	id, err := getCurrentUser(c)
	if err != nil {
		return
	}
	p.AuthorID = id

	// == 传参 进行业务处理
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost", zap.Error(err))
		Response500(c, CodeServerBusy)
		return
	}
	//返回请求响应
	ResponseSuccess(c, nil)
}

// GetPostDetail
// @Summary 帖子详情
// @Description 获取帖子详情接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 令牌"
// @Param id path string true "帖子ID"
// @Success 200 {object}   ResponseData
// @Failure 400  {object}  ResponseData
// @Failure 401  {object}  ResponseData
// @Failure 500  {object}  ResponseData
// @Router /post/{id} [get]
func GetPostDetail(c *gin.Context) {
	//参数处理，从url中获取
	communityID := c.Param("id")
	id, err := strconv.ParseInt(communityID, 10, 64)
	if err != nil {
		zap.L().Error("GetPostDetail.communityID.ParseInt failed ", zap.Error(err))
		Response400(c, CodeInvalidParam)
		return
	}
	//业务处理
	data, err := logic.GetPostDetail(id)
	if err != nil {
		if err == mysql.ErrorNoRow {
			Response400(c, CodeInvalidRow)
			return
		} else {
			zap.L().Error(" logic.GetPostDetail failed ", zap.Error(err))
			Response500(c, CodeServerBusy)
			return
		}

	}
	//返回响应
	ResponseSuccess(c, data)
}

// PostList 可选时间顺序或者投票顺序返回；从redis获取排序信息
// 根据前端传来的参数动态的获取帖子列表(创建时间、或分数)
//1.获取参数
//2.去redis查询id列表
//3.根据id去数据库查询帖子详细信息

// PostList
// @Summary 帖子列表
// @Description 根据时间新旧或票赞数高低返回帖子列表
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 令牌"
// @Param page query int true "页数"
// @Param size query int true "页尺寸"
// @Param order query string false "页尺寸"
// @Success 200  {object}   ResponseData
// @Failure 400  {object}  ResponseData
// @Failure 401  {object}  ResponseData
// @Failure 500  {object}  ResponseData
// @Router /post [get]
func PostList(c *gin.Context) {
	// 处理请求参数
	// 默认值
	p := models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: orderTime,
	}

	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error("PostList2.ShouldBindQuery failed;", zap.Error(err))
		Response400(c, CodeInvalidParam)
		return
	}

	//业务处理
	data, err := logic.PostList2(&p)
	if err != nil {
		if err == mysql.ErrorNoRow {
			Response400(c, CodeInvalidRow)
			return
		} else {
			zap.L().Error(" logic.PostList failed ", zap.Error(err))
			Response500(c, CodeServerBusy)
			return
		}

	}
	//返回响应
	ResponseSuccess(c, data)
}
