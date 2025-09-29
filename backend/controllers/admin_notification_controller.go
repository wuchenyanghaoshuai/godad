package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "godad-backend/services"
)

// AdminNotificationController 管理员通知相关控制器
type AdminNotificationController struct {
    notificationService *services.NotificationService
}

func NewAdminNotificationController(notificationService *services.NotificationService) *AdminNotificationController {
    return &AdminNotificationController{ notificationService: notificationService }
}

// BroadcastSystemNotification 管理员广播系统通知
// POST /api/admin/notifications/system/broadcast
// 请求体: { title?: string, content: string }
func (c *AdminNotificationController) BroadcastSystemNotification(ctx *gin.Context) {
    // 管理员鉴权已在路由中通过中间件处理
    userID, exists := ctx.Get("user_id")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
        return
    }

    var req struct {
        Title   string `json:"title"`
        Content string `json:"content" binding:"required"`
    }
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误或缺少内容"})
        return
    }

    if err := c.notificationService.BroadcastSystemNotificationWithTitle(userID.(uint), req.Title, req.Content); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "广播失败", "detail": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "code":    200,
        "message": "系统通知广播成功",
    })
}

// ListSystemNotifications 历史广播列表
// GET /api/admin/notifications/system/history?page=&size=
func (c *AdminNotificationController) ListSystemNotifications(ctx *gin.Context) {
    page := 1
    size := 10
    if v := ctx.Query("page"); v != "" { if p, err := strconv.Atoi(v); err == nil && p > 0 { page = p } }
    if v := ctx.Query("size"); v != "" { if s, err := strconv.Atoi(v); err == nil && s > 0 && s <= 100 { size = s } }

    scope := ctx.Query("type")     // system|moderation|all
    subtype := ctx.Query("subtype") // reporter|author|""
    list, total, err := c.notificationService.ListSystemBroadcasts(page, size, scope, subtype)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取历史广播失败"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{ "code": 200, "data": gin.H{ "items": list, "total": total, "page": page, "size": size } })
}
