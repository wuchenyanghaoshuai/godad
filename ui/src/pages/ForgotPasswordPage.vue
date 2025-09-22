<template>
  <div class="min-h-screen bg-gradient-to-br from-pink-50 via-purple-50 to-orange-50 flex items-center justify-center p-4">
    <div class="max-w-md w-full">
      <!-- 标题区域 -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gradient-to-r from-pink-500 to-purple-600 mb-4">
          <KeyIcon class="w-8 h-8 text-white" />
        </div>
        <h1 class="text-3xl font-bold bg-gradient-to-r from-pink-600 to-purple-600 bg-clip-text text-transparent">
          找回密码
        </h1>
        <p class="text-gray-600 mt-2">输入您的邮箱地址，我们将发送重置链接</p>
      </div>

      <!-- 忘记密码表单 -->
      <div class="bg-white/80 backdrop-blur-sm rounded-2xl shadow-xl border border-white/20 p-8">
        <form @submit.prevent="handleForgotPassword" class="space-y-6">
          <!-- 邮箱输入 -->
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
              邮箱地址
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <MailIcon class="h-5 w-5 text-gray-400" />
              </div>
              <input
                id="email"
                v-model="form.email"
                type="email"
                required
                class="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent transition-all duration-200 bg-white/70"
                placeholder="请输入您的邮箱"
                :disabled="isSubmitting"
              />
            </div>
          </div>

          <!-- 提交按钮 -->
          <button
            type="submit"
            :disabled="isSubmitting || !form.email"
            class="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-gradient-to-r from-pink-600 to-purple-600 hover:from-pink-700 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-pink-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 transform hover:scale-105"
          >
            <LoaderIcon v-if="isSubmitting" class="w-4 h-4 mr-2 animate-spin" />
            {{ isSubmitting ? '发送中...' : '发送重置链接' }}
          </button>
        </form>

        <!-- 成功提示 -->
        <div v-if="emailSent" class="mt-6 p-4 bg-green-50 border border-green-200 rounded-lg">
          <div class="flex items-center">
            <CheckCircleIcon class="w-5 h-5 text-green-500 mr-2" />
            <p class="text-sm text-green-800">
              重置链接已发送到您的邮箱，请查收邮件并按照指引操作
            </p>
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

      <!-- 帮助信息 -->
      <div class="mt-8 text-center">
        <p class="text-xs text-gray-500">
          没有收到邮件？请检查垃圾邮件文件夹，或
          <button 
            @click="handleResendEmail" 
            :disabled="isSubmitting || !canResend"
            class="text-pink-600 hover:text-pink-500 font-medium disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ canResend ? '重新发送' : `${resendCountdown}秒后可重发` }}
          </button>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  KeyIcon, 
  MailIcon, 
  LoaderIcon, 
  CheckCircleIcon, 
  XCircleIcon 
} from 'lucide-vue-next'
import { AuthApi } from '@/api/auth'

useRouter()

// 响应式数据
const form = ref({
  email: ''
})

const isSubmitting = ref(false)
const emailSent = ref(false)
const error = ref('')
const canResend = ref(true)
const resendCountdown = ref(0)

let countdownTimer: number | null = null

// 处理忘记密码
const handleForgotPassword = async () => {
  if (!form.value.email) return

  try {
    isSubmitting.value = true
    error.value = ''
    
    await AuthApi.forgotPassword(form.value.email)
    
    emailSent.value = true
    startResendCountdown()
    
  } catch (err: any) {
    error.value = err.message || '发送重置链接失败，请稍后重试'
    console.error('忘记密码失败:', err)
  } finally {
    isSubmitting.value = false
  }
}

// 重新发送邮件
const handleResendEmail = async () => {
  if (!canResend.value) return
  await handleForgotPassword()
}

// 开始重发倒计时
const startResendCountdown = () => {
  canResend.value = false
  resendCountdown.value = 60
  
  countdownTimer = setInterval(() => {
    resendCountdown.value--
    
    if (resendCountdown.value <= 0) {
      canResend.value = true
      if (countdownTimer) {
        clearInterval(countdownTimer)
        countdownTimer = null
      }
    }
  }, 1000)
}

// 组件卸载时清理定时器
onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
  }
})
</script>
