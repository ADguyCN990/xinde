package company

import (
	"fmt"
	"gorm.io/gorm"
	"xinde/internal/store"
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
