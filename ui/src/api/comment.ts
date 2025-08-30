import { http } from './http'
import type { Comment, CommentCreateRequest, CommentUpdateRequest, PaginatedResponse, ApiResponse } from './types'

/**
 * 评论相关API接口
 */
export const CommentApi = {
  /**
   * 获取文章评论列表
   * @param articleId 文章ID
   * @param params 查询参数
   */
  getComments: (articleId: number, params?: {
    page?: number
    limit?: number
    sort?: 'newest' | 'oldest' | 'most_liked'
    parent_id?: number
  }): Promise<ApiResponse<PaginatedResponse<Comment>>> => {
    return http.get(`/article-comments/${articleId}`, { params })
  },

  /**
   * 获取评论详情
   * @param commentId 评论ID
   */
  getComment: (commentId: number): Promise<ApiResponse<Comment>> => {
    return http.get(`/comments/${commentId}`)
  },

  /**
   * 创建评论
   * @param data 评论数据（包含article_id）
   */
  createComment: (data: CommentCreateRequest): Promise<ApiResponse<Comment>> => {
    return http.post(`/comments`, data)
  },

  /**
   * 更新评论
   * @param commentId 评论ID
   * @param data 更新数据
   */
  updateComment: (commentId: number, data: CommentUpdateRequest): Promise<ApiResponse<Comment>> => {
    return http.put(`/comments/${commentId}`, data)
  },

  /**
   * 删除评论
   * @param commentId 评论ID
   */
  deleteComment: (commentId: number): Promise<ApiResponse<void>> => {
    return http.delete(`/comments/${commentId}`)
  },

  /**
   * 点赞评论
   * @param commentId 评论ID
   */
  likeComment: (commentId: number): Promise<ApiResponse<void>> => {
    return http.post(`/comments/${commentId}/like`)
  },

  /**
   * 取消点赞评论
   * @param commentId 评论ID
   */
  unlikeComment: (commentId: number): Promise<ApiResponse<void>> => {
    return http.post(`/comments/${commentId}/unlike`)
  },

  /**
   * 获取评论的回复列表
   * @param parentId 父评论ID
   * @param params 查询参数
   */
  getReplies: (parentId: number, params?: {
    page?: number
    limit?: number
  }): Promise<ApiResponse<PaginatedResponse<Comment>>> => {
    return http.get(`/comments/replies/${parentId}`, { params })
  },

  /**
   * 获取当前用户评论列表
   */
  getMyComments: (params?: {
    page?: number
    limit?: number
  }): Promise<ApiResponse<PaginatedResponse<Comment>>> => {
    return http.get('/comments/my', { params })
  }
}