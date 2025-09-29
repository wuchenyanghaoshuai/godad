package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ErrorHandler provides unified error handling methods
type ErrorHandler struct{}

// NewErrorHandler creates a new error handler instance
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// HandleJSONBind handles JSON binding errors with standardized response
func (h *ErrorHandler) HandleJSONBind(ctx *gin.Context, req interface{}) bool {
	if err := ctx.ShouldBindJSON(req); err != nil {
		Error(ctx, CodeBadRequest, "请求参数错误: "+err.Error())
		return false
	}
	return true
}

// HandleQueryBind handles query parameter binding errors
func (h *ErrorHandler) HandleQueryBind(ctx *gin.Context, req interface{}) bool {
	if err := ctx.ShouldBindQuery(req); err != nil {
		Error(ctx, CodeBadRequest, "查询参数错误: "+err.Error())
		return false
	}
	return true
}

// RequireAuth checks user authentication and returns userID
func (h *ErrorHandler) RequireAuth(ctx *gin.Context) (uint, bool) {
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		Error(ctx, CodeUnauthorized, "用户未登录")
		return 0, false
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		Error(ctx, CodeUnauthorized, "用户ID格式错误")
		return 0, false
	}

	return userID, true
}

// ParseIDParam parses and validates ID parameter from URL
func (h *ErrorHandler) ParseIDParam(ctx *gin.Context, paramName string) (uint, bool) {
	idStr := ctx.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		Error(ctx, CodeBadRequest, paramName+"格式错误")
		return 0, false
	}
	return uint(id), true
}

// ParsePaginationParams parses and validates pagination parameters
func (h *ErrorHandler) ParsePaginationParams(ctx *gin.Context) (int, int) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	return page, size
}

// HandleServiceError handles service layer errors with appropriate HTTP status
func (h *ErrorHandler) HandleServiceError(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	errorMessage := err.Error()

	// Define error mappings
	errorMappings := map[string]int{
		"用户未登录":        CodeUnauthorized,
		"用户不存在":        CodeNotFound,
		"文章不存在":        CodeNotFound,
		"评论不存在":        CodeNotFound,
		"帖子不存在":        CodeNotFound,
		"回复不存在":        CodeNotFound,
		"分类不存在":        CodeNotFound,
		"资源不存在":        CodeNotFound,
		"无权限":           CodeForbidden,
		"无权限修改":        CodeForbidden,
		"无权限删除":        CodeForbidden,
		"无权限访问":        CodeForbidden,
		"权限不足":         CodeForbidden,
		"参数错误":         CodeBadRequest,
		"参数格式错误":      CodeBadRequest,
		"请求参数错误":      CodeBadRequest,
		"邮箱已存在":        CodeBadRequest,
		"用户名已存在":      CodeBadRequest,
		"密码错误":         CodeBadRequest,
		"验证码错误":        CodeBadRequest,
		"验证码已过期":      CodeBadRequest,
		"已关注该用户":      CodeBadRequest,
		"未关注该用户":      CodeBadRequest,
		"已收藏该文章":      CodeBadRequest,
		"未收藏该文章":      CodeBadRequest,
		"已点赞":          CodeBadRequest,
		"未点赞":          CodeBadRequest,
	}

	// Check for specific error patterns
	for pattern, code := range errorMappings {
		if strings.Contains(errorMessage, pattern) {
			Error(ctx, code, errorMessage)
			return
		}
	}

	// Default to internal server error
	Error(ctx, CodeInternalError, errorMessage)
}

// ValidateRequired checks if required fields are not empty
func (h *ErrorHandler) ValidateRequired(ctx *gin.Context, fieldName, value string) bool {
	if strings.TrimSpace(value) == "" {
		Error(ctx, CodeBadRequest, fieldName+"不能为空")
		return false
	}
	return true
}

// ValidateStringLength validates string length
func (h *ErrorHandler) ValidateStringLength(ctx *gin.Context, fieldName, value string, min, max int) bool {
	length := len(value)
	if length < min || length > max {
		Error(ctx, CodeBadRequest, fmt.Sprintf("%s长度必须在%d-%d个字符之间", fieldName, min, max))
		return false
	}
	return true
}

// ValidateEnum validates if value is in allowed enums
func (h *ErrorHandler) ValidateEnum(ctx *gin.Context, fieldName, value string, allowedValues []string) bool {
	for _, allowed := range allowedValues {
		if value == allowed {
			return true
		}
	}
	Error(ctx, CodeBadRequest, fmt.Sprintf("%s值无效，允许的值: %s", fieldName, strings.Join(allowedValues, ", ")))
	return false
}

// HandleMultipartFormBind handles multipart form binding errors
func (h *ErrorHandler) HandleMultipartFormBind(ctx *gin.Context, req interface{}) bool {
	if err := ctx.ShouldBind(req); err != nil {
		Error(ctx, CodeBadRequest, "表单参数错误: "+err.Error())
		return false
	}
	return true
}

// ValidateFileUpload validates uploaded file
func (h *ErrorHandler) ValidateFileUpload(ctx *gin.Context, fieldName string, allowedTypes []string, maxSize int64) bool {
	file, header, err := ctx.Request.FormFile(fieldName)
	if err != nil {
		Error(ctx, CodeBadRequest, "文件上传失败: "+err.Error())
		return false
	}
	defer file.Close()

	// Check file size
	if header.Size > maxSize {
		Error(ctx, CodeBadRequest, fmt.Sprintf("文件大小不能超过%dMB", maxSize/(1024*1024)))
		return false
	}

	// Check file type
	contentType := header.Header.Get("Content-Type")
	for _, allowedType := range allowedTypes {
		if strings.Contains(contentType, allowedType) {
			return true
		}
	}

	Error(ctx, CodeBadRequest, fmt.Sprintf("文件类型不支持，支持的类型: %s", strings.Join(allowedTypes, ", ")))
	return false
}

// ParseOptionalIDParam parses optional ID parameter from URL
func (h *ErrorHandler) ParseOptionalIDParam(ctx *gin.Context, paramName string) (uint, bool, bool) {
	idStr := ctx.Param(paramName)
	if idStr == "" {
		return 0, true, false // no error, but no ID provided
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		Error(ctx, CodeBadRequest, paramName+"格式错误")
		return 0, false, false
	}
	return uint(id), true, true
}