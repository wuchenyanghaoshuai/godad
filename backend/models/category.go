package models

import (
	"time"

	"gorm.io/gorm"
)

// Category 文章分类模型
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"type:varchar(50);uniqueIndex;not null;comment:分类名称"`
	Slug        string         `json:"slug" gorm:"type:varchar(50);uniqueIndex;not null;comment:分类别名"`
	Description string         `json:"description" gorm:"type:text;comment:分类描述"`
	Icon        string         `json:"icon" gorm:"type:varchar(255);comment:分类图标"`
	Color       string         `json:"color" gorm:"type:varchar(20);comment:分类颜色"`
	Sort        int            `json:"sort" gorm:"type:int;default:0;comment:排序"`
	Status      int8           `json:"status" gorm:"type:tinyint;default:1;comment:状态 0-禁用 1-启用"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联关系
	Articles []Article `json:"articles,omitempty" gorm:"foreignKey:CategoryID"`
}

// CategoryCreateRequest 分类创建请求
type CategoryCreateRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50" example:"育儿知识"`
	Slug        string `json:"slug" binding:"required,min=1,max=50" example:"parenting"`
	Description string `json:"description" binding:"max=500" example:"关于育儿的知识分享"`
	Icon        string `json:"icon" binding:"max=255" example:"baby"`
	Color       string `json:"color" binding:"max=20" example:"#FF6B6B"`
	Sort        int    `json:"sort" binding:"min=0" example:"1"`
}

// CategoryUpdateRequest 分类更新请求
type CategoryUpdateRequest struct {
	Name        string `json:"name" binding:"min=1,max=50"`
	Slug        string `json:"slug" binding:"min=1,max=50"`
	Description string `json:"description" binding:"max=500"`
	Icon        string `json:"icon" binding:"max=255"`
	Color       string `json:"color" binding:"max=20"`
	Sort        int    `json:"sort" binding:"min=0"`
	Status      int8   `json:"status" binding:"min=0,max=1"`
}

// CategoryResponse 分类响应
type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Color       string    `json:"color"`
	Sort        int       `json:"sort"`
	Status      int8      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ArticleCount int64    `json:"article_count,omitempty"` // 文章数量
}

// ToResponse 转换为响应格式
func (c *Category) ToResponse() *CategoryResponse {
	return &CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		Icon:        c.Icon,
		Color:       c.Color,
		Sort:        c.Sort,
		Status:      c.Status,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}