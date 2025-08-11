package account

type RegisterReq struct {
	Name              string `form:"name" json:"name" binding:"required"`
	Username          string `form:"username" json:"username" binding:"required"`
	Password          string `form:"password" json:"password" binding:"required"`
	ConfirmedPassword string `form:"confirmed_password" json:"confirmed_password" binding:"required"`
	Phone             string `form:"phone" json:"phone" binding:"required"`
	CompanyName       string `form:"company_name" json:"company_name" binding:"required"`
	CompanyAddress    string `form:"company_address" json:"company_address"`
	Email             string `form:"email" json:"email"`
}

type RegisterResp struct {
	Msg    string `json:"msg"`
	Status bool   `json:"status"`
}
