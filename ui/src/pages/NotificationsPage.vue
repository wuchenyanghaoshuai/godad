<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 导航栏 -->
    <BaseHeader :showNotifications="true" />

    <div class="max-w-7xl mx-auto px-4 py-6">
      <!-- 页面标题 -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-gray-900">消息中心</h1>
        <p class="text-gray-600 mt-1">查看所有互动消息和通知</p>
      </div>

      <!-- 水平布局容器 -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden flex h-[calc(100vh-var(--header-h))] horizontal-layout">
        <!-- 左侧面板：通知列表 -->
        <div class="w-2/5 min-w-[320px] border-r border-gray-200 flex flex-col left-panel">
          <!-- 消息头部 -->
            <div class="p-4 border-b border-gray-200 sticky top-0 z-10 bg-white">
            <div class="flex items-center justify-between">
              <h2 class="font-semibold text-gray-900 flex items-center space-x-2">
                <svg class="w-5 h-5 text-pink-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5 5v-5zM4 19h5v-5l-5 5zm0-13h5V1L4 6zm11-5v5h5l-5-5z"/>
                </svg>
                <span>互动消息</span>
                <span v-if="unreadMessagesCount > 0" class="bg-red-500 text-white text-xs rounded-full px-2 py-0.5">
                  {{ unreadMessagesCount }}
                </span>
              </h2>
              <div v-if="allNotifications && allNotifications.length > 0" class="flex items-center space-x-3">
                <button
                  @click="markAllMessagesAsRead"
                  class="text-sm text-pink-600 hover:text-pink-700 disabled:opacity-50 disabled:cursor-not-allowed"
                  :disabled="isMarkingAll"
                >
                  {{ isMarkingAll ? '标记中…' : '全部标记为已读' }}
                </button>
                <button
                  @click="clearAllNotifications"
                  class="text-sm text-red-600 hover:text-red-700 disabled:opacity-50 disabled:cursor-not-allowed"
                  :disabled="isClearingAll"
                >
                  {{ isClearingAll ? '清空中…' : '清空所有消息' }}
                </button>
              </div>
            </div>
            <!-- 错误提示条 -->
            <div v-if="notificationsError" class="px-4 py-2 bg-red-50 text-red-600 text-sm flex items-center justify-between">
              <span class="truncate">{{ notificationsError }}</span>
              <button
                class="p-1 rounded hover:bg-red-100 disabled:opacity-50 disabled:cursor-not-allowed"
                @click="fetchNotifications(1)"
                :disabled="notificationsLoading"
                title="重试"
                aria-label="重试"
              >
                <svg v-if="!notificationsLoading" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                <svg v-else class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path class="opacity-25" stroke-width="4" d="M12 4v0" />
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke-width="4"></circle>
                </svg>
              </button>
            </div>
            <!-- 筛选页签 -->
            <div class="mt-3 flex items-center gap-2 text-sm">
              <button
                v-for="tab in tabs"
                :key="tab.key"
                @click="activeTab = tab.key"
                class="px-3 py-1 rounded-full border transition-colors"
                :class="activeTab === tab.key
                  ? 'bg-pink-50 border-pink-200 text-pink-700'
                  : 'bg-white border-gray-200 text-gray-600 hover:bg-gray-50'"
              >
                {{ tab.label }}
                <span v-if="tab.count > 0" class="ml-1 text-xs text-gray-400">{{ tab.count }}</span>
              </button>
            </div>
          </div>

          <!-- 所有消息列表 -->
          <div class="flex-1 overflow-y-auto">
            <!-- 通知/私信列表 -->
            <div
              v-for="notification in displayNotifications"
              :key="notification.id"
              :class="[
                'p-4 border-b border-gray-100 hover:bg-gray-50 cursor-pointer notification-item relative group',
                !notification.is_read ? 'bg-blue-50' : '',
                selectedNotification?.id === notification.id ? 'notification-selected' : ''
              ]"
              @click="handleNotificationClick(notification)"
            >
                  <div class="flex items-start space-x-3">
                <!-- 头像 -->
                <div class="flex-shrink-0">
                  <img
                    v-if="notification.actor_avatar"
                    :src="notification.actor_avatar"
                    :alt="notification.actor_nickname || notification.actor_username"
                    class="w-9 h-9 rounded-full object-cover border"
                  />
                  <div v-else class="w-9 h-9 rounded-full bg-gradient-to-br from-pink-400 to-rose-400 text-white flex items-center justify-center text-sm font-semibold">
                    {{ (notification.actor_nickname || notification.actor_username || 'U').charAt(0).toUpperCase() }}
                  </div>
                </div>
                <!-- 内容 -->
                <div class="flex-1 min-w-0">
                  <div class="flex items-start justify-between">
                    <div class="flex-1 pr-2">
                      <!-- 标题行 -->
                      <p class="text-sm text-gray-900 font-medium">
                        {{ getNotificationTitle(notification) }}
                      </p>
                      <!-- 次行：摘要 + 元信息 -->
                      <div class="mt-1">
                        <p class="text-xs text-gray-600 line-clamp-1">{{ getNotificationSummary(notification) }}</p>
                        <div class="flex items-center gap-3 mt-1">
                          <span class="text-xs text-gray-400">{{ formatNotificationTime(notification.created_at) }}</span>
                          <span class="text-xs text-gray-400">{{ notificationTypeMap[notification.type] || '通知' }}</span>
                          <!-- 私信分组计数徽标 -->
                          <span v-if="notification.type === 'message' && messageGroupMeta[notification.id]" class="text-xs text-gray-500 bg-gray-100 px-1.5 py-0.5 rounded">
                            {{ messageGroupMeta[notification.id].count }} 条<span v-if="messageGroupMeta[notification.id].unread > 0">，未读 {{ messageGroupMeta[notification.id].unread }}</span>
                          </span>
                        </div>
                      </div>
                    </div>
                    <div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button
                        v-if="!notification.is_read"
                        @click.stop="markOneAsRead(notification)"
                        class="p-1 text-gray-400 hover:text-green-600"
                        title="标记为已读"
                      >
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                        </svg>
                      </button>
                      <button
                        @click.stop="deleteNotification(notification)"
                        class="p-1 text-gray-400 hover:text-red-600"
                        :title="notification.type === 'message' ? '删除私信' : '删除通知'"
                        :disabled="deletingIds.has(notification.id)"
                      >
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                        </svg>
                      </button>
                      <span v-if="!notification.is_read" class="w-2 h-2 bg-blue-500 rounded-full"></span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div v-if="!allNotifications || allNotifications.length === 0" class="p-8 text-center text-gray-500">
              <svg class="w-12 h-12 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5 5v-5zM4 19h5v-5l-5 5zm0-13h5V1L4 6zm11-5v5h5l-5-5z"/>
              </svg>
              <p>暂无互动消息</p>
            </div>
            <!-- 加载更多 -->
            <div v-if="hasMore && !notificationsLoading" class="p-4 text-center border-t border-gray-100">
              <button
                @click="loadMoreNotifications"
                class="px-4 py-2 text-sm text-pink-600 border border-pink-200 rounded-lg hover:bg-pink-50 transition-colors"
              >
                加载更多
              </button>
            </div>
          </div>
        </div>

        <!-- 右侧面板：通知详情/聊天界面 -->
        <div class="w-3/5 flex flex-col right-panel detail-panel">
          <!-- 无选中通知时的占位符 -->
          <div v-if="!selectedNotification" class="h-full flex items-center justify-center text-gray-500">
            <div class="text-center">
              <svg class="w-16 h-16 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5 5v-5zM4 19h5v-5l-5 5zm0-13h5V1L4 6zm11-5v5h5l-5-5z"/>
              </svg>
              <p class="text-lg font-medium mb-2">选择一个通知查看详情</p>
              <p class="text-sm">点击左侧的通知项目查看详细信息或进行私信对话</p>
            </div>
          </div>

          <!-- 私信界面 -->
          <div v-else-if="selectedNotification.type === 'message'" class="flex flex-col h-full">
            <!-- 聊天头部 -->
            <div class="p-4 border-b border-gray-200 bg-gray-50">
              <div class="flex items-center justify-between">
                <div class="flex items-center space-x-3">
                  <span class="text-xl">💌</span>
                  <div>
                    <h3 class="font-semibold text-gray-900">
                      与 {{ selectedNotification.actor_nickname || selectedNotification.actor_username }} 的对话
                    </h3>
                    <p class="text-sm text-gray-500">私信</p>
                  </div>
                </div>
                <button
                  @click="selectedNotification = null"
                  class="text-gray-400 hover:text-gray-600 p-1 rounded-full hover:bg-gray-200"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                  </svg>
                </button>
              </div>
            </div>

            <!-- 聊天消息列表 -->
            <div class="flex-1" ref="chatContainer">
              <ChatMessageList
                v-if="currentConversation"
                ref="chatMessageListRef"
                :conversation="currentConversation"
                @message-sent="handleMessageSent"
                class="h-full"
              />
              <div v-else class="flex items-center justify-center h-full text-gray-500">
                <div class="text-center">
                  <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-pink-500 mx-auto mb-2"></div>
                  <p>正在加载对话...</p>
                </div>
              </div>
            </div>

            <!-- 聊天输入框（吸底） -->
            <div v-if="currentConversation" class="sticky bottom-0 z-10 bg-white border-t">
              <ChatInputBox
                :conversation="currentConversation"
                @message-sent="handleMessageSent"
              />
            </div>
          </div>

          <!-- 其他通知的详情显示 -->
          <div v-else class="flex flex-col h-full">
            <!-- 详情头部 -->
            <div class="p-4 border-b border-gray-200 bg-gray-50">
              <div class="flex items-center justify-between">
                <div class="flex items-center space-x-3">
                  <span class="text-xl">
                    {{ notificationIconMap[selectedNotification.type] || '📢' }}
                  </span>
                  <div>
                    <h3 class="font-semibold text-gray-900">
                      {{ notificationTypeMap[selectedNotification.type] || '通知' }}
                    </h3>
                    <p class="text-sm text-gray-500">
                      {{ formatNotificationTime(selectedNotification.created_at) }}
                    </p>
                  </div>
                </div>
                <button
                  @click="selectedNotification = null"
                  class="text-gray-400 hover:text-gray-600 p-1 rounded-full hover:bg-gray-200"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                  </svg>
                </button>
              </div>
            </div>

            <!-- 详情内容 -->
            <div class="flex-1 overflow-y-auto p-4">
              <div class="bg-white rounded-md border border-gray-200 p-4">
                <div class="flex items-start space-x-3">
                  <div class="flex-shrink-0">
                    <div class="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center">
                      <span class="text-sm">{{ notificationIconMap[selectedNotification.type] || '📢' }}</span>
                    </div>
                  </div>
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center justify-between mb-2">
                      <h4 class="text-base font-medium text-gray-900">
                        {{ selectedNotification.actor_nickname || selectedNotification.actor_username }}
                      </h4>
                      <span class="text-sm text-gray-500">
                        {{ formatNotificationTime(selectedNotification.created_at) }}
                      </span>
                    </div>
                    <p class="text-gray-700 mb-4">{{ selectedNotification.message }}</p>

                    <!-- 操作按钮区域 -->
                    <div class="flex space-x-2">
                      <!-- 点赞通知按钮 -->
                      <template v-if="selectedNotification.type === 'like'">
                        <button
                          v-if="selectedNotification.resource_id"
                          @click="viewArticle(selectedNotification.resource_id)"
                          class="bg-pink-600 text-white px-3 py-2 text-sm rounded hover:bg-pink-700 transition-colors"
                        >
                          查看文章
                        </button>
                      </template>

                      <!-- 评论通知按钮 -->
                      <template v-else-if="selectedNotification.type === 'comment'">
                        <button
                          v-if="selectedNotification.resource_id"
                          @click="router.push({ path: `/articles/${selectedNotification.resource_id}` , query: { focus: 'comments' } })"
                          class="bg-pink-600 text-white px-3 py-2 text-sm rounded hover:bg-pink-700 transition-colors"
                        >
                          查看并回复
                        </button>
                        <button
                          v-if="selectedNotification.resource_id"
                          @click="viewArticle(selectedNotification.resource_id)"
                          class="bg-gray-500 text-white px-3 py-2 text-sm rounded hover:bg-gray-600 transition-colors"
                        >
                          查看文章
                        </button>
                      </template>

                      <!-- 关注通知按钮 -->
                      <template v-else-if="selectedNotification.type === 'follow'">
                        <button
                          @click="viewUserProfile(selectedNotification)"
                          class="bg-blue-600 text-white px-3 py-2 text-sm rounded hover:bg-blue-700 transition-colors"
                        >
                          查看用户资料
                        </button>
                        <button
                          @click="startPrivateChat(selectedNotification)"
                          class="bg-green-600 text-white px-3 py-2 text-sm rounded hover:bg-green-700 transition-colors"
                        >
                          发送私信
                        </button>
                      </template>

                      <!-- 收藏通知按钮 -->
                      <template v-else-if="selectedNotification.type === 'bookmark'">
                        <button
                          v-if="selectedNotification.resource_id"
                          @click="viewArticle(selectedNotification.resource_id)"
                          class="bg-green-600 text-white px-3 py-2 text-sm rounded hover:bg-green-700 transition-colors"
                        >
                          查看文章
                        </button>
                      </template>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { NotificationApi, type Notification, formatNotificationTime, notificationIconMap, notificationTypeMap } from '../api/notification'
