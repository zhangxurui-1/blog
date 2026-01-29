package request

// ImageDelete 删除图片请求
type ImageDelete struct {
	IDs []int `json:"ids"`
}

// ImageList 图片列表请求
type ImageList struct {
	Name     *string `json:"name" form:"name"`
	Category *string `json:"category" form:"category"`
	Storage  *string `json:"storage" form:"storage"`
	PageInfo
}

type ImageUploadCallback struct {
	Key    string `json:"key"`
	Hash   string `json:"hash"`
	Fsize  string `json:"fsize"`
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}
