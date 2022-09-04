package models

import "time"

// 内存对齐，同类型的字段放邻近位置

type PostDB struct {
	Id          int64     `json:"id" db:"post_id"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id"`
	Title       string    `json:"title" db:"title"`
	Context     string    `json:"context" db:"context"`
	Status      int32     `json:"status" db:"status"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}
