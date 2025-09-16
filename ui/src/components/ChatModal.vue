<template>
  <!-- 聊天对话框模态框 -->
  <div
    v-if="isVisible"
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    @click.self="close"
  >
    <div class="bg-white rounded-lg shadow-xl w-96 h-[500px] flex flex-col overflow-hidden">
      <!-- 对话框头部 -->
      <div class="flex items-center justify-between p-4 border-b border-gray-200 bg-gray-50">
        <div class="flex items-center space-x-3">
          <img
            v-if="targetUser?.avatar"
            :src="targetUser.avatar"
            :alt="targetUser.nickname || targetUser.username"
            class="w-8 h-8 rounded-full"
          >
          <div
            v-else
            class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center text-white font-medium text-sm"
          >
            {{ (targetUser?.nickname || targetUser?.username || '?').charAt(0).toUpperCase() }}
          </div>
          <div>
            <h3 class="font-semibold text-gray-900">
              {{ targetUser?.nickname || targetUser?.username || '未知用户' }}
            </h3>
            <p class="text-xs text-gray-500">私信对话</p>
          </div>
        </div>
        <button
          @click="close"
          class="text-gray-400 hover:text-gray-600 transition-colors"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- 消息列表区域 -->
      <div
        ref="messageContainer"
        class="flex-1 overflow-y-auto p-4 space-y-3 bg-gray-50"
      >
        <!-- 加载状态 -->
        <div v-if="loading" class="flex justify-center items-center h-full">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-pink-600"></div>
        </div>

        <!-- 消息列表 -->
        <div v-else-if="messages.length > 0" class="space-y-3">
          <div
            v-for="message in messages"
            :key="message.id"
            class="flex"
            :class="message.sender_id === authStore.user?.id ? 'justify-end' : 'justify-start'"
          >
            <div
              class="max-w-xs px-3 py-2 rounded-lg"
              :class="message.sender_id === authStore.user?.id
                ? 'bg-pink-600 text-white rounded-br-sm'
                : 'bg-white border border-gray-200 text-gray-900 rounded-bl-sm'"
            >
              <p class="text-sm">{{ message.content }}</p>
              <p
                class="text-xs mt-1 opacity-75"
                :class="message.sender_id === authStore.user?.id ? 'text-pink-100' : 'text-gray-500'"
              >
                {{ formatTime(message.created_at) }}
              </p>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="flex items-center justify-center h-full">
          <div class="text-center">
            <svg class="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-3.582 8-8 8a8.013 8.013 0 01-7-4L2 20l4-4a8.014 8.014 0 01-2-5c0-4.418 3.582-8 8-8s8 3.582 8 8z" />
            </svg>
            <p class="text-gray-500 text-sm">还没有消息，开始聊天吧！</p>
          </div>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="border-t border-gray-200 p-4 bg-white">
        <div class="flex space-x-2">
          <input
            v-model="newMessage"
            type="text"
            placeholder="输入消息..."
            class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-pink-600 focus:border-transparent text-sm"
            @keyup.enter="sendMessage"
            :disabled="sendLoading"
          >
          <button
            @click="sendMessage"
            :disabled="!newMessage.trim() || sendLoading"
            class="px-4 py-2 bg-pink-600 text-white rounded-lg hover:bg-pink-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors text-sm font-medium"
          >
            <svg v-if="sendLoading" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span v-else>发送</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, watch, onUnmounted } from 'vue'
import { ChatAPI, type ConversationResponse, type ChatMessage, type User } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'

interface Props {
  isVisible: boolean
  targetUser: User | null
}

interface Emits {
  (e: 'close'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const authStore = useAuthStore()
const { showToast } = useToast()

// 响应式数据
const conversation = ref<ConversationResponse | null>(null)
const messages = ref<ChatMessage[]>([])
const newMessage = ref('')
const loading = ref(false)
const sendLoading = ref(false)
const messageContainer = ref<HTMLElement>()

// 关闭对话框
const close = () => {
  emit('close')
}

// 格式化时间
const formatTime = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  if (diff < 60000) { // 1分钟内
    return '刚刚'
  } else if (diff < 3600000) { // 1小时内
    return `${Math.floor(diff / 60000)}分钟前`
  } else if (date.toDateString() === now.toDateString()) { // 今天
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  } else {
    return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
  }
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messageContainer.value) {
      messageContainer.value.scrollTop = messageContainer.value.scrollHeight
    }
  })
}

// 获取或创建对话
const getOrCreateConversation = async () => {
  if (!props.targetUser?.id) return

  try {
    loading.value = true
    const response = await ChatAPI.getOrCreateConversation({ other_user_id: props.targetUser.id })
    conversation.value = response.data || response
    await loadMessages()
  } catch (error: any) {
    console.error('创建对话失败:', error)
    showToast('无法创建对话', 'error')
  } finally {
    loading.value = false
  }
}

// 加载消息
const loadMessages = async () => {
  if (!conversation.value?.id) return

  try {
    const response = await ChatAPI.getMessages(conversation.value.id, { page: 1, limit: 50 })
    messages.value = (response.data?.messages || response.messages || []).reverse()
    scrollToBottom()
  } catch (error: any) {
    console.error('加载消息失败:', error)
    showToast('加载消息失败', 'error')
  }
}

// 发送消息
const sendMessage = async () => {
  if (!newMessage.value.trim() || !conversation.value?.id || sendLoading.value || !props.targetUser?.id) return

  const content = newMessage.value.trim()
  newMessage.value = ''

  try {
    sendLoading.value = true
    const response = await ChatAPI.sendMessage({
      sender_id: authStore.user?.id!,
      receiver_id: props.targetUser.id,
      content,
      message_type: 'text'
    })

    const message = response.data || response
    messages.value.push(message)
    scrollToBottom()
  } catch (error: any) {
    console.error('发送消息失败:', error)
    showToast('发送消息失败', 'error')
    newMessage.value = content // 恢复消息内容
  } finally {
    sendLoading.value = false
  }
}

// 监听对话框显示状态
watch(() => props.isVisible, (visible) => {
  if (visible && props.targetUser) {
    getOrCreateConversation()
  } else {
    // 重置状态
    conversation.value = null
    messages.value = []
    newMessage.value = ''
  }
})

// ESC键关闭对话框
const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && props.isVisible) {
    close()
  }
}

// 添加/移除键盘监听
watch(() => props.isVisible, (visible) => {
  if (visible) {
    document.addEventListener('keydown', handleKeydown)
  } else {
    document.removeEventListener('keydown', handleKeydown)
  }
})

// 组件卸载时清理
onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
/* 消息容器滚动条样式 */
.overflow-y-auto::-webkit-scrollbar {
  width: 4px;
}

.overflow-y-auto::-webkit-scrollbar-track {
  background: transparent;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 2px;
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}
</style>