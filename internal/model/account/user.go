package account

import (
	"gorm.io/gorm"
	"time"
)

// User represents the t_user table in the database.
type User struct {
	// 核心字段
	UID       uint           `gorm:"primaryKey;column:uid;autoIncrement"`
	Username  string         `gorm:"column:username;not null;comment:用户账号"`
	Password  string         `gorm:"column:password;comment:用户密码"` // 密码通常是 string 类型，即使数据库中是 NULL
	Phone     string         `gorm:"column:phone;comment:用户电话号码"`
	IsAdmin   int8           `gorm:"column:is_admin;not null;default:0;comment:是否为管理员"` // tinyint -> int8
	Remarks   *string        `gorm:"column:remarks;comment:备注"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"` // GORM 软删除

	// 上次访问信息
	RecentSearchAt *time.Time `gorm:"column:recent_search_at;comment:上次访问时间"` //
	SearchDevice   *string    `gorm:"column:search_device;comment:上次访问的设备"`

	// 公司信息
	CompanyName    string  `gorm:"column:company_name;not null;comment:公司名称"`
	CompanyAddress *string `gorm:"column:company_address;comment:公司地址"`
	CompanyID      uint    `gorm:"column:company_id;comment:用户对应的公司ID"` // 使用指针 *uint 来处理可为 NULL 的情况

	// 用户个人信息
	Name      string  `gorm:"column:name;not null;comment:用户真实姓名"`
	UserEmail *string `gorm:"column:user_email;comment:用户邮箱"` // 注意字段名中的连字符

	// 注册与审核信息
	UpdatedAt time.Time  `gorm:"column:updated_at;not null"`
	HandledAt *time.Time `gorm:"column:handled_at;comment:注册申请处理时间"` // 指针处理 NULL
	IsUser    int        `gorm:"column:is_user;not null;comment:是否审核通过,0为未处理, 1为通过, 2为拒绝"`
	Why       *string    `gorm:"column:why;comment:审核拒绝的原因"`
}

// TableName specifies the table name for the User model.
func (User) TableName() string {
	return "t_user"
}
