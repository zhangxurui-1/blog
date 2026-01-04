package api

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/utils"
	"time"
)

type UserApi struct {
}

// Register 用户注册
func (userApi *UserApi) Register(c *gin.Context) {
	var req request.Register
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 从会话中获取数据
	// 之所以可以，是因为之前在发送邮箱验证码时使用 session 存储了数据，sessions 包会在响应头加上 set-cookie 和 session-id
	session := sessions.Default(c)
	// 验证邮箱
	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("The email does not match the email to be verified", c)
		return
	}
	// 验证码
	savedCode := session.Get("verification_code")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		response.FailWithMessage("Invalid verification code", c)
		return
	}
	// 验证码是否过期
	savedTime := session.Get("expiry_time")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("The verification code has expired, please resent", c)
		return
	}

	// 创建新用户
	u := database.User{Username: req.Username, Password: req.Password, Email: req.Email}
	user, err := userService.Register(u)
	if err != nil {
		global.Log.Error("Failed to register user:", zap.Error(err))
		response.FailWithMessage("Failed to register user", c)
		return
	}

	userApi.TokenNext(c, user)
}

// Login 用户登录
func (userApi *UserApi) Login(c *gin.Context) {
	switch c.Query("flag") {
	case "email":
		userApi.EmailLogin(c)
	case "qq":
		userApi.QQLogin(c)
	default:
		userApi.EmailLogin(c)
	}
}

