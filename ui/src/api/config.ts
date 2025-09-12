// API配置文件
export const API_CONFIG = {
  // 后端API基础URL
  BASE_URL: 'http://127.0.0.1:8888',
  
  // API版本
  API_VERSION: '/api',
  
  // 请求超时时间（毫秒）
  TIMEOUT: 10000,
  
  // 认证相关
  AUTH: {
    TOKEN_KEY: 'godad_token',
    REFRESH_TOKEN_KEY: 'godad_refresh_token',
    USER_INFO_KEY: 'godad_user_info'
  },
  
  // API端点
  ENDPOINTS: {
    // 认证相关
    AUTH: {
      REGISTER: '/auth/register',
      LOGIN: '/auth/login',
      LOGOUT: '/auth/logout',
      REFRESH: '/auth/refresh-token',
      FORGOT_PASSWORD: '/auth/forgot-password',
      RESET_PASSWORD: '/auth/reset-password'
    },
    
    // 用户相关
    USER: {
      PROFILE: '/user/profile',
      UPDATE_PROFILE: '/user/profile',
      CHANGE_PASSWORD: '/user/change-password',
      PUBLIC_INFO: '/user',
      CHECK_NICKNAME: '/user/check-nickname',
      REFRESH_TOKEN: '/user/refresh-token',
      GENERATE_NICKNAME: '/user/generate-nickname'
    },
    
    // 文章相关
    ARTICLE: {
      LIST: '/articles',
      DETAIL: '/articles',
      CREATE: '/articles',
      UPDATE: '/articles',
      DELETE: '/articles',
      PUBLISH: '/articles',
      UNPUBLISH: '/articles',
      MY: '/articles/my'
    },

    // 点赞相关
    LIKE: {
      TOGGLE: '/likes/toggle',
      STATUS: '/likes/status',
      MY: '/likes/my',
      BATCH_STATUS: '/likes/batch-status',
      LIST: '/likes/list',
      POPULAR: '/likes/popular'
    },
    
    // 评论相关
    COMMENT: {
      LIST: '/comments',
      ARTICLE_COMMENTS: '/article-comments',
      CREATE: '/comments',
      UPDATE: '/comments',
      DELETE: '/comments',
      MY: '/comments/my',
      LIKE: '/comments',
      UNLIKE: '/comments',
      REPLIES: '/comments/replies'
    },
    
    // 分类相关
    CATEGORY: {
      LIST: '/categories',
      ALL: '/categories/all',
      WITH_COUNT: '/categories/with-count',
      DETAIL: '/categories',
      BY_SLUG: '/categories/slug',
      CREATE: '/admin/categories',
      UPDATE: '/admin/categories',
      DELETE: '/admin/categories',
      SORT: '/admin/categories',
      TOGGLE: '/admin/categories'
    },
    
    // 收藏相关
    FAVORITE: {
      TOGGLE: '/favorites/toggle',
      STATUS: '/favorites/status',
      BATCH_STATUS: '/favorites/batch-status',
      MY: '/favorites/my',
      DELETE: '/favorites',
      ARTICLE_FAVORITES: '/favorites/article',
      POPULAR: '/favorites/popular'
    },

    // 上传相关
    UPLOAD: {
      IMAGE: '/upload/image'
    }
  }
}

// 构建完整的API URL
export const buildApiUrl = (endpoint: string): string => {
  return `${API_CONFIG.BASE_URL}${API_CONFIG.API_VERSION}${endpoint}`
}