import { ChatAPI, type ConversationResponse } from '@/api'
import { useToast } from '@/composables/useToast'
import BaseHeader from '@/components/BaseHeader.vue'
import ChatMessageList from '@/components/ChatMessageList.vue'
import ChatInputBox from '@/components/ChatInputBox.vue'
import { useNotificationSync } from '@/composables/useNotificationSync'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { showToast } = useToast()

// 当前选中的通知
const selectedNotification = ref<Notification | null>(null)

// 所有未读消息数量
const unreadMessagesCount = ref(0)

// 聊天相关状态
const currentConversation = ref<ConversationResponse | null>(null)
const chatContainer = ref<HTMLElement>()
const chatMessageListRef = ref<any>()
const { triggerRefresh, onNotificationEvent } = useNotificationSync()
let autoRefreshTimer: number | null = null

// 通知相关数据
const notifications = ref<Notification[]>([])
const isMarkingAll = ref(false)
const isClearingAll = ref(false)
const deletingIds = ref<Set<number>>(new Set())
const notificationsLoading = ref(false)
const notificationsError = ref('')
const unreadNotificationsCount = ref(0)

// 分页相关
const notificationsPage = ref(1)
const notificationsTotalPages = ref(1)
const limit = ref(10)

// 筛选页签
type TabKey = 'all' | 'notify' | 'message' | 'unread'
const activeTab = ref<TabKey>('all')
const tabs = computed(() => [
  { key: 'all' as TabKey, label: '全部', count: notifications.value.length },
  { key: 'notify' as TabKey, label: '通知', count: notifications.value.filter(n => n.type !== 'message').length },
  { key: 'message' as TabKey, label: '私信', count: notifications.value.filter(n => n.type === 'message').length },
  { key: 'unread' as TabKey, label: '未读', count: notifications.value.filter(n => !n.is_read).length },
])

