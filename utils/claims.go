package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"net"
	"server/global"
	"server/model/appTypes"
	"server/model/request"
)

// SetRefreshToken 设置 Refresh Token 的 cookie
func SetRefreshToken(c *gin.Context, token string, maxAge int) {
	// c.Request.Host 表示当前 HTTP 请求的主机部分（host），即请求的目标服务器的域名或 IP 地址
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	setCookie(c, "x-refresh-token", token, maxAge, host)
}

// ClearRefreshToken 清除 Refresh Token 的 cookie
func ClearRefreshToken(c *gin.Context) {
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	setCookie(c, "x-refresh-token", "", -1, host)
}

// setCookie 设置响应头中的 set-cookie
func setCookie(c *gin.Context, name, value string, maxAge int, host string) {
	if net.ParseIP(host) != nil {
		// 如果是 IP 地址，则不设置 domain
		c.SetCookie(name, value, maxAge, "/", "", false, true)
	} else {
		// 如果是域名，设置 cookie 的 domain 为域名
		c.SetCookie(name, value, maxAge, "/", host, false, true)
	}
}

// GetAccessToken 从请求头获取 access token
func GetAccessToken(c *gin.Context) string {
	return c.Request.Header.Get("x-access-token")
}

// GetRefreshToken 从请求头获取 refresh token
func GetRefreshToken(c *gin.Context) string {
	token, _ := c.Cookie("x-refresh-token")
	return token
}

// GetClaims 从 Context 中获取 access token 并解析
func GetClaims(c *gin.Context) (*request.JwtCustomClaims, error) {
	// 获取请求头中的 access token
	token := GetAccessToken(c)

	j := NewJWT()
	// 解析 access token
	claims, err := j.ParseAccessToken(token)
	if err != nil {
		global.Log.Error("Failed to retrieve JWT parsing information from Gin's Context. "+
			"Please check if the request header contains 'x-access-token' and if the claims structure is correct.",
			zap.Error(err),
		)
	}
	return claims, err
}

// GetRefreshClaims 从 Context 中获取 refresh token 并解析
func GetRefreshClaims(c *gin.Context) (*request.JwtCustomRefreshClaims, error) {
	token := GetRefreshToken(c)
	j := NewJWT()
	claims, err := j.ParseRefreshToken(token)
	if err != nil {
		global.Log.Error("Failed to retrieve JWT parsing information from Gin's Context. "+
			"Please check if the request header contains 'x-refresh-token' and if the claims structure is correct.",
			zap.Error(err),
		)
	}
	return claims, err
}

// GetUserInfo 从 Context 中获取用户信息
func GetUserInfo(c *gin.Context) *request.JwtCustomClaims {
	// 首先尝试从 Context 中直接获取 claims
	if claims, exists := c.Get("claims"); !exists {
		// 如果不存在，则重新解析 Access Token
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		// 如果已存在 claims，则直接返回
		return claims.(*request.JwtCustomClaims)
	}
}

// GetUserID 从 Context 中获取解析出来的用户 ID
func GetUserID(c *gin.Context) uint {
	// 首先尝试从 Context 中直接获取 claims
	if claims, exists := c.Get("claims"); !exists {
		// 如果不存在，则重新解析 Access Token
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.UserID
		}
	} else {
		// 如果已存在 claims，则直接返回
		return claims.(*request.JwtCustomClaims).UserID
	}
}

// GetUUID 从 Context 中获取解析出来的用户 UUID
func GetUUID(c *gin.Context) uuid.UUID {
	// 首先尝试从 Context 中获取 claims
	if claims, exists := c.Get("claims"); !exists {
		// 如果不存在，则重新解析 Access Token
		if cl, err := GetClaims(c); err != nil {
			return uuid.UUID{}
		} else {
			return cl.UUID
		}
	} else {
		// 如果已存在 claims，则直接返回 UUID
		return claims.(*request.JwtCustomClaims).UUID
	}
}

// GetRoleID 从 Context 中获取解析出来的用户 RoleID
func GetRoleID(c *gin.Context) appTypes.RoleID {
	// 首先尝试从 Context 中获取 claims
	if claims, exists := c.Get("claims"); !exists {
		// 如果不存在，则重新解析 Access Token
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.RoleID
		}
	} else {
		// 如果已存在 claims，则直接返回 RoleID
		return claims.(*request.JwtCustomClaims).RoleID
	}
}
