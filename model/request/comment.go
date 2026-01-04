package request

import "github.com/gofrs/uuid"

type CommentInfoByArticleID struct {
	ArticleID string `json:"article_id" form:"article_id" uri:"article_id" binding:"required"`
}

// CommentCreate 创建评论请求结构体
type CommentCreate struct {
	UserUUID  uuid.UUID `json:"-"`
	ArticleID string    `json:"article_id" binding:"required"`
	PID       *uint     `json:"pid"`
	Content   string    `json:"content" binding:"required,max=320"`
}

type CommentDelete struct {
	IDs []uint `json:"ids"`
}

type CommentList struct {
	ArticleID *string `json:"article_id" form:"article_id"`
	UserUUID  *string `json:"user_uuid" form:"user_uuid"`
	Content   *string `json:"content" form:"content"`
	PageInfo
}
