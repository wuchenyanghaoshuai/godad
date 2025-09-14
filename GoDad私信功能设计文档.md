# GoDad类QQ聊天私信功能设计文档

## 1. 功能概述

### 1.1 功能介绍
为GoDad育儿知识分享平台添加类似QQ的完整私信聊天功能，支持富媒体消息、表情、实时通信等。

### 1.2 核心特性
- **多媒体消息**: 文本、图片、表情、文件、语音消息
- **实时通信**: WebSocket实现即时消息推送
- **聊天会话管理**: 会话列表、删除会话、置顶会话
- **消息状态**: 发送中、已送达、已读、发送失败
- **表情系统**: 内置表情包、自定义表情
- **消息操作**: 撤回、删除、复制、转发
- **基于关注关系的发送限制机制**
- **聊天记录本地存储和云端同步**

## 2. 业务规则

### 2.1 发送限制规则
| 关注关系状态 | 每日发送限制 | 说明 |
|-------------|-------------|------|
| 互相关注 | 无限制 | 双方都关注对方，正常发送 |
| 单向关注 | 3条/天 | 仅一方关注另一方，限制发送 |
| 未关注 | 3条/天 | 双方都未关注，限制发送 |

### 2.2 限制计算方式
- 限制按**发送者**维度计算（A向B发送的条数）
- 重置时间：每日00:00
- 超出限制时提示用户关注对方或请求对方关注

### 2.3 消息类型定义
- **text**: 纯文本消息
- **image**: 图片消息（支持多图）
- **emoji**: 表情消息
- **file**: 文件消息
- **voice**: 语音消息
- **system**: 系统消息（如撤回通知）

### 2.4 消息状态
- `sending`: 发送中
- `sent`: 已发送到服务器
- `delivered`: 已送达对方
- `read`: 已被对方读取
- `failed`: 发送失败
- `recalled`: 已撤回

### 2.5 会话管理规则
- **删除会话**: 仅在当前用户侧删除，对方仍可见
- **置顶会话**: 个人设置，不影响对方
- **清空聊天记录**: 物理删除历史消息
- **免打扰**: 接收消息但不推送通知

## 3. 数据库设计

### 3.1 私信消息表 (chat_messages)
```sql
CREATE TABLE chat_messages (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    conversation_id BIGINT UNSIGNED NOT NULL COMMENT '会话ID',
    sender_id BIGINT UNSIGNED NOT NULL COMMENT '发送者ID',
    receiver_id BIGINT UNSIGNED NOT NULL COMMENT '接收者ID',
    message_type ENUM('text', 'image', 'emoji', 'file', 'voice', 'system') DEFAULT 'text' COMMENT '消息类型',
    content TEXT COMMENT '消息内容',
    media_info JSON COMMENT '媒体信息(图片、文件等)',
    reply_to_id BIGINT UNSIGNED NULL COMMENT '回复的消息ID',
    status ENUM('sending', 'sent', 'delivered', 'read', 'failed', 'recalled') DEFAULT 'sent' COMMENT '消息状态',
    delivered_at DATETIME(3) NULL COMMENT '送达时间',
    read_at DATETIME(3) NULL COMMENT '已读时间',
    recalled_at DATETIME(3) NULL COMMENT '撤回时间',
    local_id VARCHAR(50) NULL COMMENT '客户端本地ID',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL COMMENT '软删除时间',
    
    INDEX idx_conversation_id (conversation_id),
    INDEX idx_sender_receiver (sender_id, receiver_id),
    INDEX idx_created_at (created_at),
    INDEX idx_status (status),
    INDEX idx_reply_to (reply_to_id),
    INDEX idx_deleted_at (deleted_at),
    FOREIGN KEY (conversation_id) REFERENCES chat_conversations(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (reply_to_id) REFERENCES chat_messages(id) ON DELETE SET NULL
) COMMENT='聊天消息表';
```

### 3.2 聊天会话表 (chat_conversations)
```sql
CREATE TABLE chat_conversations (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user1_id BIGINT UNSIGNED NOT NULL COMMENT '用户1ID（较小的用户ID）',
    user2_id BIGINT UNSIGNED NOT NULL COMMENT '用户2ID（较大的用户ID）',
    last_message_id BIGINT UNSIGNED NULL COMMENT '最后一条消息ID',
    last_message_at DATETIME(3) NOT NULL COMMENT '最后消息时间',
    user1_unread_count INT UNSIGNED DEFAULT 0 COMMENT '用户1未读消息数',
    user2_unread_count INT UNSIGNED DEFAULT 0 COMMENT '用户2未读消息数',
    user1_deleted_at DATETIME(3) NULL COMMENT '用户1删除对话时间',
    user2_deleted_at DATETIME(3) NULL COMMENT '用户2删除对话时间',
    user1_pinned BOOLEAN DEFAULT FALSE COMMENT '用户1是否置顶',
    user2_pinned BOOLEAN DEFAULT FALSE COMMENT '用户2是否置顶',
    user1_muted BOOLEAN DEFAULT FALSE COMMENT '用户1是否免打扰',
    user2_muted BOOLEAN DEFAULT FALSE COMMENT '用户2是否免打扰',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    
    UNIQUE KEY uk_users (user1_id, user2_id),
    INDEX idx_user1 (user1_id),
    INDEX idx_user2 (user2_id),
    INDEX idx_last_message_at (last_message_at),
    FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (last_message_id) REFERENCES chat_messages(id) ON DELETE SET NULL
) COMMENT='聊天会话表';
```

