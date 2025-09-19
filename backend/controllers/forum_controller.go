package controllers

import (
	"strconv"

	"godad-backend/middleware"
	"godad-backend/models"
	"godad-backend/services"
	"godad-backend/utils"

	"github.com/gin-gonic/gin"
)

// ForumController 论坛控制器
type ForumController struct {
	forumService *services.ForumService
}

// NewForumController 创建论坛控制器实例
func NewForumController() *ForumController {
	return &ForumController{
		forumService: services.NewForumService(),
	}
}

// CreatePost 创建帖子
// @Summary 创建帖子
// @Description 创建新的论坛帖子
// @Tags 论坛管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param post body models.ForumPostCreateRequest true "帖子信息"
// @Success 200 {object} utils.Response{data=models.ForumPostResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/forum/posts [post]
func (c *ForumController) CreatePost(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 绑定请求参数
	var req models.ForumPostCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数格式错误: "+err.Error())
		return
	}

	// 验证话题是否有效
	if !models.IsValidTopic(req.Topic) {
		utils.Error(ctx, utils.CodeBadRequest, "无效的话题分类")
		return
	}

	// 调用服务层创建帖子
	post, err := c.forumService.CreatePost(&req, userID)
	if err != nil {
		utils.Error(ctx, utils.CodeInternalServerError, "创建帖子失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "创建帖子成功", post.ToResponse(true))
}

// GetPostList 获取帖子列表
// @Summary 获取帖子列表
// @Description 获取论坛帖子列表，支持分页和筛选
// @Tags 论坛管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param topic query string false "话题分类"
// @Param keyword query string false "搜索关键词"
// @Param sort query string false "排序方式" default("created_at desc")
// @Success 200 {object} utils.Response{data=utils.PagedResponse{items=[]models.ForumPostResponse}}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/forum/posts [get]
func (c *ForumController) GetPostList(ctx *gin.Context) {
	// 绑定查询参数
	var req models.ForumPostListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数格式错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.Sort == "" {
		req.Sort = "created_at desc"
	}

	// 调用服务层获取帖子列表
	posts, total, err := c.forumService.GetPostList(&req)
	if err != nil {
		utils.Error(ctx, utils.CodeInternalServerError, "获取帖子列表失败: "+err.Error())
		return
	}

	// 转换为响应格式
	var postResponses []models.ForumPostResponse
	for _, post := range posts {
		postResponses = append(postResponses, *post.ToResponse(false))
	}

	// 构造分页响应
	response := utils.PagedResponse{
		Items: postResponses,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
		Pages: (total + int64(req.Size) - 1) / int64(req.Size),
	}

	utils.SuccessWithMessage(ctx, "获取帖子列表成功", response)
}

// GetPost 获取帖子详情
// @Summary 获取帖子详情
// @Description 根据ID获取论坛帖子详情
// @Tags 论坛管理
// @Accept json
// @Produce json
// @Param id path int true "帖子ID"
// @Success 200 {object} utils.Response{data=models.ForumPostResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/forum/posts/{id} [get]
func (c *ForumController) GetPost(ctx *gin.Context) {
	// 获取帖子ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "无效的帖子ID")
		return
	}

	// 调用服务层获取帖子详情
	post, err := c.forumService.GetPostByID(uint(id))
	if err != nil {
		utils.Error(ctx, utils.CodeNotFound, "帖子不存在")
		return
	}

	utils.SuccessWithMessage(ctx, "获取帖子详情成功", post.ToResponse(true))
}

// UpdatePost 更新帖子
// @Summary 更新帖子
// @Description 更新论坛帖子信息
// @Tags 论坛管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "帖子ID"
// @Param post body models.ForumPostUpdateRequest true "更新信息"
// @Success 200 {object} utils.Response{data=models.ForumPostResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/forum/posts/{id} [put]
func (c *ForumController) UpdatePost(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 获取帖子ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "无效的帖子ID")
		return
	}

	// 绑定请求参数
	var req models.ForumPostUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数格式错误: "+err.Error())
		return
	}

	// 验证话题是否有效（如果提供了）
	if req.Topic != "" && !models.IsValidTopic(req.Topic) {
		utils.Error(ctx, utils.CodeBadRequest, "无效的话题分类")
		return
	}

	// 调用服务层更新帖子
	post, err := c.forumService.UpdatePost(uint(id), &req, userID)
	if err != nil {
		if err.Error() == "帖子不存在" {
			utils.Error(ctx, utils.CodeNotFound, "帖子不存在")
		} else if err.Error() == "无权限修改此帖子" {
			utils.Error(ctx, utils.CodeForbidden, "无权限修改此帖子")
		} else {
			utils.Error(ctx, utils.CodeInternalServerError, "更新帖子失败: "+err.Error())
		}
		return
	}

	utils.SuccessWithMessage(ctx, "更新帖子成功", post.ToResponse(true))
}

