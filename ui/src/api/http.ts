// HTTP客户端封装
import { API_CONFIG } from './config'
import { errorHandler } from '@/utils/errorHandler'

// 定义响应数据结构
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 定义请求配置
export interface RequestConfig {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  data?: any
  params?: any
  headers?: Record<string, string>
  timeout?: number
}

// 简化的HTTP客户端类
export class HttpClient {
  private baseURL: string
  private timeout: number

  constructor() {
    this.baseURL = API_CONFIG.BASE_URL + API_CONFIG.API_VERSION
    this.timeout = API_CONFIG.TIMEOUT
  }

  // Cookie 鉴权下不再依赖本地注入 Authorization 头
  // 保留空实现以兼容旧调用
  private getToken(): string | null { return null }
  public setToken(_token: string): void { /* no-op in cookie mode */ }
  public clearToken(): void {
    // Cookie-only：仅清理本地缓存的用户信息（若有）
    if (typeof localStorage !== 'undefined') {
      try {
        localStorage.removeItem(API_CONFIG.AUTH.USER_INFO_KEY)
      } catch (e) {
        // ignore storage errors
      }
    }
  }

  // 构建请求头
  private buildHeaders(customHeaders?: Record<string, string>): Record<string, string> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...customHeaders
    }

    // Cookie 模式不再附加 Authorization 头

    return headers
  }

  // 发送请求的核心方法
  private async request<T>(config: RequestConfig): Promise<ApiResponse<T>> {
    const { url, method = 'GET', data, params, headers, timeout = this.timeout } = config
    
    const fullUrl = url.startsWith('http') ? url : `${this.baseURL}${url}`
    const requestHeaders = this.buildHeaders(headers)

    // 构建查询参数 - 过滤掉undefined值
    let queryString = ''
    if (params) {
      const filteredParams: Record<string, string> = {}
      Object.keys(params).forEach(key => {
        const value = params[key]
        if (value !== undefined && value !== null && value !== '') {
          filteredParams[key] = String(value)
        }
      })
      queryString = new URLSearchParams(filteredParams).toString()
    }
    const requestUrl = queryString ? `${fullUrl}?${queryString}` : fullUrl

    const requestInit: RequestInit = {
      method,
      headers: requestHeaders,
      signal: AbortSignal.timeout(timeout),
      credentials: 'include'
    }

    if (data && (method === 'POST' || method === 'PUT')) {
      if (data instanceof FormData) {
        // 如果是FormData，移除Content-Type让浏览器自动设置
        delete requestHeaders['Content-Type']
        requestInit.body = data
      } else {
        requestInit.body = JSON.stringify(data)
      }
    }

    try {
      const doFetch = async (): Promise<Response> => fetch(requestUrl, requestInit)
      let response = await doFetch()

      if (response.status === 401) {
        // 尝试静默刷新并重试一次
        try {
          await this.refreshSession()
          response = await doFetch()
        } catch (e) {
          // 刷新失败，按未授权处理
        }
      }

      if (!response.ok) {
        // 创建包含响应信息的错误对象
        const errorData = {
          response: {
            status: response.status,
            statusText: response.statusText,
            data: null as any
          },
          request: { url: requestUrl, method },
          message: `HTTP ${response.status}: ${response.statusText}`
        }
        
        // 尝试解析错误响应体
        try {
          errorData.response.data = await response.json()
        } catch {
          // 忽略JSON解析错误
        }
        
        // 使用错误处理器处理API错误
        const appError = errorHandler.handleApiError(errorData)
        throw appError
      }

      const result = await response.json()
      // 统一判定成功：后端约定 code=200 或 0 视为成功
      if (typeof result === 'object' && result) {
        const code = (result.code ?? 200)
        if (code !== 200 && code !== 0) {
          const appError = errorHandler.handleBusinessError(String(code), result.message || '接口返回错误', result)
          ;(appError as any).isAppError = true
          throw appError
        }
      }
      return result as ApiResponse<T>
    } catch (error) {
      // 如果不是我们的AppError，则通过错误处理器处理
      if (!(error as any).isAppError) {
        const appError = errorHandler.handleApiError(error)
        throw appError
      }
      throw error
    }
  }

  // 刷新登录状态（基于 httpOnly refresh cookie）
  private async refreshSession(): Promise<void> {
    try {
      const url = `${this.baseURL}${API_CONFIG.ENDPOINTS.AUTH.REFRESH}`
      const resp = await fetch(url, { method: 'POST', credentials: 'include' })
      if (!resp.ok) throw new Error('refresh failed')
      const data = await resp.json()
      const code = (data?.code ?? 200)
      if (code !== 200 && code !== 0) throw new Error(data?.message || 'refresh error')
    } catch (e) {
      this.clearToken()
      throw e
    }
  }

  // GET请求
  public get<T>(url: string, params?: any, headers?: Record<string, string>): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: 'GET', params, headers })
  }

  // POST请求
  public post<T>(url: string, data?: any, headers?: Record<string, string>): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: 'POST', data, headers })
  }

  // PUT请求
  public put<T>(url: string, data?: any, headers?: Record<string, string>): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: 'PUT', data, headers })
  }

  // DELETE请求
  public delete<T>(url: string, params?: any, headers?: Record<string, string>): Promise<ApiResponse<T>> {
    return this.request<T>({ url, method: 'DELETE', params, headers })
  }
}

// 创建HTTP客户端实例
export const http = new HttpClient()

// 导出类型
// RequestConfig类型已在上面定义并导出
