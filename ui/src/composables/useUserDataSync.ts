// 用户数据同步组合式函数
import { ref, computed, reactive } from 'vue'
import { UserApi } from '@/api/user'
import { FollowApi } from '@/api/follow'
import { ArticleApi } from '@/api/article'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import type { User, Article } from '@/api/types'

// 用户统计信息接口
export interface UserStats {
  following_count: number
  followers_count: number
  article_count: number
  like_count: number
  comment_count: number
  view_count: number
}

// 关注用户信息接口
export interface FollowUser extends User {
  followed_at?: string
  is_following?: boolean
  is_mutual_follow?: boolean
}

// 用户数据同步状态
export interface UserDataState {
  profile: User | null
  stats: UserStats
  articles: Article[]
  followingList: FollowUser[]
  followersList: FollowUser[]
  mutualFollowsList: FollowUser[]
  isLoading: boolean
  error: string | null
}

// 创建默认状态
const createDefaultState = (): UserDataState => ({
  profile: null,
  stats: {
    following_count: 0,
    followers_count: 0,
    article_count: 0,
    like_count: 0,
    comment_count: 0,
    view_count: 0
  },
  articles: [],
  followingList: [],
  followersList: [],
  mutualFollowsList: [],
  isLoading: false,
  error: null
})

// 全局状态存储，用于跨组件数据同步
const globalUserData = reactive<Map<string, UserDataState>>(new Map())

