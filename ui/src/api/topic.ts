import { http } from './http'
import { API_CONFIG } from './config'
import { normalizePageResponse } from './pagination'

export interface TopicItem {
  id: number
  name: string
  display_name: string
  description?: string
  color?: string
  icon?: string
  sort?: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface TopicCreateRequest {
  name: string
  display_name: string
  description?: string
  color?: string
  icon?: string
  sort?: number
  is_active?: boolean
}

export interface TopicUpdateRequest extends TopicCreateRequest {}

export class TopicApi {
  static async getActiveTopics() {
    return http.get<TopicItem[]>(API_CONFIG.ENDPOINTS.TOPIC.ACTIVE)
  }

  static async getAdminTopics(params?: { page?: number; size?: number; all?: boolean }) {
    const res = await http.get<any>(API_CONFIG.ENDPOINTS.TOPIC.ADMIN_BASE, params)
    // 后端返回 { data: { items, total, page, page_size, total_page } }
    const page = normalizePageResponse<TopicItem>(res)
    return { code: 200, message: 'success', data: page }
  }

  static async createTopic(data: TopicCreateRequest) {
    return http.post<TopicItem>(API_CONFIG.ENDPOINTS.TOPIC.ADMIN_BASE, data)
  }

  static async updateTopic(id: number, data: TopicUpdateRequest) {
    return http.put<TopicItem>(`${API_CONFIG.ENDPOINTS.TOPIC.ADMIN_BASE}/${id}`, data)
  }

  static async deleteTopic(id: number) {
    return http.delete<null>(`${API_CONFIG.ENDPOINTS.TOPIC.ADMIN_BASE}/${id}`)
  }
}

export default TopicApi

