package controllers

import (
    "net/http"
    "strings"

    "godad-backend/middleware"
    "godad-backend/models"
    "godad-backend/config"

    "github.com/gin-gonic/gin"
)

// MentionController @ 提及相关
type MentionController struct{}

func NewMentionController() *MentionController { return &MentionController{} }

// Suggest 获取可@的好友列表（互相关注）
// GET /api/mentions/suggestions?query=
func (mc *MentionController) Suggest(ctx *gin.Context) {
    userID, ok := middleware.GetCurrentUserID(ctx)
    if !ok {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
        return
    }
    q := strings.TrimSpace(ctx.Query("query"))
    db := config.GetDB()

    type UserItem struct {
        ID       uint   `json:"id"`
        Username string `json:"username"`
        Nickname string `json:"nickname"`
        Avatar   string `json:"avatar"`
    }
    var users []UserItem

    base := db.Table("users u").
        Select("u.id, u.username, u.nickname, u.avatar").
        Joins("INNER JOIN follows f1 ON f1.follower_id = ? AND f1.followee_id = u.id AND f1.deleted_at IS NULL", userID).
        Joins("INNER JOIN follows f2 ON f2.follower_id = u.id AND f2.followee_id = ? AND f2.deleted_at IS NULL", userID).
        Where("u.status = 1")

    if q != "" {
        like := "%" + q + "%"
        base = base.Where("u.username LIKE ? OR u.nickname LIKE ?", like, like)
    }

    if err := base.Order("u.updated_at DESC").Limit(20).Scan(&users).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取候选列表失败"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": users})
}

