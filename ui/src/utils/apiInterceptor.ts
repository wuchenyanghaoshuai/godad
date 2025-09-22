// API拦截器（Cookie 模式）- 只负责全局鉴权错误的统一引导
import { errorHandler } from './errorHandler'
import type { AppError } from './errorHandler'

function shouldHandleAuthError(error: AppError): boolean {
  return (
    error.code === 'UNAUTHORIZED' ||
    error.code === 'AUTH_REQUIRED' ||
    error.code === 'AUTH_EXPIRED' ||
    error.code === 'FORBIDDEN' ||
    error.code === 'PERMISSION_DENIED'
  )
}

function redirectToLogin(): void {
  // 保存当前路径，登录后可以返回
  const currentPath = window.location.pathname + window.location.search
  if (currentPath !== '/login') {
    localStorage.setItem('redirect_after_login', currentPath)
  }
  window.location.href = '/login'
}

export class ApiInterceptor {
  static init(): void {
    // 监听未处理的 Promise 拒绝，统一鉴权错误引导
    window.addEventListener('unhandledrejection', (event) => {
      const error = event.reason
      if (error?.isAppError && shouldHandleAuthError(error)) {
        if (error.code === 'FORBIDDEN' || error.code === 'PERMISSION_DENIED') {
          errorHandler.showError('您没有权限执行此操作', { redirect: '/404' })
        } else {
          errorHandler.showError('请先登录', { redirect: '/login' })
          redirectToLogin()
        }
        event.preventDefault()
      }
    })
  }

  static async handleResponseError(error: AppError): Promise<void> {
    if (shouldHandleAuthError(error)) {
      if (error.code === 'FORBIDDEN' || error.code === 'PERMISSION_DENIED') {
        errorHandler.showError('您没有权限执行此操作', { redirect: '/404' })
      } else {
        errorHandler.showError('请先登录', { redirect: '/login' })
        redirectToLogin()
      }
    }
  }
}

export const handleApiError = ApiInterceptor.handleResponseError
export const initApiInterceptor = ApiInterceptor.init
