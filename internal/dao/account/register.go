package account

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"xinde/internal/model/account"
	"xinde/internal/store"
	"xinde/pkg/util"
)

type Dao struct {
	db *gorm.DB
}

func NewRegisterDao() (*Dao, error) {
	db := store.GetDB()
	if db == nil {
		return nil, errors.New("数据库连接未初始化，请先调用 store.InitDB()")
	}

	return &Dao{
		db: db,
	}, nil
}

// IsExistUser 根据username判断user是否已经存在
func (registerDao *Dao) IsExistUser(name string) (bool, error) {
	if registerDao == nil || registerDao.db == nil {
		return false, fmt.Errorf("Dao 或数据库连接为空")
	}

	if name == "" {
		return false, fmt.Errorf("用户名不能为空")
	}

	var user account.User
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
func (registerDao *Dao) CreateUser(username, email, name, companyName, companyAddress, password, phone string) (uint, error) {
	if registerDao == nil || registerDao.db == nil {
		return 0, fmt.Errorf("Dao 或数据库连接为空")
	}

	user := &account.User{
		Username:    username,
		Name:        name,
		UserEmail:   util.StringToPointer(email),
		CompanyName: companyName,
		CompanyArea: util.StringToPointer(companyAddress),
		Password:    password,
		Phone:       phone,
	}

	if err := registerDao.db.Create(user).Error; err != nil {
		return 0, fmt.Errorf("创建用户失败: %w", err)
	}
	return user.UID, nil
}
