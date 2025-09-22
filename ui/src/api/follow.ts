// 关注相关API服务
import { http } from './http'
import { API_CONFIG } from './config'
import type { ApiResponse, FollowStats, FollowListResponse } from './types'

// 关注相关API
export class FollowApi {
  // 关注用户
  static async followUser(userId: number): Promise<ApiResponse<any>> {
    return http.post(`${API_CONFIG.ENDPOINTS.FOLLOW.BASE}/${userId}`)
  }

  // 取消关注用户
  static async unfollowUser(userId: number): Promise<ApiResponse<any>> {
    return http.delete(`${API_CONFIG.ENDPOINTS.FOLLOW.BASE}/${userId}`)
  }

  // 检查关注状态
  static async checkFollowStatus(userId: number): Promise<ApiResponse<{ is_following: boolean }>> {
    return http.get(`${API_CONFIG.ENDPOINTS.FOLLOW.STATUS}/${userId}`)
  }

  // 获取关注列表（我关注的人）
  static async getFollowing(params?: { page?: number; limit?: number }): Promise<ApiResponse<FollowListResponse>> {
    return http.get(API_CONFIG.ENDPOINTS.FOLLOW.FOLLOWING, params)
  }

  // 获取粉丝列表（关注我的人）
  static async getFollowers(params?: { page?: number; limit?: number }): Promise<ApiResponse<FollowListResponse>> {
    return http.get(API_CONFIG.ENDPOINTS.FOLLOW.FOLLOWERS, params)
  }

  // 获取关注统计信息（当前用户）
  static async getFollowStats(): Promise<ApiResponse<FollowStats>> {
    return http.get(API_CONFIG.ENDPOINTS.FOLLOW.STATS)
  }

  // 获取指定用户的关注统计信息（公开接口）
  static async getUserFollowStats(userId: number): Promise<ApiResponse<FollowStats>> {
    return http.get(`${API_CONFIG.ENDPOINTS.FOLLOW.STATS}/${userId}`)
  }

  // 获取互相关注列表
  static async getMutualFollows(params?: { page?: number; limit?: number }): Promise<ApiResponse<FollowListResponse>> {
    return http.get(API_CONFIG.ENDPOINTS.FOLLOW.MUTUAL, params)
  }
}
