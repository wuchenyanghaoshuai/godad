<template>
  <div class="min-h-screen bg-gray-50">
    <BaseHeader />

    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- 页面导航 -->
      <div class="mb-6">
        <nav class="flex items-center space-x-2 text-sm text-gray-600">
          <router-link to="/" class="hover:text-pink-600 transition-colors">首页</router-link>
          <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
          </svg>
          <router-link to="/community" class="hover:text-pink-600 transition-colors">社区</router-link>
          <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
          </svg>
          <span class="text-gray-900">{{ post?.title || '帖子详情' }}</span>
        </nav>
      </div>

      <!-- 加载状态 -->
      <div v-if="isLoading" class="flex flex-col items-center justify-center py-16">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-pink-600 mb-4"></div>
        <p class="text-gray-500">加载中...</p>
      </div>

      <!-- 帖子未找到 -->
      <div v-else-if="!post && !isLoading" class="flex flex-col items-center justify-center py-16">
        <div class="text-gray-400 mb-4">
          <svg class="h-16 w-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900 mb-2">帖子不存在</h3>
        <p class="text-gray-500 mb-4">您访问的帖子可能已被删除或不存在</p>
        <router-link
          to="/community"
          class="bg-pink-600 text-white px-4 py-2 rounded-lg hover:bg-pink-700 transition-colors"
        >
          返回社区
        </router-link>
      </div>

      <!-- 帖子内容 -->
      <div v-else-if="post" class="bg-white rounded-lg shadow-sm border border-gray-200">
        <!-- 帖子头部 -->
        <div class="p-6 border-b border-gray-200">
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <h1 class="text-2xl font-bold text-gray-900 mb-3">{{ post.title }}</h1>
              <div class="flex items-center space-x-4 text-sm text-gray-500">
                <div class="flex items-center space-x-2">
                  <UserAvatar :avatar="post.author?.avatar || ''" :name="post.author?.nickname || post.author?.username || 'U'" :size="32" />
                  <span class="font-medium">{{ post.author?.nickname || post.author?.username }}</span>
                </div>
                <span>发布于 {{ formatDate(post.created_at) }}</span>
                <span class="bg-pink-100 text-pink-800 px-2 py-1 rounded-full text-xs">
                  {{ getTopicLabel(post.topic) }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 帖子正文 -->
        <div class="p-6">
          <div class="prose max-w-none text-gray-800 leading-relaxed whitespace-pre-wrap">{{ post.content }}</div>
        </div>

        <!-- 互动按钮 -->
        <div class="px-6 py-4 border-t border-gray-200">
          <div class="flex items-center justify-between">
            <!-- 统计信息 -->
            <div class="flex items-center space-x-6 text-sm text-gray-600">
              <div class="flex items-center space-x-1">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                <span>{{ post.view_count }} 浏览</span>
              </div>
              <div class="flex items-center space-x-1">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
                <span>{{ post.reply_count }} 回复</span>
              </div>
            </div>

            <!-- 操作按钮 -->
            <div class="flex items-center space-x-3">
              <button
                @click="toggleLike"
                :disabled="isLiking"
                :class="[
                  'flex items-center space-x-1 px-3 py-2 rounded-lg transition-all duration-200',
                  isLiked
                    ? 'bg-pink-100 text-pink-600 hover:bg-pink-200'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                ]"
              >
                <svg class="h-4 w-4" :fill="isLiked ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                </svg>
                <span>{{ post.like_count }}</span>
              </button>

              <button
                v-if="!post.is_locked"
                @click="showReplyForm = !showReplyForm"
                class="flex items-center space-x-1 px-3 py-2 bg-pink-600 text-white rounded-lg hover:bg-pink-700 transition-colors"
              >
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
                <span>回复</span>
              </button>
              <!-- 锁定提示 -->
              <div v-if="post.is_locked" class="flex items-center space-x-1 px-3 py-2 bg-gray-200 text-gray-500 rounded-lg">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
                <span>帖子已锁定</span>
              </div>

              <!-- 举报按钮 -->
              <button
                @click="openReportModal"
                class="group flex items-center px-4 py-2 bg-white text-gray-700 border border-gray-300 rounded-md font-medium transition-all duration-300 text-sm min-w-[72px] justify-center hover:border-red-300 hover:text-red-600 hover:bg-red-50"
              >
                <svg class="h-4 w-4 mr-1.5" viewBox="0 0 24 24" fill="none" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 8a6 6 0 00-12 0v5a3 3 0 006 0V9m-6 8h12"/></svg>
                举报
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 回复表单 -->
      <div v-if="post && showReplyForm && !post.is_locked" class="mt-6 bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 class="text-lg font-medium text-gray-900 mb-4">发表回复</h3>
        <form @submit.prevent="submitReply">
          <textarea
            v-model="replyContent"
            placeholder="写下你的回复..."
            rows="4"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-pink-500 focus:border-transparent resize-none"
            required
          ></textarea>
          <div class="mt-3 flex justify-end space-x-3">
            <button
              type="button"
              @click="showReplyForm = false"
              class="px-4 py-2 text-gray-600 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors"
            >
              取消
            </button>
            <button
              type="submit"
              :disabled="isSubmitting || !replyContent.trim()"
              class="px-4 py-2 bg-pink-600 text-white rounded-lg hover:bg-pink-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ isSubmitting ? '发布中...' : '发布回复' }}
            </button>
          </div>
        </form>
      </div>


      <!-- 回复列表 -->
      <div v-if="post" class="mt-6 bg-white rounded-lg shadow-sm border border-gray-200">
        <div class="p-6 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900">全部回复 ({{ replies.length }})</h3>
        </div>
        <div v-if="replies.length > 0" class="divide-y divide-gray-200">
          <div v-for="reply in replies" :key="reply.id" class="p-6">
            <div class="flex space-x-3">
              <img
                :src="reply.author?.avatar || '/default-avatar.png'"
                :alt="reply.author?.nickname || reply.author?.username"
                class="w-10 h-10 rounded-full"
              />
              <div class="flex-1">
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center space-x-2">
                    <span class="font-medium text-gray-900">{{ reply.author?.nickname || reply.author?.username }}</span>
                    <span class="text-sm text-gray-500">{{ formatDate(reply.created_at) }}</span>
                  </div>
                  <!-- 删除按钮 -->
                  <button
                    v-if="canDeleteReply(reply)"
                    @click="handleDeleteReply(reply.id)"
                    class="flex items-center space-x-1 text-xs text-red-500 hover:text-red-700 transition-all duration-200 px-2 py-1 rounded-lg hover:bg-red-50"
                  >
                    <svg class="h-3 w-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                    <span>删除</span>
                  </button>
                </div>
                <p class="text-gray-800 whitespace-pre-wrap">{{ reply.content }}</p>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="p-12 text-center">
          <div class="text-gray-400 mb-2">
            <svg class="h-12 w-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
          </div>
          <p class="text-gray-500">暂无回复，来发第一个回复吧！</p>
        </div>
      </div>

      <!-- 返回按钮 -->
      <div v-if="post" class="mt-6 flex justify-center">
        <button
          @click="goBack"
          class="bg-gray-100 text-gray-700 px-6 py-2 rounded-lg hover:bg-gray-200 transition-colors font-medium"
        >
          返回社区
        </button>
      </div>
    </div>
    <!-- 举报模态框 -->
    <div v-if="showReport" class="fixed inset-0 bg-black/30 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-md shadow">
        <h3 class="text-lg font-semibold mb-4">举报帖子</h3>
        <div class="space-y-3">
          <div>
            <label class="block text-sm mb-1">举报理由</label>
            <select v-model="reportForm.reason" class="w-full border rounded px-3 py-2">
              <option value="">请选择</option>
              <option>垃圾广告</option>
              <option>不实信息</option>
              <option>辱骂/人身攻击</option>
              <option>涉黄/违法</option>
              <option>版权问题</option>
              <option>其他</option>
            </select>
          </div>
          <div>
            <label class="block text-sm mb-1">补充说明</label>
            <textarea v-model="reportForm.description" rows="3" class="w-full border rounded px-3 py-2" placeholder="请描述问题..." />
          </div>
          <div>
            <label class="block text-sm mb-1">证据链接（可选）</label>
            <input v-model="reportForm.evidence" class="w-full border rounded px-3 py-2" placeholder="https://" />
          </div>
        </div>
        <div class="mt-4 flex justify-end gap-2">
          <button @click="showReport = false" class="px-4 py-2 border rounded">取消</button>
          <button @click="submitReport" class="px-4 py-2 bg-red-600 text-white rounded">提交举报</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ForumApi } from '@/api/forum'