// 计算属性 - 所有通知，按时间倒序
const allNotifications = computed(() => {
  if (!notifications.value) return []
  const filtered = notifications.value.filter(n => ['like', 'comment', 'message', 'follow', 'bookmark', 'system'].includes(n.type))
  return filtered.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
})

// 对 message 做聚合（按 actor_id-receiver_id）
const groupMessageNotifications = (list: Notification[]) => {
  const messageGroups = new Map<string, Notification[]>()
  const others: Notification[] = []
  for (const n of list) {
    if (n.type === 'message') {
      const key = `${n.actor_id}-${n.receiver_id}`
      if (!messageGroups.has(key)) messageGroups.set(key, [])
      messageGroups.get(key)!.push(n)
    } else {
      others.push(n)
    }
  }
  const grouped: Notification[] = []
  const meta: Record<number, { count: number; unread: number }> = {}
  for (const [_key, group] of messageGroups) {
    const latest = group.reduce((a, b) => new Date(a.created_at) > new Date(b.created_at) ? a : b)
    const count = group.length
    const unread = group.filter(n => !n.is_read).length
    meta[latest.id] = { count, unread }
    grouped.push(latest)
  }
  return { grouped, others, meta }
}

// 分类过滤（来自路由的 category 参数）
const categoryFilter = ref<string>('')

