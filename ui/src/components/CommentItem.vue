<template>
  <div class="comment-item">
    <!-- 主评论 -->
    <div class="flex items-start space-x-2 sm:space-x-3">
      <!-- 用户头像 -->
      <img
        :src="comment.user?.avatar || '/default-avatar.png'"
        :alt="comment.user?.username"
        class="w-8 h-8 sm:w-10 sm:h-10 rounded-full object-cover flex-shrink-0 border-2 border-gray-100"
      />
      
      <!-- 评论内容 -->
      <div class="flex-1 min-w-0">
        <!-- 用户信息和时间 -->
        <div class="flex flex-col sm:flex-row sm:items-center space-y-1 sm:space-y-0 sm:space-x-2 mb-2">
          <div class="flex items-center space-x-2">
            <span class="font-medium text-gray-900 text-sm sm:text-base">
              {{ comment.user?.username || '匿名用户' }}
            </span>
            <span v-if="comment.user?.role === 'admin'" class="px-1.5 sm:px-2 py-0.5 text-xs bg-red-100 text-red-700 rounded-full">
              管理员
            </span>
          </div>
          <span class="text-xs sm:text-sm text-gray-500">
            {{ formatDate(comment.created_at) }}
          </span>
        </div>
        
        <!-- 评论文本 -->
        <div class="text-gray-800 mb-3 whitespace-pre-wrap break-words text-sm sm:text-base leading-relaxed">
          {{ comment.content }}
        </div>
        
        <!-- 操作按钮 -->
        <div class="flex items-center space-x-3 sm:space-x-4">
          <!-- 点赞按钮 -->
          <button
            v-if="authStore.isAuthenticated"
            @click="handleLike"
            :class="[
              'flex items-center space-x-1 text-xs sm:text-sm transition-all duration-200 px-2 py-1 rounded-lg hover:bg-gray-50',
              comment.is_liked
                ? 'text-pink-600 hover:text-pink-700'
                : 'text-gray-500 hover:text-gray-700'
            ]"
          >
            <HeartIcon :class="['h-3 w-3 sm:h-4 sm:w-4', comment.is_liked ? 'fill-current' : '']"/>
            <span>{{ comment.likes || 0 }}</span>
          </button>
          
          <!-- 回复按钮 -->
          <button
            v-if="authStore.isAuthenticated"
            @click="toggleReply"
            class="flex items-center space-x-1 text-xs sm:text-sm text-gray-500 hover:text-gray-700 transition-all duration-200 px-2 py-1 rounded-lg hover:bg-gray-50"
          >
            <MessageCircleIcon class="h-3 w-3 sm:h-4 sm:w-4" />
            <span>回复</span>
          </button>
          
          <!-- 删除按钮 -->
          <button
            v-if="canDelete"
            @click="handleDelete"
            class="flex items-center space-x-1 text-xs sm:text-sm text-red-500 hover:text-red-700 transition-all duration-200 px-2 py-1 rounded-lg hover:bg-red-50"
          >
            <TrashIcon class="h-3 w-3 sm:h-4 sm:w-4" />
            <span>删除</span>
          </button>
        </div>
        
        <!-- 回复输入框 -->
        <div v-if="showReplyInput" class="mt-3 sm:mt-4">
          <div class="flex items-start space-x-2 sm:space-x-3">
            <img
              :src="authStore.user?.avatar || '/default-avatar.png'"
              :alt="authStore.user?.username"
              class="w-6 h-6 sm:w-8 sm:h-8 rounded-full object-cover flex-shrink-0 border border-gray-200"
            />
            <div class="flex-1">
              <textarea
                v-model="replyContent"
                placeholder="写下您的回复..."
                rows="2"
                maxlength="500"
                class="w-full px-3 py-2 text-xs sm:text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent resize-none bg-gray-50 focus:bg-white transition-colors"
              ></textarea>
              <div class="flex flex-col sm:flex-row sm:items-center justify-between mt-2 space-y-2 sm:space-y-0">
                <div class="text-xs text-gray-500">
                  {{ replyContent.length }}/500
                </div>
                <div class="flex space-x-2">
                  <button
                    @click="cancelReply"
                    class="px-2 sm:px-3 py-1 text-xs text-gray-600 hover:text-gray-800 transition-all duration-200 rounded-lg hover:bg-gray-100"
                  >
                    取消
                  </button>
                  <button
                    @click="submitReply"
                    :disabled="!replyContent.trim() || isReplying"
                    :class="[
                      'px-2 sm:px-3 py-1 text-xs rounded-lg transition-all duration-200 flex items-center font-medium',
                      replyContent.trim() && !isReplying
                        ? 'bg-gradient-to-r from-pink-600 to-purple-600 text-white hover:from-pink-700 hover:to-purple-700 shadow-sm hover:shadow-md transform hover:scale-105'
                        : 'bg-gray-300 text-gray-500 cursor-not-allowed'
                    ]"
                  >
                    <LoaderIcon v-if="isReplying" class="h-3 w-3 mr-1 animate-spin" />
                    回复
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 子评论 -->
        <div v-if="comment.replies && comment.replies.length > 0" class="mt-3 sm:mt-4 space-y-2 sm:space-y-3">
          <div
            v-for="reply in comment.replies"
            :key="reply.id"
            class="flex items-start space-x-2 sm:space-x-3 bg-gradient-to-r from-gray-50 to-blue-50 p-2 sm:p-3 rounded-lg border border-gray-100 hover:shadow-sm transition-all duration-200"
          >
            <img
              :src="reply.user.avatar || '/default-avatar.png'"
              :alt="reply.user.username"
              class="w-5 h-5 sm:w-6 sm:h-6 rounded-full object-cover flex-shrink-0 border border-gray-200"
            />
            <div class="flex-1 min-w-0">
              <div class="flex flex-wrap items-center gap-1 sm:gap-2">
                <span class="font-medium text-xs sm:text-sm text-gray-900">{{ reply.user.username }}</span>
                <span v-if="reply.user.role === 'admin'" class="px-1.5 py-0.5 text-xs bg-gradient-to-r from-pink-100 to-purple-100 text-pink-800 rounded-full border border-pink-200">
                  管理员
                </span>
                <span class="text-xs text-gray-500">{{ formatDate(reply.created_at) }}</span>
              </div>
              <p class="mt-1 text-xs sm:text-sm text-gray-700 leading-relaxed">{{ reply.content }}</p>
              <div class="flex items-center space-x-3 sm:space-x-4 mt-2">
                <button
                  @click="$emit('like', reply.id)"
                  :class="[
                    'flex items-center space-x-1 text-xs transition-all duration-200 px-2 py-1 rounded-lg hover:bg-white/50',
                    reply.is_liked ? 'text-pink-600 bg-pink-50' : 'text-gray-500 hover:text-pink-600'
                  ]"
                >
                  <HeartIcon :class="['h-3 w-3', reply.is_liked ? 'fill-current' : '']" />
                  <span class="font-medium">{{ reply.likes || 0 }}</span>
                </button>
                <button
                  v-if="canDeleteReply(reply)"
                  @click="$emit('comment-deleted', reply.id)"
                  class="flex items-center space-x-1 text-xs text-gray-500 hover:text-red-600 transition-all duration-200 px-2 py-1 rounded-lg hover:bg-red-50"
                >
                  <TrashIcon class="h-3 w-3" />
                  <span>删除</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, defineProps, defineEmits } from 'vue'
