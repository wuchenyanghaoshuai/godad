package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"godad-backend/models"
	"godad-backend/services"
	"godad-backend/utils"
)

// LikeController 点赞控制器
type LikeController struct {
	likeService *services.LikeService
}

// NewLikeController 创建点赞控制器实例
func NewLikeController(likeService *services.LikeService) *LikeController {
	return &LikeController{
		likeService: likeService,
	}
}

// ToggleLike 切换点赞状态
// @Summary 切换点赞状态
// @Description 点赞或取消点赞文章/评论
// @Tags likes
// @Accept json
// @Produce json
// @Param request body models.LikeRequest true "点赞请求"
// @Success 200 {object} utils.Response{data=models.LikeResponse} "点赞成功"
// @Success 200 {object} utils.Response{data=nil} "取消点赞成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未登录"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/likes/toggle [post]
func (lc *LikeController) ToggleLike(c *gin.Context) {
	// 获取当前用户ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "请先登录")
		return
	}
	userID := userIDInterface.(uint)

	// 解析请求参数
	var req models.LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 切换点赞状态
	result, err := lc.likeService.ToggleLike(userID, req.TargetType, req.TargetID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if result == nil {
		utils.SuccessResponse(c, nil, "取消点赞成功")
	} else {
		utils.SuccessResponse(c, result, "点赞成功")
	}
}

// GetLikeStatus 获取点赞状态
// @Summary 获取点赞状态
// @Description 获取指定对象的点赞状态和点赞数
// @Tags likes
// @Produce json
// @Param target_type query string true "目标类型" Enums(article,comment)
// @Param target_id query int true "目标ID"
// @Success 200 {object} utils.Response{data=models.LikeStatusResponse} "获取成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/likes/status [get]
func (lc *LikeController) GetLikeStatus(c *gin.Context) {
	// 获取参数
	targetType := c.Query("target_type")
	targetIDStr := c.Query("target_id")

	if targetType == "" || targetIDStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "target_type 和 target_id 不能为空")
		return
	}

	targetID, err := strconv.ParseUint(targetIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "target_id 必须是数字")
		return
	}

	// 获取当前用户ID（可选）
	var userID uint
	if userIDInterface, exists := c.Get("user_id"); exists {
		userID = userIDInterface.(uint)
	}

	// 获取点赞状态
	status, err := lc.likeService.GetLikeStatus(userID, targetType, uint(targetID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, status, "获取点赞状态成功")
}

// GetUserLikes 获取用户点赞列表
// @Summary 获取用户点赞列表
// @Description 获取用户的点赞记录
// @Tags likes
// @Produce json
// @Param target_type query string false "目标类型" Enums(article,comment)
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} utils.Response{data=utils.PaginationResponse} "获取成功"
// @Failure 401 {object} utils.Response "未登录"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/likes/my [get]
func (lc *LikeController) GetUserLikes(c *gin.Context) {
	// 获取当前用户ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "请先登录")
		return
	}
	userID := userIDInterface.(uint)

	// 获取参数
	targetType := c.DefaultQuery("target_type", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取用户点赞列表
	likes, total, err := lc.likeService.GetUserLikes(userID, targetType, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 构造分页响应
	response := utils.PaginationResponse{
		Data:       likes,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (int(total) + pageSize - 1) / pageSize,
	}

	utils.SuccessResponse(c, response, "获取用户点赞列表成功")
}

// GetLikesByTarget 获取目标对象的点赞列表
// @Summary 获取点赞列表
// @Description 获取指定对象的点赞用户列表
// @Tags likes
// @Produce json
// @Param target_type query string true "目标类型" Enums(article,comment)
// @Param target_id query int true "目标ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} utils.Response{data=utils.PaginationResponse} "获取成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/likes/list [get]
func (lc *LikeController) GetLikesByTarget(c *gin.Context) {
	// 获取参数
	targetType := c.Query("target_type")
	targetIDStr := c.Query("target_id")

	if targetType == "" || targetIDStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "target_type 和 target_id 不能为空")
		return
	}

	targetID, err := strconv.ParseUint(targetIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "target_id 必须是数字")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取点赞列表
	likes, total, err := lc.likeService.GetLikesByTarget(targetType, uint(targetID), page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 构造分页响应
	response := utils.PaginationResponse{
		Data:       likes,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (int(total) + pageSize - 1) / pageSize,
	}

	utils.SuccessResponse(c, response, "获取点赞列表成功")
}

// GetPopularContent 获取热门内容
// @Summary 获取热门内容
// @Description 根据点赞数获取热门内容
// @Tags likes
// @Produce json
// @Param target_type query string true "目标类型" Enums(article,comment)
// @Param limit query int false "数量限制" default(10)
// @Param days query int false "时间范围（天）" default(7)
// @Success 200 {object} utils.Response{data=[]map[string]interface{}} "获取成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/likes/popular [get]
func (lc *LikeController) GetPopularContent(c *gin.Context) {
	// 获取参数
	targetType := c.Query("target_type")
	if targetType == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "target_type 不能为空")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	// 参数校验
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if days < 1 || days > 365 {
		days = 7
	}

	// 获取热门内容
	results, err := lc.likeService.GetPopularContent(targetType, limit, days)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, results, "获取热门内容成功")
}

// BatchGetLikeStatus 批量获取点赞状态
// @Summary 批量获取点赞状态
// @Description 批量获取多个对象的点赞状态
// @Tags likes
// @Accept json
// @Produce json
// @Param targets body []map[string]interface{} true "目标列表"
// @Success 200 {object} utils.Response{data=map[string]bool} "获取成功"
// @Failure 400 {object} utils.Response "请求参数错误"
// @Failure 401 {object} utils.Response "未登录"
// @Failure 500 {object} utils.Response "服务器错误"
// @Router /api/likes/batch-status [post]
func (lc *LikeController) BatchGetLikeStatus(c *gin.Context) {
	// 获取当前用户ID
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "请先登录")
		return
	}
	userID := userIDInterface.(uint)

	// 解析请求参数
	var targets []map[string]interface{}
	if err := c.ShouldBindJSON(&targets); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 批量获取点赞状态
	result, err := lc.likeService.BatchGetLikeStatus(userID, targets)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, result, "批量获取点赞状态成功")
}