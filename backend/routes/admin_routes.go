package routes

import (
    "godad-backend/config"
    "godad-backend/controllers"
    "godad-backend/middleware"
    "godad-backend/services"

    "github.com/gin-gonic/gin"
)

// SetupAdminRoutes 设置管理员路由
func SetupAdminRoutes(router *gin.Engine) {
    adminController := controllers.NewAdminController()
    // 通知服务与控制器
    notificationController := controllers.NewAdminNotificationController(services.NewNotificationService(config.GetDB()))

	// 管理员路由组
	admin := router.Group("/api/admin")
    admin.Use(middleware.AuthMiddleware())
    admin.Use(controllers.AdminMiddleware())

    {
        // 统计数据
        admin.GET("/stats", adminController.GetStats)

		// 文章管理
		admin.GET("/articles", adminController.GetArticles)
		admin.PUT("/articles/:id/status", adminController.UpdateArticleStatus)
		admin.DELETE("/articles/:id", adminController.DeleteArticle)

        // 用户管理
        admin.GET("/users", adminController.GetUsers)
        admin.PUT("/users/:id/status", adminController.UpdateUserStatus)

        // 系统通知
        admin.POST("/notifications/system/broadcast", notificationController.BroadcastSystemNotification)
        admin.GET("/notifications/system/history", notificationController.ListSystemNotifications)
    }
}
