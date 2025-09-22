package device

import (
	"errors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"mime/multipart"
	model "xinde/internal/model/device"
)

func (s *Service) CreateFilterImage(adminUID, deviceTypeID uint, filterValue string, imageFile *multipart.FileHeader) error {

	var newFilterImage *model.FilterImage
	// 1. PG事务，创建filter_image记录
	err := s.dao.DB().Transaction(func(tx *gorm.DB) error {
		// a. 检查是否存在deviceType
		_, err := s.dao.GetDeviceTypeByID(tx, deviceTypeID)
		if err != nil {
			return err
		}

		// b. 检查是否存在相同的配置，如果存在还要删掉(附件也要删掉)
		isExists, filterImageID, err := s.dao.CheckFilterImageExists(tx, deviceTypeID, filterValue)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if isExists {
			businessType := viper.GetString("business_type.filter_image")
			err = s.attachmentDao.DeleteAttachmentsByBusinessTypeAndIDs(s.attachmentDao.DB(), businessType, []uint{filterImageID})
			if err != nil {
				return err
			}
			err := s.dao.DeleteFilterImageByID(tx, filterImageID)
			if err != nil {
				return err
			}
		}

		// c. 在t_filter_image表中创建记录
		newFilterImage = &model.FilterImage{
			DeviceTypeID: deviceTypeID,
			FilterValue:  filterValue,
		}
		err = s.dao.CreateFilterImage(tx, newFilterImage)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 2. MySQL事务，创建attachment记录
	err = s.attachmentDao.DB().Transaction(func(tx *gorm.DB) error {

		// 由于已经在前一个事务中删除了可能存在的附件记录，所以直接创建就行
		businessType := viper.GetString("business_type.filter_image")
		newRecord, err := s.getNewAttachmentRecord(imageFile, adminUID, newFilterImage.ID, businessType)
		if err != nil {
			return err
		}
		err = s.attachmentDao.Create(tx, newRecord)
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
