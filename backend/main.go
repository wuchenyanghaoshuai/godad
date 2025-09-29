package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"godad-backend/config"
	"godad-backend/models"
	"godad-backend/observability"
	"godad-backend/routes"
	"godad-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("启动GoDad后端服务...")

	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Printf("警告: 无法加载.env文件: %v", err)
	}

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化 OpenTelemetry
	shutdownTracing := observability.InitTracing(context.Background(), cfg)
	defer func() {
		if err := shutdownTracing(context.Background()); err != nil {
			log.Printf("关闭 tracing 失败: %v", err)
		}
	}()

	// 打印数据库连接信息（仅开发环境）
	if cfg.Server.Environment == "development" {
		log.Printf("数据库连接信息: Host=%s, Port=%d, User=%s, DBName=%s",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.DBName)
	}

	// 初始化数据库
	if err := config.InitDatabase(cfg); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer config.CloseDatabase()

	// 初始化Redis缓存
	if err := services.InitRedis(); err != nil {
		log.Printf("警告: Redis缓存初始化失败: %v", err)
		log.Println("缓存功能已禁用，相关特性将退化运行")
	} else {
		log.Println("Redis缓存初始化成功")
	}

	// 检查是否需要自动迁移（仅开发环境）
	if cfg.Server.Environment == "development" && cfg.Database.AutoMigrate {
		log.Println("开发环境：执行数据库自动迁移...")
		if err := models.AutoMigrate(config.GetDB()); err != nil {
			log.Fatalf("数据库迁移失败: %v", err)
		}
		log.Println("数据库迁移完成")
	} else {
		log.Printf("跳过数据库自动迁移（环境=%s, AutoMigrate=%v）", cfg.Server.Environment, cfg.Database.AutoMigrate)
	}

	// 根据环境设置Gin运行模式
	switch cfg.Server.Environment {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "testing":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// 设置路由
	router := routes.SetupRoutes()

	// 创建HTTP服务器
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 启动服务器
	go func() {
		log.Printf("GoDad Backend服务启动成功，监听端口: %d", cfg.Server.Port)
		log.Printf("服务地址: http://127.0.0.1:%d", cfg.Server.Port)
		log.Printf("API文档: http://127.0.0.1:%d/api", cfg.Server.Port)
		log.Printf("健康检查: http://127.0.0.1:%d/health", cfg.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 优雅关闭服务器，等待现有连接完成
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}

	log.Println("服务器已关闭")
}
