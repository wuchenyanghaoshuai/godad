package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"godad-backend/config"
	"godad-backend/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// CommentService 评论服务
type CommentService struct {
	db                  *gorm.DB
	notificationService *NotificationService
	redisClient         *redis.Client
}

// NewCommentService 创建评论服务实例
func NewCommentService() *CommentService {
	db := config.GetDB()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	
	return &CommentService{
		db:                  db,
		notificationService: NewNotificationService(db),
		redisClient:         redisClient,
	}
}

// CreateComment 创建评论
func (s *CommentService) CreateComment(userID uint, req *models.CommentCreateRequest) (*models.Comment, error) {
	// 验证请求参数
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// 验证文章是否存在
	var article models.Article
	if err := s.db.Where("id = ? AND status = ?", req.ArticleID, 1).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在或已下线")
		}
		return nil, err
	}

	// 如果是回复评论，验证父评论是否存在
	if req.ParentID != nil && *req.ParentID > 0 {
		var parentComment models.Comment
		if err := s.db.Where("id = ? AND article_id = ? AND status = ?", *req.ParentID, req.ArticleID, 1).First(&parentComment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("父评论不存在")
			}
			return nil, err
		}
	}

	// 创建评论
	comment := &models.Comment{
		ArticleID: req.ArticleID,
		UserID:    userID,
		ParentID:  req.ParentID,
		Content:   strings.TrimSpace(req.Content),
		Status:    1, // 默认启用
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 保存评论
	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新文章评论数
	if err := tx.Model(&article).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 如果是回复评论，更新父评论的回复数
	if req.ParentID != nil && *req.ParentID > 0 {
		if err := tx.Model(&models.Comment{}).Where("id = ?", *req.ParentID).UpdateColumn("reply_count", gorm.Expr("reply_count + ?", 1)).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 发送通知
	if req.ParentID != nil && *req.ParentID > 0 {
		// 这是回复评论，需要通知原评论作者
		var parentComment models.Comment
		if err := s.db.Where("id = ?", *req.ParentID).First(&parentComment).Error; err == nil {
			// 通知原评论作者（如果不是自己）
			if parentComment.UserID != userID {
				if err := s.notificationService.CreateCommentReplyNotification(userID, parentComment.UserID, req.ArticleID, *req.ParentID, req.Content); err != nil {
					fmt.Printf("发送回复通知失败: %v\n", err)
				}
			}
		}
		
		// 如果原评论作者不是文章作者，也要通知文章作者
		if article.AuthorID != userID && article.AuthorID != parentComment.UserID {
			if err := s.notificationService.CreateCommentNotification(userID, article.AuthorID, req.ArticleID, req.Content); err != nil {
				fmt.Printf("发送文章评论通知失败: %v\n", err)
			}
		}
	} else {
		// 这是顶级评论，通知文章作者
		if article.AuthorID != userID {
			if err := s.notificationService.CreateCommentNotification(userID, article.AuthorID, req.ArticleID, req.Content); err != nil {
				fmt.Printf("发送评论通知失败: %v\n", err)
			}
		}
	}

    // 预加载关联数据
    if err := s.db.Preload("User").Preload("Article").First(comment, comment.ID).Error; err != nil {
        return nil, err
    }

    // 处理@提及（互相关注限制）
    if len(req.Mentions) > 0 {
        // 去重并过滤
        seen := make(map[uint]struct{})
        for _, mid := range req.Mentions {
            if mid == 0 || mid == userID { continue }
            if _, ok := seen[mid]; ok { continue }
            seen[mid] = struct{}{}
            if s.isMutualFollow(userID, mid) {
                _ = s.notificationService.CreateMentionNotification(userID, mid, req.ArticleID, comment.ID, req.Content)
            }
        }
    }

    return comment, nil
}

// isMutualFollow 判断两个用户是否互相关注
func (s *CommentService) isMutualFollow(a, b uint) bool {
    var cnt int64
    // 存在 a->b
    s.db.Model(&models.Follow{}).Where("follower_id = ? AND followee_id = ? AND deleted_at IS NULL", a, b).Count(&cnt)
    if cnt == 0 { return false }
    cnt = 0
    s.db.Model(&models.Follow{}).Where("follower_id = ? AND followee_id = ? AND deleted_at IS NULL", b, a).Count(&cnt)
    return cnt > 0
}

// UpdateComment 更新评论
func (s *CommentService) UpdateComment(commentID, userID uint, req *models.CommentUpdateRequest) (*models.Comment, error) {
	// 验证请求参数
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// 查找评论
	var comment models.Comment
	if err := s.db.Where("id = ?", commentID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("评论不存在")
		}
		return nil, err
	}

	// 检查权限（只能修改自己的评论）
	if comment.UserID != userID {
		return nil, errors.New("无权限修改此评论")
	}

	// 检查评论状态
	if comment.Status != 1 {
		return nil, errors.New("评论已被删除或禁用")
	}

	// 更新评论内容
	updateData := map[string]interface{}{
		"content":    strings.TrimSpace(req.Content),
		"updated_at": time.Now(),
	}

	if err := s.db.Model(&comment).Updates(updateData).Error; err != nil {
		return nil, err
	}

	// 预加载关联数据
	if err := s.db.Preload("User").Preload("Article").First(&comment, comment.ID).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

// DeleteComment 删除评论（软删除）
func (s *CommentService) DeleteComment(commentID, userID uint) error {
	// 查找评论
	var comment models.Comment
	if err := s.db.Where("id = ?", commentID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		return err
	}

	// 检查权限（只能删除自己的评论）
	if comment.UserID != userID {
		return errors.New("无权限删除此评论")
	}

	// 检查评论状态
	if comment.Status == 0 {
		return errors.New("评论已被删除")
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 计算要删除的总评论数（包括子评论）
	totalDeleteCount, err := s.countCommentsRecursively(tx, comment.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	totalDeleteCount += 1 // 加上主评论本身
	
	// 硬删除子评论（递归删除所有回复）
	if err := s.deleteChildCommentsRecursively(tx, comment.ID); err != nil {
		tx.Rollback()
		return err
	}
	
	// 硬删除评论本身
	if err := tx.Delete(&comment).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新文章评论数（减去所有被删除的评论数）
	if err := tx.Model(&models.Article{}).Where("id = ?", comment.ArticleID).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", totalDeleteCount)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果是回复评论，更新父评论的回复数
	if comment.ParentID != nil && *comment.ParentID > 0 {
		if err := tx.Model(&models.Comment{}).Where("id = ?", *comment.ParentID).UpdateColumn("reply_count", gorm.Expr("reply_count - ?", 1)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}
	
	// 清除Redis缓存（如果有的话）
	s.clearCommentCache(comment.ArticleID, comment.ID)
	
	return nil
}

// GetComment 获取评论详情
func (s *CommentService) GetComment(commentID uint) (*models.Comment, error) {
	var comment models.Comment
	if err := s.db.Preload("User").Preload("Article").Where("id = ? AND status = ?", commentID, 1).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("评论不存在")
		}
		return nil, err
	}

	return &comment, nil
}

// GetCommentsByArticleWithSort 获取文章的评论列表（带排序）
func (s *CommentService) GetCommentsByArticleWithSort(articleID uint, page, size int, sort string) ([]*models.Comment, int64, error) {
	// 验证文章是否存在
	var article models.Article
	if err := s.db.Where("id = ? AND status = ?", articleID, 1).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errors.New("文章不存在或已下线")
		}
		return nil, 0, err
	}

	// 构建查询条件
	query := s.db.Model(&models.Comment{}).Where("article_id = ? AND status = ? AND parent_id IS NULL", articleID, 1)

	// 获取顶级评论总数（用于分页）
	var topLevelTotal int64
	if err := query.Count(&topLevelTotal).Error; err != nil {
		return nil, 0, err
	}

	// 获取文章的总评论数（包含回复）
	var total int64 = int64(article.CommentCount)

	// 根据排序类型设置排序规则
	var orderBy string
	switch sort {
	case "newest":
		orderBy = "created_at DESC"
	case "oldest":
		orderBy = "created_at ASC"
	case "most_liked":
		orderBy = "like_count DESC, created_at DESC"
	default:
		orderBy = "like_count DESC, created_at DESC"
	}

	// 获取评论列表（只获取顶级评论）
	var comments []*models.Comment
	offset := (page - 1) * size
	if err := query.Preload("User").Order(orderBy).Offset(offset).Limit(size).Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// 为每个顶级评论加载回复
	for _, comment := range comments {
		if err := s.loadReplies(comment); err != nil {
			return nil, 0, err
		}
	}

	return comments, total, nil
}

