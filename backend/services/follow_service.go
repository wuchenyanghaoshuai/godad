package services

import (
	"errors"
	"godad-backend/models"

	"gorm.io/gorm"
)

type FollowService struct {
	db *gorm.DB
}

func NewFollowService(db *gorm.DB) *FollowService {
	return &FollowService{db: db}
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

	return s.db.Create(&follow).Error
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

func (s *FollowService) GetFollowing(userID uint, page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit

	err := s.db.Model(&models.Follow{}).
		Where("follower_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Table("users").
		Select("users.*").
		Joins("INNER JOIN follows ON follows.followee_id = users.id").
		Where("follows.follower_id = ?", userID).
		Order("follows.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

func (s *FollowService) GetFollowers(userID uint, page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit

	err := s.db.Model(&models.Follow{}).
		Where("followee_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Table("users").
		Select("users.*").
		Joins("INNER JOIN follows ON follows.follower_id = users.id").
		Where("follows.followee_id = ?", userID).
		Order("follows.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

func (s *FollowService) GetFollowStats(userID uint) (*models.FollowStats, error) {
	var stats models.FollowStats

	err := s.db.Model(&models.Follow{}).
		Where("follower_id = ?", userID).
		Count(&stats.FollowingCount).Error
	if err != nil {
		return nil, err
	}

	err = s.db.Model(&models.Follow{}).
		Where("followee_id = ?", userID).
		Count(&stats.FollowersCount).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (s *FollowService) GetMutualFollows(userID uint, page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit

	query := `
		SELECT users.* 
		FROM users 
		INNER JOIN follows f1 ON f1.followee_id = users.id 
		INNER JOIN follows f2 ON f2.follower_id = users.id 
		WHERE f1.follower_id = ? AND f2.followee_id = ?
		ORDER BY f1.created_at DESC
		LIMIT ? OFFSET ?
	`

	countQuery := `
		SELECT COUNT(*) 
		FROM users 
		INNER JOIN follows f1 ON f1.followee_id = users.id 
		INNER JOIN follows f2 ON f2.follower_id = users.id 
		WHERE f1.follower_id = ? AND f2.followee_id = ?
	`

	err := s.db.Raw(countQuery, userID, userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Raw(query, userID, userID, limit, offset).Find(&users).Error

	return users, total, err
}