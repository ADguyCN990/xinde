package account

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		return false, fmt.Errorf(stderr.ErrorDbNil)
	}

	if name == "" {
		return false, fmt.Errorf("用户名不能为空")
	}

	var user account.User
	err := tx.Where("username = ?", name).First(&user).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// IsExistUserByID 根据ID判断用户是否存在
func (d *Dao) IsExistUserByID(tx *gorm.DB, uid uint) (bool, error) {
	if tx == nil {
		return false, fmt.Errorf(stderr.ErrorDbNil)
	}
	var count int64
	// GORM 的查询会自动处理 `deleted_at IS NULL`
	err := tx.Model(&account.User{}).Where("uid = ?", uid).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *Dao) GetUserByID(tx *gorm.DB, uid uint) (*account.User, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}

	var user *account.User
	err := tx.Model(&account.User{}).Where("uid = ?", uid).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("根据id查找user失败: " + err.Error())
	}
	return user, nil
}

// GetUserByIDForUpdate 带行级锁，根据ID查找用户
func (d *Dao) GetUserByIDForUpdate(tx *gorm.DB, uid uint) (*account.User, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}

	var user account.User
	// 行级锁，防止其他管理员同时审批这个用户
	err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "NOWAIT"}).Where("uid = ?", uid).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindUserByUsername 根据username查找用户
func (d *Dao) FindUserByUsername(tx *gorm.DB, username string) (*account.User, error) {
	var user account.User
	err := tx.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(stderr.ErrorUserUnauthorized)
		} else {
			return nil, fmt.Errorf("FindUserByUsername查询用户失败: %w", err)
		}
	}
	return &user, nil
}

// CountUserWithStatus 查询【未审批/已通过/以拒绝】用户的个数
func (d *Dao) CountUserWithStatus(tx *gorm.DB, status int) (int64, error) {
	if tx == nil {
		return 0, fmt.Errorf(stderr.ErrorDbNil)
	}
	var count int64
	err := tx.Model(&account.User{}).Where("is_user = ?", status).Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, fmt.Errorf("统计用户总数失败: %w", err)
	}
	return count, nil
}

// FindUserListWithPagination 分页查找用户列表
func (d *Dao) FindUserListWithPagination(tx *gorm.DB, page, pageSize, status int) ([]*account.User, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}

	var users []*account.User

	// 执行查询
	offset := (page - 1) * pageSize
	err := tx.
		Model(&account.User{}).
		Order("t_user.uid asc").
		Select("t_user.*, t_company.price_level").
		Joins("LEFT JOIN t_company ON t_user.company_id = t_company.id").
		Where("t_user.is_user = ?", status).
		Limit(pageSize).
		Offset(offset).
		Find(&users).
		Error
	if err != nil {
		return nil, fmt.Errorf("统计用户列表失败: %w", err)
	}

	return users, nil
}

// CreateUser 在`t_user`表中创建用户
func (d *Dao) CreateUser(tx *gorm.DB, username, email, name, companyName, companyAddress, password, phone string, companyID uint) (uint, error) {
	if d == nil || d.db == nil || tx == nil {
		return 0, fmt.Errorf(stderr.ErrorDbNil)
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

// UpdateUser 更新用户
func (d *Dao) UpdateUser(tx *gorm.DB, uid uint, updateData map[string]interface{}) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}
	err := tx.Model(account.User{}).Where("uid = ?", uid).Updates(updateData).Error
	if err != nil {
		return err
	}
	return nil
}

// FindOrCreateCompany 尝试根据Name查找公司，如果没有则创建一个新的公司
func (d *Dao) FindOrCreateCompany(tx *gorm.DB, name, address string) (uint, error) {
	if d == nil || d.db == nil || tx == nil {
		return 0, fmt.Errorf(stderr.ErrorDbNil)
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

// DeleteUserByID 根据ID删除用户
func (d *Dao) DeleteUserByID(tx *gorm.DB, uid uint) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}
	var user account.User
	err := tx.Where("uid = ?", uid).Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// Transaction 添加事务包装方法
func (d *Dao) Transaction(fn func(*gorm.DB) error) error {
	tx := d.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("开启事务失败: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // 重新抛出 panic
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// UserPrice is a temporary struct to hold the result of the price query.
type UserPrice struct {
	ProductCode string  `gorm:"column:product_code"`
	Price1      float64 `gorm:"column:price_1"`
	Price2      float64 `gorm:"column:price_2"`
	Price3      float64 `gorm:"column:price_3"`
	Price4      float64 `gorm:"column:price_4"`
	PriceLevel  string  `gorm:"column:price_level"`
}

// FindPricesForUser retrieves prices for a list of product codes based on a user's price level.
func (d *Dao) FindPricesForUser(tx *gorm.DB, uid uint, productCodes []string) ([]*UserPrice, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	if len(productCodes) == 0 {
		return []*UserPrice{}, nil
	}

	var prices []*UserPrice

	// GORM 无法完美构建 CROSS JOIN，所以我们使用原生 SQL 以确保查询正确
	sql := `
        SELECT
            p.product_code,
            p.price_1,
            p.price_2,
            p.price_3,
            p.price_4,
            IFNULL(c.price_level, 'price_1') AS price_level 
        FROM
            t_user u
        LEFT JOIN
            t_company c ON u.company_id = c.id
		CROSS JOIN
            t_price p
        WHERE
	    p.product_code IN (?)
            AND u.uid = ?;
    `

	err := tx.Raw(sql, productCodes, uid).Scan(&prices).Error
	if err != nil {
		return nil, err
	}

	return prices, nil
}
