<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 导航栏 -->
    <Navbar />

    <div class="max-w-7xl mx-auto px-4 py-6">
      <!-- 页面标题 -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-gray-900">消息中心</h1>
        <p class="text-gray-600 mt-1">查看所有互动消息和通知</p>
      </div>

      <!-- 水平布局容器 -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden flex h-[600px] horizontal-layout">
        <!-- 左侧面板：通知列表 -->
        <div class="w-2/5 border-r border-gray-200 flex flex-col left-panel">
          <!-- 消息头部 -->
          <div class="p-4 border-b border-gray-200">
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
                  class="text-sm text-pink-600 hover:text-pink-700"
                >
                  全部标记为已读
                </button>
                <button
                  @click="clearAllNotifications"
                  class="text-sm text-red-600 hover:text-red-700"
                >
                  清空所有消息
                </button>
              </div>
            </div>
          </div>

          <!-- 所有消息列表 -->
          <div class="flex-1 overflow-y-auto">
            <!-- 当前私信对话项（如果有活跃对话） -->
            <div
              v-if="currentConversation && selectedNotification?.type === 'message'"
              :class="[
                'p-4 border-b border-gray-100 cursor-pointer notification-item relative group bg-pink-50 border-l-4 border-pink-500'
              ]"
            >
              <div class="flex items-start space-x-3">
                <div class="flex-shrink-0">
                  <span class="text-lg">💌</span>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center justify-between">
                    <p class="text-sm font-medium text-gray-900">
                      与 {{ currentConversation.other_user?.nickname || currentConversation.other_user?.username }} 的对话
                    </p>
                    <div class="flex items-center space-x-2">
                      <!-- 删除对话按钮 -->
                      <button
                        @click.stop="deleteCurrentConversation"
                        class="text-red-500 hover:text-red-700 hover:bg-red-50 p-1 rounded-full transition-colors"
                        title="删除对话"
                      >
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                        </svg>
                      </button>
                      <span class="text-xs text-gray-500">
                        私信聊天
                      </span>
                    </div>
                  </div>
                  <p class="text-sm text-gray-600 mt-1">点击右侧进行对话</p>
                </div>
              </div>
            </div>

            <!-- 其他通知列表 -->
            <div
              v-for="notification in allNotifications"
              :key="notification.id"
              :class="[
                'p-4 border-b border-gray-100 hover:bg-gray-50 cursor-pointer notification-item relative group',
                !notification.is_read ? 'bg-blue-50' : '',
                selectedNotification?.id === notification.id ? 'notification-selected' : ''
              ]"
              @click="handleNotificationClick(notification)"
            >
              <div class="flex items-start space-x-3">
                <div class="flex-shrink-0">
                  <span class="text-lg">
                    {{ getNotificationIcon(notification.type) }}
                  </span>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center justify-between">
                    <p class="text-sm font-medium text-gray-900">
                      {{ getNotificationTitle(notification.type) }}
                    </p>
                    <div class="flex items-center space-x-2">
                      <!-- 删除按钮（对所有通知显示） -->
                      <button
                        @click.stop="deleteNotification(notification)"
                        class="text-red-500 hover:text-red-700 hover:bg-red-50 p-1 rounded-full transition-colors opacity-0 group-hover:opacity-100"
                        :title="notification.type === 'message' ? '删除私信' : '删除通知'"
                      >
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                        </svg>
                      </button>
                      <span v-if="!notification.is_read" class="w-2 h-2 bg-blue-500 rounded-full"></span>
                      <span class="text-xs text-gray-500">
                        {{ formatTime(notification.created_at) }}
                      </span>
                    </div>
                  </div>
                  <p class="text-sm text-gray-600 mt-1 line-clamp-2">{{ notification.message }}</p>
                </div>
              </div>
            </div>

            <div v-if="!allNotifications || allNotifications.length === 0" class="p-8 text-center text-gray-500">
              <svg class="w-12 h-12 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5 5v-5zM4 19h5v-5l-5 5zm0-13h5V1L4 6zm11-5v5h5l-5-5z"/>
              </svg>
              <p>暂无互动消息</p>
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

          <!-- 私信聊天界面 -->
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
                    <p class="text-sm text-gray-500">私信聊天</p>
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
            <div class="flex-1 overflow-y-auto" ref="chatContainer">
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

            <!-- 聊天输入框 -->
            <ChatInputBox
              v-if="currentConversation"
              :conversation="currentConversation"
              @message-sent="handleMessageSent"
            />
          </div>

          <!-- 其他通知的详情显示 -->
          <div v-else class="flex flex-col h-full">
            <!-- 详情头部 -->
            <div class="p-4 border-b border-gray-200 bg-gray-50">
              <div class="flex items-center justify-between">
                <div class="flex items-center space-x-3">
                  <span class="text-xl">
                    {{ getNotificationIcon(selectedNotification.type) }}
                  </span>
                  <div>
                    <h3 class="font-semibold text-gray-900">
                      {{ getNotificationTitle(selectedNotification.type) }}
                    </h3>
                    <p class="text-sm text-gray-500">
                      {{ formatTime(selectedNotification.created_at) }}
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
                      <span class="text-sm">{{ getNotificationIcon(selectedNotification.type) }}</span>
                    </div>
                  </div>
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center justify-between mb-2">
                      <h4 class="text-base font-medium text-gray-900">
                        {{ selectedNotification.actor_nickname || selectedNotification.actor_username }}
                      </h4>
                      <span class="text-sm text-gray-500">
                        {{ formatTime(selectedNotification.created_at) }}
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
                          @click="router.push(`/articles/${selectedNotification.resource_id}#comments`)"
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
import { ref, onMounted, computed, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { NotificationApi, type NotificationWithDetails } from '../api/notification'
import { ChatAPI, type ConversationResponse } from '@/api'
import { useToast } from '@/composables/useToast'
import Navbar from '@/components/Navbar.vue'
import ChatMessageList from '@/components/ChatMessageList.vue'
import ChatInputBox from '@/components/ChatInputBox.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { showToast } = useToast()

// 当前选中的通知
const selectedNotification = ref<NotificationWithDetails | null>(null)

// 所有未读消息数量
const unreadMessagesCount = ref(0)

// 聊天相关状态
const currentConversation = ref<ConversationResponse | null>(null)
const chatContainer = ref<HTMLElement>()
const chatMessageListRef = ref<any>()

// 通知相关数据
const notifications = ref<NotificationWithDetails[]>([])
const notificationsLoading = ref(false)
const notificationsError = ref('')
const unreadNotificationsCount = ref(0)

// 分页相关
const notificationsPage = ref(1)
const notificationsTotalPages = ref(1)
const limit = ref(10)

// 计算属性 - 所有通知，按优先级排序
const allNotifications = computed(() => {
  if (!notifications.value) return []

  // 过滤需要的通知类型并按优先级排序
  const filteredNotifications = notifications.value.filter(n => ['like', 'comment', 'message', 'follow', 'bookmark'].includes(n.type))

  // 定义优先级：点赞 > 评论 > 私信 > 关注 > 收藏
  const typePriority = {
    'like': 1,
    'comment': 2,
    'message': 3,
    'follow': 4,
    'bookmark': 5
  }

  return filteredNotifications.sort((a, b) => {
    // 首先按类型优先级排序
    const priorityDiff = typePriority[a.type] - typePriority[b.type]
    if (priorityDiff !== 0) return priorityDiff

    // 同类型内按时间倒序排序（最新的在前）
    return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  })
})

// 获取通知列表
const fetchNotifications = async (page = 1) => {
  notificationsLoading.value = true
  notificationsError.value = ''

  try {
    const response = await NotificationApi.getNotifications({
      page,
      limit: limit.value
    })

    notifications.value = response.data.notifications
    notificationsPage.value = response.data.pagination.current_page
    notificationsTotalPages.value = response.data.pagination.total_pages

    // 统计未读通知数量
    unreadNotificationsCount.value = notifications.value ? notifications.value.filter(n => !n.is_read).length : 0
    unreadMessagesCount.value = unreadNotificationsCount.value
  } catch (err: any) {
    notificationsError.value = err.response?.data?.error || '获取通知列表失败'
  } finally {
    notificationsLoading.value = false
  }
}

// 获取通知图标
const getNotificationIcon = (type: string) => {
  switch (type) {
    case 'like':
      return '❤️'
    case 'comment':
      return '💬'
    case 'bookmark':
      return '⭐'
    case 'follow':
      return '👤'
    case 'message':
      return '💌'
    default:
      return '📢'
  }
}

// 获取通知标题
const getNotificationTitle = (type: string) => {
  switch (type) {
    case 'like':
      return '点赞'
    case 'comment':
      return '评论'
    case 'message':
      return '私信'
    case 'follow':
      return '关注'
    case 'bookmark':
      return '收藏'
    default:
      return '通知'
  }
}

// 处理通知点击
const handleNotificationClick = async (notification: NotificationWithDetails) => {
  // 如果通知未读，先标记为已读
  if (!notification.is_read) {
    try {
      await NotificationApi.markAsRead([notification.id])
      notification.is_read = true
      unreadNotificationsCount.value = Math.max(0, unreadNotificationsCount.value - 1)
      unreadMessagesCount.value = Math.max(0, unreadMessagesCount.value - 1)
    } catch (error) {
      console.error('标记通知已读失败:', error)
    }
  }

  selectedNotification.value = notification

  // 如果是私信通知，加载对应的对话
  if (notification.type === 'message' && notification.actor_id) {
    await loadConversationByUserId(notification.actor_id)
  } else {
    // 非私信通知，清空对话
    currentConversation.value = null
  }
}

// 查看文章
const viewArticle = (articleId: number) => {
  router.push(`/articles/${articleId}`)
}

const viewUserProfile = (notification: NotificationWithDetails) => {
  const username = notification.actor_username
  if (username) {
    router.push(`/users/${username}`)
  } else {
    showToast('无法获取用户信息', 'error')
  }
}

// 标记所有通知为已读
const markAllMessagesAsRead = async () => {
  try {
    await NotificationApi.markAllAsRead()
    notifications.value.forEach(n => n.is_read = true)
    unreadNotificationsCount.value = 0
    unreadMessagesCount.value = 0
  } catch (error) {
    console.error('标记所有通知已读失败:', error)
  }
}

// 格式化时间
const formatTime = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`

  return date.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 删除单个通知
const deleteNotification = async (notification: NotificationWithDetails) => {
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

  } catch (error: any) {
    console.error('删除通知失败:', error)
    showToast(error.response?.data?.message || '删除通知失败', 'error')

    // 如果删除失败，重新获取通知列表恢复UI状态
    await fetchNotifications()
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
    // 立即清空UI中的通知
    notifications.value = []
    selectedNotification.value = null
    unreadNotificationsCount.value = 0
    unreadMessagesCount.value = 0

    // 调用删除所有通知API
    await NotificationApi.deleteAllNotifications()

    showToast('所有消息已清空', 'success')

  } catch (error: any) {
    console.error('清空所有通知失败:', error)
    showToast(error.response?.data?.message || '清空消息失败', 'error')

    // 如果删除失败，重新获取通知列表恢复UI状态
    await fetchNotifications()
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

// 处理消息发送
const handleMessageSent = async (message: any) => {
  console.log('NotificationsPage: 新消息已发送:', message)

  // 立即将新消息添加到聊天列表中
  if (chatMessageListRef.value && chatMessageListRef.value.addMessage) {
    console.log('NotificationsPage: 调用 addMessage')
    chatMessageListRef.value.addMessage(message)
  } else {
    console.log('NotificationsPage: chatMessageListRef 不可用')
  }

  // 强制刷新消息列表确保同步
  if (chatMessageListRef.value && chatMessageListRef.value.refreshMessages) {
    console.log('NotificationsPage: 延迟刷新消息列表')
    setTimeout(() => {
      if (chatMessageListRef.value && chatMessageListRef.value.refreshMessages) {
        chatMessageListRef.value.refreshMessages()
      }
    }, 1000) // 给后端更多时间处理
  }
}

// 开始私信聊天（从关注通知等地方触发）
const startPrivateChat = async (notification: NotificationWithDetails) => {
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

    showToast(`开始与 ${currentConversation.value.other_user?.nickname || currentConversation.value.other_user?.username} 的对话`, 'success')
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
      selectedNotification: selectedNotification.value,
      currentConversation: currentConversation.value,
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

    // 恢复对话状态
    if (chatState.selectedNotification && chatState.currentConversation) {
      // 验证对话是否仍然存在
      try {
        const response = await ChatAPI.getMessages(chatState.currentConversation.id, {
          page: 1,
          limit: 1
        })

        // 如果对话仍然存在，恢复状态
        selectedNotification.value = chatState.selectedNotification
        currentConversation.value = chatState.currentConversation
        return true
      } catch (error) {
        // 对话不存在，清除保存的状态
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