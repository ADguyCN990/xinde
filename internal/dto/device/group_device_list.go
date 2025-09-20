package device

type GroupDeviceListData struct {
	ID         uint   `json:"device_id" example:"1"`
	DeviceName string `json:"name" example:"刀柄"`
	ImageURL   string `json:"image_url" `
	GroupName  string `json:"group_name"`
}

type GroupDeviceListResp struct {
	Code    int                    `json:"code" example:"200"`
	Message string                 `json:"message" example:"操作成功"`
	Success bool                   `json:"success" example:"true"`
	Data    []*GroupDeviceListData `json:"data"`
}
