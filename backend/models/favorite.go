package models

import (
	"time"

	"gorm.io/gorm"
)

// Favorite 收藏模型
type Favorite struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint           `json:"user_id" gorm:"not null;index;comment:用户ID"`
	ArticleID uint           `json:"article_id" gorm:"not null;index;comment:文章ID"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联关系
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Article Article `json:"article,omitempty" gorm:"foreignKey:ArticleID"`
}

// FavoriteCreateRequest 收藏创建请求
type FavoriteCreateRequest struct {
	ArticleID uint `json:"article_id" binding:"required,min=1" example:"1"`
}

// FavoriteListRequest 收藏列表请求
type FavoriteListRequest struct {
	Page   int    `form:"page" binding:"min=1" example:"1"`
	Size   int    `form:"size" binding:"min=1,max=100" example:"10"`
	UserID uint   `form:"user_id" example:"1"`
	Sort   string `form:"sort" example:"created_at desc"`
}

// FavoriteResponse 收藏响应
type FavoriteResponse struct {
	ID        uint             `json:"id"`
	UserID    uint             `json:"user_id"`
	ArticleID uint             `json:"article_id"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	User      *UserResponse    `json:"user,omitempty"`
	Article   *ArticleResponse `json:"article,omitempty"`
}

// ToResponse 转换为响应格式
func (f *Favorite) ToResponse() *FavoriteResponse {
	resp := &FavoriteResponse{
		ID:        f.ID,
		UserID:    f.UserID,
		ArticleID: f.ArticleID,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}

	// 包含用户信息
	if f.User.ID != 0 {
		resp.User = f.User.ToResponse()
	}

	// 包含文章信息
	if f.Article.ID != 0 {
		resp.Article = f.Article.ToResponse(false) // 不包含文章内容
	}

	return resp
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorites"
}

// BeforeCreate 创建前的钩子函数
func (f *Favorite) BeforeCreate(tx *gorm.DB) error {
	// 检查是否已经收藏过
	var count int64
	tx.Model(&Favorite{}).Where("user_id = ? AND article_id = ?", f.UserID, f.ArticleID).Count(&count)
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}
	return nil
}