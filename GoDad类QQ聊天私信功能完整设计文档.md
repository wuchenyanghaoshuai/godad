# GoDad类QQ聊天私信功能完整设计文档

## 1. 功能概述

### 1.1 功能介绍
为GoDad育儿知识分享平台添加类似QQ的完整私信聊天功能，支持富媒体消息、表情、实时通信等。

### 1.2 核心特性
- **多媒体消息**: 文本、图片、表情、文件、语音消息
- **实时通信**: WebSocket实现即时消息推送
- **聊天会话管理**: 会话列表、删除会话、置顶会话、免打扰
- **消息状态**: 发送中、已送达、已读、发送失败、已撤回
- **表情系统**: 内置表情包、自定义表情、常用表情
- **消息操作**: 撤回、删除、复制、转发、回复
- **文件管理**: 图片预览、文件下载、缩略图生成
- **基于关注关系的发送限制机制**
- **聊天记录本地存储和云端同步**

## 2. 业务规则

### 2.1 发送限制规则
| 关注关系状态 | 每日发送限制 | 说明 |
|-------------|-------------|------|
| 互相关注 | 无限制 | 双方都关注对方，正常发送 |
| 单向关注 | 3条/天 | 仅一方关注另一方，限制发送 |
| 未关注 | 3条/天 | 双方都未关注，限制发送 |

### 2.2 消息类型定义
- **text**: 纯文本消息，支持emoji和@提及
- **image**: 图片消息（支持多图，最多9张）
- **emoji**: 表情消息（动态表情包）
- **file**: 文件消息（文档、压缩包等）
- **voice**: 语音消息（最长60秒）
- **system**: 系统消息（如撤回通知、入群通知）

### 2.3 消息状态流转
```
发送中(sending) -> 已发送(sent) -> 已送达(delivered) -> 已读(read)
                              \-> 发送失败(failed)
                              \-> 已撤回(recalled)
```

### 2.4 文件存储策略

#### 2.4.1 图片存储
- **原图**: 保存用户上传的原始图片
- **缩略图**: 自动生成多种尺寸缩略图（150x150, 300x300）
- **压缩图**: 移动端展示用，控制在200KB以内
- **存储路径**: `/uploads/chat/images/{year}/{month}/{hash}.{ext}`

#### 2.4.2 文件存储
- **文件分类**: 按类型存储到不同目录（docs、media、other）
- **重复检测**: 通过文件hash避免重复存储
- **病毒扫描**: 上传文件自动扫描病毒
- **存储路径**: `/uploads/chat/files/{type}/{year}/{month}/{hash}.{ext}`

#### 2.4.3 语音存储
- **格式转换**: 统一转为mp3格式
- **压缩**: 采样率16kHz，比特率32kbps
- **时长限制**: 最长60秒
- **存储路径**: `/uploads/chat/voice/{year}/{month}/{hash}.mp3`

### 2.5 表情系统设计

#### 2.5.1 默认表情包
- **系统表情**: 微信风格emoji，支持皮肤颜色
- **经典表情**: QQ黄脸表情包（呲牙、流汗、鼓掌等）
- **动态表情**: GIF格式，文件大小控制在500KB内

#### 2.5.2 表情存储
- **CDN加速**: 表情图片使用CDN分发
- **本地缓存**: 客户端缓存常用表情
- **延迟加载**: 按需加载表情包内容

## 3. 数据库设计

### 3.1 聊天消息表 (chat_messages)
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

## 4. 前端交互设计

### 4.1 聊天主界面布局 (/chat)

```
┌─────────────────────────────────────────────────────┐
│ [≡] GoDad聊天                              [🔔] [👤]  │
├─────────────┬───────────────────────────────────────┤
│ 聊天列表     │ 聊天窗口                               │
│            │ ┌─────────────────────────────────────┐ │
│ 🔍 搜索     │ │ [←] 莉莉妈妈        [📞][🎥][⋯]  │ │
│            │ ├─────────────────────────────────────┤ │
│ 📌 置顶对话  │ │                                   │ │
│ ● 莉莉妈妈   │ │  你好，想请教个问题      10:25    │ │
│   谢谢分享！ │ │                              [我]  │ │
│   2条未读    │ │                                   │ │
│            │ │  什么问题呢？ 😊        10:26    │ │
│ ● 王爸爸     │ │  [莉莉妈妈]                      │ │
│   好的，了解 │ │                                   │ │
│   已读       │ │  关于宝宝辅食...        10:27    │ │
│            │ │                              [我] ✓│ │
│ 李妈妈      │ │                                   │ │
│ 张医生      │ │                                   │ │
│            │ ├─────────────────────────────────────┤ │
│            │ │ [📎] [输入框]        [😊] [🎤] [📤] │ │
│            │ └─────────────────────────────────────┘ │
└─────────────┴───────────────────────────────────────┘
```