// 当前展示的数据源
const displayNotifications = computed(() => {
  const list = allNotifications.value
  const { grouped, others } = groupMessageNotifications(list)
  const filterOthersByCategory = (items: Notification[]) => {
    if (!categoryFilter.value) return items
    const allowed = ['like', 'comment', 'follow', 'bookmark', 'system', 'mention']
    if (!allowed.includes(categoryFilter.value)) return items
    return items.filter(n => n.type === (categoryFilter.value as any))
  }
  switch (activeTab.value) {
    case 'notify':
      return filterOthersByCategory(others)
    case 'message':
      return grouped.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    case 'unread':
      return list.filter(n => !n.is_read)
    case 'all':
    default:
      // 全部：非私信 + 聚合后的私信，按时间倒序
      return [...others, ...grouped].sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
  }
})

// 分组元信息映射（key: 通知id, value: 分组数量/未读数量）
const messageGroupMeta = computed(() => {
  const list = allNotifications.value
  return groupMessageNotifications(list).meta
})

// 获取通知列表
const fetchNotifications = async (page = 1, append = false) => {
  notificationsLoading.value = true
  notificationsError.value = ''

  try {
    const response = await NotificationApi.getNotificationsPage({
      page,
      limit: limit.value
    })

    const pageData = response.data
    if (append && notifications.value?.length) {
      const map = new Map<number, any>()
      ;[...notifications.value, ...pageData.items].forEach((n: any) => map.set(n.id, n))
      notifications.value = Array.from(map.values()) as any
    } else {
      notifications.value = pageData.items as any
    }
    notificationsPage.value = pageData.page
    notificationsTotalPages.value = pageData.total_pages

    // 统计未读通知数量
    unreadNotificationsCount.value = notifications.value ? notifications.value.filter(n => !n.is_read).length : 0
    unreadMessagesCount.value = unreadNotificationsCount.value
  } catch (err: any) {
    notificationsError.value = err.response?.data?.error || '获取通知列表失败'
  } finally {
    notificationsLoading.value = false
  }
}

