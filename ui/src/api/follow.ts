// 关注相关API服务
import { http } from './http'
import type { ApiResponse, FollowStats, FollowUser, FollowListResponse } from './types'

// 关注相关API
export class FollowApi {
  // 关注用户
  static async followUser(userId: number): Promise<ApiResponse<any>> {
    return http.post(`/follows/${userId}`)
  }

  // 取消关注用户
  static async unfollowUser(userId: number): Promise<ApiResponse<any>> {
    return http.delete(`/follows/${userId}`)
  }

  // 检查关注状态
  static async checkFollowStatus(userId: number): Promise<ApiResponse<{ is_following: boolean }>> {
    return http.get(`/follows/status/${userId}`)
  }

  // 获取关注列表（我关注的人）
  static async getFollowing(params?: { page?: number; limit?: number }): Promise<ApiResponse<FollowListResponse>> {
    return http.get('/follows/following', params)
  }

  // 获取粉丝列表（关注我的人）
  static async getFollowers(params?: { page?: number; limit?: number }): Promise<ApiResponse<FollowListResponse>> {
    return http.get('/follows/followers', params)
  }

  // 获取关注统计信息（当前用户）
  static async getFollowStats(): Promise<ApiResponse<FollowStats>> {
    return http.get('/follows/stats')
  }

  // 获取指定用户的关注统计信息（公开接口）
  static async getUserFollowStats(userId: number): Promise<ApiResponse<FollowStats>> {
    return http.get(`/follows/stats/${userId}`)
  }

  // 获取互相关注列表
  static async getMutualFollows(params?: { page?: number; limit?: number }): Promise<ApiResponse<FollowListResponse>> {
    return http.get('/follows/mutual', params)
  }
}