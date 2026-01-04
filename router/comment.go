package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
)

type CommentRouter struct{}

func (c *CommentRouter) InitCommentRouter(privateGroup, publicGroup, adminGroup *gin.RouterGroup) {
	commentPrivateRouter := privateGroup.Group("comment")
	commentPublicRouter := publicGroup.Group("comment")
	commentAdminRouter := adminGroup.Group("comment")

	commentApi := api.ApiGroupApp.CommentApi

	{
		commentPrivateRouter.POST("create", commentApi.CommentCreate)
		commentPrivateRouter.DELETE("delete", commentApi.CommentDelete)
		commentPrivateRouter.GET("info", commentApi.CommentInfo)
	}
	{
		commentPublicRouter.GET(":article_id", commentApi.CommentInfoByArticleID)
		commentPublicRouter.GET("new", commentApi.CommentNew)
	}
	{
		commentAdminRouter.GET("list", commentApi.CommentList)
	}

}
