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

  // 刷新token
  static async refreshToken(): Promise<ApiResponse<{ token: string; refresh_token: string }>> {
    const refreshToken = localStorage.getItem(API_CONFIG.AUTH.REFRESH_TOKEN_KEY)
    return http.post<{ token: string; refresh_token: string }>(
      API_CONFIG.ENDPOINTS.AUTH.REFRESH,
      { refresh_token: refreshToken }
    )
  }

  // 获取当前用户信息
  static async getCurrentUser(): Promise<ApiResponse<User>> {
    return http.get<User>(API_CONFIG.ENDPOINTS.USER.PROFILE)
  }

  // 忘记密码
  static async forgotPassword(email: string): Promise<ApiResponse<null>> {
    return http.post<null>(API_CONFIG.ENDPOINTS.AUTH.FORGOT_PASSWORD, { email })
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
  // 保存登录信息
  static saveAuthData(user: User, token: string): void {
    localStorage.setItem(API_CONFIG.AUTH.TOKEN_KEY, token)
    localStorage.setItem(API_CONFIG.AUTH.USER_INFO_KEY, JSON.stringify(user))
    http.setToken(token)
  }

  // 保存token
  static saveToken(token: string): void {
    localStorage.setItem(API_CONFIG.AUTH.TOKEN_KEY, token)
    http.setToken(token)
  }

  // 保存用户信息
  static saveUser(user: User): void {
    localStorage.setItem(API_CONFIG.AUTH.USER_INFO_KEY, JSON.stringify(user))
  }

  // 清除登录信息
  static clearAuthData(): void {
    localStorage.removeItem(API_CONFIG.AUTH.TOKEN_KEY)
    localStorage.removeItem(API_CONFIG.AUTH.REFRESH_TOKEN_KEY)
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

  // 检查是否已登录
  static isLoggedIn(): boolean {
    const token = localStorage.getItem(API_CONFIG.AUTH.TOKEN_KEY)
    return !!token
  }

  // 获取token
  static getToken(): string | null {
    return localStorage.getItem(API_CONFIG.AUTH.TOKEN_KEY)
  }
}