const hasMore = computed(() => notificationsPage.value < notificationsTotalPages.value)

const loadMoreNotifications = async () => {
  if (notificationsLoading.value || !hasMore.value) return
  await fetchNotifications(notificationsPage.value + 1, true)
}

// 图标/标题统一使用 notificationIconMap / notificationTypeMap

// 处理通知点击
const handleNotificationClick = async (notification: Notification) => {
  // 如果通知未读，先标记为已读
  if (!notification.is_read) {
    try {
      await NotificationApi.markAsRead([notification.id])
      notification.is_read = true
      unreadNotificationsCount.value = Math.max(0, unreadNotificationsCount.value - 1)
      unreadMessagesCount.value = Math.max(0, unreadMessagesCount.value - 1)
      triggerRefresh()
    } catch (error) {
      console.error('标记通知已读失败:', error)
    }
  }

  selectedNotification.value = notification

  // 如果是私信通知，优先按会话ID(resource_id) 打开；无则回退按用户ID
  if (notification.type === 'message') {
    if (notification.resource_id && !isNaN(Number(notification.resource_id))) {
      await loadConversationByConversationId(Number(notification.resource_id))
      if (!currentConversation.value && notification.actor_id) {
        await loadConversationByUserId(notification.actor_id)
      }
    } else if (notification.actor_id) {
      await loadConversationByUserId(notification.actor_id)
    } else {
      currentConversation.value = null
    }
  } else {
    // 非私信通知，清空对话
    currentConversation.value = null
  }
}

// 快速标记单条为已读
const markOneAsRead = async (notification: Notification) => {
  try {
    if (notification.is_read) return
    await NotificationApi.markAsRead([notification.id])
    notification.is_read = true
    unreadNotificationsCount.value = Math.max(0, unreadNotificationsCount.value - 1)
    unreadMessagesCount.value = Math.max(0, unreadMessagesCount.value - 1)
    triggerRefresh()
  } catch (e) {
    // 静默
  }
}

// 查看文章
const viewArticle = (articleId: number) => {
  router.push(`/articles/${articleId}`)
}

const viewUserProfile = (notification: Notification) => {
  const username = notification.actor_username
  if (username) {
    router.push(`/users/${username}`)
  } else {
    showToast('无法获取用户信息', 'error')
  }
}

