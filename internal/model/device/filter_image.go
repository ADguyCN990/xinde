package device

import (
	"gorm.io/gorm"
	"time"
	"xinde/internal/model/attachment"
)

type FilterImage struct {
	ID           uint                   `gorm:"primaryKey;column:id"`
	DeviceTypeID uint                   `gorm:"index;column:device_type_id;not null;comment:关联的设备类型ID"`
	FilterValue  string                 `gorm:"type:varchar(255);column:filter_value;not null;comment:筛选条件的值"`
	Image        *attachment.Attachment `gorm:"-"`
	CreatedAt    time.Time              `gorm:"column:created_at"`
	UpdatedAt    time.Time              `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt         `gorm:"index;column:deleted_at"` // 增加软删除
}

func (FilterImage) TableName() string {
	return "t_filter_image"
}
