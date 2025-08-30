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

// CommentService 评论服务
type CommentService struct {
	db *gorm.DB
}

// NewCommentService 创建评论服务实例
func NewCommentService() *CommentService {
	return &CommentService{
		db: config.GetDB(),
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

	// 预加载关联数据
	if err := s.db.Preload("User").Preload("Article").First(comment, comment.ID).Error; err != nil {
		return nil, err
	}

	return comment, nil
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

	// 软删除评论（设置状态为0）
	if err := tx.Model(&comment).Updates(map[string]interface{}{
		"status":     0,
		"deleted_at": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新文章评论数
	if err := tx.Model(&models.Article{}).Where("id = ?", comment.ArticleID).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
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
	return tx.Commit().Error
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

// GetCommentsByArticle 获取文章的评论列表
func (s *CommentService) GetCommentsByArticle(articleID uint, page, size int) ([]*models.Comment, int64, error) {
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

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

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

// loadReplies 加载评论的回复（递归加载，最多3层）
func (s *CommentService) loadReplies(comment *models.Comment) error {
	// 加载直接回复（最多显示5条）
	var replies []*models.Comment
	if err := s.db.Where("parent_id = ? AND status = ?", comment.ID, 1).Preload("User").Order("created_at ASC").Limit(5).Find(&replies).Error; err != nil {
		return err
	}

	// 为每个回复递归加载子回复（最多2层）
	for _, reply := range replies {
		var subReplies []*models.Comment
		if err := s.db.Where("parent_id = ? AND status = ?", reply.ID, 1).Preload("User").Order("created_at ASC").Limit(3).Find(&subReplies).Error; err != nil {
			return err
		}
		// 转换指针切片为值切片
		subRepliesValue := make([]models.Comment, len(subReplies))
		for i, subReply := range subReplies {
			subRepliesValue[i] = *subReply
		}
		reply.Replies = subRepliesValue
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

// GetCommentCount 获取评论总数
func (s *CommentService) GetCommentCount() (int64, error) {
	var count int64
	err := s.db.Model(&models.Comment{}).Where("deleted_at IS NULL").Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("获取评论总数失败: %v", err)
	}
	return count, nil
}