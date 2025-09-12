<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-4 sm:p-6">
    <!-- 评论标题 -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between mb-4 sm:mb-6 space-y-3 sm:space-y-0">
      <h3 class="text-base sm:text-lg font-semibold text-gray-900">
        评论 ({{ totalComments }})
      </h3>
      <div class="flex items-center space-x-2">
        <select
          v-model="sortBy"
          @change="() => loadComments(true)"
          class="text-xs sm:text-sm border border-gray-300 rounded-lg px-2 sm:px-3 py-1 focus:ring-2 focus:ring-pink-500 focus:border-transparent bg-white"
        >
          <option value="newest">最新</option>
          <option value="oldest">最早</option>
          <option value="most_liked">最多点赞</option>
        </select>
      </div>
    </div>

    <!-- 发表评论 -->
    <div v-if="authStore.isAuthenticated" class="mb-6 sm:mb-8">
      <div class="flex items-start space-x-2 sm:space-x-3">
        <img
          :src="authStore.user?.avatar || '/default-avatar.png'"
          :alt="authStore.user?.username"
          class="w-8 h-8 sm:w-10 sm:h-10 rounded-full object-cover border-2 border-pink-100"
        />
        <div class="flex-1">
          <textarea
            ref="commentTextareaRef"
            v-model="newComment"
            placeholder="写下您的评论..."
            rows="3"
            maxlength="500"
            class="w-full px-3 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent resize-none bg-gray-50 focus:bg-white transition-colors"
          ></textarea>
          <div class="flex flex-col sm:flex-row sm:items-center justify-between mt-2 space-y-2 sm:space-y-0">
            <div class="text-xs sm:text-sm text-gray-500">
              {{ newComment.length }}/500
            </div>
            <div class="flex space-x-2">
              <button
                v-if="newComment.trim()"
                @click="clearComment"
                class="px-3 py-1 text-xs sm:text-sm text-gray-600 hover:text-gray-800 transition-colors rounded-lg hover:bg-gray-100"
              >
                取消
              </button>
              <button
                @click="submitComment"
                :disabled="!newComment.trim() || isSubmitting"
                :class="[
                  'px-3 sm:px-4 py-1 sm:py-2 text-xs sm:text-sm rounded-lg transition-all duration-200 flex items-center font-medium',
                  newComment.trim() && !isSubmitting
                    ? 'bg-gradient-to-r from-pink-600 to-purple-600 text-white hover:from-pink-700 hover:to-purple-700 shadow-sm hover:shadow-md transform hover:scale-105'
                    : 'bg-gray-300 text-gray-500 cursor-not-allowed'
                ]"
              >
                <LoaderIcon v-if="isSubmitting" class="h-3 w-3 sm:h-4 sm:w-4 mr-1 animate-spin" />
                发表评论
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 未登录提示 -->
    <div v-else class="mb-6 sm:mb-8 p-3 sm:p-4 bg-gradient-to-r from-pink-50 to-purple-50 border border-pink-200 rounded-lg text-center">
      <p class="text-gray-600 mb-2 text-sm sm:text-base">登录后可以发表评论</p>
      <router-link
        to="/login"
        class="inline-flex items-center px-3 py-1 text-sm font-medium text-pink-600 hover:text-pink-700 bg-white rounded-lg border border-pink-200 hover:border-pink-300 transition-all duration-200 hover:shadow-sm"
      >
        立即登录
      </router-link>
    </div>

    <!-- 加载状态 -->
    <div v-if="isLoading" class="flex justify-center py-6 sm:py-8">
      <div class="animate-spin rounded-full h-6 w-6 sm:h-8 sm:w-8 border-b-2 border-pink-600"></div>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="text-center py-6 sm:py-8">
      <p class="text-red-600 mb-4 text-sm sm:text-base">{{ error }}</p>
      <button
        @click="() => loadComments(true)"
        class="px-4 py-2 text-sm font-medium text-pink-600 bg-pink-50 hover:bg-pink-100 rounded-lg border border-pink-200 hover:border-pink-300 transition-all duration-200"
      >
        重试
      </button>
    </div>

    <!-- 评论列表 -->
    <div v-else-if="comments.length > 0" class="space-y-4 sm:space-y-6">
      <CommentItem
        v-for="comment in comments"
        :key="comment.id"
        :comment="comment"
        :article-id="props.articleId"
        :article-author-id="props.articleAuthorId"
        :all-comments="comments"
        :depth="0"
        :max-depth="3"
        @reply-added="handleReplyAdded"
        @comment-deleted="handleCommentDeleted"
        @like="handleLike"
      />

      <!-- 加载更多 -->
      <div v-if="hasMore" class="text-center pt-4 sm:pt-6">
        <button
          @click="loadMoreComments"
          :disabled="isLoadingMore"
          class="px-4 sm:px-6 py-2 text-sm sm:text-base text-pink-600 border border-pink-600 rounded-lg hover:bg-pink-50 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center mx-auto font-medium hover:shadow-sm"
        >
          <LoaderIcon v-if="isLoadingMore" class="h-3 w-3 sm:h-4 sm:w-4 mr-2 animate-spin" />
          {{ isLoadingMore ? '加载中...' : '加载更多评论' }}
        </button>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="text-center py-8 sm:py-12">
      <MessageCircleIcon class="h-12 w-12 sm:h-16 sm:w-16 text-gray-300 mx-auto mb-3 sm:mb-4" />
      <p class="text-gray-500 text-base sm:text-lg mb-2">暂无评论</p>
      <p class="text-gray-400 text-sm sm:text-base">成为第一个评论的人吧！</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, defineProps, defineEmits } from 'vue'
