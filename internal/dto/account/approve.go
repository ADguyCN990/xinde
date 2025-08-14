package account

type ApproveReq struct {
	Why    string `json:"why" form:"why" binding:"required" example:"批准/拒绝用户申请的理由"`
	Status string `json:"status" form:"status" binding:"required,oneof=approve reject" example:"approve或者reject"`
}

type ApproveResp struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"操作成功"`
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data"`
}
