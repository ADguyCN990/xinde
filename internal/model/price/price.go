package price

import (
	"gorm.io/gorm"
	"time"
)

// Price represents the t_price table in the database.
type Price struct {
	ID          uint   `gorm:"primaryKey;column:id;autoIncrement"`
	ProductCode string `gorm:"column:product_code;unique;not null"`
	Unit        string `gorm:"column:unit;not null"`
	SpecCode    string `gorm:"column:spec_code;not null"`

	// 在 Go 中，decimal 通常用 string 或专门的 decimal 类型来精确表示
	// 使用 float64 可能会有精度问题，但对于业务计算通常也够用。
	// 如果需要高精度计算，推荐使用 shopspring/decimal
	Price1 float64 `gorm:"column:price_1;type:decimal(10,2);not null"`
	Price2 float64 `gorm:"column:price_2;type:decimal(10,2);not null"`
	Price3 float64 `gorm:"column:price_3;type:decimal(10,2);not null"`
	Price4 float64 `gorm:"column:price_4;type:decimal(10,2);not null"`

	// 标准时间戳与软删除字段
	CreatedAt time.Time      `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName explicitly sets the table name.
func (Price) TableName() string {
	return "t_price"
}
