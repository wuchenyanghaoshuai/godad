package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Content   string         `json:"content" gorm:"type:text;not null;comment:评论内容"`
	UserID    uint           `json:"user_id" gorm:"not null;index;comment:用户ID"`
	ArticleID uint           `json:"article_id" gorm:"not null;index;comment:文章ID"`
	ParentID  *uint          `json:"parent_id" gorm:"index;comment:父评论ID"`
	ReplyToID *uint          `json:"reply_to_id" gorm:"index;comment:回复的评论ID"`
	LikeCount int64          `json:"like_count" gorm:"type:bigint;default:0;comment:点赞次数"`
	Status    int8           `json:"status" gorm:"type:tinyint;default:1;comment:状态 0-已删除 1-正常 2-待审核"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联关系
	User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Article  Article   `json:"article,omitempty" gorm:"foreignKey:ArticleID"`
	Parent   *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	ReplyTo  *Comment  `json:"reply_to,omitempty" gorm:"foreignKey:ReplyToID"`
	Replies  []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}

// CommentCreateRequest 评论创建请求
type CommentCreateRequest struct {
	Content   string `json:"content" binding:"required,min=1,max=1000" example:"这篇文章很有用，谢谢分享！"`
	ArticleID uint   `json:"article_id" binding:"required,min=1" example:"1"`
	ParentID  *uint  `json:"parent_id" example:"1"`
	ReplyToID *uint  `json:"reply_to_id" example:"2"`
}

// CommentUpdateRequest 评论更新请求
type CommentUpdateRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// CommentListRequest 评论列表请求
type CommentListRequest struct {
	Page      int  `form:"page" binding:"min=1" example:"1"`
	Size      int  `form:"size" binding:"min=1,max=100" example:"10"`
	ArticleID uint `form:"article_id" binding:"required,min=1" example:"1"`
	ParentID  *uint `form:"parent_id" example:"1"`
	UserID    uint `form:"user_id" example:"1"`
	Status    int8 `form:"status" binding:"min=0,max=2" example:"1"`
	Sort      string `form:"sort" example:"created_at desc"`
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID        uint              `json:"id"`
	Content   string            `json:"content"`
	UserID    uint              `json:"user_id"`
	ArticleID uint              `json:"article_id"`
	ParentID  *uint             `json:"parent_id"`
	ReplyToID *uint             `json:"reply_to_id"`
	LikeCount int64             `json:"like_count"`
	Status    int8              `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	User      *UserResponse     `json:"user,omitempty"`
	ReplyTo   *CommentResponse  `json:"reply_to,omitempty"`
	Replies   []CommentResponse `json:"replies,omitempty"`
	ReplyCount int64            `json:"reply_count,omitempty"` // 回复数量
}

// ToResponse 转换为响应格式
func (c *Comment) ToResponse(includeReplies bool) *CommentResponse {
	resp := &CommentResponse{
		ID:        c.ID,
		Content:   c.Content,
		UserID:    c.UserID,
		ArticleID: c.ArticleID,
		ParentID:  c.ParentID,
		ReplyToID: c.ReplyToID,
		LikeCount: c.LikeCount,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	// 包含用户信息
	if c.User.ID != 0 {
		resp.User = c.User.ToResponse()
	}

	// 包含回复的评论信息
	if c.ReplyTo != nil && c.ReplyTo.ID != 0 {
		resp.ReplyTo = c.ReplyTo.ToResponse(false)
	}

	// 根据需要包含回复列表
	if includeReplies && len(c.Replies) > 0 {
		resp.Replies = make([]CommentResponse, len(c.Replies))
		for i, reply := range c.Replies {
			resp.Replies[i] = *reply.ToResponse(false)
		}
		resp.ReplyCount = int64(len(c.Replies))
	}

	return resp
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}