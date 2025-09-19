<template>
  <div class="min-h-screen bg-gray-50">
    <BaseHeader :showCreateButton="true" />

    <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- 页面导航 -->
      <div class="mb-6">
        <nav class="flex items-center space-x-2 text-sm text-gray-600">
          <router-link to="/" class="hover:text-pink-600 transition-colors">首页</router-link>
          <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
          </svg>
          <router-link to="/community" class="hover:text-pink-600 transition-colors">社区</router-link>
          <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
          </svg>
          <span class="text-gray-900">发布帖子</span>
        </nav>
      </div>

      <!-- 发帖表单 -->
      <ForumPostForm
        @success="handleSuccess"
        @cancel="handleCancel"
      />

      <!-- 发帖提示 -->
      <div class="mt-8 bg-blue-50 border border-blue-200 rounded-lg p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-blue-400" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-blue-800">发帖须知</h3>
            <div class="mt-2 text-sm text-blue-700">
              <ul class="list-disc list-inside space-y-1">
                <li>请选择合适的话题分类，方便其他用户查找</li>
                <li>标题要简洁明了，能够准确概括帖子内容</li>
                <li>内容要详细具体，便于其他家长理解和回复</li>
                <li>请保持友善和尊重，营造良好的交流氛围</li>
                <li>避免发布广告、垃圾信息或不当内容</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import BaseHeader from '@/components/BaseHeader.vue'
import ForumPostForm from '@/components/ForumPostForm.vue'
import type { ForumPost } from '@/api/types'

const router = useRouter()
const authStore = useAuthStore()
const { showToast } = useToast()

// 检查登录状态
if (!authStore.isAuthenticated) {
  showToast('请先登录后再发帖', 'warning')
  router.push('/login')
}

// 处理发帖成功
const handleSuccess = (post: ForumPost) => {
  // 跳转到新创建的帖子详情页面
  router.push(`/community/posts/${post.id}`)
}

// 处理取消发帖
const handleCancel = () => {
  router.push('/community')
}
</script>