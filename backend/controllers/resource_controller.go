package controllers

import (
	"net/http"
	"strconv"

	"godad-backend/services"
	"godad-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ResourceController 资源控制器
type ResourceController struct {
	resourceService *services.ResourceService
}

// NewResourceController 创建资源控制器实例
func NewResourceController(db *gorm.DB) *ResourceController {
	return &ResourceController{
		resourceService: services.NewResourceService(db),
	}
}

// CreateResource 创建资源
func (ctrl *ResourceController) CreateResource(c *gin.Context) {
	var req services.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 如果是管理员，可以设置上传者ID；如果是普通用户，设置为当前用户
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	userRole, _ := c.Get("user_role")
	if userRole != 2 { // 如果不是管理员
		uid := userID.(uint)
		req.UploaderID = &uid
		req.Status = 0 // 普通用户上传的资源默认为待审核状态
	}

	resource, err := ctrl.resourceService.CreateResource(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建资源失败")
		return
	}

	utils.SuccessResponse(c, resource, "创建资源成功")
}

// GetResource 获取单个资源
func (ctrl *ResourceController) GetResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的资源ID")
		return
	}

	resource, err := ctrl.resourceService.GetResourceByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "资源不存在")
		return
	}

	utils.SuccessResponse(c, resource, "获取资源成功")
}

// GetResources 获取资源列表
func (ctrl *ResourceController) GetResources(c *gin.Context) {
	var req services.GetResourcesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	response, err := ctrl.resourceService.GetResources(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取资源列表失败")
		return
	}

	utils.SuccessResponse(c, response, "获取资源列表成功")
}

// GetPublishedResources 获取已发布的资源列表（前端公开接口）
func (ctrl *ResourceController) GetPublishedResources(c *gin.Context) {
	var req services.GetResourcesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	response, err := ctrl.resourceService.GetPublishedResources(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取资源列表失败")
		return
	}

	utils.SuccessResponse(c, response, "获取资源列表成功")
}

// UpdateResource 更新资源
func (ctrl *ResourceController) UpdateResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的资源ID")
		return
	}

	var req services.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	resource, err := ctrl.resourceService.UpdateResource(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新资源失败")
		return
	}

	utils.SuccessResponse(c, resource, "更新资源成功")
}

// UpdateResourceStatus 更新资源状态
func (ctrl *ResourceController) UpdateResourceStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的资源ID")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 验证状态值
	if req.Status < 0 || req.Status > 2 {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的状态值")
		return
	}

	err = ctrl.resourceService.UpdateResourceStatus(uint(id), req.Status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新资源状态失败")
		return
	}

	utils.SuccessResponse(c, nil, "更新资源状态成功")
}

// DeleteResource 删除资源
func (ctrl *ResourceController) DeleteResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的资源ID")
		return
	}

	err = ctrl.resourceService.DeleteResource(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除资源失败")
		return
	}

	utils.SuccessResponse(c, nil, "删除资源成功")
}

// DownloadResource 下载资源（增加下载次数）
func (ctrl *ResourceController) DownloadResource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的资源ID")
		return
	}

	// 获取资源信息
	resource, err := ctrl.resourceService.GetResourceByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "资源不存在")
		return
	}

	// 检查资源状态
	if resource.Status != 1 {
		utils.ErrorResponse(c, http.StatusForbidden, "资源不可用")
		return
	}

	// 增加下载次数
	if err := ctrl.resourceService.IncrementDownloadCount(uint(id)); err != nil {
		// 即使增加下载次数失败，也不影响下载功能
		// 只记录错误，不返回错误响应
	}

	// 返回资源文件URL
	utils.SuccessResponse(c, gin.H{
		"file_url": resource.FileURL,
		"title":    resource.Title,
	}, "获取下载链接成功")
}

// GetResourceStats 获取资源统计信息
func (ctrl *ResourceController) GetResourceStats(c *gin.Context) {
	stats, err := ctrl.resourceService.GetResourceStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取统计信息失败")
		return
	}

	utils.SuccessResponse(c, stats, "获取统计信息成功")
}