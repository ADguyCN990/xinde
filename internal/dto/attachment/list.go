package attachment

type ListReq struct {
	Page     int    `json:"page" form:"page" binding:"omitempty" example:"1"`
	PageSize int    `json:"page_size" form:"page_size" binding:"omitempty,min=1,max=100" example:"1-100，可选"`
	Filename string `json:"filename" form:"filename" binding:"omitempty" example:"价格表"`
}

type ListData struct {
	ID            uint   `json:"id" example:"1"`
	Filename      string `json:"filename" example:"价格表"`
	FileType      string `json:"file_type" example:"xxxxyyyzzz.xlsx"`
	FileSize      string `json:"file_size" example:"1.2MB"` //
	StorageDriver string `json:"storage_driver" example:"local或oss或其他"`
	BusinessType  string `json:"business_type" example:"price_import"`
	UploadedBy    string `json:"uploaded_by" example:"admin"` // 显示上传者用户名，而不是ID
	CreatedAt     string `json:"created_at"`
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
