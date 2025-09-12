package group

type CreateReq struct {
	Name     string `form:"name" json:"name" binding:"required,max=100"`
	ParentID uint   `form:"parent_id" json:"parent_id" binding:"required,min=1"` // 父ID必须大于等于1 (root的ID)
}
