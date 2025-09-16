package services

import (
	"errors"
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"

	"godad-backend/config"
	"godad-backend/models"

	"gorm.io/gorm"
)

// ArticleService 文章服务
type ArticleService struct {
	db           *gorm.DB
	cacheService *CacheService
	pointsService *PointsService
	likeService  *LikeService
}

// NewArticleService 创建文章服务实例
func NewArticleService() *ArticleService {
	return &ArticleService{
		db:           config.GetDB(),
		cacheService: NewCacheService(),
		pointsService: NewPointsService(config.GetDB()),
		likeService:  NewLikeService(config.GetDB()),
	}
}

// CreateArticle 创建文章
func (s *ArticleService) CreateArticle(userID uint, req *models.ArticleCreateRequest) (*models.Article, error) {
	// 验证输入
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// 验证分类是否存在
	if req.CategoryID > 0 {
		var category models.Category
		if err := s.db.Where("id = ? AND status = ?", req.CategoryID, 1).First(&category).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("分类不存在")
			}
			return nil, fmt.Errorf("查询分类失败: %v", err)
		}
	}

	// 生成唯一的slug
	slug, err := s.generateUniqueSlug(req.Slug, req.Title)
	if err != nil {
		return nil, fmt.Errorf("生成文章别名失败: %v", err)
	}

	// 创建文章
	article := &models.Article{
		Title:       req.Title,
		Slug:        slug,
		Summary:     req.Summary,
		Content:     req.Content,
		CoverImage:  req.CoverImage,
		AuthorID:    userID,
		CategoryID:  req.CategoryID,
		Tags:        req.Tags,
		IsTop:       req.IsTop,
		IsRecommend: req.IsRecommend,
		Status:      req.Status,
	}

	// 如果没有提供摘要，自动生成
	if article.Summary == "" {
		article.Summary = s.generateSummary(article.Content)
	}

	// 保存到数据库
	if err := s.db.Create(article).Error; err != nil {
		return nil, fmt.Errorf("创建文章失败: %v", err)
	}

	// 预加载关联数据
	if err := s.db.Preload("Author").Preload("Category").First(article, article.ID).Error; err != nil {
		return nil, fmt.Errorf("加载文章数据失败: %v", err)
	}

	// 清理相关缓存
	s.cacheService.DeletePattern("articles:*")
	s.cacheService.DeletePattern("search:*")

	// 奖励发布文章积分
	go func() {
		err := s.pointsService.AwardPoints(userID, "publish_article", "article", article.ID, "发布文章")
		if err != nil {
			fmt.Printf("发布文章积分奖励失败: %v\n", err)
		}
	}()

	return article, nil
}

// UpdateArticle 更新文章
func (s *ArticleService) UpdateArticle(articleID, userID uint, req *models.ArticleUpdateRequest) (*models.Article, error) {
	// 获取文章
	article, err := s.GetArticleByID(articleID, userID, true)
	if err != nil {
		return nil, err
	}

	// 检查权限（作者或管理员可以编辑）
	userService := NewUserService()
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// 只有作者或管理员可以编辑
	if article.AuthorID != userID && user.Role != 2 {
		return nil, errors.New("无权限编辑此文章")
	}

	// 验证输入
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// 验证分类是否存在
	if req.CategoryID > 0 {
		var category models.Category
		if err := s.db.Where("id = ? AND status = ?", req.CategoryID, 1).First(&category).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("分类不存在")
			}
			return nil, fmt.Errorf("查询分类失败: %v", err)
		}
	}

	// 更新字段
	updateData := make(map[string]interface{})
	if req.Title != "" {
		updateData["title"] = req.Title
	}
	if req.Content != "" {
		updateData["content"] = req.Content
		updateData["summary"] = req.Summary
		if req.Summary == "" {
			updateData["summary"] = s.generateSummary(req.Content)
		}
	}
	if req.CoverImage != "" {
		updateData["cover_image"] = req.CoverImage
	}
	if req.CategoryID > 0 {
		updateData["category_id"] = req.CategoryID
	}
	if req.Tags != "" {
		updateData["tags"] = s.processTags(req.Tags)
	}
	if req.Status != nil {
		updateData["status"] = *req.Status
	}
	updateData["is_top"] = req.IsTop

	// 执行更新
	if err := s.db.Model(article).Updates(updateData).Error; err != nil {
		return nil, fmt.Errorf("更新文章失败: %v", err)
	}

	// 清理相关缓存
	s.cacheService.Delete(fmt.Sprintf("article:%d", articleID))
	s.cacheService.DeletePattern("articles:*")
	s.cacheService.DeletePattern("search:*")

	// 重新加载文章数据
	return s.GetArticleByID(articleID, userID, true)
}

