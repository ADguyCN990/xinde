package device

// in internal/dto/device/device.go

// ... (之前的 DTOs)

// ChangeGroupReq 定义了更换分组时的请求体
type ChangeGroupReq struct {
	GroupID uint `json:"new_group_id" form:"group_id" binding:"required,min=1" example:"1"`
}
