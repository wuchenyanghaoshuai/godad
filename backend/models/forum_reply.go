package models

import (
	"time"

	"gorm.io/gorm"
)

// ForumReply 论坛回复模型
type ForumReply struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	PostID    uint           `json:"post_id" gorm:"not null;index;comment:帖子ID"`
	AuthorID  uint           `json:"author_id" gorm:"not null;index;comment:回复者ID"`
	ParentID  *uint          `json:"parent_id" gorm:"index;comment:父回复ID,用于嵌套回复"`
	Content   string         `json:"content" gorm:"type:text;not null;comment:回复内容"`
	LikeCount int64          `json:"like_count" gorm:"type:bigint;default:0;comment:点赞次数"`
	Status    int8           `json:"status" gorm:"type:tinyint;default:1;comment:状态 0-草稿 1-已发布 2-已删除"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联关系
	Post     ForumPost     `json:"post,omitempty" gorm:"foreignKey:PostID"`
	Author   User          `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Parent   *ForumReply   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []ForumReply  `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// ForumReplyCreateRequest 创建回复请求
type ForumReplyCreateRequest struct {
	PostID   uint   `json:"post_id" binding:"required,min=1" example:"1"`
	ParentID *uint  `json:"parent_id" example:"2"` // 可选，用于嵌套回复
	Content  string `json:"content" binding:"required,min=1,max=5000" example:"我觉得你可以尝试..."`
}

// ForumReplyUpdateRequest 更新回复请求
type ForumReplyUpdateRequest struct {
	Content string `json:"content" binding:"min=1,max=5000"`
	Status  *int8  `json:"status" binding:"omitempty,min=0,max=2"`
}

// ForumReplyListRequest 回复列表请求
type ForumReplyListRequest struct {
	Page     int  `form:"page" binding:"min=1" example:"1"`
	Size     int  `form:"size" binding:"min=1,max=100" example:"20"`
	PostID   uint `form:"post_id" binding:"required,min=1" example:"1"`
	ParentID *uint `form:"parent_id" example:"2"` // 获取某个回复的子回复
	Sort     string `form:"sort" example:"created_at asc"` // created_at asc/desc, like_count desc
}

// ForumReplyResponse 回复响应
type ForumReplyResponse struct {
	ID        uint              `json:"id"`
	PostID    uint              `json:"post_id"`
	AuthorID  uint              `json:"author_id"`
	ParentID  *uint             `json:"parent_id"`
	Content   string            `json:"content"`
	LikeCount int64             `json:"like_count"`
	Status    int8              `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Author    *UserResponse     `json:"author,omitempty"`
	Parent    *ForumReplyResponse `json:"parent,omitempty"`
	Children  []ForumReplyResponse `json:"children,omitempty"`
	TimeAgo   string            `json:"time_ago,omitempty"` // 前端显示用的相对时间
}

// ToResponse 转换为响应格式
func (fr *ForumReply) ToResponse() *ForumReplyResponse {
	resp := &ForumReplyResponse{
		ID:        fr.ID,
		PostID:    fr.PostID,
		AuthorID:  fr.AuthorID,
		ParentID:  fr.ParentID,
		Content:   fr.Content,
		LikeCount: fr.LikeCount,
		Status:    fr.Status,
		CreatedAt: fr.CreatedAt,
		UpdatedAt: fr.UpdatedAt,
	}

	// 包含关联数据
	if fr.Author.ID != 0 {
		resp.Author = fr.Author.ToResponse()
	}

	if fr.Parent != nil && fr.Parent.ID != 0 {
		resp.Parent = fr.Parent.ToResponse()
	}

	// 处理子回复
	if len(fr.Children) > 0 {
		resp.Children = make([]ForumReplyResponse, len(fr.Children))
		for i, child := range fr.Children {
			if childResp := child.ToResponse(); childResp != nil {
				resp.Children[i] = *childResp
			}
		}
	}

	return resp
}

// TableName 指定表名
func (ForumReply) TableName() string {
	return "forum_replies"
}