package controllers

import (
	"strconv"

	"godad-backend/models"
	"godad-backend/services"
	"godad-backend/utils"

	"github.com/gin-gonic/gin"
)

// CategoryController 分类控制器
type CategoryController struct {
	categoryService *services.CategoryService
}

// NewCategoryController 创建分类控制器实例
func NewCategoryController() *CategoryController {
	return &CategoryController{
		categoryService: services.NewCategoryService(),
	}
}

// CreateCategory 创建分类
// @Summary 创建分类
// @Description 创建新分类（管理员权限）
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param category body models.CategoryCreateRequest true "分类信息"
// @Success 200 {object} utils.Response{data=models.CategoryResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/admin/categories [post]
func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	// 绑定请求参数
	var req models.CategoryCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数格式错误: "+err.Error())
		return
	}

	// 创建分类
	category, err := c.categoryService.CreateCategory(&req)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(ctx, category.ToResponse())
}

// UpdateCategory 更新分类
// @Summary 更新分类
// @Description 更新分类信息（管理员权限）
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "分类ID"
// @Param category body models.CategoryUpdateRequest true "分类信息"
// @Success 200 {object} utils.Response{data=models.CategoryResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/admin/categories/{id} [put]
func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	// 获取分类ID
	categoryID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "分类ID格式错误")
		return
	}

	// 绑定请求参数
	var req models.CategoryUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数格式错误: "+err.Error())
		return
	}

	// 更新分类
	category, err := c.categoryService.UpdateCategory(uint(categoryID), &req)
	if err != nil {
		if err.Error() == "分类不存在" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	utils.Success(ctx, category.ToResponse())
}

// DeleteCategory 删除分类
// @Summary 删除分类
// @Description 删除分类（管理员权限）
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "分类ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/admin/categories/{id} [delete]
func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	// 获取分类ID
	categoryID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "分类ID格式错误")
		return
	}

	// 删除分类
	err = c.categoryService.DeleteCategory(uint(categoryID))
	if err != nil {
		if err.Error() == "分类不存在" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	utils.Success(ctx, "删除分类成功")
}

// GetCategory 获取分类详情
// @Summary 获取分类详情
// @Description 根据ID获取分类详情
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} utils.Response{data=models.CategoryResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/categories/{id} [get]
func (c *CategoryController) GetCategory(ctx *gin.Context) {
	// 获取分类ID
	categoryID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "分类ID格式错误")
		return
	}

	// 获取分类
	category, err := c.categoryService.GetCategoryByID(uint(categoryID))
	if err != nil {
		if err.Error() == "分类不存在" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	utils.Success(ctx, category.ToResponse())
}

// GetCategoryBySlug 根据别名获取分类详情
// @Summary 根据别名获取分类详情
// @Description 根据别名获取分类详情
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param slug path string true "分类别名"
// @Success 200 {object} utils.Response{data=models.CategoryResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/categories/slug/{slug} [get]
func (c *CategoryController) GetCategoryBySlug(ctx *gin.Context) {
	// 获取分类别名
	slug := ctx.Param("slug")
	if slug == "" {
		utils.Error(ctx, utils.CodeBadRequest, "分类别名不能为空")
		return
	}

	// 获取分类
	category, err := c.categoryService.GetCategoryBySlug(slug)
	if err != nil {
		if err.Error() == "分类不存在" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	utils.Success(ctx, category.ToResponse())
}

// GetCategoryList 获取分类列表
// @Summary 获取分类列表
// @Description 获取分类列表，支持分页和状态筛选
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param status query int false "分类状态" Enums(-1, 0, 1)
// @Success 200 {object} utils.Response{data=utils.PageResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/categories [get]
func (c *CategoryController) GetCategoryList(ctx *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	status, _ := strconv.Atoi(ctx.DefaultQuery("status", "-1"))

	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	// 获取分类列表
	categories, total, err := c.categoryService.GetCategoryList(status, page, size)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, err.Error())
		return
	}

	// 转换为响应格式
	var categoryResponses []*models.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, category.ToResponse())
	}

	// 使用分页响应
	utils.SuccessPage(ctx, categoryResponses, total, page, size)
}

// GetAllCategories 获取所有启用的分类
// @Summary 获取所有启用的分类
// @Description 获取所有启用的分类（不分页）
// @Tags 分类管理
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.CategoryResponse}
// @Failure 500 {object} utils.Response
// @Router /api/categories/all [get]
func (c *CategoryController) GetAllCategories(ctx *gin.Context) {
	// 获取所有启用的分类
	categories, err := c.categoryService.GetAllCategories()
	if err != nil {
		utils.Error(ctx, utils.CodeInternalError, err.Error())
		return
	}

	// 转换为响应格式
	var categoryResponses []*models.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, category.ToResponse())
	}

	utils.Success(ctx, categoryResponses)
}

// GetCategoriesWithCount 获取分类及文章数量
// @Summary 获取分类及文章数量
// @Description 获取所有启用的分类及其文章数量
// @Tags 分类管理
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.CategoryResponse}
// @Failure 500 {object} utils.Response
// @Router /api/categories/with-count [get]
func (c *CategoryController) GetCategoriesWithCount(ctx *gin.Context) {
	// 获取分类及文章数量
	categories, err := c.categoryService.GetCategoriesWithArticleCount()
	if err != nil {
		utils.Error(ctx, utils.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, categories)
}

// UpdateCategorySort 更新分类排序
// @Summary 更新分类排序
// @Description 更新分类排序（管理员权限）
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "分类ID"
// @Param sort body object{sort=int} true "排序值"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/admin/categories/{id}/sort [put]
func (c *CategoryController) UpdateCategorySort(ctx *gin.Context) {
	// 获取分类ID
	categoryID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "分类ID格式错误")
		return
	}

	// 绑定请求参数
	var req struct {
		Sort int `json:"sort" binding:"required,min=0"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数格式错误: "+err.Error())
		return
	}

	// 更新排序
	err = c.categoryService.UpdateCategorySort(uint(categoryID), req.Sort)
	if err != nil {
		if err.Error() == "分类不存在" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	utils.Success(ctx, "更新排序成功")
}

// ToggleCategoryStatus 切换分类状态
// @Summary 切换分类状态
// @Description 切换分类启用/禁用状态（管理员权限）
// @Tags 分类管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "分类ID"
// @Success 200 {object} utils.Response{data=models.CategoryResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/admin/categories/{id}/toggle [put]
func (c *CategoryController) ToggleCategoryStatus(ctx *gin.Context) {
	// 获取分类ID
	categoryID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "分类ID格式错误")
		return
	}

	// 切换状态
	category, err := c.categoryService.ToggleCategoryStatus(uint(categoryID))
	if err != nil {
		if err.Error() == "分类不存在" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	utils.Success(ctx, category.ToResponse())
}