package response

// 为了让 Swagger 能够识别响应结构，需要定义具体的响应模型

// BaseResponse 基础响应结构（用于 Swagger 文档）
type BaseResponse struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"操作成功"`
	Success bool   `json:"success" example:"true"`
}

// RegisterSuccessResponse 注册成功响应
type RegisterSuccessResponse struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"用户注册成功"`
	Data    struct {
		UserID uint `json:"user_id" example:"12345"`
	} `json:"data"`
	Success bool `json:"success" example:"true"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"参数错误"`
	Success bool   `json:"success" example:"false"`
}

// ValidationErrorResponse 参数验证错误响应
type ValidationErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"参数格式错误: Username is required"`
	Success bool   `json:"success" example:"false"`
}

// ConflictErrorResponse 冲突错误响应（如用户已存在）
type ConflictErrorResponse struct {
	Code    int    `json:"code" example:"409"`
	Message string `json:"message" example:"用户已存在，注册失败"`
	Success bool   `json:"success" example:"false"`
}

// InternalErrorResponse 服务器内部错误响应
type InternalErrorResponse struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"服务器内部错误"`
	Success bool   `json:"success" example:"false"`
}

// PageResponse 分页响应结构
type PageResponse struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"操作成功"`
	Data    struct {
		List     interface{} `json:"list"` // 实际使用时替换为具体类型
		Total    int64       `json:"total" example:"100"`
		Page     int         `json:"page" example:"1"`
		PageSize int         `json:"page_size" example:"10"`
		Pages    int         `json:"pages" example:"10"`
	} `json:"data"`
	Success bool `json:"success" example:"true"`
}

// UserListResponse 用户列表响应示例
type UserListResponse struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"操作成功"`
	Data    struct {
		List     []UserInfo `json:"list"` // 用户列表
		Total    int64      `json:"total" example:"100"`
		Page     int        `json:"page" example:"1"`
		PageSize int        `json:"page_size" example:"10"`
		Pages    int        `json:"pages" example:"10"`
	} `json:"data"`
	Success bool `json:"success" example:"true"`
}

// UserInfo 用户信息（示例）
type UserInfo struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"john_doe"`
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
}
