package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// RequestLoggerMiddleware 详细请求日志中间件
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// 记录请求体（仅用于调试，生产环境建议关闭）
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}
		
		// 处理请求
		c.Next()
		
		// 计算处理时间
		latency := time.Since(start)
		
		// 记录日志
		log.Printf("[%s] %s %s - Status: %d, Latency: %v, IP: %s, UserAgent: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.URL.RawQuery,
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.UserAgent(),
		)
		
		// 如果有错误，记录错误信息
		if len(c.Errors) > 0 {
			log.Printf("Request errors: %v", c.Errors)
		}
	}
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Panic recovered: %v", recovered)
		c.JSON(500, gin.H{
			"code":    500,
			"message": "服务器内部错误",
			"data":    nil,
		})
	})
}

