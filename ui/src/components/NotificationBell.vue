<template>
  <div class="notification-bell relative" ref="notificationBellRef">
    <button
      @click="toggleDropdown"
      class="relative p-2 text-gray-600 hover:text-pink-600 hover:bg-pink-50 rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-pink-500 focus:ring-opacity-50"
      :class="{ 'text-pink-600 bg-pink-50': showDropdown }"
    >
      <!-- 邮箱图标 -->
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
      </svg>
      
      <!-- 未读数量小红点 -->
      <span
        v-if="unreadCount > 0"
        class="absolute -top-1 -right-1 inline-flex items-center justify-center w-5 h-5 text-xs font-bold text-white bg-red-500 rounded-full animate-pulse"
      >
        {{ unreadCount > 99 ? '99+' : unreadCount }}
      </span>
    </button>

    <!-- 下拉菜单 -->
    <div
      v-if="showDropdown"
      class="absolute right-0 mt-2 w-80 bg-white border border-gray-200 rounded-lg shadow-lg z-[70] max-h-96 overflow-hidden"
    >
      <!-- 菜单头部 -->
      <div class="px-4 py-3 bg-gray-50 border-b border-gray-200">
        <div class="flex items-center justify-between">
          <h3 class="text-sm font-semibold text-gray-800">消息通知</h3>
          <div class="flex items-center space-x-2">
            <span v-if="unreadCount > 0" class="text-xs text-gray-500">
              {{ unreadCount }}条未读
            </span>
            <button
              v-if="unreadCount > 0"
              @click="markAllAsRead"
              class="text-xs text-pink-600 hover:text-pink-700"
            >
              全部已读
            </button>
          </div>
        </div>
      </div>

      <!-- 通知列表 -->
      <div class="max-h-72 overflow-y-auto">
        <!-- 加载状态 -->
        <div v-if="loading" class="flex justify-center py-6">
          <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-pink-500"></div>
        </div>

        <!-- 空状态 -->
        <div v-else-if="!notifications || notifications.length === 0" class="text-center py-8">
          <div class="w-12 h-12 mx-auto mb-2 text-gray-300">
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
          </div>
          <p class="text-sm text-gray-500">暂无新消息</p>
        </div>

        <!-- 通知项 -->
        <div v-else>
          <div
            v-for="notification in displayNotifications"
            :key="notification.id"
            class="px-4 py-3 border-b border-gray-100 hover:bg-gray-50 cursor-pointer transition-colors"
            :class="{
              'bg-blue-50': !notification.is_read
            }"
            @click="handleNotificationClick(notification)"
          >
            <div class="flex items-start space-x-3">
              <!-- 用户头像 -->
              <div class="flex-shrink-0">
                <img 
                  v-if="notification.actor_avatar"
                  :src="notification.actor_avatar" 
                  :alt="notification.actor_nickname || notification.actor_username"
                  class="w-8 h-8 rounded-full object-cover"
                >
                <div 
                  v-else
                  class="w-8 h-8 bg-gradient-to-br from-pink-400 to-rose-400 rounded-full flex items-center justify-center text-white font-semibold text-xs"
                >
                  {{ (notification.actor_nickname || notification.actor_username || 'U').charAt(0).toUpperCase() }}
                </div>
              </div>

              <!-- 通知内容 -->
              <div class="flex-1 min-w-0">
                <div class="flex items-start justify-between">
                  <div class="flex-1">
                    <p class="text-xs text-gray-700 line-clamp-2">
                      <span class="font-medium text-pink-600">
                        {{ notification.actor_nickname || notification.actor_username }}
                      </span>
                      {{ getShortMessage(notification.message) }}
                    </p>
                    <p class="text-xs text-gray-500 mt-1">
                      {{ formatNotificationTime(notification.created_at) }}
                    </p>
                  </div>
                  <!-- 未读指示器 -->
                  <div
                    v-if="!notification.is_read"
                    class="w-2 h-2 bg-pink-500 rounded-full ml-2 mt-1 flex-shrink-0"
                  ></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 菜单底部 -->
      <div class="px-4 py-2 bg-gray-50 border-t border-gray-200">
        <router-link
          to="/user-center?tab=notifications"
          @click="closeDropdown"
          class="block text-center text-sm text-pink-600 hover:text-pink-700 py-1"
        >
          查看所有消息
        </router-link>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { NotificationApi, type Notification, type NotificationStats, formatNotificationTime } from '@/api/notification'

