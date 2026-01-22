package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ua-parser/uap-go/uaparser"
	"go.uber.org/zap"
	"server/global"
	"server/model/database"
	"server/service"
)

func LoginRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 异步记录日志
		go func() {
			gaodeService := service.ServiceGroupApp.GaodeService
			var userID uint
			var address string
			ip := c.ClientIP()
			// 若未传递flag参数，则默认为 "email"
			loginMethod := c.DefaultQuery("flag", "email")
			userAgent := c.Request.UserAgent()

			// 从请求上下文中获取 UserID，确保获取到的是当前请求的正确 UserID
			if value, exists := c.Get("user_id"); exists {
				
				if id, ok := value.(uint); ok {
					userID = id
				}
			}
			// 获取用户的地理位置
			address = getAddressFromIP(ip, gaodeService)
			// 解析 user-agent
			os, deviceInfo, browserInfo := parseUserAgent(userAgent)
			// 创建登录记录
			login := database.Login{
				UserID:      userID,
				LoginMethod: loginMethod,
				Address:     address,
				IP:          ip,
				OS:          os,
				DeviceInfo:  deviceInfo,
				BrowserInfo: browserInfo,
				Status:      c.Writer.Status(),
			}
			// 将登录记录存储到数据库
			if err := global.DB.Create(&login).Error; err != nil {
				global.Log.Error("Failed to record login", zap.Error(err))
			}
		}()
	}

}

// getAddressFromIP 获取 IP 对应的地理位置
func getAddressFromIP(ip string, gaodeService service.GaodeService) string {
	location, err := gaodeService.GetLocationByIP(ip)
	if err != nil || location.Province == "" {
		return "未知"
	}
	if location.City != "" && location.Province != location.City {
		return location.Province + "-" + location.City
	}
	return location.Province
}

// parseUserAgent 解析 UserAgent，返回操作系统、设备信息和浏览器信息
func parseUserAgent(userAgent string) (os, deviceInfo, browserInfo string) {
	os = userAgent
	deviceInfo = userAgent
	browserInfo = userAgent

	parser := uaparser.NewFromSaved()
	cli := parser.Parse(userAgent)

	os = cli.Os.Family
	deviceInfo = cli.Device.Family
	browserInfo = cli.UserAgent.Family
	return
}