// EmailLogin 用户邮箱登录
func (userApi *UserApi) EmailLogin(c *gin.Context) {
	var req request.Login
	if err := c.ShouldBind(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if store.Verify(req.CaptchaID, req.Captcha, true) {
		u := database.User{Email: req.Email, Password: req.Password}
		user, err := userService.EmailLogin(u)
		if err != nil {
			global.Log.Error("Failed to Login:", zap.Error(err))
			response.FailWithMessage("Failed to Login", c)
			return
		}

		userApi.TokenNext(c, user)
		return
	}
	response.FailWithMessage("Incorrect verification code", c)
}

// QQLogin 用户 QQ 登录
func (userApi *UserApi) QQLogin(c *gin.Context) {
	code := c.Query("code")
	// 未获取到 Authorization code
	if code == "" {
		response.FailWithMessage("code is required", c)
		return
	}

	// 获取 access token
	accessTokenResponse, err := qqService.GetAccessTokenByCode(code)
	if err != nil || accessTokenResponse.OpenId == "" {
		global.Log.Error("Invalid code", zap.Error(err))
		response.FailWithMessage("Invalid code", c)
		return
	}

	// 根据 token 进行登录
	user, err := userService.QQLogin(accessTokenResponse)
	if err != nil {
		global.Log.Error("Failed to Login:", zap.Error(err))
		response.FailWithMessage("Failed to Login", c)
		return
	}

	// 登录成功后生成 token
	userApi.TokenNext(c, user)
}

// TokenNext 为用户发放 token
func (userApi *UserApi) TokenNext(c *gin.Context, user database.User) {
	// 检查用户是否被冻结
	if user.Freeze {
		response.FailWithMessage("This user is frozen, please contact the administrator", c)
		return
	}

	baseClaims := request.BaseClaims{
		UserID: user.ID,
		UUID:   user.UUID,
		RoleID: user.RoleID,
	}
	j := utils.NewJWT()

	// 创建 access token
	accessTokenClaims := j.CreateAccessClaims(baseClaims)
	accessToken, err := j.CreateAccessToken(accessTokenClaims)
	if err != nil {
		global.Log.Error("Failed to create access token", zap.Error(err))
		response.FailWithMessage("Failed to create access token", c)
		return
	}

	// 创建刷新令牌
	refreshTokenClaims := j.CreateRefreshClaims(baseClaims)
	refreshToken, err := j.CreateRefreshToken(refreshTokenClaims)
	if err != nil {
		global.Log.Error("Failed to create refresh token", zap.Error(err))
		response.FailWithMessage("Failed to create refresh token", c)
		return
	}
	// 是否开启多地点登录拦截
	if !global.Config.System.UseMultipoint {
		// 设置 refresh token 并返回（在响应中写入 set-cookie 和 refresh token 信息）
		utils.SetRefreshToken(c, refreshToken, int(refreshTokenClaims.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetail(
			response.Login{
				User:                 user,
				AccessToken:          accessToken,
				AccessTokenExpiresAt: accessTokenClaims.ExpiresAt.Unix() * 1000,
			}, "Login successfully", c)
		return
	}

	// 如果开启了多点登录拦截，则检查 Redis 中是否已存在该用户的 jwt
	if jwtStr, err := jwtService.GetRedisJWT(user.UUID); errors.Is(err, redis.Nil) {
		// 如果 redis 中没有，则添加
		if err2 := jwtService.SetRedisJWT(refreshToken, user.UUID); err2 != nil {
			global.Log.Error("Failed to set jwt to redis", zap.Error(err2))
			response.FailWithMessage("Failed to set jwt to redis", c)
			return
		}
		// 设置 refresh token 并返回
		utils.SetRefreshToken(c, refreshToken, int(refreshTokenClaims.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetail(
			response.Login{
				User:                 user,
				AccessToken:          accessToken,
				AccessTokenExpiresAt: accessTokenClaims.ExpiresAt.Unix() * 1000,
			}, "Login successfully", c)
	} else if err != nil {
		// 在读取 redis 时出现了错误
		global.Log.Error("Failed to get jwt from redis", zap.Error(err))
		response.FailWithMessage("Failed to get jwt from redis", c)
	} else {
		// redis 中已存在 jwt，将旧的 jwt 加入黑名单，并设置新的 token
		var blacklist database.JwtBlacklist
		blacklist.Jwt = jwtStr
		if err := jwtService.JoinInBlackList(blacklist); err != nil {
			global.Log.Error("Failed to invalidate jwt:", zap.Error(err))
			response.FailWithMessage("Failed to invalidate jwt:", c)
			return
		}
		// 设置新的 jwt 到 redis
		if err2 := jwtService.SetRedisJWT(refreshToken, user.UUID); err2 != nil {
			global.Log.Error("Failed to set jwt to redis", zap.Error(err2))
			response.FailWithMessage("Failed to set jwt to redis", c)
			return
		}
		// 设置 refresh token 并返回
		utils.SetRefreshToken(c, refreshToken, int(refreshTokenClaims.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetail(
			response.Login{
				User:                 user,
				AccessToken:          accessToken,
				AccessTokenExpiresAt: accessTokenClaims.ExpiresAt.Unix() * 1000,
			}, "Login successfully", c)
	}
}

// ForgotPassword 找回密码
func (userApi *UserApi) ForgotPassword(c *gin.Context) {
	var req request.ForgotPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 获取会话信息，参考 Register
	session := sessions.Default(c)
	savedEmail := session.Get("email")

	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("This email does not match the email to be verified", c)
		return
	}

	savedCode := session.Get("verification_code")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		response.FailWithMessage("Invalid verification code", c)
		return
	}

	savedTime := session.Get("expiry_time")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("This verification code has expired", c)
		return
	}

	err = userService.ForgotPassword(req)
	if err != nil {
		global.Log.Error("Failed to reset password", zap.Error(err))
		response.FailWithMessage("Failed to reset password", c)
		return
	}

	response.OkWithMessage("Successfully reset password", c)
}

// UserCard 获取用户卡片信息
func (userApi *UserApi) UserCard(c *gin.Context) {
	var req request.UserCard
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userCard, err := userService.UserCard(req)
	if err != nil {
		global.Log.Error("Failed to get user card", zap.Error(err))
		response.FailWithMessage("Failed to get user card", c)
		return
	}
	response.OkWithData(userCard, c)
}

// Logout 登出
func (userApi *UserApi) Logout(c *gin.Context) {
	userService.Logout(c)
	response.OkWithMessage("Successfully logout", c)
}

// UserResetPassword 修改密码
func (userApi *UserApi) UserResetPassword(c *gin.Context) {
	var req request.UserResetPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// jwt 认证中间件已经把 claims 放到上下文了
	req.UserID = utils.GetUserID(c)

	if err := userService.UserResetPassword(req); err != nil {
		global.Log.Error("Failed to modify", zap.Error(err))
		response.FailWithMessage("Failed to modify, original password does not match the current account", c)
		return
	}

	response.OkWithMessage("Successfully reset password, please login again", c)
	userService.Logout(c)
}

// UserInfo 获取个人信息
func (userApi *UserApi) UserInfo(c *gin.Context) {
	userID := utils.GetUserID(c)
	user, err := userService.UserInfo(userID)
	if err != nil {
		global.Log.Error("Failed to get user information", zap.Error(err))
		response.FailWithMessage("Failed to get user information", c)
		return
	}

	response.OkWithData(user, c)
}

// UserChangeInfo 修改个人信息
func (userApi *UserApi) UserChangeInfo(c *gin.Context) {
	var req request.UserChangeInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	req.UserID = utils.GetUserID(c)
	if err := userService.UserChangeInfo(req); err != nil {
		global.Log.Error("Failed to change user information", zap.Error(err))
		response.FailWithMessage("Failed to change user information", c)
		return
	}

	response.OkWithMessage("Successfully changed user information", c)
}

// UserWeather 获取天气
func (userApi *UserApi) UserWeather(c *gin.Context) {
	ip := c.ClientIP()
	weather, err := userService.UserWeather(ip)
	if err != nil {
		global.Log.Error("Failed to get user weather", zap.Error(err))
		response.FailWithMessage("Failed to get user weather", c)
		return
	}
	response.OkWithData(weather, c)
}

// UserChart 获取用户图表数据，登录和注册人数
func (userApi *UserApi) UserChart(c *gin.Context) {
	var req request.UserChart
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data, err := userService.UserChart(req)
	if err != nil {
		global.Log.Error("Failed to get user chart", zap.Error(err))
		response.FailWithMessage("Failed to get user chart", c)
		return
	}

	response.OkWithData(data, c)
}

// UserList 获取用户列表
func (userApi *UserApi) UserList(c *gin.Context) {
	var pageInfo request.UserList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := userService.UserList(pageInfo)

	if err != nil {
		global.Log.Error("Failed to get user list", zap.Error(err))
		response.FailWithMessage("Failed to get user list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

// UserFreeze 冻结用户
func (userApi *UserApi) UserFreeze(c *gin.Context) {
	var req request.UserOperation
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := userService.UserFreeze(req); err != nil {
		global.Log.Error("Failed to freeze user", zap.Error(err))
		response.FailWithMessage("Failed to freeze user", c)
		return
	}

	response.OkWithMessage("Successfully freeze user", c)
}

// UserUnfreeze 解冻用户
func (userApi *UserApi) UserUnfreeze(c *gin.Context) {
	var req request.UserOperation
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := userService.UserUnfreeze(req); err != nil {
		global.Log.Error("Failed to unfreeze user", zap.Error(err))
		response.FailWithMessage("Failed to unfreeze user", c)
		return
	}
	response.OkWithMessage("Successfully unfreeze user", c)
}

// UserLoginList 获取登录日志列表
func (userApi *UserApi) UserLoginList(c *gin.Context) {
	var pageInfo request.UserLoginList
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := userService.UserLoginList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get user list", zap.Error(err))
		response.FailWithMessage("Failed to get user list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}
