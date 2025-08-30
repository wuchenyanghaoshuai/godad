package services

import (
	"fmt"
	"strings"
	"time"

	"godad-backend/config"
	"godad-backend/models"

	"gorm.io/gorm"
)

type SearchService struct {
	db           *gorm.DB
	cacheService *CacheService
}

func NewSearchService() *SearchService {
	return &SearchService{
		db:           config.GetDB(),
		cacheService: NewCacheService(),
	}
}

type SearchRequest struct {
	Keyword    string `json:"keyword" form:"keyword"`
	CategoryID uint   `json:"category_id" form:"category_id"`
	Sort       string `json:"sort" form:"sort"`
	Page       int    `json:"page" form:"page"`
	Size       int    `json:"size" form:"size"`
}

type SearchResult struct {
	Articles    []*models.Article `json:"articles"`
	Total       int64             `json:"total"`
	Keyword     string            `json:"keyword"`
	Suggestions []string          `json:"suggestions"`
}

func (s *SearchService) Search(req *SearchRequest) (*SearchResult, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.Size > 50 {
		req.Size = 50
	}

	keyword := strings.TrimSpace(req.Keyword)
	if keyword == "" {
		return &SearchResult{
			Articles: []*models.Article{},
			Total:    0,
			Keyword:  keyword,
		}, nil
	}

	// 记录搜索关键词
	go func() {
		s.cacheService.RecordSearchKeyword(keyword)
	}()

	// 生成缓存key
	cacheKey := fmt.Sprintf("search:%s:%d:%s:%d", keyword, req.CategoryID, req.Sort, req.Page)
	
	// 尝试从缓存获取搜索结果
	if cachedArticles, err := s.cacheService.GetSearchResults(cacheKey); err == nil {
		// 获取总数（可能需要实时查询以保证准确性）
		var total int64
		countQuery := s.buildSearchQuery(keyword, req.CategoryID)
		countQuery.Count(&total)
		
		// 转换为指针切片
		articlesPtr := make([]*models.Article, len(cachedArticles))
		for i := range cachedArticles {
			articlesPtr[i] = &cachedArticles[i]
		}

		return &SearchResult{
			Articles: articlesPtr,
			Total:    total,
			Keyword:  keyword,
		}, nil
	}

	var articles []*models.Article
	var total int64

	// 构建搜索查询
	query := s.buildSearchQuery(keyword, req.CategoryID)
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取搜索结果总数失败: %v", err)
	}

	// 应用排序
	orderBy := s.buildOrderBy(req.Sort)
	query = query.Order(orderBy)

	// 分页查询
	offset := (req.Page - 1) * req.Size
	if err := query.Offset(offset).Limit(req.Size).Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("搜索失败: %v", err)
	}

	// 缓存搜索结果
	if len(articles) > 0 {
		// 转换为值切片进行缓存
		articleValues := make([]models.Article, len(articles))
		for i, article := range articles {
			articleValues[i] = *article
		}
		s.cacheService.SetSearchResults(cacheKey, articleValues, 10*time.Minute)
	}

	// 获取搜索建议
	suggestions := s.getSearchSuggestions(keyword)

	return &SearchResult{
		Articles:    articles,
		Total:       total,
		Keyword:     keyword,
		Suggestions: suggestions,
	}, nil
}

func (s *SearchService) buildSearchQuery(keyword string, categoryID uint) *gorm.DB {
	query := s.db.Model(&models.Article{}).
		Preload("Author").
		Preload("Category").
		Where("status = ?", 1)

	// 关键词搜索
	if keyword != "" {
		// 支持标题、内容、标签搜索
		searchCondition := "title LIKE ? OR content LIKE ? OR tags LIKE ?"
		searchValue := "%" + keyword + "%"
		query = query.Where(searchCondition, searchValue, searchValue, searchValue)
	}

	// 分类筛选
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	return query
}

func (s *SearchService) buildOrderBy(sort string) string {
	switch sort {
	case "view_count":
		return "view_count DESC, created_at DESC"
	case "like_count":
		return "like_count DESC, created_at DESC"
	case "comment_count":
		return "comment_count DESC, created_at DESC"
	case "updated_at":
		return "updated_at DESC"
	default:
		return "created_at DESC"
	}
}

func (s *SearchService) getSearchSuggestions(keyword string) []string {
	// 基于关键词生成简单的搜索建议
	suggestions := []string{}
	
	// 这里可以实现更复杂的搜索建议逻辑
	// 比如基于历史搜索、热门关键词等
	
	if len(keyword) >= 2 {
		// 查找相似的文章标题
		var titles []string
		s.db.Model(&models.Article{}).
			Where("status = ? AND title LIKE ?", 1, "%"+keyword+"%").
			Limit(5).
			Pluck("title", &titles)
		
		suggestions = append(suggestions, titles...)
	}
	
	return suggestions
}

func (s *SearchService) GetHotKeywords(limit int) ([]string, error) {
	if limit <= 0 {
		limit = 10
	}

	// 尝试从缓存获取热门关键词
	if keywords, err := s.cacheService.GetHotKeywords(); err == nil {
		if len(keywords) > limit {
			return keywords[:limit], nil
		}
		return keywords, nil
	}

	// 从统计数据获取热门搜索关键词
	keywords, err := s.cacheService.GetTopSearchKeywords(int64(limit))
	if err != nil {
		// 如果缓存失败，返回一些默认热门关键词
		return []string{"育儿", "教育", "健康", "营养", "早教"}, nil
	}

	// 缓存热门关键词
	s.cacheService.SetHotKeywords(keywords, 1*time.Hour)

	return keywords, nil
}

func (s *SearchService) GetSearchSuggestions(keyword string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 5
	}
	
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []string{}, nil
	}

	var suggestions []string
	
	// 从文章标题中搜索相关建议
	err := s.db.Model(&models.Article{}).
		Where("status = ? AND title LIKE ?", 1, "%"+keyword+"%").
		Limit(limit).
		Pluck("title", &suggestions)
	
	if err != nil {
		return []string{}, fmt.Errorf("获取搜索建议失败: %v", err)
	}

	// 去重并处理建议
	uniqueSuggestions := s.processSearchSuggestions(suggestions, keyword)
	
	if len(uniqueSuggestions) > limit {
		return uniqueSuggestions[:limit], nil
	}
	
	return uniqueSuggestions, nil
}

func (s *SearchService) processSearchSuggestions(suggestions []string, keyword string) []string {
	seen := make(map[string]bool)
	var result []string
	
	for _, suggestion := range suggestions {
		// 简单处理：提取包含关键词的部分作为建议
		if !seen[suggestion] && strings.Contains(strings.ToLower(suggestion), strings.ToLower(keyword)) {
			seen[suggestion] = true
			result = append(result, suggestion)
		}
	}
	
	return result
}

func (s *SearchService) ClearSearchCache() error {
	return s.cacheService.DeletePattern("search:*")
}