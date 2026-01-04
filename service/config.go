package service

import (
	"gorm.io/gorm"
	"server/config"
	"server/global"
	"server/model/appTypes"
	"server/utils"
)

type ConfigService struct{}

// UpdateWebsite 更新网站配置
func (configService *ConfigService) UpdateWebsite(website config.Website) error {
	oldArr := []string{
		global.Config.Website.Logo,
		global.Config.Website.FullLogo,
		global.Config.Website.QQImage,
		global.Config.Website.WechatImage,
	}

	newArr := []string{
		website.Logo,
		website.FullLogo,
		website.QQImage,
		website.WechatImage,
	}

	added, removed := utils.DiffArrays(oldArr, newArr)
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := utils.InitImagesCategory(global.DB, removed); err != nil {
			return err
		}
		if err := utils.ChangeImagesCategory(global.DB, added, appTypes.System); err != nil {
			return err
		}
		global.Config.Website = website
		if err := utils.SaveYAML(); err != nil {
			return err
		}
		return nil
	})
}

// UpdateSystem 更新系统设置
func (configService *ConfigService) UpdateSystem(system config.System) error {
	global.Config.System.UseMultipoint = system.UseMultipoint
	global.Config.System.SessionsSecret = system.SessionsSecret
	global.Config.System.OssType = system.OssType
	return utils.SaveYAML()
}

// UpdateEmail 更新邮箱配置
func (configService *ConfigService) UpdateEmail(email config.Email) error {
	global.Config.Email = email
	return utils.SaveYAML()
}

// UpdateQQ 更新 QQ 配置
func (configService *ConfigService) UpdateQQ(QQ config.QQ) error {
	global.Config.QQ = QQ
	return utils.SaveYAML()
}

// UpdateQiniu 更新七牛配置
func (configService *ConfigService) UpdateQiniu(req config.Qiniu) error {
	global.Config.Qiniu = req
	return utils.SaveYAML()
}

// UpdateJwt 更新 jwt 配置
func (configService *ConfigService) UpdateJwt(jwt config.Jwt) error {
	global.Config.Jwt = jwt
	return utils.SaveYAML()
}

// UpdateGaode 更新高德配置
func (configService *ConfigService) UpdateGaode(gaode config.Gaode) error {
	global.Config.Gaode = gaode
	return utils.SaveYAML()
}
