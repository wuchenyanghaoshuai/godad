package services

import (
	"fmt"
	"godad-backend/models"
	"gorm.io/gorm"
)

// TopicService 话题服务
type TopicService struct {
	db           *gorm.DB
	cacheService *CacheService
}

// NewTopicService 创建话题服务实例
func NewTopicService(db *gorm.DB) *TopicService {
	return &TopicService{
		db:           db,
		cacheService: NewCacheService(),
	}
}

// CreateTopic 创建话题
func (s *TopicService) CreateTopic(req *models.TopicCreateRequest) (*models.Topic, error) {
	// 检查话题名称是否已存在
	var existingTopic models.Topic
	if err := s.db.Where("name = ?", req.Name).First(&existingTopic).Error; err == nil {
		return nil, fmt.Errorf("话题名称已存在")
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("检查话题名称失败: %v", err)
	}

	topic := &models.Topic{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		Sort:        req.Sort,
		IsActive:    true,
	}

	if req.IsActive != nil {
		topic.IsActive = *req.IsActive
	}

	if topic.Color == "" {
		topic.Color = "#6366f1" // 默认颜色
	}

	if err := s.db.Create(topic).Error; err != nil {
		return nil, fmt.Errorf("创建话题失败: %v", err)
	}

	// 清除话题列表缓存
	s.clearTopicCache()

	return topic, nil
}

// UpdateTopic 更新话题
func (s *TopicService) UpdateTopic(id uint, req *models.TopicUpdateRequest) (*models.Topic, error) {
	var topic models.Topic
	if err := s.db.First(&topic, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("话题不存在")
		}
		return nil, fmt.Errorf("查找话题失败: %v", err)
	}

	// 如果要更新名称，检查是否与其他话题重复
	if req.Name != "" && req.Name != topic.Name {
		var existingTopic models.Topic
		if err := s.db.Where("name = ? AND id != ?", req.Name, id).First(&existingTopic).Error; err == nil {
			return nil, fmt.Errorf("话题名称已存在")
		} else if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("检查话题名称失败: %v", err)
		}
		topic.Name = req.Name
	}

	// 更新其他字段
	if req.DisplayName != "" {
		topic.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		topic.Description = req.Description
	}
	if req.Color != "" {
		topic.Color = req.Color
	}
	if req.Icon != "" {
		topic.Icon = req.Icon
	}
	if req.Sort != nil {
		topic.Sort = *req.Sort
	}
	if req.IsActive != nil {
		topic.IsActive = *req.IsActive
	}

	if err := s.db.Save(&topic).Error; err != nil {
		return nil, fmt.Errorf("更新话题失败: %v", err)
	}

	// 清除话题缓存
	s.clearTopicCache()

	return &topic, nil
}

// DeleteTopic 删除话题
func (s *TopicService) DeleteTopic(id uint) error {
	var topic models.Topic
	if err := s.db.First(&topic, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("话题不存在")
		}
		return fmt.Errorf("查找话题失败: %v", err)
	}

	// 检查是否有帖子使用了这个话题
	var postCount int64
	if err := s.db.Model(&models.ForumPost{}).Where("topic = ?", topic.Name).Count(&postCount).Error; err != nil {
		return fmt.Errorf("检查话题使用情况失败: %v", err)
	}

	if postCount > 0 {
		return fmt.Errorf("该话题下还有 %d 个帖子，无法删除", postCount)
	}

	if err := s.db.Delete(&topic).Error; err != nil {
		return fmt.Errorf("删除话题失败: %v", err)
	}

	// 清除话题缓存
	s.clearTopicCache()

	return nil
}

// GetTopic 获取话题详情
func (s *TopicService) GetTopic(id uint) (*models.Topic, error) {
	var topic models.Topic
	if err := s.db.First(&topic, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("话题不存在")
		}
		return nil, fmt.Errorf("查找话题失败: %v", err)
	}

	return &topic, nil
}

// GetTopicList 获取话题列表
func (s *TopicService) GetTopicList(page, pageSize int, showAll bool) ([]models.Topic, int64, error) {
	var topics []models.Topic
	var total int64

	query := s.db.Model(&models.Topic{})

	// 如果不是显示全部，只显示启用的话题
	if !showAll {
		query = query.Where("is_active = ?", true)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取话题总数失败: %v", err)
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Order("sort ASC, created_at DESC").Offset(offset).Limit(pageSize).Find(&topics).Error; err != nil {
		return nil, 0, fmt.Errorf("获取话题列表失败: %v", err)
	}

	// 更新帖子数量
	for i := range topics {
		var postCount int64
		s.db.Model(&models.ForumPost{}).Where("topic = ?", topics[i].Name).Count(&postCount)
		topics[i].PostCount = int(postCount)
	}

	return topics, total, nil
}

// GetActiveTopics 获取所有启用的话题（用于发帖时选择）
func (s *TopicService) GetActiveTopics() ([]models.Topic, error) {
	var topics []models.Topic

	if err := s.db.Where("is_active = ?", true).Order("sort ASC, display_name ASC").Find(&topics).Error; err != nil {
		return nil, fmt.Errorf("获取启用话题失败: %v", err)
	}

	// 更新帖子数量
	for i := range topics {
		var postCount int64
		s.db.Model(&models.ForumPost{}).Where("topic = ?", topics[i].Name).Count(&postCount)
		topics[i].PostCount = int(postCount)
	}

	return topics, nil
}

// UpdateTopicPostCount 更新话题帖子数量
func (s *TopicService) UpdateTopicPostCount(topicName string) error {
	var postCount int64
	if err := s.db.Model(&models.ForumPost{}).Where("topic = ?", topicName).Count(&postCount); err != nil {
		return fmt.Errorf("统计帖子数量失败: %v", err)
	}

	if err := s.db.Model(&models.Topic{}).Where("name = ?", topicName).Update("post_count", postCount).Error; err != nil {
		return fmt.Errorf("更新话题帖子数量失败: %v", err)
	}

	// 清除话题缓存
	s.clearTopicCache()

	return nil
}

// clearTopicCache 清除话题相关缓存
func (s *TopicService) clearTopicCache() {
	s.cacheService.DeletePattern("topics:*")
	s.cacheService.DeletePattern("forum:*")
}