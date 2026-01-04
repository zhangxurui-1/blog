package database

import "server/global"

// FriendLink 友链表
type FriendLink struct {
	global.MODEL
	Logo        string `json:"logo" gorm:"size:255"`
	Image       Image  `json:"image" gorm:"foreignKey:Logo;references:URL"`
	Link        string `json:"link"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
