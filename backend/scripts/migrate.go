// 独立的数据库迁移脚本
// 运行: go run scripts/migrate.go
package main

import (
	"flag"
	"log"

	"godad-backend/config"
	"godad-backend/models"

	"github.com/joho/godotenv"
)

func main() {
	// 命令行参数
	var action = flag.String("action", "up", "迁移操作: up(向上迁移) 或 status(查看状态)")
	var force = flag.Bool("force", false, "强制执行迁移（生产环境需要此标志）")
	flag.Parse()

	log.Println("=== GoDad 数据库迁移工具 ===")

	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Printf("警告: 无法加载.env文件: %v", err)
	}

	// 加载配置
	cfg := config.LoadConfig()

	// 生产环境安全检查
	if cfg.Server.Environment == "production" && !*force {
		log.Fatal("错误: 生产环境需要使用 --force 标志来确认迁移操作")
	}

	log.Printf("当前环境: %s", cfg.Server.Environment)
	log.Printf("数据库: %s@%s:%d/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	// 初始化数据库连接
	if err := config.InitDatabase(cfg); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer config.CloseDatabase()

	switch *action {
	case "up":
		runMigration()
	case "status":
		checkMigrationStatus()
	default:
		log.Fatalf("未知的操作: %s，支持的操作: up, status", *action)
	}
}

// runMigration 执行数据库迁移
func runMigration() {
	log.Println("开始执行数据库表结构迁移...")

	db := config.GetDB()
	if err := models.AutoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("✅ 数据库迁移完成!")
	
	// 显示迁移后的表信息
	showTableInfo()
}

// checkMigrationStatus 检查迁移状态
func checkMigrationStatus() {
	log.Println("检查数据库表状态...")
	
	db := config.GetDB()
	
	// 检查各个表是否存在
	tables := []struct {
		name  string
		model interface{}
	}{
		{"users", &models.User{}},
		{"articles", &models.Article{}},
		{"categories", &models.Category{}},
		{"comments", &models.Comment{}},
		{"uploads", &models.Upload{}},
	}

	log.Println("表名\t\t状态")
	log.Println("-------------------")
	
	for _, table := range tables {
		if db.Migrator().HasTable(table.model) {
			log.Printf("%s\t\t✅ 存在", table.name)
		} else {
			log.Printf("%s\t\t❌ 不存在", table.name)
		}
	}
}

// showTableInfo 显示表信息
func showTableInfo() {
	db := config.GetDB()
	
	log.Println("\n=== 数据库表信息 ===")
	
	// 获取表列表
	var tableNames []string
	if err := db.Raw("SHOW TABLES").Scan(&tableNames).Error; err != nil {
		log.Printf("获取表列表失败: %v", err)
		return
	}
	
	log.Printf("数据库中共有 %d 个表:", len(tableNames))
	for i, name := range tableNames {
		log.Printf("%d. %s", i+1, name)
	}
}