package price

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"xinde/internal/dao/common"
	model "xinde/internal/model/price"
	"xinde/internal/store"
	"xinde/pkg/stderr"
)

type Dao struct {
	db        *gorm.DB
	commonDao *common.Dao
}

func NewPriceDao() (*Dao, error) {
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

// DB 返回原始的 gorm.DB 实例，以便 Service 层可以开启事务
func (d *Dao) DB() *gorm.DB {
	return d.db
}

// CountPrices 统计价格总数
func (d *Dao) CountPrices(tx *gorm.DB) (int64, error) {
	if tx == nil {
		return 0, fmt.Errorf(stderr.ErrorDbNil)
	}

	var count int64
	err := tx.Model(model.Price{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("统计价格总数失败: " + err.Error())
	}
	return count, nil
}

func (d *Dao) FindPriceListWithPagination(tx *gorm.DB, page, pageSize int) ([]*model.Price, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}

	var list []*model.Price
	offset := (page - 1) * pageSize
	err := tx.Model(&model.Price{}).Limit(pageSize).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("分页查找价格列表失败: " + err.Error())
	}
	return list, nil
}

func (d *Dao) UpsertPrices(tx *gorm.DB, price *model.Price) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "product_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"unit", "price_1", "price_2", "price_3", "price_4"}),
	}).Create(price).Error
	
}
