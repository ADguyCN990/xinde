package common

import (
	"fmt"
	"gorm.io/gorm"
	"xinde/internal/store"
)

type Dao struct {
	db *gorm.DB
}

func NewCommonDao() (*Dao, error) {
	db := store.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接未初始化，请先调用 store.InitDB()")
	}

	return &Dao{
		db: db,
	}, nil
}

func NewCommonPostgresDao() (*Dao, error) {
	db := store.GetPDB()
	if db == nil {
		return nil, fmt.Errorf("数据库链接未初始化，请先调用 store.InitDB()")
	}
	return &Dao{
		db: db,
	}, nil
}

// DB 返回原始的 gorm.DB 实例，以便 Service 层可以开启事务
func (d *Dao) DB() *gorm.DB {
	return d.db
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
