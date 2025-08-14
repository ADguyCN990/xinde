package account

import (
	"fmt"
	"gorm.io/gorm"
	"xinde/pkg/stderr"
)

func (s *Service) ResetRemark(uid uint, remark string) error {
	return s.dao.Transaction(func(tx *gorm.DB) error {
		// 检查用户是否存在
		isExist, err := s.dao.IsExistUserByID(tx, uid)
		if err != nil {
			return err
		}
		if !isExist {
			return fmt.Errorf(stderr.ErrorUserNotFound)
		}

		// 调用dao进行update
		updateData := map[string]interface{}{
			"remarks": remark,
		}
		err = s.dao.UpdateUser(tx, uid, updateData)
		if err != nil {
			return err
		}

		return nil
	})
}
