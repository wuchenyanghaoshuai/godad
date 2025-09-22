<template>
  <div class="hot-articles">
    <!-- 头部标题 -->
    <div class="flex items-center justify-between mb-4">
      <div>
        <h2 class="text-xl font-bold text-gray-900">🔥 热门文章</h2>
        <p v-if="isFallback" class="text-xs text-gray-500 mt-1">暂无热门，已为你推荐最新文章</p>
      </div>
      <div class="flex space-x-2">
        <select
          v-model="selectedPeriod"
          @change="loadHotArticles"
          class="text-sm border border-gray-300 rounded-md px-3 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="today">今日</option>
          <option value="week">本周</option>
          <option value="month">本月</option>
          <option value="all">全部</option>
        </select>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="text-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
      <p class="mt-4 text-gray-500">加载中...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="text-center py-8 text-red-500">
      <p>{{ error }}</p>
      <button
        @click="loadHotArticles"
        class="mt-4 bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
      >
        重试
      </button>
    </div>

    <!-- 热门文章列表 -->
    <div v-else-if="articles && articles.length > 0" class="space-y-3">
      <div
        v-for="(article, index) in articles"
        :key="article.id"
        class="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-lg transition-all duration-200 cursor-pointer"
        @click="goToArticle(article.id)"
      >
        <div class="flex items-start space-x-3">
          <!-- 排名标识 -->
          <div class="flex-shrink-0">
            <span
              class="inline-flex items-center justify-center w-8 h-8 rounded-full text-sm font-bold"
              :class="{
                'bg-gradient-to-r from-yellow-400 to-yellow-600 text-white': index === 0,
                'bg-gradient-to-r from-gray-300 to-gray-500 text-white': index === 1,
                'bg-gradient-to-r from-orange-400 to-orange-600 text-white': index === 2,
                'bg-gray-100 text-gray-600': index > 2
              }"
            >
              {{ index + 1 }}
            </span>
          </div>

          <!-- 文章信息 -->
          <div class="flex-1 min-w-0">
            <!-- 标题 -->
            <h3 class="font-medium text-gray-900 line-clamp-2 mb-2 hover:text-blue-600 transition-colors">
              {{ article.title }}
            </h3>

            <!-- 元信息 -->
            <div class="flex items-center space-x-4 text-sm text-gray-500">
              <span class="flex items-center">
                <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                {{ article.view_count }}
              </span>
              <span class="flex items-center">
                <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                </svg>
                {{ article.like_count }}
              </span>
              <span class="flex items-center">
                <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
                {{ article.comment_count }}
              </span>
              <span>{{ article.author?.username }}</span>
              <span>{{ formatDate(article.created_at) }}</span>
            </div>

            <!-- 分类标签 -->
            <div v-if="article.category" class="mt-2">
              <span class="inline-block bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded-full">
                {{ article.category.name }}
              </span>
            </div>
          </div>

          <!-- 热度值 -->
          <div class="flex-shrink-0 text-right">
            <div class="text-lg font-bold text-orange-600">
              {{ calculateHeatScore(article) }}
            </div>
            <div class="text-xs text-gray-500">热度值</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="text-center py-8 text-gray-500">
      <svg class="w-16 h-16 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      <p>暂无热门文章</p>
      <p class="text-sm mt-2">{{ getPeriodText(selectedPeriod) }}还没有热门文章</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArticleApi, type Article } from '@/api'

// Props
interface Props {
  limit?: number
  showPeriodSelector?: boolean
  defaultPeriod?: 'today' | 'week' | 'month' | 'all'
}

const props = withDefaults(defineProps<Props>(), {
  limit: 10,
  showPeriodSelector: true,
  defaultPeriod: 'today'
})

// 响应式数据
const articles = ref<Article[]>([])
const loading = ref(false)
const error = ref<string>('')
const selectedPeriod = ref(props.defaultPeriod)
const isFallback = ref(false)

// Router
const router = useRouter()

// 方法
const loadHotArticles = async () => {
  try {
    loading.value = true
    error.value = ''
    isFallback.value = false

    const response = await ArticleApi.getHotArticles({
      period: selectedPeriod.value,
      limit: props.limit
    })

    const hot = Array.isArray(response.data) ? response.data : []
    if (hot.length > 0) {
      articles.value = hot
    } else {
      // Fallback: 无热门数据时，取最新已发布文章随机填充
      const pageResp = await ArticleApi.getArticlePage({ page: 1, size: Math.max(12, props.limit), status: 1 })
      const items = Array.isArray(pageResp.data?.items) ? pageResp.data.items : []
      articles.value = pickRandom(items, Math.min(props.limit, items.length))
      isFallback.value = articles.value.length > 0
      if (!isFallback.value) {
        error.value = '暂无热门文章'
      }
    }
  } catch (err: any) {
    // 出错时也尝试 fallback
    try {
      const pageResp = await ArticleApi.getArticlePage({ page: 1, size: Math.max(12, props.limit), status: 1 })
      const items = Array.isArray(pageResp.data?.items) ? pageResp.data.items : []
      articles.value = pickRandom(items, Math.min(props.limit, items.length))
      isFallback.value = articles.value.length > 0
      if (!isFallback.value) {
        error.value = err?.message || '加载热门文章失败'
      } else {
        error.value = ''
      }
    } catch (e: any) {
      console.error('加载热门文章失败:', err)
      error.value = e?.message || err?.message || '加载热门文章失败'
    }
  } finally {
    loading.value = false
  }
}

const goToArticle = (articleId: number) => {
  router.push(`/articles/${articleId}`)
}

const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  const now = new Date()
  const diffInMinutes = Math.floor((now.getTime() - date.getTime()) / (1000 * 60))

  if (diffInMinutes < 60) {
    return `${diffInMinutes}分钟前`
  } else if (diffInMinutes < 1440) {
    const hours = Math.floor(diffInMinutes / 60)
    return `${hours}小时前`
  } else if (diffInMinutes < 10080) {
    const days = Math.floor(diffInMinutes / 1440)
    return `${days}天前`
  } else {
    return date.toLocaleDateString('zh-CN', {
      month: 'short',
      day: 'numeric'
    })
  }
}

const calculateHeatScore = (article: Article): number => {
  // 使用与后端相同的算法计算热度值
  return Math.round((article.view_count || 0) * 0.6 + (article.like_count || 0) * 0.3 + (article.comment_count || 0) * 0.1)
}

const getPeriodText = (period: string): string => {
  const periodMap: Record<string, string> = {
    today: '今日',
    week: '本周',
    month: '本月',
    all: '全部时间'
  }
  return periodMap[period] || period
}

// 暴露方法给父组件
defineExpose({
  refresh: loadHotArticles,
  articles
})

// 生命周期
onMounted(() => {
  loadHotArticles()
})

function pickRandom<T>(arr: T[], n: number): T[] {
  if (!arr || arr.length === 0) return []
  if (n >= arr.length) return arr.slice()
  const a = arr.slice()
  for (let i = a.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[a[i], a[j]] = [a[j], a[i]]
  }
  return a.slice(0, n)
}
</script>

<style scoped>
.line-clamp-2 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.hot-articles {
  @apply select-none;
}
</style>
