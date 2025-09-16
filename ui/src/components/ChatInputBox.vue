<template>
  <div class="chat-input-box bg-white border-t p-4">
    <!-- 图片预览区域 -->
    <div v-if="selectedImages.length > 0" class="mb-3 flex flex-wrap gap-2">
      <div 
        v-for="(image, index) in selectedImages" 
        :key="index"
        class="relative inline-block"
      >
        <img 
          :src="image.preview" 
          :alt="`预览图片 ${index + 1}`"
          class="w-16 h-16 object-cover rounded-lg border"
        >
        <button 
          @click="removeImage(index)"
          class="absolute -top-1 -right-1 w-5 h-5 bg-red-500 text-white rounded-full flex items-center justify-center text-xs hover:bg-red-600 transition-colors"
        >
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 消息限制警告 -->
    <div v-if="showMessageLimitWarning" class="mb-3 p-3 bg-amber-50 border border-amber-200 rounded-lg">
      <div class="flex items-center">
        <svg class="w-5 h-5 text-amber-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
        </svg>
        <span class="text-amber-700 text-sm">
          非互相关注的用户每日只能发送3条消息，您今日已发送 {{ messageCount }}/3 条消息
        </span>
      </div>
    </div>

    <!-- 输入区域 -->
    <div class="flex items-end space-x-2">
      <!-- 工具栏 -->
      <div class="flex items-center space-x-1">
        <!-- 图片上传按钮 -->
        <label
          class="p-2 text-gray-500 hover:text-blue-500 hover:bg-blue-50 rounded-lg cursor-pointer transition-colors"
          :class="{ 'opacity-50 cursor-not-allowed': !canSendMessage }"
        >
          <input
            ref="imageInput"
            type="file"
            multiple
            accept="image/*"
            class="hidden"
            @change="handleImageSelect"
            :disabled="!canSendMessage"
          >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </label>

        <!-- 表情按钮 -->
        <button
          @click="toggleEmojiPicker"
          :disabled="!canSendMessage"
          class="p-2 text-gray-500 hover:text-blue-500 hover:bg-blue-50 rounded-lg transition-colors"
          :class="{
            'text-blue-500 bg-blue-50': showEmojiPicker,
            'opacity-50 cursor-not-allowed': !canSendMessage
          }"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </button>
      </div>

      <!-- 输入框 -->
      <div class="flex-1 relative">
        <textarea
          ref="messageInput"
          v-model="messageText"
          @keydown="handleKeydown"
          @compositionstart="handleCompositionStart"
          @compositionend="handleCompositionEnd"
          @input="adjustTextareaHeight"
          :placeholder="placeholder"
          :disabled="sending || !canSendMessage"
          class="w-full px-3 py-2 border border-gray-300 rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          :class="{ 'opacity-50': !canSendMessage }"
          :style="{ minHeight: '40px', maxHeight: '120px' }"
          rows="1"
        />
      </div>

      <!-- 发送按钮 -->
      <button
        @click="sendMessage"
        :disabled="!canSend || sending || !canSendMessage"
        class="p-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex-shrink-0"
      >
        <svg v-if="sending" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
        </svg>
        <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
        </svg>
      </button>
    </div>

    <!-- 表情选择器 -->
    <div v-if="showEmojiPicker" class="absolute bottom-16 right-4 bg-white border rounded-lg shadow-lg p-4 w-72 z-10">
      <div class="flex justify-between items-center mb-3">
        <h4 class="font-medium text-gray-900">选择表情</h4>
        <button 
          @click="showEmojiPicker = false"
          class="text-gray-400 hover:text-gray-600"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      
      <div v-if="loadingEmojis" class="flex justify-center py-4">
        <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-500"></div>
      </div>
      
      <div v-else-if="emojis.length > 0" class="grid grid-cols-8 gap-1 max-h-48 overflow-y-auto">
        <button
          v-for="emoji in emojis"
          :key="emoji.id"
          @click="selectEmoji(emoji)"
          class="p-2 hover:bg-gray-100 rounded transition-colors text-2xl"
          :title="emoji.name"
        >
          {{ emoji.image_url }}
        </button>
      </div>
      
      <div v-else class="text-center py-4 text-gray-500">
        暂无可用表情
      </div>
    </div>

    <!-- 背景遮罩（点击关闭表情选择器） -->
    <div 
      v-if="showEmojiPicker" 
      class="fixed inset-0 z-0" 
      @click="showEmojiPicker = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, onMounted, onBeforeUnmount, watch } from 'vue'
