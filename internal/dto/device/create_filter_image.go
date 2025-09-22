package device

type CreateFilterImageReq struct {
	DeviceTypeID uint   `form:"device_type_id" json:"device_type_id" binding:"required,min=1"`
	FilterValue  string `form:"filter_value" json:"filter_value" binding:"required"`
}