import { LoaderIcon, MessageCircleIcon } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { CommentApi } from '@/api/comment'
import CommentItem from './CommentItem.vue'
import type { Comment, CommentCreateRequest } from '@/api/types'

// 组件属性
interface Props {
  articleId: number
  articleAuthorId?: number
}

const props = defineProps<Props>()

// 组件事件
interface Emits {
  (e: 'commentAdded'): void
  (e: 'commentDeleted'): void
}

const emit = defineEmits<Emits>()

// 状态管理
const authStore = useAuthStore()

// 响应式数据
const comments = ref<Comment[]>([])
const newComment = ref('')
const isLoading = ref(false)
const isSubmitting = ref(false)
const isLoadingMore = ref(false)
const error = ref('')
const sortBy = ref<'newest' | 'oldest' | 'most_liked'>('most_liked')
const currentPage = ref(1)
const hasMore = ref(false)
const totalComments = ref(0)

// 模板引用
const commentTextareaRef = ref<HTMLTextAreaElement>()

// 清空评论输入
const clearComment = () => {
  newComment.value = ''
}

// 提交评论
const submitComment = async () => {
  if (!newComment.value.trim() || isSubmitting.value) return
  
  try {
    isSubmitting.value = true
    
    const commentData: CommentCreateRequest = {
       content: newComment.value.trim(),
       article_id: props.articleId
     }
     
     const comment = await CommentApi.createComment(commentData)
    
    // 将新评论添加到列表顶部
    comments.value.unshift(comment.data)
    newComment.value = ''
    totalComments.value++
    
    // 通知父组件
    emit('commentAdded')
    
  } catch (err: any) {
    console.error('发表评论失败:', err)
    alert('发表评论失败：' + (err.message || '未知错误'))
  } finally {
    isSubmitting.value = false
  }
}

// 处理回复添加
const handleReplyAdded = (reply: Comment) => {
  // 回复已经在CommentItem组件中处理，这里只需要触发事件
  emit('commentAdded')
}

// 处理评论删除
const handleCommentDeleted = (commentId: number) => {
  // 首先尝试从主评论中查找
  const mainCommentIndex = comments.value.findIndex(c => c.id === commentId)
  if (mainCommentIndex > -1) {
    // 如果是主评论，直接删除
    comments.value.splice(mainCommentIndex, 1)
    totalComments.value--
  } else {
    // 如果不是主评论，说明是回复，从各个主评论的回复中查找并删除
    let found = false
    for (const comment of comments.value) {
      if (comment.replies) {
        const replyIndex = comment.replies.findIndex(r => r.id === commentId)
        if (replyIndex > -1) {
          comment.replies.splice(replyIndex, 1)
          // 更新回复数量
          if (comment.reply_count !== undefined && comment.reply_count > 0) {
            comment.reply_count--
          }
          found = true
          break
        }
      }
    }
    
    // 如果删除的是回复，总评论数也要减少
    if (found) {
      totalComments.value--
    }
  }
  
  // 触发事件
  emit('commentDeleted')
}

// 处理点赞
const handleLike = async (commentId: number) => {
  try {
    const comment = comments.value.find(c => c.id === commentId)
    if (!comment) return
    
    if (comment.is_liked) {
      await CommentApi.unlikeComment(commentId)
      comment.is_liked = false
      comment.likes = Math.max(0, (comment.likes || 1) - 1)
    } else {
      await CommentApi.likeComment(commentId)
      comment.is_liked = true
      comment.likes = (comment.likes || 0) + 1
    }
    
  } catch (err: any) {
    console.error('点赞失败:', err)
  }
}

// 加载评论列表
const loadComments = async (reset = false) => {
  try {
    isLoading.value = true
    error.value = ''
    
    const page = reset ? 1 : currentPage.value
    const response = await CommentApi.getComments(props.articleId, {
        page: page,
        limit: 20,
        sort: sortBy.value
      })
      
      // 确保 response.data 是有效数组
      const commentsData = Array.isArray(response.data) ? response.data : []
      
      if (reset) {
        comments.value = commentsData
      } else {
        comments.value.push(...commentsData)
      }
      
      totalComments.value = response.total || 0
      hasMore.value = commentsData.length === 20 && comments.value.length < (response.total || 0)
    
    if (!reset) {
      currentPage.value++
    }
    
  } catch (err: any) {
    error.value = err.message || '加载评论失败'
    console.error('加载评论失败:', err)
  } finally {
    isLoading.value = false
  }
}

// 加载更多评论
const loadMoreComments = async () => {
  if (!hasMore.value || isLoadingMore.value) return
  
  isLoadingMore.value = true
  currentPage.value++
  
  try {
    await loadComments(false)
  } catch (err: any) {
    currentPage.value-- // 回退页码
  } finally {
    isLoadingMore.value = false
  }
}

// 聚焦评论输入框
const focusCommentInput = () => {
  console.log('focusCommentInput被调用')
  console.log('commentTextareaRef.value:', commentTextareaRef.value)
  
  if (commentTextareaRef.value) {
    console.log('找到评论输入框，开始聚焦')
    commentTextareaRef.value.focus()
    console.log('评论输入框聚焦完成')
  } else {
    console.warn('评论输入框引用不存在')
  }
}

// 暴露方法给父组件
defineExpose({
  focusCommentInput
})

// 组件挂载时加载评论
onMounted(() => {
  loadComments()
})
</script>