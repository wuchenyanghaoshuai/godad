package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 设置数据库环境变量
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root", "123456", "127.0.0.1", "3307", "godad")

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("修复孤儿评论...")

	// 查找所有孤儿评论（parent_id指向已删除的评论）
	orphanRows, err := db.Raw(`
		SELECT c1.id, c1.article_id, c1.parent_id, c1.content, c1.user_id
		FROM comments c1
		WHERE c1.status = 1 AND c1.deleted_at IS NULL
		AND c1.parent_id IS NOT NULL AND c1.parent_id > 0
		AND c1.parent_id NOT IN (
			SELECT id FROM comments c2 
			WHERE c2.status = 1 AND c2.deleted_at IS NULL
		)
	`).Rows()
	
	if err != nil {
		log.Fatal("Failed to query orphan comments:", err)
	}
	defer orphanRows.Close()
	
	orphanCount := 0
	for orphanRows.Next() {
		var id, articleID, parentID, userID int
		var content string
		orphanRows.Scan(&id, &articleID, &parentID, &content, &userID)
		
		fmt.Printf("发现孤儿评论 ID:%d, Article:%d, Parent:%d (已删除), User:%d, Content:%s\n", 
			id, articleID, parentID, userID, content)
		
		// 方案1：将孤儿评论转换为主评论（设置parent_id为NULL）
		result := db.Exec("UPDATE comments SET parent_id = NULL WHERE id = ?", id)
		if result.Error != nil {
			log.Printf("修复孤儿评论 %d 失败: %v", id, result.Error)
		} else {
			fmt.Printf("  已将评论 %d 转换为主评论\n", id)
		}
		
		orphanCount++
	}
	
	if orphanCount == 0 {
		fmt.Println("没有发现孤儿评论")
	} else {
		fmt.Printf("修复了 %d 个孤儿评论\n", orphanCount)
		
		// 重新计算文章评论数
		fmt.Println("重新计算评论数...")
		result := db.Exec(`
			UPDATE articles a 
			JOIN (
				SELECT article_id, COUNT(*) as comment_count
				FROM comments 
				WHERE status = 1 AND deleted_at IS NULL
				GROUP BY article_id
			) c ON a.id = c.article_id
			SET a.comment_count = c.comment_count
		`)
		
		if result.Error != nil {
			log.Printf("重新计算评论数失败: %v", result.Error)
		} else {
			fmt.Printf("重新计算了 %d 篇文章的评论数\n", result.RowsAffected)
		}
	}

	// 检查文章464的最终状态
	fmt.Println("\n文章464修复后的评论:")
	finalRows, err := db.Raw(`
		SELECT id, parent_id, content, user_id 
		FROM comments 
		WHERE article_id = 464 AND status = 1 AND deleted_at IS NULL 
		ORDER BY created_at
	`).Rows()
	
	if err != nil {
		log.Fatal("Failed to query final comments:", err)
	}
	defer finalRows.Close()
	
	for finalRows.Next() {
		var id, parentID, userID int
		var content string
		finalRows.Scan(&id, &parentID, &content, &userID)
		parentStr := fmt.Sprintf("%d", parentID)
		if parentID == 0 {
			parentStr = "NULL"
		}
		fmt.Printf("  ID:%d, Parent:%s, User:%d, Content:%s\n", id, parentStr, userID, content)
	}

	var finalCommentCount int64
	db.Raw("SELECT comment_count FROM articles WHERE id = 464").Scan(&finalCommentCount)
	fmt.Printf("\n文章464的最终评论计数: %d\n", finalCommentCount)

	fmt.Println("修复完成！")
}