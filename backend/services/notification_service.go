package services

import (
    "fmt"
    "godad-backend/models"
    "log"
    "time"

    "gorm.io/gorm"
)

type NotificationService struct {
    db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

// CreateNotification 创建通知
func (s *NotificationService) CreateNotification(notification *models.Notification) error {
	// 避免自己给自己发通知
	if notification.ReceiverID == notification.ActorID {
		return nil
	}

	// 检查是否已存在相同的通知（避免重复通知）
	var existingNotification models.Notification
	err := s.db.Where(
		"receiver_id = ? AND actor_id = ? AND type = ? AND resource_id = ? AND created_at > DATE_SUB(NOW(), INTERVAL 5 MINUTE)",
		notification.ReceiverID, notification.ActorID, notification.Type, notification.ResourceID,
	).First(&existingNotification).Error

	if err == nil {
		// 如果存在相同通知，更新时间即可
		return s.db.Model(&existingNotification).Update("created_at", gorm.Expr("NOW()")).Error
	}

	if err != gorm.ErrRecordNotFound {
		return err
	}

	return s.db.Create(notification).Error
}

// CreateLikeNotification 创建点赞通知
func (s *NotificationService) CreateLikeNotification(actorID, receiverID, articleID uint) error {
	var article models.Article
	if err := s.db.First(&article, articleID).Error; err != nil {
		return err
	}

	notification := &models.Notification{
		ReceiverID: receiverID,
		ActorID:    actorID,
		Type:       models.NotificationTypeLike,
		ResourceID: articleID,
		Message:    fmt.Sprintf("点赞了你的文章《%s》", article.Title),
	}

	return s.CreateNotification(notification)
}

// CreateCommentNotification 创建评论通知
func (s *NotificationService) CreateCommentNotification(actorID, receiverID, articleID uint, commentContent string) error {
	var article models.Article
	if err := s.db.First(&article, articleID).Error; err != nil {
		return err
	}

	// 限制评论内容长度
	if len(commentContent) > 100 {
		commentContent = commentContent[:100] + "..."
	}

	notification := &models.Notification{
		ReceiverID: receiverID,
		ActorID:    actorID,
		Type:       models.NotificationTypeComment,
		ResourceID: articleID,
		Message:    fmt.Sprintf("评论了你的文章《%s》：%s", article.Title, commentContent),
	}

	return s.CreateNotification(notification)
}

// CreateCommentReplyNotification 创建回复评论通知  
func (s *NotificationService) CreateCommentReplyNotification(actorID, receiverID, articleID, commentID uint, commentContent string) error {
	var article models.Article
	if err := s.db.First(&article, articleID).Error; err != nil {
		return err
	}

	// 限制评论内容长度
	if len(commentContent) > 100 {
		commentContent = commentContent[:100] + "..."
	}

	notification := &models.Notification{
		ReceiverID: receiverID,
		ActorID:    actorID,
		Type:       models.NotificationTypeComment,
		ResourceID: articleID,
		Message:    fmt.Sprintf("回复了你的评论：%s", commentContent),
	}

	return s.CreateNotification(notification)
}

// CreateBookmarkNotification 创建收藏通知
func (s *NotificationService) CreateBookmarkNotification(actorID, receiverID, articleID uint) error {
	var article models.Article
	if err := s.db.First(&article, articleID).Error; err != nil {
		return err
	}

	notification := &models.Notification{
		ReceiverID: receiverID,
		ActorID:    actorID,
		Type:       models.NotificationTypeBookmark,
		ResourceID: articleID,
		Message:    fmt.Sprintf("收藏了你的文章《%s》", article.Title),
	}

	return s.CreateNotification(notification)
}

// CreateFollowNotification 创建关注通知
func (s *NotificationService) CreateFollowNotification(actorID, receiverID uint) error {
	var actor models.User
	if err := s.db.First(&actor, actorID).Error; err != nil {
		return err
	}

	displayName := actor.Nickname
	if displayName == "" {
		displayName = actor.Username
	}

	notification := &models.Notification{
		ReceiverID: receiverID,
		ActorID:    actorID,
		Type:       models.NotificationTypeFollow,
		Message:    fmt.Sprintf("%s 关注了你", displayName),
	}

	return s.CreateNotification(notification)
}

// CreateMessageNotification 创建私信通知
func (s *NotificationService) CreateMessageNotification(actorID, receiverID uint, messageContent string) error {
	var actor models.User
	if err := s.db.First(&actor, actorID).Error; err != nil {
		return err
	}

	displayName := actor.Nickname
	if displayName == "" {
		displayName = actor.Username
	}

	// 限制消息内容长度
	if len(messageContent) > 100 {
		messageContent = messageContent[:100] + "..."
	}

	// 确定用户ID的顺序（与对话表保持一致）
	var user1ID, user2ID uint
	if actorID < receiverID {
		user1ID, user2ID = actorID, receiverID
	} else {
		user1ID, user2ID = receiverID, actorID
	}

	// 查找这两个用户之间的对话
	var conversation models.ChatConversation
	err := s.db.Where("user1_id = ? AND user2_id = ?", user1ID, user2ID).
		First(&conversation).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	// 如果没有找到对话，说明这可能是一个无效的消息通知请求
	if err == gorm.ErrRecordNotFound || conversation.ID == 0 {
		// 记录警告日志，但不创建通知
		log.Printf("警告：尝试为不存在的对话创建消息通知 (actorID: %d, receiverID: %d)", actorID, receiverID)
		return nil
	}

	// 查找是否已存在基于这个对话的未读私信通知
	var existingNotification models.Notification
	err = s.db.Where("receiver_id = ? AND type = ? AND resource_id = ? AND is_read = false AND deleted_at IS NULL",
		receiverID, models.NotificationTypeMessage, conversation.ID).
		First(&existingNotification).Error

	if err == nil {
		// 存在未读通知，更新内容、发送者和时间
		existingNotification.ActorID = actorID
		existingNotification.Message = fmt.Sprintf("%s 给你发送了一条私信：%s", displayName, messageContent)
		existingNotification.UpdatedAt = time.Now()
		return s.db.Save(&existingNotification).Error
	} else if err == gorm.ErrRecordNotFound {
		// 不存在未读通知，创建新通知
		notification := &models.Notification{
			ReceiverID: receiverID,
			ActorID:    actorID,
			Type:       models.NotificationTypeMessage,
			ResourceID: conversation.ID,
			Message:    fmt.Sprintf("%s 给你发送了一条私信：%s", displayName, messageContent),
		}
		return s.CreateNotification(notification)
	}

	return err
}

// GetNotifications 获取通知列表
func (s *NotificationService) GetNotifications(userID uint, page, limit int) ([]models.NotificationWithDetails, int64, error) {
	var notifications []models.NotificationWithDetails
	var total int64

	offset := (page - 1) * limit

	// 统计总数
	err := s.db.Model(&models.Notification{}).
		Where("receiver_id = ? AND deleted_at IS NULL", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询通知详情
    query := `
        SELECT 
            n.id, n.receiver_id, n.actor_id, n.type, n.title, n.resource_id, n.comment_id, n.message, 
            n.is_read, n.created_at, n.updated_at,
            CASE WHEN n.type = 'system' THEN '系统' ELSE COALESCE(u.username, '') END as actor_username,
            CASE WHEN n.type = 'system' THEN '系统' ELSE COALESCE(u.nickname, '') END as actor_nickname,
            CASE WHEN n.type = 'system' THEN '' ELSE COALESCE(u.avatar, '') END as actor_avatar,
            COALESCE(a.title, '') as article_title, COALESCE(a.cover_image, '') as article_cover
        FROM notifications n
        LEFT JOIN users u ON n.actor_id = u.id
        LEFT JOIN articles a ON n.resource_id = a.id AND n.type IN ('like', 'comment', 'bookmark', 'mention')
        WHERE n.receiver_id = ? AND n.deleted_at IS NULL
        ORDER BY n.created_at DESC
        LIMIT ? OFFSET ?
    `

	err = s.db.Raw(query, userID, limit, offset).Scan(&notifications).Error
	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// GetNotificationStats 获取通知统计
func (s *NotificationService) GetNotificationStats(userID uint) (*models.NotificationStats, error) {
	var stats models.NotificationStats

	// 获取未读通知数
	err := s.db.Model(&models.Notification{}).
		Where("receiver_id = ? AND is_read = false AND deleted_at IS NULL", userID).
		Count(&stats.UnreadCount).Error
	if err != nil {
		return nil, err
	}

	// 获取总通知数
	err = s.db.Model(&models.Notification{}).
		Where("receiver_id = ? AND deleted_at IS NULL", userID).
		Count(&stats.TotalCount).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetNotificationStatsByType 获取各类型未读统计
func (s *NotificationService) GetNotificationStatsByType(userID uint) (*models.NotificationTypeStats, error) {
    type row struct {
        Type  string
        Count int64
    }

    var rows []row
    err := s.db.Model(&models.Notification{}).
        Select("type, COUNT(*) as count").
        Where("receiver_id = ? AND is_read = false AND deleted_at IS NULL", userID).
        Group("type").
        Scan(&rows).Error
    if err != nil {
        return nil, err
    }

    stats := &models.NotificationTypeStats{}
    for _, r := range rows {
        stats.TotalUnread += r.Count
        switch r.Type {
        case string(models.NotificationTypeMessage):
            stats.Message = r.Count
        case string(models.NotificationTypeLike):
            stats.Like = r.Count
        case string(models.NotificationTypeComment):
            stats.Comment = r.Count
        case string(models.NotificationTypeFollow):
            stats.Follow = r.Count
        case string(models.NotificationTypeBookmark):
            stats.Bookmark = r.Count
        case string(models.NotificationTypeSystem):
            stats.System = r.Count
        case string(models.NotificationTypeMention):
            stats.Mention = r.Count
        default:
            // 未知类型暂不计入具体分类，仅计入总数
        }
    }
    return stats, nil
}

// CreateMentionNotification 创建@提及通知
func (s *NotificationService) CreateMentionNotification(actorID, receiverID, articleID, commentID uint, commentContent string) error {
    if actorID == receiverID { return nil }
    var article models.Article
    if err := s.db.First(&article, articleID).Error; err != nil { return err }
    // 模板：在文章《标题》的评论中提到了你：{评论摘录}
    snippet := commentContent
    if len(snippet) > 100 {
        snippet = snippet[:100] + "..."
    }
    msg := fmt.Sprintf("在文章《%s》的评论中提到了你：%s", article.Title, snippet)
    n := &models.Notification{
        ReceiverID: receiverID,
        ActorID:    actorID,
        Type:       models.NotificationTypeMention,
        ResourceID: articleID,
        CommentID:  commentID,
        Message:    msg,
        Title:      "",
    }
    return s.CreateNotification(n)
}

// BroadcastSystemNotification 管理员广播系统通知到所有用户
func (s *NotificationService) BroadcastSystemNotification(adminID uint, message string) error {
    if message == "" {
        return fmt.Errorf("message cannot be empty")
    }

    // 生成一次广播的标识，避免5分钟窗口内被误判为重复通知
    // 使用当前时间戳作为 resource_id（秒级足够区分）
    broadcastID := uint(time.Now().Unix())

    // 拉取所有有效用户ID（排除管理员自己）
    var userIDs []uint
    if err := s.db.Model(&models.User{}).
        Where("status = ? AND id <> ?", 1, adminID).
        Pluck("id", &userIDs).Error; err != nil {
        return err
    }
    if len(userIDs) == 0 {
        return nil
    }

    // 分批插入，避免单次批量过大
    batchSize := 500
    now := time.Now()
    for i := 0; i < len(userIDs); i += batchSize {
        end := i + batchSize
        if end > len(userIDs) { end = len(userIDs) }

        batch := make([]models.Notification, 0, end-i)
        for _, uid := range userIDs[i:end] {
            // 避免自己给自己
            if uid == adminID { continue }
            batch = append(batch, models.Notification{
                ReceiverID: uid,
                ActorID:    adminID,
                Type:       models.NotificationTypeSystem,
                Title:      "",
                ResourceID: broadcastID,
                Message:    message,
                IsRead:     false,
                CreatedAt:  now,
                UpdatedAt:  now,
            })
        }

        if len(batch) == 0 { continue }

        if err := s.db.Create(&batch).Error; err != nil {
            return err
        }
    }

    return nil
}

// BroadcastSystemNotificationWithTitle 广播系统通知（带标题）
func (s *NotificationService) BroadcastSystemNotificationWithTitle(adminID uint, title, content string) error {
    if content == "" && title == "" {
        return fmt.Errorf("title or content required")
    }
    broadcastID := uint(time.Now().Unix())
    var userIDs []uint
    if err := s.db.Model(&models.User{}).
        Where("status = ? AND id <> ?", 1, adminID).
        Pluck("id", &userIDs).Error; err != nil {
        return err
    }
    if len(userIDs) == 0 { return nil }
    now := time.Now()
    batchSize := 500
    for i := 0; i < len(userIDs); i += batchSize {
        end := i + batchSize
        if end > len(userIDs) { end = len(userIDs) }
        batch := make([]models.Notification, 0, end-i)
        for _, uid := range userIDs[i:end] {
            if uid == adminID { continue }
            batch = append(batch, models.Notification{
                ReceiverID: uid,
                ActorID:    adminID,
                Type:       models.NotificationTypeSystem,
                Title:      title,
                ResourceID: broadcastID,
                Message:    content,
                IsRead:     false,
                CreatedAt:  now,
                UpdatedAt:  now,
            })
        }
        if len(batch) == 0 { continue }
        if err := s.db.Create(&batch).Error; err != nil { return err }
    }
    return nil
}

// GetSystemBroadcastCount 返回系统广播历史条数（按 resource_id 去重）
func (s *NotificationService) GetSystemBroadcastCount() (int64, error) {
    type res struct{ Cnt int64 }
    var r res
    err := s.db.Raw("SELECT COUNT(DISTINCT resource_id) AS cnt FROM notifications WHERE type = ? AND deleted_at IS NULL", models.NotificationTypeSystem).Scan(&r).Error
    return r.Cnt, err
}

// SystemBroadcastSummary 广播汇总
type SystemBroadcastSummary struct {
    BroadcastID uint      `json:"broadcast_id"`
    Title       string    `json:"title"`
    Message     string    `json:"message"`
    CreatedAt   time.Time `json:"created_at"`
    Total       int64     `json:"total"`
    ReadCount   int64     `json:"read_count"`
}

// ListSystemBroadcasts 列出系统广播历史（按 resource_id 分组）
func (s *NotificationService) ListSystemBroadcasts(page, limit int) ([]SystemBroadcastSummary, int64, error) {
    var total int64
    // 总组数
    err := s.db.Table("notifications").
        Where("type = ? AND deleted_at IS NULL", models.NotificationTypeSystem).
        Select("COUNT(DISTINCT resource_id) as cnt").
        Count(&total).Error
    if err != nil { return nil, 0, err }

    offset := (page - 1) * limit
    var list []SystemBroadcastSummary
    // 选取 MAX(created_at)、任意标题/消息（通常相同）、总数与已读数
    q := `
        SELECT 
            resource_id AS broadcast_id,
            COALESCE(MAX(NULLIF(title, '')), '') AS title,
            COALESCE(MAX(message), '') AS message,
            MAX(created_at) AS created_at,
            COUNT(*) AS total,
            SUM(CASE WHEN is_read = 1 THEN 1 ELSE 0 END) AS read_count
        FROM notifications
        WHERE type = 'system' AND deleted_at IS NULL
        GROUP BY resource_id
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?`

    if err := s.db.Raw(q, limit, offset).Scan(&list).Error; err != nil {
        return nil, 0, err
    }
    return list, total, nil
}

// MarkAsRead 标记通知为已读
func (s *NotificationService) MarkAsRead(userID uint, notificationIDs []uint) error {
	if len(notificationIDs) == 0 {
		return nil
	}

	return s.db.Model(&models.Notification{}).
		Where("receiver_id = ? AND id IN ?", userID, notificationIDs).
		Update("is_read", true).Error
}

// MarkAllAsRead 标记所有通知为已读
func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return s.db.Model(&models.Notification{}).
		Where("receiver_id = ? AND is_read = false", userID).
		Update("is_read", true).Error
}

// DeleteNotification 删除通知
func (s *NotificationService) DeleteNotification(userID uint, notificationID uint) error {
	return s.db.Where("receiver_id = ? AND id = ?", userID, notificationID).
		Delete(&models.Notification{}).Error
}

// DeleteAllNotifications 删除用户所有通知
func (s *NotificationService) DeleteAllNotifications(userID uint) error {
	return s.db.Where("receiver_id = ?", userID).
		Delete(&models.Notification{}).Error
}

// CleanupOldNotifications 清理旧通知（30天前的已读通知）
func (s *NotificationService) CleanupOldNotifications() error {
	result := s.db.Unscoped().
		Where("is_read = true AND created_at < DATE_SUB(NOW(), INTERVAL 30 DAY)").
		Delete(&models.Notification{})
	
	if result.Error != nil {
		return result.Error
	}

	log.Printf("Cleaned up %d old notifications", result.RowsAffected)
	return nil
}