### 4.2 消息气泡设计

#### 4.2.1 文本消息
```
对方消息:
┌─────────────────────────┐
│ [头像] 用户名 时间        │
│        ┌──────────────┐  │
│        │ 消息内容     │  │
│        └──────────────┘  │
└─────────────────────────┘

我的消息:
┌─────────────────────────┐
│        时间 [头像]      │
│  ┌──────────────┐       │
│  │ 消息内容     │       │
│  └──────────────┘       │
│              [✓✓] 已读  │
└─────────────────────────┘
```

#### 4.2.2 图片消息
```
┌─────────────────────────┐
│ [头像] 用户名 时间        │
│        ┌──────────────┐  │
│        │ [图片预览]   │  │
│        │ [多图网格]   │  │
│        └──────────────┘  │
│        点击查看大图      │
└─────────────────────────┘
```

#### 4.2.3 表情消息
```
┌─────────────────────────┐
│ [头像] 用户名 时间        │
│        [大表情 64x64]    │
└─────────────────────────┘
```

### 4.3 表情面板设计

```
┌─────────────────────────────────────────────────┐
│ [最近] [微信] [QQ经典] [动图] [收藏] [+添加]     │
├─────────────────────────────────────────────────┤
│ 😀 😃 😄 😁 😆 😅 😂 🤣 😊 😇 🙂 🙃 😉      │
│ 😌 😍 🥰 😘 😗 😙 😚 😋 😛 😝 😜 🤪 🤨      │
│ 🧐 🤓 😎 🤩 🥳 😏 😒 😞 😔 😟 😕 🙁 ☹️      │
│ 😣 😖 😫 😩 🥺 😢 😭 😤 😠 😡 🤬 🤯 😳      │
│ 🥵 🥶 😱 😨 😰 😥 😓 🤗 🤔 🤭 🤫 🤥 😶      │
└─────────────────────────────────────────────────┘
```

### 4.4 文件上传界面

```
拖拽上传区域:
┌─────────────────────────────────┐
│         📎 拖拽文件到此处        │
│           或点击选择文件         │
│                                │
│ 支持格式: 图片、文档、压缩包等    │
│ 单文件最大: 100MB               │
└─────────────────────────────────┘

图片选择器:
┌─────────────────────────────────┐
│ [本地图片] [相机] [表情包]       │
├─────────────────────────────────┤
│ □ □ □ □ □ □ □ □ □              │
│ □ □ □ □ □ □ □ □ □              │
│ □ □ □ □ □ □ □ □ □              │
├─────────────────────────────────┤
│     [预览] [取消] [发送(3)]      │
└─────────────────────────────────┘
```

### 4.5 消息操作菜单

```
长按/右键消息弹出菜单:
┌─────────────────┐
│ 📋 复制          │
│ ↩️ 回复          │
│ ↗️ 转发          │
│ ❤️ 收藏          │
│ 🚫 撤回          │
│ 🗑️ 删除          │
│ ℹ️ 消息信息       │
└─────────────────┘
```

### 4.6 会话设置面板

```
┌─────────────────────────────────┐
│ 聊天信息                        │
├─────────────────────────────────┤
│ [头像] 莉莉妈妈                 │
│        在线 · 上次活跃: 刚刚     │
├─────────────────────────────────┤
│ 📌 置顶聊天      [开关]         │
│ 🔕 消息免打扰    [开关]         │
│ 🔍 查找聊天记录                 │
│ 📂 查看文件                     │
├─────────────────────────────────┤
│ 🚫 举报用户                     │
│ 🗑️ 删除聊天                     │
│ 🚪 清空聊天记录                 │
└─────────────────────────────────┘
```

## 5. API设计

### 5.1 发送消息
```http
POST /api/chat/messages
Content-Type: application/json
Authorization: Bearer <token>

{
    "receiver_id": 123,
    "message_type": "text",
    "content": "你好，想请教一个育儿问题...",
    "reply_to_id": 456,
    "local_id": "local_msg_123456789"
}

Response:
{
    "code": 200,
    "message": "发送成功",
    "data": {
        "id": 789,
        "conversation_id": 12,
        "sender_id": 456,
        "receiver_id": 123,
        "message_type": "text",
        "content": "你好，想请教一个育儿问题...",
        "status": "sent",
        "local_id": "local_msg_123456789",
        "created_at": "2025-09-12T10:30:00Z"
    }
}
```

