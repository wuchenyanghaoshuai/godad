package controllers

import (
	"strconv"

	"godad-backend/middleware"
	"godad-backend/models"
	"godad-backend/services"
	"godad-backend/utils"

	"github.com/gin-gonic/gin"
)

// ArticleController 文章控制器
type ArticleController struct {
	articleService *services.ArticleService
}

// NewArticleController 创建文章控制器实例
func NewArticleController() *ArticleController {
	return &ArticleController{
		articleService: services.NewArticleService(),
	}
}

// CreateArticle 创建文章
func (c *ArticleController) CreateArticle(ctx *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "未登录")
		return
	}

	var req models.ArticleCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	article, err := c.articleService.CreateArticle(userID, &req)
	if err != nil {
		utils.Error(ctx, utils.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, article.ToResponse(true))
}

// UpdateArticle 更新文章
func (c *ArticleController) UpdateArticle(ctx *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "未登录")
		return
	}

	articleID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "文章ID格式错误")
		return
	}

	var req models.ArticleUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	article, err := c.articleService.UpdateArticle(uint(articleID), userID, &req)
	if err != nil {
		utils.Error(ctx, utils.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, article.ToResponse(true))
}

// DeleteArticle 删除文章
// @Summary 删除文章
// @Description 删除文章（软删除）
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "文章ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/articles/{id} [delete]
func (c *ArticleController) DeleteArticle(ctx *gin.Context) {
	// 获取文章ID
	articleID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "文章ID格式错误")
		return
	}

	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 删除文章
	err = c.articleService.DeleteArticle(uint(articleID), userID)
	if err != nil {
		if err.Error() == "无权限删除此文章" {
			utils.Error(ctx, utils.CodeForbidden, err.Error())
		} else if err.Error() == "文章不存在" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	utils.Success(ctx, nil)
}

// GetArticle 获取文章详情
// @Summary 获取文章详情
// @Description 根据ID获取文章详情
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} utils.Response{data=models.ArticleResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/articles/{id} [get]
func (c *ArticleController) GetArticle(ctx *gin.Context) {
	// 获取文章ID
	articleID, err := utils.ParseUintParam(ctx, "id")
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "文章ID格式错误")
		return
	}

	// 获取当前用户ID（可选）
	userID, _ := middleware.GetCurrentUserID(ctx)

	// 获取文章 - 允许用户查看自己的文章（包括草稿）
	article, err := c.articleService.GetArticleByID(articleID, userID, true)
	if err != nil {
		if err.Error() == "文章不存在" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	utils.Success(ctx, article.ToResponse(true))
}

// GetArticleList 获取文章列表
func (c *ArticleController) GetArticleList(ctx *gin.Context) {
	var req models.ArticleListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	// 解析分页参数
	req.Page, req.Size = utils.ParsePaginationParams(ctx)

	articles, total, err := c.articleService.GetArticleList(&req)
	if err != nil {
		utils.Error(ctx, utils.CodeInternalError, err.Error())
		return
	}

	// 转换为响应格式
	var responses []*models.ArticleResponse
	for _, article := range articles {
		responses = append(responses, article.ToResponse(false))
	}

	// 使用分页响应工具
	utils.RespondWithPagination(ctx, responses, total, req.Page, req.Size)
}

// GetMyArticles 获取我的文章列表
func (c *ArticleController) GetMyArticles(ctx *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "未登录")
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	status := ctx.Query("status")

	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	articles, total, err := c.articleService.GetUserArticles(userID, page, size, status)
	if err != nil {
		utils.Error(ctx, utils.CodeInternalError, err.Error())
		return
	}

	// 转换为响应格式
	var responses []*models.ArticleResponse
	for _, article := range articles {
		responses = append(responses, article.ToResponse(false))
	}

	// 使用分页响应
	utils.SuccessPage(ctx, responses, total, page, size)
}

// LikeArticle 点赞文章
// @Summary 点赞文章
// @Description 对文章进行点赞
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "文章ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/articles/{id}/like [post]
func (c *ArticleController) LikeArticle(ctx *gin.Context) {
	// 获取文章ID
	articleID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "文章ID格式错误")
		return
	}

	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 点赞文章
	err = c.articleService.LikeArticle(uint(articleID), userID)
	if err != nil {
		if err.Error() == "文章不存在" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	utils.Success(ctx, nil)
}

// GetHotArticles 获取热门文章
// @Summary 获取热门文章
// @Description 根据时间周期获取热门文章排行榜
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param period query string false "时间周期：today(今日)、week(本周)、month(本月)、all(全部)" default(today)
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} utils.Response{data=[]models.ArticleResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/articles/hot [get]
func (c *ArticleController) GetHotArticles(ctx *gin.Context) {
	// 获取参数
	period := ctx.DefaultQuery("period", "today")
	limitStr := ctx.DefaultQuery("limit", "10")

	// 验证时间周期参数
	validPeriods := map[string]bool{
		"today": true,
		"week":  true,
		"month": true,
		"all":   true,
	}
	if !validPeriods[period] {
		utils.Error(ctx, utils.CodeBadRequest, "时间周期参数无效，支持：today, week, month, all")
		return
	}

	// 解析限制数量
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10
	}

	// 获取热门文章
	articles, err := c.articleService.GetHotArticles(period, limit)
	if err != nil {
		utils.Error(ctx, utils.CodeInternalError, err.Error())
		return
	}

	// 转换为响应格式
	var responses []*models.ArticleResponse
	for _, article := range articles {
		responses = append(responses, article.ToResponse(false))
	}

	utils.Success(ctx, responses)
}