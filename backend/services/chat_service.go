package services

import (
	"fmt"
	"godad-backend/models"
	"time"

	"gorm.io/gorm"
)

type ChatService struct {
	db                  *gorm.DB
	notificationService *NotificationService
}

func NewChatService(db *gorm.DB, notificationService *NotificationService) *ChatService {
	return &ChatService{
		db:                  db,
		notificationService: notificationService,
	}
}

// 获取或创建对话
func (cs *ChatService) GetOrCreateConversation(user1ID, user2ID uint) (*models.ChatConversation, error) {
	// 确保用户ID顺序正确
	userID1, userID2 := models.GetConversationUsers(user1ID, user2ID)

	var conversation models.ChatConversation
	// 查找所有对话，包括已删除的
	err := cs.db.Where("user1_id = ? AND user2_id = ?", userID1, userID2).
		First(&conversation).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新对话
		conversation = models.ChatConversation{
			User1ID: userID1,
			User2ID: userID2,
		}
		err = cs.db.Create(&conversation).Error
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// 检查对话是否被删除
		if conversation.User1Deleted || conversation.User2Deleted {
			// 恢复对话：重置删除标记
			// 确定当前用户是user1还是user2，只重置对应的删除标记
			updates := make(map[string]interface{})

			if user1ID == userID1 {
				// 当前用户是user1
				updates["user1_deleted"] = false
			} else {
				// 当前用户是user2
				updates["user2_deleted"] = false
			}

			err = cs.db.Model(&conversation).Updates(updates).Error
			if err != nil {
				return nil, err
			}
		}
	}

	// 预加载用户信息
	err = cs.db.Preload("User1").Preload("User2").First(&conversation, conversation.ID).Error
	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

