package routes

import (
	"godad-backend/controllers"
	"godad-backend/middleware"

	"github.com/gin-gonic/gin"
)

func NotificationRoutes(router *gin.Engine, notificationController *controllers.NotificationController) {
    // 需要认证的通知路由
    auth := router.Group("/api/notifications")
    auth.Use(middleware.AuthMiddleware())
    {
        // 获取通知列表
        auth.GET("", notificationController.GetNotifications)
        
        // 获取通知统计
        auth.GET("/stats", notificationController.GetNotificationStats)

        // SSE 实时流
        auth.GET("/stream", notificationController.Stream)
        
        // 标记通知为已读
        auth.PUT("/mark-read", notificationController.MarkAsRead)
		
		// 标记所有通知为已读
		auth.PUT("/mark-all-read", notificationController.MarkAllAsRead)
		
		// 标记单个通知为已读（通过URL参数）
		auth.PUT("/:id/mark-read", notificationController.MarkAsReadByURL)
		
		// 批量标记通知为已读（通过URL参数）
		auth.PUT("/batch-mark-read", notificationController.BatchMarkAsRead)
		
		// 删除所有通知（必须在 /:id 之前注册）
		auth.DELETE("/all", notificationController.DeleteAllNotifications)
		
		// 删除通知
		auth.DELETE("/:id", notificationController.DeleteNotification)
	}
}
