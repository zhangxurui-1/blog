package upload

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
	"path/filepath"
	"server/global"
	"server/utils"
	"strings"
	"time"
)

type Qiniu struct{}

func (q *Qiniu) UploadImage(file *multipart.FileHeader) (string, string, error) {
	// 检查文件大小
	size := float64(file.Size) / (1024 * 1024)
	if size >= float64(global.Config.Upload.Size) {
		return "", "", fmt.Errorf("the image size exceeds the limitation, current size : %.2f MB, "+
			"max upload size: %d MB", size, global.Config.Upload.Size)
	}

	// 提取文件扩展名
	ext := filepath.Ext(file.Filename)
	// 提取文件名（把扩展名剔除）
	name := strings.TrimSuffix(file.Filename, ext)
	// 禁止上传不符合格式的文件
	if _, exists := WhiteImageList[ext]; !exists {
		return "", "", fmt.Errorf("the file extension %s is not allowed", ext)
	}

	// 上传策略
	putPolicy := storage.PutPolicy{Scope: global.Config.Qiniu.Bucket}

	mac := qbox.NewMac(global.Config.Qiniu.AccessKey, global.Config.Qiniu.SecretKey)
	// 生成上传凭证
	upToken := putPolicy.UploadToken(mac)

	cfg := qiniuConfig()
	// 构建表单上传对象
	formUploader := storage.NewFormUploader(cfg)
	// PutRet 为七牛标准的上传回复内容
	putRet := storage.PutRet{}
	// PutExtra 为表单上传的额外可选项
	putExtra := storage.PutExtra{Params: map[string]string{}}

	fileKey := utils.MD5V([]byte(name)) + "-" + time.Now().Format("20060102150405") + ext

	data, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer data.Close()

	// 以表单方式上传一个文件
	err = formUploader.Put(context.Background(), &putRet, upToken, fileKey, data, file.Size, &putExtra)
	if err != nil {
		return "", "", err
	}

	return global.Config.Qiniu.ImgPath + putRet.Key, putRet.Key, nil
}

// DeleteImage 删除文件
func (q *Qiniu) DeleteImage(key string) error {
	mac := qbox.NewMac(global.Config.Qiniu.AccessKey, global.Config.Qiniu.SecretKey)
	cfg := qiniuConfig()
	bucketManager := storage.NewBucketManager(mac, cfg)
	return bucketManager.Delete(global.Config.Qiniu.Bucket, key)
}

// qiniuConfig 配置
func qiniuConfig() *storage.Config {
	cfg := storage.Config{
		UseHTTPS:      global.Config.Qiniu.UseHTTPS,
		UseCdnDomains: global.Config.Qiniu.UseCdnDomains,
	}
	switch global.Config.Qiniu.Zone {
	case "z0", "ZoneHuadong":
		cfg.Zone = &storage.ZoneHuadong
	case "z1", "ZoneHuabei":
		cfg.Zone = &storage.ZoneHuabei
	case "z2", "ZoneHuanan":
		cfg.Zone = &storage.ZoneHuanan
	case "na0", "ZoneBeimei":
		cfg.Zone = &storage.ZoneBeimei
	case "as0", "ZoneXinjiapo":
		cfg.Zone = &storage.ZoneXinjiapo
	case "ZoneHuadongZheJiang2":
		cfg.Zone = &storage.ZoneHuadongZheJiang2
	}

	return &cfg
}
