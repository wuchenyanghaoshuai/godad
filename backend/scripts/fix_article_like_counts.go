package main

import (
	"fmt"
	"log"

	"godad-backend/config"
	"godad-backend/models"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	if err := config.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	db := config.GetDB()

	fmt.Println("开始修复文章点赞数统计...")

	// 获取所有文章
	var articles []models.Article
	if err := db.Find(&articles).Error; err != nil {
		log.Fatalf("Failed to get articles: %v", err)
	}

	fixedCount := 0

	for _, article := range articles {
		// 统计该文章的实际点赞数量
		var actualLikeCount int64
		if err := db.Model(&models.Like{}).Where("target_type = ? AND target_id = ?", "article", article.ID).Count(&actualLikeCount).Error; err != nil {
			fmt.Printf("Error counting likes for article %d: %v\n", article.ID, err)
			continue
		}

		// 如果数量不匹配，更新
		if article.LikeCount != actualLikeCount {
			fmt.Printf("文章ID %d: 当前like_count=%d, 实际点赞数=%d, 正在更新...\n", 
				article.ID, article.LikeCount, actualLikeCount)
			
			if err := db.Model(&models.Article{}).Where("id = ?", article.ID).Update("like_count", actualLikeCount).Error; err != nil {
				fmt.Printf("Error updating article %d: %v\n", article.ID, err)
				continue
			}
			
			fmt.Printf("文章ID %d 点赞数已更新为 %d\n", article.ID, actualLikeCount)
			fixedCount++
		}
	}

	fmt.Printf("文章点赞数统计修复完成！共修复 %d 篇文章。\n", fixedCount)
}