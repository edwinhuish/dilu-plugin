package tencent

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/baowk/dilu-plugin/file_store/config"

	"github.com/baowk/dilu-core/core"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func New(cfg *config.FSCfg) *TencentCOS {
	return &TencentCOS{
		cfg: cfg,
	}
}

type TencentCOS struct {
	cfg *config.FSCfg
}

// UploadFile upload file to COS
func (e *TencentCOS) UploadFile(file *multipart.FileHeader) (filePath string, fileKey string, err error) {
	client := NewClient(e.cfg)
	f, openError := file.Open()
	if openError != nil {
		core.Log.Error("function file.Open() Filed", openError)
		err = errors.New("function file.Open() Filed, err:" + openError.Error())
		return
	}
	defer f.Close() // 创建文件 defer 关闭
	fileKey = fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)

	_, err = client.Object.Put(context.Background(), e.cfg.PathPrefix+"/"+fileKey, f, nil)
	if err != nil {
		return
	}
	filePath = e.cfg.BaseURL + "/" + e.cfg.PathPrefix + "/" + fileKey
	return
}

// DeleteFile delete file form COS
func (e *TencentCOS) DeleteFile(key string) error {
	client := NewClient(e.cfg)
	name := e.cfg.PathPrefix + "/" + key
	_, err := client.Object.Delete(context.Background(), name)
	if err != nil {
		core.Log.Error("function bucketManager.Delete() Filed", err)
		return errors.New("function bucketManager.Delete() Filed, err:" + err.Error())
	}
	return nil
}

// NewClient init COS client
func NewClient(cfg *config.FSCfg) *cos.Client {
	urlStr, _ := url.Parse("https://" + cfg.Bucket + ".cos." + cfg.Region + ".myqcloud.com")
	baseURL := &cos.BaseURL{BucketURL: urlStr}
	client := cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})
	return client
}