import {
  HeartIcon,
  MessageCircleIcon,
  TrashIcon,
  LoaderIcon
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { CommentApi } from '@/api/comment'
import type { Comment, CommentCreateRequest } from '@/api/types'

// 组件属性
interface Props {
  comment: Comment
  articleId: number
}

const props = defineProps<Props>()

// 组件事件
interface Emits {
  (e: 'reply-added', reply: Comment): void
  (e: 'comment-deleted', commentId: number): void
  (e: 'like', commentId: number): void
}

const emit = defineEmits<Emits>()

// 状态管理
const authStore = useAuthStore()

// 响应式数据
const showReplyInput = ref(false)
const replyContent = ref('')
const isReplying = ref(false)

// 计算属性
const canDelete = computed(() => {
  if (!authStore.isAuthenticated) return false
  return authStore.isAdmin || authStore.user?.id === props.comment.user?.id
})

// 检查是否可以删除回复
const canDeleteReply = (reply: Comment) => {
  if (!authStore.isAuthenticated) return false
  return authStore.isAdmin || authStore.user?.id === reply.user?.id
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000)
  
  if (diffInSeconds < 60) {
    return '刚刚'
  } else if (diffInSeconds < 3600) {
    return `${Math.floor(diffInSeconds / 60)}分钟前`
  } else if (diffInSeconds < 86400) {
    return `${Math.floor(diffInSeconds / 3600)}小时前`
  } else if (diffInSeconds < 2592000) {
    return `${Math.floor(diffInSeconds / 86400)}天前`
  } else {
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    })
  }
}

// 切换回复输入框
const toggleReply = () => {
  showReplyInput.value = !showReplyInput.value
  if (!showReplyInput.value) {
    replyContent.value = ''
  }
}

// 取消回复
const cancelReply = () => {
  showReplyInput.value = false
  replyContent.value = ''
}

// 提交回复
const submitReply = async () => {
  if (!replyContent.value.trim() || isReplying.value) return
  
  try {
    isReplying.value = true
    
    const replyData: CommentCreateRequest = {
       content: replyContent.value.trim(),
       article_id: props.articleId,
       parent_id: props.comment.id
     }
     
     const reply = await CommentApi.createComment(replyData)
     
     if (!props.comment.replies) {
       props.comment.replies = []
     }
     props.comment.replies.push(reply.data)
     
     // 更新回复数量
     if (props.comment.reply_count !== undefined) {
       props.comment.reply_count++
     }
     
     emit('reply-added', reply.data)
    
    // 清空输入框并隐藏
    replyContent.value = ''
    showReplyInput.value = false
    
  } catch (err: any) {
    console.error('回复失败:', err)
    alert('回复失败：' + (err.message || '未知错误'))
  } finally {
    isReplying.value = false
  }
}

// 点赞评论
const handleLike = async () => {
  if (!authStore.isAuthenticated) {
    return
  }
  
  // 直接传递给父组件处理
  emit('like', props.comment.id)
}

// 删除评论
const handleDelete = async () => {
  if (!confirm('确定要删除这条评论吗？')) {
    return
  }
  
  try {
    await CommentApi.deleteComment(props.comment.id)
    
    // 触发事件让父组件处理删除
    emit('comment-deleted', props.comment.id)
    
  } catch (error) {
    console.error('删除评论失败:', error)
  }
}
</script>

<style scoped>
.comment-item {
  @apply relative;
}

.comment-item:not(:last-child)::after {
  content: '';
  @apply absolute left-5 top-12 w-px h-full bg-gray-100;
  z-index: -1;
}
</style>