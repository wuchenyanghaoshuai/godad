package controllers

import (
	"net/http"
	"strconv"

	"godad-backend/models"
	"godad-backend/services"
	"godad-backend/utils"

	"github.com/gin-gonic/gin"
)

// FavoriteController 收藏控制器
type FavoriteController struct {
	favoriteService *services.FavoriteService
}

// NewFavoriteController 创建收藏控制器
func NewFavoriteController(favoriteService *services.FavoriteService) *FavoriteController {
	return &FavoriteController{
		favoriteService: favoriteService,
	}
}

// ToggleFavorite 切换收藏状态
func (c *FavoriteController) ToggleFavorite(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}

	var req models.FavoriteCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	result, err := c.favoriteService.ToggleFavorite(userID.(uint), req.ArticleID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "操作失败: "+err.Error())
		return
	}

	if result == nil {
		// 取消收藏
		utils.SuccessResponse(ctx, gin.H{
			"is_favorited": false,
		}, "取消收藏成功")
	} else {
		// 添加收藏
		utils.SuccessResponse(ctx, gin.H{
			"is_favorited": true,
			"favorite":     result,
		}, "收藏成功")
	}
}

// GetFavoriteStatus 获取收藏状态
func (c *FavoriteController) GetFavoriteStatus(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}

	articleIDStr := ctx.Param("articleId")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "文章ID无效")
		return
	}

	isFavorited, err := c.favoriteService.GetFavoriteStatus(userID.(uint), uint(articleID))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "查询收藏状态失败: "+err.Error())
		return
	}

	utils.SuccessResponse(ctx, gin.H{
		"is_favorited": isFavorited,
	}, "获取收藏状态成功")
}

// BatchGetFavoriteStatus 批量获取收藏状态
func (c *FavoriteController) BatchGetFavoriteStatus(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}

	var req struct {
		ArticleIDs []uint `json:"article_ids" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	result, err := c.favoriteService.BatchGetFavoriteStatus(userID.(uint), req.ArticleIDs)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "批量查询收藏状态失败: "+err.Error())
		return
	}

	utils.SuccessResponse(ctx, result, "批量获取收藏状态成功")
}

// GetUserFavorites 获取用户收藏列表
func (c *FavoriteController) GetUserFavorites(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}

	var req models.FavoriteListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	favorites, total, err := c.favoriteService.GetUserFavorites(userID.(uint), req.Page, req.Size)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "获取收藏列表失败: "+err.Error())
		return
	}

	// 转换为响应格式
	var responses []*models.FavoriteResponse
	for _, favorite := range favorites {
		responses = append(responses, favorite.ToResponse())
	}

	// 计算分页信息
	totalPages := (int(total) + req.Size - 1) / req.Size

	utils.SuccessResponse(ctx, gin.H{
		"favorites": responses,
		"pagination": gin.H{
			"total":        total,
			"current_page": req.Page,
			"per_page":     req.Size,
			"total_pages":  totalPages,
		},
	}, "获取收藏列表成功")
}

// GetArticleFavorites 获取文章收藏列表
func (c *FavoriteController) GetArticleFavorites(ctx *gin.Context) {
	articleIDStr := ctx.Param("articleId")
	articleID, err := strconv.ParseUint(articleIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "文章ID无效")
		return
	}

	var req models.FavoriteListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	favorites, total, err := c.favoriteService.GetArticleFavorites(uint(articleID), req.Page, req.Size)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "获取文章收藏列表失败: "+err.Error())
		return
	}

	// 转换为响应格式
	var responses []*models.FavoriteResponse
	for _, favorite := range favorites {
		responses = append(responses, favorite.ToResponse())
	}

	// 计算分页信息
	totalPages := (int(total) + req.Size - 1) / req.Size

	utils.SuccessResponse(ctx, gin.H{
		"favorites": responses,
		"pagination": gin.H{
			"total":        total,
			"current_page": req.Page,
			"per_page":     req.Size,
			"total_pages":  totalPages,
		},
	}, "获取文章收藏列表成功")
}

// DeleteFavorite 删除收藏
func (c *FavoriteController) DeleteFavorite(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}

	favoriteIDStr := ctx.Param("id")
	favoriteID, err := strconv.ParseUint(favoriteIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "收藏ID无效")
		return
	}

	if err := c.favoriteService.DeleteFavorite(userID.(uint), uint(favoriteID)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "删除收藏失败: "+err.Error())
		return
	}

	utils.SuccessResponse(ctx, nil, "删除收藏成功")
}

// GetPopularFavorites 获取热门收藏
func (c *FavoriteController) GetPopularFavorites(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	daysStr := ctx.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 7
	}

	results, err := c.favoriteService.GetPopularFavorites(limit, days)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "获取热门收藏失败: "+err.Error())
		return
	}

	utils.SuccessResponse(ctx, results, "获取热门收藏成功")
}