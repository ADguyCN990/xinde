package group

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"mime/multipart"
	"xinde/internal/dao/attachment"
	"xinde/internal/dao/group"
	model "xinde/internal/model/attachment"
	"xinde/pkg/jwt"
	"xinde/pkg/logger"
	"xinde/pkg/util"
)

type Service struct {
	dao           *group.Dao
	j             *jwt.JWTService
	attachmentDao *attachment.Dao
}

func NewGroupService() (*Service, error) {
	dao, err := group.NewGroupDao()
	if err != nil {
		return nil, fmt.Errorf("创建Dao实例失败: %v", err.Error())
	}
	j := jwt.NewJWTService()
	attachmentDao, err := attachment.NewAttachmentDao()
	if err != nil {
		return nil, fmt.Errorf("创建Dao实例失败: %v", err.Error())
	}
	return &Service{
		dao:           dao,
		j:             j,
		attachmentDao: attachmentDao,
	}, nil
}

func (s *Service) Create(name string, parentID uint, adminID uint, iconFile *multipart.FileHeader) error {
	return s.dao.DB().Transaction(func(tx *gorm.DB) error {
		// 创建分组的数据库记录
		id, err := s.dao.Create(tx, name, parentID)
		if err != nil {
			return err
		}

		// 如果有图片文件上传，则处理文件并上传附件记录
		if iconFile != nil {
			storagePath, err := util.SaveUploadedFile(iconFile)
			if err != nil {
				return err
			}

			// 创建附件的数据库记录
			businessType := viper.GetString("business_type.group_icon")
			a := &model.Attachment{
				Filename:      iconFile.Filename,
				StoragePath:   storagePath,
				FileType:      iconFile.Header.Get("Content-Type"),
				FileSize:      uint64(iconFile.Size),
				StorageDriver: "local",
				UploadedByUID: adminID,
				BusinessType:  util.StringToPointer(businessType),
				BusinessID:    id,
			}

			if err := s.attachmentDao.Create(s.attachmentDao.DB(), a); err != nil {
				// 记录日志，但通常不因为这个失败而中断主流程
				logger.Error("记录上传附件信息到数据库失败: " + err.Error())
			}
		}
		return nil
	})
}
