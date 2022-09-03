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
			return nil, Errorsql
		}
	}
	return community, nil
}
