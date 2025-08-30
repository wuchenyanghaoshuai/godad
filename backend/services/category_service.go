package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"godad-backend/config"
	"godad-backend/models"

	"gorm.io/gorm"
)

// CategoryService 分类服务
type CategoryService struct {
	db           *gorm.DB
	cacheService *CacheService
}

// NewCategoryService 创建分类服务实例
func NewCategoryService() *CategoryService {
	return &CategoryService{
		db:           config.GetDB(),
		cacheService: NewCacheService(),
	}
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(req *models.CategoryCreateRequest) (*models.Category, error) {
	// 验证输入
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// 检查名称是否已存在
	var existingCategory models.Category
	if err := s.db.Where("name = ?", req.Name).First(&existingCategory).Error; err == nil {
		return nil, errors.New("分类名称已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("检查分类名称失败: %v", err)
	}

	// 检查别名是否已存在
	if req.Slug != "" {
		if err := s.db.Where("slug = ?", req.Slug).First(&existingCategory).Error; err == nil {
			return nil, errors.New("分类别名已存在")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("检查分类别名失败: %v", err)
		}
	} else {
		// 如果没有提供别名，自动生成
		req.Slug = s.generateSlug(req.Name)
	}

	// 获取最大排序值
	var maxSort int
	s.db.Model(&models.Category{}).Select("COALESCE(MAX(sort), 0)").Scan(&maxSort)

	// 创建分类
	category := &models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
		Sort:        maxSort + 1,
		Status:      1, // 默认启用
	}

	// 保存到数据库
	if err := s.db.Create(category).Error; err != nil {
		return nil, fmt.Errorf("创建分类失败: %v", err)
	}

	// 清理分类缓存
	s.cacheService.Delete("categories")

	return category, nil
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(categoryID uint, req *models.CategoryUpdateRequest) (*models.Category, error) {
	// 获取分类
	category, err := s.GetCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}

	// 验证输入
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// 检查名称是否已存在（排除自己）
	if req.Name != "" && req.Name != category.Name {
		var existingCategory models.Category
		if err := s.db.Where("name = ? AND id != ?", req.Name, categoryID).First(&existingCategory).Error; err == nil {
			return nil, errors.New("分类名称已存在")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("检查分类名称失败: %v", err)
		}
	}

	// 检查别名是否已存在（排除自己）
	if req.Slug != "" && req.Slug != category.Slug {
		var existingCategory models.Category
		if err := s.db.Where("slug = ? AND id != ?", req.Slug, categoryID).First(&existingCategory).Error; err == nil {
			return nil, errors.New("分类别名已存在")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("检查分类别名失败: %v", err)
		}
	}

	// 更新字段
	updateData := make(map[string]interface{})
	if req.Name != "" {
		updateData["name"] = req.Name
	}
	if req.Slug != "" {
		updateData["slug"] = req.Slug
	}
	if req.Description != "" {
		updateData["description"] = req.Description
	}
	if req.Icon != "" {
		updateData["icon"] = req.Icon
	}
	if req.Color != "" {
		updateData["color"] = req.Color
	}
	if req.Sort > 0 {
		updateData["sort"] = req.Sort
	}
	if req.Status >= 0 {
		updateData["status"] = req.Status
	}

	// 执行更新
	if err := s.db.Model(category).Updates(updateData).Error; err != nil {
		return nil, fmt.Errorf("更新分类失败: %v", err)
	}

	// 清理分类缓存
	s.cacheService.Delete("categories")

	// 重新加载分类数据
	return s.GetCategoryByID(categoryID)
}

// DeleteCategory 删除分类
func (s *CategoryService) DeleteCategory(categoryID uint) error {
	// 获取分类
	category, err := s.GetCategoryByID(categoryID)
	if err != nil {
		return err
	}

	// 检查是否有文章使用此分类
	var articleCount int64
	if err := s.db.Model(&models.Article{}).Where("category_id = ?", categoryID).Count(&articleCount).Error; err != nil {
		return fmt.Errorf("检查分类使用情况失败: %v", err)
	}

	if articleCount > 0 {
		return errors.New("该分类下还有文章，无法删除")
	}

	// 删除分类
	if err := s.db.Delete(category).Error; err != nil {
		return fmt.Errorf("删除分类失败: %v", err)
	}

	// 清理分类缓存
	s.cacheService.Delete("categories")

	return nil
}

// GetCategoryByID 根据ID获取分类
func (s *CategoryService) GetCategoryByID(categoryID uint) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, fmt.Errorf("查询分类失败: %v", err)
	}

	return &category, nil
}

// GetCategoryBySlug 根据别名获取分类
func (s *CategoryService) GetCategoryBySlug(slug string) (*models.Category, error) {
	var category models.Category
	if err := s.db.Where("slug = ? AND status = ?", slug, 1).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, fmt.Errorf("查询分类失败: %v", err)
	}

	return &category, nil
}

// GetCategoryList 获取分类列表
func (s *CategoryService) GetCategoryList(status int, page, size int) ([]*models.Category, int64, error) {
	var categories []*models.Category
	var total int64

	query := s.db.Model(&models.Category{})

	// 状态过滤
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取分类总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("sort ASC, created_at DESC").Find(&categories).Error; err != nil {
		return nil, 0, fmt.Errorf("获取分类列表失败: %v", err)
	}

	return categories, total, nil
}

