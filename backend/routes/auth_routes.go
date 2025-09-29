package routes

import (
	"godad-backend/config"
	"godad-backend/controllers"
	"godad-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes 设置认证相关路由
func SetupAuthRoutes(router *gin.Engine) {
	userController := controllers.NewUserController()
	passwordResetController := controllers.NewPasswordResetController(config.GetDB())

	// 创建认证路由组
	auth := router.Group("/api/auth")

	// 公开路由 - 不需要认证，但需要速率限制
	publicAuth := auth.Group("")
	publicAuth.Use(middleware.AuthRateLimit())
	{
		publicAuth.POST("/register", userController.Register)                       // 用户注册
		publicAuth.POST("/login", userController.Login)                             // 用户登录
		publicAuth.POST("/forgot-password", passwordResetController.ForgotPassword) // 忘记密码
		publicAuth.POST("/reset-password", passwordResetController.ResetPassword)   // 重置密码
		publicAuth.POST("/refresh-token", userController.RefreshToken)              // 基于刷新Cookie刷新令牌
	}

	// 需要认证的路由
	authenticated := auth.Group("")
	authenticated.Use(middleware.AuthMiddleware())
	{
		authenticated.POST("/logout", userController.Logout) // 用户登出
	}
}
