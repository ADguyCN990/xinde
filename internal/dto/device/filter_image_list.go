package device

// FilterImageListReq 定义了列表的请求参数
type FilterImageListReq struct {
	Page         int  `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize     int  `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`
	DeviceTypeID uint `form:"device_type_id" json:"device_type_id" binding:"required,min=1"`
}

// FilterImageListData 用于列表的单项数据
type FilterImageListData struct {
	ID             uint   `json:"id"`
	DeviceTypeName string `json:"device_type_name"` // 对应设备
	FilterValue    string `json:"filter_value"`     // 匹配的设备筛选条件的名称的值
	ImageURL       string `json:"image_url"`        // 图片
}

// FilterImageListPageData 列表的分页数据
type FilterImageListPageData struct {
	List     []*FilterImageListData `json:"list"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"pageSize"`
	Pages    int                    `json:"pages"`
}

type FilterImageListResp struct {
	Code    int                      `json:"code" example:"200"`
	Message string                   `json:"message" example:"操作成功"`
	Success bool                     `json:"success" example:"true"`
	Data    *FilterImageListPageData `json:"data"`
}
