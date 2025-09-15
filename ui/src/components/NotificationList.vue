<template>
  <div class="notification-list">
    <!-- 头部操作栏 -->
    <div class="flex items-center justify-between mb-4 p-4 bg-gray-50 rounded-lg">
      <div class="flex items-center space-x-4">
        <h3 class="text-lg font-semibold text-gray-800">消息通知</h3>
        <span class="px-2 py-1 text-xs bg-pink-100 text-pink-600 rounded-full" v-if="stats?.unread_count">
          {{ stats.unread_count }}条未读
        </span>
      </div>
      <div class="flex items-center space-x-2">
        <button 
          @click="markAllAsRead"
          v-if="stats?.unread_count"
          class="px-3 py-1 text-sm text-pink-600 hover:text-pink-700 hover:bg-pink-50 rounded-lg transition-colors"
        >
          全部已读
        </button>
        <button 
          @click="clearAllNotifications"
          v-if="stats?.total_count"
          class="px-3 py-1 text-sm text-red-600 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors"
          :disabled="clearing"
        >
          {{ clearing ? '清除中...' : '一键清除' }}
        </button>
        <button 
          @click="refreshNotifications"
          class="p-2 text-gray-600 hover:text-pink-600 hover:bg-pink-50 rounded-lg transition-colors"
          :disabled="loading"
        >
          <svg class="w-4 h-4" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading && (!notifications || notifications.length === 0)" class="flex justify-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-pink-500"></div>
    </div>

    <!-- 空状态 -->
    <div v-else-if="!loading && (!notifications || notifications.length === 0)" class="text-center py-12">
      <div class="w-16 h-16 mx-auto mb-4 text-gray-300">
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
        </svg>
      </div>
      <p class="text-gray-500 text-lg font-medium">暂无消息通知</p>
      <p class="text-gray-400 text-sm mt-1">当有人点赞、评论或关注你时，消息会显示在这里</p>
    </div>

    <!-- 通知列表 -->
    <div v-else class="space-y-2">
      <div 
        v-for="notification in notifications" 
        :key="notification.id"
        class="notification-item group relative p-4 bg-white border border-gray-200 rounded-lg hover:shadow-md transition-all duration-200 cursor-pointer"
        :class="{
          'bg-blue-50 border-blue-200': !notification.is_read,
          'bg-white border-gray-200': notification.is_read
        }"
        @click="handleNotificationClick(notification)"
      >
        <!-- 未读指示器 -->
        <div 
          v-if="!notification.is_read"
          class="absolute top-4 left-2 w-2 h-2 bg-pink-500 rounded-full"
        ></div>

        <div class="flex items-start space-x-3 ml-4">
          <!-- 用户头像 -->
          <div class="flex-shrink-0">
            <img 
              v-if="notification.actor_avatar"
              :src="notification.actor_avatar" 
              :alt="notification.actor_nickname || notification.actor_username"
              class="w-10 h-10 rounded-full object-cover border-2 border-gray-100"
            >
            <div 
              v-else
              class="w-10 h-10 bg-gradient-to-br from-pink-400 to-rose-400 rounded-full flex items-center justify-center text-white font-semibold text-sm"
            >
              {{ (notification.actor_nickname || notification.actor_username || 'U').charAt(0).toUpperCase() }}
            </div>
          </div>

          <!-- 通知内容 -->
          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <!-- 通知类型图标和消息 -->
                <div class="flex items-start space-x-2">
                  <span class="text-lg">{{ getNotificationIcon(notification.type) }}</span>
                  <div>
                    <p class="text-sm text-gray-800 leading-relaxed">
                      <span class="font-medium text-pink-600">
                        {{ notification.actor_nickname || notification.actor_username }}
                      </span>
                      {{ notification.message }}
                    </p>
                    <div class="flex items-center space-x-4 mt-2">
                      <span class="text-xs text-gray-500">
                        {{ formatNotificationTime(notification.created_at) }}
                      </span>
                      <span class="text-xs text-gray-400">
                        {{ notificationTypeMap[notification.type] }}
                      </span>
                    </div>
                  </div>
                </div>

                <!-- 文章缩略图 -->
                <div 
                  v-if="notification.article_title && (notification.type === 'like' || notification.type === 'comment')"
                  class="mt-3 p-2 bg-gray-50 rounded-lg"
                >
                  <div class="flex items-center space-x-2">
                    <img 
                      v-if="notification.article_cover"
                      :src="notification.article_cover"
                      :alt="notification.article_title"
                      class="w-12 h-8 object-cover rounded"
                    >
                    <p class="text-sm text-gray-600 line-clamp-2 flex-1">
                      {{ notification.article_title }}
                    </p>
                  </div>
                </div>
              </div>

              <!-- 操作按钮 -->
              <div class="flex items-center space-x-1 ml-4 opacity-0 group-hover:opacity-100 transition-opacity">
                <button
                  v-if="!notification.is_read"
                  @click.stop="markAsRead([notification.id])"
                  class="p-1 text-gray-400 hover:text-pink-600 transition-colors"
                  title="标记为已读"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </button>
                <button
                  @click.stop="deleteNotification(notification.id)"
                  class="p-1 text-gray-400 hover:text-red-600 transition-colors"
                  title="删除通知"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 加载更多 -->
    <div v-if="hasMore && !loading" class="text-center py-6">
      <button 
        @click="loadMore"
        class="px-6 py-2 text-sm text-pink-600 border border-pink-200 rounded-lg hover:bg-pink-50 transition-colors"
      >
        加载更多
      </button>
    </div>

    <!-- 底部加载状态 -->
    <div v-if="loading && notifications && notifications.length > 0" class="text-center py-4">
      <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-pink-500 mx-auto"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { NotificationApi, type Notification, type NotificationStats, formatNotificationTime, notificationTypeMap, notificationIconMap } from '@/api/notification'

