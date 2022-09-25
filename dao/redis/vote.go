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

func CreatePostTime(postID int64) error {
	pipeline:=Rdb.TxPipeline()
	pipeline.ZAdd(GetRedisKey(KeyPostTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 帖子初始分数值
	pipeline.ZAdd(GetRedisKey(KeyPostScore), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	_,err:=pipeline.Exec()
	return err
}

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
	ov := Rdb.ZScore(GetRedisKey(KeyPostVotedPrefix+postID), userID).Val()
	var dir float64
	if value > ov {
		dir = 1
	} else {
		dir = -1
	}
	diff := math.Abs(ov - value) //计算两次投票的插值

	pipeline:=Rdb.TxPipeline()
	pipeline.ZIncrBy(GetRedisKey(KeyPostScore), dir*diff*scorePerVote, postID)

	//3 记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(GetRedisKey(KeyPostVotedPrefix+postID), postID)
	} else {
		pipeline.ZAdd(GetRedisKey(KeyPostVotedPrefix+postID), redis.Z{
			Score:  value, //赞成票还是反对票
			Member: userID,
		})
	}
	_,err:=pipeline.Exec()
	return err
}
