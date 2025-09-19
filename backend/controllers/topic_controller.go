package controllers

import (
	"net/http"
	"strconv"

	"godad-backend/models"
	"godad-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TopicController 话题控制器
type TopicController struct {
	topicService *services.TopicService
}

// NewTopicController 创建话题控制器实例
func NewTopicController(db *gorm.DB) *TopicController {
	return &TopicController{
		topicService: services.NewTopicService(db),
	}
}

// CreateTopic 创建话题
func (tc *TopicController) CreateTopic(c *gin.Context) {
	var req models.TopicCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	topic, err := tc.topicService.CreateTopic(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "话题创建成功",
		"data":    topic,
	})
}

// UpdateTopic 更新话题
func (tc *TopicController) UpdateTopic(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的话题ID",
		})
		return
	}

	var req models.TopicUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	topic, err := tc.topicService.UpdateTopic(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "话题更新成功",
		"data":    topic,
	})
}

// DeleteTopic 删除话题
func (tc *TopicController) DeleteTopic(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的话题ID",
		})
		return
	}

	if err := tc.topicService.DeleteTopic(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "话题删除成功",
	})
}

// GetTopic 获取话题详情
func (tc *TopicController) GetTopic(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的话题ID",
		})
		return
	}

	topic, err := tc.topicService.GetTopic(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    topic,
	})
}

// GetTopicList 获取话题列表
func (tc *TopicController) GetTopicList(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	// 解析是否显示全部（包括禁用的话题）
	showAll := c.Query("all") == "true"

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	topics, total, err := tc.topicService.GetTopicList(page, pageSize, showAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"items":      topics,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetActiveTopics 获取所有启用的话题
func (tc *TopicController) GetActiveTopics(c *gin.Context) {
	topics, err := tc.topicService.GetActiveTopics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    topics,
	})
}