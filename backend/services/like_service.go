package services

import (
	"errors"
	"fmt"

	"godad-backend/models"
	"gorm.io/gorm"
)

// LikeService 点赞服务
type LikeService struct {
	db                  *gorm.DB
	cacheService        *CacheService
	notificationService *NotificationService
	pointsService       *PointsService
}

// NewLikeService 创建点赞服务实例
func NewLikeService(db *gorm.DB) *LikeService {
	return &LikeService{
		db:                  db,
		cacheService:        NewCacheService(),
		notificationService: NewNotificationService(db),
		pointsService:       NewPointsService(db),
	}
}

// ToggleLike 切换点赞状态（点赞/取消点赞）
func (s *LikeService) ToggleLike(userID uint, targetType string, targetID uint) (*models.LikeResponse, error) {
	// 检查是否已经点赞
	var existingLike models.Like
	result := s.db.Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).First(&existingLike)
	
	if result.Error == nil {
		// 已经点赞，取消点赞
		if err := s.db.Delete(&existingLike).Error; err != nil {
			return nil, fmt.Errorf("取消点赞失败: %v", err)
		}

		// 更新点赞计数
		if err := s.updateLikeCount(targetType, targetID, -1); err != nil {
			return nil, fmt.Errorf("更新点赞计数失败: %v", err)
		}

		// 取消点赞时扣除作者积分（仅对文章点赞）
		if targetType == "article" {
			var article models.Article
			if err := s.db.First(&article, targetID).Error; err == nil {
				if article.AuthorID != userID { // 避免影响自己的积分
					// 扣除文章作者被点赞的积分
					go func() {
						err := s.pointsService.DeductPoints(article.AuthorID, "article_unliked", "like", existingLike.ID, "取消点赞", 2)
						if err != nil {
							fmt.Printf("取消点赞扣除积分失败: %v\n", err)
						}
					}()
				}
			}
		}

		return nil, nil // 返回 nil 表示取消点赞
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询点赞记录失败: %v", result.Error)
	}
	
	// 没有点赞，添加点赞
	like := models.Like{
		UserID:     userID,
		TargetType: targetType,
		TargetID:   targetID,
	}
	
	if err := s.db.Create(&like).Error; err != nil {
		return nil, fmt.Errorf("添加点赞失败: %v", err)
	}
	
	// 更新点赞计数
	if err := s.updateLikeCount(targetType, targetID, 1); err != nil {
		return nil, fmt.Errorf("更新点赞计数失败: %v", err)
	}

	// 创建点赞动态和通知（仅对文章点赞）
	if targetType == "article" {
		var article models.Article
		if err := s.db.First(&article, targetID).Error; err == nil {
			// 发送点赞通知给文章作者
			if article.AuthorID != userID { // 避免给自己发通知
				if err := s.notificationService.CreateLikeNotification(userID, article.AuthorID, targetID); err != nil {
					// 通知发送失败不影响点赞操作，只记录错误
					fmt.Printf("发送点赞通知失败: %v\n", err)
				}

				// 奖励文章作者被点赞积分
				go func() {
					err := s.pointsService.AwardPoints(article.AuthorID, "article_liked", "like", like.ID, "文章被点赞")
					if err != nil {
						fmt.Printf("文章被点赞积分奖励失败: %v\n", err)
					}
				}()
			}
		}
	}
	
	response := &models.LikeResponse{
		ID:         like.ID,
		UserID:     like.UserID,
		TargetType: like.TargetType,
		TargetID:   like.TargetID,
		CreatedAt:  like.CreatedAt,
	}
	
	return response, nil
}

// GetLikeStatus 获取点赞状态
func (s *LikeService) GetLikeStatus(userID uint, targetType string, targetID uint) (*models.LikeStatusResponse, error) {
	// 检查用户是否点赞
	var count int64
	isLiked := false
	
	if userID > 0 {
		s.db.Model(&models.Like{}).Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).Count(&count)
		isLiked = count > 0
	}
	
	// 获取总点赞数
	var totalCount int64
	s.db.Model(&models.Like{}).Where("target_type = ? AND target_id = ?", targetType, targetID).Count(&totalCount)
	
	return &models.LikeStatusResponse{
		IsLiked:   isLiked,
		LikeCount: int(totalCount),
	}, nil
}

