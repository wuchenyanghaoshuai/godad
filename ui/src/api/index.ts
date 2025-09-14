// API服务统一导出
export { API_CONFIG } from './config'
export { HttpClient, http } from './http'
export * from './types'
export * from './chatTypes'

// 导入API服务类
export { AuthApi, AuthUtils } from './auth'
export { UserApi } from './user'
export { ArticleApi } from './article'
export { CommentApi } from './comment'
export { CategoryApi } from './category'
export { UploadApi, UploadUtils } from './upload'
export { ChatAPI } from './chat'
export { NotificationApi } from './notification'

// 重新导入API类以确保正确引用
import { AuthApi, AuthUtils } from './auth'
import { UserApi } from './user'
import { ArticleApi } from './article'
import { CommentApi } from './comment'
import { CategoryApi } from './category'
import { UploadApi, UploadUtils } from './upload'
import { ChatAPI } from './chat'
import { NotificationApi } from './notification'

// 统一的API服务对象
export const api = {
  auth: AuthApi,
  user: UserApi,
  article: ArticleApi,
  comment: CommentApi,
  category: CategoryApi,
  upload: UploadApi,
  chat: ChatAPI,
  notification: NotificationApi
}

// 统一的工具对象
export const utils = {
  auth: AuthUtils,
  upload: UploadUtils
}