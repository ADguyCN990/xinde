// in internal/dto/device/device.go
package device

// --- 用于 JSONB 内部结构的 DTOs ---

// ComponentDTO 对应方案中的一个组件（你提到的多个设备之一）
type ComponentDTO struct {
	Name             string  `json:"name"`
	ProductCode      string  `json:"product_code"`
	SpecCode         string  `json:"spec_code"`
	Brand            string  `json:"brand"`
	InventoryXinde   string  `json:"inventory_xinde"`   // 信德库存
	InventoryGongpin string  `json:"inventory_gongpin"` // 工品库存
	Price            float64 `json:"price"`             // 从MySQL价格表获取
	// ... 未来可以从二方API获取更多字段
}

// DetailsDTO 是整个 JSONB 字段 `details` 的 Go 结构化表示
type DetailsDTO struct {
	Filters    map[string]interface{} `json:"filters"`    // 存储所有筛选条件
	Components []*ComponentDTO        `json:"components"` // 存储组件列表
	Parameters map[string]interface{} `json:"parameters"` // 存储公共参数
}

// --- 用于数据导入的 DTO ---

// ImportDataDTO 是一个中间数据结构，用于承载从 Excel 解析后、
// 准备写入数据库的单条设备（方案）数据。
type ImportDataDTO struct {
	Name    string      `json:"name"`
	GroupID uint        `json:"group_id"`
	Details *DetailsDTO `json:"details"` // 直接使用我们定义的内部结构
}
