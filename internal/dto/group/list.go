package group

type ListReq struct {
	Page     int `json:"page" form:"page" binding:"omitempty" example:"1"`
	PageSize int `json:"page_size" form:"page_size" binding:"omitempty,min=1,max=100" example:"1-100，可选"`
}

// BackendListData 用于后台扁平化列表的单项数据
type BackendListData struct {
	ID         uint   `json:"id" example:"2"`
	Name       string `json:"name" example:"钢件加工"`
	ParentID   uint   `json:"parent_id" example:"1"`
	ParentName string `json:"parent_name,omitempty" example:"root"`
	Level      int    `json:"level" example:"1"` // 【新增】层级字段
	IconURL    string `json:"icon_url,omitempty" example:"图片url"`
	CreatedAt  string `json:"created_at" example:"2020-09-08 09:09:09"`
}

type ListPageData struct {
	List     []*BackendListData `json:"list"`
	Total    int                `json:"total" example:"137"`
	Page     int                `json:"page" example:"1"`
	PageSize int                `json:"pageSize" example:"20"`
	Pages    int                `json:"pages" example:"7"`
}

type ListResp struct {
	Code    int           `json:"code" example:"200"`
	Message string        `json:"message" example:"操作成功"`
	Success bool          `json:"success" example:"true"`
	Data    *ListPageData `json:"data"`
}