// DeleteArticle 删除文章（软删除）
func (s *ArticleService) DeleteArticle(articleID, userID uint) error {
	// 获取文章
	article, err := s.GetArticleByID(articleID, userID, true)
	if err != nil {
		return err
	}

	// 检查权限（作者或管理员可以删除）
	// 获取用户信息检查是否为管理员
	userService := NewUserService()
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return err
	}

	// 只有作者或管理员可以删除
	if article.AuthorID != userID && user.Role != 2 {
		return errors.New("无权限删除此文章")
	}

	// 软删除
	if err := s.db.Delete(article).Error; err != nil {
		return fmt.Errorf("删除文章失败: %v", err)
	}

	// 清理相关缓存
	s.cacheService.Delete(fmt.Sprintf("article:%d", articleID))
	s.cacheService.DeletePattern("articles:*")
	s.cacheService.DeletePattern("search:*")

	return nil
}

// GetArticleByID 根据ID获取文章
func (s *ArticleService) GetArticleByID(articleID, userID uint, checkOwner bool) (*models.Article, error) {
	// 尝试从缓存获取（仅对非作者查看的情况）
	if !checkOwner {
		if cachedArticle, err := s.cacheService.GetArticle(articleID); err == nil {
			// 增加浏览量缓存
			s.cacheService.IncreaseViewCount(articleID)
			// 更新数据库浏览量（异步）
			go func() {
				s.db.Model(&models.Article{}).Where("id = ?", articleID).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
			}()
			return cachedArticle, nil
		}
	}

	var article models.Article
	query := s.db.Preload("Author").Preload("Category")

	// 如果不是检查所有者，只查询已发布的文章
	if !checkOwner {
		query = query.Where("status = ?", 1)
	} else if userID > 0 {
		// 如果是检查所有者，可以查看自己的所有文章
		query = query.Where("author_id = ? OR status = ?", userID, 1)
	}

	if err := query.First(&article, articleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, fmt.Errorf("查询文章失败: %v", err)
	}

	// 如果不是作者查看，增加浏览量
	if !checkOwner || article.AuthorID != userID {
		s.db.Model(&article).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
		article.ViewCount++
		// 缓存文章（仅对已发布的文章）
		if article.Status == 1 {
			s.cacheService.SetArticle(articleID, &article, 30*time.Minute)
		}
	}

	return &article, nil
}

// GetArticleList 获取文章列表
func (s *ArticleService) GetArticleList(req *models.ArticleListRequest) ([]*models.Article, int64, error) {
	// 生成缓存key
	cacheKey := fmt.Sprintf("articles:list:%d:%d:%d:%d:%s:%s", req.CategoryID, req.AuthorID, req.Status, req.Page, req.Keyword, req.Sort)
	
	// 尝试从缓存获取（仅对默认列表查询）
	if req.Status == 0 || req.Status == 1 {
		if cachedArticles, err := s.cacheService.GetArticleList(cacheKey); err == nil {
			// 从数据库获取总数（缓存可能不准确）
			var total int64
			countQuery := s.db.Model(&models.Article{})
			if req.Status != 0 {
				countQuery = countQuery.Where("status = ?", req.Status)
			} else {
				countQuery = countQuery.Where("status = ?", 1)
			}
			if req.CategoryID > 0 {
				countQuery = countQuery.Where("category_id = ?", req.CategoryID)
			}
			if req.AuthorID > 0 {
				countQuery = countQuery.Where("author_id = ?", req.AuthorID)
			}
			if req.Keyword != "" {
				countQuery = countQuery.Where("title LIKE ? OR content LIKE ? OR tags LIKE ?",
					"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
			}
			countQuery.Count(&total)
			
			// 转换为指针切片
			articlesPtr := make([]*models.Article, len(cachedArticles))
			for i := range cachedArticles {
				articlesPtr[i] = &cachedArticles[i]
			}
			return articlesPtr, total, nil
		}
	}

	var articles []*models.Article
	var total int64

	query := s.db.Model(&models.Article{}).Preload("Author").Preload("Category")

	// 状态过滤
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	} else {
		// 默认只显示已发布的文章
		query = query.Where("status = ?", 1)
	}

	// 分类过滤
	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	// 作者过滤
	if req.AuthorID > 0 {
		query = query.Where("author_id = ?", req.AuthorID)
	}

	// 关键词搜索
	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ? OR tags LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}



	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取文章总数失败: %v", err)
	}

	// 排序
	orderBy := "created_at DESC"
	switch req.Sort {
	case "view_count":
		orderBy = "view_count DESC, created_at DESC"
	case "like_count":
		orderBy = "like_count DESC, created_at DESC"
	case "comment_count":
		orderBy = "comment_count DESC, created_at DESC"
	case "updated_at":
		orderBy = "updated_at DESC"
	}

	// 置顶文章优先
	query = query.Order("is_top DESC, " + orderBy)

	// 分页查询
	offset := (req.Page - 1) * req.Size
	if err := query.Offset(offset).Limit(req.Size).Find(&articles).Error; err != nil {
		return nil, 0, fmt.Errorf("获取文章列表失败: %v", err)
	}

	// 缓存结果（仅对已发布文章的默认查询）
	if (req.Status == 0 || req.Status == 1) && len(articles) > 0 {
		// 转换为值切片进行缓存
		articleValues := make([]models.Article, len(articles))
		for i, article := range articles {
			articleValues[i] = *article
		}
		s.cacheService.SetArticleList(cacheKey, articleValues, 15*time.Minute)
	}

	return articles, total, nil
}

