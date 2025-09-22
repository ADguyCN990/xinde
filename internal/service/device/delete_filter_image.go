package device

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func (s *Service) DeleteFilterImage(id uint) error {
	// 1. 在PG事务中删除记录
	err := s.dao.DB().Transaction(func(tx *gorm.DB) error {
		// a. 确认记录存在
		_, err := s.dao.GetFilterImageByID(tx, id)
		if err != nil {
			return err
		}

		// b. 删除记录
		err = s.dao.DeleteFilterImageByID(tx, id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 2. 删除关联的附件
	err = s.attachmentDao.DB().Transaction(func(tx *gorm.DB) error {
		businessType := viper.GetString("business_type.filter_image")
		err := s.attachmentDao.DeleteAttachmentsByBusinessTypeAndIDs(tx, businessType, []uint{id})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
