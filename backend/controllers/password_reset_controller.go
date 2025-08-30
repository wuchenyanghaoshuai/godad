package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"godad-backend/models"
	"godad-backend/services"
	"godad-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PasswordResetController 密码重置控制器
type PasswordResetController struct {
	db           *gorm.DB
	userService  *services.UserService
	emailService *services.EmailService
}

// NewPasswordResetController 创建密码重置控制器实例
func NewPasswordResetController(db *gorm.DB) *PasswordResetController {
	return &PasswordResetController{
		db:           db,
		userService:  services.NewUserService(),
		emailService: services.NewEmailService(),
	}
}

// ForgotPasswordRequest 忘记密码请求结构
type ForgotPasswordRequest struct {
	EmailOrUsername string `json:"email_or_username" binding:"required"`
}

// ResetPasswordRequest 重置密码请求结构
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ForgotPassword 处理忘记密码请求
func (c *PasswordResetController) ForgotPassword(ctx *gin.Context) {
	var req ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	// 检查用户是否存在（支持用户名或邮箱）
	var user models.User
	query := c.db.Where("status = ?", 1)
	if c.userService.IsEmail(req.EmailOrUsername) {
		query = query.Where("email = ?", req.EmailOrUsername)
	} else {
		query = query.Where("username = ?", req.EmailOrUsername)
	}
	
	if err := query.First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 为了安全考虑，即使用户不存在也返回成功
			utils.Success(ctx, gin.H{"message": "如果该用户存在，重置链接已发送到关联邮箱"})
			return
		}
		log.Printf("查询用户失败: %v", err)
		utils.Error(ctx, utils.CodeInternalError, "系统错误")
		return
	}

	// 生成重置令牌
	token, err := generateResetToken()
	if err != nil {
		log.Printf("生成重置令牌失败: %v", err)
		utils.Error(ctx, utils.CodeInternalError, "系统错误")
		return
	}

	// 创建密码重置记录（总是使用用户的邮箱）
	resetRecord := models.PasswordReset{
		Email:     user.Email,
		Token:     token,
		ExpiresAt: time.Now().Add(30 * time.Minute), // 30分钟后过期
	}

	if err := c.db.Create(&resetRecord).Error; err != nil {
		log.Printf("创建密码重置记录失败: %v", err)
		utils.Error(ctx, utils.CodeInternalError, "系统错误")
		return
	}

	// 发送重置邮件（总是发送到用户的邮箱）
	resetURL := fmt.Sprintf("http://127.0.0.1:3333/reset-password?token=%s", token)
	if err := c.emailService.SendPasswordResetEmail(user.Email, resetURL); err != nil {
		log.Printf("发送重置邮件失败: %v", err)
		// 不删除重置记录，因为即使邮件发送失败，用户仍可通过控制台获取链接
		// 只有在邮件服务完全不可用时才考虑失败处理
	}

	utils.Success(ctx, gin.H{"message": "重置链接已发送到您的邮箱"})
}

// ResetPassword 处理重置密码请求
func (c *PasswordResetController) ResetPassword(ctx *gin.Context) {
	var req ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, utils.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	// 查找重置记录
	var resetRecord models.PasswordReset
	if err := c.db.Where("token = ?", req.Token).First(&resetRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Error(ctx, utils.CodeBadRequest, "无效的重置令牌")
			return
		}
		log.Printf("查询重置记录失败: %v", err)
		utils.Error(ctx, utils.CodeInternalError, "系统错误")
		return
	}

	// 检查令牌是否有效
	if !resetRecord.IsValid() {
		if resetRecord.IsExpired() {
			utils.Error(ctx, utils.CodeBadRequest, "重置令牌已过期")
		} else {
			utils.Error(ctx, utils.CodeBadRequest, "重置令牌已使用")
		}
		return
	}

	// 查找用户
	var user models.User
	if err := c.db.Where("email = ?", resetRecord.Email).First(&user).Error; err != nil {
		log.Printf("查询用户失败: %v", err)
		utils.Error(ctx, utils.CodeInternalError, "用户不存在")
		return
	}

	// 加密新密码 - 使用与用户服务相同的加密方式
	hashedPassword := c.userService.HashPassword(req.NewPassword)

	// 开始事务
	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新用户密码
	if err := tx.Model(&user).Update("password", hashedPassword).Error; err != nil {
		tx.Rollback()
		log.Printf("更新密码失败: %v", err)
		utils.Error(ctx, utils.CodeInternalError, "密码更新失败")
		return
	}

	// 标记令牌为已使用
	resetRecord.MarkAsUsed()
	if err := tx.Save(&resetRecord).Error; err != nil {
		tx.Rollback()
		log.Printf("标记令牌失败: %v", err)
		utils.Error(ctx, utils.CodeInternalError, "系统错误")
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("事务提交失败: %v", err)
		utils.Error(ctx, utils.CodeInternalError, "系统错误")
		return
	}

	utils.Success(ctx, gin.H{"message": "密码重置成功"})
}

// generateResetToken 生成重置令牌
func generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// 移除了旧的 sendResetEmail 方法，现在使用 emailService

// CleanExpiredTokens 清理过期的重置令牌（可以作为定时任务）
func (c *PasswordResetController) CleanExpiredTokens() error {
	result := c.db.Where("expires_at < ?", time.Now()).Delete(&models.PasswordReset{})
	if result.Error != nil {
		return result.Error
	}
	
	log.Printf("清理了 %d 个过期的密码重置令牌", result.RowsAffected)
	return nil
}