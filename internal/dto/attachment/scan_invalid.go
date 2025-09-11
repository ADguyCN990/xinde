package attachment

type OrphanRecord struct {
	ID          uint   `json:"id"`
	Filename    string `json:"filename"`
	StoragePath string `json:"storage_path"`
	UploadedBy  string `json:"uploaded_by"`
	CreatedAt   string `json:"created_at"`
}

type OrphanData struct {
	OrphanRecords []*OrphanRecord `json:"orphan_records"` //数据库有，磁盘没有
	OrphanFiles   []string        `json:"orphan_files"`   //数据库没有，磁盘有
}

type ScanInvalidResp struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"操作成功"`
	Success bool        `json:"success" example:"true"`
	Data    *OrphanData `json:"data"`
}
