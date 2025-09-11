<template>
  <div class="min-h-screen bg-gradient-to-br from-purple-50 via-blue-50 to-indigo-100 flex items-center justify-center p-4">
    <!-- 背景装饰 -->
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute -top-40 -right-40 w-80 h-80 bg-purple-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-blob"></div>
      <div class="absolute -bottom-40 -left-40 w-80 h-80 bg-blue-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-blob animation-delay-2000"></div>
      <div class="absolute top-40 left-40 w-80 h-80 bg-indigo-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-blob animation-delay-4000"></div>
    </div>

    <!-- 登录卡片 -->
    <div class="relative bg-white rounded-2xl shadow-2xl p-8 sm:p-10 w-full max-w-md border border-gray-100">
      <!-- Logo 和标题 -->
      <div class="text-center mb-8">
        <div class="mx-auto w-16 h-16 bg-gradient-to-br from-purple-600 to-blue-600 rounded-2xl flex items-center justify-center mb-4 shadow-lg">
          <CogIcon class="h-8 w-8 text-white" />
        </div>
        <h1 class="text-3xl font-bold text-gray-900 mb-2">后台管理</h1>
        <p class="text-gray-600">请输入管理员账号和密码</p>
      </div>

      <!-- 登录表单 -->
      <form @submit.prevent="handleLogin" class="space-y-6">
        <!-- 错误提示 -->
        <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm">
          {{ error }}
        </div>

        <!-- 成功提示 -->
        <div v-if="success" class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg text-sm">
          {{ success }}
        </div>

        <!-- 用户名输入 -->
        <div>
          <label for="username" class="block text-sm font-semibold text-gray-700 mb-2">
            管理员账号
          </label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <UserIcon class="h-5 w-5 text-gray-400" />
            </div>
            <input
              id="username"
              v-model="form.username"
              type="text"
              required
              autocomplete="username"
              placeholder="请输入管理员账号"
              class="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 sm:text-sm transition-all duration-200"
              :disabled="isLoading"
            />
          </div>
        </div>

        <!-- 密码输入 -->
        <div>
          <label for="password" class="block text-sm font-semibold text-gray-700 mb-2">
            密码
          </label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <LockIcon class="h-5 w-5 text-gray-400" />
            </div>
            <input
              id="password"
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              required
              autocomplete="current-password"
              placeholder="请输入密码"
              class="block w-full pl-10 pr-12 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500 sm:text-sm transition-all duration-200"
              :disabled="isLoading"
            />
            <button
              type="button"
              @click="showPassword = !showPassword"
              class="absolute inset-y-0 right-0 pr-3 flex items-center"
              :disabled="isLoading"
            >
              <EyeIcon v-if="showPassword" class="h-5 w-5 text-gray-400 hover:text-gray-600" />
              <EyeOffIcon v-else class="h-5 w-5 text-gray-400 hover:text-gray-600" />
            </button>
          </div>
        </div>

        <!-- 记住登录 -->
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <input
              id="remember-me"
              v-model="form.rememberMe"
              type="checkbox"
              class="h-4 w-4 text-purple-600 focus:ring-purple-500 border-gray-300 rounded"
              :disabled="isLoading"
            />
            <label for="remember-me" class="ml-2 block text-sm text-gray-700">
              记住登录状态
            </label>
          </div>
        </div>

        <!-- 登录按钮 -->
        <button
          type="submit"
          :disabled="isLoading || !form.username || !form.password"
          class="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-semibold text-white bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 transform hover:scale-[1.02]"
        >
          <svg
            v-if="isLoading"
            class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            ></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
          {{ isLoading ? '登录中...' : '登录' }}
        </button>
      </form>

      <!-- 返回主站链接 -->
      <div class="mt-8 pt-6 border-t border-gray-200">
        <div class="text-center">
          <router-link
            to="/"
            class="text-sm text-purple-600 hover:text-purple-700 font-medium transition-colors"
          >
            ← 返回主站
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { CogIcon, UserIcon, LockIcon, EyeIcon, EyeOffIcon } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// 响应式数据
const isLoading = ref(false)
const error = ref('')
const success = ref('')
const showPassword = ref(false)

// 表单数据
const form = reactive({
  username: '',
  password: '',
  rememberMe: false
})

// 处理登录
const handleLogin = async () => {
  try {
    isLoading.value = true
    error.value = ''
    success.value = ''

    // 调用登录API
    const result = await authStore.login({
      username: form.username,
      password: form.password,
      rememberMe: form.rememberMe
    })

    // 检查是否是管理员
    if (!authStore.isAdmin) {
      error.value = '权限不足：您不是管理员，无法访问后台管理'
      await authStore.logout()
      return
    }

    success.value = '登录成功，正在跳转...'
    
    // 延迟跳转以显示成功消息
    setTimeout(() => {
      router.push('/management-dashboard')
    }, 1000)

  } catch (err: any) {
    console.error('登录失败:', err)
    error.value = err.message || '登录失败，请检查账号密码'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
@keyframes blob {
  0% {
    transform: translate(0px, 0px) scale(1);
  }
  33% {
    transform: translate(30px, -50px) scale(1.1);
  }
  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }
  100% {
    transform: translate(0px, 0px) scale(1);
  }
}

.animate-blob {
  animation: blob 7s infinite;
}

.animation-delay-2000 {
  animation-delay: 2s;
}

.animation-delay-4000 {
  animation-delay: 4s;
}
</style>