<template>
  <div class="chat-message-list h-full flex flex-col">
    <!-- 消息头部 -->
    <div v-if="currentConversation && currentConversation.other_user" class="flex items-center justify-between p-4 border-b bg-white sticky top-0 z-10">
      <div class="flex items-center space-x-3">
        <UserAvatar
          :avatar="currentConversation?.other_user?.avatar || ''"
          :name="currentConversation?.other_user?.nickname || currentConversation?.other_user?.username || 'U'"
          :size="40"
        />
        <div>
          <router-link
            :to="`/users/${currentConversation?.other_user?.username || ''}`"
            class="font-medium text-gray-900 hover:text-pink-600 transition-colors"
          >
            {{ currentConversation?.other_user?.nickname || currentConversation?.other_user?.username }}
          </router-link>
          <p class="text-sm text-gray-500">
            <router-link
              :to="`/users/${currentConversation?.other_user?.username || ''}`"
              class="hover:text-pink-600 transition-colors"
            >
              @{{ currentConversation?.other_user?.username }}
            </router-link>
          </p>
        </div>
      </div>
      <div class="flex items-center space-x-2">
        <button 
          @click="refreshMessages"
          class="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
          :disabled="loading"
          title="刷新消息"
        >
          <svg class="w-5 h-5" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 消息容器 -->
    <div ref="messagesContainer" class="flex-1 overflow-y-auto p-4 space-y-3 bg-gray-50" @scroll="handleScroll">
      <!-- 顶部自动加载锚点（被看到时自动加载更早消息） -->
      <div ref="topSentinel" class="h-1"></div>
      <!-- 加载更多按钮 -->
    <div class="text-center">
      <!-- 顶部加载骨架 -->
      <div v-if="loadingMore" class="py-2 text-gray-500 text-xs flex items-center justify-center gap-2">
        <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-500"></div>
        加载中...
      </div>
      <!-- 备用按钮（IntersectionObserver 不可用时显示） -->
      <button 
        v-else-if="hasMore && !loading && !supportsIO"
        @click="loadMore"
        class="px-4 py-2 text-sm text-blue-600 hover:text-blue-700 hover:bg-blue-50 rounded-lg transition-colors"
      >
        加载更早的消息
      </button>
    </div>

      <!-- 初始加载状态 -->
      <div v-if="loading && (!messages || messages.length === 0)" class="flex items-center justify-center py-12">
        <div class="text-center">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto mb-2"></div>
          <p class="text-sm text-gray-500">加载消息中...</p>
        </div>
      </div>

      <!-- 无消息状态 -->
      <div v-else-if="!loading && (!messages || messages.length === 0)" class="flex items-center justify-center py-12">
        <div class="text-center">
          <div class="w-16 h-16 mx-auto mb-4 text-gray-300">
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-3.582 8-8 8a8.013 8.013 0 01-7-4L2 20l4-4a8.014 8.014 0 01-2-5c0-4.418 3.582-8 8-8s8 3.582 8 8z" />
            </svg>
          </div>
          <p class="text-gray-500">开始你们的对话吧</p>
        </div>
      </div>

      <!-- 消息列表 -->
      <div v-for="(message, index) in messages" :key="message.id" class="message-group">
        <!-- 居中时间戳分隔符 -->
        <div v-if="shouldShowCenterTime(index).show" class="flex items-center justify-center my-4">
          <div class="px-3 py-1 bg-white text-xs text-gray-500 rounded-full shadow-sm">
            {{ shouldShowCenterTime(index).type === 'date' ?
                formatDateTimeSeparator(message.created_at) :
                formatTimeSeparator(message.created_at) }}
          </div>
        </div>

        <!-- 消息气泡 -->
        <div class="flex" :class="message.sender_id === currentUserId ? 'justify-end' : 'justify-start'">
          <!-- 对方头像（左侧） -->
          <div v-if="message.sender_id !== currentUserId" class="flex-shrink-0 mr-2">
            <UserAvatar :avatar="message.sender.avatar || ''" :name="message.sender.nickname || message.sender.username || 'U'" :size="32" />
          </div>

          <!-- 消息内容 -->
          <div class="max-w-xs lg:max-w-md group relative">
            <!-- 悬停时间戳（所有消息都有） -->
            <div
              class="absolute top-0 opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none z-10 text-xs text-gray-500 bg-white px-2 py-1 rounded shadow-lg whitespace-nowrap"
              :class="message.sender_id === currentUserId ? '-left-20' : '-right-20'"
            >
              {{ formatDetailTime(message.created_at) }}
            </div>

            <div
              class="relative rounded-2xl px-4 py-2 shadow-sm"
              :class="message.sender_id === currentUserId
                ? 'bg-blue-500 text-white'
                : 'bg-white text-gray-900'"
            >
              <!-- 文本消息 -->
              <div v-if="message.message_type === 'text'" class="break-words">
                {{ message.content }}
              </div>

              <!-- 图片消息 -->
              <div v-else-if="message.message_type === 'image'" class="space-y-2">
                <div 
                  v-for="(image, index) in message.images" 
                  :key="index"
                  class="relative"
                >
                  <img 
                    :src="image.url" 
                    :alt="`图片 ${index + 1}`"
                    class="max-w-full rounded-lg cursor-pointer hover:opacity-90 transition-opacity"
                    :style="{ maxWidth: '200px', maxHeight: '200px' }"
                    @click="openImagePreview(image.url)"
                    @error="handleImageError"
                    loading="lazy"
                  >
                  <div class="text-xs mt-1 opacity-75">
                    {{ formatImageSize(image.size) }}
                  </div>
                </div>
              </div>

              <!-- 表情消息 -->
              <div v-else-if="message.message_type === 'emoji'" class="text-2xl">
                {{ message.emoji?.image_url || '😀' }}
              </div>

              <!-- 轻量私信：不展示已读双勾，减少占位 -->
            </div>
          </div>

          <!-- 自己的头像（右侧） -->
          <div v-if="message.sender_id === currentUserId" class="flex-shrink-0 ml-2">
            <img 
              :src="message.sender.avatar || '/default-avatar.png'" 
              :alt="message.sender.nickname || message.sender.username"
              class="w-8 h-8 rounded-full object-cover"
              @error="handleAvatarError"
            >
          </div>
        </div>
      </div>
      <!-- 新消息提示（不在底部时出现） -->
      <div v-if="newMessageCount > 0 && !isAtBottom" class="sticky bottom-2 flex justify-center">
        <button @click="jumpToBottom" class="px-3 py-1.5 text-xs bg-blue-600 text-white rounded-full shadow hover:bg-blue-700">
          有 {{ newMessageCount }} 条新消息，点击查看
        </button>
      </div>
      <!-- 回到底部按钮（无新消息但不在底部时显示） -->
      <div v-else-if="!isAtBottom" class="sticky bottom-2 flex justify-center">
        <button @click="jumpToBottom" class="px-2.5 py-1.5 text-xs bg-gray-700 text-white rounded-full shadow hover:bg-gray-800 flex items-center gap-1">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
          回到底部
        </button>
      </div>
    </div>

    <!-- 图片预览模态框 -->
    <div v-if="imagePreview.show" class="fixed inset-0 bg-black bg-opacity-90 flex items-center justify-center z-50" @click="closeImagePreview">
      <div class="relative max-w-4xl max-h-full p-4" @click.stop>
        <img 
          :src="imagePreview.url" 
          alt="预览图片"
          class="max-w-full max-h-full object-contain"
        >
        <button 
          @click="closeImagePreview"
          class="absolute top-2 right-2 w-8 h-8 bg-black bg-opacity-50 text-white rounded-full flex items-center justify-center hover:bg-opacity-70 transition-colors"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onBeforeUnmount, onMounted } from 'vue'