// DeletePost 删除帖子
// @Summary 删除帖子
// @Description 删除论坛帖子
// @Tags 论坛管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "帖子ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/forum/posts/{id} [delete]
func (c *ForumController) DeletePost(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 获取帖子ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "无效的帖子ID")
		return
	}

	// 调用服务层删除帖子
	err = c.forumService.DeletePost(uint(id), userID)
	if err != nil {
		if err.Error() == "帖子不存在" {
			utils.Error(ctx, utils.CodeNotFound, "帖子不存在")
		} else if err.Error() == "无权限删除此帖子" {
			utils.Error(ctx, utils.CodeForbidden, "无权限删除此帖子")
		} else {
			utils.Error(ctx, utils.CodeInternalServerError, "删除帖子失败: "+err.Error())
		}
		return
	}

	utils.SuccessWithMessage(ctx, "删除帖子成功", nil)
}

// GetTopics 获取话题列表
// @Summary 获取话题列表
// @Description 获取所有有效的话题分类
// @Tags 论坛管理
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]string}
// @Router /api/forum/topics [get]
func (c *ForumController) GetTopics(ctx *gin.Context) {
	topics := models.GetValidTopics()
	utils.SuccessWithMessage(ctx, "获取话题列表成功", topics)
}

// IncrementPostView 增加帖子浏览量
// @Summary 增加帖子浏览量
// @Description 增加指定帖子的浏览量
// @Tags 论坛管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "帖子ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/forum/posts/{id}/view [post]
func (c *ForumController) IncrementPostView(ctx *gin.Context) {
	// 获取帖子ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "无效的帖子ID")
		return
	}

	// 调用服务层增加浏览量
	err = c.forumService.IncrementPostView(uint(id))
	if err != nil {
		utils.Error(ctx, utils.CodeInternalServerError, "增加浏览量失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "增加浏览量成功", nil)
}

// CreateReply 创建回复
// @Summary 创建回复
// @Description 创建帖子回复
// @Tags 论坛管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param reply body models.ForumReplyCreateRequest true "回复信息"
// @Success 200 {object} utils.Response{data=models.ForumReplyResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/forum/replies [post]
func (c *ForumController) CreateReply(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 绑定请求参数
	var req models.ForumReplyCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数格式错误: "+err.Error())
		return
	}

	// 调用服务层创建回复
	reply, err := c.forumService.CreateReply(&req, userID)
	if err != nil {
		utils.Error(ctx, utils.CodeInternalServerError, "创建回复失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "创建回复成功", reply.ToResponse())
}

// GetPostReplies 获取帖子回复列表
// @Summary 获取帖子回复列表
// @Description 获取指定帖子的回复列表
// @Tags 论坛管理
// @Accept json
// @Produce json
// @Param id path int true "帖子ID"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Param sort query string false "排序方式" default("created_at asc")
// @Success 200 {object} utils.Response{data=utils.PagedResponse{items=[]models.ForumReplyResponse}}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/forum/posts/{id}/replies [get]
func (c *ForumController) GetPostReplies(ctx *gin.Context) {
	// 获取帖子ID
	idStr := ctx.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "无效的帖子ID")
		return
	}

	// 绑定查询参数
	var req models.ForumReplyListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数格式错误: "+err.Error())
		return
	}

	// 设置帖子ID和默认值
	req.PostID = uint(postID)
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	if req.Sort == "" {
		req.Sort = "created_at asc"
	}

	// 调用服务层获取回复列表
	replies, total, err := c.forumService.GetReplyList(&req)
	if err != nil {
		utils.Error(ctx, utils.CodeInternalServerError, "获取回复列表失败: "+err.Error())
		return
	}

	// 转换为响应格式
	var replyResponses []models.ForumReplyResponse
	for _, reply := range replies {
		replyResponses = append(replyResponses, *reply.ToResponse())
	}

	// 构造分页响应
	response := utils.PagedResponse{
		Items: replyResponses,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
		Pages: (total + int64(req.Size) - 1) / int64(req.Size),
	}

	utils.SuccessWithMessage(ctx, "获取回复列表成功", response)
}

