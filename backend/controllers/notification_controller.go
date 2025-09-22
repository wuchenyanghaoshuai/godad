package controllers

import (
    "godad-backend/services"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
)

type NotificationController struct {
	notificationService *services.NotificationService
}

func NewNotificationController(notificationService *services.NotificationService) *NotificationController {
	return &NotificationController{
		notificationService: notificationService,
	}
}

// GetNotifications 获取通知列表
func (c *NotificationController) GetNotifications(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	notifications, total, err := c.notificationService.GetNotifications(userID.(uint), page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notifications"})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Notifications retrieved successfully",
		"data": gin.H{
			"notifications": notifications,
			"pagination": gin.H{
				"current_page": page,
				"per_page":     limit,
				"total":        total,
				"total_pages":  totalPages,
			},
		},
	})
}

// GetNotificationStats 获取通知统计
func (c *NotificationController) GetNotificationStats(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	stats, err := c.notificationService.GetNotificationStats(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notification stats"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Notification stats retrieved successfully",
		"data":    stats,
	})
}

// Stream 通知SSE流（开发阶段：定时推送未读数）
func (c *NotificationController) Stream(ctx *gin.Context) {
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 设置SSE头
    ctx.Header("Content-Type", "text/event-stream")
    ctx.Header("Cache-Control", "no-cache")
    ctx.Header("Connection", "keep-alive")

    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    // 首次推送
    if stats, err := c.notificationService.GetNotificationStats(userID.(uint)); err == nil {
        ctx.SSEvent("message", stats)
        ctx.Writer.Flush()
    }

    // 循环推送，直到客户端断开
    for {
        select {
        case <-ctx.Request.Context().Done():
            return
        case <-ticker.C:
            stats, err := c.notificationService.GetNotificationStats(userID.(uint))
            if err != nil {
                continue
            }
            ctx.SSEvent("message", stats)
            ctx.Writer.Flush()
        }
    }
}

// MarkAsRead 标记通知为已读
func (c *NotificationController) MarkAsRead(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		NotificationIDs []uint `json:"notification_ids"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if len(req.NotificationIDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "notification_ids cannot be empty"})
		return
	}

	err := c.notificationService.MarkAsRead(userID.(uint), req.NotificationIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notifications as read"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Notifications marked as read successfully",
		"data":    nil,
	})
}

// MarkAllAsRead 标记所有通知为已读
func (c *NotificationController) MarkAllAsRead(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := c.notificationService.MarkAllAsRead(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark all notifications as read"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "All notifications marked as read successfully",
		"data":    nil,
	})
}

// DeleteNotification 删除通知
func (c *NotificationController) DeleteNotification(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	notificationIDStr := ctx.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	err = c.notificationService.DeleteNotification(userID.(uint), uint(notificationID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Notification deleted successfully",
		"data":    nil,
	})
}

// MarkAsReadByURL 通过URL参数标记单个通知为已读
func (c *NotificationController) MarkAsReadByURL(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	notificationIDStr := ctx.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	err = c.notificationService.MarkAsRead(userID.(uint), []uint{uint(notificationID)})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notification as read"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Notification marked as read successfully",
		"data":    nil,
	})
}

// BatchMarkAsRead 批量标记通知为已读（通过URL参数）
func (c *NotificationController) BatchMarkAsRead(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idsStr := ctx.Query("ids")
	if idsStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ids parameter is required"})
		return
	}

	idStrs := strings.Split(idsStr, ",")
	var notificationIDs []uint

	for _, idStr := range idStrs {
		id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID format"})
			return
		}
		notificationIDs = append(notificationIDs, uint(id))
	}

	err := c.notificationService.MarkAsRead(userID.(uint), notificationIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notifications as read"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Notifications marked as read successfully",
		"data":    nil,
	})
}

// DeleteAllNotifications 删除所有通知
func (c *NotificationController) DeleteAllNotifications(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := c.notificationService.DeleteAllNotifications(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete all notifications"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "All notifications deleted successfully",
		"data":    nil,
	})
}
