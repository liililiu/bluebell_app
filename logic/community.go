package logic

import (
	"bluebell_app/dao/mysql"
	"bluebell_app/models"
)

func GetCommunityList() ([]*models.Community, error) {
	return mysql.GetCommunityList()
}

// 获取社区分类详情

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	//调用dao层来实现
	return mysql.GetCommunityDeatilByID(id)
}

func CreatePost(p *models.Post) error {
	// 调用dao层来实现
	return mysql.CreatePost(p)
}

// 获取帖子详情

func GetPostDetail(id int64) (*models.PostDB, error) {
	// 调用dao层实现
	return mysql.GetPostDetail(id)
}