import { ChatAPI, UploadApi, type ChatEmoji, type ConversationResponse, type ImageInfo } from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'

// Props
interface Props {
  conversation?: ConversationResponse | null
  disabled?: boolean
}
const props = withDefaults(defineProps<Props>(), {
  disabled: false
})

// Emits
const emit = defineEmits<{
  'message-sent': [message: any]
}>()

// Store & Composables
const authStore = useAuthStore()
const { showToast } = useToast()

// 响应式数据
const messageText = ref('')
const selectedImages = ref<Array<{file: File, preview: string, imageInfo?: ImageInfo}>>([])
const sending = ref(false)
const showEmojiPicker = ref(false)
const emojis = ref<ChatEmoji[]>([])
const loadingEmojis = ref(false)
const isComposing = ref(false)

// 消息限制相关
const canSendMessage = ref(true)
const mutualFollow = ref(false)
const messageCount = ref(0)
const checkingLimit = ref(false)

// DOM refs
const messageInput = ref<HTMLTextAreaElement>()
const imageInput = ref<HTMLInputElement>()

// 计算属性
const placeholder = computed(() => {
  if (props.disabled) return '请先选择一个对话'
  if (!props.conversation) return '请先选择一个对话'
  return `给 ${props.conversation.other_user?.nickname || props.conversation.other_user?.username || '对方'} 发消息... (Ctrl+Enter发送)`
})

const canSend = computed(() => {
  if (props.disabled || !props.conversation || sending.value) return false
  return messageText.value.trim().length > 0 || selectedImages.value.length > 0
})

const showMessageLimitWarning = computed(() => {
  return props.conversation && !mutualFollow.value && messageCount.value >= 3
})

// 方法
const handleKeydown = (event: KeyboardEvent) => {
  // Ctrl+Enter 或 Cmd+Enter 发送消息
  if (event.key === 'Enter' && (event.ctrlKey || event.metaKey)) {
    event.preventDefault()
    if (canSend.value) {
      sendMessage()
    }
    return
  }

  // 单纯 Enter 键：只有在不是输入法状态下才发送
  if (event.key === 'Enter' && !event.shiftKey && !isComposing.value) {
    event.preventDefault()
    if (canSend.value) {
      sendMessage()
    }
  }
}

const handleCompositionStart = () => {
  isComposing.value = true
}

const handleCompositionEnd = () => {
  isComposing.value = false
}

const adjustTextareaHeight = () => {
  if (!messageInput.value) return
  
  messageInput.value.style.height = 'auto'
  const scrollHeight = messageInput.value.scrollHeight
  const maxHeight = 120
  messageInput.value.style.height = `${Math.min(scrollHeight, maxHeight)}px`
}

const handleImageSelect = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const files = target.files
  
  if (!files || files.length === 0) return
  
  // 检查图片数量限制
  if (selectedImages.value.length + files.length > 9) {
    showToast('最多只能选择9张图片', 'warning')
    return
  }
  
  for (const file of Array.from(files)) {
    // 检查文件类型
    if (!file.type.startsWith('image/')) {
      showToast(`${file.name} 不是有效的图片文件`, 'error')
      continue
    }
    
    // 检查文件大小 (10MB)
    if (file.size > 10 * 1024 * 1024) {
      showToast(`${file.name} 文件过大，请选择小于10MB的图片`, 'error')
      continue
    }
    
    // 创建预览
    const preview = URL.createObjectURL(file)
    selectedImages.value.push({ file, preview })
  }
  
  // 清空input
  target.value = ''
}

const removeImage = (index: number) => {
  const image = selectedImages.value[index]
  URL.revokeObjectURL(image.preview) // 释放内存
  selectedImages.value.splice(index, 1)
}

const uploadImages = async (): Promise<ImageInfo[]> => {
  const uploadPromises = selectedImages.value.map(async (item) => {
    try {
      const formData = new FormData()
      formData.append('image', item.file)
      formData.append('usage', 'chat')
      
      const response = await UploadApi.uploadImage(formData)
      
      // 构造ImageInfo对象
      const imageInfo: ImageInfo = {
        url: response.data.public_url,
        width: 0, // 这些信息可能需要从图片文件中获取
        height: 0,
        size: item.file.size,
        thumbnail: response.data.public_url // 假设缩略图和原图相同，实际可能不同
      }
      
      return imageInfo
    } catch (error: any) {
      showToast(`上传图片失败: ${error.message}`, 'error')
      throw error
    }
  })
  
  return Promise.all(uploadPromises)
}

