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
		&Notification{},
		&ChatConversation{},
		&ChatMessage{},
		&ChatEmoji{},
		&ChatDailyLimit{},
		&ForumPost{},
		&ForumReply{},
		&Topic{},
		&Resource{},
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

	// 确保枚举类型字段更新（例如：notifications.type 增加 system）
	if err := ensureEnumColumns(db); err != nil {
		log.Printf("更新枚举字段失败: %v", err)
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

// ensureEnumColumns 确保枚举字段包含最新取值
func ensureEnumColumns(db *gorm.DB) error {
    // MySQL: 扩展 notifications.type 枚举，加入 'system'
    db.Exec("ALTER TABLE notifications MODIFY COLUMN type ENUM('like','comment','bookmark','follow','message','system','mention') NOT NULL")
    // 新增标题列（如果不存在）
    db.Exec("ALTER TABLE notifications ADD COLUMN IF NOT EXISTS title VARCHAR(255) NULL AFTER type")
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

	// 论坛帖子表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_forum_posts_author_id ON forum_posts(author_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_forum_posts_topic ON forum_posts(topic)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_forum_posts_status ON forum_posts(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_forum_posts_is_top ON forum_posts(is_top)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_forum_posts_is_hot ON forum_posts(is_hot)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_forum_posts_created_at ON forum_posts(created_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_forum_posts_last_reply_at ON forum_posts(last_reply_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_forum_posts_view_count ON forum_posts(view_count)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_forum_posts_reply_count ON forum_posts(reply_count)")

	// 论坛回复表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_forum_replies_post_id ON forum_replies(post_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_forum_replies_author_id ON forum_replies(author_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_forum_replies_parent_id ON forum_replies(parent_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_forum_replies_status ON forum_replies(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_forum_replies_created_at ON forum_replies(created_at)")

	// 话题表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_topics_is_active ON topics(is_active)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_topics_sort ON topics(sort)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_topics_created_at ON topics(created_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_topics_display_name ON topics(display_name)")

	// 资源表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_status ON resources(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_type ON resources(type)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_category ON resources(category)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_uploader_id ON resources(uploader_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_download_count ON resources(download_count)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_created_at ON resources(created_at)")

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

	// 论坛帖子表外键
	db.Exec("ALTER TABLE forum_posts ADD CONSTRAINT fk_forum_posts_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE")

	// 论坛回复表外键
	db.Exec("ALTER TABLE forum_replies ADD CONSTRAINT fk_forum_replies_post FOREIGN KEY (post_id) REFERENCES forum_posts(id) ON DELETE CASCADE ON UPDATE CASCADE")
	db.Exec("ALTER TABLE forum_replies ADD CONSTRAINT fk_forum_replies_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE")
	db.Exec("ALTER TABLE forum_replies ADD CONSTRAINT fk_forum_replies_parent FOREIGN KEY (parent_id) REFERENCES forum_replies(id) ON DELETE CASCADE ON UPDATE CASCADE")

	// 资源表外键
	db.Exec("ALTER TABLE resources ADD CONSTRAINT fk_resources_uploader FOREIGN KEY (uploader_id) REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE")

	return nil
}

// initBaseData 初始化基础数据
func initBaseData(db *gorm.DB) error {
	// 检查分类数据
	var categoryCount int64
	db.Model(&Category{}).Count(&categoryCount)
	if categoryCount > 0 {
		log.Println("分类数据已存在，跳过分类初始化")
	} else {
		// 初始化分类数据
		if err := initCategories(db); err != nil {
			return err
		}
	}

	// 检查话题数据
	var topicCount int64
	db.Model(&Topic{}).Count(&topicCount)
	if topicCount > 0 {
		log.Println("话题数据已存在，跳过话题初始化")
	} else {
		// 初始化话题数据
		if err := initTopics(db); err != nil {
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

// initCategories 初始化分类数据
func initCategories(db *gorm.DB) error {

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

	log.Println("分类数据初始化完成")
	return nil
}

// initTopics 初始化话题数据
func initTopics(db *gorm.DB) error {
	topics := []Topic{
		{
			Name:        "Baby Care",
			DisplayName: "婴儿护理",
			Description: "新生儿和婴儿护理相关话题",
			Color:       "#FF6B6B",
			Icon:        "baby",
			Sort:        1,
			IsActive:    true,
		},
		{
			Name:        "Feeding",
			DisplayName: "喂养",
			Description: "母乳喂养、辅食添加等话题",
			Color:       "#4ECDC4",
			Icon:        "utensils",
			Sort:        2,
			IsActive:    true,
		},
		{
			Name:        "Sleep",
			DisplayName: "睡眠",
			Description: "婴幼儿睡眠问题和建议",
			Color:       "#45B7D1",
			Icon:        "moon",
			Sort:        3,
			IsActive:    true,
		},
		{
			Name:        "Health",
			DisplayName: "健康",
			Description: "儿童健康、疫苗、疾病防护",
			Color:       "#96CEB4",
			Icon:        "heart",
			Sort:        4,
			IsActive:    true,
		},
		{
			Name:        "Development",
			DisplayName: "发育",
			Description: "儿童成长发育相关话题",
			Color:       "#FFEAA7",
			Icon:        "seedling",
			Sort:        5,
			IsActive:    true,
		},
		{
			Name:        "Activities",
			DisplayName: "活动",
			Description: "亲子活动、游戏、娱乐",
			Color:       "#DDA0DD",
			Icon:        "gamepad",
			Sort:        6,
			IsActive:    true,
		},
		{
			Name:        "Gear",
			DisplayName: "用品",
			Description: "婴幼儿用品推荐和评测",
			Color:       "#98D8C8",
			Icon:        "shopping-cart",
			Sort:        7,
			IsActive:    true,
		},
		{
			Name:        "Parenting",
			DisplayName: "育儿",
			Description: "育儿经验和心得分享",
			Color:       "#F7DC6F",
			Icon:        "users",
			Sort:        8,
			IsActive:    true,
		},
		{
			Name:        "Other",
			DisplayName: "其他",
			Description: "其他相关话题",
			Color:       "#AED6F1",
			Icon:        "comment-dots",
			Sort:        99,
			IsActive:    true,
		},
	}

	for _, topic := range topics {
		if err := db.Create(&topic).Error; err != nil {
			log.Printf("创建话题失败: %v", err)
			return err
		}
	}

	log.Println("话题数据初始化完成")
	return nil
}
