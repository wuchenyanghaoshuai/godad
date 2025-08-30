// API拦截器 - 处理认证和token刷新
import { errorHandler } from './errorHandler'
import { http } from '@/api/http'
import { AuthApi } from '@/api/auth'
import { API_CONFIG } from '@/api/config'
import type { AppError } from './errorHandler'

// Token刷新状态管理
class TokenManager {
  private isRefreshing = false
  private failedQueue: Array<{
    resolve: (value?: any) => void
    reject: (error: any) => void
  }> = []

  // 处理认证错误
  async handleAuthError(error: AppError): Promise<void> {
    if (error.code === 'UNAUTHORIZED') {
      await this.handleTokenExpired()
    } else if (error.code === 'FORBIDDEN') {
      this.handleAccessDenied()
    }
  }

  // 处理token过期
  private async handleTokenExpired(): Promise<void> {
    const refreshToken = localStorage.getItem(API_CONFIG.AUTH.REFRESH_TOKEN_KEY)
    
    if (!refreshToken) {
      this.redirectToLogin()
      return
    }

    if (this.isRefreshing) {
      // 如果正在刷新token，将请求加入队列
      return new Promise<void>((resolve, reject) => {
        this.failedQueue.push({ resolve, reject })
      })
    }

    this.isRefreshing = true

    try {
      // 尝试刷新token
      const response = await AuthApi.refreshToken()
      const newToken = response.data.token
      
      // 更新token
      http.setToken(newToken)
      localStorage.setItem(API_CONFIG.AUTH.TOKEN_KEY, newToken)
      
      // 处理队列中的请求
      this.processQueue(null, newToken)
      
      errorHandler.showSuccess('登录状态已更新')
    } catch (refreshError) {
      // 刷新失败，清除所有认证信息
      this.processQueue(refreshError, null)
      this.redirectToLogin()
      errorHandler.showError('登录已过期，请重新登录')
    } finally {
      this.isRefreshing = false
    }
  }

  // 处理访问被拒绝
  private handleAccessDenied(): void {
    errorHandler.showError('您没有权限执行此操作')
    // 可以重定向到无权限页面
    // router.push('/unauthorized')
  }

  // 重定向到登录页
  private redirectToLogin(): void {
    http.clearToken()
    
    // 保存当前路径，登录后可以返回
    const currentPath = window.location.pathname
    if (currentPath !== '/login') {
      localStorage.setItem('redirect_after_login', currentPath)
    }
    
    // 重定向到登录页
    window.location.href = '/login'
  }

  // 处理队列中的请求
  private processQueue(error: any, token: string | null): void {
    this.failedQueue.forEach(({ resolve, reject }) => {
      if (error) {
        reject(error)
      } else {
        resolve()
      }
    })
    
    this.failedQueue = []
  }

  // 检查是否需要处理认证错误
  shouldHandleAuthError(error: AppError): boolean {
    return error.code === 'UNAUTHORIZED' || error.code === 'FORBIDDEN'
  }
}

// 创建token管理器实例
export const tokenManager = new TokenManager()

// API拦截器类
export class ApiInterceptor {
  // 初始化拦截器
  static init(): void {
    // 监听全局错误事件
    window.addEventListener('unhandledrejection', (event) => {
      const error = event.reason
      if (error?.isAppError && tokenManager.shouldHandleAuthError(error)) {
        tokenManager.handleAuthError(error)
        event.preventDefault() // 阻止默认的错误处理
      }
    })
  }

  // 处理API响应错误
  static async handleResponseError(error: AppError): Promise<void> {
    if (tokenManager.shouldHandleAuthError(error)) {
      await tokenManager.handleAuthError(error)
    }
  }
}

// 导出便捷方法
export const handleApiError = ApiInterceptor.handleResponseError
export const initApiInterceptor = ApiInterceptor.init