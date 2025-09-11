package services

import (
	"fmt"
	"strings"

	"godad-backend/models"
	"gorm.io/gorm"
)

// TagService 标签服务
type TagService struct {
	db *gorm.DB
}

// NewTagService 创建标签服务实例
func NewTagService(db *gorm.DB) *TagService {
	return &TagService{db: db}
}

// CreateTag 创建标签
func (s *TagService) CreateTag(req *models.TagRequest) (*models.TagResponse, error) {
	// 检查标签名是否已存在
	var existingTag models.Tag
	if err := s.db.Where("name = ?", req.Name).First(&existingTag).Error; err == nil {
		return nil, fmt.Errorf("标签 '%s' 已存在", req.Name)
	}
	
	tag := models.Tag{
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
	}
	
	if err := s.db.Create(&tag).Error; err != nil {
		return nil, fmt.Errorf("创建标签失败: %v", err)
	}
	
	response := &models.TagResponse{
		ID:          tag.ID,
		Name:        tag.Name,
		Color:       tag.Color,
		Description: tag.Description,
		UsageCount:  tag.UsageCount,
		CreatedAt:   tag.CreatedAt,
		UpdatedAt:   tag.UpdatedAt,
	}
	
	return response, nil
}

// UpdateTag 更新标签
func (s *TagService) UpdateTag(id uint, req *models.TagRequest) (*models.TagResponse, error) {
	var tag models.Tag
	if err := s.db.First(&tag, id).Error; err != nil {
		return nil, fmt.Errorf("标签不存在")
	}
	
	// 检查标签名是否已被其他标签使用
	if req.Name != tag.Name {
		var existingTag models.Tag
		if err := s.db.Where("name = ? AND id != ?", req.Name, id).First(&existingTag).Error; err == nil {
			return nil, fmt.Errorf("标签名 '%s' 已被使用", req.Name)
		}
	}
	
	// 更新字段
	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Color != "" {
		tag.Color = req.Color
	}
	if req.Description != "" {
		tag.Description = req.Description
	}
	
	if err := s.db.Save(&tag).Error; err != nil {
		return nil, fmt.Errorf("更新标签失败: %v", err)
	}
	
	response := &models.TagResponse{
		ID:          tag.ID,
		Name:        tag.Name,
		Color:       tag.Color,
		Description: tag.Description,
		UsageCount:  tag.UsageCount,
		CreatedAt:   tag.CreatedAt,
		UpdatedAt:   tag.UpdatedAt,
	}
	
	return response, nil
}

// DeleteTag 删除标签
func (s *TagService) DeleteTag(id uint) error {
	var tag models.Tag
	if err := s.db.First(&tag, id).Error; err != nil {
		return fmt.Errorf("标签不存在")
	}
	
	// 开启事务
	tx := s.db.Begin()
	
	// 删除文章标签关联
	if err := tx.Where("tag_id = ?", id).Delete(&models.ArticleTag{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除标签关联失败: %v", err)
	}
	
	// 删除标签
	if err := tx.Delete(&tag).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除标签失败: %v", err)
	}
	
	return tx.Commit().Error
}

// GetTags 获取标签列表
func (s *TagService) GetTags(page, pageSize int, search string) ([]models.TagResponse, int64, error) {
	var tags []models.Tag
	var total int64
	
	query := s.db.Model(&models.Tag{})
	
	// 搜索
	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取标签总数失败: %v", err)
	}
	
	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Order("usage_count DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&tags).Error; err != nil {
		return nil, 0, fmt.Errorf("获取标签列表失败: %v", err)
	}
	
	// 转换响应格式
	var responses []models.TagResponse
	for _, tag := range tags {
		responses = append(responses, models.TagResponse{
			ID:          tag.ID,
			Name:        tag.Name,
			Color:       tag.Color,
			Description: tag.Description,
			UsageCount:  tag.UsageCount,
			CreatedAt:   tag.CreatedAt,
			UpdatedAt:   tag.UpdatedAt,
		})
	}
	
	return responses, total, nil
}

// GetTagByID 获取标签详情
func (s *TagService) GetTagByID(id uint) (*models.TagResponse, error) {
	var tag models.Tag
	if err := s.db.First(&tag, id).Error; err != nil {
		return nil, fmt.Errorf("标签不存在")
	}
	
	response := &models.TagResponse{
		ID:          tag.ID,
		Name:        tag.Name,
		Color:       tag.Color,
		Description: tag.Description,
		UsageCount:  tag.UsageCount,
		CreatedAt:   tag.CreatedAt,
		UpdatedAt:   tag.UpdatedAt,
	}
	
	return response, nil
}

// GetPopularTags 获取热门标签
func (s *TagService) GetPopularTags(limit int) ([]models.PopularTag, error) {
	var tags []models.Tag
	
	if err := s.db.Where("usage_count > 0").Order("usage_count DESC, created_at DESC").Limit(limit).Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("获取热门标签失败: %v", err)
	}
	
	var results []models.PopularTag
	for _, tag := range tags {
		results = append(results, models.PopularTag{
			Tag:        tag,
			UsageCount: tag.UsageCount,
		})
	}
	
	return results, nil
}

