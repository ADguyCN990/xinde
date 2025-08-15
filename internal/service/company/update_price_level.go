package company

import (
	"fmt"
	"gorm.io/gorm"
	"xinde/pkg/stderr"
)

func (s *Service) UpdatePriceLevel(id uint, priceLevel string) error {
	return s.dao.DB().Transaction(func(tx *gorm.DB) error {
		// 检查公司是否存在
		isExist, err := s.dao.IsExistCompanyByID(tx, id)
		if err != nil {
			return err
		}
		if !isExist {
			return fmt.Errorf(stderr.ErrorCompanyNotFound)
		}

		// 调用dao修改公司的价格等级
		updateData := map[string]interface{}{
			"price_level": priceLevel,
		}
		err = s.dao.UpdateCompany(tx, id, updateData)
		if err != nil {
			return err
		}

		return nil
	})
}
