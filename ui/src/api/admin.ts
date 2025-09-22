import { http } from './http'
import { normalizePageResponse } from './pagination'

export interface AdminArticleListParams {
  page?: number
  size?: number
  status?: string | number
  keyword?: string
}

export interface AdminUserListParams {
  page?: number
  size?: number
  keyword?: string
}

export const AdminApi = {
  async getStats() {
    return http.get<any>('/admin/stats')
  },

  async getArticlesPage(params: AdminArticleListParams = {}) {
    const res = await http.get<any>('/admin/articles', params)
    // 后端返回 { data: { articles, total } }
    const raw = res?.data ? { data: res.data.articles, total: res.data.total, page: params.page, size: params.size } : res
    const page = normalizePageResponse<any>(raw)
    return { code: 200, message: 'success', data: page }
  },

  async updateArticleStatus(id: number, status: number) {
    return http.put(`/admin/articles/${id}/status`, { status })
  },

  async deleteArticle(id: number) {
    return http.delete(`/admin/articles/${id}`)
  },

  async getUsersPage(params: AdminUserListParams = {}) {
    const res = await http.get<any>('/admin/users', params)
    // 后端返回 { data: { users, total } }
    const raw = res?.data ? { data: res.data.users, total: res.data.total, page: params.page, size: params.size } : res
    const page = normalizePageResponse<any>(raw)
    return { code: 200, message: 'success', data: page }
  },

  async updateUserStatus(id: number, status: number) {
    return http.put(`/admin/users/${id}/status`, { status })
  },
}

export default AdminApi

