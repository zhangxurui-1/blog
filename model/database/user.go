package database

import (
	"server/global"
	"server/model/appTypes"

	"github.com/gofrs/uuid"
)

type User struct {
	global.MODEL
	UUID      uuid.UUID         `json:"uuid" gorm:"type:char(36);unique"`
	Username  string            `json:"username"`
	Password  string            `json:"-"` // 传递 json 时忽略该字段
	Email     string            `json:"email"`
	Openid    string            `json:"openid"`
	Avatar    string            `json:"avatar" gorm:"size:255"`
	Address   string            `json:"address"`
	Signature string            `json:"signature" gorm:"default:'default signature'"`
	RoleID    appTypes.RoleID   `json:"role_id"`  // 用户角色
	Register  appTypes.Register `json:"register"` // 用户的注册类型（邮箱/QQ）
	Freeze    bool              `json:"freeze"`   // 是否被冻结
}
