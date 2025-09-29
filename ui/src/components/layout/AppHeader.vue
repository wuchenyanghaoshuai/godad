<template>
  <nav class="bg-white shadow-sm sticky top-0 z-50 border-b border-gray-100">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-14 sm:h-16">
        <!-- Logo 和导航链接 -->
        <div class="flex items-center space-x-4 sm:space-x-8">
          <!-- Logo -->
          <router-link to="/" class="flex items-center group">
            <img
              src="@/assets/images/logo/GoDad_logo.png"
              alt="GoDad Logo"
              class="h-14 sm:h-16 w-auto transition-all duration-300 group-hover:scale-105"
              style="background: white; border-radius: 8px; padding: 2px;"
            />
          </router-link>

          <!-- 桌面端导航链接 -->
          <div v-if="showNavigation" class="hidden md:flex items-center space-x-6 lg:space-x-8">
            <router-link
              v-for="navItem in navItems"
              :key="navItem.path"
              :to="navItem.path"
              class="text-gray-700 hover:text-pink-600 transition-all duration-300 text-sm font-medium relative group"
              :class="{ 'text-pink-600': isActiveRoute(navItem.path) }"
            >
              {{ navItem.label }}
              <span class="absolute -bottom-1 left-0 w-0 h-0.5 bg-pink-600 transition-all duration-300 group-hover:w-full"
                    :class="{ 'w-full': isActiveRoute(navItem.path) }"></span>
            </router-link>
          </div>
        </div>

        <!-- 中间区域：搜索框 -->
        <SearchBox :show-search="showSearch" />

        <!-- 右侧区域 -->
        <div class="flex items-center space-x-2 sm:space-x-4">
          <!-- 搜索按钮（移动端和中等屏幕） -->
          <button
            v-if="showSearch"
            class="lg:hidden p-2 rounded-lg text-gray-600 hover:text-pink-600 hover:bg-pink-50 transition-all duration-200"
            @click="$router.push('/search')"
          >
            <SearchIcon class="h-5 w-5" />
          </button>

          <!-- 移动端菜单按钮 -->
          <button
            v-if="showNavigation"
            class="md:hidden p-2 rounded-lg text-gray-600 hover:text-pink-600 hover:bg-pink-50 transition-all duration-200"
            @click="showMobileMenu = !showMobileMenu"
          >
            <MenuIcon v-if="!showMobileMenu" class="h-5 w-5" />
            <XIcon v-else class="h-5 w-5" />
          </button>

          <!-- 用户区域 -->
          <div class="flex items-center space-x-2 sm:space-x-3">
            <!-- 已登录状态 -->
            <div v-if="authStore.isAuthenticated" class="flex items-center space-x-2 sm:space-x-3">
              <!-- 发布文章按钮 -->
              <router-link
                v-if="showCreateButton"
                to="/articles/create"
                class="hidden sm:flex items-center px-3 py-1.5 bg-gradient-to-r from-pink-500 to-rose-400 text-white rounded-lg hover:from-pink-600 hover:to-rose-500 transition-all duration-300 text-sm font-medium shadow-md hover:shadow-lg transform hover:scale-105"
              >
                <PlusIcon class="h-4 w-4 mr-1.5" />
                发布文章
              </router-link>

              <!-- 消息通知按钮 + 下拉菜单 -->
              <NotificationDropdown
                :show-notifications="showNotifications"
                :total-unread-count="totalUnreadCount"
                :unread-by-type="unreadByType"
                @fetch-unread-by-type="fetchUnreadByType"
              />

              <!-- 用户头像菜单 -->
              <div v-if="showUserMenu" class="relative">
                <button
                  ref="userMenuButton"
                  @click="showUserMenuDropdown = !showUserMenuDropdown"
                  class="flex items-center space-x-2 p-1.5 rounded-lg hover:bg-gray-50 transition-all duration-200"
                >
                  <UserAvatar
                    :avatar="authStore.user?.avatar || ''"
                    :name="authStore.user?.nickname || authStore.user?.username || 'U'"
                    :size="36"
                    wrapper-class="border-2 border-pink-200 hover:border-pink-300 transition-colors"
                  />
                  <ChevronDownIcon class="hidden sm:block h-4 w-4 text-gray-400" />
                </button>

                <!-- 用户下拉菜单 -->
                <div
                  v-if="showUserMenuDropdown"
                  class="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-100 py-1 z-50"
                >
                  <div class="px-4 py-2 border-b border-gray-100">
                    <div class="text-sm font-medium text-gray-900 truncate">{{ authStore.user?.nickname || authStore.user?.username }}</div>
                    <div class="text-xs text-gray-500 truncate" :title="authStore.user?.email">{{ authStore.user?.email }}</div>
                    <!-- 用户积分显示 -->
                    <div v-if="showUserPoints" class="mt-2">
                      <UserPointsDisplay mode="simple" size="sm" />
                    </div>
                  </div>
                  <router-link
                    to="/user-center"
                    class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-pink-50 hover:text-pink-600 transition-colors"
                    @click="showUserMenuDropdown = false"
                  >
                    <UserIcon class="h-4 w-4 mr-3" />
                    个人中心
                  </router-link>
                  <router-link
                    :to="`/users/${authStore.user?.username}`"
                    class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-pink-50 hover:text-pink-600 transition-colors"
                    @click="showUserMenuDropdown = false"
                  >
                    <UserIcon class="h-4 w-4 mr-3" />
                    我的资料
                  </router-link>
                  <router-link
                    v-if="showCreateButton"
                    to="/articles/create"
                    class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-pink-50 hover:text-pink-600 transition-colors"
                    @click="showUserMenuDropdown = false"
                  >
                    <PlusIcon class="h-4 w-4 mr-3" />
                    发布文章
                  </router-link>
                  <button
                    @click="handleLogout"
                    class="flex items-center w-full px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors"
                  >
                    <LogOutIcon class="h-4 w-4 mr-3" />
                    退出登录
                  </button>
                </div>
              </div>
            </div>

            <!-- 未登录状态 -->
            <div v-else class="flex items-center space-x-2">
              <router-link
                to="/login"
                class="px-3 py-1.5 text-gray-700 hover:text-pink-600 transition-colors text-sm font-medium"
              >
                登录
              </router-link>
              <router-link
                to="/register"
                class="px-3 py-1.5 bg-gradient-to-r from-pink-500 to-rose-400 text-white rounded-lg hover:from-pink-600 hover:to-rose-500 transition-all duration-300 text-sm font-medium shadow-md hover:shadow-lg transform hover:scale-105"
              >
                注册
              </router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- 移动端导航菜单 -->
      <div
        v-if="showNavigation && showMobileMenu"
        class="md:hidden border-t border-gray-100 py-2"
      >
        <div class="flex flex-col space-y-1">
          <router-link
            v-for="navItem in navItems"
            :key="navItem.path"
            :to="navItem.path"
            class="px-4 py-2 text-gray-700 hover:text-pink-600 hover:bg-pink-50 transition-colors text-sm font-medium rounded-lg"
            :class="{ 'text-pink-600 bg-pink-50': isActiveRoute(navItem.path) }"
            @click="showMobileMenu = false"
          >
            {{ navItem.label }}
          </router-link>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  MenuIcon,
  XIcon,
  UserIcon,
  LogOutIcon,
  ChevronDownIcon,
  PlusIcon,
  SearchIcon
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { NotificationApi } from '@/api/notification'
import { useNotificationSync } from '@/composables/useNotificationSync'
import UserPointsDisplay from '../UserPointsDisplay.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import SearchBox from '@/components/SearchBox.vue'
import NotificationDropdown from '@/components/NotificationDropdown.vue'

