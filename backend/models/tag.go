package models

import (
	"time"
	"gorm.io/gorm"
)

// Tag 标签表
type Tag struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"not null;unique;size:50"`
	Color       string    `json:"color" gorm:"size:7;default:#3B82F6"` // 标签颜色，默认蓝色
	Description string    `json:"description" gorm:"size:200"`
	UsageCount  int       `json:"usage_count" gorm:"default:0"` // 使用次数
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	Articles []Article `json:"articles,omitempty" gorm:"many2many:article_tags"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}

// ArticleTag 文章标签关联表
type ArticleTag struct {
	ArticleID uint `json:"article_id" gorm:"primaryKey"`
	TagID     uint `json:"tag_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	Article Article `json:"article,omitempty" gorm:"foreignKey:ArticleID"`
	Tag     Tag     `json:"tag,omitempty" gorm:"foreignKey:TagID"`
}

// TableName 指定表名
func (ArticleTag) TableName() string {
	return "article_tags"
}

// TagRequest 标签请求
type TagRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Color       string `json:"color" binding:"omitempty,hexcolor"`
	Description string `json:"description" binding:"omitempty,max=200"`
}

// TagResponse 标签响应
type TagResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	Description string    `json:"description"`
	UsageCount  int       `json:"usage_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PopularTag 热门标签
type PopularTag struct {
	Tag        Tag `json:"tag"`
	UsageCount int `json:"usage_count"`
}

// BeforeCreate GORM钩子：创建前
func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	if t.Color == "" {
		t.Color = "#3B82F6" // 默认蓝色
	}
	return nil
}

// BeforeUpdate GORM钩子：更新前
func (t *Tag) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate 文章标签关联表创建前
func (at *ArticleTag) BeforeCreate(tx *gorm.DB) error {
	at.CreatedAt = time.Now()
	return nil
}