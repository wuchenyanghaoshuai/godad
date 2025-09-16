package services

import (
	"fmt"
	"godad-backend/models"
	"time"

	"gorm.io/gorm"
)

type PointsService struct {
	db *gorm.DB
}

func NewPointsService(db *gorm.DB) *PointsService {
	return &PointsService{
		db: db,
	}
}

// AwardPoints 奖励积分
func (ps *PointsService) AwardPoints(userID uint, action string, sourceType string, sourceID uint, description string) error {
	// 获取积分规则
	rule, err := models.GetPointsRule(ps.db, action)
	if err != nil {
		return fmt.Errorf("积分规则不存在: %v", err)
	}

	// 检查每日限制
	if rule.DailyLimit > 0 {
		canAward, err := ps.checkDailyLimit(userID, action, rule.DailyLimit)
		if err != nil {
			return err
		}
		if !canAward {
			return nil // 达到每日限制，静默返回
		}
	}

	// 开始事务
	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建积分交易记录
	transaction := &models.PointsTransaction{
		UserID:      userID,
		Action:      action,
		Points:      rule.Points,
		Description: description,
		SourceType:  sourceType,
		SourceID:    sourceID,
		CreatedAt:   time.Now(),
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新用户积分
	err = ps.updateUserPoints(tx, userID, rule.Points)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// DeductPoints 扣除积分
func (ps *PointsService) DeductPoints(userID uint, action string, sourceType string, sourceID uint, description string, points int64) error {
	// 开始事务
	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建积分交易记录（负值）
	transaction := &models.PointsTransaction{
		UserID:      userID,
		Action:      action,
		Points:      -points, // 负值表示扣除
		Description: description,
		SourceType:  sourceType,
		SourceID:    sourceID,
		CreatedAt:   time.Now(),
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新用户积分（扣除）
	err := ps.updateUserPoints(tx, userID, -points)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetUserPoints 获取用户积分信息
func (ps *PointsService) GetUserPoints(userID uint) (*models.UserPoints, error) {
	userPoints, err := models.GetUserPoints(ps.db, userID)
	if err == gorm.ErrRecordNotFound {
		// 如果用户积分记录不存在，创建一个
		return models.CreateUserPoints(ps.db, userID)
	}
	return userPoints, err
}

// GetPointsHistory 获取用户积分历史记录
func (ps *PointsService) GetPointsHistory(userID uint, page, limit int) ([]*models.PointsTransaction, int64, error) {
	var transactions []*models.PointsTransaction
	var total int64

	query := ps.db.Where("user_id = ?", userID)

	// 计算总数
	err := query.Model(&models.PointsTransaction{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询记录
	err = query.Preload("User").
		Order("created_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&transactions).Error

	return transactions, total, err
}

// GetLevels 获取所有等级配置
func (ps *PointsService) GetLevels() ([]*models.UserLevel, error) {
	return models.GetAllLevels(ps.db)
}

// GetPointsRules 获取积分规则列表
func (ps *PointsService) GetPointsRules() ([]*models.PointsRule, error) {
	var rules []*models.PointsRule
	err := ps.db.Where("status = 1").Order("id ASC").Find(&rules).Error
	return rules, err
}

// updateUserPoints 更新用户积分和等级
func (ps *PointsService) updateUserPoints(tx *gorm.DB, userID uint, points int64) error {
	// 获取或创建用户积分记录
	var userPoints models.UserPoints
	err := tx.Where("user_id = ?", userID).First(&userPoints).Error
	if err == gorm.ErrRecordNotFound {
		// 创建新的积分记录
		userPoints = models.UserPoints{
			UserID:      userID,
			TotalPoints: points,
		}
		// 确保新用户的积分不为负
		if userPoints.TotalPoints < 0 {
			userPoints.TotalPoints = 0
		}
	} else if err != nil {
		return err
	} else {
		// 更新积分
		userPoints.TotalPoints += points
	}

	// 确保积分不会变成负数
	if userPoints.TotalPoints < 0 {
		userPoints.TotalPoints = 0
	}

	// 计算新等级
	newLevel, err := models.GetLevelByPoints(tx, userPoints.TotalPoints)
	if err != nil {
		return err
	}

	oldLevel := userPoints.CurrentLevel
	userPoints.CurrentLevel = newLevel.Level

	// 计算下一等级需要的积分
	nextLevel, err := ps.getNextLevel(tx, newLevel.Level)
	if err == nil {
		userPoints.NextLevelPoints = nextLevel.MinPoints - userPoints.TotalPoints
	} else {
		userPoints.NextLevelPoints = 0 // 已经是最高等级
	}

	// 保存或更新记录
	if userPoints.ID == 0 {
		err = tx.Create(&userPoints).Error
	} else {
		err = tx.Save(&userPoints).Error
	}

	if err != nil {
		return err
	}

	// 如果等级发生变化，可以在这里添加升级通知逻辑
	if oldLevel != 0 && oldLevel != userPoints.CurrentLevel {
		// TODO: 发送升级通知
		fmt.Printf("用户 %d 从等级 %d 升级到等级 %d\n", userID, oldLevel, userPoints.CurrentLevel)
	}

	return nil
}

// checkDailyLimit 检查每日积分获取限制
func (ps *PointsService) checkDailyLimit(userID uint, action string, dailyLimit int64) (bool, error) {
	today := time.Now().Format("2006-01-02")
	startTime, _ := time.Parse("2006-01-02", today)
	endTime := startTime.Add(24 * time.Hour)

	var count int64
	err := ps.db.Model(&models.PointsTransaction{}).
		Where("user_id = ? AND action = ? AND created_at >= ? AND created_at < ?",
			userID, action, startTime, endTime).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count < dailyLimit, nil
}

// getNextLevel 获取下一等级
func (ps *PointsService) getNextLevel(tx *gorm.DB, currentLevel int64) (*models.UserLevel, error) {
	var nextLevel models.UserLevel
	err := tx.Where("level > ? AND status = 1", currentLevel).
		Order("level ASC").
		First(&nextLevel).Error
	return &nextLevel, err
}

// GetPointsStats 获取积分统计信息
func (ps *PointsService) GetPointsStats(userID uint) (map[string]interface{}, error) {
	userPoints, err := ps.GetUserPoints(userID)
	if err != nil {
		return nil, err
	}

	// 获取今日获得的积分
	today := time.Now().Format("2006-01-02")
	startTime, _ := time.Parse("2006-01-02", today)
	endTime := startTime.Add(24 * time.Hour)

	var todayPoints int64
	err = ps.db.Model(&models.PointsTransaction{}).
		Select("COALESCE(SUM(points), 0)").
		Where("user_id = ? AND created_at >= ? AND created_at < ?",
			userID, startTime, endTime).
		Scan(&todayPoints).Error

	if err != nil {
		return nil, err
	}

	// 获取用户排名
	var rank int64
	err = ps.db.Model(&models.UserPoints{}).
		Where("total_points > ?", userPoints.TotalPoints).
		Count(&rank).Error

	if err != nil {
		return nil, err
	}
	rank++ // 排名从1开始

	stats := map[string]interface{}{
		"total_points":      userPoints.TotalPoints,
		"current_level":     userPoints.CurrentLevel,
		"next_level_points": userPoints.NextLevelPoints,
		"today_points":      todayPoints,
		"rank":              rank,
		"level_info":        userPoints.Level,
	}

	return stats, nil
}