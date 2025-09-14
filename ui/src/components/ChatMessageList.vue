<template>
  <div class="chat-message-list h-full flex flex-col">
    <!-- 消息头部 -->
    <div v-if="currentConversation && currentConversation.other_user" class="flex items-center justify-between p-4 border-b bg-white">
      <div class="flex items-center space-x-3">
        <img
          :src="currentConversation?.other_user?.avatar || '/default-avatar.png'"
          :alt="currentConversation?.other_user?.nickname || currentConversation?.other_user?.username"
          class="w-10 h-10 rounded-full object-cover"
          @error="handleAvatarError"
        >
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
    <div ref="messagesContainer" class="flex-1 overflow-y-auto p-4 space-y-3 bg-gray-50">
      <!-- 加载更多按钮 -->
      <div v-if="hasMore && !loading" class="text-center">
        <button 
          @click="loadMore"
          class="px-4 py-2 text-sm text-blue-600 hover:text-blue-700 hover:bg-blue-50 rounded-lg transition-colors"
          :disabled="loadingMore"
        >
          {{ loadingMore ? '加载中...' : '加载更早的消息' }}
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
      <div v-for="message in messages" :key="message.id" class="message-group">
        <!-- 日期分隔符 -->
        <div v-if="shouldShowDateSeparator(message)" class="flex items-center justify-center my-4">
          <div class="px-3 py-1 bg-white text-xs text-gray-500 rounded-full shadow-sm">
            {{ formatDateSeparator(message.created_at) }}
          </div>
        </div>

        <!-- 消息气泡 -->
        <div class="flex" :class="message.sender_id === currentUserId ? 'justify-end' : 'justify-start'">
          <!-- 对方头像（左侧） -->
          <div v-if="message.sender_id !== currentUserId" class="flex-shrink-0 mr-2">
            <img 
              :src="message.sender.avatar || '/default-avatar.png'" 
              :alt="message.sender.nickname || message.sender.username"
              class="w-8 h-8 rounded-full object-cover"
              @error="handleAvatarError"
            >
          </div>

          <!-- 消息内容 -->
          <div class="max-w-xs lg:max-w-md">
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
                {{ message.emoji?.unicode || '😀' }}
              </div>

              <!-- 消息状态指示器 -->
              <div class="flex items-center justify-end mt-1 space-x-1">
                <span class="text-xs opacity-75">
                  {{ formatMessageTime(message.created_at) }}
                </span>
                <!-- 已读状态（仅发送的消息显示） -->
                <div v-if="message.sender_id === currentUserId" class="flex items-center">
                  <svg 
                    class="w-3 h-3" 
                    :class="message.is_read ? 'text-blue-200' : 'text-blue-300'"
                    fill="none" 
                    stroke="currentColor" 
                    viewBox="0 0 24 24"
                  >
                    <path 
                      stroke-linecap="round" 
                      stroke-linejoin="round" 
                      stroke-width="2" 
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                  <svg 
                    v-if="message.is_read" 
                    class="w-3 h-3 -ml-1 text-blue-200" 
                    fill="none" 
                    stroke="currentColor" 
                    viewBox="0 0 24 24"
                  >
                    <path 
                      stroke-linecap="round" 
                      stroke-linejoin="round" 
                      stroke-width="2" 
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                </div>
              </div>
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
import { ref, computed, watch, nextTick, onMounted, onBeforeUnmount } from 'vue'
import { ChatAPI, type ChatMessage, type ConversationResponse } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'

// Props
interface Props {
  conversation?: ConversationResponse | null
}
const props = defineProps<Props>()

// Store & Composables
const authStore = useAuthStore()
const { showToast } = useToast()