// 获取用户的对话列表
func (cs *ChatService) GetConversationsByUser(userID uint, page, limit int) ([]*models.ChatConversation, int64, error) {
	var conversations []*models.ChatConversation
	var total int64
	
	query := cs.db.Where("(user1_id = ? OR user2_id = ?) AND ((user1_id = ? AND user1_deleted = 0) OR (user2_id = ? AND user2_deleted = 0))",
		userID, userID, userID, userID)
	
	// 计算总数
	err := query.Model(&models.ChatConversation{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 查询数据
	err = query.Preload("User1").Preload("User2").
		Order("last_message_at DESC").
		Offset((page - 1) * limit).Limit(limit).
		Find(&conversations).Error
	if err != nil {
		return nil, 0, err
	}
	
	return conversations, total, nil
}

// 发送消息
func (cs *ChatService) SendMessage(req SendMessageRequest) (*models.ChatMessage, error) {
	// 检查互相关注关系
	mutualFollow, err := cs.checkMutualFollow(req.SenderID, req.ReceiverID)
	if err != nil {
		return nil, err
	}
	
	// 如果不是互相关注，检查每日发送限制
	if !mutualFollow {
		canSend, err := cs.checkDailyLimit(req.SenderID, req.ReceiverID)
		if err != nil {
			return nil, err
		}
		if !canSend {
			return nil, fmt.Errorf("已达到每日消息发送限制")
		}
	}
	
	// 获取或创建对话
	conversation, err := cs.GetOrCreateConversation(req.SenderID, req.ReceiverID)
	if err != nil {
		return nil, err
	}
	
	// 开始事务
	tx := cs.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// 创建消息
	message := &models.ChatMessage{
		ConversationID: conversation.ID,
		SenderID:       req.SenderID,
		ReceiverID:     req.ReceiverID,
		MessageType:    req.MessageType,
		Content:        req.Content,
		EmojiID:        req.EmojiID,
	}
	
	// 设置图片信息
	if req.Images != nil && len(req.Images) > 0 {
		// 转换类型
		images := make([]models.ImageInfo, len(req.Images))
		for i, img := range req.Images {
			images[i] = *img
		}
		err = message.SetImages(images)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	
	err = tx.Create(message).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	
	// 更新对话信息
	err = cs.updateConversationAfterMessage(tx, conversation, message)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	
	// 如果不是互相关注，更新每日限制计数
	if !mutualFollow {
		err = cs.updateDailyLimit(tx, req.SenderID, req.ReceiverID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	
	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}
	
	// 加载关联数据
	cs.db.Preload("Sender").Preload("Receiver").Preload("Emoji").First(message, message.ID)
	
	// 创建私信通知
	if cs.notificationService != nil && message.Content != nil {
		go func() {
			err := cs.notificationService.CreateMessageNotification(req.SenderID, req.ReceiverID, *message.Content)
			if err != nil {
				// 记录错误但不影响主流程
				fmt.Printf("Failed to create message notification: %v\n", err)
			}
		}()
	}
	
	return message, nil
}

// 获取对话中的消息列表
func (cs *ChatService) GetMessagesByConversation(conversationID uint, userID uint, page, limit int) ([]*models.ChatMessage, int64, error) {
	var messages []*models.ChatMessage
	var total int64
	
	// 验证用户是否属于这个对话
	var conversation models.ChatConversation
	err := cs.db.Where("id = ? AND (user1_id = ? OR user2_id = ?)", conversationID, userID, userID).
		First(&conversation).Error
	if err != nil {
		return nil, 0, fmt.Errorf("对话不存在或无权限访问")
	}
	
	// 构建查询条件：用户可以看到未被自己删除的消息
	query := cs.db.Where("conversation_id = ?", conversationID)
	
	// 添加删除条件
	query = query.Where("(sender_id = ? AND is_deleted_by_sender = false) OR (receiver_id = ? AND is_deleted_by_receiver = false)", userID, userID)
	
	// 计算总数
	err = query.Model(&models.ChatMessage{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 查询消息 - 先按时间倒序获取最新消息，然后再正序排列
	err = query.Preload("Sender").Preload("Receiver").Preload("Emoji").
		Order("created_at DESC").
		Offset((page - 1) * limit).Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}

	// 将消息重新按时间正序排列，确保UI中旧消息在上，新消息在下
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	
	return messages, total, nil
}

// 标记消息为已读
func (cs *ChatService) MarkMessagesAsRead(conversationID uint, userID uint) error {
	// 验证用户权限
	var conversation models.ChatConversation
	err := cs.db.Where("id = ? AND (user1_id = ? OR user2_id = ?)", conversationID, userID, userID).
		First(&conversation).Error
	if err != nil {
		return fmt.Errorf("对话不存在或无权限访问")
	}
	
	tx := cs.db.Begin()
	
	// 更新消息为已读
	now := time.Now()
	err = tx.Model(&models.ChatMessage{}).
		Where("conversation_id = ? AND receiver_id = ? AND is_read = false", conversationID, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
	
	if err != nil {
		tx.Rollback()
		return err
	}
	
	// 更新对话的未读计数
	if conversation.User1ID == userID {
		err = tx.Model(&conversation).Update("user1_unread_count", 0).Error
	} else {
		err = tx.Model(&conversation).Update("user2_unread_count", 0).Error
	}
	
	if err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit().Error
}

// 删除对话（软删除）
func (cs *ChatService) DeleteConversation(conversationID uint, userID uint) error {
	var conversation models.ChatConversation
	err := cs.db.Where("id = ? AND (user1_id = ? OR user2_id = ?)", conversationID, userID, userID).
		First(&conversation).Error
	if err != nil {
		return fmt.Errorf("对话不存在或无权限访问")
	}
	
	// 开始事务
	tx := cs.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// 标记用户删除对话
	updates := make(map[string]interface{})
	var otherUserID uint
	if conversation.User1ID == userID {
		updates["user1_deleted"] = true
		otherUserID = conversation.User2ID
	} else {
		updates["user2_deleted"] = true
		otherUserID = conversation.User1ID
	}
	
	// 更新对话状态
	err = tx.Model(&conversation).Updates(updates).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	
	// 删除当前用户接收到的来自对方用户的私信通知
	err = tx.Where("receiver_id = ? AND actor_id = ? AND type = ?", userID, otherUserID, "message").
		Delete(&models.Notification{}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("删除私信通知失败: %v", err)
	}
	
	return tx.Commit().Error
}

// 获取表情列表
func (cs *ChatService) GetEmojis() ([]*models.ChatEmoji, error) {
	var emojis []*models.ChatEmoji
	err := cs.db.Where("is_active = true").Order("sort_order ASC, created_at ASC").Find(&emojis).Error
	return emojis, err
}

// 辅助函数：检查互相关注关系
func (cs *ChatService) checkMutualFollow(userID1, userID2 uint) (bool, error) {
	var count int64
	
	// 检查是否互相关注
	err := cs.db.Model(&models.Follow{}).
		Where("((follower_id = ? AND followee_id = ?) OR (follower_id = ? AND followee_id = ?)) AND deleted_at IS NULL",
			userID1, userID2, userID2, userID1).
		Count(&count).Error
	
	if err != nil {
		return false, err
	}
	
	return count == 2, nil // 需要双方都关注才算互相关注
}

// 辅助函数：检查每日发送限制
func (cs *ChatService) checkDailyLimit(senderID, receiverID uint) (bool, error) {
	today := models.GetTodayDate()
	
	var limit models.ChatDailyLimit
	err := cs.db.Where("sender_id = ? AND receiver_id = ? AND date = ?", senderID, receiverID, today).
		First(&limit).Error
	
	if err == gorm.ErrRecordNotFound {
		// 没有记录，可以发送
		return true, nil
	} else if err != nil {
		return false, err
	}
	
	return limit.CanSendMessage(3), nil  // 默认每日限制3条消息
}

// 辅助函数：更新每日限制计数
func (cs *ChatService) updateDailyLimit(tx *gorm.DB, senderID, receiverID uint) error {
	today := models.GetTodayDate()
	
	var limit models.ChatDailyLimit
	err := tx.Where("sender_id = ? AND receiver_id = ? AND date = ?", senderID, receiverID, today).
		First(&limit).Error
	
	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		limit = models.ChatDailyLimit{
			SenderID:     senderID,
			ReceiverID:   receiverID,
			Date:         today,
			MessageCount: 1,
		}
		return tx.Create(&limit).Error
	} else if err != nil {
		return err
	} else {
		// 更新计数
		limit.IncrementMessageCount()
		return tx.Save(&limit).Error
	}
}

// 辅助函数：更新对话信息
func (cs *ChatService) updateConversationAfterMessage(tx *gorm.DB, conversation *models.ChatConversation, message *models.ChatMessage) error {
	// 获取消息预览内容
	content := message.GetDisplayContent()
	now := time.Now()
	
	updates := map[string]interface{}{
		"last_message_id":      message.ID,
		"last_message_content": content,
		"last_message_type":    message.MessageType,
		"last_message_at":      now,
	}
	
	// 增加接收方的未读计数
	if conversation.User1ID == message.ReceiverID {
		updates["user1_unread_count"] = gorm.Expr("user1_unread_count + 1")
	} else {
		updates["user2_unread_count"] = gorm.Expr("user2_unread_count + 1")
	}
	
	return tx.Model(conversation).Updates(updates).Error
}

// CheckMessageLimit 检查消息发送限制
func (cs *ChatService) CheckMessageLimit(senderID, receiverID uint) (bool, bool, int, error) {
	// 检查互相关注关系
	mutualFollow, err := cs.checkMutualFollow(senderID, receiverID)
	if err != nil {
		return false, false, 0, err
	}

	// 如果是互相关注，没有限制
	if mutualFollow {
		return true, true, 0, nil
	}

	// 获取今日发送的消息数量
	today := models.GetTodayDate()
	var limit models.ChatDailyLimit
	err = cs.db.Where("sender_id = ? AND receiver_id = ? AND date = ?", senderID, receiverID, today).
		First(&limit).Error

	messageCount := 0
	if err == nil {
		messageCount = int(limit.MessageCount)
	} else if err != gorm.ErrRecordNotFound {
		return false, false, 0, err
	}

	// 检查是否还能发送
	canSend := messageCount < 3

	return canSend, false, messageCount, nil
}

// 请求结构体
type SendMessageRequest struct {
	SenderID    uint                     `json:"sender_id"`
	ReceiverID  uint                     `json:"receiver_id"`
	MessageType string                   `json:"message_type"`
	Content     *string                  `json:"content,omitempty"`
	Images      []*models.ImageInfo      `json:"images,omitempty"`
	EmojiID     *uint                    `json:"emoji_id,omitempty"`
}