package services

import (
    "errors"
    "fmt"
    "strings"
    "time"

    "godad-backend/config"
    "godad-backend/models"
    "gorm.io/gorm"
)

type AppealService struct { db *gorm.DB }

func NewAppealService() *AppealService { return &AppealService{ db: config.GetDB() } }

type CreateAppealRequest struct {
    ReportID   uint   `json:"report_id"`
    TargetID   uint   `json:"target_id"` // 后备：仅有内容ID时由服务端匹配最新的reviewed举报
    Reason     string `json:"reason"`
    Evidence   string `json:"evidence"`
}

func (s *AppealService) CreateAppeal(appellantID uint, req *CreateAppealRequest) (*models.Appeal, error) {
    if strings.TrimSpace(req.Reason) == "" {
        return nil, errors.New("申诉原因不能为空")
    }

    // 查找可申诉的举报
    var report models.Report
    if req.ReportID > 0 {
        if err := s.db.Where("id = ? AND status = ?", req.ReportID, "reviewed").First(&report).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) { return nil, errors.New("举报不可申诉或不存在") }
            return nil, err
        }
    } else if req.TargetID > 0 {
        // 同一内容可能有多个举报，选择最新通过的一条，且确保申诉人为该内容作者
        // 尝试文章
        var r1 models.Report
        err1 := s.db.Where("target_id = ? AND target_type = ? AND status = ?", req.TargetID, "article", "reviewed").Order("updated_at DESC").First(&r1).Error
        if err1 == nil {
            var a models.Article
            if err := s.db.Select("id, author_id").Where("id = ?", req.TargetID).First(&a).Error; err == nil && a.AuthorID == appellantID {
                report = r1
            }
        }
        // 尝试帖子
        if report.ID == 0 {
            var r2 models.Report
            err2 := s.db.Where("target_id = ? AND target_type = ? AND status = ?", req.TargetID, "forum_post", "reviewed").Order("updated_at DESC").First(&r2).Error
            if err2 == nil {
                var p models.ForumPost
                if err := s.db.Select("id, author_id").Where("id = ?", req.TargetID).First(&p).Error; err == nil && p.AuthorID == appellantID {
                    report = r2
                }
            }
        }
        if report.ID == 0 {
            return nil, errors.New("未找到可申诉的举报或无权申诉")
        }
    } else {
        return nil, errors.New("缺少report_id或target_id")
    }

    // 校验申诉人必须为被举报内容作者
    targetType := strings.ToLower(report.TargetType)
    if targetType == "article" {
        var a models.Article
        if err := s.db.Select("id, author_id").Where("id = ?", report.TargetID).First(&a).Error; err != nil { return nil, errors.New("文章不存在") }
        if a.AuthorID != appellantID { return nil, errors.New("无权对该内容发起申诉") }
    } else if targetType == "forum_post" {
        var p models.ForumPost
        if err := s.db.Select("id, author_id").Where("id = ?", report.TargetID).First(&p).Error; err != nil { return nil, errors.New("帖子不存在") }
        if p.AuthorID != appellantID { return nil, errors.New("无权对该内容发起申诉") }
    } else {
        return nil, errors.New("不支持的举报类型")
    }

    // 检查是否已申诉
    var exists models.Appeal
    if err := s.db.Where("report_id = ? AND appellant_id = ?", report.ID, appellantID).First(&exists).Error; err == nil {
        return nil, errors.New("已提交过申诉")
    }

    appeal := &models.Appeal{
        ReportID:    report.ID,
        AppellantID: appellantID,
        Reason:      strings.TrimSpace(req.Reason),
        Evidence:    strings.TrimSpace(req.Evidence),
        Status:      "pending",
    }
    if err := s.db.Create(appeal).Error; err != nil { return nil, err }

    // 通知申诉人
    ns := NewNotificationService(s.db)
    msg := "您的申诉已提交，我们将在48小时内处理"
    _ = ns.CreateNotification(&models.Notification{
        ReceiverID: appellantID,
        ActorID:    appellantID,
        Type:       models.NotificationTypeSystem,
        Title:      "申诉受理通知",
        ResourceID: report.TargetID,
        Message:    msg,
        IsRead:     false,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    })

    return appeal, nil
}

type AppealListParams struct {
    Page   int
    Size   int
    Status string
}

func (s *AppealService) ListAppeals(p *AppealListParams) ([]models.Appeal, int64, error) {
    if p.Page <= 0 { p.Page = 1 }
    if p.Size <= 0 || p.Size > 100 { p.Size = 20 }
    q := s.db.Model(&models.Appeal{})
    if p.Status != "" { q = q.Where("status = ?", p.Status) }
    var total int64
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }
    var items []models.Appeal
    if err := q.Order("created_at DESC").Offset((p.Page-1)*p.Size).Limit(p.Size).Find(&items).Error; err != nil { return nil, 0, err }
    return items, total, nil
}

func (s *AppealService) UpdateStatus(id uint, adminID uint, status string, note string) (*models.Appeal, error) {
    if status != "reviewed" && status != "rejected" { return nil, errors.New("无效状态") }
    var a models.Appeal
    if err := s.db.Where("id = ?", id).First(&a).Error; err != nil { return nil, err }
    if a.Status != "pending" { return nil, errors.New("申诉已处理") }
    if err := s.db.Model(&a).Updates(map[string]interface{}{
        "status": status,
        "handled_by": adminID,
        "handled_note": strings.TrimSpace(note),
    }).Error; err != nil { return nil, err }
    // 通知申诉人
    ns := NewNotificationService(s.db)
    result := map[string]string{"reviewed":"已通过","rejected":"已驳回"}[status]
    msg := fmt.Sprintf("您的申诉%s", result)
    if note = strings.TrimSpace(note); note != "" { msg = msg + "。备注：" + note }
    _ = ns.CreateNotification(&models.Notification{
        ReceiverID: a.AppellantID,
        ActorID:    adminID,
        Type:       models.NotificationTypeSystem,
        Title:      "申诉处理结果",
        Message:    msg,
        IsRead:     false,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    })
    return &a, nil
}

