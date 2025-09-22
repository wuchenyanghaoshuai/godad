package models

import (
	"time"
)

// Resource 资源模型
type Resource struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Title         string    `gorm:"size:255;not null" json:"title"`                 // 资源标题
	Description   string    `gorm:"type:text" json:"description"`                   // 资源描述
	Type          string    `gorm:"size:50;not null" json:"type"`                   // 资源类型: e-book, video, tool
	Category      string    `gorm:"size:100;not null" json:"category"`              // 分类: E-books, Videos, Tools
	Image         string    `gorm:"size:500" json:"image"`                          // 封面图片URL
	FileURL       string    `gorm:"size:500;not null" json:"file_url"`              // 资源文件URL
	ButtonText    string    `gorm:"size:100;default:'立即下载'" json:"button_text"`     // 按钮文本
	Status        int       `gorm:"default:0" json:"status"`                        // 状态: 0-待审核, 1-已发布, 2-已拒绝
	DownloadCount int       `gorm:"default:0" json:"download_count"`                // 下载次数
	UploaderID    *uint     `json:"uploader_id"`                                    // 上传者ID (可为空，管理员上传)
	Uploader      *User     `gorm:"foreignKey:UploaderID" json:"uploader,omitempty"` // 上传者信息
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ResourceStatus 资源状态常量
const (
	ResourceStatusPending  = 0 // 待审核
	ResourceStatusApproved = 1 // 已发布
	ResourceStatusRejected = 2 // 已拒绝
)

// ResourceType 资源类型常量
const (
	ResourceTypeEBook = "e-book"
	ResourceTypeVideo = "video"
	ResourceTypeTool  = "tool"
)

// ResourceCategory 资源分类常量
const (
	ResourceCategoryEBooks = "E-books"
	ResourceCategoryVideos = "Videos"
	ResourceCategoryTools  = "Tools"
)

// BeforeCreate 创建前的钩子
func (r *Resource) BeforeCreate() error {
	if r.ButtonText == "" {
		switch r.Type {
		case ResourceTypeEBook:
			r.ButtonText = "立即下载"
		case ResourceTypeVideo:
			r.ButtonText = "立即观看"
		case ResourceTypeTool:
			r.ButtonText = "开始使用"
		default:
			r.ButtonText = "立即下载"
		}
	}
	return nil
}

// TableName 设置表名
func (Resource) TableName() string {
	return "resources"
}