package device

import (
	"fmt"
	"xinde/internal/dao/device"
	"xinde/pkg/jwt"
)

type Service struct {
	Dao *device.Dao
	j   *jwt.JWTService
}

func NewDeviceService() (*Service, error) {
	dao, err := device.NewDeviceDao()
	if err != nil {
		return nil, fmt.Errorf("创建Dao实例失败: " + err.Error())
	}
	j := jwt.NewJWTService()
	return &Service{
		Dao: dao,
		j:   j,
	}, nil
}