// Props - 允许每个页面自定义配置
interface NavItem {
  path: string
  label: string
}

withDefaults(defineProps<{
  navItems?: NavItem[]
  showSearch?: boolean
  showCreateButton?: boolean
  showNotifications?: boolean
  showUserPoints?: boolean
  showNavigation?: boolean
  showUserMenu?: boolean
}>(), {
  navItems: () => [
    { path: '/articles', label: '文章' },
    { path: '/community', label: '社区' },
    { path: '/resources', label: '资源' }
  ],
  showSearch: true,
  showCreateButton: true,
  showNotifications: true,
  showUserPoints: true,
  showNavigation: true,
  showUserMenu: true
})

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const { onNotificationEvent } = useNotificationSync()

// 响应式数据
const showUserMenuDropdown = ref(false)
const showMobileMenu = ref(false)
const userMenuButton = ref<HTMLElement>()
const notifMenuWrapper = ref<HTMLElement | null>(null)
const showNotifMenu = ref(false)
let notifMenuHideTimer: number | null = null
const unreadMessagesCount = ref(0)
const unreadNotificationsCount = ref(0)
const unreadByType = ref<{ [k: string]: number }>({ message: 0, like: 0, comment: 0, follow: 0, bookmark: 0, system: 0, mention: 0, total_unread: 0 })
let unreadByTypeFetchedAt = 0


// 方法
const isActiveRoute = (path: string): boolean => {
  if (path === '/articles') {
    return route.path.startsWith('/articles')
  }
  return route.path === path
}

// 获取所有未读计数
const fetchUnreadCounts = async () => {
  if (!authStore.isAuthenticated) {
    unreadMessagesCount.value = 0
    unreadNotificationsCount.value = 0
    return
  }

  try {
    const response = await NotificationApi.getNotificationStats()
    const totalCount = response.data.unread_count || 0
    unreadMessagesCount.value = totalCount
    unreadNotificationsCount.value = totalCount
  } catch (error) {
    console.error('获取未读计数失败:', error)
    unreadMessagesCount.value = 0
    unreadNotificationsCount.value = 0
  }
}

