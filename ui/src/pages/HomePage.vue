<template>
  <AppLayout
    background-class="bg-gradient-to-br from-pink-50 to-orange-50"
    :show-footer="true"
  >
    <!-- 英雄区域 -->
    <PageContainer
      background="gradient"
      padding="none"
      :show-max-width="false"
      custom-class="relative min-h-screen flex items-center justify-center overflow-hidden"
    >
      <!-- 装饰性背景元素 -->
      <div class="absolute inset-0 overflow-hidden">
        <div class="absolute -top-20 sm:-top-40 -right-20 sm:-right-40 w-40 sm:w-80 h-40 sm:h-80 bg-gradient-to-br from-pink-400 to-orange-400 rounded-full opacity-20 animate-pulse"></div>
        <div class="absolute -bottom-20 sm:-bottom-40 -left-20 sm:-left-40 w-48 sm:w-96 h-48 sm:h-96 bg-gradient-to-tr from-orange-400 to-yellow-400 rounded-full opacity-20 animate-pulse delay-1000"></div>
        <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-32 sm:w-64 h-32 sm:h-64 bg-gradient-to-r from-pink-300 to-orange-300 rounded-full opacity-10 animate-spin" style="animation-duration: 20s;"></div>
      </div>

      <div class="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <h1 class="text-3xl sm:text-5xl lg:text-7xl font-bold text-gray-900 mb-4 sm:mb-6 leading-tight">
          <span class="bg-gradient-to-r from-pink-600 to-orange-600 bg-clip-text text-transparent">GoDad</span>
          <br class="hidden sm:block" />
          <span class="text-gray-800 text-2xl sm:text-4xl lg:text-6xl">育儿知识分享平台</span>
        </h1>

        <p class="text-base sm:text-lg lg:text-2xl text-gray-600 mb-6 sm:mb-8 max-w-3xl mx-auto leading-relaxed px-2">
          专业的育儿指导，温暖的社区支持，陪伴每一个家庭的成长之路
        </p>
        <div class="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center items-center px-4 sm:px-0">
          <router-link
            to="/articles"
            class="w-full sm:w-auto inline-flex items-center justify-center bg-gradient-to-r from-pink-600 to-orange-600 text-white px-6 sm:px-8 py-3 sm:py-4 rounded-full text-base sm:text-lg font-semibold hover:shadow-xl transform hover:scale-105 transition-all duration-300 shadow-lg"
          >
            开始阅读
            <ArrowRightIcon class="ml-2 h-4 w-4 sm:h-5 sm:w-5" />
          </router-link>
          <router-link
            to="/register"
            class="w-full sm:w-auto inline-flex items-center justify-center border-2 border-pink-600 text-pink-600 px-6 sm:px-8 py-3 sm:py-4 rounded-full text-base sm:text-lg font-semibold hover:bg-pink-50 hover:shadow-lg transition-all duration-300 bg-white/80 backdrop-blur-sm"
          >
            加入社区
            <UsersIcon class="ml-2 h-4 w-4 sm:h-5 sm:w-5" />
          </router-link>
        </div>

        <!-- 统计数据预览 -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-3 sm:gap-4 md:gap-6 mt-12 sm:mt-16 px-4 max-w-4xl mx-auto">
          <div class="text-center bg-white/80 backdrop-blur-sm rounded-lg sm:rounded-xl p-3 sm:p-4 md:p-6 border border-white/50 hover:bg-white/90 transition-all duration-300 shadow-lg">
            <div class="text-xl sm:text-2xl md:text-3xl font-bold text-pink-600 mb-1 sm:mb-2">10K+</div>
            <div class="text-xs sm:text-sm md:text-base text-gray-600 font-medium">注册用户</div>
          </div>
          <div class="text-center bg-white/80 backdrop-blur-sm rounded-lg sm:rounded-xl p-3 sm:p-4 md:p-6 border border-white/50 hover:bg-white/90 transition-all duration-300 shadow-lg">
            <div class="text-xl sm:text-2xl md:text-3xl font-bold text-orange-600 mb-1 sm:mb-2">500+</div>
            <div class="text-xs sm:text-sm md:text-base text-gray-600 font-medium">优质文章</div>
          </div>
          <div class="text-center bg-white/80 backdrop-blur-sm rounded-lg sm:rounded-xl p-3 sm:p-4 md:p-6 border border-white/50 hover:bg-white/90 transition-all duration-300 shadow-lg">
            <div class="text-xl sm:text-2xl md:text-3xl font-bold text-purple-600 mb-1 sm:mb-2">50+</div>
            <div class="text-xs sm:text-sm md:text-base text-gray-600 font-medium">专业作者</div>
          </div>
          <div class="text-center bg-white/80 backdrop-blur-sm rounded-lg sm:rounded-xl p-3 sm:p-4 md:p-6 border border-white/50 hover:bg-white/90 transition-all duration-300 shadow-lg">
            <div class="text-xl sm:text-2xl md:text-3xl font-bold text-green-600 mb-1 sm:mb-2">24/7</div>
            <div class="text-xs sm:text-sm md:text-base text-gray-600 font-medium">在线支持</div>
          </div>
        </div>
      </div>
    </PageContainer>

    <!-- 特色功能 -->
    <PageContainer background="white" padding="xl">
      <div class="text-center mb-8 sm:mb-12 md:mb-16">
        <h2 class="text-2xl sm:text-3xl md:text-4xl font-bold text-gray-900 mb-3 sm:mb-4">
          为什么选择 <span class="bg-gradient-to-r from-pink-600 to-orange-600 bg-clip-text text-transparent">GoDad</span>
        </h2>
        <p class="text-base sm:text-lg md:text-xl text-gray-600 max-w-3xl mx-auto px-4">
          我们致力于为每一位父母提供最专业、最贴心的育儿支持
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 sm:gap-8">
        <!-- 专业可靠 -->
        <div class="text-center p-6 sm:p-8 rounded-xl sm:rounded-2xl bg-gradient-to-br from-pink-50 to-orange-50 hover:shadow-xl hover:scale-105 transition-all duration-300 group transform hover:-translate-y-1">
          <div class="w-12 h-12 sm:w-16 sm:h-16 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center mx-auto mb-4 sm:mb-6 group-hover:scale-110 transition-transform duration-300 shadow-lg">
            <HeartIcon class="h-6 w-6 sm:h-8 sm:w-8 text-white" />
          </div>
          <h3 class="text-xl sm:text-2xl font-bold text-gray-900 mb-3 sm:mb-4">专业可靠</h3>
          <p class="text-gray-600 leading-relaxed text-sm sm:text-base">
            汇聚儿科医生、营养师、教育专家的专业知识，确保每一条建议都科学可靠
          </p>
        </div>

        <!-- 温暖社区 -->
        <div class="text-center p-6 sm:p-8 rounded-xl sm:rounded-2xl bg-gradient-to-br from-blue-50 to-purple-50 hover:shadow-xl hover:scale-105 transition-all duration-300 group transform hover:-translate-y-1">
          <div class="w-12 h-12 sm:w-16 sm:h-16 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full flex items-center justify-center mx-auto mb-4 sm:mb-6 group-hover:scale-110 transition-transform duration-300 shadow-lg">
            <UsersIcon class="h-6 w-6 sm:h-8 sm:w-8 text-white" />
          </div>
          <h3 class="text-xl sm:text-2xl font-bold text-gray-900 mb-3 sm:mb-4">温暖社区</h3>
          <p class="text-gray-600 leading-relaxed text-sm sm:text-base">
            连接全球父母，分享育儿经验，在这里你永远不会感到孤单
          </p>
        </div>

        <!-- 持续更新 -->
        <div class="text-center p-6 sm:p-8 rounded-xl sm:rounded-2xl bg-gradient-to-br from-green-50 to-teal-50 hover:shadow-xl hover:scale-105 transition-all duration-300 group transform hover:-translate-y-1">
          <div class="w-12 h-12 sm:w-16 sm:h-16 bg-gradient-to-r from-green-500 to-teal-500 rounded-full flex items-center justify-center mx-auto mb-4 sm:mb-6 group-hover:scale-110 transition-transform duration-300 shadow-lg">
            <TrendingUpIcon class="h-6 w-6 sm:h-8 sm:w-8 text-white" />
          </div>
          <h3 class="text-xl sm:text-2xl font-bold text-gray-900 mb-3 sm:mb-4">持续更新</h3>
          <p class="text-gray-600 leading-relaxed text-sm sm:text-base">
            紧跟最新育儿研究成果，定期更新内容，让您的育儿知识始终保持前沿
          </p>
        </div>
      </div>
    </PageContainer>

    <!-- 统计数据（浅色渐变背景 + 白卡片） -->
    <PageContainer
      background="tint"
      padding="xl"
      :show-max-width="false"
      custom-class="relative overflow-hidden"
    >
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-8 sm:mb-12">
          <h2 class="text-2xl sm:text-3xl md:text-4xl font-bold text-gray-900 mb-3">
            数据见证我们的成长
          </h2>
          <p class="text-base sm:text-lg text-gray-600 max-w-2xl mx-auto">
            每一个数字背后都是我们对育儿事业的坚持与热爱
          </p>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-6">
          <div class="rounded-xl border border-gray-100 bg-white p-5 sm:p-6 shadow-sm hover:shadow-md transition-all duration-300 text-center">
            <div class="text-3xl sm:text-4xl font-bold text-pink-600">10K+</div>
            <div class="mt-1 text-sm sm:text-base text-gray-700 font-medium">注册用户</div>
            <div class="mt-0.5 text-xs text-gray-500">信任之选</div>
          </div>

          <div class="rounded-xl border border-gray-100 bg-white p-5 sm:p-6 shadow-sm hover:shadow-md transition-all duration-300 text-center">
            <div class="text-3xl sm:text-4xl font-bold text-orange-600">5K+</div>
            <div class="mt-1 text-sm sm:text-base text-gray-700 font-medium">优质文章</div>
            <div class="mt-0.5 text-xs text-gray-500">知识宝库</div>
          </div>

          <div class="rounded-xl border border-gray-100 bg-white p-5 sm:p-6 shadow-sm hover:shadow-md transition-all duration-300 text-center">
            <div class="text-3xl sm:text-4xl font-bold text-purple-600">100+</div>
            <div class="mt-1 text-sm sm:text-base text-gray-700 font-medium">专业作者</div>
            <div class="mt-0.5 text-xs text-gray-500">权威专家</div>
          </div>

          <div class="rounded-xl border border-gray-100 bg-white p-5 sm:p-6 shadow-sm hover:shadow-md transition-all duration-300 text-center">
            <div class="text-3xl sm:text-4xl font-bold text-green-600">24/7</div>
            <div class="mt-1 text-sm sm:text-base text-gray-700 font-medium">在线支持</div>
            <div class="mt-0.5 text-xs text-gray-500">贴心服务</div>
          </div>
        </div>
      </div>
    </PageContainer>

    <!-- 热门文章 - 只在有数据时显示 -->
    <PageContainer v-if="hasHotArticles" background="gray" padding="xl">
      <HotArticles ref="hotArticlesRef" />
    </PageContainer>

    <!-- 最新文章预览 -->
    <PageContainer background="white" padding="xl">
      <div class="text-center mb-8 sm:mb-12 md:mb-16">
        <h2 class="text-2xl sm:text-3xl md:text-4xl font-bold text-gray-900 mb-3 sm:mb-4">最新文章</h2>
        <p class="text-base sm:text-lg md:text-xl text-gray-600 max-w-3xl mx-auto px-4 sm:px-0">
          发现最新的育儿知识和经验分享
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 sm:gap-8 mb-8 sm:mb-12">
        <article
          v-for="(a, idx) in weeklyHot3"
          :key="a.id"
          class="bg-white rounded-xl sm:rounded-2xl shadow-lg overflow-hidden hover:shadow-2xl hover:scale-105 transition-all duration-300 group cursor-pointer"
          @click="$router.push(`/articles/${a.id}`)"
        >
          <div class="h-40 sm:h-48 relative overflow-hidden">
            <img
              v-if="a.cover_image && !(a as any)._imageError"
              :src="a.cover_image"
              :alt="a.title"
              loading="lazy"
              class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
              @error="handleImageError(a)"
              @load="handleImageLoad(a)"
            />
            <!-- 无封面或加载失败：采用与 /articles 同风格的占位配色 -->
            <div v-else class="w-full h-full bg-gradient-warm"></div>
            <div class="absolute inset-0 bg-black/0 group-hover:bg-black/10 transition-all duration-300"></div>
            <div class="absolute top-3 sm:top-4 left-3 sm:left-4">
              <span class="bg-white/90 text-pink-800 text-xs font-medium px-2 sm:px-3 py-1 rounded-full backdrop-blur-sm">本周热门</span>
            </div>
          </div>
          <div class="p-4 sm:p-6">
            <div class="flex items-center mb-2 sm:mb-3">
              <span v-if="a.category?.name" class="bg-pink-100 text-pink-800 text-xs font-medium px-2 sm:px-2.5 py-0.5 rounded-full">{{ a.category.name }}</span>
              <span class="text-gray-500 text-xs sm:text-sm ml-auto">{{ formatDate(a.created_at) }}</span>
            </div>
            <h3 class="text-base sm:text-lg md:text-xl font-bold text-gray-900 mb-2 sm:mb-3 single-line-ellipsis group-hover:text-pink-600 transition-colors duration-300" :title="a.title">
              {{ a.title }}
            </h3>
            <p class="text-gray-600 mb-3 sm:mb-4 line-clamp-2 text-xs sm:text-sm md:text-base">
              {{ getHomePreviewSummary(a) }}
            </p>
            <div class="flex items-center justify-between">
              <div class="flex items-center">
                <div class="w-6 h-6 sm:w-8 sm:h-8 bg-gradient-to-r from-pink-500 to-orange-500 rounded-full flex items-center justify-center shadow-md">
                  <span class="text-white text-xs sm:text-sm font-medium">{{ (a.author?.nickname || a.author?.username || 'U').slice(0,1).toUpperCase() }}</span>
                </div>
                <span class="text-gray-700 text-xs sm:text-sm ml-2 font-medium">{{ a.author?.nickname || a.author?.username || '匿名' }}</span>
              </div>
              <div class="flex items-center text-gray-500 text-xs sm:text-sm">
                <HeartIcon class="h-3 w-3 sm:h-4 sm:w-4 mr-1 group-hover:text-pink-500 transition-colors duration-300" />
                <span>{{ a.like_count || 0 }}</span>
              </div>
            </div>
          </div>
        </article>
      </div>

      <div class="text-center">
        <router-link
          to="/articles"
          class="inline-flex items-center px-4 sm:px-6 md:px-8 py-2 sm:py-3 bg-gradient-to-r from-pink-600 to-orange-600 text-white font-medium rounded-full hover:from-pink-700 hover:to-orange-700 transition-all duration-300 hover:scale-105 shadow-lg hover:shadow-xl text-xs sm:text-sm md:text-base"
        >
          查看更多文章
          <ArrowRightIcon class="ml-1 sm:ml-2 h-3 w-3 sm:h-4 sm:w-4 md:h-5 md:w-5" />
        </router-link>
      </div>
    </PageContainer>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, defineAsyncComponent } from 'vue'
