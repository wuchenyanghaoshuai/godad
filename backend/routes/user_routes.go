package routes

import (
	"godad-backend/controllers"
	"godad-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes 设置用户相关路由
func SetupUserRoutes(router *gin.Engine) {
	userController := controllers.NewUserController()

	// 用户API路由组
	userGroup := router.Group("/api/user")
	{
		// 公开路由（无需认证）
		userGroup.GET("/:id", userController.GetUserByID)        // 获取用户公开信息
		userGroup.GET("/:id/articles", userController.GetUserArticles) // 获取用户文章列表
		userGroup.GET("/profile/:username", userController.GetUserByUsername) // 根据用户名获取用户信息
		userGroup.GET("/profile/:username/articles", userController.GetUserArticlesByUsername) // 根据用户名获取用户文章列表
		userGroup.GET("/check-nickname", userController.CheckNickname) // 检查昵称是否可用
		userGroup.POST("/generate-nickname", userController.GenerateRandomNickname) // 生成随机昵称

		// 需要认证的路由
		authGroup := userGroup.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			authGroup.GET("/profile", userController.GetProfile)           // 获取当前用户信息
			authGroup.PUT("/profile", userController.UpdateProfile)        // 更新当前用户信息
			authGroup.POST("/change-password", userController.ChangePassword) // 修改密码
		}

		// 管理员路由
		adminGroup := userGroup.Group("")
		adminGroup.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			adminGroup.GET("/list", userController.GetUserList) // 获取用户列表（管理员）
		}
	}
}