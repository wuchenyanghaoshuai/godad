// 聊天相关类型定义
import type { User } from './types'

// 图片信息
export interface ImageInfo {
  url: string
  width: number
  height: number
  size: number
  thumbnail?: string
}

// 聊天表情
export interface ChatEmoji {
  id: number
  name: string
  code: string
  image_url: string
  category: string
  sort_order: number
  is_active: boolean
  created_at: string
  updated_at: string
}

// 聊天对话
export interface ChatConversation {
  id: number
  user1_id: number
  user2_id: number
  user1: User
  user2: User
  last_message_id?: number
  last_message_content?: string
  last_message_type?: string
  last_message_time?: string
  user1_unread_count: number
  user2_unread_count: number
  user1_deleted: boolean
  user2_deleted: boolean
  created_at: string
  updated_at: string
}

// 对话响应（包含对方用户信息）
export interface ConversationResponse {
  id: number
  user1_id: number
  user2_id: number
  last_message_id?: number
  last_message_content?: string
  last_message_type?: string
  last_message_time?: string
  user1_unread_count: number
  user2_unread_count: number
  user1_deleted: boolean
  user2_deleted: boolean
  created_at: string
  updated_at: string
  other_user: User
  unread_count: number
}

// 聊天消息
export interface ChatMessage {
  id: number
  conversation_id: number
  sender_id: number
  receiver_id: number
  message_type: 'text' | 'image' | 'emoji'
  content?: string
  images?: ImageInfo[]
  emoji_id?: number
  emoji?: ChatEmoji
  is_read: boolean
  read_at?: string
  is_deleted_by_sender: boolean
  is_deleted_by_receiver: boolean
  sender: User
  receiver: User
  created_at: string
  updated_at: string
}

// 发送消息请求
export interface SendMessageRequest {
  sender_id: number
  receiver_id: number
  message_type: 'text' | 'image' | 'emoji'
  content?: string
  images?: ImageInfo[]
  emoji_id?: number
}

// 创建对话请求
export interface CreateConversationRequest {
  other_user_id: number
}

// 聊天列表查询参数
export interface ChatListParams {
  page?: number
  limit?: number
}

// 对话列表响应
export interface ConversationListResponse {
  conversations: ConversationResponse[]
  pagination: {
    page: number
    limit: number
    total: number
    total_pages: number
  }
}

// 消息列表响应
export interface MessageListResponse {
  messages: ChatMessage[]
  pagination: {
    page: number
    limit: number
    total: number
    total_pages: number
  }
}