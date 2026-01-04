// Package api 层只做参数传递和错误处理，具体业务逻辑放在 service 层
package api

import "server/service"

type ApiGroup struct {
	BaseApi
	UserApi
	ImageApi
	ArticleApi
	CommentApi
	AdvertisementApi
	FriendLinkApi
	FeedbackApi
	WebsiteApi
	ConfigApi
}

var ApiGroupApp = new(ApiGroup)
var (
	baseService          = service.ServiceGroupApp.BaseService
	userService          = service.ServiceGroupApp.UserService
	qqService            = service.ServiceGroupApp.QQService
	jwtService           = service.ServiceGroupApp.JwtService
	imageService         = service.ServiceGroupApp.ImageService
	articleService       = service.ServiceGroupApp.ArticleService
	commentService       = service.ServiceGroupApp.CommentService
	advertisementService = service.ServiceGroupApp.AdvertisementService
	friendLinkService    = service.ServiceGroupApp.FriendLinkService
	feedbackService      = service.ServiceGroupApp.FeedbackService
	websiteService       = service.ServiceGroupApp.WebsiteService
	configService        = service.ServiceGroupApp.ConfigService
)