export function useUserDataSync(username?: string) {
  const authStore = useAuthStore()
  const { toast } = useToast()
  
  // 确定用户标识符，优先使用传入的username，否则使用当前登录用户
  const userKey = computed(() => username || authStore.user?.username || '')

  // 获取或创建用户数据状态 - 确保同一个用户返回同一个响应式对象
  const state = computed(() => {
    const key = userKey.value
    if (!key) return reactive(createDefaultState())

    if (!globalUserData.has(key)) {
      globalUserData.set(key, reactive(createDefaultState()))
    }
    return globalUserData.get(key)!
  })
  
  // 是否为当前用户
  const isCurrentUser = computed(() => 
    userKey.value === authStore.user?.username
  )
  
  // 加载用户基本信息
  const loadUserProfile = async (forceRefresh = false) => {
    if (!userKey.value || (state.value.profile && !forceRefresh)) return

    try {
      state.value.isLoading = true
      state.value.error = null

      const response = isCurrentUser.value
        ? await UserApi.getProfile()
        : await UserApi.getUserByUsername(userKey.value)

      state.value.profile = response.data

      // 如果是当前用户，同步到认证store
      if (isCurrentUser.value) {
        authStore.updateUserInfo(response.data)
      }
    } catch (error: any) {
      state.value.error = error.response?.data?.message || '加载用户信息失败'
      console.error('加载用户信息失败:', error)
    } finally {
      state.value.isLoading = false
    }
  }
  
  // 加载用户统计信息
  const loadUserStats = async (forceRefresh = false) => {
    if (!userKey.value) return

    try {
      // 获取用户的关注统计（适用于所有用户）
      if (state.value.profile?.id) {
        try {
          let response
          if (isCurrentUser.value) {
            // 当前用户使用认证接口
            response = await FollowApi.getFollowStats()
          } else {
            // 其他用户使用公开接口
            response = await FollowApi.getUserFollowStats(state.value.profile.id)
          }
          const stats = response.data || response
          state.value.stats.following_count = stats.following_count || 0
          state.value.stats.followers_count = stats.followers_count || 0
          console.log('成功获取用户关注统计:', stats)
        } catch (error) {
          console.error('获取关注统计失败:', error)
        }
      }

      // 根据文章数组设置统计信息
      if (state.value.articles.length > 0) {
        state.value.stats.article_count = state.value.articles.length

        // 计算总点赞数
        state.value.stats.like_count = state.value.articles.reduce((total, article) =>
          total + (article.like_count || 0), 0)

        // 计算总阅读数
        state.value.stats.view_count = state.value.articles.reduce((total, article) =>
          total + (article.view_count || 0), 0)
      }
      
    } catch (error) {
      console.error('加载用户统计失败:', error)
    }
  }
  
  // 加载用户文章
  const loadUserArticles = async (params = { page: 1, size: 20 }) => {
    if (!userKey.value) return

    try {
      // 设置加载状态
      state.value.isLoading = true
      state.value.error = null

      let response

      if (isCurrentUser.value) {
        // 当前用户使用 /articles/my 端点
        response = await ArticleApi.getMyArticles(params)
      } else {
        // 其他用户使用用户名查询
        response = await UserApi.getUserArticlesByUsername(userKey.value, params)
      }

      // 处理后端响应格式：
      // 后端使用 utils.SuccessPage() 返回: {code, message, data: [...articles...], total, page, size}
      // 而不是前端期望的 {code, message, data: {items: [...articles...], total, page, size}}
      // 所以需要直接使用 response.data 作为文章数组
      let articles = []

      if (Array.isArray(response.data)) {
        // 如果 response.data 直接是数组（当前后端返回格式）
        articles = response.data
      } else if (response.data && Array.isArray(response.data.items)) {
        // 如果是标准的 PaginatedResponse 格式
        articles = response.data.items
      } else if (response && Array.isArray(response)) {
        // 直接是数组的情况
        articles = response
      } else {
        // 其他格式，尝试提取数组
        articles = response.data || []
      }

      state.value.articles = articles

      // 更新文章数量统计
      state.value.stats.article_count = articles.length
    } catch (error) {
      state.value.error = error.message || '加载文章失败'
      console.error('加载用户文章失败:', error)
      throw error
    } finally {
      // 重置加载状态
      state.value.isLoading = false
    }
  }
  
  // 加载关注列表
  const loadFollowingList = async (params = { page: 1, limit: 50 }) => {
    if (!isCurrentUser.value) return
    
    try {
      const response = await FollowApi.getFollowing(params)
      // 处理不同的响应格式
      let users = []
      if (response.data?.users) {
        users = response.data.users
      } else if (Array.isArray(response.data)) {
        users = response.data
      } else {
        users = response.users || []
      }
      state.value.followingList = users
      state.value.stats.following_count = users.length
    } catch (error) {
      console.error('加载关注列表失败:', error)
    }
  }
  
  // 加载粉丝列表
  const loadFollowersList = async (params = { page: 1, limit: 50 }) => {
    if (!isCurrentUser.value) return
    
    try {
      const response = await FollowApi.getFollowers(params)
      // 处理不同的响应格式
      let users = []
      if (response.data?.users) {
        users = response.data.users
      } else if (Array.isArray(response.data)) {
        users = response.data
      } else {
        users = response.users || []
      }
      state.value.followersList = users
      state.value.stats.followers_count = users.length
    } catch (error) {
      console.error('加载粉丝列表失败:', error)
    }
  }
  
  // 加载互关列表
  const loadMutualFollowsList = async (params = { page: 1, limit: 50 }) => {
    if (!isCurrentUser.value) return
    
    try {
      const response = await FollowApi.getMutualFollows(params)
      // 处理不同的响应格式
      let users = []
      if (response.data?.users) {
        users = response.data.users
      } else if (Array.isArray(response.data)) {
        users = response.data
      } else {
        users = response.users || []
      }
      state.value.mutualFollowsList = users
    } catch (error) {
      console.error('加载互关列表失败:', error)
    }
  }
  
  // 更新用户信息
  const updateUserProfile = async (updateData: any) => {
    if (!isCurrentUser.value) {
      throw new Error('只能更新自己的信息')
    }
    
    try {
      const response = await UserApi.updateProfile(updateData)
      
      // 更新本地状态
      state.value.profile = response.data
      
      // 同步到认证store
      authStore.updateUserInfo(response.data)
      
      // 通知其他可能监听该用户数据的组件
      await loadUserProfile(true)
      
      toast.success('个人信息更新成功')
      return response.data
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || '更新失败'
      toast.error(errorMessage)
      throw error
    }
  }
  
  // 关注用户
  const followUser = async (userId: number) => {
    if (!isCurrentUser.value) return
    
    try {
      await UserApi.followUser(userId)
      
      // 更新关注统计
      state.value.stats.following_count += 1
      
      // 重新加载关注列表
      await loadFollowingList()
      
      toast.success('关注成功')
    } catch (error: any) {
      toast.error(error.response?.data?.message || '关注失败')
      throw error
    }
  }
  
  // 取消关注用户
  const unfollowUser = async (userId: number) => {
    if (!isCurrentUser.value) return
    
    try {
      await UserApi.unfollowUser(userId)
      
      // 更新关注统计
      state.value.stats.following_count = Math.max(0, state.value.stats.following_count - 1)
      
      // 从关注列表中移除
      state.value.followingList = state.value.followingList.filter(user => user.id !== userId)
      
      // 重新加载关注列表以确保数据一致性
      await loadFollowingList()
      
      toast.success('取消关注成功')
    } catch (error: any) {
      toast.error(error.response?.data?.message || '取消关注失败')
      throw error
    }
  }
  
  // 添加文章后更新统计
  const onArticleAdded = (article: Article) => {
    state.value.articles.unshift(article)
    state.value.stats.article_count += 1
  }
  
  // 删除文章后更新统计
  const onArticleDeleted = (articleId: number) => {
    const index = state.value.articles.findIndex(article => article.id === articleId)
    if (index > -1) {
      state.value.articles.splice(index, 1)
      state.value.stats.article_count = Math.max(0, state.value.stats.article_count - 1)
    }
  }
  
  // 更新文章后刷新数据
  const onArticleUpdated = (updatedArticle: Article) => {
    const index = state.value.articles.findIndex(article => article.id === updatedArticle.id)
    if (index > -1) {
      state.value.articles[index] = updatedArticle
    }
  }
  
  // 初始化数据
  const initUserData = async () => {
    if (!userKey.value) return

    // 首先加载用户基本信息和文章数据
    await Promise.all([
      loadUserProfile(),
      loadUserArticles()
    ])

    // 然后根据文章数据计算统计信息
    await loadUserStats()

    // 如果是当前用户，加载关注相关数据
    if (isCurrentUser.value) {
      await Promise.all([
        loadFollowingList(),
        loadFollowersList(),
        loadMutualFollowsList()
      ])
    }
  }
  
  // 刷新所有数据
  const refreshAllData = async () => {
    await Promise.all([
      loadUserProfile(true),
      loadUserStats(true),
      loadUserArticles()
    ])
    
    if (isCurrentUser.value) {
      await Promise.all([
        loadFollowingList(),
        loadFollowersList(),
        loadMutualFollowsList()
      ])
    }
  }
  
  // 清除用户数据（用于登出等场景）
  const clearUserData = (targetUsername?: string) => {
    const key = targetUsername || userKey.value
    if (key) {
      globalUserData.delete(key)
    }
  }
  
  return {
    // 状态
    profile: computed(() => state.value.profile),
    stats: computed(() => state.value.stats),
    articles: computed(() => state.value.articles),
    followingList: computed(() => state.value.followingList),
    followersList: computed(() => state.value.followersList),
    mutualFollowsList: computed(() => state.value.mutualFollowsList),
    isLoading: computed(() => state.value.isLoading),
    error: computed(() => state.value.error),
    isCurrentUser,
    
    // 方法
    loadUserProfile,
    loadUserStats,
    loadUserArticles,
    loadFollowingList,
    loadFollowersList,
    loadMutualFollowsList,
    updateUserProfile,
    followUser,
    unfollowUser,
    onArticleAdded,
    onArticleDeleted,
    onArticleUpdated,
    initUserData,
    refreshAllData,
    clearUserData
  }
}