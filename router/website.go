package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
)

type WebsiteRouter struct{}

func (w *WebsiteRouter) InitWebsiteRouter(adminGroup, publicGroup *gin.RouterGroup) {
	websiteAdminRouter := adminGroup.Group("website")
	websitePublicRouter := publicGroup.Group("website")

	websiteApi := api.ApiGroupApp.WebsiteApi
	{
		websiteAdminRouter.POST("addCarousel", websiteApi.WebsiteAddCarousel)
		websiteAdminRouter.PUT("cancelCarousel", websiteApi.WebsiteCancelCarousel)
		websiteAdminRouter.POST("createFooterLink", websiteApi.WebsiteCreateFooterLink)
		websiteAdminRouter.DELETE("deleteFooterLink", websiteApi.WebsiteDeleteFooterLink)
	}
	{
		websitePublicRouter.GET("logo", websiteApi.WebsiteLogo)
		websitePublicRouter.GET("title", websiteApi.WebsiteTitle)
		websitePublicRouter.GET("info", websiteApi.WebsiteInfo)
		websitePublicRouter.GET("carousel", websiteApi.WebsiteCarousel)
		websitePublicRouter.GET("news", websiteApi.WebsiteNews)
		websitePublicRouter.GET("calendar", websiteApi.WebsiteCalendar)
		websitePublicRouter.GET("footerLink", websiteApi.WebsiteFooterLink)
	}
}
