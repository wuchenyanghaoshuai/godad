package routes

import (
	"godad-backend/controllers"
	"godad-backend/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupResourceRoutes 设置资源相关路由
func SetupResourceRoutes(router *gin.Engine, db *gorm.DB) {
	resourceController := controllers.NewResourceController(db)

	// 公开路由 - 不需要认证
	publicGroup := router.Group("/api/resources")
	{
		// 获取已发布的资源列表（前端公开接口）
		publicGroup.GET("", resourceController.GetPublishedResources)
		// 获取单个资源详情
		publicGroup.GET("/:id", resourceController.GetResource)
		// 下载资源（增加下载次数）
		publicGroup.POST("/:id/download", resourceController.DownloadResource)
	}

	// 需要登录的路由
	authGroup := router.Group("/api/resources")
	authGroup.Use(middleware.AuthMiddleware())
	{
		// 用户可以上传资源（待审核）
		authGroup.POST("", resourceController.CreateResource)
	}

	// 管理员路由
	adminGroup := router.Group("/api/admin/resources")
	adminGroup.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// 获取所有资源列表（包括待审核）
		adminGroup.GET("", resourceController.GetResources)
		// 获取单个资源详情
		adminGroup.GET("/:id", resourceController.GetResource)
		// 创建资源
		adminGroup.POST("", resourceController.CreateResource)
		// 更新资源
		adminGroup.PUT("/:id", resourceController.UpdateResource)
		// 更新资源状态（审核）
		adminGroup.PUT("/:id/status", resourceController.UpdateResourceStatus)
		// 删除资源
		adminGroup.DELETE("/:id", resourceController.DeleteResource)
		// 获取资源统计信息
		adminGroup.GET("/stats", resourceController.GetResourceStats)
	}
}