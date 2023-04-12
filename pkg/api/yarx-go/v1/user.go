package v1

// 登录请求表单
type LoginRequestForm struct {
	Username string `json:"username" valid:"alphanum,required,stringlength(1|255)"`
	Password string `json:"password" valid:"required,stringlength(6|18)"`
}

// 登录后响应的内容
type LoginResponse struct {
	Token string `json:"token"`
}

// 创建用户请求表单
type CreateUserRequestForm struct {
	Username string `json:"username" valid:"alphanum,required,stringlength(1|255)"`
	Password string `json:"password" valid:"required,stringlength(6|18)"`
	NickName string `json:"nicname" valid:"required,stringlength(1|255)"`
	Email    string `json:"email" valid:"required,email"`
}

// 用户详情
type UserInfo struct {
	Username  string `json:"username"`
	NickName  string `json:"nicname"`
	Email     string `json:email`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