// Props
interface Props {
  autoRefresh?: boolean
  refreshInterval?: number
}

const props = withDefaults(defineProps<Props>(), {
  autoRefresh: false,
  refreshInterval: 30000 // 30秒
})

// 路由
const router = useRouter()

// 数据状态
const notifications = ref<Notification[]>([])
const stats = ref<NotificationStats>()
const loading = ref(false)
const clearing = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const hasMore = ref(true)

// 获取通知图标
const getNotificationIcon = (type: string) => {
  return notificationIconMap[type as keyof typeof notificationIconMap] || '📢'
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
const loadNotifications = async (reset = false) => {
  loading.value = true
  try {
    if (reset) {
      currentPage.value = 1
      notifications.value = []
    }

    const response = await NotificationApi.getNotifications({
      page: currentPage.value,
      limit: pageSize.value
    })

    if (response.code === 200) {
      const { notifications: newNotifications, pagination } = response.data
      
      if (reset) {
        notifications.value = newNotifications
      } else {
        notifications.value?.push(...newNotifications)
      }
      
      hasMore.value = currentPage.value < pagination.total_pages
    }
  } catch (error) {
    console.error('加载通知列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 刷新通知
const refreshNotifications = async () => {
  await Promise.all([
    loadStats(),
    loadNotifications(true)
  ])
}

// 加载更多
const loadMore = async () => {
  if (loading.value || !hasMore.value) return
  currentPage.value++
  await loadNotifications()
}

// 标记为已读
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

// 一键清除所有通知
const clearAllNotifications = async () => {
  if (!confirm('确定要删除所有通知吗？此操作不可恢复。')) {
    return
  }

  clearing.value = true
  try {
    await NotificationApi.deleteAllNotifications()
    
    // 清空本地状态
    notifications.value = []
    stats.value = {
      unread_count: 0,
      total_count: 0
    }
    
    // 重新加载数据
    await Promise.all([
      loadStats(),
      loadNotifications(true)
    ])
  } catch (error) {
    console.error('清除所有通知失败:', error)
  } finally {
    clearing.value = false
  }
}

// 处理通知点击
const handleNotificationClick = async (notification: Notification) => {
  console.log('通知被点击:', notification)
  
  try {
    // 如果是评论通知，跳转到文章并聚焦评论区
    if (notification.type === 'comment' && notification.resource_id) {
      console.log('评论通知，准备跳转到文章:', notification.resource_id)
      
      // 先标记为已读
      if (!notification.is_read) {
        console.log('标记通知为已读')
        await markAsRead([notification.id])
      }
      
      // 跳转到文章页面，并通过URL参数指示聚焦评论区
      console.log('执行路由跳转')
      await router.push({
        path: `/articles/${notification.resource_id}`,
        query: { focus: 'comments' }
      })
      console.log('路由跳转完成')
      
    } else if (notification.type === 'like' && notification.resource_id) {
      // 点赞通知跳转到文章
      console.log('点赞通知，跳转到文章:', notification.resource_id)
      if (!notification.is_read) {
        await markAsRead([notification.id])
      }
      await router.push(`/articles/${notification.resource_id}`)
      
    } else if (notification.type === 'follow') {
      // 关注通知跳转到用户页面
      console.log('关注通知，跳转到用户页面:', notification.actor_id)
      if (!notification.is_read) {
        await markAsRead([notification.id])
      }
      await router.push(`/users/${notification.actor_id}`)

    } else if (notification.type === 'message') {
      // 私信通知跳转到消息页面
      console.log('私信通知，跳转到消息页面，发送者:', notification.actor_id)
      if (!notification.is_read) {
        await markAsRead([notification.id])
      }
      // 跳转到消息页面，并且携带用户ID参数以便打开对应的对话
      await router.push({
        path: '/messages',
        query: { user: notification.actor_id }
      })
    }
  } catch (error) {
    console.error('处理通知点击失败:', error)
  }
}

// 删除通知
const deleteNotification = async (notificationId: number) => {
  try {
    await NotificationApi.deleteNotification(notificationId)
    
    // 从列表中移除
    const index = notifications.value?.findIndex(n => n.id === notificationId)
    if (index !== -1 && notifications.value) {
      notifications.value.splice(index, 1)
    }
    
    // 更新统计
    await loadStats()
  } catch (error) {
    console.error('删除通知失败:', error)
  }
}

// 自动刷新
let refreshTimer: number | null = null

const startAutoRefresh = () => {
  if (!props.autoRefresh) return
  
  refreshTimer = window.setInterval(() => {
    loadStats() // 只更新统计，不刷新整个列表
  }, props.refreshInterval)
}

const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

// 组件挂载
onMounted(() => {
  refreshNotifications()
  startAutoRefresh()
})

// 组件卸载
import { onUnmounted } from 'vue'
onUnmounted(() => {
  stopAutoRefresh()
})

// 暴露方法给父组件
defineExpose({
  refreshNotifications,
  loadStats,
  markAllAsRead
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