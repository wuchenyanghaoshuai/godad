package routes

import (
	"godad-backend/controllers"
	"godad-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupCommentRoutes 设置评论相关路由
func SetupCommentRoutes(router *gin.Engine) {
	// 创建评论控制器实例
	commentController := controllers.NewCommentController()

	// API v1 路由组
	api := router.Group("/api")
	{
		// 评论相关路由
		comments := api.Group("/comments")
		{
			// 公开路由 - 获取评论详情
			comments.GET("/:id", commentController.GetComment)
			// 公开路由 - 获取评论回复列表
			comments.GET("/replies/:parent_id", commentController.GetCommentReplies)

			// 需要认证的路由
			authComments := comments.Group("", middleware.AuthMiddleware())
			{
				// 创建评论
				authComments.POST("", commentController.CreateComment)
				// 更新评论
				authComments.PUT("/:id", commentController.UpdateComment)
				// 删除评论
				authComments.DELETE("/:id", commentController.DeleteComment)
				// 获取当前用户评论列表
				authComments.GET("/my", commentController.GetMyComments)
				// 点赞评论
				authComments.POST("/:id/like", commentController.LikeComment)
				// 取消点赞评论
				authComments.POST("/:id/unlike", commentController.UnlikeComment)
			}
		}

		// 文章评论相关路由
		articleComments := api.Group("/article-comments")
		{
			// 公开路由 - 获取文章评论列表
			articleComments.GET("/:article_id", commentController.GetCommentsByArticle)
		}
	}
}