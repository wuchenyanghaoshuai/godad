<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { MenuIcon, XIcon, UserIcon, LogOutIcon, ArrowRightIcon, UsersIcon, HeartIcon, TrendingUpIcon, ChevronDownIcon, SettingsIcon, PlusIcon, MailIcon, PhoneIcon, MapPinIcon, CogIcon } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

// 路由
const router = useRouter()

// 认证store
const authStore = useAuthStore()

// 响应式数据
const showUserMenu = ref(false)
const showMobileMenu = ref(false)
const userMenuButton = ref<HTMLElement>()

// 退出登录
const handleLogout = async () => {
  try {
    await authStore.logout()
    showUserMenu.value = false
    router.push('/login')
  } catch (error) {
    console.error('退出登录失败:', error)
  }
}

// 点击外部关闭用户菜单
const handleClickOutside = (event: Event) => {
  if (userMenuButton.value && !userMenuButton.value.contains(event.target as Node)) {
    showUserMenu.value = false
  }
  // 点击外部关闭移动端菜单
  if (showMobileMenu.value && !(event.target as Element).closest('nav')) {
    showMobileMenu.value = false
  }
}

// 组件挂载时
onMounted(() => {
  // 初始化认证状态
  authStore.initAuth()
  
  // 添加点击外部事件监听
  document.addEventListener('click', handleClickOutside)
})

