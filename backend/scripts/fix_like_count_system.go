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

	// 初始化Redis
	if err := services.InitRedis(); err != nil {
		fmt.Printf("Redis初始化失败，跳过缓存清理: %v\n", err)
	}

	db := config.GetDB()

	fmt.Println("开始系统性修复点赞数统计问题...")

	// 1. 修复数据库中的点赞计数
	fmt.Println("步骤1: 修复数据库点赞计数...")
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
			fmt.Printf("修复文章ID %d: 当前like_count=%d -> 实际点赞数=%d\n", 
				article.ID, article.LikeCount, actualLikeCount)
			
			if err := db.Model(&models.Article{}).Where("id = ?", article.ID).Update("like_count", actualLikeCount).Error; err != nil {
				fmt.Printf("Error updating article %d: %v\n", article.ID, err)
				continue
			}
			fixedCount++
		}
	}
	fmt.Printf("数据库点赞计数修复完成！共修复 %d 篇文章。\n", fixedCount)

	// 2. 清除所有相关缓存
	fmt.Println("步骤2: 清除所有相关缓存...")
	cacheService := services.NewCacheService()
	
	// 清除文章列表缓存
	if err := cacheService.DeletePattern("articles:list:*"); err != nil {
		fmt.Printf("清除文章列表缓存失败: %v\n", err)
	} else {
		fmt.Println("✓ 文章列表缓存已清除")
	}

	// 清除单个文章缓存
	if err := cacheService.DeletePattern("article:*"); err != nil {
		fmt.Printf("清除文章缓存失败: %v\n", err)
	} else {
		fmt.Println("✓ 单个文章缓存已清除")
	}

	// 清除搜索结果缓存
	if err := cacheService.DeletePattern("search:*"); err != nil {
		fmt.Printf("清除搜索缓存失败: %v\n", err)
	} else {
		fmt.Println("✓ 搜索结果缓存已清除")
	}

	// 3. 验证修复效果
	fmt.Println("步骤3: 验证修复效果...")
	fmt.Println("检查所有文章的点赞数一致性...")
	
	var inconsistentArticles []struct {
		ID        uint
		Title     string
		LikeCount int64
		ActualLikes int64
	}

	for _, article := range articles {
		var actualLikeCount int64
		db.Model(&models.Like{}).Where("target_type = ? AND target_id = ?", "article", article.ID).Count(&actualLikeCount)
		
		// 重新加载文章数据以确保最新
		var updatedArticle models.Article
		if err := db.First(&updatedArticle, article.ID).Error; err != nil {
			continue
		}
		
		if updatedArticle.LikeCount != actualLikeCount {
			inconsistentArticles = append(inconsistentArticles, struct {
				ID        uint
				Title     string
				LikeCount int64
				ActualLikes int64
			}{
				ID:        updatedArticle.ID,
				Title:     updatedArticle.Title,
				LikeCount: updatedArticle.LikeCount,
				ActualLikes: actualLikeCount,
			})
		}
	}

	if len(inconsistentArticles) == 0 {
		fmt.Println("✓ 验证通过：所有文章点赞数已一致")
	} else {
		fmt.Printf("⚠️  仍有 %d 篇文章点赞数不一致：\n", len(inconsistentArticles))
		for _, article := range inconsistentArticles {
			fmt.Printf("  - 文章ID %d (%s): DB=%d, 实际=%d\n", 
				article.ID, article.Title, article.LikeCount, article.ActualLikes)
		}
	}

	fmt.Println("系统性点赞数修复完成！")
	fmt.Println("请刷新前端页面以查看修复效果。")
}