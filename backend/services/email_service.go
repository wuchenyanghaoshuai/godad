package services

import (
	"crypto/tls"
	"fmt"
	"log"

	"godad-backend/utils"
	"gopkg.in/gomail.v2"
)

// EmailService 邮件服务
type EmailService struct {
	host     string
	port     int
	username string
	password string
	from     string
	fromName string
}

// EmailTemplate 邮件模板数据
type EmailTemplate struct {
	Subject string
	Body    string
}

// NewEmailService 创建邮件服务实例
func NewEmailService() *EmailService {
	return &EmailService{
		host:     utils.GetEnv("SMTP_HOST", "smtp.gmail.com"),
		port:     utils.GetEnvAsInt("SMTP_PORT", 587),
		username: utils.GetEnv("SMTP_USERNAME", ""),
		password: utils.GetEnv("SMTP_PASSWORD", ""),
		from:     utils.GetEnv("SMTP_FROM", ""),
		fromName: utils.GetEnv("SMTP_FROM_NAME", "GoDad育儿平台"),
	}
}

// SendPasswordResetEmail 发送密码重置邮件
func (e *EmailService) SendPasswordResetEmail(to, resetURL string) error {
	template := e.getPasswordResetTemplate(resetURL)
	return e.sendEmail(to, template.Subject, template.Body)
}

// SendWelcomeEmail 发送欢迎邮件
func (e *EmailService) SendWelcomeEmail(to, username string) error {
	template := e.getWelcomeTemplate(username)
	return e.sendEmail(to, template.Subject, template.Body)
}

// SendCommentNotificationEmail 发送评论通知邮件
func (e *EmailService) SendCommentNotificationEmail(to, recipientName, actorName, articleTitle, commentContent string) error {
	template := e.getCommentNotificationTemplate(recipientName, actorName, articleTitle, commentContent)
	return e.sendEmail(to, template.Subject, template.Body)
}

// SendReplyNotificationEmail 发送回复通知邮件
func (e *EmailService) SendReplyNotificationEmail(to, recipientName, actorName, commentContent string) error {
	template := e.getReplyNotificationTemplate(recipientName, actorName, commentContent)
	return e.sendEmail(to, template.Subject, template.Body)
}

// sendEmail 发送邮件
func (e *EmailService) sendEmail(to, subject, body string) error {
	// 检查配置
	if e.username == "" || e.password == "" || e.from == "" {
		log.Printf("邮件服务未配置，使用控制台模拟发送")
		e.logEmailToConsole(to, subject, body)
		return nil
	}

	// 创建邮件消息
	message := gomail.NewMessage()
	message.SetAddressHeader("From", e.from, e.fromName)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	// 创建SMTP拨号器
	dialer := gomail.NewDialer(e.host, e.port, e.username, e.password)
	
	// Configure TLS based on port
	if e.port == 465 {
		// SSL/TLS mode
		dialer.SSL = true
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         e.host,
		}
	} else if e.port == 587 {
		// STARTTLS mode
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         e.host,
		}
	}

	// 发送邮件
	if err := dialer.DialAndSend(message); err != nil {
		log.Printf("邮件发送失败: %v", err)
		// 发送失败时也记录到控制台作为备用
		e.logEmailToConsole(to, subject, body)
		return fmt.Errorf("邮件发送失败: %v", err)
	}

	log.Printf("邮件发送成功: %s -> %s", e.from, to)
	return nil
}

// getFromAddress 获取发件人地址
func (e *EmailService) getFromAddress() string {
	if e.fromName != "" {
		return fmt.Sprintf("%s <%s>", e.fromName, e.from)
	}
	return e.from
}

// logEmailToConsole 将邮件信息记录到控制台（用于开发和备用）
func (e *EmailService) logEmailToConsole(to, subject, body string) {
	log.Printf("=== 邮件发送记录 ===")
	log.Printf("收件人: %s", to)
	log.Printf("主题: %s", subject)
	log.Printf("内容: %s", body)
	log.Printf("==================")
}

