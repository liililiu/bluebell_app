package redis

// redis key
// redis key 注意使用命名空间的方式,方便查询和拆分

const (
	KeyPrefix          = "bluebell:"
	KeyPostTime        = "post:time"   // 帖子发帖时间 Zset
	KeyPostScore       = "post:score"  // 帖子投票分数 Zset
	KeyPostVotedPrefix = "post:voted:" // 记录用户及投票类型 Zset
)