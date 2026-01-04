package service

import (
	"gorm.io/gorm"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"
)

type AdvertisementService struct{}

// AdvertisementInfo 获取广告信息
func (advertisementService *AdvertisementService) AdvertisementInfo() (
	ads []database.Advertisement, total int64, err error) {

	err = global.DB.Model(&database.Advertisement{}).Count(&total).Find(&ads).Error
	if err != nil {
		return nil, 0, err
	}
	return ads, total, nil
}

// AdvertisementCreate 创建广告
func (advertisementService *AdvertisementService) AdvertisementCreate(
	req request.AdvertisementCreate) error {

	advertisementToCreate := database.Advertisement{
		AdImage: req.AdImage,
		Link:    req.Link,
		Title:   req.Title,
		Content: req.Content,
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := utils.ChangeImagesCategory(tx, []string{req.AdImage}, appTypes.AdImage); err != nil {
			return err
		}
		return tx.Create(&advertisementToCreate).Error
	})
}

// AdvertisementDelete 删除广告
func (advertisementService *AdvertisementService) AdvertisementDelete(
	req request.AdvertisementDelete) error {

	if len(req.IDs) == 0 {
		return nil
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.IDs {
			var advertisementToDelete database.Advertisement
			// 先查找待删除的记录
			if err := tx.Take(&advertisementToDelete, id).Error; err != nil {
				return err
			}
			// 修改图片的类型
			if err := utils.InitImagesCategory(tx, []string{advertisementToDelete.AdImage}); err != nil {
				return err
			}
			// 删除记录
			if err := tx.Delete(&advertisementToDelete).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// AdvertisementUpdate 更新广告
func (advertisementService *AdvertisementService) AdvertisementUpdate(
	req request.AdvertisementUpdate) error {
	updates := struct {
		Link    string `json:"link"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}{
		Link:    req.Link,
		Title:   req.Title,
		Content: req.Content,
	}
	return global.DB.Take(&database.Advertisement{}, req.ID).Updates(updates).Error
}

// AdvertisementList 获取广告列表
func (advertisementService *AdvertisementService) AdvertisementList(
	info request.AdvertisementList) (interface{}, int64, error) {

	db := global.DB
	if info.Title != nil {
		db = db.Where("title LIKE ?", "%"+*info.Title+"%")
	}
	if info.Content != nil {
		db = db.Where("content LIKE ?", "%"+*info.Content+"%")
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	return utils.MySQLPagination(&database.Advertisement{}, option)
}
