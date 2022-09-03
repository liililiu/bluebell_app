package mysql

import (
	"bluebell_app/models"
	"database/sql"
	"go.uber.org/zap"
)

func GetCommunityList() (data []*models.Community, err error) {
	var community []*models.Community
	sqlStr := `select community_id ,community_name from community`
	err = db.Select(&community, sqlStr)
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db! ")
			return
		} else {
			zap.L().Error("mysql.GetCommunityList failed,", zap.Error(err))
			return nil, ErrorSql
		}
	}
	return community, nil
}

func GetCommunityDeatilByID(id int64) (data *models.CommunityDetail, err error) {
	details := new(models.CommunityDetail)
	sqlStr := `select community_id,community_name,introduction,create_time from community where community_id =?`
	err = db.Get(details, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorNoRow
		}
		return nil, err
	}
	return details, nil
}
