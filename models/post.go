package models

import "time"

// 内存对齐，同类型的字段放邻近位置

//PostDB 帖子信息
type PostDB struct {
	Id          int64     `json:"id" db:"post_id"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id"`
	Title       string    `json:"title" db:"title"`
	Context     string    `json:"context" db:"context"`
	Status      int32     `json:"status" db:"status"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// ApiPostDb 返回前端更详细的信息
type ApiPostDb struct {
	CommunityName string `json:"community_name"`
	*PostDB
	*User `json:"user"` // 给前端返回的数据分层
}
