package account

type ApprovalListReq struct {
	Page     int    `json:"page" form:"page" binding:"omitempty" example:"1"`
	PageSize int    `json:"page_size" form:"page_size" binding:"omitempty,min=1,max=100" example:"1-100，可选"`
	Status   string `json:"status" form:"status" binding:"omitempty" example:"pending或approved或rejected"` // oneof 校验参数必须是其中之一
}

type ApprovalListData struct {
	ID          uint   `json:"id" example:"2"`
	Username    string `json:"username" example:"张三，账号名称"`
	Name        string `json:"name" example:"张三，真实名称"`
	Phone       string `json:"phone" example:"13800138000"`
	Email       string `json:"email" example:"13800138000@qq.com"`
	CompanyName string `json:"company_name" example:"宁波鲍斯产业链服务有限公司"`
	CreatedAt   string `json:"created_at" example:"2020-09-08 09:08:09"`
	Why         string `json:"why" example:"同意该用户的申请；拒绝该用户的申请"`
}

type ApprovalListPageData struct {
	List     []*ApprovalListData `json:"list"`
	Total    int                 `json:"total" example:"137"`
	Page     int                 `json:"page" example:"1"`
	PageSize int                 `json:"pageSize" example:"20"`
	Pages    int                 `json:"pages" example:"7"`
}

type ApprovalListResp struct {
	Code    int                   `json:"code" example:"200"`
	Message string                `json:"message" example:"操作成功"`
	Success bool                  `json:"success" example:"true"`
	Data    *ApprovalListPageData `json:"data"`
}
