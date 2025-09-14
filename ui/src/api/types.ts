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
}

export interface CommentUpdateRequest {
  content: string
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