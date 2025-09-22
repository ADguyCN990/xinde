package device

type ChangeDeviceTypeReq struct {
	DeviceTypeID uint `json:"device_type_id" form:"device_type_id" binding:"required,min=1"`
}
