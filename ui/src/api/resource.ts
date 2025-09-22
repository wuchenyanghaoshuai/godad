import { http } from './http'
import { API_CONFIG } from './config'
import type { ApiResponse } from './types'

// 资源类型定义
export interface Resource {
  id: number
  title: string
  description: string
  type: 'e-book' | 'video' | 'tool'
  category: string
  image?: string
  file_url: string
  button_text: string
  status: number
  download_count: number
  uploader_id?: number
  uploader?: {
    id: number
    username: string
    nickname?: string
    avatar?: string
  }
  created_at: string
  updated_at: string
}

// 资源列表响应类型
export interface ResourceListResponse {
  items: Resource[]
  total: number
  page: number
  size: number
  total_pages: number
}

// 资源统计响应类型
export interface ResourceStats {
  total: number
  published: number
  pending: number
  by_type: Array<{ type: string; count: number }>
  by_category: Array<{ category: string; count: number }>
}

// 获取资源请求参数
export interface GetResourcesParams {
  page?: number
  size?: number
  status?: string
  category?: string
  type?: string
  keyword?: string
}

// 创建/更新资源请求参数
export interface CreateResourceRequest {
  title: string
  description: string
  type: string
  category: string
  image?: string
  file_url: string
  button_text?: string
  status?: number
  uploader_id?: number
}

export interface UpdateResourceRequest {
  title?: string
  description?: string
  type?: string
  category?: string
  image?: string
  file_url?: string
  button_text?: string
  status?: number
}

// 资源状态常量
export const RESOURCE_STATUS = {
  PENDING: 0,
  APPROVED: 1,
  REJECTED: 2
} as const

// 资源API
export const ResourceApi = {
  // 获取已发布的资源列表（前端公开接口）
  async getPublishedResources(params?: GetResourcesParams): Promise<ApiResponse<ResourceListResponse>> {
    return http.get(API_CONFIG.ENDPOINTS.RESOURCE.BASE, params)
  },

  // 获取单个资源详情
  async getResource(id: number): Promise<ApiResponse<Resource>> {
    return http.get(`${API_CONFIG.ENDPOINTS.RESOURCE.BASE}/${id}`)
  },

  // 下载资源（增加下载次数）
  async downloadResource(id: number): Promise<ApiResponse<{ file_url: string; title: string }>> {
    return http.post(`${API_CONFIG.ENDPOINTS.RESOURCE.BASE}/${id}/download`)
  },

  // 用户上传资源（需要登录）
  async createResource(data: CreateResourceRequest): Promise<ApiResponse<Resource>> {
    return http.post(API_CONFIG.ENDPOINTS.RESOURCE.BASE, data)
  },

  // 管理员接口 - 获取所有资源列表（包括待审核）
  async getAllResources(params?: GetResourcesParams): Promise<ApiResponse<ResourceListResponse>> {
    return http.get(API_CONFIG.ENDPOINTS.RESOURCE.ADMIN_BASE, params)
  },

  // 管理员接口 - 创建资源
  async adminCreateResource(data: CreateResourceRequest): Promise<ApiResponse<Resource>> {
    return http.post(API_CONFIG.ENDPOINTS.RESOURCE.ADMIN_BASE, data)
  },

  // 管理员接口 - 更新资源
  async updateResource(id: number, data: UpdateResourceRequest): Promise<ApiResponse<Resource>> {
    return http.put(`${API_CONFIG.ENDPOINTS.RESOURCE.ADMIN_BASE}/${id}`, data)
  },

  // 管理员接口 - 更新资源状态（审核）
  async updateResourceStatus(id: number, status: number): Promise<ApiResponse<void>> {
    return http.put(`${API_CONFIG.ENDPOINTS.RESOURCE.ADMIN_BASE}/${id}/status`, { status })
  },

  // 管理员接口 - 删除资源
  async deleteResource(id: number): Promise<ApiResponse<void>> {
    return http.delete(`${API_CONFIG.ENDPOINTS.RESOURCE.ADMIN_BASE}/${id}`)
  },

  // 管理员接口 - 获取资源统计
  async getResourceStats(): Promise<ApiResponse<ResourceStats>> {
    return http.get(API_CONFIG.ENDPOINTS.RESOURCE.STATS)
  }
}
