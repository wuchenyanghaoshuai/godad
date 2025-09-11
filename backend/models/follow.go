package models

import (
	"time"

	"gorm.io/gorm"
)

type Follow struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	FollowerID uint           `gorm:"not null;index" json:"follower_id"`
	FolloweeID uint           `gorm:"not null;index" json:"followee_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Follower User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	Followee User `gorm:"foreignKey:FolloweeID" json:"followee,omitempty"`
}

func (Follow) TableName() string {
	return "follows"
}

type FollowStats struct {
	FollowingCount int64 `json:"following_count"`
	FollowersCount int64 `json:"followers_count"`
}