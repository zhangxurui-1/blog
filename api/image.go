package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	"server/model/request"
	"server/model/response"
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

	// local 返回格式 /uploads/image/fileName | qiniu 返回格式 http(s)://image.xxx.xx/fileName
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
