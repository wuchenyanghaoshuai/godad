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

	fmt.Println("修复评论计数...")

	// 修复文章评论数 - 使用JOIN方式避免MySQL限制
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
		log.Fatal("Failed to fix article comment counts:", result.Error)
	}
	
	fmt.Printf("修复了 %d 篇文章的评论计数\n", result.RowsAffected)

	// 修复评论的回复数 - 使用JOIN方式避免MySQL限制
	result = db.Exec(`
		UPDATE comments c1 
		JOIN (
			SELECT parent_id, COUNT(*) as reply_count
			FROM comments 
			WHERE parent_id IS NOT NULL AND parent_id > 0 
			AND status = 1 AND deleted_at IS NULL
			GROUP BY parent_id
		) c2 ON c1.id = c2.parent_id
		SET c1.reply_count = c2.reply_count
		WHERE c1.status = 1 AND c1.deleted_at IS NULL
	`)
	
	if result.Error != nil {
		log.Fatal("Failed to fix comment reply counts:", result.Error)
	}
	
	fmt.Printf("修复了 %d 条评论的回复计数\n", result.RowsAffected)

	// 将没有回复的评论的reply_count设为0
	result = db.Exec(`
		UPDATE comments 
		SET reply_count = 0 
		WHERE status = 1 AND deleted_at IS NULL 
		AND id NOT IN (
			SELECT DISTINCT parent_id 
			FROM (SELECT parent_id FROM comments WHERE parent_id IS NOT NULL AND parent_id > 0 AND status = 1 AND deleted_at IS NULL) t
		)
	`)
	
	if result.Error != nil {
		log.Fatal("Failed to reset reply counts:", result.Error)
	}
	
	fmt.Printf("重置了 %d 条评论的回复计数为0\n", result.RowsAffected)

	// 检查文章464的最终计数
	var commentCount int64
	db.Raw("SELECT comment_count FROM articles WHERE id = 464").Scan(&commentCount)
	
	fmt.Printf("文章464的最终评论计数: %d\n", commentCount)
	
	// 显示有效评论的详细信息
	fmt.Println("\n文章464的有效评论:")
	rows, err := db.Raw(`
		SELECT id, parent_id, content, user_id 
		FROM comments 
		WHERE article_id = 464 AND status = 1 AND deleted_at IS NULL 
		ORDER BY created_at
	`).Rows()
	
	if err != nil {
		log.Fatal("Failed to query comments:", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var id, parentID, userID int
		var content string
		rows.Scan(&id, &parentID, &content, &userID)
		parentStr := fmt.Sprintf("%d", parentID)
		if parentID == 0 {
			parentStr = "NULL"
		}
		fmt.Printf("  ID:%d, Parent:%s, User:%d, Content:%s\n", id, parentStr, userID, content)
	}

	// 检查是否有孤儿评论（parent_id指向已删除的评论）
	fmt.Println("\n检查孤儿评论（parent_id指向已删除评论）:")
	orphanRows, err := db.Raw(`
		SELECT c1.id, c1.parent_id, c1.content, c1.user_id
		FROM comments c1
		WHERE c1.article_id = 464 
		AND c1.status = 1 AND c1.deleted_at IS NULL
		AND c1.parent_id IS NOT NULL AND c1.parent_id > 0
		AND c1.parent_id NOT IN (
			SELECT id FROM comments c2 
			WHERE c2.status = 1 AND c2.deleted_at IS NULL
		)
		ORDER BY c1.created_at
	`).Rows()
	
	if err != nil {
		log.Fatal("Failed to query orphan comments:", err)
	}
	defer orphanRows.Close()
	
	orphanCount := 0
	for orphanRows.Next() {
		var id, parentID, userID int
		var content string
		orphanRows.Scan(&id, &parentID, &content, &userID)
		fmt.Printf("  孤儿评论 ID:%d, Parent:%d (已删除), User:%d, Content:%s\n", id, parentID, userID, content)
		orphanCount++
	}
	
	if orphanCount == 0 {
		fmt.Println("  没有发现孤儿评论")
	} else {
		fmt.Printf("  发现 %d 个孤儿评论\n", orphanCount)
	}

	fmt.Println("修复完成！")
}