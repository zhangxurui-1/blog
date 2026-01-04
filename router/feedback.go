package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
)

type FeedbackRouter struct {
}

// InitArticleRouter 初始化反馈路由
func (f *FeedbackRouter) InitFeedbackRouter(privateGroup, publicGroup, adminGroup *gin.RouterGroup) {
	feedbackPrivateRouter := privateGroup.Group("feedback")
	feedbackPublicRouter := publicGroup.Group("feedback")
	feedbackAdminRouter := adminGroup.Group("feedback")

	feedbackApi := api.ApiGroupApp.FeedbackApi

	{
		feedbackPrivateRouter.POST("create", feedbackApi.FeedbackCreate)
		feedbackPrivateRouter.GET("info", feedbackApi.FeedbackInfo)
	}
	{
		feedbackPublicRouter.GET("new", feedbackApi.FeedbackNew)
	}
	{
		feedbackAdminRouter.DELETE("delete", feedbackApi.FeedbackDelete)
		feedbackAdminRouter.PUT("reply", feedbackApi.FeedbackReply)
		feedbackAdminRouter.GET("list", feedbackApi.FeedbackList)
	}
}
