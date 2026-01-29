package upload

import (
	"mime/multipart"
	"server/global"
	"server/model/appTypes"
)

// WhiteImageList 白名单，支持上传的文件格式
var WhiteImageList = map[string]struct{}{
	".jpg":  {},
	".png":  {},
	".jpeg": {},
	".ico":  {},
	".tiff": {},
	".gif":  {},
	".svg":  {},
	".webp": {},
}

// OSS 对象存储接口
type OSS interface {
	UploadImage(file *multipart.FileHeader) (string, string, error)
	DeleteImage(key string) error
	NewUpToken() (string, error)
}

// NewOss 实例化 OSS
func NewOss() OSS {
	switch global.Config.System.OssType {
	case "local":
		return &Local{}
	case "qiniu":
		return &Qiniu{}
	default:
		return &Local{}
	}
}

// NewOssWithStorage 是根据传入的存储类型返回相应的 OSS 实例
// 为什么有了 NewOSS() 还要有 NewOssWithStorage() ?
// 因为删除图片的时候，并不知道它当初是以什么方式上传的，因此要动态配置
func NewOssWithStorage(storage appTypes.Storage) OSS {
	switch storage {
	case appTypes.Local:
		return &Local{}
	case appTypes.Qiniu:
		return &Qiniu{}
	default:
		return &Local{}
	}
}
