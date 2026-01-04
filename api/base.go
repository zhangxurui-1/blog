package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"server/global"
	"server/model/request"
	"server/model/response"
)

type BaseApi struct {
}

var store = base64Captcha.DefaultMemStore

// "github.com/mojocn/base64Captcha" 图形验证码第三方包

// Captcha 生成数字验证码
func (baseApi *BaseApi) Captcha(c *gin.Context) {
	// 创建驱动
	driver := base64Captcha.NewDriverDigit(
		global.Config.Captcha.Height,
		global.Config.Captcha.Width,
		global.Config.Captcha.Length,
		global.Config.Captcha.MaxSkew,
		global.Config.Captcha.DotCount,
	)
	// 创建验证码对象
	captcha := base64Captcha.NewCaptcha(driver, store)

	// 生成验证码
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		global.Log.Error("Failed to generate captcha", zap.Error(err))
		response.FailWithMessage("Failed to generate captcha", c)
		return
	}
	// 返回
	response.OkWithData(response.Captcha{
		CaptchaID: id,
		PicPath:   b64s,
	}, c)
}

// SendEmailVerificationCode 发送邮箱验证码
func (baseApi *BaseApi) SendEmailVerificationCode(c *gin.Context) {
	var req request.SendEmailVerificationCode
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 图形验证码验证成功
	if store.Verify(req.CaptchaID, req.Captcha, true) {
		// 邮箱验证码发送失败
		if err := baseService.SendEmailVerificationCode(c, req.Email); err != nil {
			global.Log.Error("Failed to send email verification code", zap.Error(err))
			response.FailWithMessage("Failed to send email verification code", c)
			return
		}
		// 发送成功
		response.OkWithMessage("Successfully sent email verification code", c)
		return
	}
	// 图形验证码验证失败
	response.FailWithMessage("Incorrect verification code", c)
}

// QQLoginURL 返回 QQ 登录链接
func (baseApi *BaseApi) QQLoginURL(c *gin.Context) {
	url := global.Config.QQ.QQLoginURL()
	response.OkWithData(url, c)
}
