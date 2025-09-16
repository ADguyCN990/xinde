package device

type Import struct {
	GroupID uint `json:"group_id" form:"group_id" binding:"required,min=1"`
}
