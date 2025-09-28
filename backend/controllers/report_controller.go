package controllers

import (
    "fmt"
    "strconv"
    "strings"
    "godad-backend/middleware"
    "godad-backend/services"
    "godad-backend/utils"
    "godad-backend/config"
    "godad-backend/models"
    "github.com/gin-gonic/gin"
)

type ReportController struct {
    svc *services.ReportService
}

func NewReportController() *ReportController { return &ReportController{ svc: services.NewReportService() } }

// CreateReport 创建举报（需登录）
func (c *ReportController) CreateReport(ctx *gin.Context) {
    userID, ok := middleware.GetCurrentUserID(ctx)
    if !ok { utils.Error(ctx, utils.CodeUnauthorized, "请先登录"); return }
    var req services.CreateReportRequest
    if err := ctx.ShouldBindJSON(&req); err != nil { utils.Error(ctx, utils.CodeBadRequest, "参数错误"); return }
    r, err := c.svc.CreateReport(&req, userID)
    if err != nil { utils.Error(ctx, utils.CodeInternalServerError, err.Error()); return }
    utils.SuccessWithMessage(ctx, "举报已提交", r)
}

// MyReports 我的举报列表
func (c *ReportController) MyReports(ctx *gin.Context) {
    userID, ok := middleware.GetCurrentUserID(ctx)
    if !ok { utils.Error(ctx, utils.CodeUnauthorized, "请先登录"); return }
    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
    items, total, err := c.svc.ListReports(&services.ReportListParams{ Page: page, Size: size, ReporterID: userID })
    if err != nil { utils.Error(ctx, utils.CodeInternalServerError, err.Error()); return }
    resp := utils.PagedResponse{ Items: items, Total: total, Page: page, Size: size, Pages: (total + int64(size) - 1) / int64(size) }
    utils.SuccessWithMessage(ctx, "获取成功", resp)
}

// AdminList 举报列表（管理员）
func (c *ReportController) AdminList(ctx *gin.Context) {
    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
    status := ctx.Query("status")
    targetType := ctx.Query("target_type")
    keyword := ctx.Query("keyword")
    items, total, err := c.svc.ListReports(&services.ReportListParams{ Page: page, Size: size, Status: status, TargetType: targetType, Keyword: keyword })
    if err != nil { utils.Error(ctx, utils.CodeInternalServerError, err.Error()); return }

    // 收集用户ID
    idSet := map[uint]struct{}{}
    for _, r := range items {
        idSet[r.ReporterID] = struct{}{}
        if r.HandledBy != nil { idSet[*r.HandledBy] = struct{}{} }
    }
    var ids []uint
    for id := range idSet { ids = append(ids, id) }

    // 查询用户信息
    users := []models.User{}
    if len(ids) > 0 {
        if err := config.GetDB().Where("id IN ?", ids).Find(&users).Error; err != nil {
            utils.Error(ctx, utils.CodeInternalServerError, err.Error()); return
        }
    }
    uMap := map[uint]models.User{}
    for _, u := range users { uMap[u.ID] = u }

    // 构造视图
    viewItems := make([]models.ReportWithDetails, 0, len(items))
    for _, r := range items {
        v := models.ReportWithDetails{ Report: r }
        if u, ok := uMap[r.ReporterID]; ok {
            v.ReporterUsername = u.Username
            v.ReporterNickname = u.Nickname
            v.ReporterAvatar = u.Avatar
        }
        if r.HandledBy != nil {
            if u, ok := uMap[*r.HandledBy]; ok {
                v.HandlerUsername = u.Username
                v.HandlerNickname = u.Nickname
                v.HandlerAvatar = u.Avatar
            }
        }
        viewItems = append(viewItems, v)
    }

    resp := utils.PagedResponse{ Items: viewItems, Total: total, Page: page, Size: size, Pages: (total + int64(size) - 1) / int64(size) }
    utils.SuccessWithMessage(ctx, "获取成功", resp)
}

// AdminUpdateStatus 处理举报（管理员）
func (c *ReportController) AdminUpdateStatus(ctx *gin.Context) {
    id64, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
    if err != nil { utils.Error(ctx, utils.CodeBadRequest, "无效ID"); return }
    var body struct { Status string `json:"status"`; HandledNote string `json:"handled_note"` }
    if err := ctx.ShouldBindJSON(&body); err != nil { utils.Error(ctx, utils.CodeBadRequest, "参数错误"); return }
    adminID, _ := middleware.GetCurrentUserID(ctx)
    updated, err := c.svc.UpdateStatus(uint(id64), body.Status, adminID, body.HandledNote)
    if err != nil { utils.Error(ctx, utils.CodeInternalServerError, err.Error()); return }

    // 发送系统通知给举报人
    go func(r *models.Report) {
        if r == nil { return }
        db := config.GetDB()
        notif := services.NewNotificationService(db)
        // 获取目标标题
        targetTitle := fmt.Sprintf("#%d", r.TargetID)
        if strings.ToLower(r.TargetType) == "article" {
            var a models.Article
            if err := db.Select("id, title").Where("id = ?", r.TargetID).First(&a).Error; err == nil && a.Title != "" {
                targetTitle = a.Title
            }
        } else if strings.ToLower(r.TargetType) == "forum_post" {
            var p models.ForumPost
            if err := db.Select("id, title").Where("id = ?", r.TargetID).First(&p).Error; err == nil && p.Title != "" {
                targetTitle = p.Title
            }
        }
        // 组装消息
        resultText := "已处理"
        if r.Status == "reviewed" { resultText = "已通过" } else if r.Status == "rejected" { resultText = "已驳回" }
        note := strings.TrimSpace(r.HandledNote)
        if len(note) > 120 { note = note[:120] + "..." }
        msg := fmt.Sprintf("您对%s《%s》的举报%s", map[string]string{"article":"文章","forum_post":"社区帖子"}[strings.ToLower(r.TargetType)], targetTitle, resultText)
        if note != "" { msg = msg + "。备注：" + note }

        n := &models.Notification{
            ReceiverID: r.ReporterID,
            ActorID:    adminID,
            Type:       models.NotificationTypeSystem,
            Title:      "举报处理结果",
            ResourceID: r.TargetID,
            Message:    msg,
            IsRead:     false,
        }
        _ = notif.CreateNotification(n)
    }(updated)

    utils.Success(ctx, nil)
}
