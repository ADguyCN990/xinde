package group

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"mime/multipart"
	model "xinde/internal/model/attachment"
	"xinde/pkg/stderr"
	"xinde/pkg/util"
)

func (s *Service) Update(adminID uint, groupID uint, parentID uint, name string, icon *multipart.FileHeader) error {
	return s.dao.DB().Transaction(func(tx *gorm.DB) error {
		// 检查分组是否存在
		if parentID != 0 {
			// 1. 【新增】校验：分组不能成为自己的父级
			if groupID == parentID {
				return fmt.Errorf(stderr.ErrorCannotMoveGroupIntoItself)
			}

			// 2. 【新增】校验：分组不能被移动到自己的子树下
			// a. 获取当前分组的所有子孙节点ID
			descendantIDs, err := s.dao.FindAllDescendantIDs(tx, groupID)
			if err != nil {
				return err
			}

			// b. 检查新的父级ID是否在子孙列表中
			for _, descID := range descendantIDs {
				if parentID == descID {
					return fmt.Errorf(stderr.ErrorCannotMoveGroupIntoItself)
				}
			}
			_, err = s.dao.GetGroupByID(tx, parentID)
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				} else {
					return fmt.Errorf(stderr.ErrorGroupNotFound)
				}
			}
		}

		// 更新分组
		updateMap := make(map[string]interface{})
		if parentID != 0 {
			updateMap["parent_id"] = parentID
		}
		if name != "" {
			updateMap["name"] = name
		}
		err := s.dao.UpdateGroupByID(tx, groupID, updateMap)
		if err != nil {
			return err
		}

		// 如果更改了分组icon，还需要更新attachment的部分
		if icon != nil {
			// 如果有旧的图片，删除它在数据库中的记录
			businessType := viper.GetString("business_type.group_icon")
			oldIcon, err := s.attachmentDao.GetAttachmentByBusinessType(tx, businessType, groupID)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if oldIcon != nil {
				err := s.attachmentDao.DeleteAttachmentByID(tx, oldIcon.ID)
				if err != nil {
					return err
				}
			}

			// 保存新上传的文件
			storagePath, err := util.SaveUploadedFile(icon)
			if err != nil {
				return err
			}

			// 往attachment表里写入新的记录
			newRecord := &model.Attachment{
				Filename:      icon.Filename,
				StoragePath:   storagePath,
				FileType:      icon.Header.Get("Content-Type"),
				FileSize:      uint64(icon.Size),
				StorageDriver: "local",
				UploadedByUID: adminID,
				BusinessType:  util.StringToPointer(businessType),
				BusinessID:    groupID,
			}
			err = s.attachmentDao.Create(tx, newRecord)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
