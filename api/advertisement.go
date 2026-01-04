package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	"server/model/request"
	"server/model/response"
)

type AdvertisementApi struct{}

// AdvertisementInfo 获取广告信息
func (advertisementApi *AdvertisementApi) AdvertisementInfo(c *gin.Context) {
	list, total, err := advertisementService.AdvertisementInfo()
	if err != nil {
		global.Log.Error("Failed to get advertisement information:", zap.Error(err))
		response.FailWithMessage("Failed to get advertisement information", c)
		return
	}
	response.OkWithData(response.AdvertisementInfo{List: list, Total: total}, c)
}

// AdvertisementCreate 创建广告
func (advertisementApi *AdvertisementApi) AdvertisementCreate(c *gin.Context) {
	var req request.AdvertisementCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := advertisementService.AdvertisementCreate(req); err != nil {
		global.Log.Error("Failed to create advertisement:", zap.Error(err))
		response.FailWithMessage("Failed to create advertisement", c)
		return
	}

	response.OkWithMessage("Successfully created advertisement", c)
}

// AdvertisementDelete 删除广告
func (advertisementApi *AdvertisementApi) AdvertisementDelete(c *gin.Context) {
	var req request.AdvertisementDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
	}

	if err := advertisementService.AdvertisementDelete(req); err != nil {
		global.Log.Error("Failed to delete advertisement:", zap.Error(err))
		response.FailWithMessage("Failed to delete advertisement", c)
		return
	}

	response.OkWithMessage("Successfully deleted advertisement", c)
}

// AdvertisementUpdate 更新广告
func (advertisementApi *AdvertisementApi) AdvertisementUpdate(c *gin.Context) {
	var req request.AdvertisementUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := advertisementService.AdvertisementUpdate(req); err != nil {
		global.Log.Error("Failed to update advertisement:", zap.Error(err))
		response.FailWithMessage("Failed to update advertisement", c)
		return
	}

	response.OkWithMessage("Successfully updated advertisement", c)
}

// AdvertisementList 获取广告列表
func (advertisementApi *AdvertisementApi) AdvertisementList(c *gin.Context) {
	var pageInfo request.AdvertisementList
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := advertisementService.AdvertisementList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get advertisement list:", zap.Error(err))
		response.FailWithMessage("Failed to get advertisement list", c)
		return
	}

	response.OkWithData(response.PageResult{List: list, Total: total}, c)
}