// GetCommentsByArticle 获取文章的评论列表（保持兼容性）
func (s *CommentService) GetCommentsByArticle(articleID uint, page, size int) ([]*models.Comment, int64, error) {
	return s.GetCommentsByArticleWithSort(articleID, page, size, "most_liked")
}

// GetCommentsByArticleOld 获取文章的评论列表（旧版本）
func (s *CommentService) GetCommentsByArticleOld(articleID uint, page, size int) ([]*models.Comment, int64, error) {
	// 验证文章是否存在
	var article models.Article
	if err := s.db.Where("id = ? AND status = ?", articleID, 1).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errors.New("文章不存在或已下线")
		}
		return nil, 0, err
	}

	// 构建查询条件
	query := s.db.Model(&models.Comment{}).Where("article_id = ? AND status = ? AND parent_id IS NULL", articleID, 1)

	// 获取顶级评论总数（用于分页）
	var topLevelTotal int64
	if err := query.Count(&topLevelTotal).Error; err != nil {
		return nil, 0, err
	}

	// 获取文章的总评论数（包含回复）
	var total int64 = int64(article.CommentCount)

	// 获取评论列表（只获取顶级评论）
	var comments []*models.Comment
	offset := (page - 1) * size
	if err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(size).Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// 为每个顶级评论加载回复
	for _, comment := range comments {
		if err := s.loadReplies(comment); err != nil {
			return nil, 0, err
		}
	}

	return comments, total, nil
}

