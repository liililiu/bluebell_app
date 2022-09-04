package mysql

import (
	"bluebell_app/models"
	sf "bluebell_app/pkg/snowflake"
	"database/sql"
	"go.uber.org/zap"
	"strconv"
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

func CreatePost(p *models.Post) error {
	//生成uuid
	id := sf.GenID()
	communityID, _ := strconv.ParseInt(p.CommunityID, 10, 64)
	dbP := models.PostDB{
		Id:          id,
		AuthorID:    p.AuthorID,
		CommunityID: communityID,
		Title:       p.Title,
		Context:     p.Context,
		Status:      1, // 1为已创建，0为未创建
	}
	//入库
	//goland:noinspection SqlInsertValues
	sqlStr := `insert into post (post_id,title,context,author_id,community_id,status) values (
		?,?,?,?,?,?)`
	_, err := db.Exec(sqlStr, dbP.Id, dbP.Title, dbP.Context, dbP.AuthorID, dbP.CommunityID, dbP.Status)
	//返回
	return err

}

func GetPostDetail(id int64) (data *models.PostDB, err error) {
	p := new(models.PostDB)
	sqlStr := `select title,context,author_id,community_id,status from post where post_id=?`
	if err := db.Get(p, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorNoRow
		}
		return nil, err
	}
	return p, nil
}
