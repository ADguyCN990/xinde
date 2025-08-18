package price

type ListReq struct {
	Page     int `json:"page" form:"page" binding:"omitempty" example:"1"`
	PageSize int `json:"page_size" form:"page_size" binding:"omitempty,min=1,max=100" example:"1-100，可选"`
}

type ListData struct {
	ID          uint    `json:"id" example:"1"`
	ProductCode string  `json:"product_code" example:"WGC001547"`
	Price1      float64 `json:"price_1" example:"1571.30"`
	Price2      float64 `json:"price_2" example:"1445.60"`
	Price3      float64 `json:"price_3" example:"1230.00"`
	Price4      float64 `json:"price_4" example:"1127.00"`
	Unit        string  `json:"unit" example:"PCS"`
	SpecCode    string  `json:"spec_code" example:"SDQCR1212H07"`
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
