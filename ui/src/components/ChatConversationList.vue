<template>
  <div class="chat-conversation-list h-full flex flex-col">
    <!-- 头部 -->
    <div class="flex items-center justify-between p-4 border-b bg-white">
      <div class="flex items-center space-x-3">
        <h2 class="text-lg font-semibold text-gray-800">私信</h2>
        <span 
          v-if="totalUnreadCount > 0" 
          class="px-2 py-1 text-xs bg-red-100 text-red-600 rounded-full"
        >
          {{ totalUnreadCount }}
        </span>
      </div>
      <button 
        @click="refreshConversations"
        class="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
        :disabled="loading"
      >
        <svg class="w-5 h-5" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading && (!conversations || conversations.length === 0)" class="flex-1 flex items-center justify-center">
      <div class="text-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto mb-2"></div>
        <p class="text-sm text-gray-500">加载对话中...</p>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else-if="!loading && (!conversations || conversations.length === 0)" class="flex-1 flex items-center justify-center">
      <div class="text-center py-12">
        <div class="w-16 h-16 mx-auto mb-4 text-gray-300">
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-3.582 8-8 8a8.013 8.013 0 01-7-4L2 20l4-4a8.014 8.014 0 01-2-5c0-4.418 3.582-8 8-8s8 3.582 8 8z" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900 mb-1">暂无对话</h3>
        <p class="text-gray-500">开始一段新的对话吧</p>
      </div>
    </div>

    <!-- 对话列表 -->
    <div v-else class="flex-1 overflow-y-auto">
      <div 
        v-for="conversation in conversations" 
        :key="conversation.id"
        @click="selectConversation(conversation)"
        class="flex items-center p-4 hover:bg-gray-50 cursor-pointer border-b transition-colors"
        :class="{
          'bg-blue-50 border-blue-200': selectedConversationId === conversation.id
        }"
      >
        <!-- 头像 -->
        <div class="relative flex-shrink-0 mr-3">
          <img 
            :src="conversation.other_user.avatar || '/default-avatar.png'" 
            :alt="conversation.other_user.nickname || conversation.other_user.username"
            class="w-12 h-12 rounded-full object-cover"
            @error="handleAvatarError"
          >
          <div 
            v-if="conversation.unread_count > 0"
            class="absolute -top-1 -right-1 w-5 h-5 bg-red-500 text-white text-xs rounded-full flex items-center justify-center"
          >
            {{ conversation.unread_count > 99 ? '99+' : conversation.unread_count }}
          </div>
        </div>

        <!-- 对话信息 -->
        <div class="flex-1 min-w-0">
          <div class="flex items-center justify-between mb-1">
            <router-link
              :to="`/users/${conversation.other_user.username}`"
              class="font-medium text-gray-900 hover:text-pink-600 truncate transition-colors"
            >
              {{ conversation.other_user.nickname || conversation.other_user.username }}
            </router-link>
            <span class="text-xs text-gray-500 flex-shrink-0 ml-2">
              {{ formatTime(conversation.last_message_time) }}
            </span>
          </div>
          <p class="text-sm text-gray-600 truncate">
            {{ getLastMessagePreview(conversation) }}
          </p>
        </div>

        <!-- 删除按钮 -->
        <div class="flex-shrink-0 ml-2">
          <button 
            @click.stop="showDeleteConfirm(conversation)"
            class="p-1 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded transition-colors"
            :title="'删除与' + (conversation.other_user.nickname || conversation.other_user.username) + '的对话'"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
        </div>
      </div>

      <!-- 加载更多按钮 -->
      <div v-if="hasMore && !loading" class="p-4 text-center">
        <button 
          @click="loadMore"
          class="px-4 py-2 text-sm text-blue-600 hover:text-blue-700 hover:bg-blue-50 rounded-lg transition-colors"
          :disabled="loadingMore"
        >
          {{ loadingMore ? '加载中...' : '加载更多' }}
        </button>
      </div>
    </div>

    <!-- 删除确认对话框 -->
    <div v-if="showDeleteDialog" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="hideDeleteConfirm">
      <div class="bg-white rounded-lg p-6 max-w-sm mx-4" @click.stop>
        <div class="flex items-center mb-4">
          <div class="w-8 h-8 bg-red-100 text-red-600 rounded-full flex items-center justify-center mr-3">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 15.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
          </div>
          <h3 class="text-lg font-medium text-gray-900">删除对话</h3>
        </div>
        <p class="text-gray-600 mb-6">
          确定要删除与 <span class="font-medium">{{ conversationToDelete?.other_user?.nickname || conversationToDelete?.other_user?.username }}</span> 的对话吗？此操作不可恢复。
        </p>
        <div class="flex justify-end space-x-3">
          <button 
            @click="hideDeleteConfirm"
            class="px-4 py-2 text-sm text-gray-600 hover:text-gray-700 hover:bg-gray-50 rounded-lg transition-colors"
          >
            取消
          </button>
          <button 
            @click="confirmDelete"
            :disabled="deleting"
            class="px-4 py-2 text-sm bg-red-600 text-white hover:bg-red-700 rounded-lg transition-colors disabled:opacity-50"
          >
            {{ deleting ? '删除中...' : '确认删除' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { ChatAPI, type ConversationResponse } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'

// Props & Emits
const emit = defineEmits<{
  'conversation-selected': [conversation: ConversationResponse]
}>()

// Store & Composables
const authStore = useAuthStore()
const { showToast } = useToast()

// 响应式数据
const conversations = ref<ConversationResponse[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const selectedConversationId = ref<number | null>(null)
const currentPage = ref(1)
const hasMore = ref(true)

// 删除对话相关
const showDeleteDialog = ref(false)
const conversationToDelete = ref<ConversationResponse | null>(null)
const deleting = ref(false)

// 计算属性
const totalUnreadCount = computed(() => {
  return conversations.value?.reduce((total, conv) => total + conv.unread_count, 0) || 0
})

// 方法
const loadConversations = async (page = 1) => {
  if (page === 1) {
    loading.value = true
  } else {
    loadingMore.value = true
  }

  try {
    const response = await ChatAPI.getConversations({
      page,
      limit: 20
    })

    if (page === 1) {
      conversations.value = response.data.conversations || []
    } else {
      conversations.value = [...(conversations.value || []), ...(response.data.conversations || [])]
    }

    currentPage.value = page
    hasMore.value = response.data.pagination.page < response.data.pagination.total_pages
  } catch (error: any) {
    showToast(error.message || '加载对话失败', 'error')
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

const refreshConversations = () => {
  currentPage.value = 1
  hasMore.value = true
  loadConversations(1)
}

const loadMore = () => {
  if (!hasMore.value || loadingMore.value) return
  loadConversations(currentPage.value + 1)
}

const selectConversation = (conversation: ConversationResponse) => {
  selectedConversationId.value = conversation.id
  emit('conversation-selected', conversation)
}

const showDeleteConfirm = (conversation: ConversationResponse) => {
  conversationToDelete.value = conversation
  showDeleteDialog.value = true
}

const hideDeleteConfirm = () => {
  showDeleteDialog.value = false
  conversationToDelete.value = null
}

const confirmDelete = async () => {
  if (!conversationToDelete.value) return
  
  deleting.value = true
  try {
    await ChatAPI.deleteConversation(conversationToDelete.value.id)
    
    // 从列表中移除
    conversations.value = conversations.value.filter(
      conv => conv.id !== conversationToDelete.value!.id
    )
    
    // 如果删除的是当前选中的对话，清除选中状态
    if (selectedConversationId.value === conversationToDelete.value.id) {
      selectedConversationId.value = null
    }
    
    showToast('对话已删除', 'success')
    hideDeleteConfirm()
  } catch (error: any) {
    showToast(error.message || '删除对话失败', 'error')
  } finally {
    deleting.value = false
  }
}

const formatTime = (timeString?: string): string => {
  if (!timeString) return ''
  
  const time = new Date(timeString)
  const now = new Date()
  const diffMs = now.getTime() - time.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) return '刚刚'
  if (diffMins < 60) return `${diffMins}分钟前`
  if (diffHours < 24) return `${diffHours}小时前`
  if (diffDays < 7) return `${diffDays}天前`
  
  return time.toLocaleDateString()
}

const getLastMessagePreview = (conversation: ConversationResponse): string => {
  if (!conversation.last_message_content) return '暂无消息'
  
  switch (conversation.last_message_type) {
    case 'text':
      return conversation.last_message_content
    case 'image':
      return '[图片]'
    case 'emoji':
      return '[表情]'
    default:
      return '暂无消息'
  }
}

const handleAvatarError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = '/default-avatar.png'
}

// 生命周期
onMounted(() => {
  loadConversations()
})

// 暴露方法给父组件
defineExpose({
  refreshConversations,
  selectConversation
})
</script>

<style scoped>
.chat-conversation-list {
  background-color: #ffffff;
}

.chat-conversation-list ::-webkit-scrollbar {
  width: 6px;
}

.chat-conversation-list ::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.chat-conversation-list ::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.chat-conversation-list ::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>