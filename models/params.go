package models

//定义请求参数结构体

// 注册请求参数

type ParamSignUp struct {
	Username string `json:"username" binding:"required,gte=2,lte=20"` //gte大于等于、lte小于等于
	Password string `json:"password" binding:"required"`
	//Email      string `json:"email" binding:"required,email"`                           //email可以做电子邮件的格式验证
	RePassword string `json:"re_password" binding:"required,eqfield=Password"` // eqfield一个字段必须等于另一个字段
}

// 登录请求参数

type ParamLogin struct {
	Username string `json:"username" binding:"required,gte=2,lte=20"` //gte大于等于、lte小于等于
	Password string `json:"password" binding:"required"`
}

// 创建帖子请求参数

type Post struct {
	ID          int64  `json:"ID" `
	AuthorID    int64  `json:"author_id" `
	CommunityID string `json:"community_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Context     string `json:"context" binding:"required"`
}

// 投票数据

type ParamVoteData struct {
	// UserID  //用户id
	PostID    string `json:"post_id" binding:"required" `                       //帖子id
	Direction int8   `json:"direction,string" binding:"oneof=1 -1 0" ` // 赞成票1还是反对票-1取消投票0
}
