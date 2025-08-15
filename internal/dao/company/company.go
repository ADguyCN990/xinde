package company

import (
	"fmt"
	"gorm.io/gorm"
	model "xinde/internal/model/account"
	"xinde/internal/store"
	"xinde/pkg/stderr"
)

type Dao struct {
	db *gorm.DB
}

func NewCompanyDao() (*Dao, error) {
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

// CountCompanies 统计公司总数
func (d *Dao) CountCompanies(tx *gorm.DB) (int64, error) {
	if tx == nil {
		return 0, fmt.Errorf(stderr.ErrorDbNil)
	}

	var count int64
	err := tx.Model(&model.Company{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("统计公司总数失败: " + err.Error())
	}
	return count, nil
}

// FindCompanyListWithPagination 分页查找公司列表
func (d *Dao) FindCompanyListWithPagination(tx *gorm.DB, page, pageSize int) ([]*model.Company, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}

	var companies []*model.Company
	offset := (page - 1) * pageSize
	err := tx.Model(&model.Company{}).
		Limit(pageSize).Offset(offset).
		Find(&companies).Error
	if err != nil {
		return nil, fmt.Errorf("查找公司列表失败: " + err.Error())
	}
	return companies, nil
}
