package routes

import (
	"godad-backend/controllers"
	"godad-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupPointsRoutes(router *gin.Engine, pointsController *controllers.PointsController) {
	// 公共配置信息（无需认证）
	publicPointsGroup := router.Group("/api/points")
	{
		publicPointsGroup.GET("/levels", pointsController.GetLevels)
		publicPointsGroup.GET("/rules", pointsController.GetPointsRules)
	}

	// 用户积分相关（需要认证）
	authPointsGroup := router.Group("/api/points")
	authPointsGroup.Use(middleware.AuthMiddleware())
	{
		authPointsGroup.GET("/user", pointsController.GetUserPoints)
		authPointsGroup.GET("/history", pointsController.GetPointsHistory)
		authPointsGroup.GET("/stats", pointsController.GetPointsStats)
	}

	// 管理员功能
	adminPointsGroup := router.Group("/api/admin/points")
	adminPointsGroup.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		adminPointsGroup.POST("/award", pointsController.AwardPoints)
	}
}