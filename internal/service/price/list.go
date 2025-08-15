package price

import (
	"fmt"
	"xinde/internal/dao/price"
	"xinde/pkg/jwt"
)

type Service struct {
	dao *price.Dao
	jwt *jwt.JWTService
}

func NewPriceService() (*Service, error) {
	dao, err := price.NewPriceDao()
	if err != nil {
		return nil, fmt.Errorf("创建Dao实例失败: %v", err)
	}
	jwtService := jwt.NewJWTService()
	return &Service{
		dao: dao,
		jwt: jwtService,
	}, nil
}