// AttachTagsToArticle 为文章添加标签
func (s *TagService) AttachTagsToArticle(articleID uint, tagNames []string) error {
	if len(tagNames) == 0 {
		return nil
	}
	
	tx := s.db.Begin()
	
	// 删除现有的标签关联
	if err := tx.Where("article_id = ?", articleID).Delete(&models.ArticleTag{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除现有标签关联失败: %v", err)
	}
	
	// 处理每个标签
	for _, tagName := range tagNames {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}
		
		var tag models.Tag
		
		// 查找或创建标签
		if err := tx.Where("name = ?", tagName).First(&tag).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// 创建新标签
				tag = models.Tag{
					Name:       tagName,
					Color:      "#3B82F6", // 默认蓝色
					UsageCount: 0,
				}
				if err := tx.Create(&tag).Error; err != nil {
					tx.Rollback()
					return fmt.Errorf("创建标签失败: %v", err)
				}
			} else {
				tx.Rollback()
				return fmt.Errorf("查询标签失败: %v", err)
			}
		}
		
		// 创建文章标签关联
		articleTag := models.ArticleTag{
			ArticleID: articleID,
			TagID:     tag.ID,
		}
		
		if err := tx.Create(&articleTag).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建标签关联失败: %v", err)
		}
		
		// 更新标签使用次数
		if err := tx.Model(&tag).UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("更新标签使用次数失败: %v", err)
		}
	}
	
	return tx.Commit().Error
}

// DetachTagsFromArticle 移除文章的标签
func (s *TagService) DetachTagsFromArticle(articleID uint) error {
	tx := s.db.Begin()
	
	// 获取要删除的标签ID
	var tagIDs []uint
	if err := tx.Model(&models.ArticleTag{}).Where("article_id = ?", articleID).Pluck("tag_id", &tagIDs).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("获取标签ID失败: %v", err)
	}
	
	// 删除标签关联
	if err := tx.Where("article_id = ?", articleID).Delete(&models.ArticleTag{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除标签关联失败: %v", err)
	}
	
	// 更新标签使用次数
	for _, tagID := range tagIDs {
		if err := tx.Model(&models.Tag{}).Where("id = ?", tagID).UpdateColumn("usage_count", gorm.Expr("GREATEST(usage_count - 1, 0)")).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("更新标签使用次数失败: %v", err)
		}
	}
	
	return tx.Commit().Error
}

// GetArticlesByTag 获取标签下的文章
func (s *TagService) GetArticlesByTag(tagID uint, page, pageSize int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64
	
	// 通过标签关联表查询文章
	subQuery := s.db.Model(&models.ArticleTag{}).Select("article_id").Where("tag_id = ?", tagID)
	
	query := s.db.Model(&models.Article{}).Where("id IN (?) AND status = ?", subQuery, 1)
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取文章总数失败: %v", err)
	}
	
	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Preload("Author").Preload("Category").Order("published_at DESC").Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, fmt.Errorf("获取文章列表失败: %v", err)
	}
	
	return articles, total, nil
}

// GetArticleTags 获取文章的标签列表
func (s *TagService) GetArticleTags(articleID uint) ([]models.Tag, error) {
	var tags []models.Tag
	
	if err := s.db.Joins("JOIN article_tags ON tags.id = article_tags.tag_id").Where("article_tags.article_id = ?", articleID).Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("获取文章标签失败: %v", err)
	}
	
	return tags, nil
}

// SearchTags 搜索标签（用于自动完成）
func (s *TagService) SearchTags(query string, limit int) ([]models.Tag, error) {
	var tags []models.Tag
	
	if err := s.db.Where("name LIKE ?", "%"+query+"%").Order("usage_count DESC, name ASC").Limit(limit).Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("搜索标签失败: %v", err)
	}
	
	return tags, nil
}

// GetTagStats 获取标签统计信息
func (s *TagService) GetTagStats() (map[string]interface{}, error) {
	var totalTags int64
	var totalUsage int64
	var avgUsage float64
	
	// 总标签数
	if err := s.db.Model(&models.Tag{}).Count(&totalTags).Error; err != nil {
		return nil, fmt.Errorf("获取标签总数失败: %v", err)
	}
	
	// 总使用次数
	if err := s.db.Model(&models.Tag{}).Select("SUM(usage_count)").Scan(&totalUsage).Error; err != nil {
		return nil, fmt.Errorf("获取总使用次数失败: %v", err)
	}
	
	// 平均使用次数
	if totalTags > 0 {
		avgUsage = float64(totalUsage) / float64(totalTags)
	}
	
	// 未使用的标签数
	var unusedTags int64
	if err := s.db.Model(&models.Tag{}).Where("usage_count = 0").Count(&unusedTags).Error; err != nil {
		return nil, fmt.Errorf("获取未使用标签数失败: %v", err)
	}
	
	return map[string]interface{}{
		"total_tags":   totalTags,
		"total_usage":  totalUsage,
		"avg_usage":    avgUsage,
		"unused_tags":  unusedTags,
	}, nil
}