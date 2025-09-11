package main

import (
	"fmt"
	"log"

	"godad-backend/config"
	"godad-backend/services"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	if err := config.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	fmt.Println("开始清理所有相关缓存...")

	// 初始化缓存服务
	cacheService := services.NewCacheService()

	// 清理文章缓存
	cacheService.DeletePattern("article:*")
	cacheService.DeletePattern("articles:*")
	
	fmt.Println("缓存清理完成！")
}