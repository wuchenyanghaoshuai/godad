// 分类相关API服务
import { http } from './http'
import { API_CONFIG } from './config'
import type {
  Category,
  CategoryRequest,
  ApiResponse,
  ListParams
} from './types'

// 分类API服务类
export class CategoryApi {
  // 获取分类列表
  static async getCategoryList(params?: ListParams): Promise<ApiResponse<Category[]>> {
    return http.get<Category[]>(API_CONFIG.ENDPOINTS.CATEGORY.LIST, params)
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

  // 获取热门分类
  static async getPopularCategories(limit?: number): Promise<ApiResponse<Category[]>> {
    return http.get<Category[]>('/category/popular', { limit })
  }
}