// 左侧列表：标题与摘要
const getNotificationTitle = (n: Notification): string => {
  if (n.type === 'system') {
    return n.title && n.title.trim() ? n.title.trim() : '系统通知'
  }
  if (n.type === 'message') {
    const name = n.actor_nickname || n.actor_username || '私信'
    return `私信 · ${name}`
  }
  const typeName = notificationTypeMap[n.type] || '通知'
  const name = n.actor_nickname || n.actor_username || ''
  return name ? `${typeName} · ${name}` : typeName
}

const getNotificationSummary = (n: Notification): string => {
  const text = (n.message || '').trim()
  // 限制长度，避免泄露过多详情
  const max = 80
  return text.length > max ? text.slice(0, max) + '…' : text
}

// 标记所有通知为已读
const markAllMessagesAsRead = async () => {
  try {
    isMarkingAll.value = true
    await NotificationApi.markAllAsRead()
    notifications.value.forEach(n => n.is_read = true)
    unreadNotificationsCount.value = 0
    unreadMessagesCount.value = 0
    triggerRefresh()
  } catch (error) {
    console.error('标记所有通知已读失败:', error)
  } finally {
    isMarkingAll.value = false
  }
}

// 时间格式统一使用 formatNotificationTime

// 删除单个通知
const deleteNotification = async (notification: Notification) => {
  if (!authStore.user) {
    showToast('请先登录', 'error')
    return
  }

  // 确认删除
  const confirmed = confirm('确定要删除这条通知吗？')
  if (!confirmed) {
    return
  }

  try {
    deletingIds.value.add(notification.id)
    // 立即从UI中移除该通知
    notifications.value = notifications.value.filter(n => n.id !== notification.id)

    // 如果当前选中的是被删除的通知，清除选中状态
    if (selectedNotification.value?.id === notification.id) {
      selectedNotification.value = null
    }

    // 更新未读数量
    if (!notification.is_read) {
      unreadNotificationsCount.value = Math.max(0, unreadNotificationsCount.value - 1)
      unreadMessagesCount.value = Math.max(0, unreadMessagesCount.value - 1)
    }

    // 调用删除通知API
    await NotificationApi.deleteNotification(notification.id)

    showToast('通知已删除', 'success')
    triggerRefresh()

  } catch (error: any) {
    console.error('删除通知失败:', error)
    showToast(error.response?.data?.message || '删除通知失败', 'error')

    // 如果删除失败，重新获取通知列表恢复UI状态
    await fetchNotifications()
  } finally {
    deletingIds.value.delete(notification.id)
  }
}

// 清空所有通知
const clearAllNotifications = async () => {
  if (!authStore.user) {
    showToast('请先登录', 'error')
    return
  }

  // 确认删除
  const confirmed = confirm('确定要清空所有消息吗？此操作不可恢复。')
  if (!confirmed) {
    return
  }

  try {
    isClearingAll.value = true
    // 立即清空UI中的通知
    notifications.value = []
    selectedNotification.value = null
    unreadNotificationsCount.value = 0
    unreadMessagesCount.value = 0
    notificationsPage.value = 1
    notificationsTotalPages.value = 1

    // 调用删除所有通知API
    await NotificationApi.deleteAllNotifications()

    showToast('所有消息已清空', 'success')
    triggerRefresh()

  } catch (error: any) {
    console.error('清空所有通知失败:', error)
    showToast(error.response?.data?.message || '清空消息失败', 'error')

    // 如果删除失败，重新获取通知列表恢复UI状态
    await fetchNotifications()
  } finally {
    isClearingAll.value = false
  }
}

// 根据actor_id加载对话
const loadConversationByUserId = async (actorId: number) => {
  if (!authStore.user) return

  try {
    // 查找或创建与该用户的对话
    const response = await ChatAPI.getOrCreateConversation({
      other_user_id: actorId
    })

    currentConversation.value = response.data
  } catch (error: any) {
    console.error('加载对话失败:', error)
    showToast('加载对话失败', 'error')
  }
}

// 根据会话ID加载对话（通过会话列表查找）
const loadConversationByConversationId = async (conversationId: number) => {
  if (!authStore.user) return
  try {
    const resp = await ChatAPI.getConversations({ page: 1, limit: 100 })
    const found = resp.data.conversations.find(c => c.id === conversationId)
    if (found) {
      currentConversation.value = found
    } else {
      currentConversation.value = null
    }
  } catch (error: any) {
    console.error('根据会话ID加载对话失败:', error)
    currentConversation.value = null
  }
}