// GetUserArticles 获取用户的文章列表
func (s *ArticleService) GetUserArticles(userID uint, page, size int, status string) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	query := s.db.Model(&models.Article{}).Preload("Category").Where("author_id = ?", userID)

	// 状态过滤
	if status != "" {
		var statusCode int8
		switch status {
		case "draft":
			statusCode = 0
		case "published":
			statusCode = 1
		default:
			// 如果是数字字符串，尝试解析
			if status == "0" {
				statusCode = 0
			} else if status == "1" {
				statusCode = 1
			} else {
				return nil, 0, fmt.Errorf("无效的状态参数: %s", status)
			}
		}
		query = query.Where("status = ?", statusCode)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户文章总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&articles).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户文章列表失败: %v", err)
	}

	return articles, total, nil
}

// LikeArticle 点赞文章
func (s *ArticleService) LikeArticle(articleID, userID uint) error {
	// 检查文章是否存在
	var article models.Article
	if err := s.db.Where("id = ? AND status = ?", articleID, 1).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return fmt.Errorf("查询文章失败: %v", err)
	}

	// 使用 LikeService 来处理点赞逻辑，这样会自动记录点赞记录并奖励积分
	_, err := s.likeService.ToggleLike(userID, "article", articleID)
	if err != nil {
		return fmt.Errorf("点赞失败: %v", err)
	}

	return nil
}

// validateCreateRequest 验证创建请求
func (s *ArticleService) validateCreateRequest(req *models.ArticleCreateRequest) error {
	if req.Title == "" {
		return errors.New("文章标题不能为空")
	}
	if len(req.Title) > 200 {
		return errors.New("文章标题不能超过200个字符")
	}
	if req.Content == "" {
		return errors.New("文章内容不能为空")
	}
	if len(req.Content) > 100000 {
		return errors.New("文章内容不能超过100000个字符")
	}
	if req.Status == 0 {
		req.Status = 0 // 默认为草稿
	}
	if req.Status < 0 || req.Status > 2 {
		return errors.New("文章状态只能是0(草稿)、1(已发布)或2(已下架)")
	}
	return nil
}

// validateUpdateRequest 验证更新请求
func (s *ArticleService) validateUpdateRequest(req *models.ArticleUpdateRequest) error {
	if req.Title != "" && len(req.Title) > 200 {
		return errors.New("文章标题不能超过200个字符")
	}
	if req.Content != "" && len(req.Content) > 100000 {
		return errors.New("文章内容不能超过100000个字符")
	}
	if req.Status != nil && (*req.Status < 0 || *req.Status > 2) {
		return errors.New("文章状态只能是0(草稿)、1(已发布)或2(已下架)")
	}
	return nil
}

// processTags 处理标签
func (s *ArticleService) processTags(tags string) string {
	if tags == "" {
		return ""
	}
	// 分割标签，去重，去空格
	tagList := strings.Split(tags, ",")
	tagMap := make(map[string]bool)
	var result []string

	for _, tag := range tagList {
		tag = strings.TrimSpace(tag)
		if tag != "" && !tagMap[tag] {
			tagMap[tag] = true
			result = append(result, tag)
		}
	}

	return strings.Join(result, ",")
}

