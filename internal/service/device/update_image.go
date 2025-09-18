package device

import (
	"github.com/spf13/viper"
	"mime/multipart"
)

func (s *Service) UpdateImage(deviceTypeID, adminUID uint, imageFile *multipart.FileHeader) error {

	// 首先确认deviceType是否存在
	_, err := s.dao.GetDeviceTypeByID(s.dao.DB(), deviceTypeID)
	if err != nil {
		return err
	}

	// 根据business_type和business_id查找并删除的旧图片
	businessType := viper.GetString("business_type.device_icon")
	err = s.attachmentDao.DeleteAttachmentsByBusinessTypeAndIDs(s.attachmentDao.DB(), businessType, []uint{deviceTypeID})
	if err != nil {
		return err
	}

	// 保存新上传的文件
	newRecord, err := s.getNewAttachmentRecord(imageFile, adminUID, deviceTypeID, businessType)
	if err != nil {
		return err
	}

	// 在t_attachment表中创建记录
	err = s.attachmentDao.Create(s.attachmentDao.DB(), newRecord)
	if err != nil {
		return err
	}

	return nil
}
