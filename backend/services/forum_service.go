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

// ForumService 论坛服务
type ForumService struct {
    db                  *gorm.DB
    notificationService *NotificationService
}

// NewForumService 创建论坛服务实例
func NewForumService() *ForumService {
    db := config.GetDB()
    return &ForumService{
        db:                  db,
        notificationService: NewNotificationService(db),
    }
}

// GetForumStats 获取论坛统计数据
func (s *ForumService) GetForumStats() (*models.ForumStats, error) {
    stats := &models.ForumStats{}

    // 总帖子数（已发布）
    if err := s.db.Model(&models.ForumPost{}).
        Where("status = ?", 1).
        Count(&stats.ForumPostCount).Error; err != nil {
        return nil, fmt.Errorf("统计帖子数量失败: %w", err)
    }

    // 总回复数（已发布）
    if err := s.db.Model(&models.ForumReply{}).
        Where("status = ?", 1).
        Count(&stats.ForumReplyCount).Error; err != nil {
        return nil, fmt.Errorf("统计回复数量失败: %w", err)
    }

    // 精华帖子（is_hot=true）
    if err := s.db.Model(&models.ForumPost{}).
        Where("status = ? AND is_hot = ?", 1, true).
        Count(&stats.FeaturedPostCount).Error; err != nil {
        return nil, fmt.Errorf("统计精华帖子失败: %w", err)
    }

    // 置顶帖子（is_top=true）
    if err := s.db.Model(&models.ForumPost{}).
        Where("status = ? AND is_top = ?", 1, true).
        Count(&stats.PinnedPostCount).Error; err != nil {
        return nil, fmt.Errorf("统计置顶帖子失败: %w", err)
    }

    // 今日新增（按本地时区）
    now := time.Now()
    startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

    if err := s.db.Model(&models.ForumPost{}).
        Where("status = ? AND created_at >= ?", 1, startOfDay).
        Count(&stats.TodayNewPosts).Error; err != nil {
        return nil, fmt.Errorf("统计今日新增帖子失败: %w", err)
    }

    if err := s.db.Model(&models.ForumReply{}).
        Where("status = ? AND created_at >= ?", 1, startOfDay).
        Count(&stats.TodayNewReplies).Error; err != nil {
        return nil, fmt.Errorf("统计今日新增回复失败: %w", err)
    }

    // 最近7天发帖的活跃用户数（去重author_id）
    sevenDaysAgo := now.AddDate(0, 0, -7)
    if err := s.db.Model(&models.ForumPost{}).
        Where("status = ? AND created_at >= ?", 1, sevenDaysAgo).
        Distinct("author_id").
        Count(&stats.ActiveUserCount).Error; err != nil {
        return nil, fmt.Errorf("统计活跃用户失败: %w", err)
    }

    return stats, nil
}

// CreatePost 创建帖子
func (s *ForumService) CreatePost(req *models.ForumPostCreateRequest, userID uint) (*models.ForumPost, error) {
	// 验证请求参数
	if err := s.validateCreatePostRequest(req); err != nil {
		return nil, err
	}

	// 验证用户是否存在
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("验证用户失败: %w", err)
	}

	// 创建帖子
	post := &models.ForumPost{
		Title:       strings.TrimSpace(req.Title),
		Content:     strings.TrimSpace(req.Content),
		Topic:       req.Topic,
		AuthorID:    userID,
		Status:      1, // 直接发布
		LastReplyAt: nil,
	}

	// 保存到数据库
	if err := s.db.Create(post).Error; err != nil {
		return nil, fmt.Errorf("创建帖子失败: %w", err)
	}

	// 预加载关联数据
	if err := s.db.Preload("Author").First(post, post.ID).Error; err != nil {
		return nil, fmt.Errorf("加载帖子数据失败: %w", err)
	}

	return post, nil
}

