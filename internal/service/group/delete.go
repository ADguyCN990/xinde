package group

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"xinde/pkg/stderr"
)

const rootGroupID = 1

func (s *Service) Delete(groupID uint) error {
	return s.dao.DB().Transaction(func(tx *gorm.DB) error {
		// 业务场景：根分组不能被删除
		if groupID == rootGroupID {
			return fmt.Errorf(stderr.ErrorRootGroupCannotBeDeleted)
		}

		// 查找当前分组是否存在
		_, err := s.dao.GetGroupByID(tx, groupID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf(stderr.ErrorGroupNotFound)
			} else {
				return err
			}

		}

		// 获取当前分组及其所有子分组的ID
		idList, err := s.dao.FindAllDescendantIDs(tx, groupID)
		if err != nil {
			return err
		}
		idList = append(idList, groupID)

		// 【未来任务】处理关联的设备
		// ------------------------------------------------------------------
		// TODO: 当设备模块实现后，在这里添加逻辑：
		//  - 调用 device service 的方法:
		//    deviceService.MoveDevicesToGroup(tx, idsToDelete, RootGroupID)
		//  - 该方法会执行类似 UPDATE t_device SET group_id = 1 WHERE group_id IN (...) 的操作。
		//  - 因为设备存储可能不是MySQL，所以必须通过service层调用，不能直接操作DAO。
		// ------------------------------------------------------------------

		// 处理这些分组关联的图标
		var icons []uint
		businessType := viper.GetString("business_type.group_icon")
		for _, id := range idList {
			icon, err := s.attachmentDao.GetAttachmentByBusinessType(tx, businessType, id)
			if err != nil {
				return err
			}
			icons = append(icons, icon.ID)
		}
		for _, icon := range icons {
			err := s.attachmentDao.DeleteAttachmentByID(tx, icon)
			if err != nil {
				return err
			}
		}

		// 删除分组
		err = s.dao.DeleteGroupsByIDs(tx, idList)
		if err != nil {
			return err
		}
		return nil
	})
}
