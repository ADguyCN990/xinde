package account

type UpdatePasswordReq struct {
	Password string `json:"password" form:"password" binding:"required" example:"这是一个密码"`
}
