package logic

import (
	"bluebell_app/dao/redis"
	"bluebell_app/models"
	"go.uber.org/zap"
	"strconv"
)

// 投票功能
// 谁给哪个帖子投了什么票(三要素)

// 用户投票是有相关算法的，推荐阅读阮一峰博客相关内容

// 这里使用简化版的算法

//=====投票参数过来的几种情况
//direction=1
//	1.之前没有，现在赞成  -->更新分数和投票记录
//	2.之前反对，现在赞成  -->更新分数和投票记录
//direction=0
//	1.之前赞成，现在取消  -->更新分数和投票记录
//	2.之前反对，现在取消  -->更新分数和投票记录
//direction=-1
//	1.之前没有，现在反对  -->更新分数和投票记录
//	2.之前赞成，现在反对  -->更新分数和投票记录

//=====投票的限制(冷热数据处理)
//帖子自发表之日一周后就不允许再投票了
//	1.到期之后将redis保存的数据持久化到mysql中
//	2.到期之后删除 KeyPostVotedPrefix

// VoteForPost 为帖子投票的函数
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))

}
