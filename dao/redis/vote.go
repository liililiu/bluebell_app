package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"time"
)

var (
	ErrVoteTimeExpir = errors.New("投票时间已过")
)

// CreatePostTime 帖子创建时 记录创建时间
func CreatePostTime(postID int64) error {
	pipeline := Rdb.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(GetRedisKey(KeyPostTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 帖子初始分数值
	pipeline.ZAdd(GetRedisKey(KeyPostScore), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	_, err := pipeline.Exec()
	return err
}

// VoteForPost value就是投的1 -1 0三票
func VoteForPost(userID, postID string, value float64) error {
	// 1判断投票限制，时间是否超过一周
	//去redis获取帖子发表时间
	postTime := Rdb.ZScore(GetRedisKey(KeyPostTime), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpir
	}

	// 2和3放到事务中
	// 2更新帖子分数
	//先查当前用户对当前帖子之前的投票纪录
	// ov代表之前投的票的数值 1 -1 0
	ov := Rdb.ZScore(GetRedisKey(KeyPostVotedPrefix+postID), userID).Val()

	// 因为值只能为1 -1 0;这里的 例如 (-1)-(-1)=0 就决定了不能无限制投票
	var dir float64
	if value > ov {
		dir = 1
	} else {
		dir = -1
	}
	diff := math.Abs(ov - value) //计算两次投票的插值

	pipeline := Rdb.TxPipeline()
	// 更新分数,注意dir的正负取值
	pipeline.ZIncrBy(GetRedisKey(KeyPostScore), dir*diff*scorePerVote, postID)

	//3 记录用户为该帖子投票的数据
	if value == 0 {
		// 移除投票
		pipeline.ZRem(GetRedisKey(KeyPostVotedPrefix+postID), postID)
	} else {
		pipeline.ZAdd(GetRedisKey(KeyPostVotedPrefix+postID), redis.Z{
			Score:  value, //赞成票还是反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
