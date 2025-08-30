package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"godad-backend/models"

	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

type CacheService struct{}

func NewCacheService() *CacheService {
	return &CacheService{}
}

func InitRedis() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis连接失败: %v", err)
	}

	return nil
}

func (c *CacheService) SetArticleList(key string, articles []models.Article, duration time.Duration) error {
	articlesJSON, err := json.Marshal(articles)
	if err != nil {
		return fmt.Errorf("文章列表序列化失败: %v", err)
	}

	err = rdb.Set(ctx, key, articlesJSON, duration).Err()
	if err != nil {
		return fmt.Errorf("缓存文章列表失败: %v", err)
	}

	return nil
}

func (c *CacheService) GetArticleList(key string) ([]models.Article, error) {
	articlesJSON, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("缓存未找到")
		}
		return nil, fmt.Errorf("获取缓存文章列表失败: %v", err)
	}

	var articles []models.Article
	err = json.Unmarshal([]byte(articlesJSON), &articles)
	if err != nil {
		return nil, fmt.Errorf("反序列化文章列表失败: %v", err)
	}

	return articles, nil
}

func (c *CacheService) SetArticle(articleID uint, article *models.Article, duration time.Duration) error {
	key := fmt.Sprintf("article:%d", articleID)
	articleJSON, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("文章序列化失败: %v", err)
	}

	err = rdb.Set(ctx, key, articleJSON, duration).Err()
	if err != nil {
		return fmt.Errorf("缓存文章失败: %v", err)
	}

	return nil
}

func (c *CacheService) GetArticle(articleID uint) (*models.Article, error) {
	key := fmt.Sprintf("article:%d", articleID)
	articleJSON, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("缓存未找到")
		}
		return nil, fmt.Errorf("获取缓存文章失败: %v", err)
	}

	var article models.Article
	err = json.Unmarshal([]byte(articleJSON), &article)
	if err != nil {
		return nil, fmt.Errorf("反序列化文章失败: %v", err)
	}

	return &article, nil
}

func (c *CacheService) SetUser(userID uint, user *models.User, duration time.Duration) error {
	key := fmt.Sprintf("user:%d", userID)
	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("用户信息序列化失败: %v", err)
	}

	err = rdb.Set(ctx, key, userJSON, duration).Err()
	if err != nil {
		return fmt.Errorf("缓存用户信息失败: %v", err)
	}

	return nil
}

func (c *CacheService) GetUser(userID uint) (*models.User, error) {
	key := fmt.Sprintf("user:%d", userID)
	userJSON, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("缓存未找到")
		}
		return nil, fmt.Errorf("获取缓存用户信息失败: %v", err)
	}

	var user models.User
	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return nil, fmt.Errorf("反序列化用户信息失败: %v", err)
	}

	return &user, nil
}

func (c *CacheService) SetCategories(categories []models.Category, duration time.Duration) error {
	key := "categories"
	categoriesJSON, err := json.Marshal(categories)
	if err != nil {
		return fmt.Errorf("分类列表序列化失败: %v", err)
	}

	err = rdb.Set(ctx, key, categoriesJSON, duration).Err()
	if err != nil {
		return fmt.Errorf("缓存分类列表失败: %v", err)
	}

	return nil
}

func (c *CacheService) GetCategories() ([]models.Category, error) {
	key := "categories"
	categoriesJSON, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("缓存未找到")
		}
		return nil, fmt.Errorf("获取缓存分类列表失败: %v", err)
	}

	var categories []models.Category
	err = json.Unmarshal([]byte(categoriesJSON), &categories)
	if err != nil {
		return nil, fmt.Errorf("反序列化分类列表失败: %v", err)
	}

	return categories, nil
}

func (c *CacheService) SetSearchResults(keyword string, articles []models.Article, duration time.Duration) error {
	key := fmt.Sprintf("search:%s", keyword)
	articlesJSON, err := json.Marshal(articles)
	if err != nil {
		return fmt.Errorf("搜索结果序列化失败: %v", err)
	}

	err = rdb.Set(ctx, key, articlesJSON, duration).Err()
	if err != nil {
		return fmt.Errorf("缓存搜索结果失败: %v", err)
	}

	return nil
}

