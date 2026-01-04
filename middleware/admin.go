package middleware

import (
	"github.com/gin-gonic/gin"
	"server/model/appTypes"
	"server/model/response"
	"server/utils"
)

// AdminAuth 检查用户是否具有管理员权限
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户的角色ID
		roleID := utils.GetRoleID(c)

		if roleID != appTypes.Admin {
			// 如果不是管理员
			response.Forbidden("Access denied. Admin privileges are required", c)

			c.Abort()
			return
		}
		// 如果用户是管理员，继续执行后续处理
		c.Next()
	}
}
