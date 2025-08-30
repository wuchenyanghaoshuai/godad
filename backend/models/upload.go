package models

import (
	"time"

	"gorm.io/gorm"
)

// Upload 文件上传模型
type Upload struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	FileName    string         `json:"file_name" gorm:"type:varchar(255);not null;comment:原始文件名"`
	SystemName  string         `json:"system_name" gorm:"type:varchar(255);not null;comment:系统生成的文件名(UUID)"`
	FileSize    int64          `json:"file_size" gorm:"type:bigint;not null;comment:文件大小(字节)"`
	FileType    string         `json:"file_type" gorm:"type:varchar(100);not null;comment:文件类型"`
	MimeType    string         `json:"mime_type" gorm:"type:varchar(100);not null;comment:MIME类型"`
	FileHash    string         `json:"file_hash" gorm:"type:varchar(64);index;not null;index:idx_user_file_hash;comment:文件哈希值"`
	StoragePath string         `json:"storage_path" gorm:"type:varchar(500);not null;comment:存储路径"`
	PublicURL   string         `json:"public_url" gorm:"type:varchar(500);not null;comment:公开访问URL"`
	UserID      uint           `json:"user_id" gorm:"not null;index;index:idx_user_file_hash,unique;comment:上传用户ID"`
	Usage       string         `json:"usage" gorm:"type:varchar(50);comment:用途(avatar,article,comment等)"`
	Status      int8           `json:"status" gorm:"type:tinyint;default:1;comment:状态 0-已删除 1-正常"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`

	// 关联关系
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// UploadRequest 文件上传请求
type UploadRequest struct {
	Usage string `form:"usage" binding:"required" example:"avatar"` // avatar, article, comment
}

// UploadListRequest 文件列表请求
type UploadListRequest struct {
	Page     int    `form:"page" binding:"min=1" example:"1"`
	Size     int    `form:"size" binding:"min=1,max=100" example:"10"`
	UserID   uint   `form:"user_id" example:"1"`
	Usage    string `form:"usage" example:"avatar"`
	FileType string `form:"file_type" example:"image"`
	Status   int8   `form:"status" binding:"min=0,max=1" example:"1"`
	Sort     string `form:"sort" example:"created_at desc"`
}

// UploadResponse 文件上传响应
type UploadResponse struct {
	ID          uint          `json:"id"`
	FileName    string        `json:"file_name"`
	SystemName  string        `json:"system_name"`
	FileSize    int64         `json:"file_size"`
	FileType    string        `json:"file_type"`
	MimeType    string        `json:"mime_type"`
	FileHash    string        `json:"file_hash"`
	StoragePath string        `json:"storage_path"`
	PublicURL   string        `json:"public_url"`
	UserID      uint          `json:"user_id"`
	Usage       string        `json:"usage"`
	Status      int8          `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	User        *UserResponse `json:"user,omitempty"`
}

// ToResponse 转换为响应格式
func (u *Upload) ToResponse() *UploadResponse {
	resp := &UploadResponse{
		ID:          u.ID,
		FileName:    u.FileName,
		SystemName:  u.SystemName,
		FileSize:    u.FileSize,
		FileType:    u.FileType,
		MimeType:    u.MimeType,
		FileHash:    u.FileHash,
		StoragePath: u.StoragePath,
		PublicURL:   u.PublicURL,
		UserID:      u.UserID,
		Usage:       u.Usage,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}

	// 包含用户信息
	if u.User.ID != 0 {
		resp.User = u.User.ToResponse()
	}

	return resp
}

// TableName 指定表名
func (Upload) TableName() string {
	return "uploads"
}

// GetFileTypeFromMime 根据MIME类型获取文件类型
func GetFileTypeFromMime(mimeType string) string {
	switch {
	case mimeType[:5] == "image":
		return "image"
	case mimeType[:5] == "video":
		return "video"
	case mimeType[:5] == "audio":
		return "audio"
	case mimeType == "application/pdf":
		return "pdf"
	case mimeType == "application/msword" || mimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return "document"
	default:
		return "other"
	}
}