// GetUserLikes 获取用户的点赞列表
func (s *LikeService) GetUserLikes(userID uint, targetType string, page, pageSize int) ([]models.Like, int64, error) {
	var likes []models.Like
	var total int64
	
	query := s.db.Model(&models.Like{}).Where("user_id = ?", userID)
	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取点赞总数失败: %v", err)
	}
	
	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&likes).Error; err != nil {
		return nil, 0, fmt.Errorf("获取点赞列表失败: %v", err)
	}
	
	return likes, total, nil
}

// GetLikesByTarget 获取目标对象的点赞列表
func (s *LikeService) GetLikesByTarget(targetType string, targetID uint, page, pageSize int) ([]models.Like, int64, error) {
	var likes []models.Like
	var total int64
	
	query := s.db.Model(&models.Like{}).Where("target_type = ? AND target_id = ?", targetType, targetID)
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取点赞总数失败: %v", err)
	}
	
	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&likes).Error; err != nil {
		return nil, 0, fmt.Errorf("获取点赞列表失败: %v", err)
	}
	
	return likes, total, nil
}

// GetPopularContent 获取热门内容（按点赞数排序）
func (s *LikeService) GetPopularContent(targetType string, limit int, days int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	query := `
		SELECT target_id, COUNT(*) as like_count
		FROM likes 
		WHERE target_type = ? AND created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
		GROUP BY target_id 
		ORDER BY like_count DESC 
		LIMIT ?
	`
	
	if err := s.db.Raw(query, targetType, days, limit).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("获取热门内容失败: %v", err)
	}
	
	return results, nil
}

// updateLikeCount 更新目标对象的点赞计数
func (s *LikeService) updateLikeCount(targetType string, targetID uint, delta int) error {
	var err error
	switch targetType {
	case "article":
		err = s.db.Model(&models.Article{}).Where("id = ?", targetID).UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
		if err == nil {
			// 清除文章缓存
			articleCacheKey := fmt.Sprintf("article:%d", targetID)
			s.cacheService.Delete(articleCacheKey)
			// 清除所有文章列表缓存（使用通配符）
			s.cacheService.DeletePattern("articles:list:*")
			// 同时清除搜索结果缓存，因为搜索结果也包含点赞数
			s.cacheService.DeletePattern("search:*")
		}
	case "comment":
		err = s.db.Model(&models.Comment{}).Where("id = ?", targetID).UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
	case "forum_post":
		err = s.db.Model(&models.ForumPost{}).Where("id = ?", targetID).UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
		if err == nil {
			// 清除论坛帖子缓存
			postCacheKey := fmt.Sprintf("forum_post:%d", targetID)
			s.cacheService.Delete(postCacheKey)
			// 清除论坛帖子列表缓存
			s.cacheService.DeletePattern("forum:posts:*")
		}
	default:
		return fmt.Errorf("不支持的目标类型: %s", targetType)
	}
	return err
}

// BatchGetLikeStatus 批量获取点赞状态
func (s *LikeService) BatchGetLikeStatus(userID uint, targets []map[string]interface{}) (map[string]bool, error) {
	if userID == 0 {
		return make(map[string]bool), nil
	}
	
	result := make(map[string]bool)
	
	for _, target := range targets {
		targetType, ok1 := target["type"].(string)
		targetID, ok2 := target["id"].(uint)
		
		if !ok1 || !ok2 {
			continue
		}
		
		var count int64
		s.db.Model(&models.Like{}).Where("user_id = ? AND target_type = ? AND target_id = ?", userID, targetType, targetID).Count(&count)
		
		key := fmt.Sprintf("%s_%d", targetType, targetID)
		result[key] = count > 0
	}
	
	return result, nil
}