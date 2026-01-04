package request

type SendEmailVerificationCode struct {
	Email     string `json:"email" binding:"required,email"` // 字段必须，且为 email 格式
	Captcha   string `json:"captcha" binding:"required,len=6"`
	CaptchaID string `json:"captcha_id" binding:"required"`
}
