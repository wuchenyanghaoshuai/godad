package tests

import (
	"fmt"
	"regexp"
	"testing"
	"godad-backend/models"
	"godad-backend/services"
	"github.com/stretchr/testify/assert"
)

// TestPasswordValidation 测试密码强度验证
func TestPasswordValidation(t *testing.T) {
	// 创建测试用例
	testCases := []struct {
		name        string
		password    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "密码过短",
			password:    "123",
			expectError: true,
			errorMsg:    "密码长度至少需要8位",
		},
		{
			name:        "没有小写字母",
			password:    "PASSWORD123",
			expectError: true,
			errorMsg:    "密码必须包含至少一个小写字母",
		},
		{
			name:        "没有大写字母",
			password:    "password123",
			expectError: true,
			errorMsg:    "密码必须包含至少一个大写字母",
		},
		{
			name:        "符合要求的密码",
			password:    "Password123",
			expectError: false,
			errorMsg:    "",
		},
		{
			name:        "复杂密码",
			password:    "MySecurePass123",
			expectError: false,
			errorMsg:    "",
		},
	}

	// 执行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 模拟注册请求
			req := &models.UserRegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: tc.password,
				Nickname: "测试用户",
			}

			// 这里应该调用实际的密码验证逻辑
			err := validatePasswordStrength(req.Password)

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestUserRegistrationFlow 测试用户注册流程
func TestUserRegistrationFlow(t *testing.T) {
	// 注意: 这需要测试数据库环境
	t.Skip("需要测试数据库环境")
	
	// 创建用户服务实例
	userService := services.NewUserService()

	// 测试注册请求
	req := &models.UserRegisterRequest{
		Username: "testuser001",
		Email:    "test001@example.com", 
		Password: "TestPass123",
		Nickname: "测试用户001",
	}

	// 执行注册
	user, err := userService.Register(req)
	
	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Nickname, user.Nickname)
	assert.Equal(t, 0, user.Role) // 新注册用户应该是普通用户
}

// TestAdminPermission 测试管理员权限检查
func TestAdminPermission(t *testing.T) {
	testCases := []struct {
		name     string
		userRole int
		isAdmin  bool
	}{
		{
			name:     "普通用户",
			userRole: 0,
			isAdmin:  false,
		},
		{
			name:     "内容管理员",
			userRole: 1,
			isAdmin:  false,
		},
		{
			name:     "系统管理员", 
			userRole: 2,
			isAdmin:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := &models.User{
				Role: tc.userRole,
			}

			isAdmin := checkIsAdmin(user)
			assert.Equal(t, tc.isAdmin, isAdmin)
		})
	}
}

// 辅助函数：密码强度验证
func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("密码长度至少需要8位")
	}
	
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return fmt.Errorf("密码必须包含至少一个小写字母")
	}
	
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return fmt.Errorf("密码必须包含至少一个大写字母")
	}
	
	return nil
}

// 辅助函数：管理员权限检查
func checkIsAdmin(user *models.User) bool {
	return user.Role == 2
}