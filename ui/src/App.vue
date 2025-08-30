<template>
  <ErrorBoundary ref="errorBoundary">
    <router-view />
  </ErrorBoundary>
  <Toast />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import ErrorBoundary from '@/components/ErrorBoundary.vue'
import Toast from '@/components/Toast.vue'
import { errorHandler } from '@/utils/errorHandler'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const errorBoundary = ref()
const authStore = useAuthStore()

// 初始化应用
onMounted(async () => {
  // 初始化错误处理系统
  errorHandler.init(router, errorBoundary.value)
  
  // 初始化认证状态
  try {
    await authStore.initAuth()
  } catch (error) {
    console.error('认证初始化失败:', error)
  }
})
</script>