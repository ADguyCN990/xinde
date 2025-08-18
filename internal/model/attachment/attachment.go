package attachment

import (
	"gorm.io/gorm"
	"time"
)

type Attachment struct {
	ID            uint   `gorm:"primaryKey"`
	Filename      string `gorm:"type:varchar(255);not null"`
	StoragePath   string `gorm:"type:varchar(512);not null"`
	FileType      string `gorm:"type:varchar(100);not null"`
	FileSize      uint64 `gorm:"not null"`
	StorageDriver string `gorm:"type:varchar(20);not null;default:local"`
	UploadedByUID uint   `gorm:"not null;index"`
	// 使用ReadOnly标签
	UploaderName string  `gorm:"->"`
	BusinessType *string `gorm:"type:varchar(50);index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (Attachment) TableName() string {
	return "t_attachment"
}
