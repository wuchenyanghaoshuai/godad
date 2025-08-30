package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware CORS跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的源
		c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:3333")
		
		// 设置允许的请求方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		
		// 设置允许的请求头
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		
		// 设置允许携带凭证
		c.Header("Access-Control-Allow-Credentials", "true")
		
		// 设置预检请求的缓存时间
		c.Header("Access-Control-Max-Age", "86400")
		
		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// SecurityMiddleware 安全中间件
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置安全头
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		c.Next()
	}
}