import { ArrowRightIcon, UsersIcon, HeartIcon, TrendingUpIcon } from 'lucide-vue-next'
import { AppLayout, PageContainer } from '@/components/layout'
import { ArticleApi, type Article } from '@/api'

// 异步加载热门文章组件
const HotArticles = defineAsyncComponent({
  loader: () => import('@/components/HotArticles.vue'),
  loadingComponent: { template: '<div class="animate-pulse bg-gray-200 h-64 rounded-lg"></div>' },
  errorComponent: { template: '<div class="text-center text-gray-500 py-8">加载失败</div>' },
  delay: 200,
  timeout: 3000
})

// 热门文章相关
const hotArticlesRef = ref()
const hasHotArticles = computed(() => {
  return hotArticlesRef.value?.articles && hotArticlesRef.value.articles.length > 0
})

// 本周热门 Top3 用于“最新文章预览”
const weeklyHot3 = ref<Article[]>([])

const gradientByIndex = (idx: number) => {
  const gradients = [
    'bg-gradient-to-br from-pink-400 to-orange-400',
    'bg-gradient-to-br from-blue-400 to-purple-400',
    'bg-gradient-to-br from-green-400 to-teal-400'
  ]
  return gradients[idx % gradients.length]
}

const formatDate = (iso: string | Date | undefined) => {
  if (!iso) return ''
  const d = typeof iso === 'string' ? new Date(iso) : iso
  if (isNaN(d.getTime())) return ''
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

onMounted(async () => {
  try {
    const resp = await ArticleApi.getHotArticles({ period: 'week', limit: 3 })
    weeklyHot3.value = Array.isArray(resp.data) ? resp.data : []
    // Fallback: no hot articles -> take random 3 from latest published
    if (!weeklyHot3.value || weeklyHot3.value.length === 0) {
      const pageResp = await ArticleApi.getArticlePage({ page: 1, size: 12, status: 1 })
      const items = Array.isArray(pageResp.data?.items) ? pageResp.data.items : []
      if (items.length > 0) {
        weeklyHot3.value = pickRandom(items, Math.min(3, items.length))
      }
    }
  } catch (e) {
    // On error, still attempt fallback to latest articles
    try {
      const pageResp = await ArticleApi.getArticlePage({ page: 1, size: 12, status: 1 })
      const items = Array.isArray(pageResp.data?.items) ? pageResp.data.items : []
      weeklyHot3.value = pickRandom(items, Math.min(3, items.length))
    } catch {
      weeklyHot3.value = []
    }
  }
})

// 图片加载处理方法
const handleImageError = (article: any) => {
  article._imageError = true
  console.warn(`Failed to load image for article: ${article.title}`)
}

const handleImageLoad = (article: any) => {
  article._imageLoaded = true
}

// 统一PC摘要：两行显示，无摘要时从正文提取前90字
const getHomePreviewSummary = (article: Article) => {
  const limit = 90
  const s = (article.summary || '').trim()
  if (s) return s.length > limit ? s.slice(0, limit) + '...' : s
  const raw = (article as any).content || ''
  const plain = raw.replace(/<[^>]*>/g, '')
  return plain.length > limit ? plain.slice(0, limit) + '...' : plain
}

function pickRandom<T>(arr: T[], n: number): T[] {
  if (n >= arr.length) return arr.slice()
  // Fisher-Yates shuffle partial
  const a = arr.slice()
  for (let i = a.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1))
    ;[a[i], a[j]] = [a[j], a[i]]
  }
  return a.slice(0, n)
}
</script>

<style scoped>
.single-line-ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
