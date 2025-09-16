package models

import (
	"time"

	"gorm.io/gorm"
)

// UserPoints 用户积分表
type UserPoints struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"user_id" gorm:"uniqueIndex;not null"`
	TotalPoints     int64          `json:"total_points" gorm:"default:0"`
	CurrentLevel    int64          `json:"current_level" gorm:"default:1;index"`
	NextLevelPoints int64          `json:"next_level_points" gorm:"default:0"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`

	// 关联
	User  User      `json:"user" gorm:"foreignKey:UserID"`
	Level UserLevel `json:"level" gorm:"foreignKey:CurrentLevel;references:Level"`
}

// UserLevel 用户等级配置表
type UserLevel struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:50;not null"`
	Level       int64          `json:"level" gorm:"uniqueIndex;not null"`
	MinPoints   int64          `json:"min_points" gorm:"not null"`
	MaxPoints   int64          `json:"max_points" gorm:"not null"`
	Color       string         `json:"color" gorm:"size:7;default:#1976D2"`
	Icon        string         `json:"icon" gorm:"size:50"`
	Badge       string         `json:"badge" gorm:"size:100"`
	Description string         `json:"description" gorm:"size:200"`
	Privileges  string         `json:"privileges" gorm:"type:text"`
	Status      int8           `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// PointsTransaction 积分交易记录表
type PointsTransaction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"index;not null"`
	Action      string    `json:"action" gorm:"size:50;index;not null"`
	Points      int64     `json:"points" gorm:"not null"`
	Description string    `json:"description" gorm:"size:200"`
	SourceType  string    `json:"source_type" gorm:"size:50"`
	SourceID    uint      `json:"source_id"`
	CreatedAt   time.Time `json:"created_at"`

	// 关联
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// PointsRule 积分规则配置表
type PointsRule struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Action      string         `json:"action" gorm:"size:50;uniqueIndex;not null"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Points      int64          `json:"points" gorm:"not null"`
	DailyLimit  int64          `json:"daily_limit" gorm:"default:0"`
	Description string         `json:"description" gorm:"size:200"`
	Status      int8           `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// TableName 指定表名
func (UserPoints) TableName() string {
	return "user_points"
}

func (UserLevel) TableName() string {
	return "user_levels"
}

func (PointsTransaction) TableName() string {
	return "points_transactions"
}

func (PointsRule) TableName() string {
	return "points_rules"
}

// GetLevelByPoints 根据积分获取等级
func GetLevelByPoints(db *gorm.DB, points int64) (*UserLevel, error) {
	var level UserLevel
	err := db.Where("min_points <= ? AND max_points >= ? AND status = 1", points, points).
		Order("level ASC").
		First(&level).Error
	return &level, err
}

// GetAllLevels 获取所有等级配置
func GetAllLevels(db *gorm.DB) ([]*UserLevel, error) {
	var levels []*UserLevel
	err := db.Where("status = 1").Order("level ASC").Find(&levels).Error
	return levels, err
}

// GetPointsRule 根据行为获取积分规则
func GetPointsRule(db *gorm.DB, action string) (*PointsRule, error) {
	var rule PointsRule
	err := db.Where("action = ? AND status = 1", action).First(&rule).Error
	return &rule, err
}

// GetUserPoints 获取用户积分信息
func GetUserPoints(db *gorm.DB, userID uint) (*UserPoints, error) {
	var userPoints UserPoints
	err := db.Preload("User").Preload("Level").
		Where("user_id = ?", userID).
		First(&userPoints).Error
	return &userPoints, err
}

// CreateUserPoints 创建用户积分记录
func CreateUserPoints(db *gorm.DB, userID uint) (*UserPoints, error) {
	userPoints := &UserPoints{
		UserID:          userID,
		TotalPoints:     0,
		CurrentLevel:    1,
		NextLevelPoints: 100, // 默认下一级需要100积分
	}

	err := db.Create(userPoints).Error
	if err != nil {
		return nil, err
	}

	// 预加载关联数据
	err = db.Preload("User").Preload("Level").First(userPoints, userPoints.ID).Error
	return userPoints, err
}