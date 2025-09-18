package device

import (
	"errors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func (s *Service) Delete(deviceTypeID uint) error {
	// 删除postgresql里的deviceType和device
	err := s.dao.DB().Transaction(func(tx *gorm.DB) error {
		// 首先检查deviceType是否存在
		_, err := s.dao.GetDeviceTypeByID(tx, deviceTypeID)
		if err != nil {
			return err
		}

		// 删除所有关联的方案
		err = s.dao.DeleteByDeviceTypeID(tx, deviceTypeID)
		if err != nil {
			return err
		}

		// 删除deviceType本身
		err = s.dao.DeleteDeviceTypeByID(tx, deviceTypeID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 删除t_attachment表里的excel和主图文件
	err = s.attachmentDao.DB().Transaction(func(tx *gorm.DB) error {
		businessTypes := []string{
			viper.GetString("business_type.device_icon"),
			viper.GetString("business_type.device_import"),
		}

		for _, businessType := range businessTypes {
			err := s.attachmentDao.DeleteAttachmentsByBusinessTypeAndIDs(tx, businessType, []uint{deviceTypeID})
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					continue
				}
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
