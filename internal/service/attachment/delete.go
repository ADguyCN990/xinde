package attachment

import (
	"fmt"
	"gorm.io/gorm"
	"os"
)

func (s *Service) Delete(id uint) error {
	return s.dao.DB().Transaction(func(tx *gorm.DB) error {
		tx = s.dao.DB()
		attachment, err := s.dao.GetAttachmentByID(tx, id)
		if err != nil {
			return err
		}

		// 根据存储驱动，删除物理文件
		switch attachment.StorageDriver {
		case "local":
			absolutePath, err := s.getAbsolutePath(attachment)
			if err != nil {
				return err
			}
			if err := os.Remove(absolutePath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("删除本地物理文件失败: %w", err)
			}
		case "oss":
			return fmt.Errorf("目前还未支持OSS存储的删除")
		case "s3":
			return fmt.Errorf("目前还未支持S3存储的删除")
		default:
			return fmt.Errorf("未知的存储类型")
		}

		// 删除数据库中的记录
		err = s.dao.DeleteAttachmentByID(tx, id)
		if err != nil {
			return err
		}
		return nil
	})
}
