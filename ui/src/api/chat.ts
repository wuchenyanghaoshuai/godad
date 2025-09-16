// 聊天API接口
import { http } from './http'
import type { 
  ChatConversation, 
  ChatMessage, 
  ChatEmoji, 
  SendMessageRequest, 
  ConversationListResponse,
  MessageListResponse,
  CreateConversationRequest,
  ChatListParams 
} from './chatTypes'

// 聊天API类
export class ChatAPI {
  // 获取对话列表
  static async getConversations(params: ChatListParams = {}) {
    return http.get<ConversationListResponse>('/chat/conversations', {
      page: params.page || 1,
      limit: params.limit || 20
    })
  }

  // 获取或创建对话
  static async getOrCreateConversation(data: CreateConversationRequest) {
    return http.post<ChatConversation>('/chat/conversations', data)
  }

  // 删除对话
  static async deleteConversation(conversationId: number) {
    return http.delete(`/chat/conversations/${conversationId}`)
  }

  // 获取对话消息列表
  static async getMessages(conversationId: number, params: ChatListParams = {}) {
    return http.get<MessageListResponse>(`/chat/conversations/${conversationId}/messages`, {
      page: params.page || 1,
      limit: params.limit || 50
    })
  }

  // 发送消息
  static async sendMessage(data: SendMessageRequest) {
    return http.post<ChatMessage>('/chat/messages', data)
  }

  // 标记消息为已读
  static async markAsRead(conversationId: number) {
    return http.put(`/chat/conversations/${conversationId}/read`)
  }

  // 获取表情列表
  static async getEmojis() {
    return http.get<ChatEmoji[]>('/chat/emojis')
  }

  // 检查消息发送限制
  static async checkMessageLimit(receiverId: number) {
    return http.post<{
      can_send: boolean
      mutual_follow: boolean
      message_count: number
      daily_limit: number
    }>('/chat/check-limit', { receiver_id: receiverId })
  }
}

// 导出聊天相关API方法
export const {
  getConversations,
  getOrCreateConversation,
  deleteConversation,
  getMessages,
  sendMessage,
  markAsRead,
  getEmojis,
  checkMessageLimit
} = ChatAPI