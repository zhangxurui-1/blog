package api

import (
	"fmt"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/utils/upload"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ImageApi struct {
}

// ImageUpload 上传图片
func (imageApi *ImageApi) ImageUpload(c *gin.Context) {
	_, header, err := c.Request.FormFile("image")
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	// local 返回格式 ./uploads/image/fileName | qiniu 返回格式 http(s)://image.xxx.xx/fileName
	url, err := imageService.ImageUpload(header)
	if err != nil {
		global.Log.Error("Failed to upload image", zap.Error(err))
		response.FailWithMessage("Failed to upload image", c)
		return
	}

	response.OkWithDetail(response.ImageUpload{
		Url:     url,
		OssType: global.Config.System.OssType,
	}, "Successfully uploaded image", c)
}

// ImageDelete 删除图片
func (imageApi *ImageApi) ImageDelete(c *gin.Context) {
	var req request.ImageDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := imageService.ImageDelete(req); err != nil {
		global.Log.Error("Failed to delete image", zap.Error(err))
		response.FailWithMessage("Failed to delete image", c)
		return
	}
	response.OkWithMessage("Successfully deleted image", c)
}

// ImageList 获取图片列表
func (imageApi *ImageApi) ImageList(c *gin.Context) {
	var pageInfo request.ImageList
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	imageList, total, err := imageService.ImageList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get image list", zap.Error(err))
		response.FailWithMessage("Failed to get image list", c)
		return
	}

	response.OkWithData(response.PageResult{
		List:  imageList,
		Total: total,
	}, c)
}

func (imageApi *ImageApi) ImageUploadToken(c *gin.Context) {
	oss := upload.NewOss()
	upToken, err := oss.NewUpToken()
	if err != nil {
		global.Log.Error("Failed to get upload token", zap.Error(err))
		response.FailWithMessage("Failed to get upload token", c)
		return
	}
	response.OkWithDetail(response.UpToken{
		UpToken: upToken,
	}, "Successfully get upload token", c)
}

func (imageApi *ImageApi) ImageUploadCallback(c *gin.Context) {
	var req request.ImageUploadCallback
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println("req=", req)

	global.Log.Info("Client uploaded an image on Qiniu", zap.Any("image_info", req))

	var image *database.Image
	var err error
	if image, err = imageService.ImageUploadCallback(&req); err != nil {
		global.Log.Error("Failed to upload image callback", zap.Error(err))
		response.FailWithMessage("Failed to upload image callback", c)
		return
	}

	response.OkWithDetail(response.ImageUploadCallback{
		Name:     image.Name,
		URL:      image.URL,
		Category: image.Category,
		Storage:  image.Storage,
	}, "Successfully uploaded image", c)
}
