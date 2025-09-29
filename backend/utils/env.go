package utils

import (
	"os"
	"strconv"
)

// GetEnv 获取环境变量，如果不存在则返回默认值
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsInt 获取环境变量并转换为整数
func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvAsBool 获取环境变量并转换为布尔值
func GetEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// GetEnvAsFloat 获取环境变量并转换为浮点数
func GetEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// MustGetEnv 获取环境变量，如果不存在则panic
// 用于必须配置的环境变量
func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("Required environment variable " + key + " is not set")
	}
	return value
}

// IsProduction 检查当前是否为生产环境
func IsProduction() bool {
	return GetEnv("SERVER_ENV", "development") == "production"
}

// IsDevelopment 检查当前是否为开发环境
func IsDevelopment() bool {
	return GetEnv("SERVER_ENV", "development") == "development"
}

// IsTesting 检查当前是否为测试环境
func IsTesting() bool {
	return GetEnv("SERVER_ENV", "development") == "testing"
}