package group

import (
	"fmt"
	"gorm.io/gorm"
	"xinde/internal/dao/common"
	"xinde/internal/store"
)

type Dao struct {
	db        *gorm.DB
	commonDao *common.Dao
}

func NewGroupDao() (*Dao, error) {
	db := store.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接未初始化，请先调用 store.InitDB()")
	}
	commonDao, err := common.NewCommonDao()
	if err != nil {
		return nil, err
	}
	return &Dao{
		db:        db,
		commonDao: commonDao,
	}, nil
}
