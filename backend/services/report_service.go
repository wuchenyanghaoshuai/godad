package services

import (
    "errors"
    "fmt"
    "strings"

    "godad-backend/config"
    "godad-backend/models"
    "gorm.io/gorm"
)

type ReportService struct {
    db *gorm.DB
}

func NewReportService() *ReportService {
    return &ReportService{ db: config.GetDB() }
}

type CreateReportRequest struct {
    TargetType string `json:"target_type" binding:"required"`
    TargetID   uint   `json:"target_id" binding:"required"`
    Reason     string `json:"reason" binding:"required,min=1,max=100"`
    Description string `json:"description" binding:"max=500"`
    Evidence   string `json:"evidence" binding:"max=255"`
}

func (s *ReportService) CreateReport(req *CreateReportRequest, reporterID uint) (*models.Report, error) {
    rt := strings.ToLower(req.TargetType)
    if rt != "article" && rt != "forum_post" { return nil, fmt.Errorf("不支持的举报类型") }

    r := &models.Report{
        TargetType: rt,
        TargetID: req.TargetID,
        ReporterID: reporterID,
        Reason: strings.TrimSpace(req.Reason),
        Description: strings.TrimSpace(req.Description),
        Evidence: strings.TrimSpace(req.Evidence),
        Status: "pending",
    }
    if err := s.db.Create(r).Error; err != nil { return nil, err }
    return r, nil
}

type ReportListParams struct {
    Page int
    Size int
    Status string
    TargetType string
    Keyword string
    ReporterID uint
}

func (s *ReportService) ListReports(p *ReportListParams) ([]models.Report, int64, error) {
    var items []models.Report
    var total int64
    q := s.db.Model(&models.Report{})
    if p.ReporterID > 0 { q = q.Where("reporter_id = ?", p.ReporterID) }
    if p.Status != "" { q = q.Where("status = ?", p.Status) }
    if p.TargetType != "" { q = q.Where("target_type = ?", p.TargetType) }
    if p.Keyword != "" { like := "%"+p.Keyword+"%"; q = q.Where("reason LIKE ? OR description LIKE ?", like, like) }
    if err := q.Count(&total).Error; err != nil { return nil, 0, err }
    offset := (p.Page-1)*p.Size
    if err := q.Order("created_at DESC").Offset(offset).Limit(p.Size).Find(&items).Error; err != nil { return nil, 0, err }
    return items, total, nil
}

func (s *ReportService) UpdateStatus(id uint, status string, handledBy uint, note string, action string) (*models.Report, error) {
    if status != "reviewed" && status != "rejected" { return nil, fmt.Errorf("无效状态") }
    // 先读取当前状态，防止重复处理
    var current models.Report
    if err := s.db.Where("id = ?", id).First(&current).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) { return nil, fmt.Errorf("举报不存在") }
        return nil, err
    }
    if current.Status != "pending" {
        return nil, fmt.Errorf("举报已处理")
    }
    updates := map[string]interface{}{
        "status":       status,
        "handled_by":   handledBy,
        "handled_note": strings.TrimSpace(note),
    }
    if action != "" { updates["handled_action"] = strings.TrimSpace(action) }
    tx := s.db.Model(&models.Report{}).Where("id = ? AND status = ?", id, "pending").Updates(updates)
    if tx.Error != nil {
        // 兼容旧库：若 handled_action 不存在导致失败，去掉该字段重试
        if action != "" {
            delete(updates, "handled_action")
            if err2 := s.db.Model(&models.Report{}).Where("id = ? AND status = ?", id, "pending").Updates(updates).Error; err2 != nil {
                return nil, tx.Error
            }
        } else {
            return nil, tx.Error
        }
    }
    if tx.RowsAffected == 0 {
        return nil, fmt.Errorf("举报已处理")
    }
    var r models.Report
    if err := s.db.Where("id = ?", id).First(&r).Error; err != nil {
        return nil, err
    }
    return &r, nil
}
