package device_access_log

import (
	"fmt"
	"gorm.io/gorm"
	"xinde/internal/dao/common"
	model "xinde/internal/model/device_access_log"
	"xinde/internal/store"
	"xinde/pkg/stderr"
)

type Dao struct {
	commonDao *common.Dao
	dao       *gorm.DB
}

func (d *Dao) DB() *gorm.DB {
	return d.dao
}

func NewDeviceAccessLogDao() (*Dao, error) {
	commonDao, err := common.NewCommonDao()
	if err != nil {
		return nil, fmt.Errorf("创建Dao层实例失败: %v", err)
	}
	dao := store.GetDB()
	return &Dao{commonDao: commonDao, dao: dao}, nil
}

func (d *Dao) Create(tx *gorm.DB, log *model.DeviceAccessLog) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}
	err := tx.Model(model.DeviceAccessLog{}).Create(log).Error
	if err != nil {
		return err
	}
	return nil
}
