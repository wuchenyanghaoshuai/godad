package models

import (
	"time"
	"gorm.io/gorm"
)

// Topic 话题表
type Topic struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"type:varchar(50);not null;uniqueIndex:idx_topic_name"`
	DisplayName string         `json:"display_name" gorm:"type:varchar(100);not null"`
	Description string         `json:"description" gorm:"type:text"`
	Color       string         `json:"color" gorm:"type:varchar(20);default:'#6366f1'"` // 话题颜色
	Icon        string         `json:"icon" gorm:"type:varchar(50)"`                     // 话题图标
	Sort        int            `json:"sort" gorm:"default:0;index"`                      // 排序
	IsActive    bool           `json:"is_active" gorm:"default:true;index"`              // 是否启用
	PostCount   int            `json:"post_count" gorm:"default:0"`                      // 帖子数量
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (Topic) TableName() string {
	return "topics"
}

// TopicCreateRequest 创建话题请求
type TopicCreateRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
	Color       string `json:"color" binding:"omitempty,hex"`
	Icon        string `json:"icon" binding:"max=50"`
	Sort        int    `json:"sort" binding:"min=0"`
	IsActive    *bool  `json:"is_active"`
}

// TopicUpdateRequest 更新话题请求
type TopicUpdateRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=50"`
	DisplayName string `json:"display_name" binding:"omitempty,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
	Color       string `json:"color" binding:"omitempty,hex"`
	Icon        string `json:"icon" binding:"max=50"`
	Sort        *int   `json:"sort" binding:"omitempty,min=0"`
	IsActive    *bool  `json:"is_active"`
}

// TopicResponse 话题响应
type TopicResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	Icon        string    `json:"icon"`
	Sort        int       `json:"sort"`
	IsActive    bool      `json:"is_active"`
	PostCount   int       `json:"post_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate GORM钩子：创建前
func (t *Topic) BeforeCreate(tx *gorm.DB) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate GORM钩子：更新前
func (t *Topic) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}