// GetPostList 获取帖子列表
func (s *ForumService) GetPostList(req *models.ForumPostListRequest) ([]models.ForumPost, int64, error) {
	var posts []models.ForumPost
	var total int64

	// 构建查询条件
	query := s.db.Model(&models.ForumPost{}).Where("status = ?", 1)

	// 话题筛选
	if req.Topic != "" && req.Topic != "All" {
		query = query.Where("topic = ?", req.Topic)
	}

	// 作者筛选
	if req.AuthorID > 0 {
		query = query.Where("author_id = ?", req.AuthorID)
	}

	// 关键词搜索
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("title LIKE ? OR content LIKE ?", keyword, keyword)
	}

	// 置顶筛选
	if req.IsTop != nil {
		query = query.Where("is_top = ?", *req.IsTop)
	}

	// 热门筛选
	if req.IsHot != nil {
		query = query.Where("is_hot = ?", *req.IsHot)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计帖子数量失败: %w", err)
	}

	// 排序
	orderStr := s.buildOrderString(req.Sort)
	if orderStr != "" {
		query = query.Order(orderStr)
	}

	// 分页
	offset := (req.Page - 1) * req.Size
	query = query.Offset(offset).Limit(req.Size)

	// 预加载关联数据
	query = query.Preload("Author")

	// 执行查询
	if err := query.Find(&posts).Error; err != nil {
		return nil, 0, fmt.Errorf("查询帖子列表失败: %w", err)
	}

	return posts, total, nil
}

// GetAdminPostList 管理员获取帖子列表（不强制仅已发布，可按状态筛选，支持包含软删除）
func (s *ForumService) GetAdminPostList(req *models.AdminForumPostListRequest) ([]models.ForumPost, int64, error) {
    var posts []models.ForumPost
    var total int64

    // 构建查询（可选包含软删除）
    query := s.db.Model(&models.ForumPost{})
    if req.IncludeDeleted {
        query = query.Unscoped()
    }

    // 状态筛选（可选）
    if req.Status != nil {
        query = query.Where("status = ?", *req.Status)
    }

    // 话题筛选
    if req.Topic != "" && req.Topic != models.TopicAll {
        query = query.Where("topic = ?", req.Topic)
    }

    // 作者筛选
    if req.AuthorID > 0 {
        query = query.Where("author_id = ?", req.AuthorID)
    }

    // 关键词搜索
    if req.Keyword != "" {
        keyword := "%" + req.Keyword + "%"
        query = query.Where("title LIKE ? OR content LIKE ?", keyword, keyword)
    }

    // 置顶/热门筛选
    if req.IsTop != nil {
        query = query.Where("is_top = ?", *req.IsTop)
    }
    if req.IsHot != nil {
        query = query.Where("is_hot = ?", *req.IsHot)
    }

    // 统计总数
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, fmt.Errorf("统计帖子数量失败: %w", err)
    }

    // 排序：优先置顶/精华，再按指定排序
    orderStr := s.buildOrderString(req.Sort)
    orderClause := "is_top DESC, is_hot DESC"
    if orderStr != "" {
        orderClause = orderClause + ", " + orderStr
    }
    query = query.Order(orderClause)

    // 分页
    offset := (req.Page - 1) * req.Size
    query = query.Offset(offset).Limit(req.Size)

    // 预加载作者
    query = query.Preload("Author")

    // 查询
    if err := query.Find(&posts).Error; err != nil {
        return nil, 0, fmt.Errorf("查询帖子列表失败: %w", err)
    }

    return posts, total, nil
}

// GetPostByID 根据ID获取帖子详情
func (s *ForumService) GetPostByID(id uint) (*models.ForumPost, error) {
	var post models.ForumPost

	// 查询帖子并预加载关联数据
	if err := s.db.Preload("Author").Where("id = ? AND status = ?", id, 1).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("帖子不存在")
		}
		return nil, fmt.Errorf("查询帖子失败: %w", err)
	}

	return &post, nil
}

// UpdatePost 更新帖子
func (s *ForumService) UpdatePost(id uint, req *models.ForumPostUpdateRequest, userID uint, isAdmin bool) (*models.ForumPost, error) {
	// 查找帖子
	var post models.ForumPost
	if err := s.db.Where("id = ? AND status = ?", id, 1).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("帖子不存在")
		}
		return nil, fmt.Errorf("查询帖子失败: %w", err)
	}

    // 锁定检查：锁定后只有管理员可修改
    if post.IsLocked && !isAdmin {
        return nil, errors.New("帖子已锁定，无法编辑")
    }

    // 权限检查：作者或管理员
    if post.AuthorID != userID && !isAdmin {
        return nil, errors.New("无权限修改此帖子")
    }

	// 更新字段
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = strings.TrimSpace(req.Title)
	}
	if req.Content != "" {
		updates["content"] = strings.TrimSpace(req.Content)
	}
	if req.Topic != "" {
		updates["topic"] = req.Topic
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	// 如果没有更新内容，直接返回
	if len(updates) == 0 {
		return &post, nil
	}

	// 执行更新
	if err := s.db.Model(&post).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新帖子失败: %w", err)
	}

	// 重新加载帖子数据
	if err := s.db.Preload("Author").First(&post, id).Error; err != nil {
		return nil, fmt.Errorf("加载更新后的帖子数据失败: %w", err)
	}

	return &post, nil
}

