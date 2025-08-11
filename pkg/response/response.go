package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`           // 业务状态码
	Message string      `json:"message"`        // 响应消息
	Data    interface{} `json:"data,omitempty"` // 响应数据，空时不返回该字段
	Success bool        `json:"success"`        // 是否成功
}

// 预定义的业务状态码
const (
	CodeSuccess            = 200 // 成功
	CodeInvalidParams      = 400 // 参数错误
	CodeUnauthorized       = 401 // 未认证
	CodeForbidden          = 403 // 禁止访问
	CodeNotFound           = 404 // 资源不存在
	CodeConflict           = 409 // 资源冲突（如用户已存在）
	CodeInternalError      = 500 // 服务器内部错误
	CodeServiceUnavailable = 503 // 服务不可用
)

// 预定义的响应消息
const (
	MsgSuccess            = "操作成功"
	MsgInvalidParams      = "参数错误"
	MsgUnauthorized       = "未授权访问"
	MsgForbidden          = "禁止访问"
	MsgNotFound           = "资源不存在"
	MsgInternalError      = "服务器内部错误"
	MsgServiceUnavailable = "服务暂不可用"
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: MsgSuccess,
		Data:    data,
		Success: true,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
		Success: true,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpCode int, businessCode int, message string) {
	c.JSON(httpCode, Response{
		Code:    businessCode,
		Message: message,
		Success: false,
	})
}

// 便捷的错误响应方法
func BadRequest(c *gin.Context, message string) {
	if message == "" {
		message = MsgInvalidParams
	}
	Error(c, http.StatusBadRequest, CodeInvalidParams, message)
}

func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = MsgUnauthorized
	}
	Error(c, http.StatusUnauthorized, CodeUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = MsgForbidden
	}
	Error(c, http.StatusForbidden, CodeForbidden, message)
}

func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = MsgNotFound
	}
	Error(c, http.StatusNotFound, CodeNotFound, message)
}

func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, CodeConflict, message)
}

func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = MsgInternalError
	}
	Error(c, http.StatusInternalServerError, CodeInternalError, message)
}

func ServiceUnavailable(c *gin.Context, message string) {
	if message == "" {
		message = MsgServiceUnavailable
	}
	Error(c, http.StatusServiceUnavailable, CodeServiceUnavailable, message)
}
