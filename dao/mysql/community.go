package mysql

import (
	"bluebell_app/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

// GetCommunityList 返回社区名称列表
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

// GetCommunityDeatilByID 返回社区所有信息详情
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

// CreatePost 创建帖子
func CreatePost(p *models.Post) error {
	communityID, _ := strconv.ParseInt(p.CommunityID, 10, 64)
	dbP := models.PostDB{
		Id:          p.ID,
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

// GetPostDetail 获取帖子详情
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

// PostList 获取帖子详情列表分页
func PostList(page, size int64) (data []*models.PostDB, err error) {
	sqlStr := `select title,context,author_id,community_id,status from post order by create_time desc limit ? offset ?`
	p := make([]*models.PostDB, 0, size)
	// 注意分页
	if err := db.Select(&p, sqlStr, size, (page-1)*size); err != nil {
		zap.L().Error("mysql.PostList.Select failed ", zap.Error(err))
		return nil, err
	}
	return p, nil

}

//GetPostListByIDs 根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.PostDB, err error) {
	sqlStr := `select post_id,title,context,author_id,community_id,create_time from post 
    		where post_id in (?)
    		order by FIND_IN_SET(post_id,?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		zap.L().Error("sqlx.In(sqlStr, ids) failed : ", zap.Error(err))
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...) //注意写法
	return

}
