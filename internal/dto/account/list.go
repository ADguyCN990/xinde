package account

type ListReq struct {
	Page     int `json:"page" form:"page" binding:"required,min=1" example:"1"`
	PageSize int `json:"page_size" form:"page_size" binding:"omitempty,min=1,max=100" example:"1-100，可选"`
}

type ListData struct {
	ID             uint   `json:"id" example:"2"`
	Username       string `json:"username" example:"张三，账号名称"`
	Name           string `json:"name" example:"张三，真实名称"`
	Phone          string `json:"phone" example:"13800138000"`
	Email          string `json:"email" example:"13800138000@qq.com"`
	CompanyName    string `json:"company_name" example:"宁波鲍斯产业链服务有限公司"`
	PriceLevel     string `json:"price_level" example:"price_1"`
	Remark         string `json:"remark" example:"备注"`
	Role           string `json:"role" example:"普通用户"`
	CreatedAt      string `json:"created_at" example:"2020-09-08 09:08:09"`
	RecentSearchAt string `json:"recent_search_at" example:"2020-09-08 09:08:09"`
	SearchDevice   string `json:"search_device" example:"车削刀杆"`
}

type ListPageData struct {
	List     []*ListData `json:"list"`
	Total    int         `json:"total" example:"137"`
	Page     int         `json:"page" example:"1"`
	PageSize int         `json:"pageSize" example:"20"`
	Pages    int         `json:"pages" example:"7"`
}

type ListResp struct {
	Code    int           `json:"code" example:"200"`
	Message string        `json:"message" example:"操作成功"`
	Success bool          `json:"success" example:"true"`
	Data    *ListPageData `json:"data"`
}