import { useToast } from '@/composables/useToast'
import { useAuthStore } from '@/stores/auth'
import BaseHeader from '@/components/BaseHeader.vue'
import UserAvatar from '@/components/UserAvatar.vue'
import type { ForumPost, ForumReply } from '@/api/types'
import ReportApi from '@/api/report'

const route = useRoute()
const router = useRouter()
const { showToast } = useToast()
const authStore = useAuthStore()

const post = ref<ForumPost | null>(null)
const replies = ref<ForumReply[]>([])
const isLoading = ref(false)
const isLiked = ref(false)
const isLiking = ref(false)
const showReplyForm = ref(false)
const replyContent = ref('')
const isSubmitting = ref(false)
const showReport = ref(false)
const reportForm = ref({ reason: '', description: '', evidence: '' })

// 话题标签映射
const topicLabels: Record<string, string> = {
  'All': '全部',
  'Baby Care': '婴儿护理',
  'Feeding': '喂养',
  'Sleep': '睡眠',
  'Health': '健康',
  'Development': '发育',
  'Activities': '活动',
  'Gear': '用品',
  'Parenting': '育儿',
  'Family Life': '家庭生活',
  'Work & Life Balance': '工作生活平衡',
  'Relationships': '人际关系',
  'Mental Health': '心理健康',
  'Finances': '财务',
  'Legal': '法律',
  'Other': '其他'
}

