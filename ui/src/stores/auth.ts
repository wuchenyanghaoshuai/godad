// 用户认证状态管理
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { AuthApi, AuthUtils } from '@/api'
import type { User, UserLoginRequest, UserRegisterRequest } from '@/api/types'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const user = ref<User | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const isAuthenticated = computed(() => !!user.value)
  const isAdmin = computed(() => (user.value as any)?.role === 'admin' || (user.value as any)?.role === 2)
  const isContentManager = computed(() => (user.value as any)?.role === 'content_manager' || isAdmin.value || (user.value as any)?.role === 1)

  // 初始化认证状态
  const initAuth = async () => {
    try {
      // 先从本地快速恢复用户信息
      const storedUser = AuthUtils.getStoredUser()
      if (storedUser) {
        user.value = storedUser
      }

      // 总是尝试验证当前会话状态（Cookie模式）
      // 这样可以处理：页面刷新、浏览器重启等情况
      await getCurrentUserSilently().catch(() => {
        // 如果验证失败，清除本地用户信息
        user.value = null
        AuthUtils.clearAuthData()
      })
    } catch (err) {
      console.error('初始化认证状态失败:', err)
      await logout()
    }
  }

  // 获取当前用户信息
  const getCurrentUser = async () => {
    try {
      isLoading.value = true
      error.value = null

      const response = await AuthApi.getCurrentUser()
      user.value = response.data

      // 保存到本地存储
      AuthUtils.saveUser(response.data)
    } catch (err: any) {
      error.value = err instanceof Error ? err.message : '获取用户信息失败'
      // 只有在非认证错误时才打印日志
      if (err?.code !== 'AUTH_REQUIRED' && err?.status !== 401 && err?.status !== 403) {
        console.error('getCurrentUser failed:', err)
      }
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // 静默获取当前用户信息（用于初始化时检查登录状态）
  const getCurrentUserSilently = async () => {
    try {
      const response = await AuthApi.getCurrentUser(true) // 使用静默模式
      user.value = response.data
      AuthUtils.saveUser(response.data)
    } catch (err: any) {
      // 静默模式下不更新loading状态，不设置error，不打印日志
      throw err
    }
  }

  // 用户登录
  const login = async (credentials: UserLoginRequest) => {
    try {
      isLoading.value = true
      error.value = null

      const response = await AuthApi.login(credentials)
      const { user: userData } = response.data
      
      // 保存认证信息
      user.value = userData
      // Cookie 模式：仅保存用户信息以优化体验
      AuthUtils.saveAuthData(userData)
      
      return userData
    } catch (err) {
      error.value = err instanceof Error ? err.message : '登录失败'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // 用户注册
  const register = async (userData: UserRegisterRequest) => {
    try {
      isLoading.value = true
      error.value = null

      const response = await AuthApi.register(userData)
      const newUser = response.data
      
      // 注册成功，返回用户信息但不自动登录
      return newUser
    } catch (err) {
      error.value = err instanceof Error ? err.message : '注册失败'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // 用户登出
  const logout = async () => {
    try {
      // Cookie 模式：无论本地 token 是否存在都调用登出
      await AuthApi.logout()
    } catch (err) {
      console.error('登出API调用失败:', err)
    } finally {
      // 保存"记住我"的信息
      const rememberMe = localStorage.getItem('remember_me') === 'true'
      const rememberedUsername = localStorage.getItem('remembered_username')
      const rememberedPassword = localStorage.getItem('remembered_password')
      
      // 清除用户数据同步缓存
      const currentUsername = user.value?.username
      
      // 清除本地状态和存储
      user.value = null
      error.value = null
      AuthUtils.clearAuthData()

      // 清除路由守卫的初始化标记，允许重新初始化
      sessionStorage.removeItem('auth_initialized')
      
      // 动态导入并清除用户数据同步缓存
      if (currentUsername) {
        try {
          const { useUserDataSync } = await import('@/composables/useUserDataSync')
          const userDataSync = useUserDataSync()
          userDataSync.clearUserData(currentUsername)
        } catch (err) {
          console.error('清除用户数据缓存失败:', err)
        }
      }
      
      // 如果用户选择了记住我，恢复这些信息
      if (rememberMe && rememberedUsername) {
        localStorage.setItem('remember_me', 'true')
        localStorage.setItem('remembered_username', rememberedUsername)
        if (rememberedPassword) {
          localStorage.setItem('remembered_password', rememberedPassword)
        }
      }
    }
  }

  // Cookie-only：无需在前端管理 refresh token

  // 更新用户信息
  const updateUserInfo = (updatedUser: User) => {
    user.value = updatedUser
    AuthUtils.saveUser(updatedUser)
  }

  // 检查认证状态
  const checkAuth = async (): Promise<boolean> => {
    try {
      // 如果已有用户信息，直接返回
      if (user.value) return true
      // Cookie 模式：直接探测当前会话
      await getCurrentUser()
      return !!user.value
    } catch (err) {
      console.error('认证检查失败:', err)
      await logout()
      return false
    }
  }

  // 清除错误信息
  const clearError = () => {
    error.value = null
  }

  return {
    // 状态
    user,
    isLoading,
    error,
    
    // 计算属性
    isAuthenticated,
    isAdmin,
    isContentManager,
    
    // 方法
    initAuth,
    getCurrentUser,
    getCurrentUserSilently,
    login,
    register,
    logout,
    updateUserInfo,
    checkAuth,
    clearError
  }
})