### 5.2 上传文件
```http
POST /api/chat/upload
Content-Type: multipart/form-data
Authorization: Bearer <token>

FormData:
- file: [文件数据]
- type: "image" | "file" | "voice"
- receiver_id: 123

Response:
{
    "code": 200,
    "message": "上传成功", 
    "data": {
        "file_id": 456,
        "file_url": "https://cdn.godad.com/chat/images/2025/09/abc123.jpg",
        "thumbnail_url": "https://cdn.godad.com/chat/thumbs/2025/09/abc123_300x300.jpg",
        "file_size": 1024000,
        "file_type": "image/jpeg",
        "width": 1920,
        "height": 1080
    }
}
```

### 5.3 获取表情包
```http
GET /api/chat/emojis/packs
Authorization: Bearer <token>

Response:
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "packs": [
            {
                "id": 1,
                "name": "默认表情",
                "cover_image": "https://cdn.godad.com/emojis/default/cover.png",
                "is_default": true,
                "emojis": [
                    {
                        "id": 1,
                        "code": ":smile:",
                        "name": "微笑",
                        "image_url": "https://cdn.godad.com/emojis/default/smile.png",
                        "width": 24,
                        "height": 24
                    }
                ]
            }
        ]
    }
}
```

## 6. WebSocket实时通信

### 6.1 连接建立
```javascript
const ws = new WebSocket('wss://api.godad.com/ws/chat');
ws.onopen = () => {
    // 发送认证信息
    ws.send(JSON.stringify({
        type: 'auth',
        token: localStorage.getItem('token')
    }));
};
```

### 6.2 消息推送格式
```javascript
// 新消息推送
{
    "type": "new_message",
    "data": {
        "id": 789,
        "conversation_id": 12,
        "sender_id": 456,
        "message_type": "text", 
        "content": "你好",
        "created_at": "2025-09-12T10:30:00Z",
        "sender": {
            "id": 456,
            "username": "dad_wang",
            "nickname": "王爸爸",
            "avatar": "https://cdn.godad.com/avatars/wang.jpg"
        }
    }
}

// 消息状态更新
{
    "type": "message_status",
    "data": {
        "message_id": 789,
        "status": "read",
        "read_at": "2025-09-12T10:35:00Z"
    }
}

// 用户在线状态
{
    "type": "user_status", 
    "data": {
        "user_id": 123,
        "status": "online",
        "last_seen": "2025-09-12T10:30:00Z"
    }
}

// 正在输入状态
{
    "type": "typing",
    "data": {
        "conversation_id": 12,
        "user_id": 123,
        "is_typing": true
    }
}
```

## 7. 技术实现要点

### 7.1 文件上传处理
```go
func (c *ChatController) UploadFile(ctx *gin.Context) {
    file, header, err := ctx.Request.FormFile("file")
    if err != nil {
        ctx.JSON(400, gin.H{"error": "文件上传失败"})
        return
    }
    defer file.Close()
    
    // 1. 文件类型检查
    if !isAllowedFileType(header.Header.Get("Content-Type")) {
        ctx.JSON(400, gin.H{"error": "不支持的文件类型"})
        return
    }
    
    // 2. 文件大小检查
    if header.Size > maxFileSize {
        ctx.JSON(400, gin.H{"error": "文件太大"})
        return
    }
    
    // 3. 生成文件哈希
    hasher := sha256.New()
    io.Copy(hasher, file)
    fileHash := hex.EncodeToString(hasher.Sum(nil))
    
    // 4. 检查重复文件
    if existingFile := checkDuplicateFile(fileHash); existingFile != nil {
        return existingFile
    }
    
    // 5. 保存文件
    filename := generateFilename(fileHash, header.Filename)
    savePath := filepath.Join(uploadPath, filename)
    
    if err := saveFile(file, savePath); err != nil {
        ctx.JSON(500, gin.H{"error": "文件保存失败"})
        return
    }
    
    // 6. 生成缩略图 (如果是图片)
    var thumbnailPath string
    if isImage(header.Header.Get("Content-Type")) {
        thumbnailPath = generateThumbnail(savePath)
    }
    
    // 7. 保存文件信息到数据库
    fileInfo := &ChatFile{
        OriginalName: header.Filename,
        StoredName: filename,
        FilePath: savePath,
        FileSize: header.Size,
        FileHash: fileHash,
        ThumbnailPath: thumbnailPath,
    }
    
    // 返回文件信息
    ctx.JSON(200, gin.H{"data": fileInfo})
}
```

