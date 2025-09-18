package device

type ChangeNameReq struct {
	Name string `json:"name" form:"name" binding:"required" example:"车削刀杆"`
}
