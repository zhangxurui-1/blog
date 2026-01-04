package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
)

type FriendLinkRouter struct {
}

func (f *FriendLinkRouter) InitFriendLinkRouter(AdminRouter *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	friendLinkAdminRouter := AdminRouter.Group("friendLink")
	friendLinkPublicRouter := PublicRouter.Group("friendLink")

	friendLinkApi := api.ApiGroupApp.FriendLinkApi
	{
		friendLinkAdminRouter.POST("create", friendLinkApi.FriendLinkCreate)
		friendLinkAdminRouter.DELETE("delete", friendLinkApi.FriendLinkDelete)
		friendLinkAdminRouter.PUT("update", friendLinkApi.FriendLinkUpdate)
		friendLinkAdminRouter.GET("list", friendLinkApi.FriendLinkList)
	}
	{
		friendLinkPublicRouter.GET("info", friendLinkApi.FriendLinkInfo)
	}
}