// AdminSetPostTop 设置帖子置顶状态
func (s *ForumService) AdminSetPostTop(id uint, top bool) error {
    return s.db.Model(&models.ForumPost{}).Where("id = ?", id).Update("is_top", top).Error
}

// AdminSetPostHot 设置帖子热门状态
func (s *ForumService) AdminSetPostHot(id uint, hot bool) error {
    return s.db.Model(&models.ForumPost{}).Where("id = ?", id).Update("is_hot", hot).Error
}

// AdminSetPostLock 设置帖子锁定状态
func (s *ForumService) AdminSetPostLock(id uint, locked bool) error {
    return s.db.Model(&models.ForumPost{}).Where("id = ?", id).Update("is_locked", locked).Error
}

// AdminDeletePost 管理员删除帖子（软删除）
func (s *ForumService) AdminDeletePost(id uint) error {
    var post models.ForumPost
    if err := s.db.Where("id = ?", id).First(&post).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("帖子不存在")
        }
        return fmt.Errorf("查询帖子失败: %w", err)
    }
    if err := s.db.Delete(&post).Error; err != nil {
        return fmt.Errorf("删除帖子失败: %w", err)
    }
    return nil
}

// DeletePost 删除帖子
func (s *ForumService) DeletePost(id uint, userID uint) error {
	// 查找帖子
	var post models.ForumPost
	if err := s.db.Where("id = ? AND status = ?", id, 1).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("帖子不存在")
		}
		return fmt.Errorf("查询帖子失败: %w", err)
	}

	// 检查权限：只有作者本人可以删除
	if post.AuthorID != userID {
		return errors.New("无权限删除此帖子")
	}

	// 软删除帖子
	if err := s.db.Delete(&post).Error; err != nil {
		return fmt.Errorf("删除帖子失败: %w", err)
	}

	return nil
}

// IncrementPostView 增加帖子浏览量
func (s *ForumService) IncrementPostView(id uint) error {
	// 增加浏览量
	if err := s.db.Model(&models.ForumPost{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error; err != nil {
		return fmt.Errorf("增加浏览量失败: %w", err)
	}

	return nil
}

// CreateReply 创建回复
func (s *ForumService) CreateReply(req *models.ForumReplyCreateRequest, userID uint) (*models.ForumReply, error) {
	// 验证请求参数
	if err := s.validateCreateReplyRequest(req); err != nil {
		return nil, err
	}

	// 验证帖子是否存在
	var post models.ForumPost
	if err := s.db.Where("id = ? AND status = ?", req.PostID, 1).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("帖子不存在")
		}
		return nil, fmt.Errorf("验证帖子失败: %w", err)
	}

	// 如果是回复某个评论，验证父评论是否存在
	if req.ParentID != nil {
		var parentReply models.ForumReply
		if err := s.db.Where("id = ? AND post_id = ? AND status = ?", *req.ParentID, req.PostID, 1).First(&parentReply).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("父回复不存在")
			}
			return nil, fmt.Errorf("验证父回复失败: %w", err)
		}
	}

	// 创建回复
	reply := &models.ForumReply{
		PostID:   req.PostID,
		AuthorID: userID,
		ParentID: req.ParentID,
		Content:  strings.TrimSpace(req.Content),
		Status:   1, // 直接发布
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 保存回复
	if err := tx.Create(reply).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建回复失败: %w", err)
	}

	// 更新帖子的回复数量和最后回复时间
	now := time.Now()
	if err := tx.Model(&post).Updates(map[string]interface{}{
		"reply_count":   gorm.Expr("reply_count + 1"),
		"last_reply_at": &now,
	}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新帖子统计失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	// 预加载关联数据
	if err := s.db.Preload("Author").First(reply, reply.ID).Error; err != nil {
		return nil, fmt.Errorf("加载回复数据失败: %w", err)
	}

	// 发送通知给帖子作者（如果不是自己回复自己的帖子）
	if post.AuthorID != userID {
		go s.sendReplyNotification(&post, reply, userID)
	}

	return reply, nil
}