// generateSummary 生成文章摘要
func (s *ArticleService) generateSummary(content string) string {
	// 清理HTML标签和实体
	cleanContent := s.cleanHTMLContent(content)
	
	// 取前200个字符
	runes := []rune(cleanContent)
	if len(runes) <= 200 {
		return cleanContent
	}
	return string(runes[:200]) + "..."
}

// generateUniqueSlug 生成唯一的slug
func (s *ArticleService) generateUniqueSlug(slug, title string) (string, error) {
	// 如果没有提供slug，基于标题生成
	if slug == "" {
		slug = s.generateSlugFromTitle(title)
	}

	// 检查slug是否已存在
	originalSlug := slug
	counter := 0

	for {
		var count int64
		err := s.db.Model(&models.Article{}).Where("slug = ?", slug).Count(&count).Error
		if err != nil {
			return "", fmt.Errorf("检查slug重复失败: %v", err)
		}

		// 如果不存在重复，返回当前slug
		if count == 0 {
			return slug, nil
		}

		// 如果存在重复，添加数字后缀
		counter++
		slug = fmt.Sprintf("%s-%d", originalSlug, counter)
	}
}

// generateSlugFromTitle 从标题生成slug
func (s *ArticleService) generateSlugFromTitle(title string) string {
	// 简单的slug生成：转换为小写，替换特殊字符为短横线
	slug := strings.ToLower(title)
	slug = regexp.MustCompile(`[^a-z0-9\p{Han}]+`).ReplaceAllString(slug, "-")
	slug = regexp.MustCompile(`^-+|-+$`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	
	// 如果slug为空或只有短横线，使用时间戳作为备选
	if slug == "" || slug == "-" {
		slug = fmt.Sprintf("article-%d", time.Now().Unix())
	}
	
	return slug
}

// GetArticleCount 获取文章总数
func (s *ArticleService) GetArticleCount() (int64, error) {
	var count int64
	err := s.db.Model(&models.Article{}).Where("deleted_at IS NULL").Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("获取文章总数失败: %v", err)
	}
	return count, nil
}

// cleanHTMLContent 清理HTML标签和实体
func (s *ArticleService) cleanHTMLContent(content string) string {
	// 1. 解码HTML实体（如&nbsp;转为空格）
	content = html.UnescapeString(content)
	
	// 2. 移除HTML标签
	// 匹配所有HTML标签
	htmlTagReg := regexp.MustCompile(`<[^>]*>`)
	content = htmlTagReg.ReplaceAllString(content, "")
	
	// 3. 将连续的空白字符替换为单个空格
	spaceReg := regexp.MustCompile(`\s+`)
	content = spaceReg.ReplaceAllString(content, " ")
	
	// 4. 去除首尾空白
	content = strings.TrimSpace(content)
	
	return content
}

// GetHotArticles 获取热门文章
func (s *ArticleService) GetHotArticles(period string, limit int) ([]*models.Article, error) {
	var articles []*models.Article

	// 构建基本查询
	query := s.db.Preload("Author").Preload("Category").
		Where("status = ? AND deleted_at IS NULL", 1)

	// 根据时间周期添加条件
	switch period {
	case "today":
		// 今日热门
		today := time.Now().Format("2006-01-02")
		query = query.Where("DATE(created_at) = ?", today)
	case "week":
		// 本周热门
		weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1)
		weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
		query = query.Where("created_at >= ?", weekStart)
	case "month":
		// 本月热门
		monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
		query = query.Where("created_at >= ?", monthStart)
	case "all":
		// 全部时间，不添加时间限制
	default:
		// 默认为今日热门
		today := time.Now().Format("2006-01-02")
		query = query.Where("DATE(created_at) = ?", today)
	}

	// 按热度值排序：浏览量*0.6 + 点赞数*0.3 + 评论数*0.1
	// 使用原生SQL计算热度值并排序
	err := query.
		Select("articles.*, (view_count * 0.6 + like_count * 0.3 + comment_count * 0.1) as heat_score").
		Order("heat_score DESC").
		Limit(limit).
		Find(&articles).Error

	if err != nil {
		return nil, fmt.Errorf("查询热门文章失败: %v", err)
	}

	return articles, nil
}