// in internal/model/device/device.go
package device

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Device struct {
	ID           uint           `gorm:"primaryKey;column:id"`
	Name         string         `gorm:"type:varchar(255);column:name;not null;comment:方案名称"`
	DeviceTypeID uint           `gorm:"index;column:device_type_id;not null"`
	Details      datatypes.JSON `gorm:"type:jsonb;column:details;not null;comment:方案的动态详情(jsonb)"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index;column:deleted_at"`
}

// TableName 指定 Device 结构体对应的数据库表名
func (Device) TableName() string {
	return "t_device"
}
