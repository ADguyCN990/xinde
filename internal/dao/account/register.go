package account

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"xinde/internal/model"
	"xinde/internal/store"
)

type RegisterDao struct {
	db *gorm.DB
}

func NewRegisterDao() (*RegisterDao, error) {
	db := store.GetDB()
	if db == nil {
		return nil, errors.New("数据库连接未初始化，请先调用 store.InitDB()")
	}

	return &RegisterDao{
		db: db,
	}, nil
}

// IsExistUser 根据username判断user是否已经存在
func (registerDao *RegisterDao) IsExistUser(name string) (bool, error) {
	if registerDao == nil || registerDao.db == nil {
		return false, fmt.Errorf("RegisterDao 或数据库连接为空")
	}

	if name == "" {
		return false, fmt.Errorf("用户名不能为空")
	}

	var user model.User
	err := registerDao.db.Where("username = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("查询用户失败: %w", err)
	}
	return true, nil
}

// CreateUser 在`t_user`表中创建用户
func (registerDao *RegisterDao) CreateUser(username, email, name, companyName, companyAddress, password, phone string) (uint, error) {
	if registerDao == nil || registerDao.db == nil {
		return 0, fmt.Errorf("RegisterDao 或数据库连接为空")
	}

	user := &model.User{
		Username:    username,
		Name:        name,
		UserEmail:   email,
		CompanyName: companyName,
		CompanyArea: companyAddress,
		Password:    password,
		Phone:       phone,
	}

	if err := registerDao.db.Create(user).Error; err != nil {
		return 0, fmt.Errorf("创建用户失败: %w", err)
	}
	return user.UID, nil
}
