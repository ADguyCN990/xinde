package group

import (
	"fmt"
	"gorm.io/gorm"
	"xinde/internal/dao/common"
	model "xinde/internal/model/group"
	"xinde/internal/store"
	"xinde/pkg/stderr"
)

type Dao struct {
	db        *gorm.DB
	commonDao *common.Dao
}

func (d *Dao) DB() *gorm.DB {
	return d.db
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

func (d *Dao) Create(tx *gorm.DB, groupName string, parentID uint) (uint, error) {
	if tx == nil {
		return 0, fmt.Errorf(stderr.ErrorDbNil)
	}

	g := &model.Group{
		Name:     groupName,
		ParentID: parentID,
	}
	if err := tx.Model(&model.Group{}).Create(g).Error; err != nil {
		return 0, fmt.Errorf("Dao层创建分组失败: " + err.Error())
	}

	return g.ID, nil
}

func (d *Dao) GetAll(tx *gorm.DB) ([]*model.Group, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	var groups []*model.Group
	if err := tx.Model(&model.Group{}).Order("id asc").Find(&groups).Error; err != nil {
		return nil, fmt.Errorf("Dao层查找所有分组失败: " + err.Error())
	}
	return groups, nil
}
