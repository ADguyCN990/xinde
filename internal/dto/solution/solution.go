package solution

// PaginationReq (与之前相同)
type PaginationReq struct {
	Page     int `json:"page" form:"page" binding:"required,min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"required,min=1,max=100"`
}

// QueryReq (与之前相同)
type QueryReq struct {
	DeviceTypeID   uint                   `json:"device_type_id" form:"device_type_id" binding:"required,min=1"`
	CurrentFilters map[string]interface{} `json:"current_filters" form:"current_filters"`
	Pagination     PaginationReq          `json:"pagination"`
}

// --- Response DTOs ---

// ComponentData 对应方案中的一个组件，是【读取模型】，包含了聚合后的所有数据
type ComponentData struct {
	Name             string  `json:"name"`
	ProductCode      string  `json:"product_code"`
	SpecCode         string  `json:"spec_code"`
	Brand            string  `json:"brand,omitempty"`   // 来自 API
	ImageURL         string  `json:"image_url"`         // 【新增】来自 API
	InventoryXinde   string  `json:"inventory_xinde"`   // 来自 API (onhand)
	InventoryGongpin string  `json:"inventory_gongpin"` // 来自 API (bsonhand)
	Price            float64 `json:"price"`             // 来自 MySQL 价格表

}

// DetailsData 是 JSONB 字段 `details` 的【读取模型】表示
type DetailsData struct {
	Filters    map[string]interface{} `json:"filters"`
	Components []*ComponentData       `json:"components"`
	Parameters map[string]interface{} `json:"parameters"`
}

// SolutionData 代表一条返回给前端的、聚合了所有数据的方案
type SolutionData struct {
	ID      uint         `json:"id"`
	Name    string       `json:"name"`
	Details *DetailsData `json:"details"`
}

// FilterOption 代表一个可用的筛选选项
type FilterOption struct {
	Value    interface{} `json:"value"`
	ImageURL string      `json:"image_url,omitempty"`
}

// AvailableFilter 代表一个可用的筛选条件及其所有选项
type AvailableFilter struct {
	FilterName string         `json:"filter_name"`
	Options    []FilterOption `json:"options"`
}

// SolutionsPageData 对应你提供的 ListPageData 格式
type SolutionsPageData struct {
	List     []*SolutionData `json:"list"`
	Total    int64           `json:"total"` // DAO层返回的是int64，这里保持一致
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
	Pages    int64           `json:"pages"`
}

// QueryResp 是查询接口的完整响应体
type QueryResp struct {
	Solutions        *SolutionsPageData `json:"solutions"`
	AvailableFilters []*AvailableFilter `json:"available_filters"`
}