### 7.2 消息状态更新
```go
func (s *ChatService) UpdateMessageStatus(messageID uint, status string, userID uint) error {
    // 1. 验证权限 - 只有接收者可以更新状态
    message, err := s.GetMessage(messageID)
    if err != nil {
        return err
    }
    
    if message.ReceiverID != userID {
        return errors.New("无权限更新消息状态")
    }
    
    // 2. 状态流转检查
    if !isValidStatusTransition(message.Status, status) {
        return errors.New("无效的状态转换")
    }
    
    // 3. 更新状态
    updates := map[string]interface{}{
        "status": status,
    }
    
    switch status {
    case "delivered":
        updates["delivered_at"] = time.Now()
    case "read":
        updates["read_at"] = time.Now()
        // 同时更新未读计数
        s.decrementUnreadCount(message.ConversationID, userID)
    }
    
    err = s.db.Model(&ChatMessage{}).
        Where("id = ?", messageID).
        Updates(updates).Error
        
    if err != nil {
        return err
    }
    
    // 4. WebSocket推送状态更新
    s.wsService.SendToUser(message.SenderID, WSMessage{
        Type: "message_status",
        Data: map[string]interface{}{
            "message_id": messageID,
            "status": status,
            "updated_at": time.Now(),
        },
    })
    
    return nil
}
```

### 7.3 表情解析
```javascript
// 前端表情解析
class EmojiParser {
    constructor() {
        this.emojiMap = new Map();
        this.loadEmojiData();
    }
    
    // 解析文本中的表情代码
    parseEmojis(text) {
        return text.replace(/:([a-zA-Z0-9_+-]+):/g, (match, code) => {
            const emoji = this.emojiMap.get(code);
            if (emoji) {
                return `<img src="${emoji.image_url}" 
                            alt="${emoji.name}" 
                            class="inline-emoji" 
                            width="${emoji.width}" 
                            height="${emoji.height}">`;
            }
            return match;
        });
    }
    
    // 获取表情选择器HTML
    getEmojiPickerHTML() {
        let html = '<div class="emoji-picker">';
        
        this.emojiPacks.forEach(pack => {
            html += `<div class="emoji-pack" data-pack-id="${pack.id}">`;
            html += `<h3>${pack.name}</h3>`;
            html += '<div class="emoji-grid">';
            
            pack.emojis.forEach(emoji => {
                html += `<img src="${emoji.image_url}" 
                             alt="${emoji.name}"
                             data-code="${emoji.code}"
                             class="emoji-item"
                             title="${emoji.name}">`;
            });
            
            html += '</div></div>';
        });
        
        html += '</div>';
        return html;
    }
}
```

## 8. 安全防护

### 8.1 文件安全
- **病毒扫描**: 集成ClamAV扫描上传文件
- **文件类型检查**: 严格限制允许的文件类型
- **文件大小限制**: 图片10MB，文件100MB，语音5MB
- **重复文件检测**: 通过hash避免存储重复文件

### 8.2 消息安全
- **XSS防护**: 所有用户输入进行HTML转义
- **敏感词过滤**: 实时过滤敏感内容
- **频率限制**: 防止消息轰炸
- **内容审核**: 可疑内容人工审核

### 8.3 隐私保护
- **端到端加密**: 敏感消息可选加密传输
- **阅后即焚**: 支持设置消息自动删除
- **屏蔽拉黑**: 完善的用户屏蔽机制

## 9. 性能优化

### 9.1 数据库优化
- **分表策略**: 按时间分表存储历史消息
- **索引优化**: 合理设置复合索引
- **数据清理**: 定期清理过期数据

### 9.2 缓存策略
- **Redis缓存**: 活跃会话、在线状态缓存
- **CDN加速**: 静态文件和表情CDN分发
- **本地存储**: 客户端缓存聊天记录

### 9.3 前端优化
- **虚拟滚动**: 长聊天记录虚拟滚动
- **图片懒加载**: 按需加载聊天图片
- **消息分页**: 分页加载历史消息

## 10. 开发计划

### 10.1 第一阶段 (基础功能)
- [ ] 数据库表创建和迁移
- [ ] WebSocket连接和认证
- [ ] 基础文本消息发送/接收
- [ ] 会话列表和聊天界面
- [ ] 发送限制逻辑

### 10.2 第二阶段 (多媒体)
- [ ] 图片上传和展示
- [ ] 表情系统
- [ ] 文件上传和下载
- [ ] 语音消息录制和播放

### 10.3 第三阶段 (高级功能)
- [ ] 消息撤回和删除
- [ ] 消息回复和转发
- [ ] 会话管理(置顶、免打扰)
- [ ] 搜索功能

### 10.4 第四阶段 (优化完善)
- [ ] 性能优化
- [ ] 安全加固
- [ ] 移动端适配
- [ ] 数据统计和监控

---

**文档版本**: v2.0  
**创建时间**: 2025-09-12  
**更新时间**: 2025-09-12  
**负责人**: Claude Code

这个设计文档涵盖了类QQ聊天功能的所有核心要素，包括多媒体消息、表情系统、文件存储、实时通信、安全防护等。你觉得这个设计如何？有什么需要调整或补充的地方吗？