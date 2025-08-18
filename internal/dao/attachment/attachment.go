package attachment

import (
	"fmt"
	"gorm.io/gorm"
	"xinde/internal/dao/common"
	model "xinde/internal/model/attachment"
	"xinde/internal/store"
	"xinde/pkg/stderr"
)

type Dao struct {
	db        *gorm.DB
	commonDao *common.Dao
}

func NewAttachmentDao() (*Dao, error) {
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

func (d *Dao) Create(tx *gorm.DB, att *model.Attachment) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}

	err := tx.Model(&model.Attachment{}).Create(att).Error
	if err != nil {
		return fmt.Errorf("创建attachment失败: " + err.Error())
	}
	return nil
}
