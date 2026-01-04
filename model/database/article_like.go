package database

import "server/global"

// ArticleLike 文章收藏表
type ArticleLike struct {
	global.MODEL
	ArticleID string `json:"article_id"`                 // 文章 ID
	UserID    uint   `json:"user_id"`                    // 用户 ID
	User      User   `json:"-" gorm:"foreignKey:UserID"` // 让 UserID 与 User 表的主键关联
}
