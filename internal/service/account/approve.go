package account

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
	"xinde/pkg/stderr"
)

func (s *Service) ApproveUser(id uint, status int, why string) (err error) {
	// 启动一个事务
	tx := s.dao.DB().Begin()
	if tx.Error != nil {
		return fmt.Errorf("开启事务失败: %w", tx.Error)
	}

	// 使用 defer 来确保事务在函数退出时能被处理（提交或回滚）
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // 如果发生 panic，回滚事务
		}
	}()

	// 查找用户是否存在
	user, err := s.dao.GetUserByIDForUpdate(tx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf(stderr.ErrorUserNotFound)
		}
		return err
	}

	// 检查用户的`is_user`状态，此时只可能为0
	switch user.IsUser {
	case 1:
		return fmt.Errorf(stderr.ErrorUserPassed)
	case 2:
		return fmt.Errorf(stderr.ErrorUserBanned)
	}

	// 更新用户的`is_user`状态
	updateData := map[string]interface{}{
		"is_user":    status,
		"why":        why,
		"handled_at": time.Now(),
	}
	err = s.dao.UpdateUser(tx, id, updateData)
	if err != nil {
		return err
	}

	// 所有操作成功，提交事务
	if err := tx.Commit().Error; err != nil {
		// 提交失败也需要回滚（虽然很少见），并记录错误
		tx.Rollback()
		return fmt.Errorf("提交事务失败: %w", err)
	}
	return nil
}
