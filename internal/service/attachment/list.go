package attachment

import (
	"fmt"
	"xinde/internal/dao/attachment"
	"xinde/pkg/jwt"
)

type Service struct {
	jwt *jwt.JWTService
	Dao *attachment.Dao
}

func NewAttachmentService() (*Service, error) {
	jwtService := jwt.NewJWTService()
	dao, err := attachment.NewAttachmentDao()
	if err != nil {
		return nil, fmt.Errorf("创建Dao实例失败: " + err.Error())
	}
	return &Service{
		jwt: jwtService,
		Dao: dao,
	}, nil
}
