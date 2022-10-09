package redis

// redis key
// redis key 注意使用命名空间的方式,方便查询和拆分

const (
	KeyPrefix          = "bluebell:"
	KeyPostTime        = "post:time"   // 帖子发帖时间 Zset
	KeyPostScore       = "post:score"  // 帖子投票分数 Zset
	KeyPostVotedPrefix = "post:voted:" // 记录用户及投票类型 Zset,参数是帖子id

	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票的分数
)

func GetRedisKey(key string) string {
	return KeyPrefix + key
}
