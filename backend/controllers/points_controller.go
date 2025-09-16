package controllers

import (
	"godad-backend/middleware"
	"godad-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PointsController struct {
	pointsService *services.PointsService
}

func NewPointsController(pointsService *services.PointsService) *PointsController {
	return &PointsController{
		pointsService: pointsService,
	}
}

// GetUserPoints 获取用户积分信息
func (pc *PointsController) GetUserPoints(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)

	userPoints, err := pc.pointsService.GetUserPoints(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户积分失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    userPoints,
	})
}

// GetPointsHistory 获取积分历史记录
func (pc *PointsController) GetPointsHistory(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	transactions, total, err := pc.pointsService.GetPointsHistory(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取积分历史失败",
			"error":   err.Error(),
		})
		return
	}

	totalPages := (int(total) + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"transactions": transactions,
			"pagination": gin.H{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": totalPages,
			},
		},
	})
}

// GetLevels 获取等级配置列表
func (pc *PointsController) GetLevels(c *gin.Context) {
	levels, err := pc.pointsService.GetLevels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取等级配置失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    levels,
	})
}

// GetPointsRules 获取积分规则列表
func (pc *PointsController) GetPointsRules(c *gin.Context) {
	rules, err := pc.pointsService.GetPointsRules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取积分规则失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    rules,
	})
}

// GetPointsStats 获取用户积分统计
func (pc *PointsController) GetPointsStats(c *gin.Context) {
	userID, _ := middleware.GetCurrentUserID(c)

	stats, err := pc.pointsService.GetPointsStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取积分统计失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    stats,
	})
}

// AwardPoints 手动奖励积分（管理员功能）
func (pc *PointsController) AwardPoints(c *gin.Context) {
	var req struct {
		UserID      uint   `json:"user_id" binding:"required"`
		Action      string `json:"action" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数格式错误",
			"error":   err.Error(),
		})
		return
	}

	err := pc.pointsService.AwardPoints(req.UserID, req.Action, "manual", 0, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "积分奖励成功",
	})
}