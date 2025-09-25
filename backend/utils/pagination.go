package utils

import (
	"math"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PaginationResponse 分页响应结构
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Size       int         `json:"size"`
	TotalPages int         `json:"total_pages"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
}

// Paginate 分页处理函数
func Paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if size <= 0 || size > 100 {
			size = 10
		}

		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

// NewPaginationResponse 创建分页响应
func NewPaginationResponse(data interface{}, total int64, page, size int) *PaginationResponse {
	totalPages := int(math.Ceil(float64(total) / float64(size)))

	return &PaginationResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// RespondWithPagination 返回分页响应
func RespondWithPagination(c *gin.Context, data interface{}, total int64, page, size int) {
	response := NewPaginationResponse(data, total, page, size)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    response,
	})
}