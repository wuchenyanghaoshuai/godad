package routes

import (
	"godad-backend/controllers"
	"godad-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUploadRoutes 设置上传相关路由
func SetupUploadRoutes(router *gin.Engine) {
	// 创建上传控制器实例
	uploadController := controllers.NewUploadController()

	// API v1 路由组
	api := router.Group("/api")
	{
		// 上传相关路由 - 需要认证
		upload := api.Group("/upload", middleware.AuthMiddleware())
		{
			// 上传图片
			upload.POST("/image", uploadController.UploadImage)
			// 上传头像
			upload.POST("/avatar", uploadController.UploadAvatar)
			// 删除文件
			upload.DELETE("/:id", uploadController.DeleteFile)
			// 获取当前用户上传文件列表
			upload.GET("/my", uploadController.GetMyUploads)
			// 获取上传文件详情
			upload.GET("/:id", uploadController.GetUpload)
		}
	}
}