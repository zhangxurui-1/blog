package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"server/global"
	"server/utils"
	"strings"
	"time"
)

type Local struct{}

func (l *Local) UploadImage(file *multipart.FileHeader) (string, string, error) {
	// FileHeader.Size 单位是字节，这里转成 MB
	size := float64(file.Size) / float64(1024*1024)
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

	// 对文件名哈希后重新组成完整文件名，避免冲突
	filename := utils.MD5V([]byte(name)) + "-" + time.Now().Format("20060102150405") + ext
	path := global.Config.Upload.Path + "/image/"

	// 创建文件夹
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", "", err
	}

	// 在文件夹下创建文件
	filepath := path + filename
	out, err := os.Create(filepath)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	// 将上传的文件内容写入本地文件
	f, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	if _, err2 := io.Copy(out, f); err2 != nil {
		return "", "", err2
	}

	return "/" + filepath, filename, nil
}

// DeleteImage 删除文件
func (l *Local) DeleteImage(key string) error {
	path := global.Config.Upload.Path + "/image/" + key
	return os.Remove(path)
}
