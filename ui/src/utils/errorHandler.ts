import type { Router } from 'vue-router'
import { toastBridge as toast } from './toastBridge'

// 错误类型定义
export interface AppError {
  code: string
  message: string
  details?: any
  timestamp: number
}

// 错误代码映射
const ERROR_MESSAGES: Record<string, string> = {
  // 网络错误
  'NETWORK_ERROR': '网络连接失败，请检查网络设置',
  'TIMEOUT_ERROR': '请求超时，请稍后重试',
  'CORS_ERROR': '跨域请求被阻止',
  
  // 认证错误
  'AUTH_REQUIRED': '请先登录',
  'AUTH_EXPIRED': '登录已过期，请重新登录',
  'AUTH_INVALID': '登录信息无效',
  'PERMISSION_DENIED': '权限不足',
  
  // 业务错误
  'VALIDATION_ERROR': '输入信息有误',
  'RESOURCE_NOT_FOUND': '请求的资源不存在',
  'RESOURCE_CONFLICT': '资源冲突',
  'RATE_LIMIT_EXCEEDED': '请求过于频繁，请稍后重试',
  
  // 服务器错误
  'SERVER_ERROR': '服务器内部错误',
  'SERVICE_UNAVAILABLE': '服务暂时不可用',
  'MAINTENANCE': '系统维护中，请稍后访问',
  
  // 默认错误
  'UNKNOWN_ERROR': '发生未知错误'
}

// HTTP状态码映射
const HTTP_ERROR_CODES: Record<number, string> = {
  400: 'VALIDATION_ERROR',
  401: 'AUTH_REQUIRED',
  403: 'PERMISSION_DENIED',
  404: 'RESOURCE_NOT_FOUND',
  409: 'RESOURCE_CONFLICT',
  429: 'RATE_LIMIT_EXCEEDED',
  500: 'SERVER_ERROR',
  502: 'SERVICE_UNAVAILABLE',
  503: 'SERVICE_UNAVAILABLE',
  504: 'TIMEOUT_ERROR'
}

// 错误处理器类
export class ErrorHandler {
  private static instance: ErrorHandler
  private router?: Router
  private errorBoundary?: any
  
  private constructor() {}
  
  static getInstance(): ErrorHandler {
    if (!ErrorHandler.instance) {
      ErrorHandler.instance = new ErrorHandler()
    }
    return ErrorHandler.instance
  }
  
  // 初始化
  init(router: Router, errorBoundary?: any) {
    this.router = router
    this.errorBoundary = errorBoundary
    
    // 全局错误监听
    window.addEventListener('error', this.handleGlobalError.bind(this))
    window.addEventListener('unhandledrejection', this.handleUnhandledRejection.bind(this))
  }
  
  // 处理API错误
  handleApiError(error: any): AppError {
    let code = 'UNKNOWN_ERROR'
    let message = ERROR_MESSAGES.UNKNOWN_ERROR
    
    if (error.response) {
      // HTTP错误响应
      const status = error.response.status
      code = HTTP_ERROR_CODES[status] || 'SERVER_ERROR'
      
      // 尝试从响应中获取错误信息
      if (error.response.data?.message) {
        message = error.response.data.message
      } else if (error.response.data?.error) {
        message = error.response.data.error
      } else {
        message = ERROR_MESSAGES[code] || ERROR_MESSAGES.UNKNOWN_ERROR
      }
    } else if (error.request) {
      // 网络错误
      code = 'NETWORK_ERROR'
      message = ERROR_MESSAGES.NETWORK_ERROR
    } else {
      // 其他错误
      message = error.message || ERROR_MESSAGES.UNKNOWN_ERROR
    }
    
    const appError: AppError = {
      code,
      message,
      details: error,
      timestamp: Date.now()
    }
    ;(appError as any).isAppError = true
    
    this.logError(appError)
    return appError
  }
  
  // 处理业务错误
  handleBusinessError(code: string, message?: string, details?: any): AppError {
    const appError: AppError = {
      code,
      message: message || ERROR_MESSAGES[code] || ERROR_MESSAGES.UNKNOWN_ERROR,
      details,
      timestamp: Date.now()
    }
    ;(appError as any).isAppError = true
    
    this.logError(appError)
    return appError
  }
  
  // 显示错误
  showError(error: AppError | string, options?: {
    showToast?: boolean
    showBoundary?: boolean
    redirect?: string
  }) {
    const { showToast = true, showBoundary = false, redirect } = options || {}
    
    let appError: AppError
    if (typeof error === 'string') {
      appError = {
        code: 'UNKNOWN_ERROR',
        message: error,
        timestamp: Date.now()
      }
    } else {
      appError = error
    }
    
    // 显示Toast通知
    if (showToast) {
      if (appError.code.includes('AUTH')) {
        toast.error(appError.message, {
          action: {
            label: '去登录',
            onClick: () => this.router?.push('/login')
          }
        })
      } else {
        toast.error(appError.message)
      }
    }
    
    // 显示错误边界
    if (showBoundary && this.errorBoundary) {
      this.errorBoundary.showError(appError.message, JSON.stringify(appError.details, null, 2))
    }
    
    // 处理特殊错误的重定向
    if (redirect) {
      this.router?.push(redirect)
    } else if (appError.code === 'AUTH_REQUIRED' || appError.code === 'AUTH_EXPIRED') {
      this.router?.push('/login')
    } else if (appError.code === 'PERMISSION_DENIED') {
      this.router?.push('/404')
    }
  }
  
  // 处理全局错误
  public handleGlobalError(event: ErrorEvent) {
    const error = this.handleBusinessError(
      'UNKNOWN_ERROR',
      event.message,
      {
        filename: event.filename,
        lineno: event.lineno,
        colno: event.colno,
        stack: event.error?.stack
      }
    )
    
    this.showError(error, { showBoundary: true })
  }
  
  // 处理未捕获的Promise拒绝
  private handleUnhandledRejection(event: PromiseRejectionEvent) {
    const error = this.handleBusinessError(
      'UNKNOWN_ERROR',
      event.reason?.message || '未处理的Promise拒绝',
      event.reason
    )
    
    this.showError(error, { showBoundary: true })
    event.preventDefault()
  }
  
  // 显示成功信息
  public showSuccess(message: string): void {
    toast.success(message)
  }

  // 记录错误
  private logError(error: AppError) {
    // 忽略正常的认证相关错误（401、403等），避免在控制台产生误导性的错误信息
    if (error.code === 'AUTH_REQUIRED' || error.code === 'PERMISSION_DENIED') {
      return
    }

    console.error('应用错误:', error)

    // 在生产环境中，可以将错误发送到监控服务
    if (import.meta.env.PROD) {
      // 发送错误到监控服务
      // this.sendToMonitoring(error)
    }
  }
  
  // 发送错误到监控服务（示例）
  private sendToMonitoring(_error: AppError) {
    // 实现错误上报逻辑
    // 例如发送到Sentry、LogRocket等服务
  }
}

// 导出单例实例
export const errorHandler = ErrorHandler.getInstance()

// 便捷方法
export const handleApiError = (error: any) => errorHandler.handleApiError(error)
export const handleBusinessError = (code: string, message?: string, details?: any) => 
  errorHandler.handleBusinessError(code, message, details)
export const showError = (error: AppError | string, options?: any) => 
  errorHandler.showError(error, options)
