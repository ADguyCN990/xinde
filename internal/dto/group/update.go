package group

type UpdateReq struct {
	ParentID uint   `json:"parent_id" form:"parent_id" binding:"omitempty,min=1"`
	Name     string `json:"name" form:"name" binding:"omitempty"`
}
