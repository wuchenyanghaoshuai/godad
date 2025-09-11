<template>
  <div class="min-h-screen bg-gradient-to-br from-pink-50 to-orange-50 flex items-center justify-center px-4">
    <div class="max-w-md w-full space-y-8">
      <!-- 头部 -->
      <div class="text-center">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">欢迎回来</h1>
        <p class="text-gray-600">登录您的GoDad账户</p>
      </div>

      <!-- 登录表单 -->
      <div class="bg-white rounded-2xl shadow-xl p-8">
        <form @submit.prevent="handleLogin" class="space-y-6">
          <!-- 用户名输入 -->
          <div>
            <label for="username" class="block text-sm font-medium text-gray-700 mb-2">
              用户名
            </label>
            <input
              id="username"
              v-model="form.username"
              type="text"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-colors"
              placeholder="请输入您的用户名"
              :disabled="isLoading"
            />
          </div>

          <!-- 密码输入 -->
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
              密码
            </label>
            <div class="relative">
              <input
                id="password"
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                required
                autocomplete="current-password"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-colors pr-12"
                placeholder="请输入您的密码"
                :disabled="isLoading"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
                :disabled="isLoading"
              >
                <EyeIcon v-if="!showPassword" class="h-5 w-5" />
                <EyeOffIcon v-else class="h-5 w-5" />
              </button>
            </div>
          </div>

          <!-- 记住我 -->
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <input
                id="remember"
                v-model="form.remember"
                type="checkbox"
                class="h-4 w-4 text-pink-600 focus:ring-pink-500 border-gray-300 rounded"
                @change="handleRememberChange"
              />
              <label for="remember" class="ml-2 block text-sm text-gray-700">
                记住我
              </label>
            </div>
            <router-link
              to="/forgot-password"
              class="text-sm text-pink-600 hover:text-pink-500"
            >
              忘记密码？
            </router-link>
          </div>

          <!-- 成功信息 -->
          <div v-if="successMessage" class="bg-green-50 border border-green-200 rounded-lg p-3">
            <p class="text-sm text-green-600">{{ successMessage }}</p>
          </div>

          <!-- 错误信息 -->
          <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-3">
            <p class="text-sm text-red-600">{{ error }}</p>
          </div>

          <!-- 登录按钮 -->
          <button
            type="submit"
            :disabled="isLoading"
            class="w-full bg-gradient-to-r from-pink-500 to-orange-500 text-white py-3 px-4 rounded-lg font-medium hover:from-pink-600 hover:to-orange-600 focus:outline-none focus:ring-2 focus:ring-pink-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
          >
            <span v-if="isLoading" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              登录中...
            </span>
            <span v-else>登录</span>
          </button>
        </form>

        <!-- 忘记密码链接 -->
        <div class="mt-4 text-center">
          <router-link
            to="/forgot-password"
            class="text-sm text-pink-600 hover:text-pink-500 font-medium"
          >
            忘记密码？
          </router-link>
        </div>

        <!-- 注册链接 -->
        <div class="mt-6 text-center">
          <p class="text-sm text-gray-600">
            还没有账户？
            <router-link
              to="/register"
              class="text-pink-600 hover:text-pink-500 font-medium"
            >
              立即注册
            </router-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { EyeIcon, EyeOffIcon } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import type { UserLoginRequest } from '@/api/types'

// 路由
const router = useRouter()
const route = useRoute()

// 认证store
const authStore = useAuthStore()

// 响应式数据
const showPassword = ref(false)
const isLoading = ref(false)
const error = ref<string | null>(null)
const successMessage = ref<string | null>(null)

// 表单数据
const form = reactive<UserLoginRequest & { remember: boolean }>({
  username: '',
  password: '',
  remember: false
})
// 简单的加密函数（仅用于演示，生产环境应使用更强的加密）
const encryptPassword = (password: string): string => {
  return btoa(password) // Base64编码，实际应用中应使用更安全的加密
}

const decryptPassword = (encrypted: string): string => {
  try {
    return atob(encrypted) // Base64解码
  } catch {
    return ''
  }
}

// 处理登录
const handleLogin = async () => {
  try {
    isLoading.value = true
    error.value = null
    successMessage.value = null

    // 调用登录API
    await authStore.login({
      username: form.username,
      password: form.password
    })

    // 如果选择了记住我，保存用户名和密码
    if (form.remember) {
      localStorage.setItem('remembered_username', form.username)
      localStorage.setItem('remembered_password', encryptPassword(form.password))
      localStorage.setItem('remember_me', 'true')
    } else {
      localStorage.removeItem('remembered_username')
      localStorage.removeItem('remembered_password')
      localStorage.removeItem('remember_me')
    }

    // 登录成功，跳转到首页
    router.push('/')
  } catch (err) {
    error.value = err instanceof Error ? err.message : '登录失败，请重试'
  } finally {
    isLoading.value = false
  }
}

// 处理记住我选项变化
const handleRememberChange = () => {
  if (!form.remember) {
    // 如果取消勾选记住我，清除保存的信息
    localStorage.removeItem('remember_me')
    localStorage.removeItem('remembered_username')
    localStorage.removeItem('remembered_password')
  }
}

// 组件挂载时检查URL参数中的消息
onMounted(() => {
  const message = route.query.message as string
  if (message) {
    successMessage.value = message
    // 3秒后清除消息
    setTimeout(() => {
      successMessage.value = null
    }, 5000)
  }
  
  // 恢复记住的用户名和密码
  const rememberMe = localStorage.getItem('remember_me') === 'true'
  if (rememberMe) {
    const rememberedUsername = localStorage.getItem('remembered_username')
    const rememberedPassword = localStorage.getItem('remembered_password')
    
    if (rememberedUsername) {
      form.username = rememberedUsername
    }
    
    if (rememberedPassword) {
      form.password = decryptPassword(rememberedPassword)
    }
    
    form.remember = true
  }
  
  // 确保密码输入框类型正确
  setTimeout(() => {
    const passwordInput = document.getElementById('password') as HTMLInputElement
    if (passwordInput && !showPassword.value) {
      passwordInput.type = 'password'
    }
  }, 100)
})
</script>