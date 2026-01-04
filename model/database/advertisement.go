package database

import "server/global"

// Advertisement 广告表
type Advertisement struct {
	global.MODEL
	AdImage string `json:"ad_image" gorm:"size:255"`
	Image   Image  `json:"image" gorm:"foreignKey:AdImage;references:URL"`
	Link    string `json:"link"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
