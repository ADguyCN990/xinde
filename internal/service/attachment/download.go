package attachment

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io"
	"os"
	"path/filepath"
	"xinde/pkg/stderr"
)

func (s *Service) GetAttachmentForDownload(id uint) (string, string, io.ReadCloser, error) {
	tx := s.dao.DB()

	attachment, err := s.dao.GetAttachmentByID(tx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", "", nil, fmt.Errorf(stderr.ErrorAttachmentNotFound)
	}
	if err != nil {
		return "", "", nil, err
	}

	// 根据存储驱动，来决定如何获取文件
	switch attachment.StorageDriver {
	case "local":
		// 构建文件的存储路径
		savePath := viper.GetString("attachment.save_path")
		if savePath == "" {
			return "", "", nil, fmt.Errorf("save_path未配置")
		}
		absolutePath := filepath.Join(savePath, attachment.StoragePath)

		// 打开文件
		file, err := os.Open(absolutePath)
		if os.IsNotExist(err) {
			return "", "", nil, fmt.Errorf(stderr.ErrorAttachmentNotFound)
		}
		if err != nil {
			return "", "", nil, fmt.Errorf("打开文件失败: " + err.Error())
		}
		return attachment.Filename, attachment.FileType, file, nil
	case "oss":
		return "", "", nil, fmt.Errorf("目前还未支持OSS存储")
	case "s3":
		return "", "", nil, fmt.Errorf("目前还未支持S3存储")
	default:
		return "", "", nil, fmt.Errorf("未知的存储驱动")
	}
}
