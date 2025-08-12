package account

import (
	"gorm.io/gorm"
	"time"
)

// Company represents the t_company table in the database.
// It maps to the company information table.
type Company struct {
	// 主键
	ID uint `gorm:"primaryKey;column:id;autoIncrement;comment:主键默认id"`

	// 核心信息
	Name       string `gorm:"column:name;not null;comment:公司名称"`
	PriceLevel string `gorm:"column:price_level;not null;default:price_1;comment:该公司查看产品的价格等级"`

	// 可为空的字段，使用指针类型
	Address *string `gorm:"column:address;comment:公司地址"`

	// 自动管理的时间戳与软删除字段
	CreatedAt time.Time      `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index:idx_deleted_at"` // GORM 软删除字段，会自动处理索引
}

// TableName explicitly sets the table name for the Company model.
func (Company) TableName() string {
	return "t_company"
}
