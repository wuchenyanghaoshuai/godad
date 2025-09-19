package routes

import (
	"godad-backend/controllers"
	"godad-backend/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterTopicRoutes 注册话题路由
func RegisterTopicRoutes(router *gin.Engine, db *gorm.DB) {
	topicController := controllers.NewTopicController(db)

	// 公共路由组
	api := router.Group("/api")
	{
		// 获取所有启用的话题（用于发帖时选择）
		api.GET("/topics/active", topicController.GetActiveTopics)
		// 获取话题详情
		api.GET("/topics/:id", topicController.GetTopic)
	}

	// 管理员路由组
	admin := router.Group("/api/admin/topics")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.AdminMiddleware())
	{
		// 获取话题列表（包括禁用的话题）
		admin.GET("", topicController.GetTopicList)
		// 创建话题
		admin.POST("", topicController.CreateTopic)
		// 更新话题
		admin.PUT("/:id", topicController.UpdateTopic)
		// 删除话题
		admin.DELETE("/:id", topicController.DeleteTopic)
		// 获取话题详情
		admin.GET("/:id", topicController.GetTopic)
	}
}