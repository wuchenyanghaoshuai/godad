package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"godad-backend/models"
	"godad-backend/services"
	"godad-backend/utils"
)

// TagController 标签控制器
type TagController struct {
	tagService *services.TagService
}

// NewTagController 创建标签控制器实例
func NewTagController(tagService *services.TagService) *TagController {
	return &TagController{
		tagService: tagService,
	}
}

// CreateTag 创建标签
// @Summary 创建标签
// @Description 创建新的标签
// @Tags tags
// @Accept json
// @Produce json
// @Param request body models.TagRequest true "标签信息"
// @Success 201 {object} utils.Response{data=models.TagResponse} "创建成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未登录"
// @Failure 403 {object} utils.Response "权限不足"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/tags [post]
func (tc *TagController) CreateTag(c *gin.Context) {
	// 检查用户权限（需要管理员或内容管理员权限）
	userRole, exists := c.Get("user_role")
	if !exists || (userRole.(int) != 1 && userRole.(int) != 2) {
		utils.ErrorResponse(c, http.StatusForbidden, "权限不足")
		return
	}

	// 解析请求参数
	var req models.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 创建标签
	tag, err := tc.tagService.CreateTag(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, tag, "创建标签成功")
}

// UpdateTag 更新标签
// @Summary 更新标签
// @Description 更新标签信息
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "标签ID"
// @Param request body models.TagRequest true "标签信息"
// @Success 200 {object} utils.Response{data=models.TagResponse} "更新成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未登录"
// @Failure 403 {object} utils.Response "权限不足"
// @Failure 404 {object} utils.Response "标签不存在"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/tags/{id} [put]
func (tc *TagController) UpdateTag(c *gin.Context) {
	// 检查用户权限
	userRole, exists := c.Get("user_role")
	if !exists || (userRole.(int) != 1 && userRole.(int) != 2) {
		utils.ErrorResponse(c, http.StatusForbidden, "权限不足")
		return
	}

	// 获取标签ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "标签ID必须是数字")
		return
	}

	// 解析请求参数
	var req models.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 更新标签
	tag, err := tc.tagService.UpdateTag(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, tag, "更新标签成功")
}

// DeleteTag 删除标签
// @Summary 删除标签
// @Description 删除标签
// @Tags tags
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} utils.Response "删除成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未登录"
// @Failure 403 {object} utils.Response "权限不足"
// @Failure 404 {object} utils.Response "标签不存在"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/tags/{id} [delete]
func (tc *TagController) DeleteTag(c *gin.Context) {
	// 检查用户权限
	userRole, exists := c.Get("user_role")
	if !exists || (userRole.(int) != 1 && userRole.(int) != 2) {
		utils.ErrorResponse(c, http.StatusForbidden, "权限不足")
		return
	}

	// 获取标签ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "标签ID必须是数字")
		return
	}

	// 删除标签
	if err := tc.tagService.DeleteTag(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, nil, "删除标签成功")
}

// GetTags 获取标签列表
// @Summary 获取标签列表
// @Description 获取标签列表，支持搜索和分页
// @Tags tags
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param search query string false "搜索关键词"
// @Success 200 {object} utils.Response{data=utils.PaginationResponse} "获取成功"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/tags [get]
func (tc *TagController) GetTags(c *gin.Context) {
	// 获取参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.DefaultQuery("search", "")

	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取标签列表
	tags, total, err := tc.tagService.GetTags(page, pageSize, search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 构造分页响应
	response := utils.PaginationResponse{
		Data:       tags,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (int(total) + pageSize - 1) / pageSize,
	}

	utils.SuccessResponse(c, response, "获取标签列表成功")
}

// GetTagByID 获取标签详情
// @Summary 获取标签详情
// @Description 根据ID获取标签详细信息
// @Tags tags
// @Produce json
// @Param id path int true "标签ID"
// @Success 200 {object} utils.Response{data=models.TagResponse} "获取成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 404 {object} utils.Response "标签不存在"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/tags/{id} [get]
func (tc *TagController) GetTagByID(c *gin.Context) {
	// 获取标签ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "标签ID必须是数字")
		return
	}

	// 获取标签详情
	tag, err := tc.tagService.GetTagByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, tag, "获取标签详情成功")
}

// GetPopularTags 获取热门标签
// @Summary 获取热门标签
// @Description 获取使用次数最多的热门标签
// @Tags tags
// @Produce json
// @Param limit query int false "数量限制" default(20)
// @Success 200 {object} utils.Response{data=[]models.PopularTag} "获取成功"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/tags/popular [get]
func (tc *TagController) GetPopularTags(c *gin.Context) {
	// 获取参数
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// 获取热门标签
	tags, err := tc.tagService.GetPopularTags(limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, tags, "获取热门标签成功")
}

// GetArticlesByTag 获取标签下的文章
// @Summary 获取标签下的文章
// @Description 获取指定标签下的文章列表
// @Tags tags
// @Produce json
// @Param id path int true "标签ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} utils.Response{data=utils.PaginationResponse} "获取成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/tags/{id}/articles [get]
func (tc *TagController) GetArticlesByTag(c *gin.Context) {
	// 获取标签ID
	idStr := c.Param("id")
	tagID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "标签ID必须是数字")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取标签下的文章
	articles, total, err := tc.tagService.GetArticlesByTag(uint(tagID), page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 构造分页响应
	response := utils.PaginationResponse{
		Data:       articles,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (int(total) + pageSize - 1) / pageSize,
	}

	utils.SuccessResponse(c, response, "获取标签文章成功")
}

// SearchTags 搜索标签
// @Summary 搜索标签
// @Description 搜索标签，用于自动完成功能
// @Tags tags
// @Produce json
// @Param q query string true "搜索关键词"
// @Param limit query int false "数量限制" default(10)
// @Success 200 {object} utils.Response{data=[]models.Tag} "搜索成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/tags/search [get]
func (tc *TagController) SearchTags(c *gin.Context) {
	// 获取搜索关键词
	query := c.Query("q")
	if query == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "搜索关键词不能为空")
		return
	}

	// 获取数量限制
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// 搜索标签
	tags, err := tc.tagService.SearchTags(query, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, tags, "搜索标签成功")
}

// GetTagStats 获取标签统计信息
// @Summary 获取标签统计信息
// @Description 获取标签的统计信息
// @Tags tags
// @Produce json
// @Success 200 {object} utils.Response{data=map[string]interface{}} "获取成功"
// @Failure 401 {object} utils.Response "未登录"
// @Failure 403 {object} utils.Response "权限不足"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/tags/stats [get]
func (tc *TagController) GetTagStats(c *gin.Context) {
	// 检查用户权限
	userRole, exists := c.Get("user_role")
	if !exists || (userRole.(int) != 1 && userRole.(int) != 2) {
		utils.ErrorResponse(c, http.StatusForbidden, "权限不足")
		return
	}

	// 获取标签统计
	stats, err := tc.tagService.GetTagStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, stats, "获取标签统计成功")
}