package request

// Register 用户注册 request
type Register struct {
	Username         string `json:"username" binding:"required,max=20"`
	Password         string `json:"password" binding:"required,min=6,max=16"`
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verification_code" binding:"required,len=6"`
}

// Login 用户登录 request
type Login struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6,max=16"`
	Captcha   string `json:"captcha" binding:"required,len=6"`
	CaptchaID string `json:"captcha_id" binding:"required"`
}

// ForgotPassword 用户找回密码请求
type ForgotPassword struct {
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verification_code" binding:"required,len=6"`
	NewPassword      string `json:"new_password" binding:"required,min=8,max=16"`
}

// UserCard 用户卡片请求结构体
type UserCard struct {
	UUID string `json:"uuid" binding:"required" form:"uuid"`
}

// UserResetPassword 用户重置密码请求
type UserResetPassword struct {
	UserID      uint   `json:"-"`
	Password    string `json:"password" binding:"required,min=8,max=16"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=16"`
}

// UserChangeInfo 用户修改信息请求
type UserChangeInfo struct {
	UserID    uint   `json:"-"`
	Username  string `json:"username" binding:"required,max=20"`
	Address   string `json:"address" binding:"required,max=200"`
	Signature string `json:"signature" binding:"max=320"`
}

// UserChart 用户图表数据请求
type UserChart struct {
	Date int `json:"date" form:"date" binding:"required,oneof=7 30 90 180 365"`
}

// UserList 用户列表请求
type UserList struct {
	UUID     *string `json:"uuid" form:"uuid"` // *string 代表参数可传可不传，前端可以传 null
	Register *string `json:"register" form:"register"`
	PageInfo
}

// UserOperation 冻结/解冻用户请求
type UserOperation struct {
	ID uint `json:"id" binding:"required"`
}

// UserLoginList 获取登录日志列表请求
type UserLoginList struct {
	UUID *string `json:"uuid" form:"uuid"`
	PageInfo
}
