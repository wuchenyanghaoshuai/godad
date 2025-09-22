// 论坛相关API服务
import { http } from './http'
import { API_CONFIG } from './config'
import { normalizePageResponse } from './pagination'
import type {
  ForumPost,
  ForumPostCreateRequest,
  ForumPostUpdateRequest,
  ForumPostListParams,
  ForumReply,
  ForumReplyCreateRequest,
  ForumReplyListParams,
  ApiResponse,
  PaginatedResponse
} from './types'

// 论坛API服务类
export class ForumApi {
  // ========== 帖子相关 ==========

  // 获取帖子列表
  static async getPostList(params?: ForumPostListParams): Promise<ApiResponse<PaginatedResponse<ForumPost>>> {
    const resp = await http.get<any>(API_CONFIG.ENDPOINTS.FORUM.POSTS, params)
    const page = normalizePageResponse<ForumPost>(resp)
    return { code: 200, message: 'success', data: page }
  }

  // 获取帖子详情
  static async getPost(id: number): Promise<ApiResponse<ForumPost>> {
    return http.get<ForumPost>(`${API_CONFIG.ENDPOINTS.FORUM.POSTS}/${id}`)
  }

  // 获取帖子详情（别名）
  static async getPostDetail(id: number): Promise<ApiResponse<ForumPost>> {
    return this.getPost(id)
  }

  // 创建帖子
  static async createPost(data: ForumPostCreateRequest): Promise<ApiResponse<ForumPost>> {
    return http.post<ForumPost>(API_CONFIG.ENDPOINTS.FORUM.POSTS, data)
  }

  // 更新帖子
  static async updatePost(id: number, data: ForumPostUpdateRequest): Promise<ApiResponse<ForumPost>> {
    return http.put<ForumPost>(`${API_CONFIG.ENDPOINTS.FORUM.POSTS}/${id}`, data)
  }

  // 删除帖子
  static async deletePost(id: number): Promise<ApiResponse<null>> {
    return http.delete<null>(`${API_CONFIG.ENDPOINTS.FORUM.POSTS}/${id}`)
  }

  // 获取我的帖子列表
  static async getMyPosts(params?: ForumPostListParams): Promise<ApiResponse<PaginatedResponse<ForumPost>>> {
    const resp = await http.get<any>(API_CONFIG.ENDPOINTS.FORUM.MY_POSTS, params)
    const page = normalizePageResponse<ForumPost>(resp)
    return { code: 200, message: 'success', data: page }
  }

  // 增加帖子浏览量
  static async incrementPostView(id: number): Promise<ApiResponse<null>> {
    return http.post<null>(`${API_CONFIG.ENDPOINTS.FORUM.POSTS}/${id}/view`)
  }

  // 点赞帖子 (使用通用点赞API)
  static async togglePostLike(id: number): Promise<ApiResponse<any>> {
    return http.post<any>(API_CONFIG.ENDPOINTS.LIKE.TOGGLE, {
      target_type: 'forum_post',
      target_id: id
    })
  }

  // 获取帖子点赞状态
  static async getPostLikeStatus(id: number): Promise<ApiResponse<any>> {
    return http.get<any>(API_CONFIG.ENDPOINTS.LIKE.STATUS, { target_type: 'forum_post', target_id: id })
  }

  // ========== 回复相关 ==========

  // 获取帖子回复列表
  static async getPostReplies(postId: number, params?: ForumReplyListParams): Promise<ApiResponse<PaginatedResponse<ForumReply>>> {
    const queryParams = { ...params, post_id: postId }
    const resp = await http.get<any>(`${API_CONFIG.ENDPOINTS.FORUM.POSTS}/${postId}/replies`, queryParams)
    const page = normalizePageResponse<ForumReply>(resp)
    return { code: 200, message: 'success', data: page }
  }

  // 创建回复
  static async createReply(data: ForumReplyCreateRequest): Promise<ApiResponse<ForumReply>> {
    return http.post<ForumReply>(API_CONFIG.ENDPOINTS.FORUM.REPLIES, data)
  }

  // 更新回复
  static async updateReply(id: number, content: string): Promise<ApiResponse<ForumReply>> {
    return http.put<ForumReply>(`${API_CONFIG.ENDPOINTS.FORUM.REPLIES}/${id}`, { content })
  }

  // 删除回复
  static async deleteReply(id: number): Promise<ApiResponse<null>> {
    return http.delete<null>(`${API_CONFIG.ENDPOINTS.FORUM.REPLIES}/${id}`)
  }

  // 点赞回复
  static async toggleReplyLike(id: number): Promise<ApiResponse<any>> {
    return http.post<any>(`${API_CONFIG.ENDPOINTS.FORUM.REPLIES}/${id}/like`)
  }

  // ========== 话题相关 ==========

  // 获取话题列表
  static async getTopics(): Promise<ApiResponse<string[]>> {
    return http.get<string[]>(API_CONFIG.ENDPOINTS.FORUM.TOPICS)
  }

  // 获取热门帖子
  static async getHotPosts(params?: { limit?: number }): Promise<ApiResponse<ForumPost[]>> {
    return http.get<ForumPost[]>(API_CONFIG.ENDPOINTS.FORUM.HOT, params)
  }
}
