package device

import (
	"fmt"
	"gorm.io/gorm"
	"xinde/internal/dao/common"
	model "xinde/internal/model/device"
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

func NewDeviceDao() (dao *Dao, err error) {
	db := store.GetPDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接未初始化，请先调用 store.InitDB()")
	}
	commonDao, err := common.NewCommonPostgresDao()
	if err != nil {
		return nil, err
	}
	return &Dao{
		db:        db,
		commonDao: commonDao,
	}, nil
}

func (d *Dao) CreateDevice(tx *gorm.DB, dev *model.Device) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}
	if err := tx.Model(&model.Device{}).Create(dev).Error; err != nil {
		return fmt.Errorf("Dao层创建设备方案失败: " + err.Error())
	}
	return nil
}

func (d *Dao) BatchCreateDevice(tx *gorm.DB, devs []*model.Device) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}
	if len(devs) == 0 {
		return nil
	}
	if err := tx.Model(&model.Device{}).Create(devs).Error; err != nil {
		return fmt.Errorf("Dao层批量创建设备方案失败: " + err.Error())
	}
	return nil
}

// FindOrCreateDeviceType 查找或创建一个设备类型
func (d *Dao) FindOrCreateDeviceType(tx *gorm.DB, name string, groupID uint) (*model.DeviceType, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	var dt model.DeviceType
	// FirstOrCreate 会查找，如果找不到，就用给定的结构体创建
	if err := tx.Where("name = ? AND group_id = ?", name, groupID).FirstOrCreate(&dt, model.DeviceType{
		Name:    name,
		GroupID: groupID,
	}).Error; err != nil {
		return nil, fmt.Errorf("查找或创建设备类型失败: " + err.Error())
	}
	return &dt, nil
}

func (d *Dao) DeleteByDeviceTypeID(tx *gorm.DB, deviceTypeID uint) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}
	err := tx.Delete(&model.Device{}, "device_type_id = ?", deviceTypeID).Error
	if err != nil {
		return fmt.Errorf("Dao层根据DeviceTypeID删除设备失败: " + err.Error())
	}
	return nil
}
