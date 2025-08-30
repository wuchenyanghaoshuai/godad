package models

import (
	"time"

	"gorm.io/gorm"
)

// Article 文章模型
type Article struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string         `json:"title" gorm:"type:varchar(200);not null;comment:文章标题"`
	Slug        string         `json:"slug" gorm:"type:varchar(200);uniqueIndex;not null;comment:文章别名"`
	Summary     string         `json:"summary" gorm:"type:varchar(500);comment:文章摘要"`
	Content     string         `json:"content" gorm:"type:longtext;not null;comment:文章内容"`
	CoverImage  string         `json:"cover_image" gorm:"type:varchar(255);comment:封面图片"`
	AuthorID    uint           `json:"author_id" gorm:"not null;index;comment:作者ID"`
	CategoryID  uint           `json:"category_id" gorm:"not null;index;comment:分类ID"`
	Tags        string         `json:"tags" gorm:"type:varchar(500);comment:标签,逗号分隔"`
	ViewCount   int64          `json:"view_count" gorm:"type:bigint;default:0;comment:浏览次数"`
	LikeCount   int64          `json:"like_count" gorm:"type:bigint;default:0;comment:点赞次数"`
	CommentCount int64         `json:"comment_count" gorm:"type:bigint;default:0;comment:评论次数"`
	FavoriteCount int64        `json:"favorite_count" gorm:"type:bigint;default:0;comment:收藏次数"`
	IsTop       bool           `json:"is_top" gorm:"type:boolean;default:false;comment:是否置顶"`
	IsRecommend bool           `json:"is_recommend" gorm:"type:boolean;default:false;comment:是否推荐"`
	Status      int8           `json:"status" gorm:"type:tinyint;default:1;comment:状态 0-草稿 1-已发布 2-已下架"`
	PublishedAt *time.Time     `json:"published_at" gorm:"comment:发布时间"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联关系
	Author    User       `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Category  Category   `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Comments  []Comment  `json:"comments,omitempty" gorm:"foreignKey:ArticleID"`
	Favorites []Favorite `json:"favorites,omitempty" gorm:"foreignKey:ArticleID"`
}

// ArticleCreateRequest 文章创建请求
type ArticleCreateRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=200" example:"如何培养孩子的阅读习惯"`
	Slug        string `json:"slug" binding:"required,min=1,max=200" example:"how-to-develop-reading-habits"`
	Summary     string `json:"summary" binding:"max=500" example:"本文介绍了培养孩子阅读习惯的几个有效方法"`
	Content     string `json:"content" binding:"required,min=1" example:"文章内容..."`
	CoverImage  string `json:"cover_image" binding:"max=255" example:"https://example.com/cover.jpg"`
	CategoryID  uint   `json:"category_id" binding:"required,min=1" example:"1"`
	Tags        string `json:"tags" binding:"max=500" example:"育儿,阅读,习惯"`
	IsTop       bool   `json:"is_top" example:"false"`
	IsRecommend bool   `json:"is_recommend" example:"false"`
	Status      int8   `json:"status" binding:"min=0,max=2" example:"1"`
}

// ArticleUpdateRequest 文章更新请求
type ArticleUpdateRequest struct {
	Title       string `json:"title" binding:"min=1,max=200"`
	Slug        string `json:"slug" binding:"min=1,max=200"`
	Summary     string `json:"summary" binding:"max=500"`
	Content     string `json:"content" binding:"min=1"`
	CoverImage  string `json:"cover_image" binding:"max=255"`
	CategoryID  uint   `json:"category_id" binding:"min=1"`
	Tags        string `json:"tags" binding:"max=500"`
	IsTop       *bool  `json:"is_top"`
	IsRecommend *bool  `json:"is_recommend"`
	Status      *int8  `json:"status" binding:"omitempty,min=0,max=2"`
}

// ArticleListRequest 文章列表请求
type ArticleListRequest struct {
	Page       int    `form:"page" binding:"min=1" example:"1"`
	Size       int    `form:"size" binding:"min=1,max=100" example:"10"`
	CategoryID uint   `form:"category_id" example:"1"`
	AuthorID   uint   `form:"author_id" example:"1"`
	Keyword    string `form:"keyword" example:"育儿"`
	Status     int8   `form:"status" binding:"omitempty,min=0,max=2" example:"1"`
	IsTop      *bool  `form:"is_top" example:"true"`
	IsRecommend *bool `form:"is_recommend" example:"true"`
	Sort       string `form:"sort" example:"created_at desc"`
}

// ArticleResponse 文章响应
type ArticleResponse struct {
	ID            uint              `json:"id"`
	Title         string            `json:"title"`
	Slug          string            `json:"slug"`
	Summary       string            `json:"summary"`
	Content       string            `json:"content,omitempty"` // 列表时不返回内容
	CoverImage    string            `json:"cover_image"`
	AuthorID      uint              `json:"author_id"`
	CategoryID    uint              `json:"category_id"`
	Tags          string            `json:"tags"`
	ViewCount     int64             `json:"view_count"`
	LikeCount     int64             `json:"like_count"`
	CommentCount  int64             `json:"comment_count"`
	FavoriteCount int64             `json:"favorite_count"`
	IsTop         bool              `json:"is_top"`
	IsRecommend   bool              `json:"is_recommend"`
	Status        int8              `json:"status"`
	PublishedAt   *time.Time        `json:"published_at"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Author        *UserResponse     `json:"author,omitempty"`
	Category      *CategoryResponse `json:"category,omitempty"`
}

// ToResponse 转换为响应格式
func (a *Article) ToResponse(includeContent bool) *ArticleResponse {
	resp := &ArticleResponse{
		ID:            a.ID,
		Title:         a.Title,
		Slug:          a.Slug,
		Summary:       a.Summary,
		CoverImage:    a.CoverImage,
		AuthorID:      a.AuthorID,
		CategoryID:    a.CategoryID,
		Tags:          a.Tags,
		ViewCount:     a.ViewCount,
		LikeCount:     a.LikeCount,
		CommentCount:  a.CommentCount,
		FavoriteCount: a.FavoriteCount,
		IsTop:         a.IsTop,
		IsRecommend:   a.IsRecommend,
		Status:        a.Status,
		PublishedAt:   a.PublishedAt,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}

	// 根据需要包含内容
	if includeContent {
		resp.Content = a.Content
	}

	// 包含关联数据
	if a.Author.ID != 0 {
		resp.Author = a.Author.ToResponse()
	}
	if a.Category.ID != 0 {
		resp.Category = a.Category.ToResponse()
	}

	return resp
}

// TableName 指定表名
func (Article) TableName() string {
	return "articles"
}