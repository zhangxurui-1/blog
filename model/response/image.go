package response

// ImageUpload 图片上传响应
type ImageUpload struct {
	Url     string `json:"url"`
	OssType string `json:"oss_type"`
}