import { ChatAPI, type ChatMessage, type ConversationResponse } from '@/api'
import { NotificationApi } from '@/api/notification'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import { useNotificationSync } from '@/composables/useNotificationSync'
import UserAvatar from '@/components/UserAvatar.vue'

// Props
interface Props {
  conversation?: ConversationResponse | null
}
const props = defineProps<Props>()

// Store & Composables
const authStore = useAuthStore()
const { showToast } = useToast()
const { triggerRefresh } = useNotificationSync()

// 响应式数据
const messages = ref<ChatMessage[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const messagesContainer = ref<HTMLDivElement>()
const topSentinel = ref<HTMLDivElement>()
const currentPage = ref(1)
const hasMore = ref(true)
const isAtBottom = ref(true)
const newMessageCount = ref(0)
const supportsIO = ref(false)
const imagePreview = ref({
  show: false,
  url: ''
})

// 限制渲染条数，避免超长列表拖慢页面（轻量私信：保留最近10条）
const MAX_RENDERED_MESSAGES = 10
const capMessagesIfNeeded = () => {
  if (!isAtBottom.value) return
  const extra = messages.value.length - MAX_RENDERED_MESSAGES
  if (extra > 0) {
    messages.value.splice(0, extra)
  }
}

// 计算属性
const currentConversation = computed(() => props.conversation)
const currentUserId = computed(() => authStore.user?.id)

// 方法
const loadMessages = async (page = 1) => {
  if (!currentConversation.value) return

  if (page === 1) {
    loading.value = true
  } else {
    loadingMore.value = true
  }

  try {
    const response = await ChatAPI.getMessages(currentConversation.value.id, {
      page,
      limit: 10
    })

    if (page === 1) {
      messages.value = response.data.messages // 保持原有顺序：旧消息在前，新消息在后
    } else {
      // 加载更早的消息，添加到数组开头
      messages.value = [...response.data.messages, ...messages.value]
    }

    currentPage.value = page
    hasMore.value = response.data.pagination.page < response.data.pagination.total_pages

    // 确保首次加载时滚动到底部
    if (page === 1) {
      await nextTick()
      setTimeout(() => {
        scrollToBottomAnimated()
        markAsRead()
        capMessagesIfNeeded()
      }, 100)
      // 再次确保滚至底部（处理图片加载后高度变化）
      setTimeout(() => {
        scrollToBottomAnimated()
      }, 300)
    }
  } catch (error: any) {
    showToast(error.message || '加载消息失败', 'error')
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

const refreshMessages = () => {
  currentPage.value = 1
  hasMore.value = true
  loadMessages(1)
}

const loadMore = async () => {
  if (!hasMore.value || loadingMore.value) return
  
  const scrollHeight = messagesContainer.value?.scrollHeight || 0
  const scrollTop = messagesContainer.value?.scrollTop || 0
  
  await loadMessages(currentPage.value + 1)
  
  // 保持滚动位置
  await nextTick()
  if (messagesContainer.value) {
    const newScrollHeight = messagesContainer.value.scrollHeight
    messagesContainer.value.scrollTop = scrollTop + (newScrollHeight - scrollHeight)
  }
}

const markAsRead = async () => {
  if (!currentConversation.value) return

  try {
    await ChatAPI.markAsRead(currentConversation.value.id)
    // 同时清除相关的私信通知
    await clearMessageNotifications()
  } catch (error) {
    // 标记已读失败不影响用户体验，静默处理
  }
}

// 清除来自当前对话对象的私信通知
const clearMessageNotifications = async () => {
  if (!currentConversation.value) return

  try {
    // 获取当前未读通知
    const notificationsResponse = await NotificationApi.getNotifications({
      page: 1,
      limit: 100 // 获取足够多的通知来查找匹配的
    })

    if (notificationsResponse.code === 200) {
      const messageNotifications = notificationsResponse.data.notifications.filter(
        notification =>
          notification.type === 'message' &&
          notification.actor_id === currentConversation.value?.other_user?.id &&
          !notification.is_read
      )

      if (messageNotifications.length > 0) {
        const notificationIds = messageNotifications.map(n => n.id)
        await NotificationApi.markAsRead(notificationIds)
        // 触发通知刷新，更新Navbar中的通知数量
        triggerRefresh()
      }
    }
  } catch (error) {
    // 清除通知失败也不影响用户体验，静默处理
    console.debug('清除私信通知失败:', error)
  }
}

const scrollToBottomAnimated = () => {
  if (!messagesContainer.value) return

  // 强制滚动到底部，不使用动画确保立即生效
  messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  isAtBottom.value = true
  newMessageCount.value = 0
}


// 格式化详细时间（悬停显示）
const formatDetailTime = (timeString: string): string => {
  const date = new Date(timeString)
  return date.toLocaleString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 格式化日期时间分隔符（每天第一条消息）
const formatDateTimeSeparator = (timeString: string): string => {
  const date = new Date(timeString)
  const today = new Date()
  const yesterday = new Date(today)
  yesterday.setDate(today.getDate() - 1)

  if (date.toDateString() === today.toDateString()) {
    return `今天 ${date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })}`
  } else if (date.toDateString() === yesterday.toDateString()) {
    return `昨天 ${date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })}`
  } else {
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }
}

// 格式化时间分隔符（间隔1小时以上）
const formatTimeSeparator = (timeString: string): string => {
  const date = new Date(timeString)
  return date.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 判断是否应该显示居中时间戳
const shouldShowCenterTime = (currentIndex: number): { show: boolean, type: 'date' | 'time' } => {
  if (currentIndex === 0) return { show: true, type: 'date' } // 第一条消息显示日期

  const currentMessage = messages.value[currentIndex]
  const previousMessage = messages.value[currentIndex - 1]

  if (!currentMessage || !previousMessage) return { show: true, type: 'date' }

  const currentTime = new Date(currentMessage.created_at)
  const prevTime = new Date(previousMessage.created_at)

  // 检查是否跨天
  const isDifferentDay = currentTime.toDateString() !== prevTime.toDateString()
  if (isDifferentDay) {
    return { show: true, type: 'date' }
  }

  // 检查是否超过1小时
  const diffHours = (currentTime.getTime() - prevTime.getTime()) / (1000 * 60 * 60)
  if (diffHours >= 1) {
    return { show: true, type: 'time' }
  }

  return { show: false, type: 'time' }
}

const formatImageSize = (size: number): string => {
  if (size < 1024) return `${size}B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)}KB`
  return `${(size / (1024 * 1024)).toFixed(1)}MB`
}

const openImagePreview = (url: string) => {
  imagePreview.value = { show: true, url }
}

const closeImagePreview = () => {
  imagePreview.value = { show: false, url: '' }
}

const handleAvatarError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = '/default-avatar.png'
}

const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = '/placeholder-image.png' // 添加一个占位图片
}

// 添加新消息到列表
const addMessage = (message: ChatMessage) => {
  // 检查消息是否已存在，避免重复添加
  const existingMessage = messages.value.find(msg => msg.id === message.id)
  if (!existingMessage) {
    // 暂时停止轮询，避免干扰
    const wasPolling = pollingTimer !== null
    if (wasPolling) {
      stopPolling()
    }

    // 添加到消息列表末尾
    messages.value.push(message)

    // 立即滚动到底部并标记为已读
    nextTick(() => {
      if (isAtBottom.value) {
        scrollToBottomAnimated()
        markAsRead()
        capMessagesIfNeeded()
      } else {
        newMessageCount.value += 1
      }

      // 延迟重启轮询，给消息显示足够时间
      if (wasPolling) {
        setTimeout(() => {
          startPolling()
        }, 2000)
      }
    })
  }
}

// 轮询检查新消息
let pollingTimer: NodeJS.Timeout | null = null
let topObserver: IntersectionObserver | null = null
const TOP_LOAD_THRESHOLD = 20

const startPolling = () => {
  if (pollingTimer) return

  pollingTimer = setInterval(async () => {
    // 仅在页面可见时进行低频轮询
    if (document.visibilityState !== 'visible') return
    if (currentConversation.value && !loading.value) {
      try {
        const response = await ChatAPI.getMessages(currentConversation.value.id, {
          page: 1,
          limit: 10
        })

        const newMessages = response.data.messages

        if (newMessages && newMessages.length > 0) {
          const sortedNewMessages = [...newMessages] // 保持后端返回的顺序：旧消息在前，新消息在后
          const latestNewMessage = sortedNewMessages[sortedNewMessages.length - 1]
          const currentLatestMessage = messages.value[messages.value.length - 1]

          // 简单的新消息检查逻辑
          if (!currentLatestMessage || latestNewMessage.id > currentLatestMessage.id) {
            // 如果当前没有消息，直接设置所有消息
            if (messages.value.length === 0) {
              messages.value = sortedNewMessages
              await nextTick()
              scrollToBottomAnimated()
              markAsRead()
              capMessagesIfNeeded()
            } else {
              // 只添加新消息，避免重复
              const newMessagesToAdd = sortedNewMessages.filter(msg =>
                msg.id > currentLatestMessage.id
              )
              if (newMessagesToAdd.length > 0) {
                messages.value.push(...newMessagesToAdd)
                await nextTick()
                if (isAtBottom.value) {
                  scrollToBottomAnimated()
                  markAsRead()
                  capMessagesIfNeeded()
                } else {
                  newMessageCount.value += newMessagesToAdd.length
                }
              }
            }
          }
        }
      } catch (error) {
        // 静默处理错误，避免影响用户体验
        console.debug('轮询获取消息失败:', error)
      }
    }
  }, 10000) // 10s 低频轮询，降低负担
}

const stopPolling = () => {
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = null
  }
}

// 监听对话变化
watch(
  () => props.conversation,
  (newConversation, _oldConversation) => {
    stopPolling() // 停止之前的轮询

    if (newConversation) {
      messages.value = []
      currentPage.value = 1
      hasMore.value = true
      loadMessages(1)
      startPolling() // 开始轮询新消息
    } else {
      messages.value = []
    }
  },
  { immediate: true }
)

// 生命周期
onBeforeUnmount(() => {
  stopPolling()
  if (topObserver) {
    topObserver.disconnect()
    topObserver = null
  }
})

// 暴露方法给父组件
defineExpose({
  addMessage,
  refreshMessages,
  scrollToBottomAnimated
})

// 滚动与提示逻辑
const handleScroll = () => {
  if (!messagesContainer.value) return
  const el = messagesContainer.value
  const threshold = 10
  const atBottom = el.scrollHeight - el.clientHeight - el.scrollTop <= threshold
  isAtBottom.value = atBottom
  if (atBottom) {
    newMessageCount.value = 0
  }

  // 到顶自动加载更早的消息
  if (el.scrollTop <= TOP_LOAD_THRESHOLD && hasMore.value && !loadingMore.value && !loading.value) {
    // 防抖：下一帧再触发，避免一次滚动多次触发
    requestAnimationFrame(() => {
      if (messagesContainer.value && messagesContainer.value.scrollTop <= TOP_LOAD_THRESHOLD) {
        loadMore()
      }
    })
  }
}

// 顶部自动加载：IntersectionObserver 作为滚动阈值的补充，更稳定
onMounted(() => {
  supportsIO.value = 'IntersectionObserver' in window
  if (messagesContainer.value && topSentinel.value && supportsIO.value) {
    try {
      topObserver = new IntersectionObserver((entries) => {
        const entry = entries[0]
        if (entry && entry.isIntersecting) {
          if (hasMore.value && !loadingMore.value && !loading.value) {
            loadMore()
          }
        }
      }, { root: messagesContainer.value, rootMargin: '0px', threshold: 0 })
      topObserver.observe(topSentinel.value)
    } catch (_) {
      // 忽略观察器异常，回退到滚动监听
    }
  }
})

const jumpToBottom = () => {
  scrollToBottomAnimated()
  markAsRead()
}
</script>

<style scoped>
.chat-message-list {
  height: 100%;
}

.chat-message-list ::-webkit-scrollbar {
  width: 6px;
}

.chat-message-list ::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.chat-message-list ::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.chat-message-list ::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

.message-group {
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
