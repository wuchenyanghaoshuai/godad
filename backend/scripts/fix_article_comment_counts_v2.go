package main

import (
	"fmt"
	"log"

	"godad-backend/config"
	"godad-backend/models"
	"godad-backend/services"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	if err := config.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	db := config.GetDB()

	// 初始化缓存服务
	cacheService := services.NewCacheService()

	fmt.Println("开始修复文章评论数统计并清理缓存...")

	// 获取所有文章
	var articles []models.Article
	if err := db.Find(&articles).Error; err != nil {
		log.Fatalf("Failed to get articles: %v", err)
	}

	fixedCount := 0

	for _, article := range articles {
		// 统计该文章的主评论数量（parent_id = NULL）
		var actualCommentCount int64
		if err := db.Model(&models.Comment{}).Where("article_id = ? AND parent_id IS NULL", article.ID).Count(&actualCommentCount).Error; err != nil {
			fmt.Printf("Error counting comments for article %d: %v\n", article.ID, err)
			continue
		}

		// 如果数量不匹配，更新
		if article.CommentCount != actualCommentCount {
			fmt.Printf("文章ID %d: 当前comment_count=%d, 实际主评论数=%d, 正在更新...\n", 
				article.ID, article.CommentCount, actualCommentCount)
			
			if err := db.Model(&models.Article{}).Where("id = ?", article.ID).Update("comment_count", actualCommentCount).Error; err != nil {
				fmt.Printf("Error updating article %d: %v\n", article.ID, err)
				continue
			}
			
			// 清理相关缓存
			articleKey := fmt.Sprintf("article:%d", article.ID)
			cacheService.Delete(articleKey)
			
			// 清理文章列表缓存
			cacheService.DeletePattern("articles:*")
			
			fmt.Printf("文章ID %d 评论数已更新为 %d，缓存已清理\n", article.ID, actualCommentCount)
			fixedCount++
		}
	}

	fmt.Printf("文章评论数统计修复完成！共修复 %d 篇文章。\n", fixedCount)
}