package services

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"godad-backend/config"
	"godad-backend/models"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
)

// UploadService 上传服务
type UploadService struct {
	ossClient *oss.Client
	bucket     *oss.Bucket
}

// NewUploadService 创建上传服务实例
func NewUploadService() (*UploadService, error) {
	cfg := config.GetConfig()

	// 创建OSS客户端
	client, err := oss.New(cfg.OSS.Endpoint, cfg.OSS.AccessKeyID, cfg.OSS.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %v", err)
	}

	// 获取存储桶
	bucket, err := client.Bucket(cfg.OSS.BucketName)
	if err != nil {
		return nil, fmt.Errorf("获取OSS存储桶失败: %v", err)
	}

	return &UploadService{
		ossClient: client,
		bucket:    bucket,
	}, nil
}

// UploadImage 上传图片
func (s *UploadService) UploadImage(file *multipart.FileHeader, userID uint, usage string) (*models.Upload, error) {
	// 验证文件类型
	if !s.isValidImageType(file.Filename) {
		return nil, fmt.Errorf("不支持的文件类型，仅支持 jpg, jpeg, png, gif, webp")
	}

	// 验证文件大小（5MB限制）
	if file.Size > 5*1024*1024 {
		return nil, fmt.Errorf("文件大小不能超过5MB")
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	// 读取文件内容
	fileContent, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 验证文件的真实MIME类型
	if !s.isValidImageContent(fileContent) {
		return nil, fmt.Errorf("文件内容不是有效的图片格式")
	}

	// 检查文件头，防止恶意文件
	if s.containsMaliciousContent(fileContent) {
		return nil, fmt.Errorf("检测到不安全的文件内容")
	}

	// 生成系统文件名（UUID）
	systemName := s.generateSystemName(file.Filename)
	
	// 生成存储路径
	fileName := s.generateFileName(file.Filename, usage)

	// 上传到OSS
	err = s.bucket.PutObject(fileName, bytes.NewReader(fileContent))
	if err != nil {
		return nil, fmt.Errorf("上传文件到OSS失败: %v", err)
	}

	// 生成文件URL
	publicURL := s.generateFileURL(fileName)

	// 计算文件哈希
	hash := fmt.Sprintf("%x", sha256.Sum256(fileContent))
	
	// 获取MIME类型
	mimeType := http.DetectContentType(fileContent)
	
	// 检查是否已存在相同的文件（同一用户+相同哈希值）
	var existingUpload models.Upload
	err = config.GetDB().Where("user_id = ? AND file_hash = ?", userID, hash).First(&existingUpload).Error
	if err == nil {
		// 文件已存在，删除刚上传的重复文件，返回已存在的记录
		s.bucket.DeleteObject(fileName)
		return &existingUpload, nil
	}

	// 保存上传记录到数据库
	upload := &models.Upload{
		UserID:      userID,
		FileName:    file.Filename,
		SystemName:  systemName,
		FileSize:    file.Size,
		FileType:    models.GetFileTypeFromMime(mimeType),
		MimeType:    mimeType,
		FileHash:    hash,
		StoragePath: fileName,
		PublicURL:   publicURL,
		Usage:       usage,
		Status:      1,
	}

	if err := config.GetDB().Create(upload).Error; err != nil {
		// 如果数据库保存失败，删除已上传的文件
		s.bucket.DeleteObject(fileName)
		return nil, fmt.Errorf("保存上传记录失败: %v", err)
	}

	return upload, nil
}

// UploadAvatar 上传头像
func (s *UploadService) UploadAvatar(file *multipart.FileHeader, userID uint) (*models.Upload, error) {
	// 验证文件类型
	if !s.isValidImageType(file.Filename) {
		return nil, fmt.Errorf("不支持的文件类型，仅支持 jpg, jpeg, png, gif, webp")
	}

	// 验证文件大小（2MB限制）
	if file.Size > 2*1024*1024 {
		return nil, fmt.Errorf("头像文件大小不能超过2MB")
	}

	// 上传文件
	upload, err := s.UploadImage(file, userID, "avatar")
	if err != nil {
		return nil, err
	}

	// 更新用户头像
	err = config.GetDB().Model(&models.User{}).Where("id = ?", userID).Update("avatar", upload.PublicURL).Error
	if err != nil {
		return nil, fmt.Errorf("更新用户头像失败: %v", err)
	}

	return upload, nil
}

// DeleteFile 删除文件
func (s *UploadService) DeleteFile(uploadID uint, userID uint) error {
	// 查找上传记录
	var upload models.Upload
	err := config.GetDB().Where("id = ? AND user_id = ?", uploadID, userID).First(&upload).Error
	if err != nil {
		return fmt.Errorf("文件不存在或无权限删除")
	}

	// 从OSS删除文件
	err = s.bucket.DeleteObject(upload.StoragePath)
	if err != nil {
		return fmt.Errorf("删除OSS文件失败: %v", err)
	}

	// 从数据库删除记录
	err = config.GetDB().Delete(&upload).Error
	if err != nil {
		return fmt.Errorf("删除上传记录失败: %v", err)
	}

	return nil
}

// GetUserUploads 获取用户上传文件列表
func (s *UploadService) GetUserUploads(userID uint, usage string, page, size int) ([]*models.Upload, int64, error) {
	var uploads []*models.Upload
	var total int64

	// 构建查询
	query := config.GetDB().Model(&models.Upload{}).Where("user_id = ?", userID)
	if usage != "" {
		query = query.Where("usage = ?", usage)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取上传文件总数失败: %v", err)
	}

	// 分页查询
	offset := (page - 1) * size
	err = query.Order("created_at DESC").Offset(offset).Limit(size).Find(&uploads).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取上传文件列表失败: %v", err)
	}

	return uploads, total, nil
}

