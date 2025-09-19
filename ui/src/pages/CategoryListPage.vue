<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 通用头部导航 -->
    <BaseHeader />

    <!-- 分类网格 -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <!-- 加载状态 -->
      <div v-if="loading" class="text-center py-12">
        <div class="inline-flex items-center px-4 py-2 font-medium text-pink-600 bg-pink-50 rounded-lg">
          <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-pink-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          加载中...
        </div>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error" class="text-center py-12">
        <div class="bg-red-50 border border-red-200 rounded-lg p-4 inline-block">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">{{ error }}</h3>
            </div>
          </div>
        </div>
      </div>

      <!-- 分类网格 -->
      <div v-if="categories.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <router-link 
          v-for="category in categories" 
          :key="category.id"
          :to="{ path: '/articles', query: { category: category.id } }"
          class="group bg-white rounded-xl shadow-sm hover:shadow-lg transition-all duration-300 overflow-hidden border border-gray-100 hover:border-pink-200 transform hover:scale-105 cursor-pointer"
        >
          <div class="p-6">
            <!-- 分类图标和信息 -->
            <div class="flex items-start mb-4">
              <div class="w-14 h-14 bg-gradient-to-br from-pink-400 to-rose-400 rounded-xl flex items-center justify-center group-hover:scale-105 transition-transform duration-200 shadow-lg">
                <span class="text-white text-xl font-bold">{{ category.name.charAt(0) }}</span>
              </div>
              <div class="ml-4 flex-1">
                <h3 class="text-xl font-semibold text-gray-900 group-hover:text-pink-600 transition-colors mb-1">
                  {{ category.name }}
                </h3>
                <div class="flex items-center text-sm text-gray-500">
                  <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                  </svg>
                  {{ category.article_count || 0 }} 篇文章
                </div>
              </div>
            </div>
            
            <!-- 分类描述 -->
            <p v-if="category.description" class="text-gray-600 text-sm mb-4 line-clamp-2 leading-relaxed">
              {{ category.description }}
            </p>
            
            <!-- 进入箭头 -->
            <div class="flex justify-end">
              <div class="w-8 h-8 bg-gray-100 group-hover:bg-pink-100 rounded-full flex items-center justify-center transition-colors">
                <svg class="w-4 h-4 text-gray-400 group-hover:text-pink-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                </svg>
              </div>
            </div>
          </div>
        </router-link>
      </div>

      <!-- 空状态 -->
      <div v-else class="text-center py-16">
        <div class="text-gray-400 text-6xl mb-4">📚</div>
        <h3 class="text-lg font-medium text-gray-900 mb-2">暂无分类</h3>
        <p class="text-gray-600">目前还没有设置任何分类</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { CategoryApi } from '@/api/category'
import type { Category } from '@/api/types'
import BaseHeader from '@/components/BaseHeader.vue'

// 响应式数据
const loading = ref(false)
const error = ref<string | null>(null)
const categories = ref<Category[]>([])

// 加载分类列表
const loadCategories = async () => {
  try {
    loading.value = true
    error.value = null
    
    const response = await CategoryApi.getCategoriesWithCount()
    categories.value = response.data || []
  } catch (err: any) {
    error.value = err.message || '加载分类失败'
    console.error('加载分类失败:', err)
  } finally {
    loading.value = false
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadCategories()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>