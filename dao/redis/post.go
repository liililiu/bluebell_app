package redis

import "bluebell_app/models"

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
