// 论坛相关API服务
import { http } from './http'
import { API_CONFIG } from './config'
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
    return http.get<PaginatedResponse<ForumPost>>('/forum/posts', params)
  }

  // 获取帖子详情
  static async getPost(id: number): Promise<ApiResponse<ForumPost>> {
    return http.get<ForumPost>(`/forum/posts/${id}`)
  }

  // 获取帖子详情（别名）
  static async getPostDetail(id: number): Promise<ApiResponse<ForumPost>> {
    return this.getPost(id)
  }

  // 创建帖子
  static async createPost(data: ForumPostCreateRequest): Promise<ApiResponse<ForumPost>> {
    return http.post<ForumPost>('/forum/posts', data)
  }

  // 更新帖子
  static async updatePost(id: number, data: ForumPostUpdateRequest): Promise<ApiResponse<ForumPost>> {
    return http.put<ForumPost>(`/forum/posts/${id}`, data)
  }

  // 删除帖子
  static async deletePost(id: number): Promise<ApiResponse<null>> {
    return http.delete<null>(`/forum/posts/${id}`)
  }

  // 获取我的帖子列表
  static async getMyPosts(params?: ForumPostListParams): Promise<ApiResponse<PaginatedResponse<ForumPost>>> {
    return http.get<PaginatedResponse<ForumPost>>('/forum/posts/my', params)
  }

  // 增加帖子浏览量
  static async incrementPostView(id: number): Promise<ApiResponse<null>> {
    return http.post<null>(`/forum/posts/${id}/view`)
  }

  // 点赞帖子 (使用通用点赞API)
  static async togglePostLike(id: number): Promise<ApiResponse<any>> {
    return http.post<any>('/likes/toggle', {
      target_type: 'forum_post',
      target_id: id
    })
  }

  // 获取帖子点赞状态
  static async getPostLikeStatus(id: number): Promise<ApiResponse<any>> {
    return http.get<any>(`/likes/status?target_type=forum_post&target_id=${id}`)
  }

  // ========== 回复相关 ==========

  // 获取帖子回复列表
  static async getPostReplies(postId: number, params?: ForumReplyListParams): Promise<ApiResponse<PaginatedResponse<ForumReply>>> {
    const queryParams = { ...params, post_id: postId }
    return http.get<PaginatedResponse<ForumReply>>(`/forum/posts/${postId}/replies`, queryParams)
  }

  // 创建回复
  static async createReply(data: ForumReplyCreateRequest): Promise<ApiResponse<ForumReply>> {
    return http.post<ForumReply>('/forum/replies', data)
  }

  // 更新回复
  static async updateReply(id: number, content: string): Promise<ApiResponse<ForumReply>> {
    return http.put<ForumReply>(`/forum/replies/${id}`, { content })
  }

  // 删除回复
  static async deleteReply(id: number): Promise<ApiResponse<null>> {
    return http.delete<null>(`/forum/replies/${id}`)
  }

  // 点赞回复
  static async toggleReplyLike(id: number): Promise<ApiResponse<any>> {
    return http.post<any>(`/forum/replies/${id}/like`)
  }

  // ========== 话题相关 ==========

  // 获取话题列表
  static async getTopics(): Promise<ApiResponse<string[]>> {
    return http.get<string[]>('/forum/topics')
  }

  // 获取热门帖子
  static async getHotPosts(params?: { limit?: number }): Promise<ApiResponse<ForumPost[]>> {
    return http.get<ForumPost[]>('/forum/posts/hot', params)
  }
}