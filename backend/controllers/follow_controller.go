package controllers

import (
	"godad-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FollowController struct {
	followService *services.FollowService
}

func NewFollowController(followService *services.FollowService) *FollowController {
	return &FollowController{followService: followService}
}

func (c *FollowController) FollowUser(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	followeeIDStr := ctx.Param("id")
	followeeID, err := strconv.ParseUint(followeeIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = c.followService.FollowUser(userID.(uint), uint(followeeID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Successfully followed user",
		"data":    nil,
	})
}

func (c *FollowController) UnfollowUser(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	followeeIDStr := ctx.Param("id")
	followeeID, err := strconv.ParseUint(followeeIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = c.followService.UnfollowUser(userID.(uint), uint(followeeID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Successfully unfollowed user",
		"data":    nil,
	})
}

func (c *FollowController) CheckFollowStatus(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	targetIDStr := ctx.Param("id")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	isFollowing, err := c.followService.IsFollowing(userID.(uint), uint(targetID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check follow status"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":         200,
		"message":      "Follow status retrieved successfully",
		"data":         gin.H{"is_following": isFollowing},
		"is_following": isFollowing, // Keep backward compatibility
	})
}

func (c *FollowController) GetFollowing(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := c.followService.GetFollowing(userID.(uint), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get following list"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users":  users,
		"total":  total,
		"page":   page,
		"limit":  limit,
		"pages":  (total + int64(limit) - 1) / int64(limit),
	})
}

func (c *FollowController) GetFollowers(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := c.followService.GetFollowers(userID.(uint), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get followers list"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users":  users,
		"total":  total,
		"page":   page,
		"limit":  limit,
		"pages":  (total + int64(limit) - 1) / int64(limit),
	})
}

func (c *FollowController) GetFollowStats(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	stats, err := c.followService.GetFollowStats(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get follow stats"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Follow stats retrieved successfully",
		"data":    stats,
	})
}

// 获取指定用户的关注统计（公开接口）
func (c *FollowController) GetUserFollowStats(ctx *gin.Context) {
	targetIDStr := ctx.Param("id")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stats, err := c.followService.GetFollowStats(uint(targetID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get follow stats"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Follow stats retrieved successfully",
		"data":    stats,
	})
}

func (c *FollowController) GetMutualFollows(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := c.followService.GetMutualFollows(userID.(uint), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get mutual follows"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users":  users,
		"total":  total,
		"page":   page,
		"limit":  limit,
		"pages":  (total + int64(limit) - 1) / int64(limit),
	})
}

