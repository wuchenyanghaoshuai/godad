package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库实例
var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase(config *Config) error {
	// 配置GORM日志
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	
	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(config.Database.DSN), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// 获取底层sql.DB对象进行连接池配置
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)                   // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                  // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)         // 连接最大生存时间
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)  // 连接最大空闲时间
	
	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	
	log.Println("Database connected successfully")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}