package database

import (
	"server/global"
	"server/model/appTypes"
)

// Image 图片表
type Image struct {
	global.MODEL
	Name     string            `json:"name"`
	URL      string            `json:"url" gorm:"size:255;unique"`
	Category appTypes.Category `json:"category"` // 图片类别
	Storage  appTypes.Storage  `json:"storage"`  // 存储类型
}
