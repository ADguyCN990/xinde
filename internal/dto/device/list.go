package device

type ListReq struct {
	Page     int `json:"page" form:"page" binding:"omitempty" example:"1"`
	PageSize int `json:"page_size" form:"page_size" binding:"omitempty,min=1,max=100" example:"1-100，可选"`
}

type ListData struct {
	ID            uint   `json:"id"`
	GroupName     string `json:"group_name"`     // 分组名称
	Name          string `json:"name"`           // 设备类型名称
	ImageURL      string `json:"image_url"`      // 设备图片
	SolutionCount int64  `json:"solution_count"` // 条目数 (方案数量)
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
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
