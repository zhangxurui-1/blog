package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"
)

type CommentService struct{}

// CommentInfoByArticleID 根据文章 id 查询所有的评论；返回一级评论
func (commentService *CommentService) CommentInfoByArticleID(req request.CommentInfoByArticleID) ([]database.Comment, error) {
	var comments []database.Comment

	// 加载一级评论
	if err := global.DB.Where("article_id = ? and p_id is NULL", req.ArticleID).
		Preload("User", func(db *gorm.DB) *gorm.DB { // 预加载 User
			return db.Select("uuid, username, avatar, address, signature")
		}).Find(&comments).Error; err != nil {
		return nil, err
	}

	// 递归加载子评论
	for i := range comments {
		if err := commentService.LoadChildren(&comments[i]); err != nil {
			return nil, err
		}
	}

	return comments, nil
}

// CommentNew 获取文章的最新评论
func (commentService *CommentService) CommentNew() ([]database.Comment, error) {
	var comments []database.Comment
	// 查询最新的 5 条记录
	err := global.DB.Order("id desc").Limit(5).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid, username, avatar, address, signature")
	}).Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// CommentCreate 创建评论
func (commentService *CommentService) CommentCreate(req request.CommentCreate) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if req.PID != nil {
			if errors.Is(tx.Take(&database.Comment{}, req.PID).Error, gorm.ErrRecordNotFound) {
				return errors.New("parent comment not found")
			}
		}
		return tx.Create(&database.Comment{
			ArticleID: req.ArticleID,
			PID:       req.PID,
			UserUUID:  req.UserUUID,
			Content:   req.Content,
		}).Error
	})

}

// CommentDelete 删除文章评论
func (commentService *CommentService) CommentDelete(c *gin.Context, req request.CommentDelete) error {
	// 需要先删除外键关联的字段，再删除记录
	// 避免出现外键关联的记录已经被删除的情况
	if len(req.IDs) == 0 {
		return nil
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.IDs {
			var comment database.Comment
			if err := global.DB.Take(&comment, id).Error; err != nil {
				return err
			}
			// 鉴权
			userUUID := utils.GetUUID(c)
			userRoleID := utils.GetRoleID(c)
			if userUUID != comment.UserUUID && userRoleID != appTypes.Admin {
				return errors.New("do not have permission to delete the comment")
			}

			if err := commentService.DeleteCommentAndChildren(tx, id); err != nil {
				return err
			}
		}
		return nil
	})
}

// CommentInfo 获取用户评论
func (commentService *CommentService) CommentInfo(uuid uuid.UUID) ([]database.Comment, error) {
	var rawComments []database.Comment
	err := global.DB.Order("id desc").Where("user_uuid = ?", uuid).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("uuid, username, avatar, address, signature")
		}).Find(&rawComments).Error
	if err != nil {
		return nil, err
	}
	// 加载子评论
	for i := range rawComments {
		if err := commentService.LoadChildren(&rawComments[i]); err != nil {
			return nil, err
		}
	}
	// 评论去重，如果当前评论的子评论包含自己的评论，需要删除子评论
	var comments []database.Comment

	idMap := commentService.FindChildCommentsByRootCommentUserUUID(rawComments)
	for i := range rawComments {
		// 如果在这个表里，说明重复，不作为结果返回
		if _, exists := idMap[rawComments[i].ID]; !exists {
			comments = append(comments, rawComments[i])
		}
	}

	return comments, nil
}

// CommentList 获取评论列表，可根据 文章id/UserUUID/内容 筛选
func (commentService *CommentService) CommentList(
	info request.CommentList) (interface{}, int64, error) {

	db := global.DB
	if info.ArticleID != nil {
		db = db.Where("article_id = ?", *info.ArticleID)
	}
	if info.UserUUID != nil {
		db = db.Where("user_uuid = ?", *info.UserUUID)
	}
	if info.Content != nil {
		db = db.Where("content LIKE ?", "%"+*info.Content+"%")
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	return utils.MySQLPagination(&database.Comment{}, option)
}
