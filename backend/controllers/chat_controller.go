package controllers

import (
	"godad-backend/middleware"
	"godad-backend/models"
	"godad-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	chatService *services.ChatService
}

func NewChatController(chatService *services.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
	}
}

// GetConversations 获取用户的对话列表
func (cc *ChatController) GetConversations(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	conversations, total, err := cc.chatService.GetConversationsByUser(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取对话列表失败",
			"error":   err.Error(),
		})
		return
	}
	
	// 处理返回数据，添加对方用户信息
	type ConversationResponse struct {
		*models.ChatConversation
		OtherUser    *models.User `json:"other_user"`
		UnreadCount  uint         `json:"unread_count"`
	}
	
	var response []ConversationResponse
	for _, conv := range conversations {
		var otherUser *models.User
		var unreadCount uint
		
		if conv.User1ID == userID {
			otherUser = &conv.User2
			unreadCount = conv.User1UnreadCount
		} else {
			otherUser = &conv.User1
			unreadCount = conv.User2UnreadCount
		}
		
		response = append(response, ConversationResponse{
			ChatConversation: conv,
			OtherUser:        otherUser,
			UnreadCount:      unreadCount,
		})
	}
	
	totalPages := (int(total) + limit - 1) / limit
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"conversations": response,
			"pagination": gin.H{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": totalPages,
			},
		},
	})
}

// GetMessages 获取对话中的消息列表
func (cc *ChatController) GetMessages(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	
	conversationIDStr := c.Param("conversationId")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "对话ID无效",
		})
		return
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	
	messages, total, err := cc.chatService.GetMessagesByConversation(uint(conversationID), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	
	totalPages := (int(total) + limit - 1) / limit
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"messages": messages,
			"pagination": gin.H{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": totalPages,
			},
		},
	})
}

// SendMessage 发送消息
func (cc *ChatController) SendMessage(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	
	var req services.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数格式错误",
			"error":   err.Error(),
		})
		return
	}
	
	// 验证发送者ID
	if req.SenderID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "无权限发送消息",
		})
		return
	}
	
	// 验证消息类型和内容
	if req.MessageType == "text" && (req.Content == nil || *req.Content == "") {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文本消息内容不能为空",
		})
		return
	}
	
	if req.MessageType == "image" && (req.Images == nil || len(req.Images) == 0) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "图片消息必须包含图片信息",
		})
		return
	}
	
	if req.MessageType == "emoji" && req.EmojiID == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "表情消息必须包含表情ID",
		})
		return
	}
	
	message, err := cc.chatService.SendMessage(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "发送成功",
		"data":    message,
	})
}

// MarkAsRead 标记消息为已读
func (cc *ChatController) MarkAsRead(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	
	conversationIDStr := c.Param("conversationId")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "对话ID无效",
		})
		return
	}
	
	err = cc.chatService.MarkMessagesAsRead(uint(conversationID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "标记已读成功",
	})
}

// DeleteConversation 删除对话
func (cc *ChatController) DeleteConversation(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)
	
	conversationIDStr := c.Param("conversationId")
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "对话ID无效",
		})
		return
	}
	
	err = cc.chatService.DeleteConversation(uint(conversationID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除对话成功",
	})
}

// GetEmojis 获取表情列表
func (cc *ChatController) GetEmojis(c *gin.Context) {
	emojis, err := cc.chatService.GetEmojis()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取表情列表失败",
			"error":   err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    emojis,
	})
}

// GetOrCreateConversation 获取或创建对话
func (cc *ChatController) GetOrCreateConversation(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)

	var req struct {
		OtherUserID uint `json:"other_user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数格式错误",
			"error":   err.Error(),
		})
		return
	}

	if req.OtherUserID == userID {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不能与自己创建对话",
		})
		return
	}

	conversation, err := cc.chatService.GetOrCreateConversation(userID, req.OtherUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建对话失败",
			"error":   err.Error(),
		})
		return
	}

	// 处理返回数据，添加对方用户信息（与GetConversations保持一致）
	type ConversationResponse struct {
		*models.ChatConversation
		OtherUser    *models.User `json:"other_user"`
		UnreadCount  uint         `json:"unread_count"`
	}

	var otherUser *models.User
	var unreadCount uint

	if conversation.User1ID == userID {
		otherUser = &conversation.User2
		unreadCount = conversation.User1UnreadCount
	} else {
		otherUser = &conversation.User1
		unreadCount = conversation.User2UnreadCount
	}

	response := ConversationResponse{
		ChatConversation: conversation,
		OtherUser:        otherUser,
		UnreadCount:      unreadCount,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    response,
	})
}

// CheckMessageLimit 检查消息发送限制
func (cc *ChatController) CheckMessageLimit(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)

	var req struct {
		ReceiverID uint `json:"receiver_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数格式错误",
			"error":   err.Error(),
		})
		return
	}

	// 检查是否可以发送消息
	canSend, mutualFollow, messageCount, err := cc.chatService.CheckMessageLimit(userID, req.ReceiverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "检查消息限制失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "检查成功",
		"data": gin.H{
			"can_send":       canSend,
			"mutual_follow":  mutualFollow,
			"message_count":  messageCount,
			"daily_limit":    3,
		},
	})
}