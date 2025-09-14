package services

import (
	"fmt"
	"godad-backend/models"
	"log"

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

	notification := &models.Notification{
		ReceiverID: receiverID,
		ActorID:    actorID,
		Type:       models.NotificationTypeMessage,
		Message:    fmt.Sprintf("%s 给你发送了一条私信：%s", displayName, messageContent),
	}

	return s.CreateNotification(notification)
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
			n.id, n.receiver_id, n.actor_id, n.type, n.resource_id, n.message, 
			n.is_read, n.created_at, n.updated_at,
			u.username as actor_username, u.nickname as actor_nickname, u.avatar as actor_avatar,
			COALESCE(a.title, '') as article_title, COALESCE(a.cover_image, '') as article_cover
		FROM notifications n
		INNER JOIN users u ON n.actor_id = u.id
		LEFT JOIN articles a ON n.resource_id = a.id AND n.type IN ('like', 'comment', 'bookmark')
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