package middleware

import (
	"godad-backend/utils"

	"github.com/gin-gonic/gin"
)

// RequireAuth 必须登录的中间件
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := GetCurrentUserID(c)
		if !exists || userID == 0 {
			utils.Error(c, utils.CodeUnauthorized, "请先登录")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireAdmin 必须是管理员的中间件
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 首先检查是否登录
		userID, exists := GetCurrentUserID(c)
		if !exists || userID == 0 {
			utils.Error(c, utils.CodeUnauthorized, "请先登录")
			c.Abort()
			return
		}

		// 检查是否是管理员
		role, roleExists := c.Get("user_role")
		if !roleExists {
			utils.Error(c, utils.CodeForbidden, "无权访问")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok || roleStr != "admin" {
			utils.Error(c, utils.CodeForbidden, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth 可选登录的中间件
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试解析token，但不强制要求
		token := c.GetHeader("Authorization")
		if token != "" {
			// 这里可以添加token解析逻辑，但失败时不中断请求
			// 让后续处理函数自己决定是否需要用户信息
		}
		c.Next()
	}
}

// ValidateOwnershipOrAdmin 验证资源所有权或管理员权限
func ValidateOwnershipOrAdmin(userIDFromResource uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserID, exists := GetCurrentUserID(c)
		if !exists {
			utils.Error(c, utils.CodeUnauthorized, "请先登录")
			c.Abort()
			return
		}

		// 如果是资源所有者，直接通过
		if currentUserID == userIDFromResource {
			c.Next()
			return
		}

		// 如果不是所有者，检查是否是管理员
		role, roleExists := c.Get("user_role")
		if !roleExists {
			utils.Error(c, utils.CodeForbidden, "无权访问")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok || roleStr != "admin" {
			utils.Error(c, utils.CodeForbidden, "无权访问此资源")
			c.Abort()
			return
		}

		c.Next()
	}
}