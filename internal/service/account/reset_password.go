package account

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"xinde/pkg/stderr"
	"xinde/pkg/util"
)

func (s *Service) ResetPassword(uid uint) error {
	return s.dao.Transaction(func(tx *gorm.DB) error {
		// 检查用户是否存在
		isExist, err := s.dao.IsExistUserByID(tx, uid)
		if err != nil {
			return err
		}
		if !isExist {
			return fmt.Errorf(stderr.ErrorUserNotFound)
		}

		// 调用dao将用户的密码重置为123456
		hashPassword, err := util.HashPassword(viper.GetString("account.defaultPassword"))
		if err != nil {
			return fmt.Errorf("加密密码失败: %w", err)
		}
		updateData := map[string]interface{}{
			"password": hashPassword,
		}
		err = s.dao.UpdateUser(tx, uid, updateData)
		if err != nil {
			return err
		}

		return nil
	})
}
