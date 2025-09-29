// 认证相关API服务
import { http } from './http'
import { API_CONFIG } from './config'
import type {
  UserRegisterRequest,
  UserLoginRequest,
  UserLoginResponse,
  User,
  ApiResponse
} from './types'

// 认证API服务类
export class AuthApi {
  // 用户注册
  static async register(data: UserRegisterRequest): Promise<ApiResponse<User>> {
    return http.post<User>(API_CONFIG.ENDPOINTS.AUTH.REGISTER, data)
  }

  // 用户登录
  static async login(data: UserLoginRequest): Promise<ApiResponse<UserLoginResponse>> {
    return http.post<UserLoginResponse>(API_CONFIG.ENDPOINTS.AUTH.LOGIN, data)
  }

  // 用户登出
  static async logout(): Promise<ApiResponse<null>> {
    return http.post<null>(API_CONFIG.ENDPOINTS.AUTH.LOGOUT)
  }

  // 刷新token（Cookie 模式，无需传参）
  static async refreshToken(): Promise<ApiResponse<{ token?: string; refresh_token?: string }>> {
    return http.post<{ token?: string; refresh_token?: string }>(
      API_CONFIG.ENDPOINTS.AUTH.REFRESH
    )
  }

  // 获取当前用户信息
  static async getCurrentUser(silent = false): Promise<ApiResponse<User>> {
    return http.get<User>(API_CONFIG.ENDPOINTS.USER.PROFILE, undefined, undefined, silent)
  }

  // 忘记密码
  static async forgotPassword(email: string): Promise<ApiResponse<null>> {
    return http.post<null>(API_CONFIG.ENDPOINTS.AUTH.FORGOT_PASSWORD, { email_or_username: email })
  }

  // 重置密码
  static async resetPassword(token: string, newPassword: string): Promise<ApiResponse<null>> {
    return http.post<null>(API_CONFIG.ENDPOINTS.AUTH.RESET_PASSWORD, { 
      token, 
      new_password: newPassword 
    })
  }
}

// 认证工具函数
export class AuthUtils {
  // 保存登录信息（Cookie 模式仅保存用户信息）
  static saveAuthData(user: User): void {
    localStorage.setItem(API_CONFIG.AUTH.USER_INFO_KEY, JSON.stringify(user))
  }

  // Cookie 模式下不再持久化 token（保留历史 API 兼容的空操作）
  // static saveToken(_) { /* no-op */ }

  // 保存用户信息
  static saveUser(user: User): void {
    localStorage.setItem(API_CONFIG.AUTH.USER_INFO_KEY, JSON.stringify(user))
  }

  // 清除登录信息
  static clearAuthData(): void {
    localStorage.removeItem(API_CONFIG.AUTH.USER_INFO_KEY)
    http.clearToken()
  }

  // 获取存储的用户信息
  static getStoredUser(): User | null {
    const userStr = localStorage.getItem(API_CONFIG.AUTH.USER_INFO_KEY)
    if (userStr) {
      try {
        return JSON.parse(userStr)
      } catch (error) {
        console.error('Failed to parse stored user info:', error)
        return null
      }
    }
    return null
  }

  // isLoggedIn/getToken 在 Cookie 模式下不再使用
}