// 处理消息发送
const handleMessageSent = async (message: any) => {
  // 立即将新消息添加到聊天列表中（若组件可用）
  try {
    chatMessageListRef.value?.addMessage?.(message)
    chatMessageListRef.value?.scrollToBottomAnimated?.()
  } catch (e) {
    console.debug('即时添加消息失败，稍后刷新修复。', e)
  }

  // 无论如何，稍后强制刷新一次，确保与服务端状态一致
  setTimeout(() => {
    chatMessageListRef.value?.refreshMessages?.()
  }, 1000)
}

// 开始私信聊天（从关注通知等地方触发）
const startPrivateChat = async (notification: Notification) => {
  if (!notification.actor_id) {
    showToast('无法获取用户信息', 'error')
    return
  }

  await startPrivateChatByUserId(notification.actor_id)
}

// 根据用户ID开始私信聊天
const startPrivateChatByUserId = async (userId: number) => {
  try {
    // 加载或创建对话
    await loadConversationByUserId(userId)

    if (!currentConversation.value) {
      showToast('无法创建对话', 'error')
      return
    }

    // 创建一个虚拟的私信通知来显示聊天界面
    selectedNotification.value = {
      id: 0,
      type: 'message',
      message: '开始新对话',
      actor_id: userId,
      actor_username: currentConversation.value.other_user?.username || '',
      actor_nickname: currentConversation.value.other_user?.nickname || '',
      actor_avatar: currentConversation.value.other_user?.avatar || '',
      resource_id: null,
      resource_type: null,
      is_read: true,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }

    // 清除URL参数
    if (route.query.user) {
      router.replace({ query: {} })
    }

    // 已成功开始对话
  } catch (error: any) {
    console.error('开始私信聊天失败:', error)
    showToast('开始对话失败', 'error')
  }
}

// 删除当前对话
const deleteCurrentConversation = async () => {
  if (!currentConversation.value) {
    return
  }

  const otherUserName = currentConversation.value.other_user?.nickname || currentConversation.value.other_user?.username || '对方'

  // 确认删除
  const confirmed = confirm(`确定要删除与 ${otherUserName} 的对话吗？此操作不可恢复。`)
  if (!confirmed) {
    return
  }

  try {
    // 调用删除对话API
    await ChatAPI.deleteConversation(currentConversation.value.id)

    // 清空当前对话状态
    currentConversation.value = null
    selectedNotification.value = null

    // 清除保存的聊天状态
    localStorage.removeItem('godad-chat-state')

    showToast(`已删除与 ${otherUserName} 的对话`, 'success')
  } catch (error: any) {
    console.error('删除对话失败:', error)
    showToast(error.response?.data?.message || '删除对话失败', 'error')
  }
}

// 监听路由查询参数变化
watch(() => route.query.user, async (newUserId) => {
  if (newUserId && !isNaN(Number(newUserId))) {
    await startPrivateChatByUserId(Number(newUserId))
  }
})

// 监听登录状态变化
watch(() => authStore.isAuthenticated, (isAuth) => {
  if (!isAuth) {
    router.push('/login')
  } else {
    fetchNotifications()
  }
})

// 保存对话状态到localStorage
const saveChatState = () => {
  if (selectedNotification.value && currentConversation.value) {
    const chatState = {
      conversationId: currentConversation.value.id,
      otherUserId: currentConversation.value.other_user?.id,
      timestamp: Date.now()
    }
    localStorage.setItem('godad-chat-state', JSON.stringify(chatState))
  } else {
    localStorage.removeItem('godad-chat-state')
  }
}

// 从localStorage恢复对话状态
const restoreChatState = async () => {
  const savedState = localStorage.getItem('godad-chat-state')
  if (!savedState) return false

  try {
    const chatState = JSON.parse(savedState)

    // 检查状态是否过期（24小时）
    if (Date.now() - chatState.timestamp > 24 * 60 * 60 * 1000) {
      localStorage.removeItem('godad-chat-state')
      return false
    }

    // 恢复对话：按 otherUserId 重新获取/创建对话
    if (chatState.otherUserId) {
      try {
        await startPrivateChatByUserId(chatState.otherUserId)
        return true
      } catch (error) {
        localStorage.removeItem('godad-chat-state')
        return false
      }
    }
  } catch (error) {
    console.error('恢复聊天状态失败:', error)
    localStorage.removeItem('godad-chat-state')
  }

  return false
}

