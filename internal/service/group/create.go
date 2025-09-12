package group

import (
	"fmt"
	"xinde/internal/dao/attachment"
	"xinde/internal/dao/group"
	"xinde/pkg/jwt"
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
