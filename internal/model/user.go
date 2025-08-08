package model

// TUser undefined
type TUser struct {
	Uid `Int Not Null AutoIncrement",  "`username string `json:"uid` int NOT NULL AUTO_INCREMENT",  "`username" gorm:"uid` int NOT NULL AUTO_INCREMENT",  "`username"`
	Password string `json:"password" gorm:"password"`
	Usermobile string `json:"usermobile" gorm:"usermobile"`
	IsAdmin int8 `json:"is_admin" gorm:"is_admin"` // 是否为管理员
	Remarks string `json:"remarks" gorm:"remarks"` // 备注
	Lasttime int64 `json:"lasttime" gorm:"lasttime"` // 上次访问时间戳
	SearchDevice string `json:"search_device" gorm:"search_device"` // 上次访问的设备
	Comname string `json:"comname" gorm:"comname"`
	Comarea string `json:"comarea" gorm:"comarea"`
	Name string `json:"name" gorm:"name"` // 用户真实姓名
	User-email string `json:"user-email" gorm:"user-email"`
	CreatedAt int64 `json:"created_at" gorm:"created_at"` // 注册时间戳
	IsUser int64 `json:"is_user" gorm:"is_user"` // 是否审核通过，1为通过
	Why string `json:"why" gorm:"why"` // 审核拒绝的原因
	HandledAt int64 `json:"handled_at" gorm:"handled_at"` // 注册申请通过时间戳
	CompanyId int64 `json:"company_id" gorm:"company_id"` // 用户对应的公司ID
}

// TableName 表名称
func (*TUser) TableName() string {
	return "t_user"
}
