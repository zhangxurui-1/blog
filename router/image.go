package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type ImageRouter struct {
}

func (i *ImageRouter) InitImageRouter(adminRouter *gin.RouterGroup, publicRouter *gin.RouterGroup) {
	imageRouter := adminRouter.Group("image")
	imagePublicRouter := publicRouter.Group("image")
	imageApi := api.ApiGroupApp.ImageApi

	{
		imageRouter.POST("upload", imageApi.ImageUpload)
		imageRouter.DELETE("delete", imageApi.ImageDelete)
		imageRouter.GET("list", imageApi.ImageList)
		imageRouter.GET("upload_token", imageApi.ImageUploadToken)
	}
	{
		// callback url for image upload
		imagePublicRouter.POST("upload_callback", imageApi.ImageUploadCallback)
	}
}
