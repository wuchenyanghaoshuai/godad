<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 flex items-center justify-center px-4">
    <div class="max-w-md w-full text-center">
      <!-- 404 图标 -->  
      <div class="mb-8">
        <div class="w-32 h-32 mx-auto bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center shadow-2xl">
          <span class="text-6xl font-bold text-white">404</span>
        </div>
      </div>
      
      <!-- 错误信息 -->
      <h1 class="text-3xl sm:text-4xl font-bold text-gray-900 mb-4">
        页面未找到
      </h1>
      
      <p class="text-lg text-gray-600 mb-8 leading-relaxed">
        抱歉，您访问的页面不存在或已被移除。
      </p>
      
      <!-- 操作按钮 -->
      <div class="space-y-4">
        <router-link 
          to="/"
          class="inline-flex items-center justify-center w-full px-6 py-3 bg-gradient-to-r from-pink-600 to-orange-600 text-white font-medium rounded-lg hover:from-pink-700 hover:to-orange-700 transition-all duration-300 hover:scale-105 shadow-lg hover:shadow-xl"
        >
          <HomeIcon class="h-5 w-5 mr-2" />
          返回首页
        </router-link>
        
        <button 
          @click="goBack"
          class="inline-flex items-center justify-center w-full px-6 py-3 bg-gray-200 text-gray-700 font-medium rounded-lg hover:bg-gray-300 transition-all duration-300"
        >
          <ArrowLeftIcon class="h-5 w-5 mr-2" />
          返回上一页
        </button>
      </div>
      
      <!-- 建议链接 -->
      <div class="mt-12 pt-8 border-t border-gray-200">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">您可能感兴趣的页面</h3>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <router-link 
            to="/articles"
            class="p-3 bg-white rounded-lg shadow-sm hover:shadow-md transition-all duration-300 text-gray-700 hover:text-pink-600 border border-gray-100 hover:border-pink-200"
          >
            <DocumentTextIcon class="h-5 w-5 mx-auto mb-1" />
            <span class="text-sm font-medium">文章列表</span>
          </router-link>
          
          <router-link 
            to="/categories"
            class="p-3 bg-white rounded-lg shadow-sm hover:shadow-md transition-all duration-300 text-gray-700 hover:text-pink-600 border border-gray-100 hover:border-pink-200"
          >
            <TagIcon class="h-5 w-5 mx-auto mb-1" />
            <span class="text-sm font-medium">分类浏览</span>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { 
  HomeIcon, 
  ArrowLeftIcon, 
  FileTextIcon as DocumentTextIcon, 
  TagIcon 
} from 'lucide-vue-next'

const router = useRouter()

// 返回上一页
const goBack = () => {
  if (window.history.length > 1) {
    router.go(-1)
  } else {
    router.push('/')
  }
}
</script>