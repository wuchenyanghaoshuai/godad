package main

import (
	"fmt"
	"log"
	"os"

	"godad-backend/config"
	"godad-backend/services"
)

func main() {
	// 设置数据库环境变量
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_PORT", "3307") 
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWORD", "123456")
	os.Setenv("DB_NAME", "godad")

	// 初始化数据库连接
	config.InitDB()

	// 创建评论服务
	commentService := services.NewCommentService()

	fmt.Println("测试评论硬删除功能...")
	fmt.Println("删除前评论状态：")

	// 删除评论42（它是评论40的回复）
	// 用户5删除自己的评论
	err := commentService.DeleteComment(42, 5)
	if err != nil {
		log.Printf("删除评论失败: %v", err)
	} else {
		fmt.Println("评论42删除成功！")
	}

	fmt.Println("测试完成")
}