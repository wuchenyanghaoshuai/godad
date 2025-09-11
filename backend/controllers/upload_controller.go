package controllers

import (
	"fmt"
	"strconv"

	"godad-backend/middleware"
	"godad-backend/models"
	"godad-backend/services"
	"godad-backend/utils"

	"github.com/gin-gonic/gin"
)

// UploadController 上传控制器
type UploadController struct {
	uploadService *services.UploadService
}

// NewUploadController 创建上传控制器实例
func NewUploadController() *UploadController {
	uploadService, err := services.NewUploadService()
	if err != nil {
		panic("初始化上传服务失败: " + err.Error())
	}

	return &UploadController{
		uploadService: uploadService,
	}
}

// UploadImage 上传图片
// @Summary 上传图片
// @Description 上传图片文件到阿里云OSS
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param file formData file true "图片文件"
// @Param type formData string false "上传类型" Enums(article,comment,other) default(other)
// @Success 200 {object} utils.Response{data=models.UploadResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/upload/image [post]
func (c *UploadController) UploadImage(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "请选择要上传的文件")
		return
	}

	// 获取上传类型
	uploadType := ctx.DefaultPostForm("type", "other")

	// 获取文章标题（可选）
	articleTitle := ctx.PostForm("article_title")
	
	// 添加调试日志
	fmt.Printf("=== 上传调试信息 ===\n")
	fmt.Printf("uploadType: %s\n", uploadType)
	fmt.Printf("articleTitle: '%s'\n", articleTitle)
	fmt.Printf("articleTitle length: %d\n", len(articleTitle))
	
	// 打印所有表单参数
	ctx.Request.ParseMultipartForm(32 << 20) // 32MB
	if ctx.Request.MultipartForm != nil {
		fmt.Printf("所有表单参数:\n")
		for key, values := range ctx.Request.MultipartForm.Value {
			fmt.Printf("  %s: %v\n", key, values)
		}
	}
	fmt.Printf("==================\n")

	// 上传文件
	var upload *models.Upload
	if articleTitle != "" {
		upload, err = c.uploadService.UploadImage(file, userID, uploadType, articleTitle)
	} else {
		upload, err = c.uploadService.UploadImage(file, userID, uploadType)
	}
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(ctx, upload.ToResponse())
}

// UploadAvatar 上传头像
// @Summary 上传头像
// @Description 上传用户头像
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param file formData file true "头像文件"
// @Success 200 {object} utils.Response{data=models.UploadResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/upload/avatar [post]
func (c *UploadController) UploadAvatar(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "请选择要上传的头像文件")
		return
	}

	// 上传头像
	upload, err := c.uploadService.UploadAvatar(file, userID)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, err.Error())
		return
	}

	utils.Success(ctx, upload.ToResponse())
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Description 删除已上传的文件
// @Tags 文件上传
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "上传记录ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/upload/{id} [delete]
func (c *UploadController) DeleteFile(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 获取上传记录ID
	uploadID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "上传记录ID格式错误")
		return
	}

	// 删除文件
	err = c.uploadService.DeleteFile(uint(uploadID), userID)
	if err != nil {
		if err.Error() == "文件不存在或无权限删除" {
			utils.Error(ctx, utils.CodeNotFound, err.Error())
		} else {
			utils.Error(ctx, utils.CodeBadRequest, err.Error())
		}
		return
	}

	utils.Success(ctx, "删除文件成功")
}

// GetMyUploads 获取当前用户的上传文件列表
// @Summary 获取当前用户的上传文件列表
// @Description 获取当前用户的上传文件列表，支持分页和类型筛选
// @Tags 文件上传
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param type query string false "上传类型" Enums(avatar,article,comment,other)
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=utils.PageResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/upload/my [get]
func (c *UploadController) GetMyUploads(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := middleware.GetCurrentUserID(ctx)
	if !exists {
		utils.Error(ctx, utils.CodeUnauthorized, "用户未登录")
		return
	}

	// 解析查询参数
	uploadType := ctx.Query("type")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	// 验证分页参数
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 10
	}

	// 获取上传文件列表
	uploads, total, err := c.uploadService.GetUserUploads(userID, uploadType, page, size)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, err.Error())
		return
	}

	// 转换为响应格式
	var uploadResponses []*models.UploadResponse
	for _, upload := range uploads {
		uploadResponses = append(uploadResponses, upload.ToResponse())
	}

	// 使用分页响应
	utils.SuccessPage(ctx, uploadResponses, total, page, size)
}

// GetUpload 获取上传文件详情
// @Summary 获取上传文件详情
// @Description 根据ID获取上传文件详情
// @Tags 文件上传
// @Accept json
// @Produce json
// @Param id path int true "上传记录ID"
// @Success 200 {object} utils.Response{data=models.UploadResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/upload/{id} [get]
func (c *UploadController) GetUpload(ctx *gin.Context) {
	// 获取上传记录ID
	uploadID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "上传记录ID格式错误")
		return
	}

	// 获取上传记录
	upload, err := c.uploadService.GetUploadByID(uint(uploadID))
	if err != nil {
		utils.Error(ctx, utils.CodeNotFound, "上传记录不存在")
		return
	}

	utils.Success(ctx, upload.ToResponse())
}