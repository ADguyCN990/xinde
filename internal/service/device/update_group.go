package device

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"xinde/pkg/stderr"
)

func (s *Service) UpdateGroup(deviceTypeID, groupID uint) error {
	return s.dao.DB().Transaction(func(tx *gorm.DB) error {
		// 检查groupID对应的分组是否存在
		_, err := s.groupDao.GetGroupByID(s.groupDao.DB(), groupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf(stderr.ErrorGroupNotFound)
			} else {
				return err
			}
		}

		// 更新groupID
		updateMap := map[string]interface{}{
			"group_id": groupID,
		}
		err = s.dao.UpdateDeviceType(tx, deviceTypeID, updateMap)
		if err != nil {
			return err
		}
		return nil
	})
}
