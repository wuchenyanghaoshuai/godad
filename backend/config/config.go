package config

import (
	"os"
	"strconv"
	"sync"
)

var (
	globalConfig *Config
	once         sync.Once
)

// Config 应用配置结构体
type Config struct {
	Database  DatabaseConfig
	JWT       JWTConfig
	OSS       OSSConfig
	Server    ServerConfig
	RateLimit RateLimitConfig
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	DBName      string
	DSN         string
	AutoMigrate bool // 是否自动迁移数据库表结构
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string
	ExpireHours int
}

// OSSConfig 阿里云OSS配置
type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	CustomDomain    string // 自定义域名，用于生成访问URL
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host        string
	Port        int
	FrontendURL string
	Environment string // 运行环境: development, production, testing
}

// RateLimitConfig 速率限制配置
type RateLimitConfig struct {
	General GeneralLimitConfig // 通用请求限制
	Auth    AuthLimitConfig    // 认证请求限制
	Upload  UploadLimitConfig  // 上传请求限制
}

// GeneralLimitConfig 通用请求限制配置
type GeneralLimitConfig struct {
	RequestsPerSecond float64 // 每秒请求数
	BurstSize        int     // 突发请求数
	CleanupMinutes   int     // 清理间隔（分钟）
}

// AuthLimitConfig 认证请求限制配置
type AuthLimitConfig struct {
	RequestsPerMinute int // 每分钟请求数
	BurstSize        int // 突发请求数
	CleanupMinutes   int // 清理间隔（分钟）
}

// UploadLimitConfig 上传请求限制配置
type UploadLimitConfig struct {
	RequestsPerMinute int // 每分钟请求数
	BurstSize        int // 突发请求数
	CleanupMinutes   int // 清理间隔（分钟）
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	config := &Config{}
	
	// 数据库配置
	config.Database.Host = getEnv("DB_HOST", "127.0.0.1")
	config.Database.Port = getEnvAsInt("DB_PORT", 3306)
	config.Database.User = getEnv("DB_USER", "root")
	config.Database.Password = getEnv("DB_PASSWORD", "")
	config.Database.DBName = getEnv("DB_NAME", "godad")
	config.Database.AutoMigrate = getEnvAsBool("DB_AUTO_MIGRATE", true) // 默认开发环境启用
	
	// 构建DSN
	config.Database.DSN = config.Database.User + ":" + config.Database.Password + 
		"@tcp(" + config.Database.Host + ":" + strconv.Itoa(config.Database.Port) + ")/" + 
		config.Database.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	
	// JWT配置
	config.JWT.Secret = getEnv("JWT_SECRET", "your-jwt-secret-key")
	config.JWT.ExpireHours = getEnvAsInt("JWT_EXPIRE_HOURS", 24)
	
	// OSS配置
	config.OSS.Endpoint = getEnv("OSS_ENDPOINT", "oss-cn-beijing.aliyuncs.com")
	config.OSS.AccessKeyID = getEnv("OSS_ACCESS_KEY_ID", "")
	config.OSS.AccessKeySecret = getEnv("OSS_ACCESS_KEY_SECRET", "")
	config.OSS.BucketName = getEnv("OSS_BUCKET_NAME", "godad")
	config.OSS.CustomDomain = getEnv("OSS_CUSTOM_DOMAIN", "")
	
	// 服务器配置
	config.Server.Host = getEnv("SERVER_HOST", "127.0.0.1")
	config.Server.Port = getEnvAsInt("SERVER_PORT", 8888)
	config.Server.FrontendURL = getEnv("FRONTEND_URL", "http://127.0.0.1:3333")
	config.Server.Environment = getEnv("SERVER_ENV", "development") // development, production, testing
	
	// 速率限制配置
	config.RateLimit.General.RequestsPerSecond = getEnvAsFloat("RATE_LIMIT_GENERAL_RPS", 10.0)
	config.RateLimit.General.BurstSize = getEnvAsInt("RATE_LIMIT_GENERAL_BURST", 20)
	config.RateLimit.General.CleanupMinutes = getEnvAsInt("RATE_LIMIT_GENERAL_CLEANUP", 30)
	
	config.RateLimit.Auth.RequestsPerMinute = getEnvAsInt("RATE_LIMIT_AUTH_RPM", 5)
	config.RateLimit.Auth.BurstSize = getEnvAsInt("RATE_LIMIT_AUTH_BURST", 3)
	config.RateLimit.Auth.CleanupMinutes = getEnvAsInt("RATE_LIMIT_AUTH_CLEANUP", 60)
	
	config.RateLimit.Upload.RequestsPerMinute = getEnvAsInt("RATE_LIMIT_UPLOAD_RPM", 3)
	config.RateLimit.Upload.BurstSize = getEnvAsInt("RATE_LIMIT_UPLOAD_BURST", 2)
	config.RateLimit.Upload.CleanupMinutes = getEnvAsInt("RATE_LIMIT_UPLOAD_CLEANUP", 60)
	
	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为整数
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool 获取环境变量并转换为布尔值
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvAsFloat 获取环境变量并转换为浮点数
func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// GetConfig 获取全局配置（单例模式）
func GetConfig() *Config {
	once.Do(func() {
		globalConfig = LoadConfig()
	})
	return globalConfig
}