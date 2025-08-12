package account

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"xinde/internal/model/account"
	"xinde/internal/store"
	"xinde/pkg/stderr"
	"xinde/pkg/util"
)

type Dao struct {
	db *gorm.DB
}

func NewRegisterDao() (*Dao, error) {
	db := store.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接未初始化，请先调用 store.InitDB()")
	}

	return &Dao{
		db: db,
	}, nil
}

// DB 返回原始的 gorm.DB 实例，以便 Service 层可以开启事务
func (d *Dao) DB() *gorm.DB {
	return d.db
}

// IsExistUser 根据username判断user是否已经存在
func (d *Dao) IsExistUser(tx *gorm.DB, name string) (bool, error) {
	if d == nil || d.db == nil || tx == nil {
		return false, fmt.Errorf(stderr.ERROR_DB_NIL)
	}

	if name == "" {
		return false, fmt.Errorf("用户名不能为空")
	}

	var user account.User
	err := tx.Where("username = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("查询用户失败: %w", err)
	}
	return true, nil
}

// CreateUser 在`t_user`表中创建用户
func (d *Dao) CreateUser(tx *gorm.DB, username, email, name, companyName, companyAddress, password, phone string, companyID uint) (uint, error) {
	if d == nil || d.db == nil || tx == nil {
		return 0, fmt.Errorf(stderr.ERROR_DB_NIL)
	}

	user := &account.User{
		Username:       username,
		Name:           name,
		UserEmail:      util.StringToPointer(email),
		CompanyID:      companyID,
		CompanyName:    companyName,
		CompanyAddress: util.StringToPointer(companyAddress),
		Password:       password,
		Phone:          phone,
	}

	if err := tx.Create(user).Error; err != nil {
		return 0, fmt.Errorf("创建用户失败: %w", err)
	}
	return user.UID, nil
}

// FindOrCreateCompany 尝试根据Name查找公司，如果没有则创建一个新的公司
func (d *Dao) FindOrCreateCompany(tx *gorm.DB, name, address string) (uint, error) {
	if d == nil || d.db == nil || tx == nil {
		return 0, fmt.Errorf(stderr.ERROR_DB_NIL)
	}

	var company account.Company
	// GORM 的 FirstOrCreate 方法完美地解决了我们的并发问题
	// 它会在一个事务中（如果 tx 不是 nil）先尝试 First，如果找不到，再 Create。
	// 底层通常会使用可串行化的隔离级别或锁来保证原子性。
	err := tx.Where(account.Company{Name: name}).
		Attrs(account.Company{Address: util.StringToPointer(address)}).
		FirstOrCreate(&company).Error

	if err != nil {
		return 0, fmt.Errorf("查找或创建公司失败: %w", err)
	}

	return company.ID, nil
}
