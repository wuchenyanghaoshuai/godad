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

    // 收集目标ID并查询标题
    articleIDs := []uint{}
    postIDs := []uint{}
    for _, r := range items {
        if strings.ToLower(r.TargetType) == "article" {
            articleIDs = append(articleIDs, r.TargetID)
        } else if strings.ToLower(r.TargetType) == "forum_post" {
            postIDs = append(postIDs, r.TargetID)
        }
    }

    // 查询文章标题
    articleTitles := map[uint]string{}
    if len(articleIDs) > 0 {
        var articles []models.Article
        // 使用 Unscoped 包含软删除的文章，便于在举报中心显示已删除/下架文章的标题
        if err := config.GetDB().Unscoped().Select("id, title").Where("id IN ?", articleIDs).Find(&articles).Error; err == nil {
            for _, a := range articles {
                articleTitles[a.ID] = a.Title
            }
        }
    }

    // 查询帖子标题
    postTitles := map[uint]string{}
    if len(postIDs) > 0 {
        var posts []models.ForumPost
        // 使用 Unscoped 包含软删除的帖子，便于在举报中心显示已删除/锁定帖子的标题
        if err := config.GetDB().Unscoped().Select("id, title").Where("id IN ?", postIDs).Find(&posts).Error; err == nil {
            for _, p := range posts {
                postTitles[p.ID] = p.Title
            }
        }
    }

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
        // 填充目标标题
        if strings.ToLower(r.TargetType) == "article" {
            if title, ok := articleTitles[r.TargetID]; ok {
                v.TargetTitle = title
            }
        } else if strings.ToLower(r.TargetType) == "forum_post" {
            if title, ok := postTitles[r.TargetID]; ok {
                v.TargetTitle = title
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
    var body struct {
        Status      string `json:"status"`
        HandledNote string `json:"handled_note"`
        Action      string `json:"action"` // 处理动作：article: unpublish/delete/warning; forum_post: lock/delete/warning
    }
    if err := ctx.ShouldBindJSON(&body); err != nil { utils.Error(ctx, utils.CodeBadRequest, "参数错误"); return }
    adminID, _ := middleware.GetCurrentUserID(ctx)
    // 校验动作
    action := strings.ToLower(strings.TrimSpace(body.Action))
    if body.Status == "reviewed" { // 仅通过时强制需要动作
        // 获取举报以确认类型
        var r models.Report
        if err := config.GetDB().Where("id = ?", uint(id64)).First(&r).Error; err != nil {
            utils.Error(ctx, utils.CodeBadRequest, "举报不存在"); return
        }
        t := strings.ToLower(r.TargetType)
        valid := map[string]map[string]bool{
            "article":    {"unpublish": true, "delete": true, "warning": true},
            "forum_post": {"lock": true, "delete": true, "warning": true},
        }
        if _, ok := valid[t][action]; !ok {
            utils.Error(ctx, utils.CodeBadRequest, "无效处理动作"); return
        }
    }

    updated, err := c.svc.UpdateStatus(uint(id64), body.Status, adminID, body.HandledNote, action)
    if err != nil { utils.Error(ctx, utils.CodeInternalServerError, err.Error()); return }

    // 执行处理动作
    if updated.Status == "reviewed" {
        t := strings.ToLower(updated.TargetType)
        switch t {
        case "article":
            switch action {
            case "unpublish":
                // 使用文章服务以便同时清理缓存
                s := services.NewArticleService()
                st := int8(2)
                _, _ = s.UpdateArticle(updated.TargetID, adminID, &models.ArticleUpdateRequest{ Status: &st })
            case "delete":
                _ = services.NewArticleService().DeleteArticle(updated.TargetID, adminID)
            case "warning":
                // 仅通知，不改动资源
            }
        case "forum_post":
            fs := services.NewForumService()
            switch action {
            case "lock":
                _ = fs.AdminSetPostLock(updated.TargetID, true)
            case "delete":
                _ = fs.AdminDeletePost(updated.TargetID)
            case "warning":
                // 仅通知
            }
        }
    }

    // 发送系统通知给举报人
    go func(r *models.Report) {
        if r == nil { return }
        db := config.GetDB()
        notif := services.NewNotificationService(db)
        // 获取目标标题
        targetTitle := fmt.Sprintf("#%d", r.TargetID)
        var targetAuthorID uint
        if strings.ToLower(r.TargetType) == "article" {
            var a models.Article
            if err := db.Select("id, title, author_id").Where("id = ?", r.TargetID).First(&a).Error; err == nil {
                targetTitle = a.Title
                targetAuthorID = a.AuthorID
            }
        } else if strings.ToLower(r.TargetType) == "forum_post" {
            var p models.ForumPost
            if err := db.Select("id, title, author_id").Where("id = ?", r.TargetID).First(&p).Error; err == nil {
                targetTitle = p.Title
                targetAuthorID = p.AuthorID
            }
        }
        // 组装消息
        resultText := "已处理"
        if r.Status == "reviewed" { resultText = "已通过" } else if r.Status == "rejected" { resultText = "已驳回" }
        note := strings.TrimSpace(r.HandledNote)
        if len(note) > 120 { note = note[:120] + "..." }
        typeMap := map[string]string{"article":"文章","forum_post":"社区帖子"}
        // 包含举报理由与处理动作
        var extra []string
        if r.Reason != "" { extra = append(extra, "原因："+r.Reason) }
        if r.HandledAction != "" {
            actMap := map[string]string{"unpublish":"下架","delete":"删除","lock":"锁定","warning":"警告"}
            act := r.HandledAction
            if v, ok := actMap[strings.ToLower(act)]; ok { act = v }
            extra = append(extra, "动作："+act)
        }
        extras := strings.Join(extra, "，")
        base := fmt.Sprintf("您对%s《%s》的举报%s", typeMap[strings.ToLower(r.TargetType)], targetTitle, resultText)
        msg := base
        if extras != "" { msg = msg + "（" + extras + "）" }
        if note != "" { msg = msg + "。备注：" + note }

        n := &models.Notification{
            ReceiverID: r.ReporterID,
            ActorID:    adminID,
            Type:       models.NotificationTypeModeration,
            Title:      "举报处理结果",
            ResourceID: r.TargetID,
            Message:    msg,
            IsRead:     false,
        }
        _ = notif.CreateNotification(n)

        // 若举报成立（reviewed），通知被举报内容的作者
        if r.Status == "reviewed" && targetAuthorID != 0 {
            // 被举报作者通知，包含原因与动作
            var extraB []string
            if r.Reason != "" { extraB = append(extraB, "原因："+r.Reason) }
            if r.HandledAction != "" {
                actMap := map[string]string{"unpublish":"下架","delete":"删除","lock":"锁定","warning":"警告"}
                act := r.HandledAction
                if v, ok := actMap[strings.ToLower(act)]; ok { act = v }
                extraB = append(extraB, "已执行："+act)
            }
            extrasB := strings.Join(extraB, "，")
            authorBase := fmt.Sprintf("您的%s《%s》被举报并经审核确认存在违规", typeMap[strings.ToLower(r.TargetType)], targetTitle)
            authorMsg := authorBase
            if extrasB != "" { authorMsg = authorMsg + "（" + extrasB + "）" }
            if note != "" { authorMsg = authorMsg + "。处理备注：" + note }
            // 引导申诉说明（前端会在消息正文内提供“申诉”按钮）
            authorMsg = authorMsg + "。如对此结果有疑问，可在本通知中点击“申诉”发起复核"
            n2 := &models.Notification{
                ReceiverID: targetAuthorID,
                ActorID:    adminID,
                Type:       models.NotificationTypeModeration,
                Title:      "内容违规处理通知",
                ResourceID: r.TargetID,
                Message:    authorMsg,
                IsRead:     false,
            }
            _ = notif.CreateNotification(n2)
        }
    }(updated)

    utils.Success(ctx, nil)
}
