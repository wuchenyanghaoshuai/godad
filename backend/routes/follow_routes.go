package routes

import (
	"godad-backend/config"
	"godad-backend/controllers"
	"godad-backend/middleware"
	"godad-backend/services"

	"github.com/gin-gonic/gin"
)

// SetupFollowRoutes 设置关注相关路由
func SetupFollowRoutes(router *gin.Engine) {
	// 初始化服务和控制器
	followService := services.NewFollowService(config.GetDB())
	followController := controllers.NewFollowController(followService)

	// 关注相关路由组
	followGroup := router.Group("/api/follows")
	followGroup.Use(middleware.AuthMiddleware())
	{
		// 关注用户
		followGroup.POST("/:id", followController.FollowUser)
		
		// 取消关注
		followGroup.DELETE("/:id", followController.UnfollowUser)
		
		// 检查关注状态
		followGroup.GET("/status/:id", followController.CheckFollowStatus)
		
		// 获取关注列表（我关注的人）
		followGroup.GET("/following", followController.GetFollowing)
		
		// 获取粉丝列表（关注我的人）
		followGroup.GET("/followers", followController.GetFollowers)
		
		// 获取关注统计信息
		followGroup.GET("/stats", followController.GetFollowStats)
		
		// 获取互相关注列表
		followGroup.GET("/mutual", followController.GetMutualFollows)
	}
}