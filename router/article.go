package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
)

// ArticleRouter 文章路由
type ArticleRouter struct {
}

// InitArticleRouter 初始化文章路由
func (a *ArticleRouter) InitArticleRouter(Router *gin.RouterGroup, PublicGroup *gin.RouterGroup, AdminGroup *gin.RouterGroup) {
	articleRouter := Router.Group("article")
	articlePublicRouter := PublicGroup.Group("article")
	articleAdminRouter := AdminGroup.Group("article")

	articleApi := api.ApiGroupApp.ArticleApi
	{
		articleRouter.POST("like", articleApi.ArticleLike)
		articleRouter.GET("isLike", articleApi.ArticleIsLike)
		articleRouter.GET("likesList", articleApi.ArticleLikesList)
	}
	{
		articlePublicRouter.GET(":id", articleApi.ArticleInfoByID)
		articlePublicRouter.GET("search", articleApi.ArticleSearch)
		articlePublicRouter.GET("category", articleApi.ArticleCategory)
		articlePublicRouter.GET("tags", articleApi.ArticleTags)
	}
	{
		articleAdminRouter.POST("create", articleApi.ArticleCreate)
		articleAdminRouter.DELETE("delete", articleApi.ArticleDelete)
		articleAdminRouter.PUT("update", articleApi.ArticleUpdate)
		articleAdminRouter.GET("list", articleApi.ArticleList)
	}
}
