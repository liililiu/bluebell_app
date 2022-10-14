package redis

import "bluebell_app/models"

// GetPostIDsInOrder 新版获取帖子详情
func GetPostIDsInOrder(p *models.ParamPostList) []string {
	//从redis获取帖子id
	key := GetRedisKey(KeyPostScore)
	if p.Order == "time" {
		key = GetRedisKey(KeyPostTime)
	}
	// redis ZRange使用
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	ids := Rdb.ZRevRange(key, start, end).Val()

	return ids
}

// GetVoteNum 从redis中获取帖子赞同票和反对票数据
func GetVoteNum(ids []string) (yesNum, noNum []int64) {

	// 赞同票数量
	yesNum = make([]int64, 0, len(ids))
	for _, v := range ids {
		key := GetRedisKey(KeyPostVotedPrefix + v)
		yes := Rdb.ZCount(key, "1", "1").Val()
		yesNum = append(yesNum, yes)
	}
	//反对票数量
	noNum = make([]int64, 0, len(ids))
	for _, v := range ids {
		key := GetRedisKey(KeyPostVotedPrefix + v)
		yes := Rdb.ZCount(key, "-1", "-1").Val()
		noNum = append(noNum, yes)
	}
	return
}
