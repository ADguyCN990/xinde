package group

type TreeReq struct {
	Icon string `json:"icon" form:"icon" binding:"required,oneof=true false" example:"true表示树状列表返回图片，false表示树状列表不返回图片"`
}

// GroupTreeNode 用于树状结构（前台展示/父级选择）的节点
type GroupTreeNode struct {
	ID       uint             `json:"id"`
	Name     string           `json:"name"`
	ParentID uint             `json:"parent_id"`
	IconURL  string           `json:"icon_url,omitempty"` // omitempty 可以在不需要时隐藏
	Children []*GroupTreeNode `json:"children,omitempty"`
}

// TreeResp 树状结构的完整响应
type TreeResp struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Success bool             `json:"success"`
	Data    []*GroupTreeNode `json:"data"` // 保持为列表，即使只返回 root 节点
}
