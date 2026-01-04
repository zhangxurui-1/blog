package database

import (
	"github.com/gofrs/uuid"
	"server/global"
)

// 反馈表
type Feedback struct {
	global.MODEL
	UserUUID uuid.UUID `json:"user_uuid" gorm:"type:char(36)"`
	User     User      `json:"-" gorm:"foreignKey:UserUUID;references:UUID"`
	Content  string    `json:"content"` // 内容
	Reply    string    `json:"reply"`   // 回复
}
