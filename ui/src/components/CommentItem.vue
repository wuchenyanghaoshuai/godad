<template>
  <div class="comment-item" :data-author-id="comment.user?.id" :data-comment-content="comment.content">
    <!-- 主评论 -->
    <div class="flex items-start space-x-2 sm:space-x-3">
      <!-- 用户头像 -->
      <UserAvatar
        :avatar="comment.user?.avatar || ''"
        :name="comment.user?.nickname || comment.user?.username || 'U'"
        :size="40"
      />
      
      <!-- 评论内容 -->
      <div class="flex-1 min-w-0">
        <!-- 用户信息和时间 -->
        <div class="flex flex-col sm:flex-row sm:items-center space-y-1 sm:space-y-0 sm:space-x-2 mb-2">
          <div class="flex items-center space-x-2">
            <!-- 整合的用户名显示：包含回复信息和内容 -->
            <span class="font-medium text-gray-900 text-sm sm:text-base">
              <router-link
                v-if="comment.user?.username"
                :to="`/users/${comment.user.username}`"
                class="font-medium text-pink-600 hover:text-pink-700 transition-colors"
              >
                {{ comment.user?.username || '匿名用户' }}
              </router-link>
              <span v-else class="font-medium text-pink-600">{{ comment.user?.username || '匿名用户' }}</span>
              <span v-if="comment.user?.role === 'admin'" class="px-1.5 sm:px-2 py-0.5 text-xs bg-red-100 text-red-700 rounded-full ml-2">
                管理员
              </span>
              <span v-if="isArticleAuthor" class="px-1.5 sm:px-2 py-0.5 text-xs bg-blue-100 text-blue-700 rounded-full ml-2">
                博主
              </span>
              <!-- 回复信息和内容在同一行 -->
              <span v-if="depth > 0 && parentComment" class="font-normal text-gray-900">
                回复了<router-link
                  v-if="parentComment.user?.username"
                  :to="`/users/${parentComment.user.username}`"
                  class="font-medium text-blue-600 hover:text-blue-700 transition-colors"
                >{{ parentComment.user?.username || '匿名用户' }}</router-link><span 
                  v-else 
                  class="font-medium text-blue-600"
                >{{ parentComment.user?.username || '匿名用户' }}</span>：<span class="text-gray-800">{{ comment.content }}</span>
              </span>
              <!-- 主评论内容 -->
              <span v-else class="font-normal text-gray-800 ml-2">{{ comment.content }}</span>
            </span>
          </div>
          <span class="text-xs sm:text-sm text-gray-500">
            {{ formatDate(comment.created_at) }}
          </span>
        </div>
        
        <!-- 操作按钮 -->
        <div class="flex items-center space-x-3 sm:space-x-4 mb-4">
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
            v-if="authStore.isAuthenticated && depth < maxDepth"
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
        <div v-if="showReplyInput" class="mb-4 bg-gray-50 p-3 rounded-lg">
          <div class="flex items-start space-x-2 sm:space-x-3">
            <UserAvatar
              :avatar="authStore.user?.avatar || ''"
              :name="authStore.user?.nickname || authStore.user?.username || 'U'"
              :size="32"
            />
            <div class="flex-1">
              <textarea
                v-model="replyContent"
                placeholder="写下您的回复..."
                rows="2"
                maxlength="500"
                class="w-full px-3 py-2 text-xs sm:text-sm border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent resize-none bg-white transition-colors"
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
        
        <!-- 子评论区域 -->
        <div v-if="hasReplies" class="space-y-3">
          <!-- 显示前3层子评论 -->
          <div
            v-for="reply in visibleReplies"
            :key="reply.id"
            :class="[
              'border-l-2 border-gray-100 pl-4',
              depth >= 2 ? 'ml-2' : 'ml-4'
            ]"
          >
            <CommentItem
              :comment="reply"
              :article-id="props.articleId"
              :depth="depth + 1"
              :max-depth="maxDepth"
              :article-author-id="props.articleAuthorId"
              :all-comments="props.allComments"
              @reply-added="handleReplyAdded"
              @comment-deleted="handleCommentDeleted"
              @like="handleLike"
            />
          </div>
          
          <!-- 展开更多回复按钮 -->
          <div v-if="hasHiddenReplies" class="ml-4">
            <button
              @click="toggleShowMore"
              class="flex items-center space-x-2 text-sm text-blue-600 hover:text-blue-800 transition-colors px-2 py-1 rounded-lg hover:bg-blue-50"
            >
              <ChevronDownIcon v-if="!showingMore" class="h-4 w-4" />
              <ChevronUpIcon v-else class="h-4 w-4" />
              <span>
                {{ showingMore ? '收起回复' : `展开更多回复 (${hiddenRepliesCount}条)` }}
              </span>
            </button>
          </div>
          
          <!-- 当达到最大深度时的提示 -->
          <div v-if="depth >= maxDepth && authStore.isAuthenticated" class="ml-4">
            <button
              @click="replyToParent"
              class="text-xs text-gray-600 hover:text-blue-600 transition-colors px-2 py-1 rounded-lg hover:bg-gray-50"
            >
              回复到上一级
            </button>
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
  LoaderIcon,
  ChevronDownIcon,
  ChevronUpIcon
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { CommentApi } from '@/api/comment'
import type { Comment, CommentCreateRequest } from '@/api/types'
import { useToast } from '@/composables/useToast'
import UserAvatar from '@/components/UserAvatar.vue'

