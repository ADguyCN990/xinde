package account

type LoginReq struct {
	Username string `form:"username" json:"username" binding:"required" example:"金晖"`
	Password string `form:"password" json:"password" binding:"required" example:"923845797582"`
}

type LoginData struct {
	Username    string `json:"username" example:"金晖"`
	Name        string `json:"name" example:"金晖"`
	Phone       string `json:"phone" example:"13065859690"`
	Email       string `json:"email,omitempty" example:"1921771473@qq.com"`
	AccessToken string `json:"access_token" example:"jwt_token"`
}

type LoginResp struct {
	Code    int        `json:"code" example:"200"`
	Message string     `json:"message" example:"操作成功"`
	Success bool       `json:"success" example:"true"`
	Data    *LoginData `json:"data"`
}
