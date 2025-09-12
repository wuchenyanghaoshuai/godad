import { http } from './http'
import { API_CONFIG } from './config'
import type { ApiResponse, PaginatedResponse } from './types'

// 收藏相关类型定义
export interface Favorite {
  id: number
  user_id: number
  article_id: number
  created_at: string
  updated_at: string
  user?: {
    id: number
    username: string
    nickname?: string
    avatar?: string
  }
  article?: {
    id: number
    title: string
    slug: string
    excerpt?: string
    cover_image?: string
    author_id: number
    category_id?: number
    status: number
    view_count: number
    like_count: number
    comment_count: number
    favorite_count: number
    is_top: boolean
    created_at: string
    updated_at: string
  }
}

export interface FavoriteCreateRequest {
  article_id: number
}

export interface FavoriteListRequest {
  page?: number
  size?: number
  user_id?: number
  sort?: string
}

export interface FavoriteResponse {
  is_favorited: boolean
  favorite?: Favorite
}

export interface FavoriteListResponse {
  favorites: Favorite[]
  pagination: {
    total: number
    current_page: number
    per_page: number
    total_pages: number
  }
}

export interface BatchFavoriteStatusRequest {
  article_ids: number[]
}

export interface PopularFavoritesRequest {
  limit?: number
  days?: number
}

// 收藏API类
export class FavoriteApi {
  /**
   * 切换收藏状态（收藏/取消收藏）
   */
  static async toggleFavorite(data: FavoriteCreateRequest): Promise<ApiResponse<FavoriteResponse>> {
    return http.post<FavoriteResponse>(API_CONFIG.ENDPOINTS.FAVORITE.TOGGLE, data)
  }

  /**
   * 获取文章收藏状态
   */
  static async getFavoriteStatus(articleId: number): Promise<ApiResponse<{ is_favorited: boolean }>> {
    return http.get<{ is_favorited: boolean }>(`${API_CONFIG.ENDPOINTS.FAVORITE.STATUS}/${articleId}`)
  }

  /**
   * 批量获取收藏状态
   */
  static async batchGetFavoriteStatus(data: BatchFavoriteStatusRequest): Promise<ApiResponse<Record<number, boolean>>> {
    return http.post<Record<number, boolean>>(API_CONFIG.ENDPOINTS.FAVORITE.BATCH_STATUS, data)
  }

  /**
   * 获取我的收藏列表
   */
  static async getUserFavorites(params?: FavoriteListRequest): Promise<ApiResponse<FavoriteListResponse>> {
    return http.get<FavoriteListResponse>(API_CONFIG.ENDPOINTS.FAVORITE.MY, params)
  }

  /**
   * 获取文章的收藏列表
   */
  static async getArticleFavorites(articleId: number, params?: FavoriteListRequest): Promise<ApiResponse<FavoriteListResponse>> {
    return http.get<FavoriteListResponse>(`${API_CONFIG.ENDPOINTS.FAVORITE.ARTICLE_FAVORITES}/${articleId}`, params)
  }

  /**
   * 删除收藏
   */
  static async deleteFavorite(favoriteId: number): Promise<ApiResponse<null>> {
    return http.delete<null>(`${API_CONFIG.ENDPOINTS.FAVORITE.DELETE}/${favoriteId}`)
  }

  /**
   * 获取热门收藏
   */
  static async getPopularFavorites(params?: PopularFavoritesRequest): Promise<ApiResponse<any[]>> {
    return http.get<any[]>(API_CONFIG.ENDPOINTS.FAVORITE.POPULAR, params)
  }
}

// 收藏相关工具函数
export class FavoriteUtils {
  /**
   * 格式化收藏时间
   */
  static formatFavoriteTime(timeString: string): string {
    const time = new Date(timeString)
    const now = new Date()
    const diffInSeconds = Math.floor((now.getTime() - time.getTime()) / 1000)
    
    if (diffInSeconds < 60) {
      return '刚刚收藏'
    } else if (diffInSeconds < 3600) {
      const minutes = Math.floor(diffInSeconds / 60)
      return `${minutes}分钟前收藏`
    } else if (diffInSeconds < 86400) {
      const hours = Math.floor(diffInSeconds / 3600)
      return `${hours}小时前收藏`
    } else if (diffInSeconds < 2592000) {
      const days = Math.floor(diffInSeconds / 86400)
      return `${days}天前收藏`
    } else {
      return time.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      }) + '收藏'
    }
  }

  /**
   * 格式化收藏数量
   */
  static formatFavoriteCount(count: number): string {
    if (count < 1000) {
      return count.toString()
    } else if (count < 10000) {
      return `${(count / 1000).toFixed(1)}k`
    } else if (count < 100000) {
      return `${Math.floor(count / 1000)}k`
    } else {
      return `${(count / 10000).toFixed(1)}w`
    }
  }

  /**
   * 获取收藏状态的文本描述
   */
  static getFavoriteStatusText(isFavorited: boolean): string {
    return isFavorited ? '已收藏' : '收藏'
  }

  /**
   * 获取收藏操作的确认消息
   */
  static getFavoriteConfirmMessage(isFavorited: boolean): string {
    return isFavorited ? '确定要取消收藏这篇文章吗？' : '确定要收藏这篇文章吗？'
  }
}

// 导出类型和API
export default FavoriteApi
export type { Favorite, FavoriteResponse, FavoriteListResponse }