import { http } from './http'
import { API_CONFIG } from './config'
import { normalizePageResponse } from './pagination'

export interface Tag {
  id: number
  name: string
  color: string
  description: string
  usage_count: number
  created_at: string
  updated_at: string
}

export interface TagRequest {
  name: string
  color?: string
  description?: string
}

export interface PopularTag {
  tag: Tag
  usage_count: number
}

export interface TagStats {
  total_tags: number
  total_usage: number
  avg_usage: number
  unused_tags: number
}

export interface PaginatedTagResponse {
  data: Tag[]
  total: number
  page: number
  size: number
  total_pages: number
}

const tagApi = {
  // 获取标签列表
  getTags: (params?: { 
    page?: number
    size?: number
    search?: string
  }) => {
    return http.get<PaginatedTagResponse>(API_CONFIG.ENDPOINTS.TAG.BASE, params)
  },

  // 获取标签列表（统一分页结构）
  getTagsPage: async (params?: { page?: number; size?: number; search?: string }) => {
    const res = await http.get<any>(API_CONFIG.ENDPOINTS.TAG.BASE, params)
    const page = normalizePageResponse<Tag>(res)
    return { code: 200, message: 'success', data: page }
  },

  // 创建标签
  createTag: (data: TagRequest) => {
    return http.post<Tag>(API_CONFIG.ENDPOINTS.TAG.BASE, data)
  },

  // 更新标签
  updateTag: (id: number, data: TagRequest) => {
    return http.put<Tag>(`${API_CONFIG.ENDPOINTS.TAG.BASE}/${id}`, data)
  },

  // 删除标签
  deleteTag: (id: number) => {
    return http.delete(`${API_CONFIG.ENDPOINTS.TAG.BASE}/${id}`)
  },

  // 获取标签详情
  getTagById: (id: number) => {
    return http.get<Tag>(`${API_CONFIG.ENDPOINTS.TAG.BASE}/${id}`)
  },

  // 获取热门标签
  getPopularTags: (limit?: number) => {
    return http.get<PopularTag[]>(API_CONFIG.ENDPOINTS.TAG.POPULAR, { limit })
  },

  // 搜索标签
  searchTags: (query: string, limit?: number) => {
    return http.get<Tag[]>(API_CONFIG.ENDPOINTS.TAG.SEARCH, { q: query, limit })
  },

  // 获取标签统计
  getTagStats: () => {
    return http.get<TagStats>(API_CONFIG.ENDPOINTS.TAG.STATS)
  },

  // 获取标签下的文章
  getArticlesByTag: (tagId: number, params?: {
    page?: number
    size?: number
  }) => {
    return http.get(`${API_CONFIG.ENDPOINTS.TAG.BASE}/${tagId}/articles`, params)
  },

  // 获取文章的标签列表
  getArticleTags: (articleId: number) => {
    return http.get<Tag[]>(`${API_CONFIG.ENDPOINTS.ARTICLE.LIST}/${articleId}/tags`)
  }
}

export default tagApi