### 3.3 表情包表 (chat_emoji_packs)
```sql
CREATE TABLE chat_emoji_packs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '表情包名称',
    description TEXT COMMENT '表情包描述',
    cover_image VARCHAR(500) COMMENT '封面图片URL',
    is_default BOOLEAN DEFAULT FALSE COMMENT '是否默认表情包',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    sort_order INT UNSIGNED DEFAULT 0 COMMENT '排序',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    
    INDEX idx_is_default (is_default),
    INDEX idx_is_active (is_active),
    INDEX idx_sort_order (sort_order)
) COMMENT='表情包表';
```

### 3.4 表情表 (chat_emojis)
```sql
CREATE TABLE chat_emojis (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    pack_id BIGINT UNSIGNED NOT NULL COMMENT '表情包ID',
    name VARCHAR(100) NOT NULL COMMENT '表情名称',
    code VARCHAR(50) NOT NULL COMMENT '表情代码(如:smile:)',
    image_url VARCHAR(500) NOT NULL COMMENT '表情图片URL',
    image_width INT UNSIGNED DEFAULT 24 COMMENT '图片宽度',
    image_height INT UNSIGNED DEFAULT 24 COMMENT '图片高度',
    sort_order INT UNSIGNED DEFAULT 0 COMMENT '排序',
    usage_count INT UNSIGNED DEFAULT 0 COMMENT '使用次数',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    
    INDEX idx_pack_id (pack_id),
    INDEX idx_code (code),
    INDEX idx_sort_order (sort_order),
    INDEX idx_usage_count (usage_count),
    FOREIGN KEY (pack_id) REFERENCES chat_emoji_packs(id) ON DELETE CASCADE
) COMMENT='表情表';
```

### 3.5 聊天文件表 (chat_files)
```sql
CREATE TABLE chat_files (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    message_id BIGINT UNSIGNED NOT NULL COMMENT '消息ID',
    original_name VARCHAR(500) NOT NULL COMMENT '原始文件名',
    stored_name VARCHAR(200) NOT NULL COMMENT '存储文件名',
    file_path VARCHAR(1000) NOT NULL COMMENT '文件路径',
    file_size BIGINT UNSIGNED NOT NULL COMMENT '文件大小(bytes)',
    file_type VARCHAR(100) COMMENT '文件类型',
    mime_type VARCHAR(200) COMMENT 'MIME类型',
    file_hash VARCHAR(64) COMMENT '文件哈希值',
    thumbnail_path VARCHAR(1000) COMMENT '缩略图路径(图片/视频)',
    duration INT UNSIGNED COMMENT '时长(语音/视频,秒)',
    width INT UNSIGNED COMMENT '宽度(图片/视频)',
    height INT UNSIGNED COMMENT '高度(图片/视频)',
    upload_status ENUM('uploading', 'completed', 'failed') DEFAULT 'completed' COMMENT '上传状态',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    
    INDEX idx_message_id (message_id),
    INDEX idx_file_hash (file_hash),
    INDEX idx_upload_status (upload_status),
    FOREIGN KEY (message_id) REFERENCES chat_messages(id) ON DELETE CASCADE
) COMMENT='聊天文件表';
```

### 3.6 每日发送限制记录表 (chat_daily_limits)
```sql
CREATE TABLE chat_daily_limits (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    sender_id BIGINT UNSIGNED NOT NULL COMMENT '发送者ID',
    receiver_id BIGINT UNSIGNED NOT NULL COMMENT '接收者ID',
    message_count INT UNSIGNED DEFAULT 0 COMMENT '当日发送消息数量',
    limit_date DATE NOT NULL COMMENT '限制日期',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    
    UNIQUE KEY uk_sender_receiver_date (sender_id, receiver_id, limit_date),
    INDEX idx_sender_date (sender_id, limit_date),
    INDEX idx_limit_date (limit_date),
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
) COMMENT='聊天每日发送限制记录表';
```