// 响应式数据
const messages = ref<ChatMessage[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const messagesContainer = ref<HTMLDivElement>()
const currentPage = ref(1)
const hasMore = ref(true)
const imagePreview = ref({
  show: false,
  url: ''
})

// 计算属性
const currentConversation = computed(() => props.conversation)
const currentUserId = computed(() => authStore.user?.id)

// 方法
const loadMessages = async (page = 1, scrollToBottom = true) => {
  if (!currentConversation.value) return

  if (page === 1) {
    loading.value = true
  } else {
    loadingMore.value = true
  }

  try {
    const response = await ChatAPI.getMessages(currentConversation.value.id, {
      page,
      limit: 30
    })

    if (page === 1) {
      messages.value = response.data.messages.reverse() // 反转消息顺序，最新的在下面
    } else {
      // 加载更早的消息，添加到数组开头
      messages.value = [...response.data.messages.reverse(), ...messages.value]
    }

    currentPage.value = page
    hasMore.value = response.data.pagination.page < response.data.pagination.total_pages

    if (scrollToBottom && page === 1) {
      await nextTick()
      scrollToBottomAnimated()
    }

    // 标记消息为已读
    if (page === 1) {
      markAsRead()
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
  
  await loadMessages(currentPage.value + 1, false)
  
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
  } catch (error) {
    // 标记已读失败不影响用户体验，静默处理
  }
}

const scrollToBottomAnimated = () => {
  if (!messagesContainer.value) return
  
  messagesContainer.value.scrollTo({
    top: messagesContainer.value.scrollHeight,
    behavior: 'smooth'
  })
}

const shouldShowDateSeparator = (message: ChatMessage): boolean => {
  const messageIndex = messages.value.findIndex(m => m.id === message.id)
  if (messageIndex === 0) return true // 第一条消息总是显示日期
  
  const prevMessage = messages.value[messageIndex - 1]
  const prevDate = new Date(prevMessage.created_at).toDateString()
  const currentDate = new Date(message.created_at).toDateString()
  
  return prevDate !== currentDate
}

const formatDateSeparator = (timeString: string): string => {
  const date = new Date(timeString)
  const today = new Date()
  const yesterday = new Date(today)
  yesterday.setDate(yesterday.getDate() - 1)

  if (date.toDateString() === today.toDateString()) {
    return '今天'
  } else if (date.toDateString() === yesterday.toDateString()) {
    return '昨天'
  } else {
    return date.toLocaleDateString('zh-CN', { 
      month: 'long', 
      day: 'numeric' 
    })
  }
}

const formatMessageTime = (timeString: string): string => {
  const date = new Date(timeString)
  return date.toLocaleTimeString('zh-CN', { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
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
  console.log('ChatMessageList addMessage called with:', message)

  // 检查消息是否已存在，避免重复添加
  const existingMessage = messages.value.find(msg => msg.id === message.id)
  if (!existingMessage) {
    messages.value.push(message)
    console.log('Message added to list, total messages:', messages.value.length)
    nextTick(() => {
      scrollToBottomAnimated()
      markAsRead()
    })
  } else {
    console.log('Message already exists in list, skipping')
  }
}

// 轮询检查新消息
let pollingTimer: NodeJS.Timeout | null = null

const startPolling = () => {
  if (pollingTimer) return

  console.log('Starting message polling for conversation:', currentConversation.value?.id)

  pollingTimer = setInterval(async () => {
    if (currentConversation.value && !loading.value) {
      try {
        console.log('Polling for messages in conversation:', currentConversation.value.id)
        const response = await ChatAPI.getMessages(currentConversation.value.id, {
          page: 1,
          limit: 30
        })

        const newMessages = response.data.messages
        console.log('Polling response messages:', newMessages?.length || 0)

        if (newMessages && newMessages.length > 0) {
          const sortedNewMessages = newMessages.reverse() // API返回的是倒序，需要正序
          const latestNewMessage = sortedNewMessages[sortedNewMessages.length - 1]
          const currentLatestMessage = messages.value[messages.value.length - 1]

          console.log('Current latest message ID:', currentLatestMessage?.id, 'New latest message ID:', latestNewMessage?.id)

          // 检查是否有真正的新消息
          if (!currentLatestMessage || latestNewMessage.id > currentLatestMessage.id) {
            // 如果当前没有消息，直接设置所有消息
            if (messages.value.length === 0) {
              console.log('No current messages, setting all:', sortedNewMessages.length)
              messages.value = sortedNewMessages
            } else {
              // 只添加新消息，避免重复
              const newMessagesToAdd = sortedNewMessages.filter(msg =>
                msg.id > currentLatestMessage.id
              )
              console.log('New messages to add:', newMessagesToAdd.length)
              if (newMessagesToAdd.length > 0) {
                messages.value.push(...newMessagesToAdd)
                console.log('Added messages, total now:', messages.value.length)
              }
            }

            await nextTick()
            scrollToBottomAnimated()
            markAsRead()
          }
        }
      } catch (error) {
        // 静默处理错误，避免影响用户体验
        console.debug('轮询获取消息失败:', error)
      }
    }
  }, 1500) // 改为1.5秒，更及时响应
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
  (newConversation, oldConversation) => {
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
})

// 暴露方法给父组件
defineExpose({
  addMessage,
  refreshMessages,
  scrollToBottomAnimated
})
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