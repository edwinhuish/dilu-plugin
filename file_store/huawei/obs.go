package huawei

import (
	"mime/multipart"

	"github.com/baowk/dilu-plugin/file_store/config"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/pkg/errors"
)

func New(cfg *config.FSCfg) *Obs {
	return &Obs{
		cfg: cfg,
	}
}

type Obs struct {
	cfg *config.FSCfg
}

func NewHuaWeiObsClient(cfg *config.FSCfg) (client *obs.ObsClient, err error) {
	return obs.New(cfg.SecretID, cfg.SecretKey, cfg.Endpoint)
}

func (o *Obs) UploadFile(file *multipart.FileHeader) (filePath string, fileKey string, err error) {
	// var open multipart.File
	open, err := file.Open()
	if err != nil {
		return
	}
	defer open.Close()
	fileKey = file.Filename
	input := &obs.PutObjectInput{
		PutObjectBasicInput: obs.PutObjectBasicInput{
			ObjectOperationInput: obs.ObjectOperationInput{
				Bucket: o.cfg.Bucket,
				Key:    fileKey,
			},
			//ContentType: file.Header.Get("content-type"),
		},
		Body: open,
	}

	var client *obs.ObsClient
	client, err = NewHuaWeiObsClient(o.cfg)
	if err != nil {
		return
	}

	_, err = client.PutObject(input)
	if err != nil {
		return
	}
	filePath = o.cfg.PathPrefix + "/" + fileKey
	return
}

func (o *Obs) DeleteFile(key string) error {
	client, err := NewHuaWeiObsClient(o.cfg)
	if err != nil {
		return errors.Wrap(err, "获取华为对象存储对象失败!")
	}
	input := &obs.DeleteObjectInput{
		Bucket: o.cfg.Bucket,
		Key:    key,
	}
	var output *obs.DeleteObjectOutput
	output, err = client.DeleteObject(input)
	if err != nil {
		return errors.Wrapf(err, "删除对象(%s)失败!, output: %v", key, output)
	}
	return nil
}