// getPasswordResetTemplate 获取密码重置邮件模板
func (e *EmailService) getPasswordResetTemplate(resetURL string) EmailTemplate {
	subject := "【GoDad育儿平台】密码重置"
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>密码重置</title>
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
            background: linear-gradient(135deg, #e76f51 0%%, #f4a261 100%%);
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
        .button { 
            display: inline-block; 
            background: linear-gradient(135deg, #e76f51 0%%, #f4a261 100%%);
            color: white; 
            text-decoration: none; 
            padding: 15px 30px; 
            border-radius: 5px; 
            margin: 20px 0;
            font-weight: bold;
        }
        .footer { 
            text-align: center; 
            margin-top: 30px; 
            padding-top: 20px; 
            border-top: 1px solid #e0e0e0; 
            color: #666; 
            font-size: 14px;
        }
        .warning {
            background: #fff3cd;
            border: 1px solid #ffeaa7;
            color: #856404;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>🔒 密码重置</h1>
        <p>GoDad育儿知识分享平台</p>
    </div>
    
    <div class="content">
        <h2>您好！</h2>
        <p>我们收到了您重置密码的请求。如果这是您本人的操作，请点击下面的按钮重置密码：</p>
        
        <div style="text-align: center;">
            <a href="%s" class="button">重置密码</a>
        </div>
        
        <div class="warning">
            <p><strong>⚠️ 安全提示：</strong></p>
            <ul>
                <li>此链接30分钟内有效，过期后需重新申请</li>
                <li>此链接仅可使用一次</li>
                <li>如果不是您本人的操作，请忽略此邮件</li>
            </ul>
        </div>
        
        <p>如果按钮无法点击，请复制以下链接到浏览器地址栏：</p>
        <p style="word-break: break-all; background: #f5f5f5; padding: 10px; border-radius: 5px; font-family: monospace;">
            %s
        </p>
        
        <p>如果您没有请求重置密码，请忽略此邮件，您的账户是安全的。</p>
    </div>
    
    <div class="footer">
        <p>此邮件由系统自动发送，请勿回复</p>
        <p>© 2025 GoDad育儿知识分享平台</p>
    </div>
</body>
</html>
`, resetURL, resetURL)

	return EmailTemplate{
		Subject: subject,
		Body:    body,
	}
}

// getWelcomeTemplate 获取欢迎邮件模板
func (e *EmailService) getWelcomeTemplate(username string) EmailTemplate {
	subject := "【GoDad育儿平台】欢迎加入我们！"
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>欢迎加入</title>
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
            background: linear-gradient(135deg, #e76f51 0%%, #f4a261 100%%);
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
        .button { 
            display: inline-block; 
            background: linear-gradient(135deg, #e76f51 0%%, #f4a261 100%%);
            color: white; 
            text-decoration: none; 
            padding: 15px 30px; 
            border-radius: 5px; 
            margin: 20px 0;
            font-weight: bold;
        }
        .footer { 
            text-align: center; 
            margin-top: 30px; 
            padding-top: 20px; 
            border-top: 1px solid #e0e0e0; 
            color: #666; 
            font-size: 14px;
        }
        .features {
            display: flex;
            flex-wrap: wrap;
            gap: 15px;
            margin: 20px 0;
        }
        .feature {
            flex: 1;
            min-width: 200px;
            background: white;
            padding: 15px;
            border-radius: 8px;
            border: 1px solid #e0e0e0;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>🎉 欢迎加入GoDad！</h1>
        <p>专业的育儿知识分享平台</p>
    </div>
    
    <div class="content">
        <h2>亲爱的 %s，您好！</h2>
        <p>热烈欢迎您加入GoDad育儿知识分享平台大家庭！我们很高兴您选择与我们一起，共同学习和分享育儿经验。</p>
        
        <div class="features">
            <div class="feature">
                <h3>📚 丰富内容</h3>
                <p>专业育儿知识、经验分享</p>
            </div>
            <div class="feature">
                <h3>👥 温暖社区</h3>
                <p>与其他父母交流互动</p>
            </div>
            <div class="feature">
                <h3>🎯 个性推荐</h3>
                <p>根据您的需求推荐内容</p>
            </div>
        </div>
        
        <div style="text-align: center;">
            <a href="http://127.0.0.1:3333" class="button">开始探索</a>
        </div>
        
        <p><strong>接下来您可以：</strong></p>
        <ul>
            <li>完善个人资料，获得更精准的内容推荐</li>
            <li>浏览热门文章，学习育儿知识</li>
            <li>关注感兴趣的作者，第一时间获取更新</li>
            <li>参与讨论，分享您的育儿经验</li>
        </ul>
        
        <p>如果您有任何问题或建议，随时联系我们。祝您在GoDad度过愉快的时光！</p>
    </div>
    
    <div class="footer">
        <p>此邮件由系统自动发送，请勿回复</p>
        <p>© 2025 GoDad育儿知识分享平台</p>
    </div>
</body>
</html>
`, username)

	return EmailTemplate{
		Subject: subject,
		Body:    body,
	}
}

// getCommentNotificationTemplate 获取评论通知邮件模板
func (e *EmailService) getCommentNotificationTemplate(recipientName, actorName, articleTitle, commentContent string) EmailTemplate {
	subject := "【GoDad育儿平台】有人评论了您的文章"
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>评论通知</title>
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
            background: linear-gradient(135deg, #e76f51 0%%, #f4a261 100%%);
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
        .button { 
            display: inline-block; 
            background: linear-gradient(135deg, #e76f51 0%%, #f4a261 100%%);
            color: white; 
            text-decoration: none; 
            padding: 15px 30px; 
            border-radius: 5px; 
            margin: 20px 0;
            font-weight: bold;
        }
        .footer { 
            text-align: center; 
            margin-top: 30px; 
            padding-top: 20px; 
            border-top: 1px solid #e0e0e0; 
            color: #666; 
            font-size: 14px;
        }
        .comment-box {
            background: #fff;
            border-left: 4px solid #e76f51;
            padding: 15px;
            margin: 20px 0;
            border-radius: 5px;
            font-style: italic;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>💬 新评论通知</h1>
        <p>GoDad育儿知识分享平台</p>
    </div>
    
    <div class="content">
        <h2>亲爱的 %s，您好！</h2>
        <p><strong>%s</strong> 评论了您的文章《<strong>%s</strong>》：</p>
        
        <div class="comment-box">
            "%s"
        </div>
        
        <div style="text-align: center;">
            <a href="http://127.0.0.1:3333" class="button">查看详情并回复</a>
        </div>
        
        <p>您可以登录平台查看完整评论并进行回复，与其他育儿伙伴进行交流。</p>
    </div>
    
    <div class="footer">
        <p>此邮件由系统自动发送，请勿回复</p>
        <p>© 2025 GoDad育儿知识分享平台</p>
    </div>
</body>
</html>
`, recipientName, actorName, articleTitle, commentContent)

	return EmailTemplate{
		Subject: subject,
		Body:    body,
	}
}

// getReplyNotificationTemplate 获取回复通知邮件模板
func (e *EmailService) getReplyNotificationTemplate(recipientName, actorName, commentContent string) EmailTemplate {
	subject := "【GoDad育儿平台】有人回复了您的评论"
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>回复通知</title>
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
            background: linear-gradient(135deg, #e76f51 0%%, #f4a261 100%%);
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
        .button { 
            display: inline-block; 
            background: linear-gradient(135deg, #e76f51 0%%, #f4a261 100%%);
            color: white; 
            text-decoration: none; 
            padding: 15px 30px; 
            border-radius: 5px; 
            margin: 20px 0;
            font-weight: bold;
        }
        .footer { 
            text-align: center; 
            margin-top: 30px; 
            padding-top: 20px; 
            border-top: 1px solid #e0e0e0; 
            color: #666; 
            font-size: 14px;
        }
        .comment-box {
            background: #fff;
            border-left: 4px solid #e76f51;
            padding: 15px;
            margin: 20px 0;
            border-radius: 5px;
            font-style: italic;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>🔄 回复通知</h1>
        <p>GoDad育儿知识分享平台</p>
    </div>
    
    <div class="content">
        <h2>亲爱的 %s，您好！</h2>
        <p><strong>%s</strong> 回复了您的评论：</p>
        
        <div class="comment-box">
            "%s"
        </div>
        
        <div style="text-align: center;">
            <a href="http://127.0.0.1:3333" class="button">查看详情并回复</a>
        </div>
        
        <p>您可以登录平台查看完整对话并继续参与讨论。</p>
    </div>
    
    <div class="footer">
        <p>此邮件由系统自动发送，请勿回复</p>
        <p>© 2025 GoDad育儿知识分享平台</p>
    </div>
</body>
</html>
`, recipientName, actorName, commentContent)

	return EmailTemplate{
		Subject: subject,
		Body:    body,
	}
}


// IsConfigured 检查邮件服务是否已配置
func (e *EmailService) IsConfigured() bool {
	return e.username != "" && e.password != "" && e.from != ""
}

// GetSupportedProviders 获取支持的邮件服务提供商配置示例
func GetSupportedProviders() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"gmail": {
			"host": "smtp.gmail.com",
			"port": 587,
			"note": "需要开启两步验证并使用应用专用密码",
		},
		"qq": {
			"host": "smtp.qq.com", 
			"port": 587,
			"note": "需要开启SMTP服务并获取授权码",
		},
		"163": {
			"host": "smtp.163.com",
			"port": 587,
			"note": "需要开启SMTP服务并获取授权码",
		},
		"outlook": {
			"host": "smtp-mail.outlook.com",
			"port": 587,
			"note": "使用Outlook/Hotmail账户",
		},
	}
}