// GetMyPosts 获取我的帖子列表
// @Summary 获取我的帖子列表
// @Description 获取当前用户发布的帖子列表
// @Tags 论坛管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=utils.PagedResponse{items=[]models.ForumPostResponse}}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/forum/posts/my [get]
func (c *ForumController) GetMyPosts(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 绑定查询参数
	var req models.ForumPostListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数格式错误: "+err.Error())
		return
	}

	// 设置用户ID和默认值
	req.AuthorID = userID
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.Sort == "" {
		req.Sort = "created_at desc"
	}

	// 调用服务层获取帖子列表
	posts, total, err := c.forumService.GetPostList(&req)
	if err != nil {
		utils.Error(ctx, utils.CodeInternalServerError, "获取我的帖子列表失败: "+err.Error())
		return
	}

	// 转换为响应格式
	var postResponses []models.ForumPostResponse
	for _, post := range posts {
		postResponses = append(postResponses, *post.ToResponse(true))
	}

	// 构造分页响应
	response := utils.PagedResponse{
		Items: postResponses,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
		Pages: (total + int64(req.Size) - 1) / int64(req.Size),
	}

	utils.SuccessWithMessage(ctx, "获取我的帖子列表成功", response)
}

// UpdateReply 更新回复
func (c *ForumController) UpdateReply(ctx *gin.Context) {
	// TODO: 实现更新回复逻辑
	utils.SuccessWithMessage(ctx, "功能开发中", nil)
}

// DeleteReply 删除回复
func (c *ForumController) DeleteReply(ctx *gin.Context) {
	// TODO: 实现删除回复逻辑
	utils.SuccessWithMessage(ctx, "功能开发中", nil)
}

// LikePost 点赞帖子
func (c *ForumController) LikePost(ctx *gin.Context) {
	// TODO: 实现点赞帖子逻辑
	utils.SuccessWithMessage(ctx, "功能开发中", nil)
}

// LikeReply 点赞回复
func (c *ForumController) LikeReply(ctx *gin.Context) {
	// TODO: 实现点赞回复逻辑
	utils.SuccessWithMessage(ctx, "功能开发中", nil)
}

// GetHotPosts 获取热门帖子
func (c *ForumController) GetHotPosts(ctx *gin.Context) {
	// TODO: 实现获取热门帖子逻辑
	utils.SuccessWithMessage(ctx, "功能开发中", nil)
}

// GetMyReplies 获取我的回复列表
func (c *ForumController) GetMyReplies(ctx *gin.Context) {
	// TODO: 实现获取我的回复列表逻辑
	utils.SuccessWithMessage(ctx, "功能开发中", nil)
}

// TogglePostTop 置顶/取消置顶帖子
func (c *ForumController) TogglePostTop(ctx *gin.Context) {
	// TODO: 实现置顶帖子逻辑
	utils.SuccessWithMessage(ctx, "功能开发中", nil)
}

// TogglePostHot 标记/取消标记热门帖子
func (c *ForumController) TogglePostHot(ctx *gin.Context) {
	// TODO: 实现标记热门帖子逻辑
	utils.SuccessWithMessage(ctx, "功能开发中", nil)
}

// BatchDeletePosts 批量删除帖子
func (c *ForumController) BatchDeletePosts(ctx *gin.Context) {
	// TODO: 实现批量删除帖子逻辑
	utils.SuccessWithMessage(ctx, "功能开发中", nil)
}

// BatchDeleteReplies 批量删除回复
func (c *ForumController) BatchDeleteReplies(ctx *gin.Context) {
	// TODO: 实现批量删除回复逻辑
	utils.SuccessWithMessage(ctx, "功能开发中", nil)
}