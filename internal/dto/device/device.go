// in internal/dto/device/device.go
package device

// --- 用于 JSONB 内部结构的 DTOs ---

// ImportComponentDTO 对应 Excel 中解析出的一个组件，只包含原始数据
type ImportComponentDTO struct {
	Name        string `json:"name"`
	ProductCode string `json:"product_code"`
	SpecCode    string `json:"spec_code"`
}

// ImportDetailsDTO 是导入时，JSONB 字段 `details` 的 Go 结构化表示
type ImportDetailsDTO struct {
	Filters    map[string]interface{} `json:"filters"`
	Components []*ImportComponentDTO  `json:"components"`
	Parameters map[string]interface{} `json:"parameters"`
}

// ImportDataDTO 是一个中间数据结构，承载从 Excel 解析后、
// 准备写入数据库的单条设备（方案）数据。
type ImportDataDTO struct {
	Name    string            `json:"name"`
	GroupID uint              `json:"group_id"`
	Details *ImportDetailsDTO `json:"details"`
}

// --- DTOs for API Read/Response (将在未来开发选型接口时使用) ---
/*
// 这是我们未来会用到的“读取模型”DTO，现在先注释掉
type ComponentDTO struct {
	Name             string  `json:"name"`
	ProductCode      string  `json:"product_code"`
	SpecCode         string  `json:"spec_code"`
	Brand            string  `json:"brand"`
	InventoryXinde   string  `json:"inventory_xinde"`
	InventoryGongpin string  `json:"inventory_gongpin"`
	Price            float64 `json:"price"`
	// ... from 2nd-party API
}
// ... etc.
*/