// 组件卸载时
onUnmounted(() => {
  // 移除事件监听
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-pink-50 to-orange-50">
    <!-- 导航栏 -->
    <nav class="bg-white/95 backdrop-blur-sm shadow-sm sticky top-0 z-50 border-b border-gray-100">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo -->
          <div class="flex items-center">
            <router-link to="/" class="flex items-center space-x-2 group">
              <div class="w-8 h-8 bg-gradient-to-r from-pink-500 to-orange-500 rounded-lg flex items-center justify-center group-hover:scale-105 transition-transform">
                <span class="text-white font-bold text-lg">G</span>
              </div>
              <span class="text-xl font-bold bg-gradient-to-r from-pink-600 to-orange-600 bg-clip-text text-transparent">GoDad</span>
            </router-link>
          </div>

          <!-- 导航链接 -->
          <div class="hidden md:flex items-center space-x-8">
            <router-link 
              to="/" 
              class="text-gray-700 hover:text-pink-600 transition-all duration-200 font-medium relative group"
            >
              首页
              <span class="absolute -bottom-1 left-0 w-0 h-0.5 bg-pink-600 transition-all duration-200 group-hover:w-full"></span>
            </router-link>
            <router-link 
              to="/articles" 
              class="text-gray-700 hover:text-pink-600 transition-all duration-200 font-medium relative group"
            >
              文章
              <span class="absolute -bottom-1 left-0 w-0 h-0.5 bg-pink-600 transition-all duration-200 group-hover:w-full"></span>
            </router-link>
            <router-link 
              to="/categories" 
              class="text-gray-700 hover:text-pink-600 transition-all duration-200 font-medium relative group"
            >
              分类
              <span class="absolute -bottom-1 left-0 w-0 h-0.5 bg-pink-600 transition-all duration-200 group-hover:w-full"></span>
            </router-link>
            <router-link 
              to="/about" 
              class="text-gray-700 hover:text-pink-600 transition-all duration-200 font-medium relative group"
            >
              关于
              <span class="absolute -bottom-1 left-0 w-0 h-0.5 bg-pink-600 transition-all duration-200 group-hover:w-full"></span>
            </router-link>
          </div>

          <!-- 用户操作区 -->
          <div class="flex items-center space-x-4">
            <!-- 未登录状态 -->
            <div v-if="!authStore.isAuthenticated" class="hidden md:flex items-center space-x-3">
              <router-link
                to="/login"
                class="text-gray-700 hover:text-pink-600 transition-all duration-200 font-medium"
              >
                登录
              </router-link>
              <router-link
                to="/register"
                class="bg-gradient-to-r from-pink-600 to-orange-600 text-white px-6 py-2 rounded-full hover:shadow-lg hover:scale-105 transition-all duration-200 font-medium"
              >
                注册
              </router-link>
            </div>

            <!-- 已登录状态 -->
            <div v-else class="relative">
              <button
                @click="showUserMenu = !showUserMenu"
                class="flex items-center space-x-2 text-gray-700 hover:text-pink-600 transition-all duration-200 p-2 rounded-lg hover:bg-gray-50"
                ref="userMenuButton"
              >
                <div class="w-8 h-8 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center shadow-md">
                  <UserIcon class="h-5 w-5 text-white" />
                </div>
                <span class="hidden md:block font-medium">{{ authStore.user?.username }}</span>
                <ChevronDownIcon class="h-4 w-4 transition-transform duration-200" :class="{ 'rotate-180': showUserMenu }" />
              </button>

              <!-- 用户菜单 -->
              <div
                v-if="showUserMenu"
                class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-xl border border-gray-100 py-2 z-50 animate-in slide-in-from-top-2 duration-200"
              >
                <router-link
                  to="/user-center"
                  class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
                  @click="showUserMenu = false"
                >
                  <UserIcon class="h-4 w-4 mr-3" />
                  个人中心
                </router-link>
                <router-link
                  to="/articles/create"
                  class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
                  @click="showUserMenu = false"
                >
                  <PlusIcon class="h-4 w-4 mr-3" />
                  发布文章
                </router-link>
                <hr class="my-2 border-gray-100" />
                <button
                  @click="handleLogout"
                  class="flex items-center w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors"
                >
                  <LogOutIcon class="h-4 w-4 mr-3" />
                  退出登录
                </button>
              </div>
            </div>

            <!-- 移动端菜单按钮 -->
            <button 
              @click="showMobileMenu = !showMobileMenu"
              class="md:hidden p-2 text-gray-700 hover:text-pink-600 hover:bg-gray-50 rounded-lg transition-all duration-200"
            >
              <MenuIcon v-if="!showMobileMenu" class="h-6 w-6" />
              <XIcon v-else class="h-6 w-6" />
            </button>
          </div>
        </div>
      </div>
    </nav>

    <!-- 主要内容 -->
    <main>
      <!-- 英雄区域 -->
      <section class="relative min-h-screen flex items-center justify-center overflow-hidden bg-gradient-to-br from-pink-50 via-orange-50 to-yellow-50">
        <!-- 装饰性背景元素 -->
        <div class="absolute inset-0 overflow-hidden">
          <div class="absolute -top-20 sm:-top-40 -right-20 sm:-right-40 w-40 sm:w-80 h-40 sm:h-80 bg-gradient-to-br from-pink-400 to-orange-400 rounded-full opacity-20 animate-pulse"></div>
          <div class="absolute -bottom-20 sm:-bottom-40 -left-20 sm:-left-40 w-48 sm:w-96 h-48 sm:h-96 bg-gradient-to-tr from-orange-400 to-yellow-400 rounded-full opacity-20 animate-pulse delay-1000"></div>
          <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-32 sm:w-64 h-32 sm:h-64 bg-gradient-to-r from-pink-300 to-orange-300 rounded-full opacity-10 animate-spin" style="animation-duration: 20s;"></div>
        </div>
        
        <div class="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h1 class="text-3xl sm:text-5xl lg:text-7xl font-bold text-gray-900 mb-4 sm:mb-6 leading-tight">
            <span class="bg-gradient-to-r from-pink-600 to-orange-600 bg-clip-text text-transparent">GoDad</span>
            <br class="hidden sm:block" />
            <span class="text-gray-800 text-2xl sm:text-4xl lg:text-6xl">育儿知识分享平台</span>
          </h1>
          
          <p class="text-base sm:text-lg lg:text-2xl text-gray-600 mb-6 sm:mb-8 max-w-3xl mx-auto leading-relaxed px-2">
            专业的育儿指导，温暖的社区支持，陪伴每一个家庭的成长之路
          </p>
          <div class="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center items-center px-4 sm:px-0">
            <router-link
              to="/articles"
              class="w-full sm:w-auto inline-flex items-center justify-center bg-gradient-to-r from-pink-600 to-orange-600 text-white px-6 sm:px-8 py-3 sm:py-4 rounded-full text-base sm:text-lg font-semibold hover:shadow-xl transform hover:scale-105 transition-all duration-300 shadow-lg"
            >
              开始阅读
              <ArrowRightIcon class="ml-2 h-4 w-4 sm:h-5 sm:w-5" />
            </router-link>
            <router-link
              to="/register"
              class="w-full sm:w-auto inline-flex items-center justify-center border-2 border-pink-600 text-pink-600 px-6 sm:px-8 py-3 sm:py-4 rounded-full text-base sm:text-lg font-semibold hover:bg-pink-50 hover:shadow-lg transition-all duration-300 bg-white/80 backdrop-blur-sm"
            >
              加入社区
              <UsersIcon class="ml-2 h-4 w-4 sm:h-5 sm:w-5" />
            </router-link>
          </div>
          
          <!-- 统计数据预览 -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-3 sm:gap-4 md:gap-6 mt-12 sm:mt-16 px-4 max-w-4xl mx-auto">
            <div class="text-center bg-white/80 backdrop-blur-sm rounded-lg sm:rounded-xl p-3 sm:p-4 md:p-6 border border-white/50 hover:bg-white/90 transition-all duration-300 shadow-lg">
              <div class="text-xl sm:text-2xl md:text-3xl font-bold text-pink-600 mb-1 sm:mb-2">10K+</div>
              <div class="text-xs sm:text-sm md:text-base text-gray-600 font-medium">注册用户</div>
            </div>
            <div class="text-center bg-white/80 backdrop-blur-sm rounded-lg sm:rounded-xl p-3 sm:p-4 md:p-6 border border-white/50 hover:bg-white/90 transition-all duration-300 shadow-lg">
              <div class="text-xl sm:text-2xl md:text-3xl font-bold text-orange-600 mb-1 sm:mb-2">500+</div>
              <div class="text-xs sm:text-sm md:text-base text-gray-600 font-medium">优质文章</div>
            </div>
            <div class="text-center bg-white/80 backdrop-blur-sm rounded-lg sm:rounded-xl p-3 sm:p-4 md:p-6 border border-white/50 hover:bg-white/90 transition-all duration-300 shadow-lg">
              <div class="text-xl sm:text-2xl md:text-3xl font-bold text-purple-600 mb-1 sm:mb-2">50+</div>
              <div class="text-xs sm:text-sm md:text-base text-gray-600 font-medium">专业作者</div>
            </div>
            <div class="text-center bg-white/80 backdrop-blur-sm rounded-lg sm:rounded-xl p-3 sm:p-4 md:p-6 border border-white/50 hover:bg-white/90 transition-all duration-300 shadow-lg">
              <div class="text-xl sm:text-2xl md:text-3xl font-bold text-green-600 mb-1 sm:mb-2">24/7</div>
              <div class="text-xs sm:text-sm md:text-base text-gray-600 font-medium">在线支持</div>
            </div>
          </div>
        </div>
      </section>

      <!-- 特色功能 -->
      <section class="py-12 sm:py-16 md:py-20 bg-white">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="text-center mb-8 sm:mb-12 md:mb-16">
            <h2 class="text-2xl sm:text-3xl md:text-4xl font-bold text-gray-900 mb-3 sm:mb-4">
              为什么选择 <span class="bg-gradient-to-r from-pink-600 to-orange-600 bg-clip-text text-transparent">GoDad</span>
            </h2>
            <p class="text-base sm:text-lg md:text-xl text-gray-600 max-w-3xl mx-auto px-4">
              我们致力于为每一位父母提供最专业、最贴心的育儿支持
            </p>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6 sm:gap-8">
            <!-- 专业可靠 -->
            <div class="text-center p-6 sm:p-8 rounded-xl sm:rounded-2xl bg-gradient-to-br from-pink-50 to-orange-50 hover:shadow-xl hover:scale-105 transition-all duration-300 group transform hover:-translate-y-1">
              <div class="w-12 h-12 sm:w-16 sm:h-16 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center mx-auto mb-4 sm:mb-6 group-hover:scale-110 transition-transform duration-300 shadow-lg">
                <HeartIcon class="h-6 w-6 sm:h-8 sm:w-8 text-white" />
              </div>
              <h3 class="text-xl sm:text-2xl font-bold text-gray-900 mb-3 sm:mb-4">专业可靠</h3>
              <p class="text-gray-600 leading-relaxed text-sm sm:text-base">
                汇聚儿科医生、营养师、教育专家的专业知识，确保每一条建议都科学可靠
              </p>
            </div>
            
            <!-- 温暖社区 -->
            <div class="text-center p-6 sm:p-8 rounded-xl sm:rounded-2xl bg-gradient-to-br from-blue-50 to-purple-50 hover:shadow-xl hover:scale-105 transition-all duration-300 group transform hover:-translate-y-1">
              <div class="w-12 h-12 sm:w-16 sm:h-16 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full flex items-center justify-center mx-auto mb-4 sm:mb-6 group-hover:scale-110 transition-transform duration-300 shadow-lg">
                <UsersIcon class="h-6 w-6 sm:h-8 sm:w-8 text-white" />
              </div>
              <h3 class="text-xl sm:text-2xl font-bold text-gray-900 mb-3 sm:mb-4">温暖社区</h3>
              <p class="text-gray-600 leading-relaxed text-sm sm:text-base">
                连接全球父母，分享育儿经验，在这里你永远不会感到孤单
              </p>
            </div>
            
            <!-- 持续更新 -->
            <div class="text-center p-6 sm:p-8 rounded-xl sm:rounded-2xl bg-gradient-to-br from-green-50 to-teal-50 hover:shadow-xl hover:scale-105 transition-all duration-300 group transform hover:-translate-y-1">
              <div class="w-12 h-12 sm:w-16 sm:h-16 bg-gradient-to-r from-green-500 to-teal-500 rounded-full flex items-center justify-center mx-auto mb-4 sm:mb-6 group-hover:scale-110 transition-transform duration-300 shadow-lg">
                <TrendingUpIcon class="h-6 w-6 sm:h-8 sm:w-8 text-white" />
              </div>
              <h3 class="text-xl sm:text-2xl font-bold text-gray-900 mb-3 sm:mb-4">持续更新</h3>
              <p class="text-gray-600 leading-relaxed text-sm sm:text-base">
                紧跟最新育儿研究成果，定期更新内容，让您的育儿知识始终保持前沿
              </p>
            </div>
          </div>
        </div>
      </section>

      <!-- 统计数据 -->
      <section class="py-16 sm:py-20 bg-gradient-to-r from-pink-600 via-purple-600 to-orange-600 relative overflow-hidden">
        <!-- 背景装饰 -->
        <div class="absolute inset-0 bg-black bg-opacity-10"></div>
        <div class="absolute top-0 left-0 w-full h-full">
          <div class="absolute top-10 left-10 w-20 h-20 bg-white bg-opacity-10 rounded-full animate-pulse"></div>
          <div class="absolute bottom-10 right-10 w-16 h-16 bg-white bg-opacity-10 rounded-full animate-pulse delay-1000"></div>
          <div class="absolute top-1/2 left-1/4 w-12 h-12 bg-white bg-opacity-10 rounded-full animate-pulse delay-500"></div>
        </div>
        
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
          <div class="text-center mb-12">
            <h2 class="text-2xl sm:text-3xl md:text-4xl font-bold text-white mb-4">
              数据见证我们的成长
            </h2>
            <p class="text-lg text-pink-100 max-w-2xl mx-auto">
              每一个数字背后都是我们对育儿事业的坚持与热爱
            </p>
          </div>
          
          <div class="grid grid-cols-2 md:grid-cols-4 gap-6 sm:gap-8">
            <div class="text-center group">
              <div class="bg-white bg-opacity-20 backdrop-blur-sm rounded-2xl p-6 sm:p-8 hover:bg-opacity-30 transition-all duration-300 hover:scale-105">
                <div class="text-3xl sm:text-4xl md:text-5xl font-bold text-white mb-2 group-hover:scale-110 transition-transform duration-300">10K+</div>
                <div class="text-pink-100 text-sm sm:text-lg font-medium">注册用户</div>
                <div class="text-pink-200 text-xs mt-1">信任之选</div>
              </div>
            </div>
            
            <div class="text-center group">
              <div class="bg-white bg-opacity-20 backdrop-blur-sm rounded-2xl p-6 sm:p-8 hover:bg-opacity-30 transition-all duration-300 hover:scale-105">
                <div class="text-3xl sm:text-4xl md:text-5xl font-bold text-white mb-2 group-hover:scale-110 transition-transform duration-300">5K+</div>
                <div class="text-pink-100 text-sm sm:text-lg font-medium">优质文章</div>
                <div class="text-pink-200 text-xs mt-1">知识宝库</div>
              </div>
            </div>
            
            <div class="text-center group">
              <div class="bg-white bg-opacity-20 backdrop-blur-sm rounded-2xl p-6 sm:p-8 hover:bg-opacity-30 transition-all duration-300 hover:scale-105">
                <div class="text-3xl sm:text-4xl md:text-5xl font-bold text-white mb-2 group-hover:scale-110 transition-transform duration-300">100+</div>
                <div class="text-pink-100 text-sm sm:text-lg font-medium">专业作者</div>
                <div class="text-pink-200 text-xs mt-1">权威专家</div>
              </div>
            </div>
            
            <div class="text-center group">
              <div class="bg-white bg-opacity-20 backdrop-blur-sm rounded-2xl p-6 sm:p-8 hover:bg-opacity-30 transition-all duration-300 hover:scale-105">
                <div class="text-3xl sm:text-4xl md:text-5xl font-bold text-white mb-2 group-hover:scale-110 transition-transform duration-300">24/7</div>
                <div class="text-pink-100 text-sm sm:text-lg font-medium">在线支持</div>
                <div class="text-pink-200 text-xs mt-1">贴心服务</div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 最新文章预览 -->
      <section class="py-12 sm:py-16 md:py-20 bg-gray-50">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="text-center mb-8 sm:mb-12 md:mb-16">
            <h2 class="text-2xl sm:text-3xl md:text-4xl font-bold text-gray-900 mb-3 sm:mb-4">最新文章</h2>
            <p class="text-base sm:text-lg md:text-xl text-gray-600 max-w-3xl mx-auto px-4 sm:px-0">
              发现最新的育儿知识和经验分享
            </p>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 sm:gap-8 mb-8 sm:mb-12">
            <!-- 文章卡片示例 -->
            <article class="bg-white rounded-xl sm:rounded-2xl shadow-lg overflow-hidden hover:shadow-2xl hover:scale-105 transition-all duration-300 group">
              <div class="h-40 sm:h-48 bg-gradient-to-br from-pink-400 to-orange-400 relative overflow-hidden">
                <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-10 transition-all duration-300"></div>
                <div class="absolute top-3 sm:top-4 left-3 sm:left-4">
                  <span class="bg-white bg-opacity-90 text-pink-800 text-xs font-medium px-2 sm:px-3 py-1 rounded-full backdrop-blur-sm">热门</span>
                </div>
              </div>
              <div class="p-4 sm:p-6">
                <div class="flex items-center mb-2 sm:mb-3">
                  <span class="bg-pink-100 text-pink-800 text-xs font-medium px-2 sm:px-2.5 py-0.5 rounded-full">新生儿护理</span>
                  <span class="text-gray-500 text-xs sm:text-sm ml-auto">2024-01-15</span>
                </div>
                <h3 class="text-base sm:text-lg md:text-xl font-bold text-gray-900 mb-2 sm:mb-3 line-clamp-2 group-hover:text-pink-600 transition-colors duration-300">
                  新生儿睡眠指南：如何帮助宝宝建立健康的睡眠习惯
                </h3>
                <p class="text-gray-600 mb-3 sm:mb-4 line-clamp-3 text-xs sm:text-sm md:text-base">
                  良好的睡眠对新生儿的成长发育至关重要。本文将为您详细介绍如何帮助宝宝建立规律的睡眠模式...
                </p>
                <div class="flex items-center justify-between">
                  <div class="flex items-center">
                    <div class="w-6 h-6 sm:w-8 sm:h-8 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center shadow-md">
                      <span class="text-white text-xs sm:text-sm font-medium">李</span>
                    </div>
                    <span class="text-gray-700 text-xs sm:text-sm ml-2 font-medium">李医生</span>
                  </div>
                  <div class="flex items-center text-gray-500 text-xs sm:text-sm">
                    <HeartIcon class="h-3 w-3 sm:h-4 sm:w-4 mr-1 group-hover:text-pink-500 transition-colors duration-300" />
                    <span>128</span>
                  </div>
                </div>
              </div>
            </article>
            
            <article class="bg-white rounded-xl sm:rounded-2xl shadow-lg overflow-hidden hover:shadow-2xl hover:scale-105 transition-all duration-300 group">
              <div class="h-40 sm:h-48 bg-gradient-to-br from-blue-400 to-purple-400 relative overflow-hidden">
                <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-10 transition-all duration-300"></div>
                <div class="absolute top-3 sm:top-4 left-3 sm:left-4">
                  <span class="bg-white bg-opacity-90 text-blue-800 text-xs font-medium px-2 sm:px-3 py-1 rounded-full backdrop-blur-sm">推荐</span>
                </div>
              </div>
              <div class="p-4 sm:p-6">
                <div class="flex items-center mb-2 sm:mb-3">
                  <span class="bg-blue-100 text-blue-800 text-xs font-medium px-2 sm:px-2.5 py-0.5 rounded-full">营养健康</span>
                  <span class="text-gray-500 text-xs sm:text-sm ml-auto">2024-01-14</span>
                </div>
                <h3 class="text-base sm:text-lg md:text-xl font-bold text-gray-900 mb-2 sm:mb-3 line-clamp-2 group-hover:text-blue-600 transition-colors duration-300">
                  婴幼儿辅食添加时间表：科学喂养从这里开始
                </h3>
                <p class="text-gray-600 mb-3 sm:mb-4 line-clamp-3 text-xs sm:text-sm md:text-base">
                  正确的辅食添加时机和方法对宝宝的健康成长非常重要。让我们一起了解科学的辅食添加方案...
                </p>
                <div class="flex items-center justify-between">
                  <div class="flex items-center">
                    <div class="w-6 h-6 sm:w-8 sm:h-8 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full flex items-center justify-center shadow-md">
                      <span class="text-white text-xs sm:text-sm font-medium">王</span>
                    </div>
                    <span class="text-gray-700 text-xs sm:text-sm ml-2 font-medium">王营养师</span>
                  </div>
                  <div class="flex items-center text-gray-500 text-xs sm:text-sm">
                    <HeartIcon class="h-3 w-3 sm:h-4 sm:w-4 mr-1 group-hover:text-blue-500 transition-colors duration-300" />
                    <span>95</span>
                  </div>
                </div>
              </div>
            </article>
            
            <article class="bg-white rounded-xl sm:rounded-2xl shadow-lg overflow-hidden hover:shadow-2xl hover:scale-105 transition-all duration-300 group">
              <div class="h-40 sm:h-48 bg-gradient-to-br from-green-400 to-teal-400 relative overflow-hidden">
                <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-10 transition-all duration-300"></div>
                <div class="absolute top-3 sm:top-4 left-3 sm:left-4">
                  <span class="bg-white bg-opacity-90 text-green-800 text-xs font-medium px-2 sm:px-3 py-1 rounded-full backdrop-blur-sm">精选</span>
                </div>
              </div>
              <div class="p-4 sm:p-6">
                <div class="flex items-center mb-2 sm:mb-3">
                  <span class="bg-green-100 text-green-800 text-xs font-medium px-2 sm:px-2.5 py-0.5 rounded-full">早期教育</span>
                  <span class="text-gray-500 text-xs sm:text-sm ml-auto">2024-01-13</span>
                </div>
                <h3 class="text-base sm:text-lg md:text-xl font-bold text-gray-900 mb-2 sm:mb-3 line-clamp-2 group-hover:text-green-600 transition-colors duration-300">
                  0-3岁宝宝语言发展里程碑及促进方法
                </h3>
                <p class="text-gray-600 mb-3 sm:mb-4 line-clamp-3 text-xs sm:text-sm md:text-base">
                  语言发展是宝宝成长的重要指标。了解各个阶段的发展特点，帮助宝宝更好地发展语言能力...
                </p>
                <div class="flex items-center justify-between">
                  <div class="flex items-center">
                    <div class="w-6 h-6 sm:w-8 sm:h-8 bg-gradient-to-r from-green-500 to-teal-500 rounded-full flex items-center justify-center shadow-md">
                      <span class="text-white text-xs sm:text-sm font-medium">张</span>
                    </div>
                    <span class="text-gray-700 text-xs sm:text-sm ml-2 font-medium">张老师</span>
                  </div>
                  <div class="flex items-center text-gray-500 text-xs sm:text-sm">
                    <HeartIcon class="h-3 w-3 sm:h-4 sm:w-4 mr-1 group-hover:text-green-500 transition-colors duration-300" />
                    <span>156</span>
                  </div>
                </div>
              </div>
            </article>
          </div>
          
          <div class="text-center">
            <router-link 
              to="/articles" 
              class="inline-flex items-center px-4 sm:px-6 md:px-8 py-2 sm:py-3 bg-gradient-to-r from-pink-600 to-orange-600 text-white font-medium rounded-full hover:from-pink-700 hover:to-orange-700 transition-all duration-300 hover:scale-105 shadow-lg hover:shadow-xl text-xs sm:text-sm md:text-base"
            >
              查看更多文章
              <ArrowRightIcon class="ml-1 sm:ml-2 h-3 w-3 sm:h-4 sm:w-4 md:h-5 md:w-5" />
            </router-link>
          </div>
        </div>
      </section>
    </main>

    <!-- 页脚 -->
    <footer class="bg-gray-900 text-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 sm:py-16">
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6 sm:gap-8">
          <!-- 品牌信息 -->
          <div class="sm:col-span-2 md:col-span-1">
            <div class="flex items-center mb-3 sm:mb-4">
              <div class="w-6 h-6 sm:w-8 sm:h-8 bg-gradient-to-r from-pink-500 to-orange-500 rounded-lg flex items-center justify-center mr-2 sm:mr-3 shadow-lg">
                <span class="text-white font-bold text-sm sm:text-lg">G</span>
              </div>
              <span class="text-lg sm:text-xl font-bold">GoDad</span>
            </div>
            <p class="text-gray-400 mb-3 sm:mb-6 text-sm sm:text-base leading-relaxed">
              专业的育儿知识分享平台，陪伴每一个家庭的成长之路。
            </p>
            <div class="flex space-x-2 sm:space-x-3">
              <div class="w-8 h-8 sm:w-10 sm:h-10 bg-gray-800 rounded-full flex items-center justify-center hover:bg-gradient-to-r hover:from-pink-600 hover:to-orange-600 transition-all duration-300 cursor-pointer group shadow-lg">
                <span class="text-xs sm:text-sm font-medium group-hover:text-white">微</span>
              </div>
              <div class="w-8 h-8 sm:w-10 sm:h-10 bg-gray-800 rounded-full flex items-center justify-center hover:bg-gradient-to-r hover:from-pink-600 hover:to-orange-600 transition-all duration-300 cursor-pointer group shadow-lg">
                <span class="text-xs sm:text-sm font-medium group-hover:text-white">博</span>
              </div>
              <div class="w-8 h-8 sm:w-10 sm:h-10 bg-gray-800 rounded-full flex items-center justify-center hover:bg-gradient-to-r hover:from-pink-600 hover:to-orange-600 transition-all duration-300 cursor-pointer group shadow-lg">
                <span class="text-xs sm:text-sm font-medium group-hover:text-white">抖</span>
              </div>
            </div>
          </div>
          
          <!-- 快速链接 -->
          <div>
            <h3 class="text-base sm:text-lg font-semibold mb-3 sm:mb-4 text-white">快速链接</h3>
            <ul class="space-y-2 sm:space-y-3">
              <li><a href="#" class="text-gray-400 hover:text-pink-400 transition-colors duration-300 text-sm sm:text-base hover:translate-x-1 transform inline-block">关于我们</a></li>
              <li><a href="#" class="text-gray-400 hover:text-pink-400 transition-colors duration-300 text-sm sm:text-base hover:translate-x-1 transform inline-block">联系我们</a></li>
              <li><a href="#" class="text-gray-400 hover:text-pink-400 transition-colors duration-300 text-sm sm:text-base hover:translate-x-1 transform inline-block">隐私政策</a></li>
              <li><a href="#" class="text-gray-400 hover:text-pink-400 transition-colors duration-300 text-sm sm:text-base hover:translate-x-1 transform inline-block">服务条款</a></li>
              <li><a href="#" class="text-gray-400 hover:text-pink-400 transition-colors duration-300 text-sm sm:text-base hover:translate-x-1 transform inline-block">帮助中心</a></li>
            </ul>
          </div>
          
          <!-- 热门分类 -->
          <div>
            <h3 class="text-base sm:text-lg font-semibold mb-3 sm:mb-4 text-white">热门分类</h3>
            <ul class="space-y-2 sm:space-y-3">
              <li><a href="#" class="text-gray-400 hover:text-orange-400 transition-colors duration-300 text-sm sm:text-base hover:translate-x-1 transform inline-block">新生儿护理</a></li>
              <li><a href="#" class="text-gray-400 hover:text-orange-400 transition-colors duration-300 text-sm sm:text-base hover:translate-x-1 transform inline-block">营养健康</a></li>
              <li><a href="#" class="text-gray-400 hover:text-orange-400 transition-colors duration-300 text-sm sm:text-base hover:translate-x-1 transform inline-block">早期教育</a></li>
              <li><a href="#" class="text-gray-400 hover:text-orange-400 transition-colors duration-300 text-sm sm:text-base hover:translate-x-1 transform inline-block">心理发展</a></li>
              <li><a href="#" class="text-gray-400 hover:text-orange-400 transition-colors duration-300 text-sm sm:text-base hover:translate-x-1 transform inline-block">安全防护</a></li>
            </ul>
          </div>
          
          <!-- 联系我们 -->
          <div>
            <h3 class="text-base sm:text-lg font-semibold mb-3 sm:mb-4 text-white">联系我们</h3>
            <div class="space-y-3 sm:space-y-4">
              <div class="flex items-center group">
                <MailIcon class="h-4 w-4 sm:h-5 sm:w-5 text-gray-400 mr-2 sm:mr-3 group-hover:text-pink-400 transition-colors duration-300" />
                <span class="text-gray-400 text-sm sm:text-base group-hover:text-white transition-colors duration-300">support@godad.com</span>
              </div>
              <div class="flex items-center group">
                <PhoneIcon class="h-4 w-4 sm:h-5 sm:w-5 text-gray-400 mr-2 sm:mr-3 group-hover:text-pink-400 transition-colors duration-300" />
                <span class="text-gray-400 text-sm sm:text-base group-hover:text-white transition-colors duration-300">400-123-4567</span>
              </div>
              <div class="flex items-start group">
                <MapPinIcon class="h-4 w-4 sm:h-5 sm:w-5 text-gray-400 mr-2 sm:mr-3 mt-0.5 group-hover:text-pink-400 transition-colors duration-300" />
                <span class="text-gray-400 text-sm sm:text-base group-hover:text-white transition-colors duration-300">北京市朝阳区科技园区</span>
              </div>
            </div>
          </div>
        </div>
        
        <div class="border-t border-gray-800 pt-6 sm:pt-8">
          <div class="flex flex-col sm:flex-row justify-between items-center space-y-3 sm:space-y-0">
            <p class="text-gray-400 text-xs sm:text-sm text-center sm:text-left">
              © 2024 GoDad. 保留所有权利。
            </p>
            <div class="flex flex-col sm:flex-row space-y-1 sm:space-y-0 sm:space-x-4 text-xs sm:text-sm text-center sm:text-right">
              <a href="#" class="text-gray-400 hover:text-white transition-colors duration-300">备案号：京ICP备12345678号</a>
              <a href="#" class="text-gray-400 hover:text-white transition-colors duration-300">京公网安备 11010502012345号</a>
            </div>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>