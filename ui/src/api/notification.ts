import { http } from './http'

// 通知类型
export type NotificationType = 'like' | 'comment' | 'bookmark' | 'follow'

// 通知接口
export interface Notification {
  id: number
  receiver_id: number
  actor_id: number
  type: NotificationType
  resource_id?: number
  message: string
  is_read: boolean
  created_at: string
  updated_at: string
  
  // 扩展字段
  actor_username: string
  actor_nickname: string
  actor_avatar: string
  article_title?: string
  article_cover?: string
  comment_content?: string
}

// 通知统计
export interface NotificationStats {
  unread_count: number
  total_count: number
}

// 通知列表响应
export interface NotificationListResponse {
  notifications: Notification[]
  pagination: {
    current_page: number
    per_page: number
    total: number
    total_pages: number
  }
}

// 通知API类
export class NotificationApi {
  // 获取通知列表
  static async getNotifications(params?: {
    page?: number
    limit?: number
  }): Promise<{ code: number; message: string; data: NotificationListResponse }> {
    const queryParams = new URLSearchParams()
    if (params?.page) queryParams.append('page', params.page.toString())
    if (params?.limit) queryParams.append('limit', params.limit.toString())
    
    const url = `/notifications${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
    return http.get(url)
  }

  // 获取通知统计
  static async getNotificationStats(): Promise<{ code: number; message: string; data: NotificationStats }> {
    return http.get('/notifications/stats')
  }

  // 标记通知为已读
  static async markAsRead(notificationIds: number[]): Promise<{ code: number; message: string }> {
    return http.put('/notifications/mark-read', { 
      notification_ids: notificationIds 
    })
  }

  // 标记单个通知为已读
  static async markSingleAsRead(notificationId: number): Promise<{ code: number; message: string }> {
    return http.put(`/notifications/${notificationId}/mark-read`)
  }

  // 标记所有通知为已读
  static async markAllAsRead(): Promise<{ code: number; message: string }> {
    return http.put('/notifications/mark-all-read')
  }

  // 批量标记通知为已读（通过URL参数）
  static async batchMarkAsRead(notificationIds: number[]): Promise<{ code: number; message: string }> {
    return http.put(`/notifications/batch-mark-read?ids=${notificationIds.join(',')}`)
  }

  // 删除通知
  static async deleteNotification(notificationId: number): Promise<{ code: number; message: string }> {
    return http.delete(`/notifications/${notificationId}`)
  }

  // 删除所有通知
  static async deleteAllNotifications(): Promise<{ code: number; message: string }> {
    return http.delete('/notifications/all')
  }
}

// 通知类型中文映射
export const notificationTypeMap: Record<NotificationType, string> = {
  like: '点赞',
  comment: '评论', 
  bookmark: '收藏',
  follow: '关注'
}

// 通知类型图标映射
export const notificationIconMap: Record<NotificationType, string> = {
  like: '❤️',
  comment: '💬',
  bookmark: '🔖', 
  follow: '👤'
}

// 格式化通知时间
export const formatNotificationTime = (time: string): string => {
  const now = new Date()
  const notificationTime = new Date(time)
  const diff = now.getTime() - notificationTime.getTime()
  
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  
  return notificationTime.toLocaleDateString('zh-CN', {
    month: 'numeric',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

export default NotificationApi