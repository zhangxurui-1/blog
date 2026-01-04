package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
)

// AdvertisementRouter 广告路由
type AdvertisementRouter struct{}

// InitAdvertisementRouter 初始化
func (a *AdvertisementRouter) InitAdvertisementRouter(AdminRouter, PublicRouter *gin.RouterGroup) {
	advertisementAdminRouter := AdminRouter.Group("advertisement")
	advertisementPublicRouter := PublicRouter.Group("advertisement")

	advertisementApi := api.ApiGroupApp.AdvertisementApi
	{
		advertisementAdminRouter.POST("create", advertisementApi.AdvertisementCreate)
		advertisementAdminRouter.DELETE("delete", advertisementApi.AdvertisementDelete)
		advertisementAdminRouter.PUT("update", advertisementApi.AdvertisementUpdate)
		advertisementAdminRouter.GET("list", advertisementApi.AdvertisementList)
	}
	{
		advertisementPublicRouter.GET("info", advertisementApi.AdvertisementInfo)
	}
}
