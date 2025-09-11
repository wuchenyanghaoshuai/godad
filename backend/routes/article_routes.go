package routes

import (
	"godad-backend/controllers"
	"godad-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupArticleRoutes 设置文章相关路由
func SetupArticleRoutes(router *gin.Engine) {
	// 创建控制器实例
	articleController := controllers.NewArticleController()
	categoryController := controllers.NewCategoryController()

	// API v1 路由组
	v1 := router.Group("/api")

	// 公开的文章路由（无需认证）
	articlePublic := v1.Group("/articles")
	{
		// 获取文章列表
		articlePublic.GET("", articleController.GetArticleList)
		// 获取文章详情
		articlePublic.GET("/:id", articleController.GetArticle)
	}

	// 需要认证的文章路由
	articleAuth := v1.Group("/articles")
	articleAuth.Use(middleware.AuthMiddleware())
	{
		// 创建文章
		articleAuth.POST("", articleController.CreateArticle)
		// 更新文章
		articleAuth.PUT("/:id", articleController.UpdateArticle)
		// 删除文章
		articleAuth.DELETE("/:id", articleController.DeleteArticle)
		// 点赞文章
		articleAuth.POST("/:id/like", articleController.LikeArticle)
		// 获取当前用户的文章列表
		articleAuth.GET("/my", articleController.GetMyArticles)
	}

	// 公开的分类路由（无需认证）
	categoryPublic := v1.Group("/categories")
	{
		// 获取分类列表
		categoryPublic.GET("", categoryController.GetCategoryList)
		// 获取所有启用的分类
		categoryPublic.GET("/all", categoryController.GetAllCategories)
		// 获取分类及文章数量
		categoryPublic.GET("/with-count", categoryController.GetCategoriesWithCount)
		// 获取分类详情
		categoryPublic.GET("/:id", categoryController.GetCategory)
		// 根据别名获取分类详情
		categoryPublic.GET("/slug/:slug", categoryController.GetCategoryBySlug)
	}

	// 管理员文章路由
	articleAdmin := v1.Group("/admin/articles")
	articleAdmin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// 管理员可以获取所有文章（包括其他用户的）
		// 这里可以添加管理员专用的文章管理接口
	}

	// 管理员分类路由
	categoryAdmin := v1.Group("/admin/categories")
	categoryAdmin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// 获取分类列表（管理员）
		categoryAdmin.GET("", categoryController.GetCategoryList)
		// 创建分类
		categoryAdmin.POST("", categoryController.CreateCategory)
		// 更新分类
		categoryAdmin.PUT("/:id", categoryController.UpdateCategory)
		// 删除分类
		categoryAdmin.DELETE("/:id", categoryController.DeleteCategory)
		// 更新分类排序
		categoryAdmin.PUT("/:id/sort", categoryController.UpdateCategorySort)
		// 切换分类状态
		categoryAdmin.PUT("/:id/toggle", categoryController.ToggleCategoryStatus)
	}
}