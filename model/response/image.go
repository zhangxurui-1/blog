package response

import "server/model/appTypes"

// ImageUpload 图片上传响应
type ImageUpload struct {
	Url     string `json:"url"`
	OssType string `json:"oss_type"`
}

type UpToken struct {
	UpToken string `json:"up_token"`
}

type ImageUploadCallback struct {
	Name     string            `json:"name"`
	URL      string            `json:"url" gorm:"size:255;unique"`
	Category appTypes.Category `json:"category"`
	Storage  appTypes.Storage  `json:"storage"`
}
