// 用户认证状态管理
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { AuthApi, AuthUtils } from '@/api'
import type { User, UserLoginRequest, UserRegisterRequest } from '@/api/types'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 2)
  const isContentManager = computed(() => user.value?.role === 1 || isAdmin.value)

  // 初始化认证状态
  const initAuth = async () => {
    try {
      // 从本地存储获取token
      const storedToken = AuthUtils.getToken()
      if (storedToken) {
        token.value = storedToken
        
        // 首先从本地存储加载用户信息（快速显示）
        const storedUser = AuthUtils.getStoredUser()
        if (storedUser) {
          user.value = storedUser
          console.log('Loaded user from localStorage:', storedUser)
        }
        
        // 然后从API获取最新用户信息
        await getCurrentUser()
      }
    } catch (err) {
      console.error('初始化认证状态失败:', err)
      // 清除无效的认证信息
      await logout()
    }
  }

  // 获取当前用户信息
  const getCurrentUser = async () => {
    try {
      isLoading.value = true
      error.value = null
      
      const response = await AuthApi.getCurrentUser()
      console.log('API response user data:', response.data)
      user.value = response.data
      
      // 保存到本地存储
      AuthUtils.saveUser(response.data)
      console.log('User saved to localStorage, role:', response.data.role)
    } catch (err) {
      error.value = err instanceof Error ? err.message : '获取用户信息失败'
      console.error('getCurrentUser failed:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // 用户登录
  const login = async (credentials: UserLoginRequest) => {
    try {
      isLoading.value = true
      error.value = null

      const response = await AuthApi.login(credentials)
      const { user: userData, token: userToken } = response.data
      
      // 保存认证信息
      user.value = userData
      token.value = userToken
      AuthUtils.saveAuthData(userData, userToken)
      
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
      // 调用登出API
      if (token.value) {
        await AuthApi.logout()
      }
    } catch (err) {
      console.error('登出API调用失败:', err)
    } finally {
      // 保存"记住我"的信息
      const rememberMe = localStorage.getItem('remember_me') === 'true'
      const rememberedUsername = localStorage.getItem('remembered_username')
      const rememberedPassword = localStorage.getItem('remembered_password')
      
      // 清除本地状态和存储
      user.value = null
      token.value = null
      error.value = null
      AuthUtils.clearAuthData()
      
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

  // 刷新token
  const refreshToken = async () => {
    try {
      const response = await AuthApi.refreshToken()
      const { token: newToken } = response.data
      token.value = newToken
      AuthUtils.saveToken(newToken)
      return newToken
    } catch (err) {
      // 刷新失败，清除认证状态
      await logout()
      throw err
    }
  }

  // 更新用户信息
  const updateUserInfo = (updatedUser: User) => {
    user.value = updatedUser
    AuthUtils.saveUser(updatedUser)
  }

  // 检查认证状态
  const checkAuth = async (): Promise<boolean> => {
    try {
      // 如果已经有用户信息，直接返回true
      if (user.value && token.value) {
        return true
      }
      
      // 尝试从本地存储恢复认证状态
      const storedToken = AuthUtils.getToken()
      if (storedToken) {
        token.value = storedToken
        await getCurrentUser()
        return !!user.value
      }
      
      return false
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
    token,
    isLoading,
    error,
    
    // 计算属性
    isAuthenticated,
    isAdmin,
    isContentManager,
    
    // 方法
    initAuth,
    getCurrentUser,
    login,
    register,
    logout,
    refreshToken,
    updateUserInfo,
    checkAuth,
    clearError
  }
})