// GetCommentReplies 获取评论的回复列表
func (s *CommentService) GetCommentReplies(parentID uint, page, size int) ([]*models.Comment, int64, error) {
	// 验证父评论是否存在
	var parentComment models.Comment
	if err := s.db.Where("id = ? AND status = ?", parentID, 1).First(&parentComment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errors.New("父评论不存在")
		}
		return nil, 0, err
	}

	// 构建查询条件
	query := s.db.Model(&models.Comment{}).Where("parent_id = ? AND status = ?", parentID, 1)

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取回复列表
	var replies []*models.Comment
	offset := (page - 1) * size
	if err := query.Preload("User").Order("created_at ASC").Offset(offset).Limit(size).Find(&replies).Error; err != nil {
		return nil, 0, err
	}

	return replies, total, nil
}

// GetUserComments 获取用户的评论列表
func (s *CommentService) GetUserComments(userID uint, page, size int) ([]*models.Comment, int64, error) {
	// 构建查询条件
	query := s.db.Model(&models.Comment{}).Where("user_id = ? AND status = ?", userID, 1)

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取评论列表
	var comments []*models.Comment
	offset := (page - 1) * size
	if err := query.Preload("User").Preload("Article").Order("created_at DESC").Offset(offset).Limit(size).Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// LikeComment 点赞评论
func (s *CommentService) LikeComment(commentID, userID uint) error {
	// 验证评论是否存在
	var comment models.Comment
	if err := s.db.Where("id = ? AND status = ?", commentID, 1).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		return err
	}

	// 检查是否已经点赞
	var favorite models.Favorite
	err := s.db.Where("user_id = ? AND article_id = ?", userID, comment.ArticleID).First(&favorite).Error
	if err == nil {
		return errors.New("已经点赞过了")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建点赞记录
	favorite = models.Favorite{
		UserID:    userID,
		ArticleID: comment.ArticleID,
		CreatedAt: time.Now(),
	}
	if err := tx.Create(&favorite).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新评论点赞数
	if err := tx.Model(&comment).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// UnlikeComment 取消点赞评论
func (s *CommentService) UnlikeComment(commentID, userID uint) error {
	// 验证评论是否存在
	var comment models.Comment
	if err := s.db.Where("id = ? AND status = ?", commentID, 1).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在")
		}
		return err
	}

	// 检查是否已经点赞
	var favorite models.Favorite
	if err := s.db.Where("user_id = ? AND article_id = ?", userID, comment.ArticleID).First(&favorite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("还没有点赞")
		}
		return err
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除点赞记录
	if err := tx.Delete(&favorite).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新评论点赞数
	if err := tx.Model(&comment).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// loadReplies 加载评论的回复（递归加载，最多5层）
func (s *CommentService) loadReplies(comment *models.Comment) error {
	return s.loadRepliesRecursive(comment, 0, 5)
}

// loadRepliesRecursive 递归加载评论回复
func (s *CommentService) loadRepliesRecursive(comment *models.Comment, currentDepth, maxDepth int) error {
	if currentDepth >= maxDepth {
		return nil
	}

	// 加载直接回复
	var replies []*models.Comment
	if err := s.db.Where("parent_id = ? AND status = ?", comment.ID, 1).Preload("User").Order("created_at ASC").Find(&replies).Error; err != nil {
		return err
	}

	// 为每个回复递归加载子回复
	for _, reply := range replies {
		if err := s.loadRepliesRecursive(reply, currentDepth+1, maxDepth); err != nil {
			return err
		}
	}

	// 转换指针切片为值切片
	repliesValue := make([]models.Comment, len(replies))
	for i, reply := range replies {
		repliesValue[i] = *reply
	}
	comment.Replies = repliesValue
	return nil
}

// validateCreateRequest 验证创建评论请求
func (s *CommentService) validateCreateRequest(req *models.CommentCreateRequest) error {
	if req.ArticleID == 0 {
		return errors.New("文章ID不能为空")
	}

	if strings.TrimSpace(req.Content) == "" {
		return errors.New("评论内容不能为空")
	}

	if len(req.Content) > 1000 {
		return errors.New("评论内容不能超过1000个字符")
	}

	return nil
}

// validateUpdateRequest 验证更新评论请求
func (s *CommentService) validateUpdateRequest(req *models.CommentUpdateRequest) error {
	if strings.TrimSpace(req.Content) == "" {
		return errors.New("评论内容不能为空")
	}

	if len(req.Content) > 1000 {
		return errors.New("评论内容不能超过1000个字符")
	}

	return nil
}

// countCommentsRecursively 递归统计子评论数量
func (s *CommentService) countCommentsRecursively(tx *gorm.DB, parentID uint) (int64, error) {
	var childComments []models.Comment
	if err := tx.Where("parent_id = ?", parentID).Find(&childComments).Error; err != nil {
		return 0, err
	}
	
	count := int64(len(childComments))
	
	for _, child := range childComments {
		// 递归统计子评论的子评论
		childCount, err := s.countCommentsRecursively(tx, child.ID)
		if err != nil {
			return 0, err
		}
		count += childCount
	}
	
	return count, nil
}

// deleteChildCommentsRecursively 递归删除子评论
func (s *CommentService) deleteChildCommentsRecursively(tx *gorm.DB, parentID uint) error {
	var childComments []models.Comment
	if err := tx.Where("parent_id = ?", parentID).Find(&childComments).Error; err != nil {
		return err
	}
	
	for _, child := range childComments {
		// 递归删除子评论的子评论
		if err := s.deleteChildCommentsRecursively(tx, child.ID); err != nil {
			return err
		}
		
		// 删除子评论
		if err := tx.Delete(&child).Error; err != nil {
			return err
		}
	}
	
	return nil
}

// clearCommentCache 清除评论相关的Redis缓存
func (s *CommentService) clearCommentCache(articleID, commentID uint) {
	ctx := context.Background()
	
	// 清除可能的评论缓存key
	cacheKeys := []string{
		fmt.Sprintf("comments:article:%d", articleID),
		fmt.Sprintf("comment:%d", commentID),
		fmt.Sprintf("comments:article:%d:*", articleID), // 使用通配符清除分页缓存
	}
	
	// 删除每个缓存key（忽略错误，因为key可能不存在）
	for _, key := range cacheKeys {
		if strings.Contains(key, "*") {
			// 对于通配符key，需要先获取匹配的keys再删除
			keys, err := s.redisClient.Keys(ctx, key).Result()
			if err == nil && len(keys) > 0 {
				s.redisClient.Del(ctx, keys...)
			}
		} else {
			s.redisClient.Del(ctx, key)
		}
	}
}

// GetCommentCount 获取评论总数
func (s *CommentService) GetCommentCount() (int64, error) {
	var count int64
	err := s.db.Model(&models.Comment{}).Where("deleted_at IS NULL").Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("获取评论总数失败: %v", err)
	}
	return count, nil
}
