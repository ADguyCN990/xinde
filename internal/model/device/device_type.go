package device

import (
	"gorm.io/gorm"
	"time"
)

type DeviceType struct {
	ID        uint           `gorm:"primaryKey;column:id"`
	Name      string         `gorm:"column:name;not null"`
	GroupID   uint           `gorm:"column:group_id;not null"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (DeviceType) TableName() string {
	return "t_device_type"
}
