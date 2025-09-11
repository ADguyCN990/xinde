package attachment

type FixOrphanReq struct {
	FilePath string `json:"file_path" binding:"required" form:"file_path" example:"孤儿文件的相对文件位置"`
	Action   string `json:"action" form:"action" binding:"required,oneof=sync delete" example:"sync表示在数据库中追加该记录，delete表示在磁盘中删除该文件"`
}
