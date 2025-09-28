// 用户相关API服务
import { http } from './http'
import { API_CONFIG } from './config'
import type {
  User,
  UserUpdateRequest,
  ChangePasswordRequest,
  ApiResponse,
  PaginatedResponse,
  ListParams
} from './types'
import { normalizePageResponse } from './pagination'

// 用户API服务类
export class UserApi {
  // 获取用户个人信息
  static async getProfile(): Promise<ApiResponse<User>> {
    return http.get<User>(API_CONFIG.ENDPOINTS.USER.PROFILE)
  }

  // 更新用户个人信息
  static async updateProfile(data: UserUpdateRequest): Promise<ApiResponse<User>> {
    return http.put<User>(API_CONFIG.ENDPOINTS.USER.UPDATE_PROFILE, data)
  }

  // 修改密码
  static async changePassword(data: ChangePasswordRequest): Promise<ApiResponse<null>> {
    return http.put<null>(API_CONFIG.ENDPOINTS.USER.CHANGE_PASSWORD, data)
  }

  // 获取用户公开信息
  static async getPublicInfo(userId: number): Promise<ApiResponse<User>> {
    return http.get<User>(`${API_CONFIG.ENDPOINTS.USER.PUBLIC_INFO}/${userId}`)
  }

  // 获取用户列表（管理员功能）
  static async getUserList(params?: ListParams): Promise<ApiResponse<PaginatedResponse<User>>> {
    return http.get<PaginatedResponse<User>>('/user/list', params)
  }

  // 检查昵称是否可用
  static async checkNickname(nickname: string): Promise<ApiResponse<{ nickname: string; exists: boolean; available: boolean }>> {
    return http.get<{ nickname: string; exists: boolean; available: boolean }>(
      `${API_CONFIG.ENDPOINTS.USER.CHECK_NICKNAME}?nickname=${encodeURIComponent(nickname)}`
    )
  }

  // 获取用户的文章列表
  static async getUserArticles(userId: number, params?: ListParams): Promise<ApiResponse<PaginatedResponse<any>>> {
    return http.get<PaginatedResponse<any>>(`${API_CONFIG.ENDPOINTS.USER.PUBLIC_INFO}/${userId}/articles`, params)
  }
  static async getUserArticlesPage(userId: number, params?: ListParams): Promise<ApiResponse<PaginatedResponse<any>>> {
    const resp = await http.get<any>(`${API_CONFIG.ENDPOINTS.USER.PUBLIC_INFO}/${userId}/articles`, params)
    const page = normalizePageResponse<any>(resp)
    return { code: 200, message: 'success', data: page }
  }

  // 获取用户的关注状态
  static async getFollowStatus(userId: number): Promise<ApiResponse<{ is_following: boolean }>> {
    return http.get<{ is_following: boolean }>(`${API_CONFIG.ENDPOINTS.USER.PUBLIC_INFO}/${userId}/follow-status`)
  }

  // 关注用户
  static async followUser(userId: number): Promise<ApiResponse<null>> {
    return http.post<null>(`${API_CONFIG.ENDPOINTS.USER.PUBLIC_INFO}/${userId}/follow`)
  }

  // 取消关注用户
  static async unfollowUser(userId: number): Promise<ApiResponse<null>> {
    return http.delete<null>(`${API_CONFIG.ENDPOINTS.USER.PUBLIC_INFO}/${userId}/follow`)
  }

  // 根据用户名获取用户信息（新的安全API）
  static async getUserByUsername(username: string): Promise<ApiResponse<User>> {
    return http.get<User>(`${API_CONFIG.ENDPOINTS.USER.BY_USERNAME}/${username}`)
  }

  // 根据用户名获取用户文章列表（新的安全API）
  static async getUserArticlesByUsername(username: string, params?: ListParams): Promise<ApiResponse<PaginatedResponse<any>>> {
    return http.get<PaginatedResponse<any>>(`${API_CONFIG.ENDPOINTS.USER.ARTICLES_BY_USERNAME}/${username}/articles`, params)
  }
  static async getUserArticlesByUsernamePage(username: string, params?: ListParams): Promise<ApiResponse<PaginatedResponse<any>>> {
    const resp = await http.get<any>(`${API_CONFIG.ENDPOINTS.USER.ARTICLES_BY_USERNAME}/${username}/articles`, params)
    const page = normalizePageResponse<any>(resp)
    return { code: 200, message: 'success', data: page }
  }

  // 公共用户搜索（仅返回安全字段）
  static async searchUsers(params: { keyword: string; page?: number; size?: number }) {
    return http.get<any>(`/user/search`, params)
  }
}
