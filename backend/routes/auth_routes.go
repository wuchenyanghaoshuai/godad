package routes

import (
	"godad-backend/controllers"
	"godad-backend/middleware"
	"godad-backend/config"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes 设置认证相关路由
func SetupAuthRoutes(router *gin.Engine) {
	userController := controllers.NewUserController()
	passwordResetController := controllers.NewPasswordResetController(config.GetDB())

	// 创建认证路由组
	auth := router.Group("/api/auth")
	{
		// 公开路由 - 不需要认证
		auth.POST("/register", userController.Register)     // 用户注册
		auth.POST("/login", userController.Login)           // 用户登录
		auth.POST("/forgot-password", passwordResetController.ForgotPassword) // 忘记密码
		auth.POST("/reset-password", passwordResetController.ResetPassword)   // 重置密码
		
		// 需要认证的路由
		authenticated := auth.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			authenticated.POST("/refresh-token", userController.RefreshToken) // 刷新令牌
			authenticated.POST("/logout", userController.Logout)              // 用户登出
		}
	}
}