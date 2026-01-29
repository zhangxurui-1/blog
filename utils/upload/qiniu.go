package upload

import (
	"context"
	"fmt"
	"mime/multipart"
	"server/global"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
)

type Qiniu struct{}

func (q *Qiniu) UploadImage(file *multipart.FileHeader) (string, string, error) {
	return "", "", fmt.Errorf(`Current storage mode is qiniu, which should be done on client side, please 
	(1) get upload token via /api/image/upload_token
	(2) upload directly to qiniu cloud with the token`)
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

func (q *Qiniu) NewUpToken() (string, error) {
	mac := credentials.NewCredentials(global.Config.Qiniu.AccessKey, global.Config.Qiniu.SecretKey)
	bucket := global.Config.Qiniu.Bucket
	putPolicy, err := uptoken.NewPutPolicy(bucket, time.Now().Add(1*time.Hour))

	if err != nil {
		return "", err
	}

	putPolicy.SetCallbackUrl(global.Config.System.Domain + "/api/image/upload_callback").
		SetCallbackBody(`{
			"key":"${key}",
			"hash":"${etag}",
			"fsize":"${fsize}",
			"bucket":"${bucket}",
			"name":"${fname}"
		}`).
		SetCallbackBodyType("application/json")

	upToken, err := uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())

	if err != nil {
		return "", err
	}
	return upToken, nil
}
