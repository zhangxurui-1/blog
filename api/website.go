package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"time"
)

type WebsiteApi struct{}

// WebsiteLogo 获取网站 logo
func (websiteApi *WebsiteApi) WebsiteLogo(c *gin.Context) {
	// 重定向到 logo 的 url
	if global.Config.Website.Logo != "" {
		c.Redirect(http.StatusMovedPermanently, global.Config.Website.Logo)
	} else {
		c.Redirect(http.StatusMovedPermanently, "/image/logo.png")
	}
}

// WebsiteTitle 获取网站标题
func (websiteApi *WebsiteApi) WebsiteTitle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"title": global.Config.Website.Title})
}

// WebsiteInfo 获取网站信息
func (websiteApi *WebsiteApi) WebsiteInfo(c *gin.Context) {
	response.OkWithData(global.Config.Website, c)
}

// WebsiteCarousel 获取网站首页走马灯
func (websiteApi *WebsiteApi) WebsiteCarousel(c *gin.Context) {
	urls := websiteService.WebsiteCarousel()
	response.OkWithData(urls, c)
}

// WebsiteNews 获取新闻数据
func (websiteApi *WebsiteApi) WebsiteNews(c *gin.Context) {
	sourceStr := c.Query("source")
	hotSearchData, err := websiteService.WebsiteNews(sourceStr)
	if err != nil {
		global.Log.Error("Failed to get news:", zap.Error(err))
		response.FailWithMessage("Failed to get news", c)
		return
	}

	response.OkWithData(hotSearchData, c)
}

// WebsiteCalendar 获取日历
func (websiteApi *WebsiteApi) WebsiteCalendar(c *gin.Context) {
	dateStr := time.Now().Format("2006/0102")
	calendar, err := websiteService.WebsiteCalendar(dateStr)
	if err != nil {
		global.Log.Error("Failed to get calendar:", zap.Error(err))
		response.FailWithMessage("Failed to get calendar", c)
		return
	}
	response.OkWithData(calendar, c)
}

// WebsiteFooterLink 获取页脚链接
func (websiteApi *WebsiteApi) WebsiteFooterLink(c *gin.Context) {
	footerLinks := websiteService.WebsiteFooterLink()
	response.OkWithData(footerLinks, c)
}

// WebsiteAddCarousel 添加背景图片（这里的添加是指将对应的图片类型设为"背景"）
func (websiteApi *WebsiteApi) WebsiteAddCarousel(c *gin.Context) {
	var req request.WebsiteCarouselOperation
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := websiteService.WebsiteAddCarousel(req); err != nil {
		global.Log.Error("Failed to add carousel", zap.Error(err))
		response.FailWithMessage("Failed to add carousel", c)
		return
	}
	response.OkWithMessage("Successfully added carousel", c)
}

// WebsiteCancelCarousel 移除首页背景（这里的移除是指将对应的图片类型设为空）
func (websiteApi *WebsiteApi) WebsiteCancelCarousel(c *gin.Context) {
	var req request.WebsiteCarouselOperation
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := websiteService.WebsiteCancelCarousel(req); err != nil {
		global.Log.Error("Failed to cancel carousel", zap.Error(err))
		response.FailWithMessage("Failed to cancel carousel", c)
	}

	response.OkWithMessage("Successfully cancelled carousel", c)
}

// WebsiteCreateFooterLink 创建页脚链接
func (websiteApi *WebsiteApi) WebsiteCreateFooterLink(c *gin.Context) {
	var req database.FooterLink
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := websiteService.WebsiteCreateFooterLink(req); err != nil {
		global.Log.Error("Failed to create footer link", zap.Error(err))
		response.FailWithMessage("Failed to create footer link", c)
		return
	}

	response.OkWithData("Successfully created footer link", c)
}

// WebsiteDeleteFooterLink 删除页脚链接
func (websiteApi *WebsiteApi) WebsiteDeleteFooterLink(c *gin.Context) {
	var req database.FooterLink
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := websiteService.WebsiteDeleteFooterLink(req); err != nil {
		global.Log.Error("Failed to delete footer link", zap.Error(err))
		response.FailWithMessage("Failed to delete footer link", c)
		return
	}
	response.OkWithMessage("Successfully deleted footer link", c)
}
