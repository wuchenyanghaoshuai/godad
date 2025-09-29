package container

import (
	"sync"

	"godad-backend/config"
	"godad-backend/services"
	"gorm.io/gorm"
)

// Container 依赖注入容器
type Container struct {
	db           *gorm.DB
	services     map[string]interface{}
	servicesMux  sync.RWMutex
	singletons   map[string]interface{}
	singletonMux sync.RWMutex
}

// NewContainer 创建新的容器实例
func NewContainer() *Container {
	return &Container{
		services:   make(map[string]interface{}),
		singletons: make(map[string]interface{}),
	}
}

// SetDB 设置数据库连接
func (c *Container) SetDB(db *gorm.DB) {
	c.db = db
}

// GetDB 获取数据库连接
func (c *Container) GetDB() *gorm.DB {
	return c.db
}

// RegisterSingleton 注册单例服务
func (c *Container) RegisterSingleton(name string, factory func(*Container) interface{}) {
	c.singletonMux.Lock()
	defer c.singletonMux.Unlock()

	if _, exists := c.singletons[name]; !exists {
		c.singletons[name] = factory(c)
	}
}

// RegisterService 注册服务
func (c *Container) RegisterService(name string, factory func(*Container) interface{}) {
	c.servicesMux.Lock()
	defer c.servicesMux.Unlock()

	c.services[name] = factory
}

// GetSingleton 获取单例服务
func (c *Container) GetSingleton(name string) interface{} {
	c.singletonMux.RLock()
	defer c.singletonMux.RUnlock()

	return c.singletons[name]
}

// GetService 获取服务实例
func (c *Container) GetService(name string) interface{} {
	c.servicesMux.RLock()
	factory, exists := c.services[name]
	c.servicesMux.RUnlock()

	if !exists {
		return nil
	}

	if factoryFunc, ok := factory.(func(*Container) interface{}); ok {
		return factoryFunc(c)
	}

	return factory
}

// GetUserService 获取用户服务
func (c *Container) GetUserService() *services.UserService {
	service := c.GetSingleton("user_service")
	if service == nil {
		return nil
	}
	return service.(*services.UserService)
}

// GetArticleService 获取文章服务
func (c *Container) GetArticleService() *services.ArticleService {
	service := c.GetSingleton("article_service")
	if service == nil {
		return nil
	}
	return service.(*services.ArticleService)
}

// GetCommentService 获取评论服务
func (c *Container) GetCommentService() *services.CommentService {
	service := c.GetSingleton("comment_service")
	if service == nil {
		return nil
	}
	return service.(*services.CommentService)
}

// GetForumService 获取论坛服务
func (c *Container) GetForumService() *services.ForumService {
	service := c.GetSingleton("forum_service")
	if service == nil {
		return nil
	}
	return service.(*services.ForumService)
}

// GetCacheService 获取缓存服务
func (c *Container) GetCacheService() *services.CacheService {
	service := c.GetSingleton("cache_service")
	if service == nil {
		return nil
	}
	return service.(*services.CacheService)
}

// GetNotificationService 获取通知服务
func (c *Container) GetNotificationService() *services.NotificationService {
	service := c.GetSingleton("notification_service")
	if service == nil {
		return nil
	}
	return service.(*services.NotificationService)
}

// GetEmailService 获取邮件服务
func (c *Container) GetEmailService() *services.EmailService {
	service := c.GetSingleton("email_service")
	if service == nil {
		return nil
	}
	return service.(*services.EmailService)
}

// GetUploadService 获取上传服务
func (c *Container) GetUploadService() *services.UploadService {
	service := c.GetSingleton("upload_service")
	if service == nil {
		return nil
	}
	return service.(*services.UploadService)
}

// GetLikeService 获取点赞服务
func (c *Container) GetLikeService() *services.LikeService {
	service := c.GetSingleton("like_service")
	if service == nil {
		return nil
	}
	return service.(*services.LikeService)
}

// GetFavoriteService 获取收藏服务
func (c *Container) GetFavoriteService() *services.FavoriteService {
	service := c.GetSingleton("favorite_service")
	if service == nil {
		return nil
	}
	return service.(*services.FavoriteService)
}

// GetFollowService 获取关注服务
func (c *Container) GetFollowService() *services.FollowService {
	service := c.GetSingleton("follow_service")
	if service == nil {
		return nil
	}
	return service.(*services.FollowService)
}

// InitializeServices 初始化所有服务
func (c *Container) InitializeServices() {
	// 注册基础服务
	c.RegisterSingleton("cache_service", func(container *Container) interface{} {
		return services.NewCacheService()
	})

	c.RegisterSingleton("email_service", func(container *Container) interface{} {
		return services.NewEmailService()
	})

	// 注册业务服务
	c.RegisterSingleton("user_service", func(container *Container) interface{} {
		return services.NewUserServiceWithDI(
			container.GetDB(),
			container.GetCacheService(),
		)
	})

	c.RegisterSingleton("article_service", func(container *Container) interface{} {
		return services.NewArticleServiceWithDI(
			container.GetDB(),
			container.GetCacheService(),
		)
	})
}

// GlobalContainer 全局容器实例
var GlobalContainer *Container
var once sync.Once

// GetGlobalContainer 获取全局容器实例
func GetGlobalContainer() *Container {
	once.Do(func() {
		GlobalContainer = NewContainer()
		GlobalContainer.SetDB(config.GetDB())
		GlobalContainer.InitializeServices()
	})
	return GlobalContainer
}