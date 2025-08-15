package company

type ListReq struct {
	Page     int `json:"page" form:"page" binding:"omitempty" example:"1"`
	PageSize int `json:"page_size" form:"page_size" binding:"omitempty,min=1,max=100" example:"1-100，可选"`
}

type ListData struct {
	ID         uint   `json:"id" example:"1"`
	Name       string `json:"name" example:"宁波鲍斯产业链有限公司"`
	Address    string `json:"address" example:"浙江省宁波市奉化区江口街道聚潮路55号"`
	PriceLevel string `json:"price_level" example:"折扣等级1 折扣等级2 折扣等级3 折扣等级4"`
	CreatedAt  string `json:"created_at" example:"2021-09-09 09:09:09"`
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
