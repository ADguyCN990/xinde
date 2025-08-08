package model

// User represents the t_user table in the database.
type User struct {
	// 核心字段
	UID        uint   `gorm:"primaryKey;column:uid;autoIncrement"`
	Username   string `gorm:"column:username;not null;comment:用户账号"`
	Password   string `gorm:"column:password;comment:用户密码"` // 密码通常是 string 类型，即使数据库中是 NULL
	UserMobile string `gorm:"column:usermobile;comment:用户电话号码"`
	IsAdmin    int8   `gorm:"column:is_admin;not null;default:0;comment:是否为管理员"` // tinyint -> int8
	Remarks    string `gorm:"column:remarks;comment:备注"`

	// 上次访问信息
	LastTime     int64  `gorm:"column:lasttime;comment:上次访问时间戳"` // 使用指针 *int64 来处理可为 NULL 的情况
	SearchDevice string `gorm:"column:search_device;comment:上次访问的设备"`

	// 公司信息
	ComName   string `gorm:"column:comname;not null;comment:公司名称"`
	ComArea   string `gorm:"column:comarea;comment:公司地址"`
	CompanyID uint   `gorm:"column:company_id;comment:用户对应的公司ID"` // 使用指针 *uint 来处理可为 NULL 的情况

	// 用户个人信息
	Name      string `gorm:"column:name;not null;comment:用户真实姓名"`
	UserEmail string `gorm:"column:user-email;comment:用户邮箱"` // 注意字段名中的连字符

	// 注册与审核信息
	CreatedAt int64  `gorm:"column:created_at;not null;comment:注册时间戳"`
	IsUser    int    `gorm:"column:is_user;not null;comment:是否审核通过，1为通过"`
	Why       string `gorm:"column:why;comment:审核拒绝的原因"`
	HandledAt int64  `gorm:"column:handled_at;comment:注册申请通过时间戳"`
}

// TableName specifies the table name for the User model.
func (User) TableName() string {
	return "t_user"
}
