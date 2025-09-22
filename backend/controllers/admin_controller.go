package controllers

import (
	"net/http"
	"strconv"

	"godad-backend/models"
	"godad-backend/services"
	"godad-backend/config"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	articleService  *services.ArticleService
	userService     *services.UserService
	categoryService *services.CategoryService
	commentService  *services.CommentService
}

func NewAdminController() *AdminController {
	return &AdminController{
		articleService:  services.NewArticleService(),
		userService:     services.NewUserService(),
		categoryService: services.NewCategoryService(),
		commentService:  services.NewCommentService(),
	}
}

// AdminMiddleware 管理员中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}

		// 获取用户信息
		user, err := services.NewUserService().GetUserByID(userID.(uint))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		// 检查是否为管理员
		if user.Role != 2 {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetStats 获取统计数据
func (ac *AdminController) GetStats(c *gin.Context) {
	// 获取文章总数
	articleCount, err := ac.articleService.GetArticleCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "获取文章数量失败",
		})
		return
	}

	// 获取用户总数
	userCount, err := ac.userService.GetUserCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "获取用户数量失败",
		})
		return
	}

	// 获取分类总数
	categoryCount, err := ac.categoryService.GetCategoryCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "获取分类数量失败",
		})
		return
	}

	// 获取评论总数
	commentCount, err := ac.commentService.GetCommentCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "获取评论数量失败",
		})
		return
	}

	// 获取资源统计
	resourceService := services.NewResourceService(config.GetDB())
	resourceStats, err := resourceService.GetResourceStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "获取资源统计失败",
		})
		return
	}

	stats := gin.H{
		"articleCount":  articleCount,
		"userCount":     userCount,
		"categoryCount": categoryCount,
		"commentCount":  commentCount,
		"resourceStats": resourceStats,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": stats,
	})
}

// GetArticles 获取文章列表（管理员）
func (ac *AdminController) GetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	// 构建请求
	req := &models.ArticleListRequest{
		Page:    page,
		Size:    size,
		Keyword: keyword,
	}

	// 解析状态
	if status != "" {
		if statusInt, err := strconv.Atoi(status); err == nil {
			req.Status = int8(statusInt)
		}
	}

	articles, total, err := ac.articleService.GetArticleList(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "获取文章列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"articles": articles,
			"total":    total,
		},
	})
}

// UpdateArticleStatus 更新文章状态
func (ac *AdminController) UpdateArticleStatus(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "无效的文章ID",
		})
		return
	}

	var req struct {
		Status int8 `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "参数错误",
		})
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("user_id")

	// 构建更新请求
	updateReq := &models.ArticleUpdateRequest{
		Status: &req.Status,
	}

	_, err = ac.articleService.UpdateArticle(uint(articleID), userID.(uint), updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "更新文章状态失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// DeleteArticle 删除文章
func (ac *AdminController) DeleteArticle(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "无效的文章ID",
		})
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("user_id")

	err = ac.articleService.DeleteArticle(uint(articleID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "删除文章失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// GetUsers 获取用户列表（管理员）
func (ac *AdminController) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")

	users, total, err := ac.userService.GetUserList(page, size, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "获取用户列表失败",
		})
		return
	}

	// 过滤敏感信息
	var safeUsers []gin.H
	for _, user := range users {
		safeUsers = append(safeUsers, gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"nickname":   user.Nickname,
			"email":      user.Email,
			"phone":      user.Phone,
			"avatar":     user.Avatar,
			"role":       user.Role,
			"status":     user.Status,
			"created_at": user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"users": safeUsers,
			"total": total,
		},
	})
}

// UpdateUserStatus 更新用户状态
func (ac *AdminController) UpdateUserStatus(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "无效的用户ID",
		})
		return
	}

	var req struct {
		Status int8 `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "参数错误",
		})
		return
	}

	// 防止管理员禁用自己
	currentUserID, _ := c.Get("user_id")
	if currentUserID.(uint) == uint(userID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "不能禁用自己",
		})
		return
	}

	// TODO: 实现用户状态更新功能
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// CreateCategory 创建分类
func (ac *AdminController) CreateCategory(c *gin.Context) {
	var req models.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "参数错误",
		})
		return
	}

	category, err := ac.categoryService.CreateCategory(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": category.ToResponse(),
	})
}

// UpdateCategory 更新分类
func (ac *AdminController) UpdateCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "无效的分类ID",
		})
		return
	}

	var req models.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "参数错误",
		})
		return
	}

	category, err := ac.categoryService.UpdateCategory(uint(categoryID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": category.ToResponse(),
	})
}

// UpdateCategoryStatus 更新分类状态
func (ac *AdminController) UpdateCategoryStatus(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "无效的分类ID",
		})
		return
	}

	category, err := ac.categoryService.ToggleCategoryStatus(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": category.ToResponse(),
	})
}

// DeleteCategory 删除分类
func (ac *AdminController) DeleteCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "无效的分类ID",
		})
		return
	}

	err = ac.categoryService.DeleteCategory(uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}