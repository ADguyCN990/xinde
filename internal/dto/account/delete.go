package account

type DeleteResp struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"操作成功"`
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data"`
}
