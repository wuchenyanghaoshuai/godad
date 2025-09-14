package models

import (
	"encoding/json"
	"time"
)

type ImageInfo struct {
	URL       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Size      int    `json:"size"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

type ChatEmoji struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Code      string    `json:"unicode" gorm:"column:code;not null;unique"`
	ImageURL  string    `json:"image_url" gorm:"not null"`
	Category  string    `json:"category" gorm:"type:enum('default','custom');default:'default'"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

type ChatConversation struct {
	ID                    uint       `json:"id" gorm:"primaryKey"`
	User1ID               uint       `json:"user1_id" gorm:"not null;index"`
	User2ID               uint       `json:"user2_id" gorm:"not null;index"`
	User1                 User       `json:"user1" gorm:"foreignKey:User1ID"`
	User2                 User       `json:"user2" gorm:"foreignKey:User2ID"`
	LastMessageID         *uint      `json:"last_message_id,omitempty"`
	LastMessageContent    *string    `json:"last_message_content,omitempty"`
	LastMessageType       string     `json:"last_message_type" gorm:"type:enum('text','image','emoji');default:'text'"`
	LastMessageAt         *time.Time `json:"last_message_time,omitempty" gorm:"column:last_message_at"`
	User1UnreadCount      uint       `json:"user1_unread_count" gorm:"default:0"`
	User2UnreadCount      uint       `json:"user2_unread_count" gorm:"default:0"`
	User1Deleted          bool       `json:"user1_deleted" gorm:"column:user1_deleted;default:false"`
	User2Deleted          bool       `json:"user2_deleted" gorm:"column:user2_deleted;default:false"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

type ChatMessage struct {
	ID                    uint       `json:"id" gorm:"primaryKey"`
	ConversationID        uint       `json:"conversation_id" gorm:"not null;index"`
	SenderID              uint       `json:"sender_id" gorm:"not null;index"`
	ReceiverID            uint       `json:"receiver_id" gorm:"not null;index"`
	MessageType           string     `json:"message_type" gorm:"type:enum('text','image','emoji');default:'text'"`
	Content               *string    `json:"content,omitempty"`
	Images                *string    `json:"images,omitempty"` // JSON stored as string
	EmojiID               *uint      `json:"emoji_id,omitempty"`
	IsRead                bool       `json:"is_read" gorm:"default:false"`
	ReadAt                *time.Time `json:"read_at,omitempty"`
	IsDeletedBySender     bool       `json:"is_deleted_by_sender" gorm:"default:false"`
	IsDeletedByReceiver   bool       `json:"is_deleted_by_receiver" gorm:"default:false"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	
	// Relations
	Conversation ChatConversation `json:"conversation,omitempty" gorm:"foreignKey:ConversationID"`
	Sender       User             `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
	Receiver     User             `json:"receiver,omitempty" gorm:"foreignKey:ReceiverID"`
	Emoji        *ChatEmoji       `json:"emoji,omitempty" gorm:"foreignKey:EmojiID"`
}

// GetImagesArray parses the Images JSON string to ImageInfo array
func (cm *ChatMessage) GetImagesArray() ([]ImageInfo, error) {
	if cm.Images == nil || *cm.Images == "" {
		return nil, nil
	}
	
	var images []ImageInfo
	err := json.Unmarshal([]byte(*cm.Images), &images)
	if err != nil {
		return nil, err
	}
	
	return images, nil
}

// SetImagesArray converts ImageInfo array to JSON string
func (cm *ChatMessage) SetImagesArray(images []ImageInfo) error {
	if images == nil {
		cm.Images = nil
		return nil
	}
	
	jsonData, err := json.Marshal(images)
	if err != nil {
		return err
	}
	
	jsonString := string(jsonData)
	cm.Images = &jsonString
	return nil
}

type ChatDailyLimit struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	SenderID     uint      `json:"sender_id" gorm:"not null;index"`
	ReceiverID   uint      `json:"receiver_id" gorm:"not null;index"`
	Date         time.Time `json:"date" gorm:"not null;index"`
	MessageCount uint      `json:"message_count" gorm:"default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Helper functions
func GetConversationUsers(user1ID, user2ID uint) (uint, uint) {
	if user1ID > user2ID {
		return user2ID, user1ID
	}
	return user1ID, user2ID
}

func GetTodayDate() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// SetImages method for ChatMessage (alias to SetImagesArray)
func (cm *ChatMessage) SetImages(images []ImageInfo) error {
	return cm.SetImagesArray(images)
}

// GetDisplayContent returns content for display purposes
func (cm *ChatMessage) GetDisplayContent() string {
	switch cm.MessageType {
	case "text":
		if cm.Content != nil {
			return *cm.Content
		}
		return ""
	case "image":
		return "[图片]"
	case "emoji":
		return "[表情]"
	default:
		return ""
	}
}

// CanSendMessage checks if user can send message based on daily limit
func (cdl *ChatDailyLimit) CanSendMessage(maxMessages uint) bool {
	return cdl.MessageCount < maxMessages
}

// IncrementMessageCount increments the message count
func (cdl *ChatDailyLimit) IncrementMessageCount() {
	cdl.MessageCount++
}

// DefaultEmojis 默认表情数据
var DefaultEmojis = []ChatEmoji{
	{Name: "微笑", Code: "smile", ImageURL: "😀", Category: "default", SortOrder: 1, IsActive: true},
	{Name: "开心", Code: "happy", ImageURL: "😊", Category: "default", SortOrder: 2, IsActive: true},
	{Name: "大笑", Code: "laugh", ImageURL: "😂", Category: "default", SortOrder: 3, IsActive: true},
	{Name: "眨眼", Code: "wink", ImageURL: "😉", Category: "default", SortOrder: 4, IsActive: true},
	{Name: "色彩", Code: "cool", ImageURL: "😎", Category: "default", SortOrder: 5, IsActive: true},
	{Name: "可爱", Code: "cute", ImageURL: "🥰", Category: "default", SortOrder: 6, IsActive: true},
	{Name: "亲吻", Code: "kiss", ImageURL: "😘", Category: "default", SortOrder: 7, IsActive: true},
	{Name: "思考", Code: "thinking", ImageURL: "🤔", Category: "default", SortOrder: 8, IsActive: true},
	{Name: "哭泣", Code: "cry", ImageURL: "😭", Category: "default", SortOrder: 9, IsActive: true},
	{Name: "生气", Code: "angry", ImageURL: "😠", Category: "default", SortOrder: 10, IsActive: true},
	{Name: "爱心", Code: "heart", ImageURL: "❤️", Category: "default", SortOrder: 11, IsActive: true},
	{Name: "点赞", Code: "thumbs_up", ImageURL: "👍", Category: "default", SortOrder: 12, IsActive: true},
	{Name: "鼓掌", Code: "clap", ImageURL: "👏", Category: "default", SortOrder: 13, IsActive: true},
	{Name: "加油", Code: "muscle", ImageURL: "💪", Category: "default", SortOrder: 14, IsActive: true},
	{Name: "OK", Code: "ok", ImageURL: "👌", Category: "default", SortOrder: 15, IsActive: true},
}