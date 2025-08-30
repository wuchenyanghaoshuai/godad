<template>
  <div class="min-h-screen bg-gradient-to-br from-pink-50 via-purple-50 to-orange-50 flex items-center justify-center p-4">
    <div class="max-w-md w-full">
      <!-- 标题区域 -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gradient-to-r from-pink-500 to-purple-600 mb-4">
          <LockIcon class="w-8 h-8 text-white" />
        </div>
        <h1 class="text-3xl font-bold bg-gradient-to-r from-pink-600 to-purple-600 bg-clip-text text-transparent">
          重置密码
        </h1>
        <p class="text-gray-600 mt-2">请设置您的新密码</p>
      </div>

      <!-- 重置密码表单 -->
      <div class="bg-white/80 backdrop-blur-sm rounded-2xl shadow-xl border border-white/20 p-8">
        <form @submit.prevent="handleResetPassword" class="space-y-6">
          <!-- 新密码输入 -->
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
              新密码
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
                minlength="6"
                class="block w-full pl-10 pr-10 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-all duration-200 bg-white/70"
                placeholder="请输入新密码（至少6位）"
                :disabled="isSubmitting"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
              >
                <EyeIcon v-if="showPassword" class="h-5 w-5" />
                <EyeOffIcon v-else class="h-5 w-5" />
              </button>
            </div>
          </div>

          <!-- 确认密码输入 -->
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-2">
              确认新密码
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <LockIcon class="h-5 w-5 text-gray-400" />
              </div>
              <input
                id="confirmPassword"
                v-model="form.confirmPassword"
                :type="showConfirmPassword ? 'text' : 'password'"
                required
                minlength="6"
                class="block w-full pl-10 pr-10 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-all duration-200 bg-white/70"
                placeholder="请再次输入新密码"
                :disabled="isSubmitting"
              />
              <button
                type="button"
                @click="showConfirmPassword = !showConfirmPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
              >
                <EyeIcon v-if="showConfirmPassword" class="h-5 w-5" />
                <EyeOffIcon v-else class="h-5 w-5" />
              </button>
            </div>
            <!-- 密码匹配提示 -->
            <div v-if="form.confirmPassword && form.password !== form.confirmPassword" class="mt-2 text-sm text-red-600">
              两次输入的密码不一致
            </div>
          </div>

          <!-- 密码强度提示 -->
          <div class="text-xs text-gray-500 space-y-1">
            <p>密码要求：</p>
            <ul class="list-disc list-inside space-y-1 ml-2">
              <li :class="form.password.length >= 6 ? 'text-green-600' : 'text-gray-400'">
                至少6个字符
              </li>
              <li :class="/[A-Za-z]/.test(form.password) && /[0-9]/.test(form.password) ? 'text-green-600' : 'text-gray-400'">
                包含字母和数字（推荐）
              </li>
            </ul>
          </div>

          <!-- 提交按钮 -->
          <button
            type="submit"
            :disabled="isSubmitting || !isFormValid"
            class="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-gradient-to-r from-pink-600 to-purple-600 hover:from-pink-700 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-pink-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 transform hover:scale-105"
          >
            <LoaderIcon v-if="isSubmitting" class="w-4 h-4 mr-2 animate-spin" />
            {{ isSubmitting ? '重置中...' : '重置密码' }}
          </button>
        </form>

        <!-- 成功提示 -->
        <div v-if="resetSuccess" class="mt-6 p-4 bg-green-50 border border-green-200 rounded-lg">
          <div class="flex items-center">
            <CheckCircleIcon class="w-5 h-5 text-green-500 mr-2" />
            <div>
              <p class="text-sm text-green-800 font-medium">密码重置成功！</p>
              <p class="text-xs text-green-700 mt-1">请使用新密码登录您的账户</p>
            </div>
          </div>
          <div class="mt-3">
            <button
              @click="goToLogin"
              class="text-sm bg-green-100 hover:bg-green-200 text-green-800 px-3 py-1 rounded transition-colors duration-200"
            >
              立即登录
            </button>
          </div>
        </div>

        <!-- 错误提示 -->
        <div v-if="error" class="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
          <div class="flex items-center">
            <XCircleIcon class="w-5 h-5 text-red-500 mr-2" />
            <p class="text-sm text-red-800">{{ error }}</p>
          </div>
        </div>

        <!-- 返回登录 -->
        <div class="mt-6 text-center">
          <router-link
            to="/login"
            class="text-sm text-pink-600 hover:text-pink-500 font-medium transition-colors duration-200"
          >
            ← 返回登录
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  LockIcon, 
  EyeIcon, 
  EyeOffIcon, 
  LoaderIcon, 
  CheckCircleIcon, 
  XCircleIcon 
} from 'lucide-vue-next'
import { AuthApi } from '@/api/auth'

const route = useRoute()
const router = useRouter()

// 响应式数据
const form = ref({
  password: '',
  confirmPassword: ''
})

const showPassword = ref(false)
const showConfirmPassword = ref(false)
const isSubmitting = ref(false)
const resetSuccess = ref(false)
const error = ref('')

// 重置令牌
const resetToken = ref('')

// 表单验证
const isFormValid = computed(() => {
  return form.value.password.length >= 6 && 
         form.value.password === form.value.confirmPassword
})

// 组件挂载时检查令牌
onMounted(() => {
  const token = route.query.token as string
  
  if (!token) {
    error.value = '无效的重置链接，请重新申请密码重置'
    return
  }
  
  resetToken.value = token
})

// 处理重置密码
const handleResetPassword = async () => {
  if (!isFormValid.value || !resetToken.value) return

  try {
    isSubmitting.value = true
    error.value = ''
    
    await AuthApi.resetPassword(resetToken.value, form.value.password)
    
    resetSuccess.value = true
    
  } catch (err: any) {
    error.value = err.message || '密码重置失败，请稍后重试'
    console.error('重置密码失败:', err)
  } finally {
    isSubmitting.value = false
  }
}

// 跳转到登录页
const goToLogin = () => {
  router.push('/login')
}
</script>