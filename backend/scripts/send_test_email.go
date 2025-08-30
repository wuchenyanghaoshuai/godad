package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"
	
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}
	
	fmt.Println("=== 发送测试邮件 ===")
	
	// Get SMTP configuration
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")
	fromName := os.Getenv("SMTP_FROM_NAME")
	
	if host == "" || portStr == "" || username == "" || password == "" || from == "" {
		log.Fatal("SMTP配置不完整，请检查环境变量")
	}
	
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("SMTP端口配置错误")
	}
	
	fmt.Printf("SMTP配置:\n")
	fmt.Printf("- 服务器: %s:%d\n", host, port)
	fmt.Printf("- 用户名: %s\n", username)
	fmt.Printf("- 发件人: %s <%s>\n", fromName, from)
	
	// Test email details
	testEmail := "15106400242@163.com"
	subject := "【GoDad育儿平台】测试邮件 - 系统功能验证"
	body := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>测试邮件</title>
    <style>
        body { 
            font-family: Arial, sans-serif; 
            line-height: 1.6; 
            color: #333; 
            max-width: 600px; 
            margin: 0 auto; 
            padding: 20px;
        }
        .header { 
            background: linear-gradient(135deg, #ec4899 0%, #f97316 100%); 
            color: white; 
            padding: 30px; 
            text-align: center; 
            border-radius: 10px 10px 0 0;
        }
        .content { 
            background: #f9f9f9; 
            padding: 30px; 
            border-radius: 0 0 10px 10px;
            border: 1px solid #e0e0e0;
        }
        .footer { 
            text-align: center; 
            margin-top: 30px; 
            padding-top: 20px; 
            border-top: 1px solid #e0e0e0; 
            color: #666; 
            font-size: 14px;
        }
        .success {
            background: #d4edda;
            border: 1px solid #c3e6cb;
            color: #155724;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>✅ 邮件系统测试</h1>
        <p>GoDad育儿知识分享平台</p>
    </div>
    
    <div class="content">
        <h2>邮件系统测试成功！</h2>
        <p>这是一封测试邮件，用于验证GoDad育儿平台的邮件发送功能是否正常工作。</p>
        
        <div class="success">
            <p><strong>🎉 测试结果：</strong></p>
            <ul>
                <li>SMTP连接正常</li>
                <li>邮件模板渲染正常</li>
                <li>邮件发送功能正常</li>
                <li>系统配置正确</li>
            </ul>
        </div>
        
        <p><strong>系统信息：</strong></p>
        <ul>
            <li>发送时间: 2025-08-29</li>
            <li>邮件服务: 163邮箱</li>
            <li>平台: GoDad育儿知识分享平台</li>
            <li>测试类型: 功能验证</li>
        </ul>
        
        <p>如果您收到这封邮件，说明邮件系统配置正确，可以正常发送邮件了！</p>
    </div>
    
    <div class="footer">
        <p>此邮件由系统自动发送用于测试</p>
        <p>© 2025 GoDad育儿知识分享平台</p>
    </div>
</body>
</html>`
	
	// Create email message
	message := gomail.NewMessage()
	message.SetAddressHeader("From", from, fromName)
	message.SetHeader("To", testEmail)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)
	
	// Create SMTP dialer with TLS settings
	dialer := gomail.NewDialer(host, port, username, password)
	
	// Try different configurations for 163 SMTP
	if port == 587 {
		// STARTTLS
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         host,
		}
		fmt.Println("使用STARTTLS模式 (端口587)")
	} else if port == 465 {
		// SSL/TLS
		dialer.SSL = true
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         host,
		}
		fmt.Println("使用SSL/TLS模式 (端口465)")
	}
	
	// Send email
	fmt.Printf("\n正在发送测试邮件到: %s\n", testEmail)
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Printf("❌ 邮件发送失败: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("✅ 测试邮件发送成功！\n")
	fmt.Printf("请检查邮箱 %s 是否收到邮件\n", testEmail)
	fmt.Println("\n=== 测试完成 ===")
}