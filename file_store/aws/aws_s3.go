package aws

import (
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/baowk/dilu-plugin/file_store/config"

	"github.com/baowk/dilu-core/core"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
)

func New(cfg *config.FSCfg) *AwsS3 {
	return &AwsS3{
		cfg: cfg,
	}
}

type AwsS3 struct {
	cfg *config.FSCfg
}

// @object: *AwsS3
// @function: UploadFile
// @description: Upload file to Aws S3 using aws-sdk-go. See https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html#s3-examples-bucket-ops-upload-file-to-bucket
// @param: file *multipart.FileHeader
// @return: string, string, error
func (e *AwsS3) UploadFile(file *multipart.FileHeader) (filePath string, fileKey string, err error) {
	session := newSession(e.cfg)
	uploader := s3manager.NewUploader(session)

	fileKey = fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)
	filename := e.cfg.PathPrefix + "/" + fileKey
	f, openError := file.Open()
	if openError != nil {
		core.Log.Error("function file.Open() Filed", zap.Any("err", openError.Error()))
		return
	}
	defer f.Close() // 创建文件 defer 关闭

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(e.cfg.Bucket),
		Key:    aws.String(filename),
		Body:   f,
	})
	if err != nil {
		core.Log.Error("function uploader.Upload() Filed", zap.Any("err", err.Error()))
		return
	}

	filePath = e.cfg.BaseURL + "/" + filename
	return
}

// @object: *AwsS3
// @function: DeleteFile
// @description: Delete file from Aws S3 using aws-sdk-go. See https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html#s3-examples-bucket-ops-delete-bucket-item
// @param: file *multipart.FileHeader
// @return: string, string, error
func (e *AwsS3) DeleteFile(key string) error {
	session := newSession(e.cfg)
	svc := s3.New(session)
	filename := e.cfg.PathPrefix + "/" + key
	bucket := e.cfg.Bucket

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		core.Log.Error("function svc.DeleteObject() Filed", zap.Any("err", err.Error()))
		return errors.New("function svc.DeleteObject() Filed, err:" + err.Error())
	}

	_ = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	return nil
}

// newSession Create S3 session
func newSession(cfg *config.FSCfg) *session.Session {
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String(cfg.Region),
		Endpoint:         aws.String(cfg.Endpoint), //minio在这里设置地址,可以兼容
		S3ForcePathStyle: aws.Bool(cfg.S3ForcePathStyle),
		DisableSSL:       aws.Bool(cfg.DisableSSL),
		Credentials: credentials.NewStaticCredentials(
			cfg.SecretID,
			cfg.SecretKey,
			"",
		),
	})
	return sess
}
