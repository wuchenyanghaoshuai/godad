package main

import (
	"fmt"
	"log"
	"godad-backend/services"
)

func main() {
	fmt.Println("=== 邮件服务测试 ===")
	
	// 创建邮件服务实例
	emailService := services.NewEmailService()
	
	// 检查邮件服务配置状态
	if emailService.IsConfigured() {
		fmt.Println("✅ 邮件服务已配置")
		fmt.Println("开始测试真实邮件发送...")
		
		// 测试密码重置邮件
		testEmail := "test@example.com" // 替换为真实邮箱进行测试
		resetURL := "http://127.0.0.1:3333/reset-password?token=test-token-123"
		
		err := emailService.SendPasswordResetEmail(testEmail, resetURL)
		if err != nil {
			fmt.Printf("❌ 邮件发送失败: %v\n", err)
		} else {
			fmt.Printf("✅ 密码重置邮件发送成功到: %s\n", testEmail)
		}
		
		// 测试欢迎邮件
		err = emailService.SendWelcomeEmail(testEmail, "测试用户")
		if err != nil {
			fmt.Printf("❌ 欢迎邮件发送失败: %v\n", err)
		} else {
			fmt.Printf("✅ 欢迎邮件发送成功到: %s\n", testEmail)
		}
		
	} else {
		fmt.Println("⚠️  邮件服务未配置")
		fmt.Println("邮件将通过控制台模拟发送")
		
		// 显示支持的邮件服务提供商
		fmt.Println("\n支持的邮件服务提供商配置：")
		providers := services.GetSupportedProviders()
		for name, config := range providers {
			fmt.Printf("\n%s:\n", name)
			for key, value := range config {
				fmt.Printf("  %s: %v\n", key, value)
			}
		}
		
		fmt.Println("\n要启用真实邮件发送，请配置以下环境变量：")
		fmt.Println("SMTP_HOST=你的SMTP服务器")
		fmt.Println("SMTP_PORT=587")
		fmt.Println("SMTP_USERNAME=你的邮箱")
		fmt.Println("SMTP_PASSWORD=你的密码或应用专用密码")
		fmt.Println("SMTP_FROM=你的发件邮箱")
		fmt.Println("SMTP_FROM_NAME=GoDad育儿平台")
		
		// 测试控制台模拟发送
		fmt.Println("\n测试控制台模拟发送...")
		testEmail := "test@example.com"
		resetURL := "http://127.0.0.1:3333/reset-password?token=test-token-123"
		
		err := emailService.SendPasswordResetEmail(testEmail, resetURL)
		if err != nil {
			log.Printf("模拟发送失败: %v", err)
		}
	}
	
	fmt.Println("\n=== 测试完成 ===")
}