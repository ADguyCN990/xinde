package group

import (
	"gorm.io/gorm"
	"time"
	"xinde/internal/model/attachment"
)

type Group struct {
	ID        uint                   `gorm:"primaryKey;column:id"`
	Name      string                 `gorm:"type:varchar(255);column:name;not null;comment:分组名称"`
	ParentID  uint                   `gorm:"not null;column:parent_id;comment:父级分组ID"`
	Icon      *attachment.Attachment `gorm:"-"` // 不在数据库中创建字段，仅用于业务逻辑
	CreatedAt time.Time              `gorm:"column:created_at"`
	UpdatedAt time.Time              `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt         `gorm:"index;column:deleted_at"`
}