## 4. API设计

### 4.1 发送私信
```http
POST /api/private-messages
Content-Type: application/json
Authorization: Bearer <token>

{
    "receiver_id": 123,
    "content": "你好，想请教一个育儿问题...",
    "message_type": "text"
}

Response:
{
    "code": 200,
    "message": "发送成功",
    "data": {
        "id": 456,
        "sender_id": 789,
        "receiver_id": 123,
        "content": "你好，想请教一个育儿问题...",
        "message_type": "text",
        "status": "sent",
        "created_at": "2025-09-12T10:30:00Z"
    }
}

Error Response (超出限制):
{
    "code": 403,
    "message": "今日向该用户发送私信已达上限(3条)，请关注对方或邀请对方关注您",
    "data": {
        "daily_limit": 3,
        "sent_today": 3,
        "is_mutual_follow": false,
        "can_send": false
    }
}
```

### 4.2 获取会话列表
```http
GET /api/private-conversations?page=1&limit=20
Authorization: Bearer <token>

Response:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "conversations": [
            {
                "id": 1,
                "other_user": {
                    "id": 123,
                    "username": "mom_lily",
                    "nickname": "莉莉妈妈",
                    "avatar": "https://example.com/avatar.jpg"
                },
                "last_message": {
                    "id": 456,
                    "content": "谢谢分享！",
                    "message_type": "text",
                    "sender_id": 123,
                    "created_at": "2025-09-12T10:30:00Z"
                },
                "unread_count": 2,
                "last_message_at": "2025-09-12T10:30:00Z"
            }
        ],
        "pagination": {
            "current_page": 1,
            "per_page": 20,
            "total": 5,
            "total_pages": 1
        }
    }
}
```

### 4.3 获取会话消息历史
```http
GET /api/private-conversations/:conversation_id/messages?page=1&limit=50
Authorization: Bearer <token>

Response:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "messages": [
            {
                "id": 456,
                "sender_id": 789,
                "receiver_id": 123,
                "content": "你好，想请教一个育儿问题...",
                "message_type": "text",
                "status": "read",
                "read_at": "2025-09-12T10:35:00Z",
                "created_at": "2025-09-12T10:30:00Z",
                "sender": {
                    "id": 789,
                    "username": "dad_wang",
                    "nickname": "王爸爸",
                    "avatar": "https://example.com/avatar2.jpg"
                }
            }
        ],
        "pagination": {
            "current_page": 1,
            "per_page": 50,
            "total": 25,
            "total_pages": 1
        }
    }
}
```

### 4.4 标记消息已读
```http
PUT /api/private-messages/mark-read
Content-Type: application/json
Authorization: Bearer <token>

{
    "conversation_id": 1,
    "message_ids": [456, 457, 458]
}

Response:
{
    "code": 200,
    "message": "标记成功",
    "data": {
        "marked_count": 3
    }
}
```

### 4.5 检查发送限制
```http
GET /api/private-messages/send-limit/:receiver_id
Authorization: Bearer <token>

Response:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "can_send": true,
        "daily_limit": 3,
        "sent_today": 1,
        "remaining": 2,
        "is_mutual_follow": false,
        "reset_at": "2025-09-13T00:00:00Z"
    }
}
```

### 4.6 删除会话
```http
DELETE /api/private-conversations/:conversation_id
Authorization: Bearer <token>

Response:
{
    "code": 200,
    "message": "删除成功"
}
```

## 5. 前端交互设计

### 5.1 入口设计
- **导航栏**: 添加私信图标，显示未读消息数量红点
- **用户主页**: 添加"发私信"按钮
- **文章页**: 作者信息区域添加"私信作者"功能

### 5.2 页面结构

#### 5.2.1 私信列表页 (/messages)
```
┌─────────────────────────────────┐
│ 私信                 [搜索框]    │
├─────────────────────────────────┤
│ ● 莉莉妈妈           谢谢分享！  │
│   2条未读           10:30      │
├─────────────────────────────────┤
│   王爸爸            好的，了解了 │
│   已读              昨天       │
├─────────────────────────────────┤
│   ...                          │
└─────────────────────────────────┘
```

#### 5.2.2 私信对话页 (/messages/:conversation_id)
```
┌─────────────────────────────────┐
│ ← 莉莉妈妈                      │
├─────────────────────────────────┤
│              你好，想请教问题    │
│              10:25 [我]        │
│                                │
│ 好的，什么问题呢？              │
│ 10:26 [莉莉妈妈]               │
│                                │
│              关于宝宝辅食...    │
│              10:27 [我] ✓已读  │
├─────────────────────────────────┤
│ [输入框]              [发送]    │
└─────────────────────────────────┘
```

