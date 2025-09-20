package device

import (
	"errors"
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

// CountDeviceTypes 查找DeviceType数量
func (d *Dao) CountDeviceTypes(tx *gorm.DB) (int64, error) {
	if tx == nil {
		return 0, fmt.Errorf(stderr.ErrorDbNil)
	}
	var count int64
	err := tx.Model(&model.DeviceType{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("Dao层查找DeviceType数量失败: " + err.Error())
	}
	return count, nil
}

// CountSolutionsWithDeviceTypeID 根据DeviceTypeID查找设备方案数量
func (d *Dao) CountSolutionsWithDeviceTypeID(tx *gorm.DB, deviceTypeID string) (int64, error) {
	if tx == nil {
		return 0, fmt.Errorf(stderr.ErrorDbNil)
	}
	var count int64
	queryBuilder := tx.Model(&model.Device{}).Where("t_device.device_type_id = ?", deviceTypeID)
	err := queryBuilder.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("Dao层根据deviceTypeID查找设备方案数量失败: " + err.Error())
	}
	return count, nil
}

// FindOrCreateDeviceType 查找或创建一个设备类型
func (d *Dao) FindOrCreateDeviceType(tx *gorm.DB, name string, groupID uint) (*model.DeviceType, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	var dt *model.DeviceType
	// FirstOrCreate 会查找，如果找不到，就用给定的结构体创建
	if err := tx.Where("name = ? AND group_id = ?", name, groupID).FirstOrCreate(&dt, model.DeviceType{
		Name:    name,
		GroupID: groupID,
	}).Error; err != nil {
		return nil, fmt.Errorf("查找或创建设备类型失败: " + err.Error())
	}
	return dt, nil
}

func (d *Dao) GetDeviceTypeByID(tx *gorm.DB, id uint) (*model.DeviceType, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	var dt *model.DeviceType
	if err := tx.Model(&model.DeviceType{}).Where("id = ?", id).First(&dt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		} else {
			return nil, fmt.Errorf("Dao层根据ID查找设备类型失败" + err.Error())
		}
	}
	return dt, nil
}

func (d *Dao) GetDeviceTypesByGroupID(tx *gorm.DB, groupID uint) ([]*model.DeviceType, error) {
	if tx == nil {
		return nil, fmt.Errorf(stderr.ErrorDbNil)
	}
	var dts []*model.DeviceType
	if err := tx.Model(&model.DeviceType{}).Where("group_id = ?", groupID).Find(&dts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		} else {
			return nil, fmt.Errorf("Dao层根据GroupID查找设备类型失败: " + err.Error())
		}
	}
	return dts, nil
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

func (d *Dao) DeleteDeviceTypeByID(tx *gorm.DB, id uint) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}
	err := tx.Delete(&model.DeviceType{}, "id = ?", id).Error
	if err != nil {
		return fmt.Errorf("Dao层根据DeviceTypeID删除设备失败: " + err.Error())
	}
	return nil
}

// RawDeviceType is a temporary struct to hold the result of the complex query.
type RawDeviceType struct {
	model.DeviceType
	SolutionCount int64 `gorm:"column:solution_count"`
}

func (d *Dao) GetDeviceTypeListPage(tx *gorm.DB, page, pageSize int) (int64, []*RawDeviceType, error) {
	var total int64
	var list []*RawDeviceType
	// 使用子查询来计算每个 device_type 的 solution 数量
	query := tx.Model(&model.DeviceType{}).
		Select("t_device_type.*, (SELECT count(*) FROM t_device WHERE t_device.device_type_id = t_device_type.id AND t_device.deleted_at IS NULL) as solution_count")

	// 1. 先计算总数
	if err := query.Count(&total).Error; err != nil {
		return 0, nil, fmt.Errorf("分页查询DeviceType失败: " + err.Error())
	}

	// 2. 再获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Order("id desc").Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
		return 0, nil, fmt.Errorf("分页查询DeviceType失败: " + err.Error())
	}

	return total, list, nil
}

func (d *Dao) UpdateDeviceType(tx *gorm.DB, id uint, updateData map[string]interface{}) error {
	if tx == nil {
		return fmt.Errorf(stderr.ErrorDbNil)
	}
	if err := tx.Model(&model.DeviceType{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		return fmt.Errorf("更新分组类型字段失败: " + err.Error())
	}
	return nil
}
