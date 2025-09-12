package attachment

import (
	"gorm.io/gorm"
	"time"
)

type Attachment struct {
	ID            uint   `gorm:"primaryKey;column:id"`
	Filename      string `gorm:"type:varchar(255);column:filename;not null"`
	StoragePath   string `gorm:"type:varchar(512);column:storage_path;not null"`
	FileType      string `gorm:"type:varchar(100);column:file_type;not null"`
	FileSize      uint64 `gorm:"not null;column:file_size"`
	StorageDriver string `gorm:"type:varchar(20);column:storage_driver;not null;default:local"`
	UploadedByUID uint   `gorm:"not null;column:uploaded_by_uid;index"`
	// 使用ReadOnly标签
	UploaderName string         `gorm:"->"`
	BusinessType *string        `gorm:"type:varchar(50);column:business_type;index"`
	BusinessID   uint           `gorm:"index;column:business_id;comment:业务ID"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Attachment) TableName() string {
	return "t_attachment"
}