// GetReplyList 获取回复列表
func (s *ForumService) GetReplyList(req *models.ForumReplyListRequest) ([]models.ForumReply, int64, error) {
	var replies []models.ForumReply
	var total int64

	// 构建查询条件
	query := s.db.Model(&models.ForumReply{}).Where("post_id = ? AND status = ?", req.PostID, 1)

	// 如果指定了父回复ID，则只查询该回复的子回复
	if req.ParentID != nil {
		query = query.Where("parent_id = ?", *req.ParentID)
	} else {
		// 否则只查询顶级回复（没有父回复的）
		query = query.Where("parent_id IS NULL")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计回复数量失败: %w", err)
	}

    // 排序：优先置顶/精华，再按指定排序
    orderStr := s.buildOrderString(req.Sort)
    orderClause := "is_top DESC, is_hot DESC"
    if orderStr != "" {
        orderClause = orderClause + ", " + orderStr
    }
    query = query.Order(orderClause)

	// 分页
	offset := (req.Page - 1) * req.Size
	query = query.Offset(offset).Limit(req.Size)

	// 预加载关联数据
	query = query.Preload("Author")

	// 执行查询
	if err := query.Find(&replies).Error; err != nil {
		return nil, 0, fmt.Errorf("查询回复列表失败: %w", err)
	}

	return replies, total, nil
}

// DeleteReply 删除回复
func (s *ForumService) DeleteReply(id uint, userID uint) error {
	// 查找回复
	var reply models.ForumReply
	if err := s.db.Where("id = ? AND status = ?", id, 1).First(&reply).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("回复不存在")
		}
		return fmt.Errorf("查询回复失败: %w", err)
	}

	// 检查权限：只有回复作者本人和管理员可以删除
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("验证用户失败: %w", err)
	}

	// 检查权限：作者本人或管理员（角色为2）
	if reply.AuthorID != userID && user.Role != 2 {
		return errors.New("无权限删除此回复")
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 软删除回复
	if err := tx.Delete(&reply).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除回复失败: %w", err)
	}

	// 更新帖子的回复数量
	if err := tx.Model(&models.ForumPost{}).Where("id = ?", reply.PostID).UpdateColumn("reply_count", gorm.Expr("reply_count - 1")).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新帖子统计失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// 验证创建帖子请求
func (s *ForumService) validateCreatePostRequest(req *models.ForumPostCreateRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return errors.New("帖子标题不能为空")
	}
	if len(req.Title) > 200 {
		return errors.New("帖子标题不能超过200个字符")
	}
	if strings.TrimSpace(req.Content) == "" {
		return errors.New("帖子内容不能为空")
	}
	if len(req.Content) > 10000 {
		return errors.New("帖子内容不能超过10000个字符")
	}
	if !models.IsValidTopic(req.Topic) {
		return errors.New("无效的话题分类")
	}
	return nil
}

// 验证创建回复请求
func (s *ForumService) validateCreateReplyRequest(req *models.ForumReplyCreateRequest) error {
	if req.PostID == 0 {
		return errors.New("帖子ID不能为空")
	}
	if strings.TrimSpace(req.Content) == "" {
		return errors.New("回复内容不能为空")
	}
	if len(req.Content) > 5000 {
		return errors.New("回复内容不能超过5000个字符")
	}
	return nil
}

// 构建排序字符串
func (s *ForumService) buildOrderString(sort string) string {
	if sort == "" {
		return "created_at DESC"
	}

	// 验证排序参数
	validSorts := map[string]string{
		"created_at desc":    "created_at DESC",
		"created_at asc":     "created_at ASC",
		"reply_count desc":   "reply_count DESC",
		"view_count desc":    "view_count DESC",
		"last_reply_at desc": "last_reply_at DESC",
		"like_count desc":    "like_count DESC",
	}

	if orderStr, exists := validSorts[strings.ToLower(sort)]; exists {
		return orderStr
	}

	return "created_at DESC"
}

// 发送回复通知
func (s *ForumService) sendReplyNotification(post *models.ForumPost, reply *models.ForumReply, replyUserID uint) {
	// 获取回复用户信息
	var replyUser models.User
	if err := s.db.Where("id = ?", replyUserID).First(&replyUser).Error; err != nil {
		return
	}

	// 创建通知
	notification := &models.Notification{
		ReceiverID: post.AuthorID,
		ActorID:    replyUserID,
		Type:       "comment", // 使用现有的comment类型
		Message:    fmt.Sprintf("%s 回复了您的帖子《%s》", replyUser.Username, post.Title),
		IsRead:     false,
		ResourceID: post.ID,
	}

	s.notificationService.CreateNotification(notification)
}
