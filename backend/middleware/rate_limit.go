package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"godad-backend/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	limiter *rate.Limiter
	lastSeen time.Time
}

// IPRateLimiter IP速率限制器管理器
type IPRateLimiter struct {
	limiters map[string]*RateLimiter
	mutex    sync.RWMutex
	rate     rate.Limit
	capacity int
	cleanup  time.Duration
}

// NewIPRateLimiter 创建IP速率限制器
func NewIPRateLimiter(r rate.Limit, capacity int, cleanup time.Duration) *IPRateLimiter {
	limiter := &IPRateLimiter{
		limiters: make(map[string]*RateLimiter),
		rate:     r,
		capacity: capacity,
		cleanup:  cleanup,
	}

	// 启动清理协程
	go limiter.cleanupRoutine()
	
	return limiter
}

// GetLimiter 获取指定IP的限制器
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	limiter, exists := i.limiters[ip]
	if !exists {
		limiter = &RateLimiter{
			limiter: rate.NewLimiter(i.rate, i.capacity),
			lastSeen: time.Now(),
		}
		i.limiters[ip] = limiter
	} else {
		limiter.lastSeen = time.Now()
	}

	return limiter.limiter
}

// cleanupRoutine 清理长时间未使用的限制器
func (i *IPRateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(i.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		i.mutex.Lock()
		for ip, limiter := range i.limiters {
			if time.Since(limiter.lastSeen) > i.cleanup {
				delete(i.limiters, ip)
			}
		}
		i.mutex.Unlock()
	}
}

// 全局限制器实例 - 延迟初始化
var (
	generalLimiter *IPRateLimiter
	authLimiter    *IPRateLimiter
	uploadLimiter  *IPRateLimiter
	initOnce       sync.Once
)

// initLimiters 初始化限制器
func initLimiters() {
	initOnce.Do(func() {
		cfg := config.GetConfig()
		
		// 普通请求限制器
		generalLimiter = NewIPRateLimiter(
			rate.Limit(cfg.RateLimit.General.RequestsPerSecond),
			cfg.RateLimit.General.BurstSize,
			time.Duration(cfg.RateLimit.General.CleanupMinutes)*time.Minute,
		)
		
		// 登录限制器
		authLimiter = NewIPRateLimiter(
			rate.Every(time.Minute/time.Duration(cfg.RateLimit.Auth.RequestsPerMinute)),
			cfg.RateLimit.Auth.BurstSize,
			time.Duration(cfg.RateLimit.Auth.CleanupMinutes)*time.Minute,
		)
		
		// 上传限制器
		uploadLimiter = NewIPRateLimiter(
			rate.Every(time.Minute/time.Duration(cfg.RateLimit.Upload.RequestsPerMinute)),
			cfg.RateLimit.Upload.BurstSize,
			time.Duration(cfg.RateLimit.Upload.CleanupMinutes)*time.Minute,
		)
	})
}

// RateLimitMiddleware 通用速率限制中间件
func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := GetClientIP(c)
		rateLimiter := limiter.GetLimiter(ip)
		
		if !rateLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
				"code": 429,
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// GeneralRateLimit 通用请求速率限制
func GeneralRateLimit() gin.HandlerFunc {
	initLimiters()
	return RateLimitMiddleware(generalLimiter)
}

// AuthRateLimit 认证请求速率限制
func AuthRateLimit() gin.HandlerFunc {
	initLimiters()
	return RateLimitMiddleware(authLimiter)
}

// UploadRateLimit 上传请求速率限制
func UploadRateLimit() gin.HandlerFunc {
	initLimiters()
	return RateLimitMiddleware(uploadLimiter)
}

// GetClientIP 获取客户端真实IP
func GetClientIP(c *gin.Context) string {
	// 检查 X-Forwarded-For 头
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		return xff
	}
	
	// 检查 X-Real-IP 头
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}
	
	// 返回RemoteAddr
	return c.ClientIP()
}

// CustomRateLimit 自定义速率限制
func CustomRateLimit(requestsPerSecond float64, burstSize int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(rate.Limit(requestsPerSecond), burstSize, 30*time.Minute)
	return RateLimitMiddleware(limiter)
}

// SlidingWindowRateLimit 滑动窗口速率限制
func SlidingWindowRateLimit(maxRequests int, windowDuration time.Duration) gin.HandlerFunc {
	type requestLog struct {
		timestamps []time.Time
		mutex      sync.RWMutex
	}
	
	ipLogs := sync.Map{}
	
	// 清理过期记录的协程
	go func() {
		ticker := time.NewTicker(windowDuration)
		defer ticker.Stop()
		
		for range ticker.C {
			now := time.Now()
			ipLogs.Range(func(key, value interface{}) bool {
				log := value.(*requestLog)
				log.mutex.Lock()
				
				// 移除过期的时间戳
				cutoff := now.Add(-windowDuration)
				validTimestamps := make([]time.Time, 0, len(log.timestamps))
				for _, ts := range log.timestamps {
					if ts.After(cutoff) {
						validTimestamps = append(validTimestamps, ts)
					}
				}
				log.timestamps = validTimestamps
				
				log.mutex.Unlock()
				return true
			})
		}
	}()
	
	return func(c *gin.Context) {
		ip := GetClientIP(c)
		now := time.Now()
		
		// 获取或创建该IP的请求日志
		logInterface, _ := ipLogs.LoadOrStore(ip, &requestLog{
			timestamps: make([]time.Time, 0),
		})
		log := logInterface.(*requestLog)
		
		log.mutex.Lock()
		defer log.mutex.Unlock()
		
		// 移除过期的时间戳
		cutoff := now.Add(-windowDuration)
		validTimestamps := make([]time.Time, 0, len(log.timestamps))
		for _, ts := range log.timestamps {
			if ts.After(cutoff) {
				validTimestamps = append(validTimestamps, ts)
			}
		}
		log.timestamps = validTimestamps
		
		// 检查是否超过限制
		if len(log.timestamps) >= maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("在%v内最多允许%d个请求", windowDuration, maxRequests),
				"code": 429,
			})
			c.Abort()
			return
		}
		
		// 记录当前请求时间
		log.timestamps = append(log.timestamps, now)
		
		c.Next()
	}
}