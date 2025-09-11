import { http } from './http'

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
  page_size: number
  total_pages: number
}

const tagApi = {
  // 获取标签列表
  getTags: (params?: { 
    page?: number
    size?: number
    search?: string
  }) => {
    return http.get<PaginatedTagResponse>('/tags', { params })
  },

  // 创建标签
  createTag: (data: TagRequest) => {
    return http.post<Tag>('/tags', data)
  },

  // 更新标签
  updateTag: (id: number, data: TagRequest) => {
    return http.put<Tag>(`/tags/${id}`, data)
  },

  // 删除标签
  deleteTag: (id: number) => {
    return http.delete(`/tags/${id}`)
  },

  // 获取标签详情
  getTagById: (id: number) => {
    return http.get<Tag>(`/tags/${id}`)
  },

  // 获取热门标签
  getPopularTags: (limit?: number) => {
    return http.get<PopularTag[]>('/tags/popular', { limit })
  },

  // 搜索标签
  searchTags: (query: string, limit?: number) => {
    return http.get<Tag[]>('/tags/search', { q: query, limit })
  },

  // 获取标签统计
  getTagStats: () => {
    return http.get<TagStats>('/tags/stats')
  },

  // 获取标签下的文章
  getArticlesByTag: (tagId: number, params?: {
    page?: number
    size?: number
  }) => {
    return http.get(`/tags/${tagId}/articles`, { params })
  },

  // 获取文章的标签列表
  getArticleTags: (articleId: number) => {
    return http.get<Tag[]>(`/articles/${articleId}/tags`)
  }
}

export default tagApi