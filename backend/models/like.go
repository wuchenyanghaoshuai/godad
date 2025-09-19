package models

import (
	"time"
	"gorm.io/gorm"
)

// Like 点赞表
type Like struct {
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     uint   `json:"user_id" gorm:"not null;index"`
	TargetType string `json:"target_type" gorm:"not null;size:20"` // article, comment, forum_post
	TargetID   uint   `json:"target_id" gorm:"not null;index"`
	CreatedAt  time.Time `json:"created_at"`

	// 关联
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (Like) TableName() string {
	return "likes"
}

// LikeRequest 点赞请求
type LikeRequest struct {
	TargetType string `json:"target_type" binding:"required,oneof=article comment forum_post"`
	TargetID   uint   `json:"target_id" binding:"required"`
}

// LikeResponse 点赞响应
type LikeResponse struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	TargetType string    `json:"target_type"`
	TargetID   uint      `json:"target_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// LikeStatus 点赞状态响应
type LikeStatusResponse struct {
	IsLiked   bool `json:"is_liked"`
	LikeCount int  `json:"like_count"`
}

// BeforeCreate GORM钩子：创建前
func (l *Like) BeforeCreate(tx *gorm.DB) error {
	l.CreatedAt = time.Now()
	return nil
}