// API数据类型定义

// 用户相关类型
export interface User {
  id: number
  username: string
  email: string
  nickname: string
  avatar?: string
  bio?: string
  role: 'user' | 'admin' | 'content_manager'
  created_at: string
  updated_at: string
}

// 关注相关类型
export interface FollowStats {
  following_count: number
  followers_count: number
}

export interface FollowUser extends User {
  is_mutual_follow?: boolean
  is_following?: boolean
}

export interface FollowListResponse {
  users: FollowUser[]
  total: number
  page: number
  limit: number
  pages: number
}

// 用户注册请求
export interface UserRegisterRequest {
  username: string
  email: string
  password: string
  nickname: string
}

// 用户登录请求
export interface UserLoginRequest {
  username: string
  password: string
}

// 用户登录响应
export interface UserLoginResponse {
  user: User
  token: string
  refresh_token: string
}

// 用户更新请求
export interface UserUpdateRequest {
  nickname?: string
  avatar?: string
  bio?: string
}

// 修改密码请求
export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

// 分类相关类型
export interface Category {
  id: number
  name: string
  description?: string
  created_at: string
  updated_at: string
}

// 分类创建/更新请求
export interface CategoryRequest {
  name: string
  description?: string
}

// 文章相关类型
export interface Article {
  id: number
  title: string
  content: string
  summary?: string
  cover_image?: string
  category_id: number
  category?: Category
  author_id: number
  author?: User
  status: 'draft' | 'published'
  view_count: number
  like_count: number
  comment_count: number
  tags: string[]
  created_at: string
  updated_at: string
}

// 文章创建请求
export interface ArticleCreateRequest {
  title: string
  slug?: string
  content: string
  summary?: string
  cover_image?: string
  category_id: number
  tags?: string // 后端期望逗号分隔的字符串
  status?: number // 后端期望数字：0-草稿，1-已发布
  is_recommend?: boolean // 后端字段名为is_recommend
}

// 文章更新请求（后端接口）
export interface ArticleUpdateRequest {
  title: string // 后端必需，min=1,max=200
  slug: string // 后端必需，min=1,max=200
  content: string // 后端必需，min=1
  summary?: string // 可选，max=500
  cover_image?: string // 可选，max=255
  category_id: number // 后端必需，min=1
  tags?: string // 可选，后端期望逗号分隔的字符串，max=500
  is_top?: boolean
  is_recommend?: boolean
  status?: number // 可选，后端期望数字：0-草稿，1-已发布，max=2
}

// 文章表单数据（前端使用）
export interface ArticleFormData {
  title: string
  content: string
  summary: string
  category_id: number
  cover_image: string
  tags: string[]
  status: 'draft' | 'published'
}

// 文章列表查询参数
export interface ArticleListParams {
  page?: number
  size?: number
  category_id?: number
  author_id?: number
  status?: number // 0-草稿，1-已发布，2-已下架
  keyword?: string
}

// 分页响应
export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  size: number
  total_pages: number
}

// 评论相关类型
export interface Comment {
  id: number
  article_id: number
  user_id: number
  parent_id?: number
  content: string
  likes: number
  is_liked?: boolean
  created_at: string
  updated_at: string
  user: {
    id: number
    username: string
    avatar?: string
    role?: 'user' | 'admin' | 'content_manager'
  }
  replies?: Comment[]
  reply_count?: number
}

// 评论列表查询参数
export interface CommentListParams {
  article_id: number
  page?: number
  size?: number
}

export interface CommentCreateRequest {
  content: string
  article_id: number
  parent_id?: number
  mentions?: number[]
}

export interface CommentUpdateRequest {
  content: string
  mentions?: number[]
}

// 图片上传响应
export interface ImageUploadResponse {
  id: number
  file_name: string
  system_name: string
  file_size: number
  file_type: string
  mime_type: string
  file_hash: string
  storage_path: string
  public_url: string
  user_id: number
  usage: string
  status: number
  created_at: string
  updated_at: string
}

// 通用列表查询参数
export interface ListParams {
  page?: number
  size?: number
  keyword?: string
}

// API错误响应
export interface ApiError {
  code: number
  message: string
  details?: any
}

// API通用响应
export interface ApiResponse<T = any> {
  data: T
  message?: string
  code?: number
}

// ========== 论坛相关类型 ==========

// 论坛帖子
export interface ForumPost {
  id: number
  title: string
  content: string
  topic: string
  author_id: number
  view_count: number
  reply_count: number
  like_count: number
  is_top: boolean
  is_hot: boolean
  is_locked?: boolean
  status: number // 0-草稿 1-已发布 2-已删除
  last_reply_at?: string
  created_at: string
  updated_at: string
  author?: User
  recent_reply?: ForumReply
  time_ago?: string
}

// 论坛回复
export interface ForumReply {
  id: number
  post_id: number
  author_id: number
  parent_id?: number
  content: string
  like_count: number
  status: number // 0-草稿 1-已发布 2-已删除
  created_at: string
  updated_at: string
  author?: User
  parent?: ForumReply
  children?: ForumReply[]
  time_ago?: string
}

// 创建帖子请求
export interface ForumPostCreateRequest {
  title: string
  content: string
  topic: string
}

// 更新帖子请求
export interface ForumPostUpdateRequest {
  title?: string
  content?: string
  topic?: string
  status?: number
}

// 帖子列表查询参数
export interface ForumPostListParams {
  page?: number
  size?: number
  topic?: string
  author_id?: number
  keyword?: string
  sort?: string // 'created_at desc' | 'reply_count desc' | 'view_count desc' | 'last_reply_at desc'
  is_top?: boolean
  is_hot?: boolean
}

// 创建回复请求
export interface ForumReplyCreateRequest {
  post_id: number
  parent_id?: number
  content: string
}

// 回复列表查询参数
export interface ForumReplyListParams {
  page?: number
  size?: number
  post_id: number
  parent_id?: number
  sort?: string // 'created_at asc' | 'created_at desc' | 'like_count desc'
}

// 话题类型
export type ForumTopic =
  | 'Baby Care'
  | 'Feeding'
  | 'Sleep'
  | 'Health'
  | 'Development'
  | 'Activities'
  | 'Gear'
  | 'Parenting'
  | 'Family Life'
  | 'Work & Life Balance'
  | 'Relationships'
  | 'Mental Health'
  | 'Finances'
  | 'Legal'
  | 'Other'

// 话题配置
export interface TopicConfig {
  key: ForumTopic
  label: string
  icon?: string
  color?: string
}

// 论坛统计
export interface ForumStats {
  total_posts: number
  total_replies: number
  active_users: number
  hot_topics: string[]
}
