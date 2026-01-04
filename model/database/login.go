package database

import "server/global"

type Login struct {
	global.MODEL
	UserID      uint   `json:"user_id"`
	User        User   `json:"user" gorm:"foreignKey:UserID"` // 让 UserID 作为外键关联至 User 表中的主键 global.MODEL.ID
	LoginMethod string `json:"login_method"`
	IP          string `json:"ip"`
	Address     string `json:"address"`
	OS          string `json:"os"`
	DeviceInfo  string `json:"device_info"` // 设备信息
	BrowserInfo string `json:"browser_info"`
	Status      int    `json:"status"` // 登录状态
}
