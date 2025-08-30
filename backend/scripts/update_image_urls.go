package main

import (
	"fmt"
	"log"

	"godad-backend/config"
	"godad-backend/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("开始更新图片URL格式...")

	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// 加载配置
	cfg := config.GetConfig()

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	oldDomain := fmt.Sprintf("https://%s.%s/", cfg.OSS.BucketName, cfg.OSS.Endpoint)
	newDomain := fmt.Sprintf("https://%s/", cfg.OSS.CustomDomain)

	fmt.Printf("将URL从 %s 更新为 %s\n", oldDomain, newDomain)

	// 更新文章封面图片
	result := db.Model(&models.Article{}).
		Where("cover_image LIKE ?", oldDomain+"%").
		Update("cover_image", gorm.Expr("REPLACE(cover_image, ?, ?)", oldDomain, newDomain))

	if result.Error != nil {
		log.Printf("更新文章封面失败: %v\n", result.Error)
	} else {
		fmt.Printf("更新了 %d 个文章封面图片URL\n", result.RowsAffected)
	}

	// 更新用户头像
	result = db.Model(&models.User{}).
		Where("avatar LIKE ?", oldDomain+"%").
		Update("avatar", gorm.Expr("REPLACE(avatar, ?, ?)", oldDomain, newDomain))

	if result.Error != nil {
		log.Printf("更新用户头像失败: %v\n", result.Error)
	} else {
		fmt.Printf("更新了 %d 个用户头像URL\n", result.RowsAffected)
	}

	fmt.Println("URL更新完成！")
}