// 监听状态变化，自动保存
watch([selectedNotification, currentConversation], () => {
  saveChatState()
}, { deep: true })

// 组件挂载时加载数据
onMounted(async () => {
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }

  // 解析路由查询参数，设置初始筛选
  const initFromQuery = () => {
    const tab = String(route.query.tab || '').toLowerCase()
    const cat = String(route.query.category || '').toLowerCase()
    const tabKeys = ['all', 'notify', 'message', 'unread']
    if (tab && tabKeys.includes(tab)) {
      activeTab.value = tab as any
    }
    if (cat) {
      categoryFilter.value = cat
      // 有明确分类时，默认切到通知大类
      if (activeTab.value !== 'message' && activeTab.value !== 'unread') {
        activeTab.value = 'notify'
      }
    }
  }
  initFromQuery()

  await fetchNotifications()

  // 先尝试恢复聊天状态
  const restored = await restoreChatState()

  // 如果没有恢复成功，检查是否有user参数
  if (!restored) {
    const userId = route.query.user
    if (userId && !isNaN(Number(userId))) {
      await startPrivateChatByUserId(Number(userId))
    }
  }

  // 监听来自其他组件的刷新事件（如标记已读/删除）
  onNotificationEvent('refresh', () => {
    fetchNotifications(notificationsPage.value)
  })

  // 页面可见时的轻量自动刷新（每 20s 刷新当前页）
  const startAutoRefresh = () => {
    stopAutoRefresh()
    autoRefreshTimer = window.setInterval(() => {
      if (document.visibilityState === 'visible') {
        fetchNotifications(notificationsPage.value)
      }
    }, 20000)
  }
  const stopAutoRefresh = () => {
    if (autoRefreshTimer) {
      clearInterval(autoRefreshTimer)
      autoRefreshTimer = null
    }
  }
  startAutoRefresh()
  const handleVisibilityChange = () => startAutoRefresh()
  document.addEventListener('visibilitychange', handleVisibilityChange)
  // 存起来以便卸载时移除
  ;(window as any)._godad_notif_handleVisibilityChange = handleVisibilityChange
})

// 监听路由参数变化，动态更新筛选
watch(() => [route.query.tab, route.query.category], () => {
  const tab = String(route.query.tab || '').toLowerCase()
  const cat = String(route.query.category || '').toLowerCase()
  const tabKeys = ['all', 'notify', 'message', 'unread']
  if (tab && tabKeys.includes(tab)) {
    activeTab.value = tab as any
  }
  categoryFilter.value = cat
  if (cat && activeTab.value !== 'message' && activeTab.value !== 'unread') {
    activeTab.value = 'notify'
  }
})

import { onBeforeUnmount } from 'vue'
onBeforeUnmount(() => {
  if (autoRefreshTimer) {
    clearInterval(autoRefreshTimer)
    autoRefreshTimer = null
  }
  const handler = (window as any)._godad_notif_handleVisibilityChange
  if (handler) {
    document.removeEventListener('visibilitychange', handler)
    delete (window as any)._godad_notif_handleVisibilityChange
  }
})
</script>

<style scoped>
/* 响应式布局 */
@media (max-width: 768px) {
  .horizontal-layout {
    flex-direction: column !important;
    height: auto !important;
  }

  .left-panel {
    width: 100% !important;
    border-right: none !important;
    border-bottom: 1px solid #e5e7eb;
  }

  .right-panel {
    width: 100% !important;
    min-height: 400px;
  }
}

/* 文本截断 */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 自定义滚动条样式 */
.overflow-y-auto::-webkit-scrollbar {
  width: 4px;
}

.overflow-y-auto::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 2px;
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* 选中状态样式 */
.notification-selected {
  background: linear-gradient(135deg, rgba(236, 72, 153, 0.1) 0%, rgba(219, 39, 119, 0.1) 100%);
  border-left: 3px solid #ec4899;
}

/* 详情面板动画 */
.detail-panel {
  transition: all 0.3s ease-in-out;
}

/* 通知项悬停效果增强 */
.notification-item {
  transition: all 0.2s ease-in-out;
}

.notification-item:hover {
  transform: translateX(4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
</style>
