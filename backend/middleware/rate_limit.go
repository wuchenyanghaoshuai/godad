package middleware

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"godad-backend/config"
	"godad-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	limiter  *rate.Limiter
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
			limiter:  rate.NewLimiter(i.rate, i.capacity),
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

type rateLimiterAdapter struct {
	local     *IPRateLimiter
	redis     *redis.Client
	window    time.Duration
	max       int
	keyPrefix string
}

var redisRateLimiterScript = redis.NewScript(`
local current = redis.call('INCR', KEYS[1])
if tonumber(current) == 1 then
    redis.call('PEXPIRE', KEYS[1], ARGV[2])
end
if tonumber(current) > tonumber(ARGV[1]) then
    return 0
end
return 1
`)

func newRateLimiterAdapter(local *IPRateLimiter, window time.Duration, max int, prefix string, client *redis.Client) *rateLimiterAdapter {
	if max < 1 {
		max = 1
	}
	return &rateLimiterAdapter{
		local:     local,
		redis:     client,
		window:    window,
		max:       max,
		keyPrefix: prefix,
	}
}

func (r *rateLimiterAdapter) Allow(ip string) bool {
	if r.redis != nil && r.max > 0 {
		allowed, err := allowWithRedis(r.redis, fmt.Sprintf("%s:%s", r.keyPrefix, ip), r.window, r.max)
		if err == nil {
			return allowed
		}
	}

	return r.local.GetLimiter(ip).Allow()
}

func allowWithRedis(client *redis.Client, key string, window time.Duration, max int) (bool, error) {
	res, err := redisRateLimiterScript.Run(context.Background(), client, []string{key}, max, window.Milliseconds()).Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

// 全局限制器实例 - 延迟初始化
var (
	generalLimiter *rateLimiterAdapter
	authLimiter    *rateLimiterAdapter
	uploadLimiter  *rateLimiterAdapter
	initOnce       sync.Once
)

// initLimiters 初始化限制器
func initLimiters() {
	initOnce.Do(func() {
		cfg := config.GetConfig()
		redisClient := services.GetRedisClient()

		generalRate := rate.Limit(cfg.RateLimit.General.RequestsPerSecond)
		if generalRate <= 0 {
			generalRate = rate.Limit(10)
		}
		generalBurst := cfg.RateLimit.General.BurstSize
		if generalBurst < 1 {
			generalBurst = 1
		}
		generalLimiter = newRateLimiterAdapter(
			NewIPRateLimiter(generalRate, generalBurst, time.Duration(cfg.RateLimit.General.CleanupMinutes)*time.Minute),
			time.Second,
			int(math.Ceil(float64(generalRate)))+generalBurst,
			"rate:general",
			redisClient,
		)

		authRPM := cfg.RateLimit.Auth.RequestsPerMinute
		if authRPM < 1 {
			authRPM = 1
		}
		authBurst := cfg.RateLimit.Auth.BurstSize
		if authBurst < 1 {
			authBurst = 1
		}
		authLimiter = newRateLimiterAdapter(
			NewIPRateLimiter(rate.Every(time.Minute/time.Duration(authRPM)), authBurst, time.Duration(cfg.RateLimit.Auth.CleanupMinutes)*time.Minute),
			time.Minute,
			authRPM+authBurst,
			"rate:auth",
			redisClient,
		)

		uploadRPM := cfg.RateLimit.Upload.RequestsPerMinute
		if uploadRPM < 1 {
			uploadRPM = 1
		}
		uploadBurst := cfg.RateLimit.Upload.BurstSize
		if uploadBurst < 1 {
			uploadBurst = 1
		}
		uploadLimiter = newRateLimiterAdapter(
			NewIPRateLimiter(rate.Every(time.Minute/time.Duration(uploadRPM)), uploadBurst, time.Duration(cfg.RateLimit.Upload.CleanupMinutes)*time.Minute),
			time.Minute,
			uploadRPM+uploadBurst,
			"rate:upload",
			redisClient,
		)
	})
}

// RateLimitMiddleware 通用速率限制中间件
func RateLimitMiddleware(adapter *rateLimiterAdapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := GetClientIP(c)
		if !adapter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
				"code":  429,
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
	if requestsPerSecond <= 0 {
		requestsPerSecond = 1
	}
	if burstSize <= 0 {
		burstSize = 1
	}
	localLimiter := NewIPRateLimiter(rate.Limit(requestsPerSecond), burstSize, 30*time.Minute)
	adapter := newRateLimiterAdapter(localLimiter, time.Second, int(math.Ceil(requestsPerSecond))+burstSize, "rate:custom", services.GetRedisClient())
	return RateLimitMiddleware(adapter)
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
				"code":  429,
			})
			c.Abort()
			return
		}

		// 记录当前请求时间
		log.timestamps = append(log.timestamps, now)

		c.Next()
	}
}
