package services

import (
	"errors"
	"fmt"
	"godad-backend/models"

	"gorm.io/gorm"
)

type FollowService struct {
	db                  *gorm.DB
	notificationService *NotificationService
}

func NewFollowService(db *gorm.DB) *FollowService {
	return &FollowService{
		db:                  db,
		notificationService: NewNotificationService(db),
	}
}

func (s *FollowService) FollowUser(followerID, followeeID uint) error {
	if followerID == followeeID {
		return errors.New("cannot follow yourself")
	}

	// 检查是否已经关注（排除软删除的记录）
	var existingFollow models.Follow
	result := s.db.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).First(&existingFollow)
	
	if result.Error == nil {
		// 找到记录，说明已经在关注中
		return errors.New("already following this user")
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	// 检查是否有软删除的记录，如果有就物理删除它
	s.db.Unscoped().Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Delete(&models.Follow{})

	// 创建新的关注记录
	follow := models.Follow{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}

	if err := s.db.Create(&follow).Error; err != nil {
		return err
	}

	// 发送关注通知
	if err := s.notificationService.CreateFollowNotification(followerID, followeeID); err != nil {
		// 通知发送失败不影响关注操作，只记录错误
		fmt.Printf("发送关注通知失败: %v\n", err)
	}

	return nil
}

func (s *FollowService) UnfollowUser(followerID, followeeID uint) error {
	result := s.db.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Delete(&models.Follow{})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return errors.New("not following this user")
	}
	
	return nil
}

func (s *FollowService) IsFollowing(followerID, followeeID uint) (bool, error) {
	var count int64
	err := s.db.Model(&models.Follow{}).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Count(&count).Error
	
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

func (s *FollowService) GetFollowing(userID uint, page, limit int) ([]models.UserWithFollowTime, int64, error) {
	var following []models.UserWithFollowTime
	var total int64

	offset := (page - 1) * limit

	// 计算总数时要排除软删除的记录
	err := s.db.Model(&models.Follow{}).
		Where("follower_id = ? AND deleted_at IS NULL", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询关注列表时也要排除软删除的记录
	err = s.db.Table("users").
		Select("users.*, follows.created_at as followed_at").
		Joins("INNER JOIN follows ON follows.followee_id = users.id").
		Where("follows.follower_id = ? AND follows.deleted_at IS NULL", userID).
		Order("follows.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&following).Error

	return following, total, err
}

func (s *FollowService) GetFollowers(userID uint, page, limit int) ([]models.UserWithFollowTime, int64, error) {
	var followers []models.UserWithFollowTime
	var total int64

	offset := (page - 1) * limit

	// 计算总数时要排除软删除的记录
	err := s.db.Model(&models.Follow{}).
		Where("followee_id = ? AND deleted_at IS NULL", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询粉丝列表时也要排除软删除的记录，同时检查是否互关
	err = s.db.Table("users").
		Select(`users.*,
			follows.created_at as followed_at,
			CASE WHEN mutual_follows.id IS NOT NULL THEN true ELSE false END as is_mutual_follow,
			CASE WHEN mutual_follows.id IS NOT NULL THEN true ELSE false END as is_following`).
		Joins("INNER JOIN follows ON follows.follower_id = users.id").
		Joins(`LEFT JOIN follows mutual_follows ON mutual_follows.follower_id = ?
			AND mutual_follows.followee_id = users.id
			AND mutual_follows.deleted_at IS NULL`, userID).
		Where("follows.followee_id = ? AND follows.deleted_at IS NULL", userID).
		Order("follows.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&followers).Error

	return followers, total, err
}

func (s *FollowService) GetFollowStats(userID uint) (*models.FollowStats, error) {
	var stats models.FollowStats

	// 统计关注数（排除软删除）
	err := s.db.Model(&models.Follow{}).
		Where("follower_id = ? AND deleted_at IS NULL", userID).
		Count(&stats.FollowingCount).Error
	if err != nil {
		return nil, err
	}

	// 统计粉丝数（排除软删除）
	err = s.db.Model(&models.Follow{}).
		Where("followee_id = ? AND deleted_at IS NULL", userID).
		Count(&stats.FollowersCount).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (s *FollowService) GetMutualFollows(userID uint, page, limit int) ([]models.UserWithFollowTime, int64, error) {
	var users []models.UserWithFollowTime
	var total int64

	offset := (page - 1) * limit

	// 修复查询，排除软删除的记录
	query := `
		SELECT users.*, f1.created_at as followed_at
		FROM users 
		INNER JOIN follows f1 ON f1.followee_id = users.id 
		INNER JOIN follows f2 ON f2.follower_id = users.id 
		WHERE f1.follower_id = ? AND f2.followee_id = ? 
		AND f1.deleted_at IS NULL AND f2.deleted_at IS NULL
		ORDER BY f1.created_at DESC
		LIMIT ? OFFSET ?
	`

	countQuery := `
		SELECT COUNT(*) 
		FROM users 
		INNER JOIN follows f1 ON f1.followee_id = users.id 
		INNER JOIN follows f2 ON f2.follower_id = users.id 
		WHERE f1.follower_id = ? AND f2.followee_id = ? 
		AND f1.deleted_at IS NULL AND f2.deleted_at IS NULL
	`

	err := s.db.Raw(countQuery, userID, userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Raw(query, userID, userID, limit, offset).Find(&users).Error

	return users, total, err
}