// 组件属性
interface Props {
  comment: Comment
  articleId: number
  depth?: number
  maxDepth?: number
  articleAuthorId?: number
  allComments?: Comment[]
}

const props = withDefaults(defineProps<Props>(), {
  depth: 0,
  maxDepth: 3,
  allComments: () => []
})

// 组件事件
interface Emits {
  (e: 'reply-added', reply: Comment): void
  (e: 'comment-deleted', commentId: number): void
  (e: 'like', commentId: number): void
  (e: 'reply-to-parent', parentId: number): void
}

const emit = defineEmits<Emits>()

// 状态管理
const authStore = useAuthStore()
const { toast } = useToast()

// 响应式数据
const showReplyInput = ref(false)
const replyContent = ref('')
const isReplying = ref(false)
const showingMore = ref(false)

// 计算属性
const canDelete = computed(() => {
  if (!authStore.isAuthenticated) return false
  return authStore.isAdmin || authStore.user?.id === props.comment.user?.id
})

const isArticleAuthor = computed(() => {
  return props.articleAuthorId && props.comment.user?.id === props.articleAuthorId
})

// 查找真正的父评论
const findParentComment = (parentId: number, comments: Comment[]): Comment | null => {
  for (const comment of comments) {
    if (comment.id === parentId) {
      return comment
    }
    if (comment.replies) {
      const found = findParentComment(parentId, comment.replies)
      if (found) return found
    }
  }
  return null
}

const parentComment = computed(() => {
  if (!props.comment.parent_id || !props.allComments) return null
  return findParentComment(props.comment.parent_id, props.allComments)
})

const hasReplies = computed(() => {
  return props.comment.replies && props.comment.replies.length > 0
})

const visibleReplies = computed(() => {
  if (!hasReplies.value) return []
  
  const replies = props.comment.replies
  if (showingMore.value) {
    return replies
  }
  
  // 默认显示前2条回复
  return replies.slice(0, 2)
})

const hasHiddenReplies = computed(() => {
  return hasReplies.value && props.comment.replies.length > 2
})

const hiddenRepliesCount = computed(() => {
  if (!hasHiddenReplies.value) return 0
  return props.comment.replies.length - 2
})

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
       parent_id: props.comment.id // 回复当前评论，所以 parent_id 就是当前评论的 ID
     }
     
     const reply = await CommentApi.createComment(replyData)
     
     // 添加到当前评论的回复中
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
    toast.error('回复失败：' + (err.message || '未知错误'))
  } finally {
    isReplying.value = false
  }
}

// 回复到上一级（当达到最大深度时）
const replyToParent = () => {
  emit('reply-to-parent', props.comment.parent_id || props.comment.id)
}

// 切换显示更多回复
const toggleShowMore = () => {
  showingMore.value = !showingMore.value
}

// 点赞评论
const handleLike = async () => {
  if (!authStore.isAuthenticated) {
    return
  }
  
  // 直接传递给父组件处理
  emit('like', props.comment.id)
}

// 处理子评论的回复添加
const handleReplyAdded = (reply: Comment) => {
  emit('reply-added', reply)
}

// 处理评论删除
const handleCommentDeleted = (commentId: number) => {
  emit('comment-deleted', commentId)
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
    toast.error('删除失败：' + (error.message || '未知错误'))
  }
}
</script>

<style scoped>
.comment-item {
  @apply relative;
}
</style>