### 5.3 限制提示设计
```
┌─────────────────────────────────┐
│ 发送限制                        │
├─────────────────────────────────┤
│ 今日向该用户发送私信已达上限     │
│ (3/3)                          │
│                                │
│ 💡 建议:                       │
│ • 关注对方可无限制发送私信       │
│ • 邀请对方关注您                │
│                                │
│ [去关注] [取消]                 │
└─────────────────────────────────┘
```

## 6. 技术实现要点

### 6.1 关注关系检查
```go
func (s *PrivateMessageService) checkMutualFollow(userID1, userID2 uint) (bool, error) {
    // 检查是否互相关注
    var count int64
    err := s.db.Model(&Follow{}).
        Where("((follower_id = ? AND followee_id = ?) OR (follower_id = ? AND followee_id = ?)) AND deleted_at IS NULL", 
              userID1, userID2, userID2, userID1).
        Count(&count).Error
    
    return count == 2, err // 需要双向关注记录
}
```

### 6.2 发送限制检查
```go
func (s *PrivateMessageService) checkDailyLimit(senderID, receiverID uint) (*DailyLimitInfo, error) {
    today := time.Now().Format("2006-01-02")
    
    // 检查是否互相关注
    isMutual, err := s.checkMutualFollow(senderID, receiverID)
    if err != nil {
        return nil, err
    }
    
    if isMutual {
        return &DailyLimitInfo{
            CanSend: true,
            IsMutualFollow: true,
            DailyLimit: 0, // 无限制
        }, nil
    }
    
    // 查询当日已发送数量
    var limit PrivateMessageDailyLimit
    s.db.FirstOrCreate(&limit, PrivateMessageDailyLimit{
        SenderID: senderID,
        ReceiverID: receiverID,
        LimitDate: today,
    })
    
    return &DailyLimitInfo{
        CanSend: limit.MessageCount < 3,
        SentToday: limit.MessageCount,
        Remaining: 3 - limit.MessageCount,
        DailyLimit: 3,
        IsMutualFollow: false,
    }, nil
}
```

### 6.3 实时通知集成
```go
// 发送私信后触发通知
func (s *PrivateMessageService) afterSendMessage(message *PrivateMessage) error {
    // 1. 更新会话信息
    s.updateConversation(message)
    
    // 2. 发送实时通知
    notification := &Notification{
        ReceiverID: message.ReceiverID,
        ActorID: message.SenderID,
        Type: "private_message",
        ResourceID: &message.ID,
        Message: "给您发送了私信",
    }
    s.notificationService.Create(notification)
    
    // 3. WebSocket推送（如果接收者在线）
    s.wsService.SendToUser(message.ReceiverID, WSMessage{
        Type: "new_private_message",
        Data: message,
    })
    
    return nil
}
```

## 7. 安全防护

### 7.1 内容安全
- **敏感词过滤**: 集成敏感词检测
- **垃圾信息防护**: 频率限制、内容相似度检查
- **举报机制**: 支持用户举报不当私信

### 7.2 权限控制
- **身份认证**: JWT令牌验证
- **操作权限**: 只能查看/操作自己参与的会话
- **隐私设置**: 支持设置"仅关注的人可发私信"

### 7.3 数据保护
- **消息加密**: 敏感信息传输加密
- **软删除**: 支持消息软删除恢复
- **数据清理**: 定期清理过期限制记录

## 8. 边界情况处理

### 8.1 用户关系变化
- **取消关注**: 实时更新发送限制状态
- **被拉黑**: 禁止发送私信
- **账号注销**: 清理相关私信数据

### 8.2 消息状态异常
- **发送失败**: 提供重试机制
- **消息丢失**: 发送状态跟踪和重传
- **离线消息**: 支持离线消息推送

### 8.3 系统限制
- **消息长度**: 限制单条消息2000字符
- **发送频率**: 每分钟最多发送10条消息
- **存储容量**: 单个会话最多保留1万条消息

## 9. 开发计划

### 9.1 第一阶段 (核心功能)
- [ ] 数据库表创建和迁移
- [ ] 基础API接口开发
- [ ] 发送限制逻辑实现
- [ ] 前端私信列表页面
- [ ] 前端对话页面

### 9.2 第二阶段 (增强功能)  
- [ ] 实时消息推送
- [ ] 消息已读状态
- [ ] 图片消息支持
- [ ] 消息搜索功能

### 9.3 第三阶段 (优化完善)
- [ ] 性能优化
- [ ] 安全防护加强  
- [ ] 用户体验优化
- [ ] 数据统计分析

---

**文档版本**: v1.0  
**创建时间**: 2025-09-12  
**更新时间**: 2025-09-12  
**负责人**: Claude Code