package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/config"
	"server/global"
	"server/model/response"
)

type ConfigApi struct{}

// GetWebsite 获取网站配置
func (configApi *ConfigApi) GetWebsite(ctx *gin.Context) {
	response.OkWithData(global.Config.Website, ctx)
}

// UpdateWebsite 更新网站配置
func (configApi *ConfigApi) UpdateWebsite(ctx *gin.Context) {
	var req config.Website
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	err := configService.UpdateWebsite(req)
	if err != nil {
		global.Log.Error("Failed to update website:", zap.Error(err))
		response.FailWithMessage("Failed to update website", ctx)
		return
	}

	response.OkWithMessage("Successfully updated website", ctx)
}

// GetSystem 获取系统配置
func (configApi *ConfigApi) GetSystem(ctx *gin.Context) {
	response.OkWithData(global.Config.System, ctx)
}

// UpdateSystem 更新系统配置
func (configApi *ConfigApi) UpdateSystem(ctx *gin.Context) {
	var req config.System
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	if err := configService.UpdateSystem(req); err != nil {
		global.Log.Error("Failed to update system:", zap.Error(err))
		response.FailWithMessage("Failed to update system", ctx)
		return
	}

	response.OkWithMessage("Successfully updated system", ctx)
}

// GetEmail 获取邮箱配置
func (configApi *ConfigApi) GetEmail(ctx *gin.Context) {
	response.OkWithData(global.Config.Email, ctx)
}

// UpdateEmail 更新邮箱配置
func (configApi *ConfigApi) UpdateEmail(ctx *gin.Context) {
	var req config.Email
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	if err := configService.UpdateEmail(req); err != nil {
		global.Log.Error("Failed to update email:", zap.Error(err))
		response.FailWithMessage("Failed to update email", ctx)
		return
	}
	response.OkWithMessage("Successfully updated email", ctx)
}

// GetQQ 获取 QQ 登录配置
func (configApi *ConfigApi) GetQQ(ctx *gin.Context) {
	response.OkWithData(global.Config.QQ, ctx)
}

// UpdateQQ 更新 QQ 登录配置
func (configApi *ConfigApi) UpdateQQ(ctx *gin.Context) {
	var req config.QQ
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	if err := configService.UpdateQQ(req); err != nil {
		global.Log.Error("Failed to update qq:", zap.Error(err))
		response.FailWithMessage("Failed to update qq", ctx)
		return
	}

	return
}

// GetQiniu 获取七牛配置
func (configApi *ConfigApi) GetQiniu(ctx *gin.Context) {
	response.OkWithData(global.Config.Qiniu, ctx)
}

// UpdateQiniu 更新七牛配置
func (configApi *ConfigApi) UpdateQiniu(ctx *gin.Context) {
	var req config.Qiniu
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	if err := configService.UpdateQiniu(req); err != nil {
		global.Log.Error("Failed to update Qiniu", zap.Error(err))
		response.FailWithMessage("Failed to update Qiniu", ctx)
		return
	}

	response.OkWithMessage("Successfully updated Qiniu", ctx)
}

// GetJwt 获取 JWT 配置
func (configApi *ConfigApi) GetJwt(ctx *gin.Context) {
	response.OkWithData(global.Config.Jwt, ctx)
}

// UpdateJwt 更新 JWT 配置
func (configApi *ConfigApi) UpdateJwt(ctx *gin.Context) {
	var req config.Jwt
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	if err := configService.UpdateJwt(req); err != nil {
		global.Log.Error("Failed to update jwt", zap.Error(err))
		response.FailWithMessage("Failed to update jwt", ctx)
		return
	}
	response.OkWithMessage("Successfully updated jwt", ctx)
}

// GetGaode 获取高德配置
func (configApi *ConfigApi) GetGaode(ctx *gin.Context) {
	response.OkWithData(global.Config.Gaode, ctx)
}

// UpdateGaode 更新高德配置
func (configApi *ConfigApi) UpdateGaode(ctx *gin.Context) {
	var req config.Gaode
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	if err := configService.UpdateGaode(req); err != nil {
		global.Log.Error("Failed to update gaode", zap.Error(err))
		response.FailWithMessage("Failed to update gaode", ctx)
		return
	}
	response.OkWithMessage("Successfully updated gaode", ctx)
}
