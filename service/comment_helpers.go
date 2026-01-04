package service

import (
	"gorm.io/gorm"
	"server/global"
	"server/model/database"
)

// LoadChildren 加载一个评论的所有子评论
func (commentService *CommentService) LoadChildren(comment *database.Comment) error {
	var children []database.Comment
	if err := global.DB.Where("p_id = ?", comment.ID).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid, username, avatar, address, signature") // 仅获取需要的字段
	}).Find(&children).Error; err != nil {
		return err
	}

	// 递归加载
	for i := range children {
		if err := commentService.LoadChildren(&children[i]); err != nil {
			return err
		}
	}
	comment.Children = children

	return nil
}

// DeleteCommentAndChildren 删除一条评论和他的所有子评论
func (commentService *CommentService) DeleteCommentAndChildren(tx *gorm.DB, commentID uint) error {
	var children []database.Comment
	// 查找所有子评论
	if err := tx.Where("p_id = ?", commentID).Find(&children).Error; err != nil {
		return err
	}

	// 递归删除
	for _, child := range children {
		if err := commentService.DeleteCommentAndChildren(tx, child.ID); err != nil {
			return err
		}
	}
	// 删除子评论后再删除这条评论本身
	if err := tx.Delete(&database.Comment{}, commentID).Error; err != nil {
		return err
	}

	return nil
}

// FindChildCommentsByRootCommentUserUUID 查找与根评论 UserUUID 相同的所有子评论并记录
func (commentService *CommentService) FindChildCommentsByRootCommentUserUUID(
	comments []database.Comment) map[uint]struct{} {

	result := make(map[uint]struct{})
	// 遍历所有根评论
	for _, rootComment := range comments {
		// 创建一个递归函数，用于查找 UserUUID 字段与根评论相同的子评论
		var findChildren func([]database.Comment)
		findChildren = func(children []database.Comment) {
			// 遍历当前子评论
			for _, child := range children {
				// 若子评论的 UserUUID 与根评论相同，则记录
				if child.UserUUID == rootComment.UserUUID {
					result[child.ID] = struct{}{}
				}

				// 若子评论还有更下层的子评论，则继续递归
				if len(child.Children) > 0 {
					findChildren(child.Children)
				}
			}
		}

		// 调用该递归函数
		findChildren(rootComment.Children)
	}

	return result
}
