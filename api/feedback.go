package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	"server/model/request"
	"server/model/response"
	"server/utils"
)

type FeedbackApi struct {
}

// FeedbackNew 获取最新的反馈路由列表
func (feedbackApi *FeedbackApi) FeedbackNew(c *gin.Context) {
	list, err := feedbackService.FeedbackNew()
	if err != nil {
		global.Log.Error("Failed to get new feedback", zap.Error(err))
		response.FailWithMessage("Failed to get new feedback", c)
		return
	}

	response.OkWithData(list, c)
}

// FeedbackCreate 创建反馈
func (feedbackApi *FeedbackApi) FeedbackCreate(c *gin.Context) {
	var req request.FeedbackCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	req.UUID = utils.GetUUID(c)
	if err := feedbackService.FeedbackCreate(req); err != nil {
		global.Log.Error("Failed to create feedback", zap.Error(err))
		response.FailWithMessage("Failed to create feedback", c)
		return
	}

	response.OkWithMessage("Successfully create feedback", c)
}

// FeedbackInfo 获取一个用户的反馈
func (feedbackApi *FeedbackApi) FeedbackInfo(c *gin.Context) {
	uuid := utils.GetUUID(c)
	list, err := feedbackService.FeedbackInfo(uuid)
	if err != nil {
		global.Log.Error("Failed to get feedback information", zap.Error(err))
		response.FailWithMessage("Failed to get feedback information", c)
		return
	}

	response.OkWithData(list, c)
}

// FeedbackDelete 删除反馈
func (feedbackApi *FeedbackApi) FeedbackDelete(c *gin.Context) {
	var req request.FeedbackDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
	}

	if err := feedbackService.FeedbackDelete(req); err != nil {
		global.Log.Error("Failed to delete feedback", zap.Error(err))
		response.FailWithMessage("Failed to delete feedback", c)
		return
	}
	response.OkWithMessage("Successfully deleted feedback", c)
}

// FeedbackReply 回复反馈
func (feedbackApi *FeedbackApi) FeedbackReply(c *gin.Context) {
	var req request.FeedbackReply
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
	}

	if err := feedbackService.FeedbackReply(req); err != nil {
		global.Log.Error("Failed to reply feedback", zap.Error(err))
		response.FailWithMessage("Failed to reply feedback", c)
		return
	}
	response.OkWithMessage("Successfully replied feedback", c)
}

// FeedbackList 获取反馈列表
func (feedbackApi *FeedbackApi) FeedbackList(c *gin.Context) {
	var pageInfo request.PageInfo
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := feedbackService.FeedbackList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get feedback list", zap.Error(err))
		response.FailWithMessage("Failed to get feedback list", c)
		return
	}

	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}
