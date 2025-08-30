// 文章相关API服务
import { http } from './http'
import { API_CONFIG } from './config'
import type {
  Article,
  ArticleCreateRequest,
  ArticleUpdateRequest,
  ArticleListParams,
  ApiResponse,
  PaginatedResponse
} from './types'

// 文章API服务类
export class ArticleApi {
  // 获取文章列表
  static async getArticleList(params?: ArticleListParams): Promise<ApiResponse<PaginatedResponse<Article>>> {
    return http.get<PaginatedResponse<Article>>(API_CONFIG.ENDPOINTS.ARTICLE.LIST, params)
  }

  // 获取文章详情
  static async getArticleDetail(id: number): Promise<ApiResponse<Article>> {
    return http.get<Article>(`${API_CONFIG.ENDPOINTS.ARTICLE.DETAIL}/${id}`)
  }

  // 创建文章
  static async createArticle(data: ArticleCreateRequest): Promise<ApiResponse<Article>> {
    return http.post<Article>(API_CONFIG.ENDPOINTS.ARTICLE.CREATE, data)
  }

  // 更新文章
  static async updateArticle(id: number, data: ArticleUpdateRequest): Promise<ApiResponse<Article>> {
    return http.put<Article>(`${API_CONFIG.ENDPOINTS.ARTICLE.UPDATE}/${id}`, data)
  }

  // 删除文章
  static async deleteArticle(id: number): Promise<ApiResponse<null>> {
    return http.delete<null>(`${API_CONFIG.ENDPOINTS.ARTICLE.DELETE}/${id}`)
  }

  // 发布文章
  static async publishArticle(id: number): Promise<ApiResponse<Article>> {
    return http.put<Article>(`${API_CONFIG.ENDPOINTS.ARTICLE.PUBLISH}/${id}`)
  }

  // 取消发布文章
  static async unpublishArticle(id: number): Promise<ApiResponse<Article>> {
    return http.put<Article>(`${API_CONFIG.ENDPOINTS.ARTICLE.UNPUBLISH}/${id}`)
  }

  // 获取我的文章列表
  static async getMyArticles(params?: ArticleListParams): Promise<ApiResponse<PaginatedResponse<Article>>> {
    return http.get<PaginatedResponse<Article>>(API_CONFIG.ENDPOINTS.ARTICLE.MY, params)
  }

  // 增加文章浏览量
  static async incrementViewCount(id: number): Promise<ApiResponse<null>> {
    return http.post<null>(`${API_CONFIG.ENDPOINTS.ARTICLE.DETAIL}/${id}/view`)
  }

  // 点赞文章
  static async likeArticle(id: number): Promise<ApiResponse<null>> {
    return http.post<null>(`${API_CONFIG.ENDPOINTS.ARTICLE.LIKE}/${id}/like`)
  }

  // 取消点赞文章
  static async unlikeArticle(id: number): Promise<ApiResponse<null>> {
    return http.delete<null>(`${API_CONFIG.ENDPOINTS.ARTICLE.LIKE}/${id}/like`)
  }

  // 获取相关文章推荐
  static async getRelatedArticles(id: number, limit: number = 5): Promise<ApiResponse<Article[]>> {
    return http.get<Article[]>(`${API_CONFIG.ENDPOINTS.ARTICLE.LIST}/${id}/related`, { limit })
  }
}