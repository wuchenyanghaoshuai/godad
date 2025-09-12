package controllers

import (
	"strconv"

	"godad-backend/middleware"
	"godad-backend/models"
	"godad-backend/services"
	"godad-backend/utils"

	"github.com/gin-gonic/gin"
)

// CommentController 评论控制器
type CommentController struct {
	commentService *services.CommentService
}

// NewCommentController 创建评论控制器实例
func NewCommentController() *CommentController {
	return &CommentController{
		commentService: services.NewCommentService(),
	}
}

// CreateComment 创建评论
// @Summary 创建评论
// @Description 创建新评论或回复
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param comment body models.CommentCreateRequest true "评论信息"
// @Success 200 {object} utils.Response{data=models.CommentResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/comments [post]
func (c *CommentController) CreateComment(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 绑定请求参数
	var req models.CommentCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数格式错误: "+err.Error())
		return
	}

	// 创建评论
	comment, err := c.commentService.CreateComment(userID, &req)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(ctx, comment.ToResponse(false))
}

// UpdateComment 更新评论
// @Summary 更新评论
// @Description 更新评论内容
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "评论ID"
// @Param comment body models.CommentUpdateRequest true "评论信息"
// @Success 200 {object} utils.Response{data=models.CommentResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/comments/{id} [put]
func (c *CommentController) UpdateComment(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 获取评论ID
	commentID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "评论ID格式错误")
		return
	}

	// 绑定请求参数
	var req models.CommentUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数格式错误: "+err.Error())
		return
	}

	// 更新评论
	comment, err := c.commentService.UpdateComment(uint(commentID), userID, &req)
	if err != nil {
		if err.Error() == "评论不存在" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else if err.Error() == "无权限修改此评论" {
			utils.Error(ctx, utils.CodeForbidden, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	utils.Success(ctx, comment.ToResponse(false))
}

// DeleteComment 删除评论
// @Summary 删除评论
// @Description 删除评论（软删除）
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "评论ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/comments/{id} [delete]
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 获取评论ID
	commentID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "评论ID格式错误")
		return
	}

	// 删除评论
	err = c.commentService.DeleteComment(uint(commentID), userID)
	if err != nil {
		if err.Error() == "评论不存在" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else if err.Error() == "无权限删除此评论" {
			utils.Error(ctx, utils.CodeForbidden, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	utils.Success(ctx, "删除评论成功")
}

// GetComment 获取评论详情
// @Summary 获取评论详情
// @Description 根据ID获取评论详情
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param id path int true "评论ID"
// @Success 200 {object} utils.Response{data=models.CommentResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/comments/{id} [get]
func (c *CommentController) GetComment(ctx *gin.Context) {
	// 获取评论ID
	commentID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "评论ID格式错误")
		return
	}

	// 获取评论
	comment, err := c.commentService.GetComment(uint(commentID))
	if err != nil {
		if err.Error() == "评论不存在" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	utils.Success(ctx, comment.ToResponse(false))
}

// GetCommentsByArticle 获取文章的评论列表
// @Summary 获取文章的评论列表
// @Description 获取指定文章的评论列表，支持分页
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param article_id path int true "文章ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=utils.PageResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/articles/{article_id}/comments [get]
func (c *CommentController) GetCommentsByArticle(ctx *gin.Context) {
	// 获取文章ID
	articleID, err := strconv.ParseUint(ctx.Param("article_id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "文章ID格式错误")
		return
	}

	// 解析查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	sort := ctx.DefaultQuery("sort", "most_liked")

	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 10
	}

	// 验证排序参数
	validSorts := map[string]bool{
		"newest":     true,
		"oldest":     true,
		"most_liked": true,
	}
	if !validSorts[sort] {
		sort = "most_liked"
	}

	// 获取评论列表
	comments, total, err := c.commentService.GetCommentsByArticleWithSort(uint(articleID), page, size, sort)
	if err != nil {
		if err.Error() == "文章不存在或已下线" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	// 转换为响应格式
	var commentResponses []*models.CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, comment.ToResponse(true))
	}

	// 使用分页响应
	utils.SuccessPage(ctx, commentResponses, total, page, size)
}

// GetCommentReplies 获取评论的回复列表
// @Summary 获取评论的回复列表
// @Description 获取指定评论的回复列表，支持分页
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param parent_id path int true "父评论ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=utils.PageResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/comments/{parent_id}/replies [get]
func (c *CommentController) GetCommentReplies(ctx *gin.Context) {
	// 获取父评论ID
	parentID, err := strconv.ParseUint(ctx.Param("parent_id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "父评论ID格式错误")
		return
	}

	// 解析查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 10
	}

	// 获取回复列表
	replies, total, err := c.commentService.GetCommentReplies(uint(parentID), page, size)
	if err != nil {
		if err.Error() == "父评论不存在" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	// 转换为响应格式
	var replyResponses []*models.CommentResponse
	for _, reply := range replies {
		replyResponses = append(replyResponses, reply.ToResponse(false))
	}

	// 使用分页响应
	utils.SuccessPage(ctx, replyResponses, total, page, size)
}

// GetMyComments 获取当前用户的评论列表
// @Summary 获取当前用户的评论列表
// @Description 获取当前用户的评论列表，支持分页
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=utils.PageResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/comments/my [get]
func (c *CommentController) GetMyComments(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 解析查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 10
	}

	// 获取评论列表
	comments, total, err := c.commentService.GetUserComments(userID, page, size)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, err.Error())
		return
	}

	// 转换为响应格式
	var commentResponses []*models.CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, comment.ToResponse(false))
	}

	// 使用分页响应
	utils.SuccessPage(ctx, commentResponses, total, page, size)
}

// LikeComment 点赞评论
// @Summary 点赞评论
// @Description 点赞指定评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "评论ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/comments/{id}/like [post]
func (c *CommentController) LikeComment(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 获取评论ID
	commentID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "评论ID格式错误")
		return
	}

	// 点赞评论
	err = c.commentService.LikeComment(uint(commentID), userID)
	if err != nil {
		if err.Error() == "评论不存在" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	utils.Success(ctx, "点赞成功")
}

// UnlikeComment 取消点赞评论
// @Summary 取消点赞评论
// @Description 取消点赞指定评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "评论ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/comments/{id}/unlike [post]
func (c *CommentController) UnlikeComment(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 获取评论ID
	commentID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "评论ID格式错误")
		return
	}

	// 取消点赞评论
	err = c.commentService.UnlikeComment(uint(commentID), userID)
	if err != nil {
		if err.Error() == "评论不存在" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	utils.Success(ctx, "取消点赞成功")
}