import { http } from './http'

export interface Like {
  id: number
  user_id: number
  target_type: 'article' | 'comment'
  target_id: number
  created_at: string
}

export interface LikeUser {
  id: number
  username: string
  avatar: string
  nickname?: string
  created_at: string
}

export interface LikeResponse {
  is_liked: boolean
  like_count: number
  message: string
}

export interface LikeStatus {
  is_liked: boolean
  like_count: number
}

export interface PopularContent {
  id: number
  title: string
  content: string
  like_count: number
  author: {
    id: number
    username: string
    avatar: string
  }
  created_at: string
}

export interface PaginatedLikeResponse {
  data: LikeUser[]
  total: number
  page: number
  size: number
  total_pages: number
}

class LikeApiService {
  // 点赞文章
  async like(articleId: number) {
    return http.post<LikeResponse>('/likes', {
      target_type: 'article',
      target_id: articleId
    })
  }

  // 取消点赞文章
  async unlike(articleId: number) {
    return http.delete<LikeResponse>(`/likes/article/${articleId}`)
  }

  // 点赞/取消点赞
  async toggleLike(targetType: 'article' | 'comment', targetId: number) {
    return http.post<LikeResponse>('/likes/toggle', {
      target_type: targetType,
      target_id: targetId
    })
  }

  // 获取点赞状态
  async getLikeStatus(targetType: 'article' | 'comment', targetId: number) {
    return http.get<LikeStatus>('/likes/status', {
      target_type: targetType,
      target_id: targetId
    })
  }

  // 获取点赞列表
  async getLikesByTarget(targetType: 'article' | 'comment', targetId: number, params?: {
    page?: number
    size?: number
  }) {
    return http.get<PaginatedLikeResponse>('/likes/list', {
      target_type: targetType,
      target_id: targetId,
      ...params
    })
  }

  // 获取用户点赞的内容
  async getUserLikes(targetType: 'article' | 'comment', params?: {
    page?: number
    size?: number
  }) {
    return http.get<PaginatedLikeResponse>('/likes/user', {
      target_type: targetType,
      ...params
    })
  }

  // 获取热门内容
  async getPopularContent(targetType: 'article' | 'comment', params?: {
    limit?: number
    days?: number
  }) {
    return http.get<PopularContent[]>('/likes/popular', {
      target_type: targetType,
      ...params
    })
  }

  // 批量获取点赞状态
  async batchGetLikeStatus(items: Array<{
    target_type: 'article' | 'comment'
    target_id: number
  }>) {
    return http.post<Record<string, LikeStatus>>('/likes/batch-status', { items })
  }
}

export const LikeApi = new LikeApiService()
export default LikeApi