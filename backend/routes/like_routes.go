package routes

import (
	"github.com/gin-gonic/gin"
	"godad-backend/controllers"
	"godad-backend/middleware"
)

// RegisterLikeRoutes 注册点赞相关路由
func RegisterLikeRoutes(router *gin.Engine, likeController *controllers.LikeController) {
	// 点赞相关路由组
	likeGroup := router.Group("/api/likes")
	{
		// 需要认证的路由
		authGroup := likeGroup.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			// 切换点赞状态
			authGroup.POST("/toggle", likeController.ToggleLike)
			
			// 获取用户点赞列表
			authGroup.GET("/my", likeController.GetUserLikes)
			
			// 批量获取点赞状态
			authGroup.POST("/batch-status", likeController.BatchGetLikeStatus)
		}

		// 公开路由（不需要认证）
		{
			// 获取点赞状态（可选认证，如果有token则提取用户信息）
			likeGroup.GET("/status", middleware.OptionalAuthMiddleware(), likeController.GetLikeStatus)
			
			// 获取目标对象的点赞列表
			likeGroup.GET("/list", likeController.GetLikesByTarget)
			
			// 获取热门内容
			likeGroup.GET("/popular", likeController.GetPopularContent)
		}
	}
}