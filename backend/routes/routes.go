package routes

import (
	"godad-backend/config"
	"godad-backend/controllers"
	"godad-backend/middleware"
	"godad-backend/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
func SetupRoutes() *gin.Engine {
	// 创建Gin引擎
	router := gin.New()

	// 添加中间件
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.SecurityMiddleware())

	// 健康检查接口
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "GoDad Backend Service is running",
			"version": "1.0.0",
		})
	})

	// API版本信息
	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":        "GoDad API",
			"version":     "1.0.0",
			"description": "GoDad育儿知识分享平台API",
			"endpoints": gin.H{
				"auth":         "/api/auth",
				"user":         "/api/user",
				"article":      "/api/article",
				"comment":      "/api/comment",
				"category":     "/api/category",
				"upload":       "/api/upload",
				"follow":       "/api/follows",
				"notification": "/api/notifications",
				"favorite":     "/api/favorites",
			},
		})
	})

	// 设置各模块路由
	SetupAuthRoutes(router)     // 认证路由
	SetupUserRoutes(router)     // 用户路由
	SetupArticleRoutes(router)  // 文章路由（包含分类路由）
	SetupCommentRoutes(router)  // 评论路由
	SetupUploadRoutes(router)   // 上传路由
	SetupAdminRoutes(router)    // 管理员路由
	SetupFollowRoutes(router)   // 关注路由
	
	// 设置点赞路由
	likeService := services.NewLikeService(config.GetDB())
	likeController := controllers.NewLikeController(likeService)
	RegisterLikeRoutes(router, likeController)

	// 设置通知路由
	notificationService := services.NewNotificationService(config.GetDB())
	notificationController := controllers.NewNotificationController(notificationService)
	NotificationRoutes(router, notificationController)

	// 设置收藏路由
	favoriteService := services.NewFavoriteService(config.GetDB())
	favoriteController := controllers.NewFavoriteController(favoriteService)
	SetupFavoriteRoutes(router, favoriteController)

	return router
}