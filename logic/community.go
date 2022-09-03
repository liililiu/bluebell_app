package logic

import (
	"bluebell_app/dao/mysql"
	"bluebell_app/models"
)

func GetCommunityList() ([]*models.Community, error) {
	return mysql.GetCommunityList()
}