// GetAllCategories 获取所有启用的分类（不分页）
func (s *CategoryService) GetAllCategories() ([]*models.Category, error) {
	// 尝试从缓存获取
	if cachedCategories, err := s.cacheService.GetCategories(); err == nil {
		// 转换为指针切片
		categoriesPtr := make([]*models.Category, len(cachedCategories))
		for i := range cachedCategories {
			categoriesPtr[i] = &cachedCategories[i]
		}
		return categoriesPtr, nil
	}

	var categories []*models.Category
	if err := s.db.Where("status = ?", 1).Order("sort ASC, created_at DESC").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("获取分类列表失败: %v", err)
	}

	// 缓存结果
	if len(categories) > 0 {
		// 转换为值切片进行缓存
		categoryValues := make([]models.Category, len(categories))
		for i, category := range categories {
			categoryValues[i] = *category
		}
		s.cacheService.SetCategories(categoryValues, 1*time.Hour)
	}

	return categories, nil
}

// GetCategoriesWithArticleCount 获取分类及其文章数量
func (s *CategoryService) GetCategoriesWithArticleCount() ([]*models.CategoryResponse, error) {
	var categories []*models.Category
	if err := s.db.Where("status = ?", 1).Order("sort ASC, created_at DESC").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("获取分类列表失败: %v", err)
	}

	var result []*models.CategoryResponse
	for _, category := range categories {
		// 获取该分类下的文章数量
		var articleCount int64
		s.db.Model(&models.Article{}).Where("category_id = ? AND status = ?", category.ID, 1).Count(&articleCount)

		categoryResp := category.ToResponse()
		categoryResp.ArticleCount = articleCount
		result = append(result, categoryResp)
	}

	return result, nil
}

// UpdateCategorySort 更新分类排序
func (s *CategoryService) UpdateCategorySort(categoryID uint, sort int) error {
	// 检查分类是否存在
	if _, err := s.GetCategoryByID(categoryID); err != nil {
		return err
	}

	// 更新排序
	if err := s.db.Model(&models.Category{}).Where("id = ?", categoryID).Update("sort", sort).Error; err != nil {
		return fmt.Errorf("更新分类排序失败: %v", err)
	}

	return nil
}

// ToggleCategoryStatus 切换分类状态
func (s *CategoryService) ToggleCategoryStatus(categoryID uint) (*models.Category, error) {
	// 获取分类
	category, err := s.GetCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}

	// 切换状态
	newStatus := 1
	if category.Status == 1 {
		newStatus = 0
	}

	// 如果要禁用分类，检查是否有文章使用
	if newStatus == 0 {
		var articleCount int64
		if err := s.db.Model(&models.Article{}).Where("category_id = ? AND status = ?", categoryID, 1).Count(&articleCount).Error; err != nil {
			return nil, fmt.Errorf("检查分类使用情况失败: %v", err)
		}
		if articleCount > 0 {
			return nil, errors.New("该分类下还有已发布的文章，无法禁用")
		}
	}

	// 更新状态
	if err := s.db.Model(category).Update("status", newStatus).Error; err != nil {
		return nil, fmt.Errorf("更新分类状态失败: %v", err)
	}

	// 清理分类缓存
	s.cacheService.Delete("categories")

	category.Status = int8(newStatus)
	return category, nil
}

// validateCreateRequest 验证创建请求
func (s *CategoryService) validateCreateRequest(req *models.CategoryCreateRequest) error {
	if req.Name == "" {
		return errors.New("分类名称不能为空")
	}
	if len(req.Name) > 50 {
		return errors.New("分类名称不能超过50个字符")
	}
	if req.Slug != "" && len(req.Slug) > 50 {
		return errors.New("分类别名不能超过50个字符")
	}
	if len(req.Description) > 200 {
		return errors.New("分类描述不能超过200个字符")
	}
	return nil
}

// validateUpdateRequest 验证更新请求
func (s *CategoryService) validateUpdateRequest(req *models.CategoryUpdateRequest) error {
	if req.Name != "" && len(req.Name) > 50 {
		return errors.New("分类名称不能超过50个字符")
	}
	if req.Slug != "" && len(req.Slug) > 50 {
		return errors.New("分类别名不能超过50个字符")
	}
	if len(req.Description) > 200 {
		return errors.New("分类描述不能超过200个字符")
	}
	if req.Status < 0 || req.Status > 1 {
		return errors.New("分类状态只能是0或1")
	}
	return nil
}

// generateSlug 生成别名
func (s *CategoryService) generateSlug(name string) string {
	// 简单的别名生成：转换为小写，替换空格为连字符
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	
	// 移除特殊字符，只保留字母、数字和连字符
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// GetCategoryCount 获取分类总数
func (s *CategoryService) GetCategoryCount() (int64, error) {
	var count int64
	err := s.db.Model(&models.Category{}).Where("deleted_at IS NULL").Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("获取分类总数失败: %v", err)
	}
	return count, nil
}