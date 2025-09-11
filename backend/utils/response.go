package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
}

// PageResponse 分页响应结构体
type PageResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`   // 总数
	Page    int         `json:"page"`    // 当前页
	Size    int         `json:"size"`    // 每页大小
}

// 状态码常量
const (
	// 成功
	CodeSuccess = 200
	
	// 客户端错误
	CodeBadRequest     = 400  // 请求参数错误
	CodeUnauthorized   = 401  // 未授权
	CodeForbidden      = 403  // 禁止访问
	CodeNotFound       = 404  // 资源不存在
	CodeValidationFail = 422  // 验证失败
	
	// 服务器错误
	CodeInternalError = 500  // 内部服务器错误
	CodeDatabaseError = 501  // 数据库错误
	CodeExternalError = 502  // 外部服务错误
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// SuccessPage 分页成功响应
func SuccessPage(c *gin.Context, data interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, PageResponse{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
		Total:   total,
		Page:    page,
		Size:    size,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	httpStatus := getHTTPStatus(code)
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, message string) {
	Error(c, CodeBadRequest, message)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
	Error(c, CodeUnauthorized, message)
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message string) {
	Error(c, CodeForbidden, message)
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string) {
	Error(c, CodeNotFound, message)
}

// ValidationFail 验证失败
func ValidationFail(c *gin.Context, message string) {
	Error(c, CodeValidationFail, message)
}

// InternalError 内部服务器错误
func InternalError(c *gin.Context, message string) {
	Error(c, CodeInternalError, message)
}

// DatabaseError 数据库错误
func DatabaseError(c *gin.Context, message string) {
	Error(c, CodeDatabaseError, message)
}

// ErrorResponse 错误响应（兼容旧版本）
func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Code:    status,
		Message: message,
		Data:    nil,
	})
}

// SuccessResponse 成功响应（兼容旧版本）
func SuccessResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// PaginationResponse 分页响应结构体
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// getHTTPStatus 根据业务状态码获取HTTP状态码
func getHTTPStatus(code int) int {
	switch {
	case code == CodeSuccess:
		return http.StatusOK
	case code >= 400 && code < 500:
		return code
	case code >= 500:
		return http.StatusInternalServerError
	default:
		return http.StatusOK
	}
}