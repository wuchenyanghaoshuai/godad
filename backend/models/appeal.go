package models

import (
    "time"
    "gorm.io/gorm"
)

// Appeal 申诉模型
type Appeal struct {
    ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
    ReportID     uint           `json:"report_id" gorm:"not null;index;comment:关联的举报ID"`
    AppellantID  uint           `json:"appellant_id" gorm:"not null;index;comment:申诉人(被举报内容作者)"`
    Reason       string         `json:"reason" gorm:"type:varchar(500);not null;comment:申诉原因"`
    Evidence     string         `json:"evidence" gorm:"type:varchar(255);comment:证据链接或说明"`
    Status       string         `json:"status" gorm:"type:varchar(20);default:'pending';index;comment:状态 pending/reviewed/rejected"`
    HandledBy    *uint          `json:"handled_by" gorm:"index;comment:处理人(管理员)"`
    HandledNote  string         `json:"handled_note" gorm:"type:varchar(255);comment:处理备注"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Appeal) TableName() string { return "appeals" }

