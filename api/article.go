package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"
	"server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ArticleApi struct {
}

func (articleApi *ArticleApi) ArticleInfoByID(c *gin.Context) {
	var req request.ArticleInfoByID
	// 文章 id 是直接在 uri 里传的，所以使用 c.ShouldBindUri(&req)
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	article, err := articleService.ArticleInfoByID(req.ID)
	if err != nil {
		global.Log.Error("Failed to get article information:", zap.Error(err))
		response.FailWithMessage("Failed to get article information:", c)
		return
	}

	response.OkWithData(article, c)
}

// ArticleSearch 文章搜索
func (articleApi *ArticleApi) ArticleSearch(c *gin.Context) {
	var req request.ArticleSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := articleService.ArticleSearch(req)
	if err != nil {
		global.Log.Error("Failed to get article search results:", zap.Error(err))
		response.FailWithMessage("Failed to get article search results", c)
		return
	}

	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

// ArticleCategory 获取所有文章分类及数量
func (articleApi *ArticleApi) ArticleCategory(c *gin.Context) {
	category, err := articleService.ArticleCategory()
	if err != nil {
		global.Log.Error("Failed to get article category:", zap.Error(err))
		response.FailWithMessage("Failed to get article category", c)
		return
	}
	response.OkWithData(category, c)
}

// ArticleTags 获取所有文章标签及数量
func (articleApi *ArticleApi) ArticleTags(c *gin.Context) {
	tags, err := articleService.ArticleTags()
	if err != nil {
		global.Log.Error("Failed to get article tags:", zap.Error(err))
		response.FailWithMessage("Failed to get article tags", c)
		return
	}

	response.OkWithData(tags, c)
}

// ArticleLike 用户收藏文章
func (articleApi *ArticleApi) ArticleLike(c *gin.Context) {
	var req request.ArticleLike
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 从 context 中获取 UserID
	req.UserID = utils.GetUserID(c)

	if err := articleService.ArticleLike(req); err != nil {
		global.Log.Error("Failed to complete the operation:", zap.Error(err))
		response.FailWithMessage("Failed to complete the operation", c)
		return
	}

	response.OkWithMessage("Successfully complete the operation", c)
}

// ArticleIsLike 判断用户是否收藏了某个文章
func (articleApi *ArticleApi) ArticleIsLike(c *gin.Context) {
	var req request.ArticleLike
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 从 context 中获取 UserID
	req.UserID = utils.GetUserID(c)

	isLike, err := articleService.ArticleIsLike(req)
	if err != nil {
		global.Log.Error("Failed to get like status:", zap.Error(err))
		response.FailWithMessage("Failed to get like status", c)
		return
	}

	response.OkWithData(isLike, c)
}

// ArticleLikesList 获取用户的收藏文章列表
func (articleApi *ArticleApi) ArticleLikesList(c *gin.Context) {
	var pageInfo request.ArticleLikesList
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	pageInfo.UserID = utils.GetUserID(c)

	list, total, err := articleService.ArticleLikesList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get article likes list:", zap.Error(err))
		response.FailWithMessage("Failed to get article likes list", c)
		return
	}

	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

// ArticleCreate publishes an article
func (articleApi *ArticleApi) ArticleCreate(c *gin.Context) {
	var req request.ArticleCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := articleService.ArticleCreate(req); err != nil {
		global.Log.Error("Failed to create article:", zap.Error(err))
		response.FailWithMessage("Failed to create article", c)
		return
	}
	response.OkWithMessage("Successfully create article", c)
}

// ArticleDelete 删除文章
func (articleApi *ArticleApi) ArticleDelete(c *gin.Context) {
	var req request.ArticleDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := articleService.ArticleDelete(req); err != nil {
		global.Log.Error("Failed to delete article:", zap.Error(err))
		response.FailWithMessage("Failed to delete article", c)
		return
	}
	response.OkWithMessage("Successfully deleted article", c)
}

// ArticleUpdate 更新文章
func (articleApi *ArticleApi) ArticleUpdate(c *gin.Context) {
	var req request.ArticleUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := articleService.ArticleUpdate(req); err != nil {
		global.Log.Error("Failed to update article:", zap.Error(err))
		response.FailWithMessage("Failed to update article", c)
		return
	}
	response.OkWithMessage("Successfully updated article", c)
}

// ArticleList 获取文章列表
func (articleApi *ArticleApi) ArticleList(c *gin.Context) {
	var pageInfo request.ArticleList
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := articleService.ArticleList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get article list:", zap.Error(err))
		response.FailWithMessage("Failed to get article list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}
