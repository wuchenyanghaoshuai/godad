<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { MenuIcon, XIcon, UserIcon, LogOutIcon, ChevronDownIcon, PlusIcon } from 'lucide-vue-next'
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
  <!-- 导航栏 -->
  <nav class="bg-white/95 backdrop-blur-sm shadow-sm sticky top-0 z-50 border-b border-gray-100">
    <div class="w-full px-4 sm:px-6 lg:px-8">
      <div class="flex items-center h-16">
        <!-- Logo - 完全靠左 -->
        <div class="flex items-center">
          <router-link to="/" class="flex items-center space-x-2 group">
            <div class="w-8 h-8 bg-gradient-to-r from-pink-500 to-orange-500 rounded-lg flex items-center justify-center group-hover:scale-105 transition-transform">
              <span class="text-white font-bold text-lg">G</span>
            </div>
            <span class="text-xl font-bold bg-gradient-to-r from-pink-600 to-orange-600 bg-clip-text text-transparent">GoDad</span>
          </router-link>
        </div>

        <!-- 右侧区域 - 完全靠右 -->
        <div class="flex items-center space-x-4 ml-auto">
          <!-- 导航链接 -->
          <div class="hidden md:flex items-center space-x-8">
            <router-link
              to="/articles"
              class="text-gray-700 hover:text-pink-600 transition-all duration-200 font-medium relative group"
            >
              文章
              <span class="absolute -bottom-1 left-0 w-0 h-0.5 bg-pink-600 transition-all duration-200 group-hover:w-full"></span>
            </router-link>
            <router-link
              to="/community"
              class="text-gray-700 hover:text-pink-600 transition-all duration-200 font-medium relative group"
            >
              社区
              <span class="absolute -bottom-1 left-0 w-0 h-0.5 bg-pink-600 transition-all duration-200 group-hover:w-full"></span>
            </router-link>
            <router-link
              to="/resources"
              class="text-gray-700 hover:text-pink-600 transition-all duration-200 font-medium relative group"
            >
              资源
              <span class="absolute -bottom-1 left-0 w-0 h-0.5 bg-pink-600 transition-all duration-200 group-hover:w-full"></span>
            </router-link>
          </div>
          <!-- 未登录状态 -->
          <div v-if="!authStore.isAuthenticated" class="hidden md:flex items-center space-x-3">
            <router-link
              to="/login"
              class="bg-gradient-to-r from-pink-600 to-orange-600 text-white px-6 py-2 rounded-full hover:shadow-lg hover:scale-105 transition-all duration-200 font-medium"
            >
              登录
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

      <!-- 移动端导航菜单 -->
      <div
        v-if="showMobileMenu"
        class="md:hidden border-t border-gray-100 py-2"
      >
        <div class="flex flex-col space-y-1">
          <router-link
            to="/articles"
            class="px-4 py-2 text-gray-700 hover:text-pink-600 hover:bg-pink-50 transition-colors text-sm font-medium rounded-lg"
            @click="showMobileMenu = false"
          >
            文章
          </router-link>
          <router-link
            to="/community"
            class="px-4 py-2 text-gray-700 hover:text-pink-600 hover:bg-pink-50 transition-colors text-sm font-medium rounded-lg"
            @click="showMobileMenu = false"
          >
            社区
          </router-link>
          <router-link
            to="/resources"
            class="px-4 py-2 text-gray-700 hover:text-pink-600 hover:bg-pink-50 transition-colors text-sm font-medium rounded-lg"
            @click="showMobileMenu = false"
          >
            资源
          </router-link>
        </div>
      </div>
    </div>
  </nav>
</template>
