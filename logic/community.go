package logic

import (
	"bluebell_app/dao/mysql"
	"bluebell_app/models"
)

func GetCommunityList() ([]*models.Community, error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	//调用dao层来实现
	return mysql.GetCommunityDeatilByID(id)
}
