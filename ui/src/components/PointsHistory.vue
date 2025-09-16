<template>
  <div class="space-y-4">
    <!-- 加载状态 -->
    <div v-if="loading" class="text-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
      <p class="mt-4 text-gray-500">加载中...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="text-center py-8 text-red-500">
      <p>{{ error }}</p>
      <button
        @click="loadHistory"
        class="mt-4 bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
      >
        重试
      </button>
    </div>

    <!-- 积分记录列表 -->
    <div v-else-if="transactions.length > 0" class="space-y-3">
      <div
        v-for="transaction in transactions"
        :key="transaction.id"
        class="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
      >
        <div class="flex items-center justify-between">
          <div class="flex-1">
            <div class="flex items-center space-x-3">
              <span class="flex items-center justify-center w-8 h-8 rounded-full text-sm font-medium"
                    :class="transaction.points > 0 ? 'bg-green-100 text-green-600' : 'bg-red-100 text-red-600'">
                {{ transaction.points > 0 ? '+' : '' }}{{ Math.abs(transaction.points) }}
              </span>
              <div>
                <h4 class="font-medium text-gray-900">{{ transaction.description }}</h4>
                <p class="text-sm text-gray-500">
                  操作类型: {{ getActionName(transaction.action) }}
                  <span v-if="transaction.source_type && transaction.source_id" class="ml-2">
                    来源: {{ transaction.source_type }}#{{ transaction.source_id }}
                  </span>
                </p>
              </div>
            </div>
          </div>
          <div class="text-right">
            <div class="text-lg font-bold" :class="transaction.points > 0 ? 'text-green-600' : 'text-red-600'">
              {{ transaction.points > 0 ? '+' : '' }}{{ transaction.points }}
            </div>
            <div class="text-xs text-gray-500">
              {{ formatDate(transaction.created_at) }}
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="pagination.total_pages > 1" class="flex justify-center mt-6">
        <nav class="flex items-center space-x-2">
          <button
            @click="loadHistory(pagination.page - 1)"
            :disabled="pagination.page <= 1 || loading"
            class="px-3 py-2 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            上一页
          </button>
          <span class="text-sm text-gray-600">
            第 {{ pagination.page }} 页，共 {{ pagination.total_pages }} 页
          </span>
          <button
            @click="loadHistory(pagination.page + 1)"
            :disabled="pagination.page >= pagination.total_pages || loading"
            class="px-3 py-2 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            下一页
          </button>
        </nav>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="text-center py-8 text-gray-500">
      <p>暂无积分记录</p>
      <p class="text-sm mt-2">开始使用系统功能来获得积分吧</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { PointsAPI, type PointsTransaction } from '@/api'

// 响应式数据
const transactions = ref<PointsTransaction[]>([])
const pagination = ref({
  page: 1,
  limit: 10,
  total: 0,
  total_pages: 0
})
const loading = ref(false)
const error = ref<string>('')

// 方法
const loadHistory = async (page = 1) => {
  try {
    loading.value = true
    error.value = ''

    const response = await PointsAPI.getPointsHistory({
      page,
      limit: pagination.value.limit
    })

    transactions.value = response.data.transactions
    pagination.value = response.data.pagination
  } catch (err: any) {
    error.value = err.message || '加载积分记录失败'
    console.error('加载积分记录失败:', err)
  } finally {
    loading.value = false
  }
}

const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  const now = new Date()
  const diffInMinutes = Math.floor((now.getTime() - date.getTime()) / (1000 * 60))

  if (diffInMinutes < 1) {
    return '刚刚'
  } else if (diffInMinutes < 60) {
    return `${diffInMinutes}分钟前`
  } else if (diffInMinutes < 1440) {
    const hours = Math.floor(diffInMinutes / 60)
    return `${hours}小时前`
  } else if (diffInMinutes < 10080) {
    const days = Math.floor(diffInMinutes / 1440)
    return `${days}天前`
  } else {
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }
}

const getActionName = (action: string): string => {
  const actionMap: Record<string, string> = {
    'login': '每日登录',
    'publish_article': '发布文章',
    'article_liked': '文章被点赞',
    'article_unliked': '文章被取消点赞',
    'comment_liked': '评论被点赞',
    'publish_comment': '发布评论',
    'share_article': '分享文章',
    'read_article': '阅读文章',
    'invite_user': '邀请新用户',
    'complete_profile': '完善资料',
    'first_post': '首次发布',
    'admin_award': '管理员奖励',
    'penalty': '违规扣分',
    'daily_login': '每日登录',
    'new_follower': '获得新粉丝'
  }
  return actionMap[action] || action
}

// 生命周期
onMounted(() => {
  loadHistory()
})
</script>

<style scoped>
/* 可以添加一些自定义样式 */
</style>