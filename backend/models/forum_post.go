package models

import (
	"time"

	"gorm.io/gorm"
)

// ForumPost 论坛帖子模型
type ForumPost struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string         `json:"title" gorm:"type:varchar(200);not null;comment:帖子标题"`
	Content     string         `json:"content" gorm:"type:text;not null;comment:帖子内容"`
	Topic       string         `json:"topic" gorm:"type:varchar(50);not null;index;comment:话题分类"`
	AuthorID    uint           `json:"author_id" gorm:"not null;index;comment:作者ID"`
	ViewCount   int64          `json:"view_count" gorm:"type:bigint;default:0;comment:浏览次数"`
	ReplyCount  int64          `json:"reply_count" gorm:"type:bigint;default:0;comment:回复次数"`
	LikeCount   int64          `json:"like_count" gorm:"type:bigint;default:0;comment:点赞次数"`
	IsTop       bool           `json:"is_top" gorm:"type:boolean;default:false;comment:是否置顶"`
    IsHot       bool           `json:"is_hot" gorm:"type:boolean;default:false;comment:是否热门"`
    IsLocked    bool           `json:"is_locked" gorm:"type:boolean;default:false;comment:是否锁定（仅管理员可编辑）"`
	Status      int8           `json:"status" gorm:"type:tinyint;default:1;comment:状态 0-草稿 1-已发布 2-已删除"`
	LastReplyAt *time.Time     `json:"last_reply_at" gorm:"comment:最后回复时间"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联关系
	Author  User         `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Replies []ForumReply `json:"replies,omitempty" gorm:"foreignKey:PostID"`
}

// ForumPostCreateRequest 创建帖子请求
type ForumPostCreateRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200" example:"新手妈妈求助：宝宝睡眠问题"`
	Content string `json:"content" binding:"required,min=1,max=10000" example:"我家宝宝4个月了，最近睡眠很不稳定..."`
	Topic   string `json:"topic" binding:"required,min=1,max=50" example:"Sleep"`
}

// ForumPostUpdateRequest 更新帖子请求
type ForumPostUpdateRequest struct {
	Title   string `json:"title" binding:"min=1,max=200"`
	Content string `json:"content" binding:"min=1,max=10000"`
	Topic   string `json:"topic" binding:"min=1,max=50"`
	Status  *int8  `json:"status" binding:"omitempty,min=0,max=2"`
}

// ForumPostListRequest 帖子列表请求
type ForumPostListRequest struct {
  Page     int    `form:"page" binding:"min=1" example:"1"`
  Size     int    `form:"size" binding:"min=1,max=100" example:"10"`
  Topic    string `form:"topic" example:"Sleep"`
  AuthorID uint   `form:"author_id" example:"1"`
  Keyword  string `form:"keyword" example:"睡眠"`
  Sort     string `form:"sort" example:"created_at desc"` // created_at desc, reply_count desc, view_count desc, last_reply_at desc
  IsTop    *bool  `form:"is_top" example:"true"`
  IsHot    *bool  `form:"is_hot" example:"true"`
  IsLocked *bool  `form:"is_locked" example:"true"`
}

// AdminForumPostListRequest 管理员帖子列表请求（可查看所有状态）
type AdminForumPostListRequest struct {
    Page           int    `form:"page" binding:"min=1" example:"1"`
    Size           int    `form:"size" binding:"min=1,max=100" example:"10"`
    Topic          string `form:"topic" example:"Sleep"`
    AuthorID       uint   `form:"author_id" example:"1"`
    Keyword        string `form:"keyword" example:"睡眠"`
    Sort           string `form:"sort" example:"created_at desc"`
    IsTop          *bool  `form:"is_top" example:"true"`
    IsHot          *bool  `form:"is_hot" example:"true"`
    IsLocked       *bool  `form:"is_locked" example:"true"`
    Status         *int8  `form:"status" example:"1"` // 0-草稿 1-已发布 2-已删除；为空时返回所有
    IncludeDeleted bool   `form:"include_deleted" example:"false"` // 是否包含软删除记录
}

// ForumPostResponse 帖子响应
type ForumPostResponse struct {
	ID          uint              `json:"id"`
	Title       string            `json:"title"`
	Content     string            `json:"content,omitempty"` // 列表时可能不返回完整内容
	Topic       string            `json:"topic"`
	AuthorID    uint              `json:"author_id"`
	ViewCount   int64             `json:"view_count"`
	ReplyCount  int64             `json:"reply_count"`
	LikeCount   int64             `json:"like_count"`
	IsTop       bool              `json:"is_top"`
    IsHot       bool              `json:"is_hot"`
    IsLocked    bool              `json:"is_locked"`
	Status      int8              `json:"status"`
	LastReplyAt *time.Time        `json:"last_reply_at"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Author      *UserResponse     `json:"author,omitempty"`
	RecentReply *ForumReplyResponse `json:"recent_reply,omitempty"` // 最新回复
	TimeAgo     string            `json:"time_ago,omitempty"`      // 前端显示用的相对时间
}

// ToResponse 转换为响应格式
func (fp *ForumPost) ToResponse(includeContent bool) *ForumPostResponse {
	resp := &ForumPostResponse{
		ID:          fp.ID,
		Title:       fp.Title,
		Topic:       fp.Topic,
		AuthorID:    fp.AuthorID,
		ViewCount:   fp.ViewCount,
		ReplyCount:  fp.ReplyCount,
		LikeCount:   fp.LikeCount,
		IsTop:       fp.IsTop,
        IsHot:       fp.IsHot,
        IsLocked:    fp.IsLocked,
		Status:      fp.Status,
		LastReplyAt: fp.LastReplyAt,
		CreatedAt:   fp.CreatedAt,
		UpdatedAt:   fp.UpdatedAt,
	}

	// 根据需要包含内容
	if includeContent {
		resp.Content = fp.Content
	}

	// 包含关联数据
	if fp.Author.ID != 0 {
		resp.Author = fp.Author.ToResponse()
	}

	return resp
}

// TableName 指定表名
func (ForumPost) TableName() string {
	return "forum_posts"
}

// Topic 常量定义
const (
	TopicAll           = "All"
	TopicBabyCare      = "Baby Care"
	TopicFeeding       = "Feeding"
	TopicSleep         = "Sleep"
	TopicHealth        = "Health"
	TopicDevelopment   = "Development"
	TopicActivities    = "Activities"
	TopicGear          = "Gear"
	TopicParenting     = "Parenting"
	TopicFamilyLife    = "Family Life"
	TopicWorkLife      = "Work & Life Balance"
	TopicRelationships = "Relationships"
	TopicMentalHealth  = "Mental Health"
	TopicFinances      = "Finances"
	TopicLegal         = "Legal"
	TopicOther         = "Other"
)

// GetValidTopics 获取有效的话题列表
func GetValidTopics() []string {
	return []string{
		TopicBabyCare,
		TopicFeeding,
		TopicSleep,
		TopicHealth,
		TopicDevelopment,
		TopicActivities,
		TopicGear,
		TopicParenting,
		TopicFamilyLife,
		TopicWorkLife,
		TopicRelationships,
		TopicMentalHealth,
		TopicFinances,
		TopicLegal,
		TopicOther,
	}
}

// IsValidTopic 检查话题是否有效
func IsValidTopic(topic string) bool {
	validTopics := GetValidTopics()
	for _, validTopic := range validTopics {
		if topic == validTopic {
			return true
		}
	}
	return false
}
