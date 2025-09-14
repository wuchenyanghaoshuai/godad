package models

import (
	"time"

	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationTypeLike     NotificationType = "like"     // 点赞
	NotificationTypeComment  NotificationType = "comment"  // 评论
	NotificationTypeBookmark NotificationType = "bookmark" // 收藏
	NotificationTypeFollow   NotificationType = "follow"   // 关注
	NotificationTypeMessage  NotificationType = "message"  // 私信
)

type Notification struct {
	ID         uint             `gorm:"primaryKey" json:"id"`
	ReceiverID uint             `gorm:"not null;index" json:"receiver_id"` // 接收者ID
	ActorID    uint             `gorm:"not null;index" json:"actor_id"`    // 行为发起者ID
	Type       NotificationType `gorm:"not null;type:enum('like','comment','bookmark','follow','message')" json:"type"`
	ResourceID uint             `json:"resource_id,omitempty"` // 资源ID（文章ID、评论ID等）
	Message    string           `gorm:"type:text" json:"message"`
	IsRead     bool             `gorm:"default:false;index" json:"is_read"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	DeletedAt  gorm.DeletedAt   `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	Receiver User     `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
	Actor    User     `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
	Article  *Article `gorm:"foreignKey:ResourceID" json:"article,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}

// NotificationStats 通知统计
type NotificationStats struct {
	UnreadCount int64 `json:"unread_count"`
	TotalCount  int64 `json:"total_count"`
}

// NotificationWithDetails 带详细信息的通知
type NotificationWithDetails struct {
	Notification
	ActorUsername  string `json:"actor_username"`
	ActorNickname  string `json:"actor_nickname"`
	ActorAvatar    string `json:"actor_avatar"`
	ArticleTitle   string `json:"article_title,omitempty"`
	ArticleCover   string `json:"article_cover,omitempty"`
	CommentContent string `json:"comment_content,omitempty"`
}