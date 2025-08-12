package account

type RegisterReq struct {
	Name              string `form:"name" json:"name" binding:"required" example:"金晖"`
	Username          string `form:"username" json:"username" binding:"required" example:"金晖"`
	Password          string `form:"password" json:"password" binding:"required" example:"923845797582"`
	ConfirmedPassword string `form:"confirmed_password" json:"confirmed_password" binding:"required" example:"923845797582"`
	Phone             string `form:"phone" json:"phone" binding:"required" example:"13065859690"`
	CompanyName       string `form:"company_name" json:"company_name" binding:"required" example:"宁波鲍斯产业链服务有限公司"`
	CompanyAddress    string `form:"company_address" json:"company_address,omitempty" example:"浙江省宁波市奉化区江口街道聚潮路55号"`
	Email             string `form:"email" json:"email,omitempty" binding:"email" example:"1921771473@qq.com"`
}

type RegisterResp struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"操作成功"`
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data"`
}
