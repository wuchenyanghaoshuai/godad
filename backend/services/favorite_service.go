package services

import (
	"errors"
	"fmt"
	"godad-backend/models"

	"gorm.io/gorm"
)

// FavoriteService 收藏服务
type FavoriteService struct {
	db                  *gorm.DB
	notificationService *NotificationService
}

// NewFavoriteService 创建收藏服务实例
func NewFavoriteService(db *gorm.DB) *FavoriteService {
	return &FavoriteService{
		db:                  db,
		notificationService: NewNotificationService(db),
	}
}

// ToggleFavorite 切换收藏状态（收藏/取消收藏）
func (s *FavoriteService) ToggleFavorite(userID, articleID uint) (*models.FavoriteResponse, error) {
	// 检查是否已经收藏
	var existingFavorite models.Favorite
	result := s.db.Where("user_id = ? AND article_id = ?", userID, articleID).First(&existingFavorite)

	if result.Error == nil {
		// 已经收藏，取消收藏
		if err := s.db.Delete(&existingFavorite).Error; err != nil {
			return nil, fmt.Errorf("取消收藏失败: %v", err)
		}

		// 更新文章收藏计数
		if err := s.updateFavoriteCount(articleID, -1); err != nil {
			return nil, fmt.Errorf("更新收藏计数失败: %v", err)
		}

		return nil, nil // 返回 nil 表示取消收藏
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询收藏记录失败: %v", result.Error)
	}

	// 没有收藏，添加收藏
	favorite := models.Favorite{
		UserID:    userID,
		ArticleID: articleID,
	}

	if err := s.db.Create(&favorite).Error; err != nil {
		return nil, fmt.Errorf("添加收藏失败: %v", err)
	}

	// 更新文章收藏计数
	if err := s.updateFavoriteCount(articleID, 1); err != nil {
		return nil, fmt.Errorf("更新收藏计数失败: %v", err)
	}

	// 创建收藏通知
	var article models.Article
	if err := s.db.First(&article, articleID).Error; err == nil {
		// 发送收藏通知给文章作者
		if article.AuthorID != userID { // 避免给自己发通知
			if err := s.notificationService.CreateBookmarkNotification(userID, article.AuthorID, articleID); err != nil {
				// 通知发送失败不影响收藏操作，只记录错误
				fmt.Printf("发送收藏通知失败: %v\n", err)
			}
		}
	}

	response := favorite.ToResponse()
	return response, nil
}

// GetFavoriteStatus 获取收藏状态
func (s *FavoriteService) GetFavoriteStatus(userID, articleID uint) (bool, error) {
	if userID == 0 {
		return false, nil
	}

	var count int64
	result := s.db.Model(&models.Favorite{}).Where("user_id = ? AND article_id = ?", userID, articleID).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("查询收藏状态失败: %v", result.Error)
	}

	return count > 0, nil
}

// GetUserFavorites 获取用户收藏列表
func (s *FavoriteService) GetUserFavorites(userID uint, page, pageSize int) ([]models.Favorite, int64, error) {
	var favorites []models.Favorite
	var total int64

	query := s.db.Model(&models.Favorite{}).Where("user_id = ?", userID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取收藏总数失败: %v", err)
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Preload("Article").Preload("Article.Author").
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&favorites).Error; err != nil {
		return nil, 0, fmt.Errorf("获取收藏列表失败: %v", err)
	}

	return favorites, total, nil
}

// BatchGetFavoriteStatus 批量获取收藏状态
func (s *FavoriteService) BatchGetFavoriteStatus(userID uint, articleIDs []uint) (map[uint]bool, error) {
	if userID == 0 || len(articleIDs) == 0 {
		return make(map[uint]bool), nil
	}

	var favorites []models.Favorite
	if err := s.db.Where("user_id = ? AND article_id IN ?", userID, articleIDs).Find(&favorites).Error; err != nil {
		return nil, fmt.Errorf("批量查询收藏状态失败: %v", err)
	}

	result := make(map[uint]bool)
	for _, favorite := range favorites {
		result[favorite.ArticleID] = true
	}

	// 确保所有请求的文章ID都有结果
	for _, articleID := range articleIDs {
		if _, exists := result[articleID]; !exists {
			result[articleID] = false
		}
	}

	return result, nil
}

// GetArticleFavorites 获取文章的收藏列表
func (s *FavoriteService) GetArticleFavorites(articleID uint, page, pageSize int) ([]models.Favorite, int64, error) {
	var favorites []models.Favorite
	var total int64

	query := s.db.Model(&models.Favorite{}).Where("article_id = ?", articleID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取文章收藏总数失败: %v", err)
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&favorites).Error; err != nil {
		return nil, 0, fmt.Errorf("获取文章收藏列表失败: %v", err)
	}

	return favorites, total, nil
}

// DeleteFavorite 删除收藏（通过ID）
func (s *FavoriteService) DeleteFavorite(userID, favoriteID uint) error {
	// 先查询收藏记录
	var favorite models.Favorite
	if err := s.db.Where("id = ? AND user_id = ?", favoriteID, userID).First(&favorite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("收藏记录不存在")
		}
		return fmt.Errorf("查询收藏记录失败: %v", err)
	}

	// 删除收藏记录
	if err := s.db.Delete(&favorite).Error; err != nil {
		return fmt.Errorf("删除收藏失败: %v", err)
	}

	// 更新文章收藏计数
	if err := s.updateFavoriteCount(favorite.ArticleID, -1); err != nil {
		return fmt.Errorf("更新收藏计数失败: %v", err)
	}

	return nil
}

// updateFavoriteCount 更新文章收藏计数
func (s *FavoriteService) updateFavoriteCount(articleID uint, delta int) error {
	return s.db.Model(&models.Article{}).Where("id = ?", articleID).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", delta)).Error
}

// GetPopularFavorites 获取热门收藏（按收藏数排序）
func (s *FavoriteService) GetPopularFavorites(limit int, days int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT article_id, COUNT(*) as favorite_count
		FROM favorites 
		WHERE created_at >= DATE_SUB(NOW(), INTERVAL ? DAY) AND deleted_at IS NULL
		GROUP BY article_id 
		ORDER BY favorite_count DESC 
		LIMIT ?
	`

	if err := s.db.Raw(query, days, limit).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("获取热门收藏失败: %v", err)
	}

	return results, nil
}