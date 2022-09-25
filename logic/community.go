package logic

import (
	"bluebell_app/dao/mysql"
	"bluebell_app/dao/redis"
	"bluebell_app/models"
	sf "bluebell_app/pkg/snowflake"
	"go.uber.org/zap"
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
	p.ID = sf.GenID()
	// 调用dao层来实现
	err := mysql.CreatePost(p)
	if err != nil {
		return err
	}
	redis.CreatePostTime(p.ID)
	return err
}

// 获取帖子详情

func GetPostDetail(id int64) (p *models.ApiPostDb, err error) {

	// 调用dao层实现，
	// 多次查询，组合返回前端所需要的数据
	// ==查询帖子细节
	data1, err := mysql.GetPostDetail(id)
	if err != nil {
		return nil, err
	}
	// ==根据user_id，查询user
	user, _ := mysql.GetUserByID(data1.AuthorID)
	// ==根据community_id,查询社区名
	communityAll, _ := mysql.GetCommunityDeatilByID(data1.CommunityID)
	//根据需要返回的信息组合需要的结构体
	p = &models.ApiPostDb{
		CommunityName: communityAll.Name,
		PostDB:        data1,
		User:          user,
	}
	return

}

func PostList(page, size int64) (data []*models.ApiPostDb, err error) {
	// 调用dao层实现，

	data1, err := mysql.PostList(page, size)
	if err != nil {
		zap.L().Error("logic.mysql.PostList failed ", zap.Error(err))
		return nil, err
	}

	t := make([]*models.ApiPostDb, 0, len(data1))
	for k, v := range data1 {
		user, _ := mysql.GetUserByID(v.AuthorID)
		//if err != nil {
		//	zap.L().Error("logic.mysql.GetUserByID failed ", zap.Error(err))
		//	return nil, err
		//}
		communityAll, _ := mysql.GetCommunityDeatilByID(v.CommunityID)
		p := &models.ApiPostDb{
			CommunityName: communityAll.Name,
			PostDB:        data1[k],
			User:          user,
		}
		t = append(t, p)
	}

	return t, nil

}
