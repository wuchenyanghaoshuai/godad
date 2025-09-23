// 分类相关API服务
import { http } from './http'
import { API_CONFIG } from './config'
import type {
  Category,
  CategoryRequest,
  ApiResponse,
  ListParams
} from './types'
import { normalizePageResponse } from './pagination'

// 分类API服务类
export class CategoryApi {
  // 获取分类列表（普通用户）
  static async getCategoryList(params?: ListParams): Promise<ApiResponse<Category[]>> {
    // 为前台使用返回“所有启用的分类（不分页）”，避免只拿到前10条
    return http.get<Category[]>(API_CONFIG.ENDPOINTS.CATEGORY.ALL)
  }

  // 获取分类列表（管理员）
  static async getAdminCategoryList(params?: ListParams): Promise<ApiResponse<Category[]>> {
    return http.get<Category[]>(API_CONFIG.ENDPOINTS.CATEGORY.ADMIN_BASE, params)
  }

  // 获取分类列表（管理员，统一分页结构）
  static async getAdminCategoryPage(params?: ListParams): Promise<ApiResponse<{ items: Category[]; total: number; page: number; size: number; total_pages: number }>> {
    const resp = await http.get<any>(API_CONFIG.ENDPOINTS.CATEGORY.ADMIN_BASE, params)
    const page = normalizePageResponse<Category>(resp)
    return { code: 200, message: 'success', data: page }
  }

  // 获取分类详情
  static async getCategoryDetail(id: number): Promise<ApiResponse<Category>> {
    return http.get<Category>(`/category/${id}`)
  }

  // 创建分类（管理员功能）
  static async createCategory(data: CategoryRequest): Promise<ApiResponse<Category>> {
    return http.post<Category>(API_CONFIG.ENDPOINTS.CATEGORY.CREATE, data)
  }

  // 更新分类（管理员功能）
  static async updateCategory(id: number, data: CategoryRequest): Promise<ApiResponse<Category>> {
    return http.put<Category>(`${API_CONFIG.ENDPOINTS.CATEGORY.UPDATE}/${id}`, data)
  }

  // 删除分类（管理员功能）
  static async deleteCategory(id: number): Promise<ApiResponse<null>> {
    return http.delete<null>(`${API_CONFIG.ENDPOINTS.CATEGORY.DELETE}/${id}`)
  }

  // 切换分类状态（管理员功能）
  static async toggleStatus(id: number, status: number): Promise<ApiResponse<null>> {
    return http.put<null>(`${API_CONFIG.ENDPOINTS.CATEGORY.ADMIN_BASE}/${id}/status`, { status })
  }

  // 获取热门分类
  static async getPopularCategories(limit?: number): Promise<ApiResponse<Category[]>> {
    return http.get<Category[]>('/category/popular', { limit })
  }

  // 获取所有启用的分类（不分页，显式方法名，供其他调用方使用）
  static async getAllCategories(): Promise<ApiResponse<Category[]>> {
    return http.get<Category[]>(API_CONFIG.ENDPOINTS.CATEGORY.ALL)
  }

  // 获取包含文章数量的分类列表（不分页）
  static async getCategoriesWithCount(): Promise<ApiResponse<(Category & { article_count?: number })[]>> {
    return http.get<(Category & { article_count?: number })[]>(API_CONFIG.ENDPOINTS.CATEGORY.WITH_COUNT)
  }
}
