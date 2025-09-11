package routes

import (
	"github.com/gin-gonic/gin"
	"godad-backend/controllers"
	"godad-backend/middleware"
)

// RegisterTagRoutes 注册标签相关路由
func RegisterTagRoutes(router *gin.Engine, tagController *controllers.TagController) {
	// 标签相关路由组
	tagGroup := router.Group("/api/tags")
	{
		// 需要管理员或内容管理员权限的路由
		adminGroup := tagGroup.Group("")
		adminGroup.Use(middleware.AuthMiddleware())
		{
			// 创建标签
			adminGroup.POST("", tagController.CreateTag)
			
			// 更新标签
			adminGroup.PUT("/:id", tagController.UpdateTag)
			
			// 删除标签
			adminGroup.DELETE("/:id", tagController.DeleteTag)
			
			// 获取标签统计信息
			adminGroup.GET("/stats", tagController.GetTagStats)
		}

		// 公开路由（不需要认证）
		{
			// 获取标签列表
			tagGroup.GET("", tagController.GetTags)
			
			// 获取标签详情
			tagGroup.GET("/:id", tagController.GetTagByID)
			
			// 获取热门标签
			tagGroup.GET("/popular", tagController.GetPopularTags)
			
			// 获取标签下的文章
			tagGroup.GET("/:id/articles", tagController.GetArticlesByTag)
			
			// 搜索标签
			tagGroup.GET("/search", tagController.SearchTags)
		}
	}
}