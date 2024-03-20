package local

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"

	"github.com/baowk/dilu-plugin/file_store/config"

	"github.com/baowk/dilu-core/common/utils/cryptos"
	"github.com/baowk/dilu-core/core"
)

func New(cfg *config.FSCfg) *Local {
	return &Local{
		cfg: cfg,
	}
}

type Local struct {
	cfg *config.FSCfg
}

//@object: *Local
//@function: UploadFile
//@description: 上传文件
//@param: file *multipart.FileHeader
//@return: string, string, error

func (e *Local) UploadFile(file *multipart.FileHeader) (filePath string, fileKey string, err error) {
	// 读取文件后缀
	ext := path.Ext(file.Filename)
	// 读取文件名并加密
	name := strings.TrimSuffix(file.Filename, ext)
	name = cryptos.MD5([]byte(name))
	// 拼接新文件名
	fileKey = name + "_" + time.Now().Format("20060102150405") + ext
	// 尝试创建此路径
	mkdirErr := os.MkdirAll(e.cfg.StorePath, os.ModePerm)
	if mkdirErr != nil {
		core.Log.Error("function os.MkdirAll() Filed", mkdirErr)
		err = mkdirErr
		return
	}
	// 拼接路径和文件名
	p := e.cfg.StorePath + "/" + fileKey
	filePath = e.cfg.PathPrefix + "/" + fileKey

	f, openError := file.Open() // 读取文件
	if openError != nil {
		core.Log.Error("function file.Open() Filed", openError)
		err = openError
		return
	}
	defer f.Close() // 创建文件 defer 关闭

	out, createErr := os.Create(p)
	if createErr != nil {
		core.Log.Error("function os.Create() Filed", createErr)
		err = createErr
		return
	}
	defer out.Close() // 创建文件 defer 关闭

	_, copyErr := io.Copy(out, f) // 传输（拷贝）文件
	if copyErr != nil {
		core.Log.Error("function io.Copy() Filed", copyErr)
		err = copyErr
		return
	}
	return
}

// @object: *Local
// @function: DeleteFile
// @description: 删除文件
// @param: key string
// @return: error
func (e *Local) DeleteFile(key string) error {
	p := e.cfg.StorePath + "/" + key
	if strings.Contains(p, e.cfg.StorePath) {
		if err := os.Remove(p); err != nil {
			return errors.New("本地文件删除失败, err:" + err.Error())
		}
	}
	return nil
}
