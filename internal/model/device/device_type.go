package device

import "time"

type DeviceType struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"column:name;not null"`
	GroupID   uint      `gorm:"column:group_id;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at"`
}
