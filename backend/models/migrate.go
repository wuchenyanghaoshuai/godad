package models

import (
	"log"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	log.Println("开始数据库迁移...")

	// 迁移所有模型
	err := db.AutoMigrate(
		&User{},
		&Category{},
		&Article{},
		&Comment{},
		&Favorite{},
		&Upload{},
		&PasswordReset{},
		&Follow{},
		&Like{},
		&Tag{},
		&Notification{},
		&ChatConversation{},
		&ChatMessage{},
		&ChatEmoji{},
		&ChatDailyLimit{},
	)

	if err != nil {
		log.Printf("数据库迁移失败: %v", err)
		return err
	}

	// 创建索引
	if err := createIndexes(db); err != nil {
		log.Printf("创建索引失败: %v", err)
		return err
	}

	// 创建外键约束
	if err := createForeignKeys(db); err != nil {
		log.Printf("创建外键约束失败: %v", err)
		return err
	}

	// 初始化基础数据
	if err := initBaseData(db); err != nil {
		log.Printf("初始化基础数据失败: %v", err)
		return err
	}

	log.Println("数据库迁移完成")
	return nil
}

// createIndexes 创建额外的索引
func createIndexes(db *gorm.DB) error {
	// 用户表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at)")

	// 分类表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_categories_status ON categories(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_categories_sort ON categories(sort)")

	// 文章表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_articles_status ON articles(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_articles_is_top ON articles(is_top)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_articles_is_recommend ON articles(is_recommend)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_articles_published_at ON articles(published_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_articles_view_count ON articles(view_count)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_articles_like_count ON articles(like_count)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at)")

	// 评论表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_status ON comments(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at)")

	// 收藏表复合索引
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_favorites_user_article ON favorites(user_id, article_id) WHERE deleted_at IS NULL")

	// 上传表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_uploads_usage ON uploads(usage)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_uploads_file_type ON uploads(file_type)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_uploads_status ON uploads(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_uploads_created_at ON uploads(created_at)")
	
	// 删除旧的唯一索引（如果存在）
	db.Exec("DROP INDEX IF EXISTS idx_uploads_file_hash ON uploads")
	
	// 创建复合唯一索引：同一用户不能上传相同哈希的文件，但不同用户可以
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_user_file_hash ON uploads(user_id, file_hash) WHERE deleted_at IS NULL")

	// 密码重置表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_password_resets_email ON password_resets(email)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_password_resets_token ON password_resets(token)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_password_resets_expires_at ON password_resets(expires_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_password_resets_created_at ON password_resets(created_at)")

	// 关注表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_follows_follower_id ON follows(follower_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_follows_followee_id ON follows(followee_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_follows_created_at ON follows(created_at)")
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_follows_unique ON follows(follower_id, followee_id) WHERE deleted_at IS NULL")

	// 通知表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_receiver_id ON notifications(receiver_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_actor_id ON notifications(actor_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_receiver_read ON notifications(receiver_id, is_read)")

	return nil
}

// createForeignKeys 创建外键约束
func createForeignKeys(db *gorm.DB) error {
	// 文章表外键
	db.Exec("ALTER TABLE articles ADD CONSTRAINT fk_articles_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE")
	db.Exec("ALTER TABLE articles ADD CONSTRAINT fk_articles_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE ON UPDATE CASCADE")

	// 评论表外键
	db.Exec("ALTER TABLE comments ADD CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE")
	db.Exec("ALTER TABLE comments ADD CONSTRAINT fk_comments_article FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE ON UPDATE CASCADE")
	db.Exec("ALTER TABLE comments ADD CONSTRAINT fk_comments_parent FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE ON UPDATE CASCADE")
	db.Exec("ALTER TABLE comments ADD CONSTRAINT fk_comments_reply_to FOREIGN KEY (reply_to_id) REFERENCES comments(id) ON DELETE CASCADE ON UPDATE CASCADE")

	// 收藏表外键
	db.Exec("ALTER TABLE favorites ADD CONSTRAINT fk_favorites_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE")
	db.Exec("ALTER TABLE favorites ADD CONSTRAINT fk_favorites_article FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE ON UPDATE CASCADE")

	// 上传表外键
	db.Exec("ALTER TABLE uploads ADD CONSTRAINT fk_uploads_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE")

	// 关注表外键
	db.Exec("ALTER TABLE follows ADD CONSTRAINT fk_follows_follower FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE")
	db.Exec("ALTER TABLE follows ADD CONSTRAINT fk_follows_followee FOREIGN KEY (followee_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE")

	// 通知表外键
	db.Exec("ALTER TABLE notifications ADD CONSTRAINT fk_notifications_receiver FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE")
	db.Exec("ALTER TABLE notifications ADD CONSTRAINT fk_notifications_actor FOREIGN KEY (actor_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE")

	return nil
}

// initBaseData 初始化基础数据
func initBaseData(db *gorm.DB) error {
	// 检查是否已有数据
	var count int64
	db.Model(&Category{}).Count(&count)
	if count > 0 {
		log.Println("基础数据已存在，跳过初始化")
		return nil
	}

	// 创建默认分类
	categories := []Category{
		{
			Name:        "育儿知识",
			Slug:        "parenting",
			Description: "关于育儿的基础知识和经验分享",
			Icon:        "baby",
			Color:       "#FF6B6B",
			Sort:        1,
			Status:      1,
		},
		{
			Name:        "健康成长",
			Slug:        "health",
			Description: "儿童健康成长相关知识",
			Icon:        "heart",
			Color:       "#4ECDC4",
			Sort:        2,
			Status:      1,
		},
		{
			Name:        "教育启蒙",
			Slug:        "education",
			Description: "早期教育和启蒙相关内容",
			Icon:        "book",
			Color:       "#45B7D1",
			Sort:        3,
			Status:      1,
		},
		{
			Name:        "营养饮食",
			Slug:        "nutrition",
			Description: "儿童营养和饮食搭配",
			Icon:        "utensils",
			Color:       "#96CEB4",
			Sort:        4,
			Status:      1,
		},
		{
			Name:        "心理发展",
			Slug:        "psychology",
			Description: "儿童心理发展和情感培养",
			Icon:        "brain",
			Color:       "#FFEAA7",
			Sort:        5,
			Status:      1,
		},
	}

	for _, category := range categories {
		if err := db.Create(&category).Error; err != nil {
			log.Printf("创建分类失败: %v", err)
			return err
		}
	}

	// 初始化聊天表情
	var emojiCount int64
	db.Model(&ChatEmoji{}).Count(&emojiCount)
	if emojiCount == 0 {
		for _, emoji := range DefaultEmojis {
			if err := db.Create(&emoji).Error; err != nil {
				log.Printf("创建表情失败: %v", err)
				return err
			}
		}
		log.Println("聊天表情初始化完成")
	}

	log.Println("基础数据初始化完成")
	return nil
}