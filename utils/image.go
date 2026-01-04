package utils

import (
	"gorm.io/gorm"
	"server/model/appTypes"
	"server/model/database"
)

// InitImagesCategory 初始化图片类别
func InitImagesCategory(tx *gorm.DB, urls []string) error {
	// 将 urls 中的图片类似设置为 "未使用"
	return tx.Model(&database.Image{}).Where("url IN (?)", urls).Update("category", appTypes.Null).Error
}

// ChangeImagesCategory 修改一组图片的类别
func ChangeImagesCategory(tx *gorm.DB, urls []string, category appTypes.Category) error {
	return tx.Model(&database.Image{}).Where("url IN (?)", urls).Update("category", category).Error
}