const getTopicLabel = (topic: string) => {
  return topicLabels[topic] || topic
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const loadPost = async () => {
  const postId = route.params.id as string
  if (!postId || isNaN(Number(postId))) {
    showToast('无效的帖子ID', 'error')
    router.push('/community')
    return
  }

  isLoading.value = true
  try {
    const response = await ForumApi.getPost(Number(postId))
    post.value = response.data

    // 增加浏览量
    ForumApi.incrementPostView(Number(postId)).catch(() => {})

    // 获取点赞状态（如果用户已登录）
    if (authStore.isAuthenticated) {
      try {
        const likeStatusResponse = await ForumApi.getPostLikeStatus(Number(postId))
        isLiked.value = likeStatusResponse.data?.is_liked || false
      } catch (error) {
        console.error('获取点赞状态失败:', error)
        isLiked.value = false
      }
    }
  } catch (error: any) {
    console.error('加载帖子失败:', error)
    if (error.response?.status === 404) {
      showToast('帖子不存在', 'error')
    } else {
      showToast('加载帖子失败', 'error')
    }
  } finally {
    isLoading.value = false
  }
}

const toggleLike = async () => {
  if (!authStore.isAuthenticated) {
    showToast('请先登录', 'warning')
    router.push('/login')
    return
  }

  if (!post.value || isLiking.value) return

  isLiking.value = true
  try {
    const response = await ForumApi.togglePostLike(post.value.id)

    // 根据后端返回的数据判断点赞状态
    // 后端逻辑：点赞成功返回点赞对象，取消点赞返回null
    if (response.data !== null && response.data !== undefined) {
      // 返回了点赞对象，说明是点赞成功
      isLiked.value = true
      post.value.like_count = (post.value.like_count || 0) + 1
    } else {
      // 返回null，说明是取消点赞成功
      isLiked.value = false
      post.value.like_count = Math.max(0, (post.value.like_count || 1) - 1)
    }
  } catch (error: any) {
    console.error('点赞操作失败:', error)
    showToast('操作失败', 'error')
  } finally {
    isLiking.value = false
  }
}

const submitReply = async () => {
  if (!authStore.isAuthenticated) {
    showToast('请先登录', 'warning')
    router.push('/login')
    return
  }

  if (!post.value || !replyContent.value.trim() || isSubmitting.value) return

  isSubmitting.value = true
  try {
    const response = await ForumApi.createReply({
      post_id: post.value.id,
      content: replyContent.value.trim()
    })

    // 重新加载回复列表以保证排序一致性
    await loadReplies()
    // 更新帖子回复数
    post.value.reply_count++

    // 重置表单
    replyContent.value = ''
    showReplyForm.value = false

    showToast('回复成功', 'success')
  } catch (error: any) {
    console.error('回复失败:', error)
    showToast('回复失败', 'error')
  } finally {
    isSubmitting.value = false
  }
}

const loadReplies = async () => {
  if (!post.value) return

  try {
    const response = await ForumApi.getPostReplies(post.value.id, {
      page: 1,
      size: 50,
      post_id: post.value.id,
      sort: 'created_at desc'
    })
    replies.value = response.data.items || []
  } catch (error: any) {
    console.error('加载回复失败:', error)
  }
}

// 检查是否可以删除回复
const canDeleteReply = (reply: ForumReply) => {
  if (!authStore.isAuthenticated) return false
  // 回复作者本人或管理员可以删除
  return authStore.user?.id === reply.author?.id || authStore.isAdmin
}

// 删除回复
const handleDeleteReply = async (replyId: number) => {
  if (!confirm('确定要删除这条回复吗？')) {
    return
  }

  try {
    await ForumApi.deleteReply(replyId)

    // 从列表中移除已删除的回复
    replies.value = replies.value.filter(reply => reply.id !== replyId)

    // 更新帖子的回复数量
    if (post.value) {
      post.value.reply_count = Math.max(0, (post.value.reply_count || 1) - 1)
    }

    showToast('删除回复成功', 'success')
  } catch (error: any) {
    console.error('删除回复失败:', error)
    showToast('删除回复失败', 'error')
  }
}

const goBack = () => {
  router.push('/community')
}

const openReportModal = () => {
  if (!authStore.isAuthenticated) {
    showToast('请先登录', 'warning')
    router.push('/login')
    return
  }
  showReport.value = true
}

const submitReport = async () => {
  if (!post.value) return
  try {
    const payload = {
      target_type: 'forum_post' as const,
      target_id: post.value.id,
      reason: reportForm.value.reason || '其他',
      description: reportForm.value.description,
      evidence: reportForm.value.evidence
    }
    await ReportApi.createReport(payload)
    showToast('举报已提交，我们会尽快核实', 'success')
    showReport.value = false
    reportForm.value = { reason: '', description: '', evidence: '' }
  } catch (e: any) {
    showToast(e.message || '提交失败，请重试', 'error')
  }
}

onMounted(async () => {
  await loadPost()
  if (post.value) {
    loadReplies()
  }
})
</script>
