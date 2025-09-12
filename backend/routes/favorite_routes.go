package routes

import (
	"godad-backend/controllers"
	"godad-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupFavoriteRoutes 设置收藏相关路由
func SetupFavoriteRoutes(router *gin.Engine, favoriteController *controllers.FavoriteController) {
	api := router.Group("/api")
	{
		// 收藏相关路由（需要登录）
		favorites := api.Group("/favorites")
		favorites.Use(middleware.AuthMiddleware())
		{
			// 切换收藏状态（收藏/取消收藏）
			favorites.POST("/toggle", favoriteController.ToggleFavorite)

			// 获取收藏状态
			favorites.GET("/status/:articleId", favoriteController.GetFavoriteStatus)

			// 批量获取收藏状态
			favorites.POST("/batch-status", favoriteController.BatchGetFavoriteStatus)

			// 获取我的收藏列表
			favorites.GET("/my", favoriteController.GetUserFavorites)

			// 删除收藏
			favorites.DELETE("/:id", favoriteController.DeleteFavorite)
		}

		// 公开的收藏相关路由
		api.GET("/favorites/article/:articleId", favoriteController.GetArticleFavorites)
		api.GET("/favorites/popular", favoriteController.GetPopularFavorites)
	}
}