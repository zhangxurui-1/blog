package service

import (
	"github.com/gofrs/uuid"
	"server/global"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"
)

type FeedbackService struct {
}

// FeedbackNew 获取最新的反馈
func (feedbackService *FeedbackService) FeedbackNew() (feedbacks []database.Feedback, err error) {
	if global.DB.Order("id desc").Find(&feedbacks).Error != nil {
		return nil, err
	}
	return feedbacks, nil
}

// FeedbackCreate 创建反馈
func (feedbackService *FeedbackService) FeedbackCreate(req request.FeedbackCreate) error {
	return global.DB.Create(&database.Feedback{
		UserUUID: req.UUID,
		Content:  req.Content,
	}).Error
}

// FeedbackInfo 查询一个用户的反馈
func (feedbackService *FeedbackService) FeedbackInfo(uuid uuid.UUID) (feedbacks []database.Feedback, err error) {
	if err = global.DB.Model(&database.Feedback{}).
		Order("id desc").Where("user_uuid = ?", uuid).Find(&feedbacks).Error; err != nil {
		return nil, err
	}

	return feedbacks, nil
}

// FeedbackDelete 删除反馈
func (feedbackService *FeedbackService) FeedbackDelete(req request.FeedbackDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}
	return global.DB.Delete(&database.Feedback{}, req.IDs).Error
}

// FeedbackReply 回复评论
func (feedbackService *FeedbackService) FeedbackReply(req request.FeedbackReply) error {
	return global.DB.Take(&database.Feedback{}, &req.ID).Update("reply", req.Reply).Error
}

// FeedbackList 获取反馈列表
func (feedbackService *FeedbackService) FeedbackList(pageInfo request.PageInfo) (list interface{}, total int64, err error) {
	option := other.MySQLOption{
		PageInfo: pageInfo,
	}

	return utils.MySQLPagination(&database.Feedback{}, option)
}
