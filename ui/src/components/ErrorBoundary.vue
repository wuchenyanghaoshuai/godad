<template>
  <div v-if="hasError" class="min-h-screen bg-gray-50 flex items-center justify-center px-4">
    <div class="max-w-md w-full text-center">
      <!-- 错误图标 -->
      <div class="mb-8">
        <div class="w-24 h-24 mx-auto bg-red-100 rounded-full flex items-center justify-center">
          <ExclamationTriangleIcon class="h-12 w-12 text-red-600" />
        </div>
      </div>
      
      <!-- 错误信息 -->
      <h1 class="text-2xl font-bold text-gray-900 mb-4">
        出现了一些问题
      </h1>
      
      <p class="text-gray-600 mb-6">
        {{ errorMessage || '应用遇到了意外错误，请稍后重试。' }}
      </p>
      
      <!-- 操作按钮 -->
      <div class="space-y-3">
        <button 
          @click="retry"
          class="w-full px-6 py-3 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 transition-colors duration-200"
        >
          重试
        </button>
        
        <button 
          @click="goHome"
          class="w-full px-6 py-3 bg-gray-200 text-gray-700 font-medium rounded-lg hover:bg-gray-300 transition-colors duration-200"
        >
          返回首页
        </button>
      </div>
      
      <!-- 错误详情（开发环境） -->
      <details v-if="isDev && errorDetails" class="mt-8 text-left">
        <summary class="cursor-pointer text-sm text-gray-500 hover:text-gray-700">
          查看错误详情
        </summary>
        <pre class="mt-2 p-4 bg-gray-100 rounded-lg text-xs text-gray-800 overflow-auto max-h-40">{{ errorDetails }}</pre>
      </details>
    </div>
  </div>
  
  <slot v-else />
</template>

<script setup lang="ts">
import { ref, onErrorCaptured } from 'vue'
import { useRouter } from 'vue-router'
import { ExclamationTriangleIcon } from '@heroicons/vue/24/outline'

const router = useRouter()

// 错误状态
const hasError = ref(false)
const errorMessage = ref('')
const errorDetails = ref('')
const isDev = import.meta.env.DEV

// 捕获子组件错误
onErrorCaptured((error: Error, instance, info) => {
  console.error('组件错误:', error, info)
  
  hasError.value = true
  errorMessage.value = error.message
  errorDetails.value = `${error.stack}\n\n组件信息: ${info}`
  
  // 阻止错误继续向上传播
  return false
})

// 重试
const retry = () => {
  hasError.value = false
  errorMessage.value = ''
  errorDetails.value = ''
  
  // 刷新当前页面
  window.location.reload()
}

// 返回首页
const goHome = () => {
  hasError.value = false
  errorMessage.value = ''
  errorDetails.value = ''
  
  router.push('/')
}

// 暴露方法供外部调用
defineExpose({
  showError: (message: string, details?: string) => {
    hasError.value = true
    errorMessage.value = message
    errorDetails.value = details || ''
  },
  clearError: () => {
    hasError.value = false
    errorMessage.value = ''
    errorDetails.value = ''
  }
})
</script>