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
