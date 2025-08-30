package models

import (
	"time"
	"gorm.io/gorm"
)

// PasswordReset 密码重置表
type PasswordReset struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"type:varchar(255);not null;index;comment:邮箱地址"`
	Token     string         `json:"token" gorm:"type:varchar(255);not null;unique;comment:重置令牌"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null;comment:过期时间"`
	UsedAt    *time.Time     `json:"used_at" gorm:"comment:使用时间"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (PasswordReset) TableName() string {
	return "password_resets"
}

// IsExpired 检查令牌是否过期
func (pr *PasswordReset) IsExpired() bool {
	return time.Now().After(pr.ExpiresAt)
}

// IsUsed 检查令牌是否已使用
func (pr *PasswordReset) IsUsed() bool {
	return pr.UsedAt != nil
}

// IsValid 检查令牌是否有效（未过期且未使用）
func (pr *PasswordReset) IsValid() bool {
	return !pr.IsExpired() && !pr.IsUsed()
}

// MarkAsUsed 标记令牌为已使用
func (pr *PasswordReset) MarkAsUsed() {
	now := time.Now()
	pr.UsedAt = &now
}