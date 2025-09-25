package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParseUintParam 解析 URL 参数中的 uint 值
func ParseUintParam(c *gin.Context, paramName string) (uint, error) {
	paramStr := c.Param(paramName)
	id, err := strconv.ParseUint(paramStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// ParseIntParam 解析 URL 参数中的 int 值
func ParseIntParam(c *gin.Context, paramName string) (int, error) {
	paramStr := c.Param(paramName)
	return strconv.Atoi(paramStr)
}

// ParsePaginationParams 解析分页参数
func ParsePaginationParams(c *gin.Context) (page, size int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ = strconv.Atoi(c.DefaultQuery("size", "10"))

	// 参数验证
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	return page, size
}

// GetUserID 从上下文中获取当前用户ID
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	switch v := userID.(type) {
	case uint:
		return v, true
	case int:
		return uint(v), true
	case float64:
		return uint(v), true
	default:
		return 0, false
	}
}

// GetUserRole 从上下文中获取当前用户角色
func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}

	roleStr, ok := role.(string)
	return roleStr, ok
}

// ValidateRequiredFields 验证必填字段
func ValidateRequiredFields(fields map[string]interface{}) []string {
	var missingFields []string

	for fieldName, fieldValue := range fields {
		if fieldValue == nil {
			missingFields = append(missingFields, fieldName)
			continue
		}

		switch v := fieldValue.(type) {
		case string:
			if v == "" {
				missingFields = append(missingFields, fieldName)
			}
		case uint:
			if v == 0 {
				missingFields = append(missingFields, fieldName)
			}
		case int:
			if v == 0 {
				missingFields = append(missingFields, fieldName)
			}
		}
	}

	return missingFields
}