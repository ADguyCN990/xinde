package device

type ImportReq struct {
	GroupID        uint   `json:"group_id" form:"group_id" binding:"required,min=1"`
	DeviceTypeName string `json:"device_type_name" form:"device_type_name" binding:"required"`
}
