package device

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"xinde/pkg/stderr"
)

func (s *Service) ChangeFilterImageDevice(id, deviceTypeID uint) error {
	// 验证要移动的记录是否存在
	filterImage, err := s.dao.GetFilterImageByID(s.dao.DB(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf(stderr.ErrorFilterImageNotFound)
		} else {
			return err
		}
	}

	// 验证目标DeviceType是否存在
	_, err = s.dao.GetDeviceTypeByID(s.dao.DB(), deviceTypeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf(stderr.ErrorDeviceNotFound)
		} else {
			return err
		}
	}

	// 检查目标设备类型下是否存在相同的配置
	isExists, _, err := s.dao.CheckFilterImageExists(s.dao.DB(), deviceTypeID, filterImage.FilterValue)
	if err != nil {
		return err
	}
	if isExists {
		return fmt.Errorf(stderr.ErrorFilterImageValueConflict)
	}

	// 一切都没有问题，更新deviceTypeID字段
	err = s.dao.UpdateFilterImage(s.dao.DB(), id, map[string]interface{}{
		"device_type_id": deviceTypeID,
	})
	if err != nil {
		return err
	}

	return nil
}