func (c *CacheService) GetSearchResults(keyword string) ([]models.Article, error) {
	key := fmt.Sprintf("search:%s", keyword)
	articlesJSON, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("缓存未找到")
		}
		return nil, fmt.Errorf("获取缓存搜索结果失败: %v", err)
	}

	var articles []models.Article
	err = json.Unmarshal([]byte(articlesJSON), &articles)
	if err != nil {
		return nil, fmt.Errorf("反序列化搜索结果失败: %v", err)
	}

	return articles, nil
}

func (c *CacheService) IncreaseViewCount(articleID uint) error {
	key := fmt.Sprintf("article:views:%d", articleID)
	err := rdb.Incr(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("增加浏览次数失败: %v", err)
	}

	rdb.Expire(ctx, key, 24*time.Hour)
	return nil
}

func (c *CacheService) GetViewCount(articleID uint) (int64, error) {
	key := fmt.Sprintf("article:views:%d", articleID)
	count, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, fmt.Errorf("获取浏览次数失败: %v", err)
	}

	viewCount, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("转换浏览次数失败: %v", err)
	}

	return viewCount, nil
}

func (c *CacheService) SetHotKeywords(keywords []string, duration time.Duration) error {
	key := "hot:keywords"
	keywordsJSON, err := json.Marshal(keywords)
	if err != nil {
		return fmt.Errorf("热门关键词序列化失败: %v", err)
	}

	err = rdb.Set(ctx, key, keywordsJSON, duration).Err()
	if err != nil {
		return fmt.Errorf("缓存热门关键词失败: %v", err)
	}

	return nil
}

func (c *CacheService) GetHotKeywords() ([]string, error) {
	key := "hot:keywords"
	keywordsJSON, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("缓存未找到")
		}
		return nil, fmt.Errorf("获取缓存热门关键词失败: %v", err)
	}

	var keywords []string
	err = json.Unmarshal([]byte(keywordsJSON), &keywords)
	if err != nil {
		return nil, fmt.Errorf("反序列化热门关键词失败: %v", err)
	}

	return keywords, nil
}

func (c *CacheService) RecordSearchKeyword(keyword string) error {
	key := "search:keyword:count"
	err := rdb.ZIncrBy(ctx, key, 1, keyword).Err()
	if err != nil {
		return err
	}

	// 设置过期时间（如果是新的有序集合）
	rdb.Expire(ctx, key, time.Hour*24*30) // 30天过期

	return nil
}

func (c *CacheService) GetTopSearchKeywords(limit int64) ([]string, error) {
	key := "search:keyword:count"
	keywords, err := rdb.ZRevRange(ctx, key, 0, limit-1).Result()
	if err != nil {
		return nil, fmt.Errorf("获取热门搜索关键词失败: %v", err)
	}

	return keywords, nil
}

func (c *CacheService) Delete(key string) error {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("删除缓存失败: %v", err)
	}
	return nil
}

func (c *CacheService) DeletePattern(pattern string) error {
	keys, err := rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("查找匹配键失败: %v", err)
	}

	if len(keys) > 0 {
		err = rdb.Del(ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("删除匹配键失败: %v", err)
		}
	}

	return nil
}

func (c *CacheService) FlushAll() error {
	err := rdb.FlushAll(ctx).Err()
	if err != nil {
		return fmt.Errorf("清空缓存失败: %v", err)
	}
	return nil
}

func (c *CacheService) SetWithExpire(key string, value interface{}, duration time.Duration) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化失败: %v", err)
	}

	err = rdb.Set(ctx, key, valueJSON, duration).Err()
	if err != nil {
		return fmt.Errorf("设置缓存失败: %v", err)
	}

	return nil
}

func (c *CacheService) Get(key string, dest interface{}) error {
	valueJSON, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("缓存未找到")
		}
		return fmt.Errorf("获取缓存失败: %v", err)
	}

	err = json.Unmarshal([]byte(valueJSON), dest)
	if err != nil {
		return fmt.Errorf("反序列化失败: %v", err)
	}

	return nil
}