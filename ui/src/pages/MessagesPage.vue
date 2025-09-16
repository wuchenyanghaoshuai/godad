<template>
  <div class="chat-page h-screen bg-gray-100">
    <div class="flex h-full">
      <!-- 对话列表侧边栏 -->
      <div class="w-80 bg-white border-r border-gray-200 flex flex-col">
        <ChatConversationList
          ref="conversationListRef"
          @conversation-selected="handleConversationSelected"
        />
      </div>

      <!-- 聊天主区域 -->
      <div class="flex-1 flex flex-col">
        <!-- 未选择对话的空状态 -->
        <div
          v-if="!selectedConversation"
          class="flex-1 flex items-center justify-center bg-gray-50"
        >
          <div class="text-center">
            <div class="w-20 h-20 mx-auto mb-4 text-gray-300">
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-3.582 8-8 8a8.013 8.013 0 01-7-4L2 20l4-4a8.014 8.014 0 01-2-5c0-4.418 3.582-8 8-8s8 3.582 8 8z"
                />
              </svg>
            </div>
            <h3 class="text-xl font-medium text-gray-700 mb-2">欢迎使用私信</h3>
            <p class="text-gray-500">选择一个对话开始聊天</p>
          </div>
        </div>

        <!-- 聊天界面 -->
        <template v-else>
          <!-- 消息列表 -->
          <div class="flex-1 overflow-hidden">
            <ChatMessageList
              ref="messageListRef"
              :conversation="selectedConversation"
            />
          </div>

          <!-- 输入框 -->
          <div class="border-t border-gray-200">
            <ChatInputBox
              :conversation="selectedConversation"
              @message-sent="handleMessageSent"
            />
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import ChatConversationList from '@/components/ChatConversationList.vue'
import ChatMessageList from '@/components/ChatMessageList.vue'
import ChatInputBox from '@/components/ChatInputBox.vue'
import { ChatAPI, type ConversationResponse, type ChatMessage } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'

// Store & Composables
const route = useRoute()
const authStore = useAuthStore()
const { showToast } = useToast()

// 响应式数据
const selectedConversation = ref<ConversationResponse | null>(null)

// DOM refs
const conversationListRef = ref<InstanceType<typeof ChatConversationList>>()
const messageListRef = ref<InstanceType<typeof ChatMessageList>>()

// 方法
const handleConversationSelected = (conversation: ConversationResponse) => {
  selectedConversation.value = conversation

  // 确保消息列表刷新
  nextTick(() => {
    if (messageListRef.value) {
      messageListRef.value.refreshMessages()
    }
  })
}

const handleMessageSent = (message: ChatMessage) => {
  // 添加消息到当前对话
  if (messageListRef.value) {
    messageListRef.value.addMessage(message)
  }

  // 刷新对话列表以更新最后消息时间
  if (conversationListRef.value) {
    conversationListRef.value.refreshConversations()
  }
}

// 处理路由参数（支持直接打开某个对话）
const handleRouteParams = async () => {
  const conversationId = route.query.conversation
  const userId = route.query.user

  if (conversationId && typeof conversationId === 'string') {
    try {
      // 根据对话ID获取对话信息
      const response = await ChatAPI.getConversations({ page: 1, limit: 100 })
      const conversation = response.data.conversations.find(c => c.id === parseInt(conversationId))

      if (conversation && conversationListRef.value) {
        conversationListRef.value.selectConversation(conversation)
      }
    } catch (error) {
      console.error('获取对话信息失败:', error)
    }
  } else if (userId && typeof userId === 'string') {
    try {
      // 先刷新对话列表，检查是否已存在对话
      if (conversationListRef.value) {
        conversationListRef.value.refreshConversations()
        await nextTick()

        // 检查列表中是否已有与目标用户的对话
        const conversations = conversationListRef.value.conversations || []
        let targetConversation = conversations.find(conv => {
          const otherUserId = conv.other_user?.id
          return otherUserId === parseInt(userId)
        })

        if (targetConversation) {
          // 如果已存在对话，直接选中
          conversationListRef.value.selectConversation(targetConversation)
        } else {
          // 如果不存在，创建新对话
          const response = await ChatAPI.getOrCreateConversation({ other_user_id: parseInt(userId) })

          if (response.data) {
            // 直接使用返回的对话数据，避免不必要的刷新
            // 创建后直接选中，因为后端已经保证不会重复创建
            const newConversation = response.data
            conversationListRef.value.selectConversation(newConversation)

            // 只在必要时刷新列表以更新UI显示
            conversationListRef.value.refreshConversations()
          }
        }
      }
    } catch (error) {
      console.error('处理对话失败:', error)
      showToast('无法与该用户开始对话', 'error')
    }
  }
}

// 生命周期
onMounted(async () => {
  // 清除可能冲突的localStorage状态
  const keysToRemove = Object.keys(localStorage).filter(key =>
    key.startsWith('chat_') || key.startsWith('conversation_') || key.startsWith('message_')
  )
  keysToRemove.forEach(key => localStorage.removeItem(key))

  // 检查用户登录状态
  if (!authStore.isAuthenticated) {
    showToast('请先登录', 'warning')
    return
  }

  // 处理路由参数
  await handleRouteParams()
})

// 页面标题
document.title = '私信 - GoDad'
</script>

<style scoped>
.chat-page {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .chat-page .w-80 {
    width: 100%;
    position: fixed;
    top: 0;
    left: 0;
    z-index: 20;
    height: 100vh;
    transform: translateX(-100%);
    transition: transform 0.3s ease;
  }

  .chat-page .w-80.show {
    transform: translateX(0);
  }
}

/* 自定义滚动条 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>