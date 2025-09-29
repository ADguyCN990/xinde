package device_access_log

import "time"

type DeviceAccessLog struct {
	ID           uint      `gorm:"primaryKey;column:id"`
	UserID       uint      `gorm:"index;column:user_id;not null;comment:访问的用户ID"`
	CompanyID    uint      `gorm:"index;column:company_id;comment:用户所属公司ID"`
	DeviceTypeID uint      `gorm:"index;column:device_type_id;not null;comment:访问的设备类型ID"`
	AccessedAt   time.Time `gorm:"column:accessed_at;not null;comment:访问时间"`
}

func (DeviceAccessLog) TableName() string {
	return "t_device_access_log"
}
