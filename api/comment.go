package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	"server/model/request"
	"server/model/response"
	"server/utils"
)

type CommentApi struct{}

// CommentInfoByArticleID 根据文章 id 获取它的所有评论
func (commentApi *CommentApi) CommentInfoByArticleID(c *gin.Context) {
	var req request.CommentInfoByArticleID
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
	}

	list, err := commentService.CommentInfoByArticleID(req)
	if err != nil {
		global.Log.Error("Failed to get comment information:", zap.Error(err))
		response.FailWithMessage("Failed to get comment information", c)
		return
	}
	response.OkWithData(list, c)
}

// CommentNew 获取文章的最新评论
func (commentApi *CommentApi) CommentNew(c *gin.Context) {
	list, err := commentService.CommentNew()
	if err != nil {
		global.Log.Error("Failed to get new comment:", zap.Error(err))
		response.FailWithMessage("Failed to get new comment", c)
		return
	}

	response.OkWithData(list, c)
}

// CommentCreate 创建评论
func (commentApi *CommentApi) CommentCreate(c *gin.Context) {
	var req request.CommentCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 从上下文获取用户的 uuid
	req.UserUUID = utils.GetUUID(c)
	if err := commentService.CommentCreate(req); err != nil {
		global.Log.Error("Failed to create comment:", zap.Error(err))
		response.FailWithMessage("Failed to create comment", c)
		return
	}

	response.OkWithMessage("Successfully created comment", c)
}

// CommentDelete 删除评论
func (commentApi *CommentApi) CommentDelete(c *gin.Context) {
	var req request.CommentDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := commentService.CommentDelete(c, req); err != nil {
		global.Log.Error("Failed to delete comment:", zap.Error(err))
		response.FailWithMessage("Failed to delete comment", c)
		return
	}
	response.OkWithMessage("Successfully deleted comment", c)
}

// CommentInfo 获取用户自己的评论
func (commentApi *CommentApi) CommentInfo(c *gin.Context) {
	uuid := utils.GetUUID(c)
	list, err := commentService.CommentInfo(uuid)
	if err != nil {
		global.Log.Error("Failed to get comment information:", zap.Error(err))
		response.FailWithMessage("Failed to get comment information", c)
		return
	}
	response.OkWithData(list, c)
}

// CommentList 获取评论列表
func (commentApi *CommentApi) CommentList(c *gin.Context) {
	var pageInfo request.CommentList
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := commentService.CommentList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get comment list:", zap.Error(err))
		response.FailWithMessage("Failed to get comment list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}
