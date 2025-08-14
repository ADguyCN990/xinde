package account

type ResetRemarkReq struct {
	Remark string `json:"remark" form:"remark" binding:"omitempty" example:"这是一条备注"`
}
