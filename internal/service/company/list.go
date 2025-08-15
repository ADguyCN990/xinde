package company

import (
	"fmt"
	"xinde/internal/dao/company"
	"xinde/pkg/jwt"
)

type Service struct {
	dao *company.Dao
	jwt *jwt.JWTService
}

func NewCompanyService() (*Service, error) {
	dao, err := company.NewCompanyDao()
	if err != nil {
		return nil, fmt.Errorf("创建 DAO 实例失败: %w", err)
	}

	jwtService := jwt.NewJWTService()

	return &Service{
		dao: dao,
		jwt: jwtService,
	}, nil
}
