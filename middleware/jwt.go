package middleware

import (
	"errors"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/service"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

var jwtService = service.ServiceGroupApp.JwtService

// JWTAuth 认证中间件，验证 jwt 是否合法
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求中的 access token 和 refresh token
		accessToken := utils.GetAccessToken(c)
		refreshToken := utils.GetRefreshToken(c)

		// 检查 Refresh Token 是否在黑名单中，如果是，则清除 Refresh Token 并返回未授权错误
		if jwtService.IsInBlackList(refreshToken) {
			utils.ClearRefreshToken(c)
			response.NoAuth("Account logged in from another location or token is invalid", c)
			c.Abort()
			return
		}

		j := utils.NewJWT()

		// 解析 access token
		claims, err := j.ParseAccessToken(accessToken)
		// 如果解析失败
		if err != nil {
			// 如果 access token 为空或过期
			if accessToken == "" || errors.Is(err, utils.TokenExpired) {
				// 尝试解析 refresh token
				refreshClaims, err := j.ParseRefreshToken(refreshToken)
				if err != nil {
					// 如果 Refresh Token 也无法解析，清除 Refresh Token 并返回未授权错误
					utils.ClearRefreshToken(c)
					response.NoAuth("Refresh token expired or invalid", c)
					c.Abort()
					return
				}
				// 如果 Refresh Token 有效，通过其 UserID 获取用户信息
				var user database.User
				if err := global.DB.Select("uuid", "role_id").Take(&user, refreshClaims.UserID).Error; err != nil {
					// 如果没有找到该用户，清除 Refresh Token 并返回未授权错误
					utils.ClearRefreshToken(c)
					response.NoAuth("The user does not exist", c)
					c.Abort()
					return
				}

				// 如果在数据库里找到了，则创建新的 access token
				newAccessClaims := j.CreateAccessClaims(request.BaseClaims{
					UserID: refreshClaims.UserID,
					UUID:   user.UUID,
					RoleID: user.RoleID,
				})
				newAccessToken, err2 := j.CreateAccessToken(newAccessClaims)
				if err2 != nil {
					// 如果生成失败，清除 Refresh Token 并返回错误
					utils.ClearRefreshToken(c)
					response.NoAuth("Failed to create access token", c)
					c.Abort()
					return
				}
				// 将新的 Access Token 和过期时间添加到响应头中
				c.Header("new-access-token", newAccessToken)
				c.Header("new-access-expires-at", strconv.FormatInt(newAccessClaims.ExpiresAt.Unix(), 10))
				// 将新的 Claims 信息存入 Context，供后续使用
				c.Set("claims", &newAccessClaims)
				c.Next()
				return
			}

			// 如果 access token 无效且不满足刷新条件，清除 refresh token并返回未授权错误
			utils.ClearRefreshToken(c)
			response.NoAuth("Invalid access token", c)
			c.Abort()
			return
		}

		// 如果 access token 合法，将其 claims 信息存入 Context
		c.Set("claims", claims)
		c.Next()
	}

}
