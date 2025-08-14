package account

import (
	"fmt"
	"gorm.io/gorm"
	"xinde/pkg/stderr"
)

func (s *Service) DeleteUser(uid uint) error {
	return s.dao.Transaction(func(tx *gorm.DB) error {
		// 检查用户是否存在
		isExist, err := s.dao.IsExistUserByID(tx, uid)
		if err != nil {
			return err
		}
		if !isExist {
			return fmt.Errorf(stderr.ErrorUserNotFound)
		}

		// 调用dao执行软删除
		err = s.dao.DeleteUserByID(tx, uid)
		if err != nil {
			return err
		}
		return nil
	})
}
