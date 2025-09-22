package services

import (
	"errors"
	"fmt"
	"strings"

	"godad-backend/models"

	"gorm.io/gorm"
)

// ResourceService 资源服务
type ResourceService struct {
	db *gorm.DB
}

// NewResourceService 创建资源服务实例
func NewResourceService(db *gorm.DB) *ResourceService {
	return &ResourceService{db: db}
}

// CreateResourceRequest 创建资源请求
type CreateResourceRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Image       string `json:"image"`
	FileURL     string `json:"file_url" binding:"required"`
	ButtonText  string `json:"button_text"`
	Status      int    `json:"status"`
	UploaderID  *uint  `json:"uploader_id"`
}

// UpdateResourceRequest 更新资源请求
type UpdateResourceRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Image       string `json:"image"`
	FileURL     string `json:"file_url"`
	ButtonText  string `json:"button_text"`
	Status      int    `json:"status"`
}

// GetResourcesRequest 获取资源列表请求
type GetResourcesRequest struct {
	Page     int    `form:"page,default=1"`
	Size     int    `form:"size,default=10"`
	Status   string `form:"status"`
	Category string `form:"category"`
	Type     string `form:"type"`
	Keyword  string `form:"keyword"`
}

// ResourceResponse 资源响应
type ResourceResponse struct {
	Items      []models.Resource `json:"items"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Size       int               `json:"size"`
	TotalPages int               `json:"total_pages"`
}

// CreateResource 创建资源
func (s *ResourceService) CreateResource(req *CreateResourceRequest) (*models.Resource, error) {
	// 验证资源类型
	if !isValidResourceType(req.Type) {
		return nil, errors.New("无效的资源类型")
	}

	// 验证资源分类
	if !isValidResourceCategory(req.Category) {
		return nil, errors.New("无效的资源分类")
	}

	// 如果没有提供按钮文本，根据类型设置默认值
	if req.ButtonText == "" {
		switch req.Type {
		case models.ResourceTypeEBook:
			req.ButtonText = "立即下载"
		case models.ResourceTypeVideo:
			req.ButtonText = "立即观看"
		case models.ResourceTypeTool:
			req.ButtonText = "开始使用"
		default:
			req.ButtonText = "立即下载"
		}
	}

	resource := &models.Resource{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Category:    req.Category,
		Image:       req.Image,
		FileURL:     req.FileURL,
		ButtonText:  req.ButtonText,
		Status:      req.Status,
		UploaderID:  req.UploaderID,
	}

	if err := s.db.Create(resource).Error; err != nil {
		return nil, fmt.Errorf("创建资源失败: %v", err)
	}

	// 加载关联的上传者信息
	if resource.UploaderID != nil {
		s.db.Preload("Uploader").First(resource, resource.ID)
	}

	return resource, nil
}

// GetResourceByID 根据ID获取资源
func (s *ResourceService) GetResourceByID(id uint) (*models.Resource, error) {
	var resource models.Resource
	err := s.db.Preload("Uploader").First(&resource, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("资源不存在")
		}
		return nil, fmt.Errorf("获取资源失败: %v", err)
	}
	return &resource, nil
}

// GetResources 获取资源列表
func (s *ResourceService) GetResources(req *GetResourcesRequest) (*ResourceResponse, error) {
	var resources []models.Resource
	var total int64

	query := s.db.Model(&models.Resource{})

	// 状态过滤
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 分类过滤
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}

	// 类型过滤
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	// 关键词搜索
	if req.Keyword != "" {
		keyword := "%" + strings.TrimSpace(req.Keyword) + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", keyword, keyword)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取资源总数失败: %v", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.Size
	err := query.Preload("Uploader").
		Order("created_at DESC").
		Offset(offset).
		Limit(req.Size).
		Find(&resources).Error

	if err != nil {
		return nil, fmt.Errorf("获取资源列表失败: %v", err)
	}

	totalPages := int((total + int64(req.Size) - 1) / int64(req.Size))

	return &ResourceResponse{
		Items:      resources,
		Total:      total,
		Page:       req.Page,
		Size:       req.Size,
		TotalPages: totalPages,
	}, nil
}

// GetPublishedResources 获取已发布的资源列表（前端使用）
func (s *ResourceService) GetPublishedResources(req *GetResourcesRequest) (*ResourceResponse, error) {
	// 只返回已发布的资源
	req.Status = "1"
	return s.GetResources(req)
}

// UpdateResource 更新资源
func (s *ResourceService) UpdateResource(id uint, req *UpdateResourceRequest) (*models.Resource, error) {
	var resource models.Resource
	if err := s.db.First(&resource, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("资源不存在")
		}
		return nil, fmt.Errorf("获取资源失败: %v", err)
	}

	// 验证资源类型
	if req.Type != "" && !isValidResourceType(req.Type) {
		return nil, errors.New("无效的资源类型")
	}

	// 验证资源分类
	if req.Category != "" && !isValidResourceCategory(req.Category) {
		return nil, errors.New("无效的资源分类")
	}

	// 更新字段
	if req.Title != "" {
		resource.Title = req.Title
	}
	if req.Description != "" {
		resource.Description = req.Description
	}
	if req.Type != "" {
		resource.Type = req.Type
	}
	if req.Category != "" {
		resource.Category = req.Category
	}
	if req.Image != "" {
		resource.Image = req.Image
	}
	if req.FileURL != "" {
		resource.FileURL = req.FileURL
	}
	if req.ButtonText != "" {
		resource.ButtonText = req.ButtonText
	}
	resource.Status = req.Status

	if err := s.db.Save(&resource).Error; err != nil {
		return nil, fmt.Errorf("更新资源失败: %v", err)
	}

	// 加载关联的上传者信息
	if resource.UploaderID != nil {
		s.db.Preload("Uploader").First(&resource, resource.ID)
	}

	return &resource, nil
}

// UpdateResourceStatus 更新资源状态
func (s *ResourceService) UpdateResourceStatus(id uint, status int) error {
	var resource models.Resource
	if err := s.db.First(&resource, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("资源不存在")
		}
		return fmt.Errorf("获取资源失败: %v", err)
	}

	if err := s.db.Model(&resource).Update("status", status).Error; err != nil {
		return fmt.Errorf("更新资源状态失败: %v", err)
	}

	return nil
}

// DeleteResource 删除资源
func (s *ResourceService) DeleteResource(id uint) error {
	var resource models.Resource
	if err := s.db.First(&resource, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("资源不存在")
		}
		return fmt.Errorf("获取资源失败: %v", err)
	}

	if err := s.db.Delete(&resource).Error; err != nil {
		return fmt.Errorf("删除资源失败: %v", err)
	}

	return nil
}

// IncrementDownloadCount 增加下载次数
func (s *ResourceService) IncrementDownloadCount(id uint) error {
	err := s.db.Model(&models.Resource{}).Where("id = ?", id).
		Update("download_count", gorm.Expr("download_count + 1")).Error
	if err != nil {
		return fmt.Errorf("更新下载次数失败: %v", err)
	}
	return nil
}

// GetResourceStats 获取资源统计
func (s *ResourceService) GetResourceStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总资源数
	var totalCount int64
	if err := s.db.Model(&models.Resource{}).Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("获取资源总数失败: %v", err)
	}
	stats["total"] = totalCount

	// 已发布资源数
	var publishedCount int64
	if err := s.db.Model(&models.Resource{}).Where("status = ?", models.ResourceStatusApproved).Count(&publishedCount).Error; err != nil {
		return nil, fmt.Errorf("获取已发布资源数失败: %v", err)
	}
	stats["published"] = publishedCount

	// 待审核资源数
	var pendingCount int64
	if err := s.db.Model(&models.Resource{}).Where("status = ?", models.ResourceStatusPending).Count(&pendingCount).Error; err != nil {
		return nil, fmt.Errorf("获取待审核资源数失败: %v", err)
	}
	stats["pending"] = pendingCount

	// 按类型统计
	var typeStats []map[string]interface{}
	err := s.db.Model(&models.Resource{}).
		Select("type, count(*) as count").
		Group("type").
		Scan(&typeStats).Error
	if err != nil {
		return nil, fmt.Errorf("获取类型统计失败: %v", err)
	}
	stats["by_type"] = typeStats

	// 按分类统计
	var categoryStats []map[string]interface{}
	err = s.db.Model(&models.Resource{}).
		Select("category, count(*) as count").
		Group("category").
		Scan(&categoryStats).Error
	if err != nil {
		return nil, fmt.Errorf("获取分类统计失败: %v", err)
	}
	stats["by_category"] = categoryStats

	return stats, nil
}

// 辅助函数：验证资源类型
func isValidResourceType(resourceType string) bool {
	validTypes := []string{
		models.ResourceTypeEBook,
		models.ResourceTypeVideo,
		models.ResourceTypeTool,
	}
	for _, validType := range validTypes {
		if resourceType == validType {
			return true
		}
	}
	return false
}

// 辅助函数：验证资源分类
func isValidResourceCategory(category string) bool {
	validCategories := []string{
		models.ResourceCategoryEBooks,
		models.ResourceCategoryVideos,
		models.ResourceCategoryTools,
	}
	for _, validCategory := range validCategories {
		if category == validCategory {
			return true
		}
	}
	return false
}