package service

import (
	"gorm.io/gorm"
	"mime/multipart"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"
	"server/utils/upload"
)

type ImageService struct{}

// ImageUpload 图片上传，返回 url
func (imageService *ImageService) ImageUpload(file *multipart.FileHeader) (string, error) {
	oss := upload.NewOss()
	url, fileName, err := oss.UploadImage(file)
	if err != nil {
		return "", err
	}

	return url, global.DB.Create(&database.Image{
		Name:     fileName,
		URL:      url,
		Category: appTypes.Null,
		Storage:  global.Config.System.Storage(),
	}).Error
}
func (imageService *ImageService) ImageDelete(req request.ImageDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}

	var images []database.Image
	// 查找待删除记录
	if err := global.DB.Find(&images, req.IDs).Error; err != nil {
		return err
	}

	for _, image := range images {
		// 开启事务
		if err := global.DB.Transaction(func(tx *gorm.DB) error {
			oss := upload.NewOssWithStorage(image.Storage)
			// 删数据库记录
			if err := global.DB.Delete(&image).Error; err != nil {
				return err
			}
			// 删存储资源
			return oss.DeleteImage(image.Name)
		}); err != nil {
			return err
		}
	}
	return nil
}

// ImageList 获取图片列表
func (imageService *ImageService) ImageList(info request.ImageList) (list interface{}, total int64, err error) {
	db := global.DB
	// 根据条件（如果有）构造查询条件
	if info.Name != nil {
		db = db.Where("name LIKE ?", "%"+*info.Name+"%")
	}
	if info.Category != nil {
		category := appTypes.ToCategory(*info.Category)
		db = db.Where("category = ?", category)
	}
	if info.Storage != nil {
		storage := appTypes.ToStorage(*info.Storage)
		db = db.Where("storage = ?", storage)
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}
	// 分页查询
	return utils.MySQLPagination(&database.Image{}, option)
}