// 各类型未读统计（A方案）
const fetchUnreadByType = async (force = false) => {
  if (!authStore.isAuthenticated) {
    unreadByType.value = { message: 0, like: 0, comment: 0, follow: 0, bookmark: 0, system: 0, mention: 0, total_unread: 0 }
    return
  }
  const now = Date.now()
  if (!force && now - unreadByTypeFetchedAt < 30000 && showNotifMenu.value) return
  try {
    const resp = await NotificationApi.getNotificationTypeStats()
    const data: any = resp.data || {}
    unreadByType.value = {
      message: data.message || 0,
      like: data.like || 0,
      comment: data.comment || 0,
      follow: data.follow || 0,
      bookmark: data.bookmark || 0,
      system: data.system || 0,
      mention: data.mention || 0,
      total_unread: data.total_unread || 0
    }
    unreadByTypeFetchedAt = now
  } catch (e) {
    // 静默失败，不影响下拉可用
  }
}

// 清除未读计数
const clearUnreadNotifications = async () => {
  if (!authStore.isAuthenticated || totalUnreadCount.value === 0) {
    return
  }

  try {
    await NotificationApi.markAllAsRead()
    unreadMessagesCount.value = 0
    unreadNotificationsCount.value = 0
  } catch (error) {
    console.error('标记通知已读失败:', error)
  }
}


// 退出登录
const handleLogout = async () => {
  try {
    await authStore.logout()
    showUserMenuDropdown.value = false
    router.push('/login')
  } catch (error) {
    console.error('退出登录失败:', error)
  }
}

// 点击外部关闭菜单
const handleClickOutside = (event: Event) => {
  if (userMenuButton.value && !userMenuButton.value.contains(event.target as Node)) {
    showUserMenuDropdown.value = false
  }
  if (notifMenuWrapper.value && !notifMenuWrapper.value.contains(event.target as Node)) {
    showNotifMenu.value = false
  }
  if (showMobileMenu.value && !(event.target as Element).closest('nav')) {
    showMobileMenu.value = false
  }
}

// 计算属性：总未读数量
const totalUnreadCount = computed(() => {
  return Math.max(unreadMessagesCount.value, unreadNotificationsCount.value)
})

// 计算属性：其他通知数量（收藏、关注等）
const otherUnreadCount = computed(() => {
  return (unreadByType.value.bookmark || 0) + (unreadByType.value.follow || 0)
})

// 监听路由变化（进入通知中心不再自动全部已读，避免误清空未读）
// 保留为空以避免自动清零未读
watch(() => route.path, async (_newPath) => {
  // no-op
})

// 监听认证状态变化
watch(() => authStore.isAuthenticated, async (isAuthenticated) => {
  if (isAuthenticated) {
    await fetchUnreadCounts()
  } else {
    unreadMessagesCount.value = 0
    unreadNotificationsCount.value = 0
  }
})

// 组件挂载时
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  fetchUnreadCounts()
  // 初始预取一次分类计数（不阻塞UI）
  fetchUnreadByType()

  onNotificationEvent('refresh', () => {
    fetchUnreadCounts()
    fetchUnreadByType(true)
  })
})

// 组件卸载时
onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

// 通知下拉控制与路由跳转
const cancelNotifMenuHide = () => {
  if (notifMenuHideTimer) {
    clearTimeout(notifMenuHideTimer)
    notifMenuHideTimer = null
  }
}
const openNotifMenuHover = () => {
  cancelNotifMenuHide()
  showNotifMenu.value = true
  fetchUnreadByType()
}
const closeNotifMenuHover = () => {
  cancelNotifMenuHide()
  notifMenuHideTimer = window.setTimeout(() => {
    showNotifMenu.value = false
    notifMenuHideTimer = null
  }, 220)
}
const toggleNotifMenu = () => { showNotifMenu.value = !showNotifMenu.value }
const goNotifications = (category: 'message' | 'like' | 'comment' | 'system' | 'mention' | 'other') => {
  showNotifMenu.value = false
  if (category === 'message') {
    router.push({ path: '/notifications', query: { tab: 'message' } })
    return
  }
  // 其他通知类型，包含收藏和关注
  if (category === 'other') {
    router.push({ path: '/notifications', query: { tab: 'notify', category: 'bookmark,follow' } })
    return
  }
  // 其余归类到通知
  const query: any = { tab: 'notify' }
  if (['like', 'comment', 'follow', 'bookmark', 'system', 'mention'].includes(category)) {
    query.category = category
  }
  router.push({ path: '/notifications', query })
}
</script>
