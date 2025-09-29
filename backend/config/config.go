package config

import (
	"log"
	"strconv"
	"sync"

	"godad-backend/utils"
)

var (
	globalConfig *Config
	once         sync.Once
)

// Config 应用配置结构体
type Config struct {
	Database      DatabaseConfig
	JWT           JWTConfig
	OSS           OSSConfig
	Server        ServerConfig
	RateLimit     RateLimitConfig
	Observability ObservabilityConfig
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
	BurstSize         int     // 突发请求数
	CleanupMinutes    int     // 清理间隔（分钟）
}

// AuthLimitConfig 认证请求限制配置
type AuthLimitConfig struct {
	RequestsPerMinute int // 每分钟请求数
	BurstSize         int // 突发请求数
	CleanupMinutes    int // 清理间隔（分钟）
}

// UploadLimitConfig 上传请求限制配置
type UploadLimitConfig struct {
	RequestsPerMinute int // 每分钟请求数
	BurstSize         int // 突发请求数
	CleanupMinutes    int // 清理间隔（分钟）
}

// ObservabilityConfig 可观测性配置
type ObservabilityConfig struct {
	OTLPExporterEndpoint string
	OTLPInsecure         bool
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	config := &Config{}

	// 服务器配置（先确定环境）
	config.Server.Environment = utils.GetEnv("SERVER_ENV", "development")
	config.Server.Host = utils.GetEnv("SERVER_HOST", "127.0.0.1")
	config.Server.Port = utils.GetEnvAsInt("SERVER_PORT", 8888)
	config.Server.FrontendURL = utils.GetEnv("FRONTEND_URL", "http://127.0.0.1:3333")

	// 数据库配置
	config.Database.Host = utils.GetEnv("DB_HOST", "127.0.0.1")
	config.Database.Port = utils.GetEnvAsInt("DB_PORT", 3306)
	config.Database.User = utils.GetEnv("DB_USER", "root")
	config.Database.Password = utils.GetEnv("DB_PASSWORD", "")
	config.Database.DBName = utils.GetEnv("DB_NAME", "godad")
	config.Database.AutoMigrate = utils.GetEnvAsBool("DB_AUTO_MIGRATE", false)

	config.Database.DSN = config.Database.User + ":" + config.Database.Password +
		"@tcp(" + config.Database.Host + ":" + strconv.Itoa(config.Database.Port) + ")/" +
		config.Database.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	// JWT配置
	config.JWT.Secret = utils.GetEnv("JWT_SECRET", "")
	config.JWT.ExpireHours = utils.GetEnvAsInt("JWT_EXPIRE_HOURS", 24)
	if config.JWT.Secret == "" {
		if config.Server.Environment == "production" {
			log.Fatal("JWT_SECRET must be provided in production environment")
		}
		log.Println("警告: 未设置 JWT_SECRET，使用仅限开发环境的临时密钥")
		config.JWT.Secret = "dev-only-change-me"
	}

	// OSS配置
	config.OSS.Endpoint = utils.GetEnv("OSS_ENDPOINT", "oss-cn-beijing.aliyuncs.com")
	config.OSS.AccessKeyID = utils.GetEnv("OSS_ACCESS_KEY_ID", "")
	config.OSS.AccessKeySecret = utils.GetEnv("OSS_ACCESS_KEY_SECRET", "")
	config.OSS.BucketName = utils.GetEnv("OSS_BUCKET_NAME", "godad")
	config.OSS.CustomDomain = utils.GetEnv("OSS_CUSTOM_DOMAIN", "")

	// 速率限制配置
	config.RateLimit.General.RequestsPerSecond = utils.GetEnvAsFloat("RATE_LIMIT_GENERAL_RPS", 10.0)
	config.RateLimit.General.BurstSize = utils.GetEnvAsInt("RATE_LIMIT_GENERAL_BURST", 20)
	config.RateLimit.General.CleanupMinutes = utils.GetEnvAsInt("RATE_LIMIT_GENERAL_CLEANUP", 30)

	config.RateLimit.Auth.RequestsPerMinute = utils.GetEnvAsInt("RATE_LIMIT_AUTH_RPM", 5)
	config.RateLimit.Auth.BurstSize = utils.GetEnvAsInt("RATE_LIMIT_AUTH_BURST", 3)
	config.RateLimit.Auth.CleanupMinutes = utils.GetEnvAsInt("RATE_LIMIT_AUTH_CLEANUP", 60)

	config.RateLimit.Upload.RequestsPerMinute = utils.GetEnvAsInt("RATE_LIMIT_UPLOAD_RPM", 3)
	config.RateLimit.Upload.BurstSize = utils.GetEnvAsInt("RATE_LIMIT_UPLOAD_BURST", 2)
	config.RateLimit.Upload.CleanupMinutes = utils.GetEnvAsInt("RATE_LIMIT_UPLOAD_CLEANUP", 60)

	// 可观测性配置
	config.Observability.OTLPExporterEndpoint = utils.GetEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	config.Observability.OTLPInsecure = utils.GetEnvAsBool("OTEL_EXPORTER_OTLP_INSECURE", false)

	return config
}

// GetConfig 获取全局配置（单例模式）
func GetConfig() *Config {
	once.Do(func() {
		globalConfig = LoadConfig()
	})
	return globalConfig
}
