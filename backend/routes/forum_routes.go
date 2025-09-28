package routes

import (
	"godad-backend/controllers"
	"godad-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupForumRoutes 设置论坛相关路由
func SetupForumRoutes(router *gin.Engine) {
	// 创建控制器实例
	forumController := controllers.NewForumController()

	// API v1 路由组
	v1 := router.Group("/api")

	// 公开的论坛路由（无需认证）
	forumPublic := v1.Group("/forum")
	{
		// 获取帖子列表
		forumPublic.GET("/posts", forumController.GetPostList)
		// 获取帖子详情
		forumPublic.GET("/posts/:id", forumController.GetPost)
		// 获取帖子回复列表
		forumPublic.GET("/posts/:id/replies", forumController.GetPostReplies)
		// 获取话题列表
		forumPublic.GET("/topics", forumController.GetTopics)
		// 获取热门帖子
		forumPublic.GET("/posts/hot", forumController.GetHotPosts)
	}

	// 需要认证的论坛路由
	forumAuth := v1.Group("/forum")
	forumAuth.Use(middleware.AuthMiddleware())
	{
		// 创建帖子
		forumAuth.POST("/posts", forumController.CreatePost)
		// 更新帖子
		forumAuth.PUT("/posts/:id", forumController.UpdatePost)
		// 删除帖子
		forumAuth.DELETE("/posts/:id", forumController.DeletePost)
		// 点赞帖子
		forumAuth.POST("/posts/:id/like", forumController.LikePost)
		// 增加浏览量
		forumAuth.POST("/posts/:id/view", forumController.IncrementPostView)

		// 回复相关
		forumAuth.POST("/replies", forumController.CreateReply)
		forumAuth.PUT("/replies/:id", forumController.UpdateReply)
		forumAuth.DELETE("/replies/:id", forumController.DeleteReply)
		forumAuth.POST("/replies/:id/like", forumController.LikeReply)

		// 获取当前用户的帖子列表
		forumAuth.GET("/posts/my", forumController.GetMyPosts)
		// 获取当前用户的回复列表
		forumAuth.GET("/replies/my", forumController.GetMyReplies)
	}

	// 管理员路由
	forumAdmin := v1.Group("/admin/forum")
	forumAdmin.Use(middleware.AuthMiddleware())
	forumAdmin.Use(middleware.AdminMiddleware())
	{
		// 论坛统计（管理员）
		forumAdmin.GET("/stats", forumController.AdminGetStats)
		// 管理员获取帖子列表（包含所有状态）
		forumAdmin.GET("/posts", forumController.AdminGetPostList)
		// 置顶/取消置顶帖子
		forumAdmin.PUT("/posts/:id/top", forumController.TogglePostTop)
		// 标记/取消标记热门帖子
		forumAdmin.PUT("/posts/:id/hot", forumController.TogglePostHot)
		// 锁定/解锁帖子
		forumAdmin.PUT("/posts/:id/lock", forumController.TogglePostLock)
		// 批量删除帖子
		forumAdmin.DELETE("/posts/batch", forumController.BatchDeletePosts)
		// 批量删除回复
		forumAdmin.DELETE("/replies/batch", forumController.BatchDeleteReplies)
		// 删除单个帖子
		forumAdmin.DELETE("/posts/:id", forumController.AdminDeletePost)
	}
}
