package models

import (
    "time"
    "gorm.io/gorm"
)

// Report 举报模型
type Report struct {
    ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
    TargetType  string         `json:"target_type" gorm:"type:varchar(50);not null;index;comment:被举报对象类型(article/forum_post)"`
    TargetID    uint           `json:"target_id" gorm:"not null;index;comment:被举报对象ID"`
    ReporterID  uint           `json:"reporter_id" gorm:"not null;index;comment:举报者ID"`
    Reason      string         `json:"reason" gorm:"type:varchar(100);not null;comment:举报理由类型"`
    Description string         `json:"description" gorm:"type:varchar(500);comment:补充说明"`
    Evidence    string         `json:"evidence" gorm:"type:varchar(255);comment:证据链接"`
    Status      string         `json:"status" gorm:"type:varchar(20);default:'pending';index;comment:处理状态 pending/reviewed/rejected"`
    HandledBy   *uint          `json:"handled_by" gorm:"index;comment:处理人ID(管理员)"`
    HandledNote string         `json:"handled_note" gorm:"type:varchar(255);comment:处理备注"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Report) TableName() string { return "reports" }

