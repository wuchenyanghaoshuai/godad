package routes

import (
	"godad-backend/controllers"
	"godad-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupChatRoutes(router *gin.Engine, chatController *controllers.ChatController) {
	chatGroup := router.Group("/api/chat")
	chatGroup.Use(middleware.AuthMiddleware())
	{
		// 对话相关
		chatGroup.GET("/conversations", chatController.GetConversations)
		chatGroup.POST("/conversations", chatController.GetOrCreateConversation)
		chatGroup.DELETE("/conversations/:conversationId", chatController.DeleteConversation)
		
		// 消息相关
		chatGroup.GET("/conversations/:conversationId/messages", chatController.GetMessages)
		chatGroup.POST("/messages", chatController.SendMessage)
		chatGroup.PUT("/conversations/:conversationId/read", chatController.MarkAsRead)
		
		// 表情相关
		chatGroup.GET("/emojis", chatController.GetEmojis)
	}
}