// 数据状态
const notifications = ref<Notification[]>([])
const stats = ref<NotificationStats>()
const loading = ref(false)
const showDropdown = ref(false)
const notificationBellRef = ref<HTMLElement>()

// 路由
const router = useRouter()

// 计算属性
const unreadCount = computed(() => stats.value?.unread_count || 0)
const displayNotifications = computed(() => notifications.value?.slice(0, 5) || []) // 只显示前5条

// 获取简短消息
const getShortMessage = (message: string): string => {
  return message.length > 30 ? message.substring(0, 30) + '...' : message
}

// 从通知消息中提取评论内容
const extractCommentContent = (message: string): string => {
  // 通知消息格式通常是: "评论了你的文章《xxx》：评论内容" 或 "回复了你的评论：评论内容"
  const colonIndex = message.lastIndexOf('：')
  if (colonIndex !== -1 && colonIndex < message.length - 1) {
    return message.substring(colonIndex + 1).trim()
  }
  return ''
}

// 切换下拉菜单
const toggleDropdown = () => {
  showDropdown.value = !showDropdown.value
  if (showDropdown.value) {
    loadNotifications()
  }
}

// 关闭下拉菜单
const closeDropdown = () => {
  showDropdown.value = false
}

// 加载通知统计
const loadStats = async () => {
  try {
    const response = await NotificationApi.getNotificationStats()
    if (response.code === 200) {
      stats.value = response.data
    }
  } catch (error) {
    console.error('加载通知统计失败:', error)
  }
}

// 加载通知列表
const loadNotifications = async () => {
  loading.value = true
  try {
    const response = await NotificationApi.getNotifications({
      page: 1,
      limit: 10 // 下拉菜单只显示最近10条
    })

    if (response.code === 200) {
      notifications.value = response.data.notifications
    }
  } catch (error) {
    console.error('加载通知列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 标记单个通知为已读
const markAsRead = async (notificationIds: number[]) => {
  try {
    await NotificationApi.markAsRead(notificationIds)
    
    // 更新本地状态
    notifications.value?.forEach(notification => {
      if (notificationIds.includes(notification.id)) {
        notification.is_read = true
      }
    })
    
    // 更新统计
    await loadStats()
  } catch (error) {
    console.error('标记已读失败:', error)
  }
}

// 全部已读
const markAllAsRead = async () => {
  try {
    await NotificationApi.markAllAsRead()
    
    // 更新本地状态
    notifications.value?.forEach(notification => {
      notification.is_read = true
    })
    
    // 更新统计
    await loadStats()
  } catch (error) {
    console.error('全部标记已读失败:', error)
  }
}

// 处理通知点击
const handleNotificationClick = async (notification: Notification) => {
  // 处理通知点击
  
  // 关闭下拉菜单
  closeDropdown()
  
  try {
    // 先标记为已读
    if (!notification.is_read) {
      await markAsRead([notification.id])
    }
    
    if (notification.type === 'comment' && notification.resource_id) {
      // 评论通知跳转到文章并聚焦到具体评论
      // 从通知消息中提取评论内容
      const commentContent = extractCommentContent(notification.message)
      await router.push({
        path: `/articles/${notification.resource_id}`,
        query: { 
          focus: 'comment',
          actor_id: notification.actor_id.toString(),
          comment_content: commentContent
        }
      })
      
    } else if (notification.type === 'like' && notification.resource_id) {
      // 点赞通知跳转到文章
      await router.push(`/articles/${notification.resource_id}`)
      
    } else if (notification.type === 'follow') {
      // 关注通知跳转到用户页面
      await router.push(`/user/${notification.actor_id}`)
    }
  } catch (error) {
    console.error('NotificationBell: 处理通知点击失败:', error)
  }
}

// 自动刷新定时器
let refreshTimer: number | null = null

// 点击外部关闭下拉菜单
const handleClickOutside = (event: Event) => {
  if (showDropdown.value && notificationBellRef.value && !notificationBellRef.value.contains(event.target as Node)) {
    showDropdown.value = false
  }
}

const startAutoRefresh = () => {
  refreshTimer = window.setInterval(() => {
    loadStats()
  }, 3000) // 每3秒刷新一次统计
}

const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

// 组件挂载
onMounted(() => {
  loadStats()
  startAutoRefresh()
  document.addEventListener('click', handleClickOutside)
})

// 组件卸载
onUnmounted(() => {
  stopAutoRefresh()
  document.removeEventListener('click', handleClickOutside)
})

// 暴露方法给父组件
defineExpose({
  loadStats,
  unreadCount: computed(() => unreadCount.value)
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notification-bell {
  position: relative;
}
</style>