const sendMessage = async () => {
  if (!canSend.value || !props.conversation || !canSendMessage.value) return

  sending.value = true
  
  try {
    let messageData: any = {
      sender_id: authStore.user?.id,
      receiver_id: props.conversation?.other_user?.id
    }
    
    // 发送图片消息
    if (selectedImages.value.length > 0) {
      const uploadedImages = await uploadImages()
      messageData = {
        ...messageData,
        message_type: 'image',
        images: uploadedImages
      }
    }
    // 发送文本消息
    else if (messageText.value.trim()) {
      messageData = {
        ...messageData,
        message_type: 'text',
        content: messageText.value.trim()
      }
    } else {
      return
    }
    
    const response = await ChatAPI.sendMessage(messageData)
    
    // 清空输入
    messageText.value = ''
    selectedImages.value.forEach(item => URL.revokeObjectURL(item.preview))
    selectedImages.value = []
    
    // 调整输入框高度
    await nextTick()
    adjustTextareaHeight()
    
    // 发送成功事件
    emit('message-sent', response.data)

    // 重新检查消息限制（因为消息数量可能已改变）
    await checkMessageLimit()

    // 聚焦输入框
    messageInput.value?.focus()
  } catch (error: any) {
    showToast(error.message || '发送消息失败', 'error')
  } finally {
    sending.value = false
  }
}

const toggleEmojiPicker = async () => {
  showEmojiPicker.value = !showEmojiPicker.value
  
  if (showEmojiPicker.value && emojis.value.length === 0) {
    await loadEmojis()
  }
}

const loadEmojis = async () => {
  loadingEmojis.value = true
  try {
    const response = await ChatAPI.getEmojis()
    emojis.value = response.data
  } catch (error: any) {
    showToast(error.message || '加载表情失败', 'error')
  } finally {
    loadingEmojis.value = false
  }
}

const selectEmoji = (emoji: ChatEmoji) => {
  // 将表情插入到输入框中
  const cursorPosition = messageInput.value?.selectionStart || messageText.value.length
  const textBefore = messageText.value.slice(0, cursorPosition)
  const textAfter = messageText.value.slice(cursorPosition)

  // 插入表情
  messageText.value = textBefore + emoji.image_url + textAfter

  // 关闭表情选择器
  showEmojiPicker.value = false

  // 聚焦输入框并设置光标位置
  nextTick(() => {
    if (messageInput.value) {
      messageInput.value.focus()
      const newPosition = cursorPosition + emoji.image_url.length
      messageInput.value.setSelectionRange(newPosition, newPosition)
    }
    adjustTextareaHeight()
  })
}

// 检查消息发送限制
const checkMessageLimit = async () => {
  if (!props.conversation?.other_user?.id || !authStore.user?.id) {
    canSendMessage.value = true
    mutualFollow.value = false
    messageCount.value = 0
    return
  }

  checkingLimit.value = true
  try {
    const response = await ChatAPI.checkMessageLimit(props.conversation.other_user.id)
    canSendMessage.value = response.data.can_send
    mutualFollow.value = response.data.mutual_follow
    messageCount.value = response.data.message_count
  } catch (error: any) {
    console.error('检查消息限制失败:', error)
    // 发生错误时默认允许发送，避免影响正常聊天
    canSendMessage.value = true
    mutualFollow.value = false
    messageCount.value = 0
  } finally {
    checkingLimit.value = false
  }
}

// 监听对话变化，检查消息限制
watch(
  () => props.conversation,
  () => {
    if (props.conversation) {
      checkMessageLimit()
    } else {
      // 没有对话时重置状态
      canSendMessage.value = true
      mutualFollow.value = false
      messageCount.value = 0
    }
  },
  { immediate: true }
)

// 生命周期
onMounted(() => {
  nextTick(() => {
    messageInput.value?.focus()
  })
})

// 清理函数
const cleanup = () => {
  selectedImages.value.forEach(item => URL.revokeObjectURL(item.preview))
  selectedImages.value = []
}

// 组件卸载时清理
onBeforeUnmount(() => {
  cleanup()
})
</script>

<style scoped>
.chat-input-box {
  position: relative;
}

.chat-input-box textarea {
  resize: none;
  overflow-y: auto;
}

.chat-input-box textarea::-webkit-scrollbar {
  width: 4px;
}

.chat-input-box textarea::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.chat-input-box textarea::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 2px;
}

.chat-input-box textarea::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>