// GetUploadByID 根据ID获取上传记录
func (s *UploadService) GetUploadByID(uploadID uint) (*models.Upload, error) {
	var upload models.Upload
	err := config.GetDB().Where("id = ?", uploadID).First(&upload).Error
	if err != nil {
		return nil, fmt.Errorf("上传记录不存在")
	}

	return &upload, nil
}

// isValidImageType 验证图片类型
func (s *UploadService) isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	for _, validType := range validTypes {
		if ext == validType {
			return true
		}
	}
	return false
}

// generateFileName 生成文件名
func (s *UploadService) generateFileName(originalName, usage string) string {
	ext := filepath.Ext(originalName)
	uuid := uuid.New().String()
	timestamp := time.Now().Format("20060102")

	// 根据上传类型生成不同的路径
	var filePath string
	switch usage {
	case "avatar":
		// UserImage下不创建日期文件夹
		filePath = fmt.Sprintf("UserImage/%s%s", uuid, ext)
	case "article":
		// ContentImage下创建日期文件夹
		filePath = fmt.Sprintf("ContentImage/%s/%s%s", timestamp, uuid, ext)
	case "comment":
		// ContentImage下创建日期文件夹
		filePath = fmt.Sprintf("ContentImage/%s/%s%s", timestamp, uuid, ext)
	default:
		// 默认使用ContentImage并创建日期文件夹
		filePath = fmt.Sprintf("ContentImage/%s/%s%s", timestamp, uuid, ext)
	}

	return filePath
}

// generateSystemName 生成系统文件名（UUID）
func (s *UploadService) generateSystemName(originalName string) string {
	ext := filepath.Ext(originalName)
	uuid := uuid.New().String()
	return fmt.Sprintf("%s%s", uuid, ext)
}

// generateFileURL 生成文件URL
func (s *UploadService) generateFileURL(fileName string) string {
	cfg := config.GetConfig()
	
	// 如果配置了自定义域名，使用自定义域名
	if cfg.OSS.CustomDomain != "" {
		return fmt.Sprintf("https://%s/%s", cfg.OSS.CustomDomain, fileName)
	}
	
	// 否则使用默认的OSS域名
	return fmt.Sprintf("https://%s.%s/%s", cfg.OSS.BucketName, cfg.OSS.Endpoint, fileName)
}

// getFileExtension 获取文件扩展名
func (s *UploadService) getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if len(ext) > 0 {
		return ext[1:] // 去掉点号
	}
	return ""
}

// isValidImageContent 验证文件内容是否为有效的图片
func (s *UploadService) isValidImageContent(content []byte) bool {
	if len(content) < 12 {
		return false
	}

	// 检测MIME类型
	mimeType := http.DetectContentType(content)
	validMimeTypes := []string{
		"image/jpeg",
		"image/jpg", 
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, validType := range validMimeTypes {
		if mimeType == validType {
			return true
		}
	}

	// 额外检查文件头
	return s.checkImageHeaders(content)
}

// checkImageHeaders 检查图片文件头
func (s *UploadService) checkImageHeaders(content []byte) bool {
	if len(content) < 12 {
		return false
	}

	// JPEG文件头: FF D8 FF
	if content[0] == 0xFF && content[1] == 0xD8 && content[2] == 0xFF {
		return true
	}

	// PNG文件头: 89 50 4E 47 0D 0A 1A 0A
	if len(content) >= 8 && 
		content[0] == 0x89 && content[1] == 0x50 && content[2] == 0x4E && content[3] == 0x47 &&
		content[4] == 0x0D && content[5] == 0x0A && content[6] == 0x1A && content[7] == 0x0A {
		return true
	}

	// GIF文件头: 47 49 46 38 (GIF8)
	if len(content) >= 6 && string(content[0:4]) == "GIF8" {
		return true
	}

	// WebP文件头: 52 49 46 46 ... 57 45 42 50 (RIFF...WEBP)
	if len(content) >= 12 && 
		string(content[0:4]) == "RIFF" && string(content[8:12]) == "WEBP" {
		return true
	}

	return false
}

// containsMaliciousContent 检查是否包含恶意内容
func (s *UploadService) containsMaliciousContent(content []byte) bool {
	// 转换为字符串进行检查（只检查前1024字节）
	checkLength := len(content)
	if checkLength > 1024 {
		checkLength = 1024
	}
	
	contentStr := strings.ToLower(string(content[:checkLength]))
	
	// 检查恶意脚本标签
	maliciousPatterns := []string{
		"<script",
		"javascript:",
		"vbscript:",
		"onload=",
		"onerror=",
		"onclick=",
		"onmouseover=",
		"<iframe",
		"<object",
		"<embed",
		"<meta",
		"<?php",
		"<%",
		"<jsp:",
	}
	
	for _, pattern := range maliciousPatterns {
		if strings.Contains(contentStr, pattern) {
			return true
		}
	}
	
	return false
}