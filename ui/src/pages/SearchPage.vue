<template>
  <div class="min-h-screen bg-gray-50">
    <Navbar />
    
    <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
      <!-- 搜索页面标题 -->
      <div class="text-center mb-8">
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-900 mb-4">全站搜索</h1>
        <p class="text-gray-600">搜索文章、用户和分类内容</p>
      </div>

      <!-- 搜索框 -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-8">
        <div class="flex flex-col sm:flex-row gap-4">
          <div class="flex-1 relative">
            <SearchIcon class="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
            <input
              v-model="searchQuery"
              @keyup.enter="performSearch"
              type="text"
              placeholder="输入关键词搜索文章内容..."
              class="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm sm:text-base"
            />
          </div>
          <div class="flex gap-2">
            <select 
              v-model="searchType" 
              class="px-3 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
            >
              <option value="articles">文章</option>
              <option value="users">用户</option>
              <option value="all">全部</option>
            </select>
            <button
              @click="performSearch"
              :disabled="!searchQuery.trim() || isSearching"
              class="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed text-sm font-medium"
            >
              <LoaderIcon v-if="isSearching" class="h-4 w-4 animate-spin" />
              <span v-else>搜索</span>
            </button>
          </div>
        </div>
        
        <!-- 搜索过滤器 -->
        <div v-if="searchQuery" class="mt-4 flex flex-wrap gap-2">
          <select v-model="selectedCategory" @change="performSearch" class="px-3 py-1 text-sm border border-gray-300 rounded">
            <option value="">所有分类</option>
            <option v-for="category in categories" :key="category.id" :value="category.id">
              {{ category.name }}
            </option>
          </select>
          <select v-model="sortBy" @change="performSearch" class="px-3 py-1 text-sm border border-gray-300 rounded">
            <option value="created_at">最新发布</option>
            <option value="view_count">浏览量</option>
            <option value="like_count">点赞数</option>
            <option value="comment_count">评论数</option>
          </select>
        </div>
      </div>

      <!-- 搜索结果统计 -->
      <div v-if="hasSearched" class="mb-6">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between bg-blue-50 border border-blue-200 rounded-lg p-4">
          <div class="text-sm text-blue-700">
            <template v-if="isSearching">
              <LoaderIcon class="inline h-4 w-4 animate-spin mr-1" />
              正在搜索...
            </template>
            <template v-else-if="searchResults.total > 0">
              找到 <span class="font-semibold">{{ searchResults.total }}</span> 条关于 
              "<span class="font-semibold">{{ currentSearchQuery }}</span>" 的结果
            </template>
            <template v-else>
              没有找到关于 "<span class="font-semibold">{{ currentSearchQuery }}</span>" 的结果
            </template>
          </div>
          <div class="text-xs text-blue-600 mt-2 sm:mt-0">
            搜索时间: {{ searchTime }}ms
          </div>
        </div>
      </div>

      <!-- 搜索结果 -->
      <div v-if="hasSearched && !isSearching">
        <!-- 文章结果 -->
        <div v-if="searchResults.total > 0" class="space-y-6">
          <article
            v-for="article in searchResults.items"
            :key="article.id"
            class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 hover:shadow-md transition-all duration-200 cursor-pointer"
            @click="goToArticle(article.id)"
          >
            <!-- 文章头部 -->
            <div class="flex flex-col sm:flex-row sm:items-start justify-between mb-4">
              <div class="flex-1">
                <div class="flex items-center space-x-2 mb-2">
                  <span 
                    v-if="article.category"
                    class="px-2 py-1 text-xs bg-blue-100 text-blue-600 rounded-full"
                  >
                    {{ article.category.name }}
                  </span>
                  <span class="text-xs text-gray-500">{{ formatDate(article.created_at) }}</span>
                </div>
                <h2 class="text-lg sm:text-xl font-semibold text-gray-900 mb-2 hover:text-blue-600 transition-colors">
                  <span v-html="highlightKeyword(article.title)"></span>
                </h2>
                <p class="text-gray-600 text-sm leading-relaxed mb-3">
                  <span v-html="highlightKeyword(getArticleSummary(article))"></span>
                </p>
              </div>
              <div v-if="article.cover_image" class="ml-4 flex-shrink-0">
                <img 
                  :src="article.cover_image" 
                  :alt="article.title"
                  class="w-20 h-20 sm:w-24 sm:h-24 object-cover rounded-lg"
                />
              </div>
            </div>
            
            <!-- 文章元信息 -->
            <div class="flex flex-wrap items-center justify-between text-xs text-gray-500 border-t pt-3">
              <div class="flex items-center space-x-4">
                <div class="flex items-center space-x-1">
                  <UserIcon class="h-3 w-3" />
                  <router-link
                    v-if="article.author?.username"
                    :to="`/users/${article.author.username}`"
                    @click.stop
                    class="hover:text-pink-600 transition-colors cursor-pointer"
                  >
                    {{ article.author.username }}
                  </router-link>
                  <span v-else>{{ article.author?.username }}</span>
                </div>
                <div class="flex items-center space-x-1">
                  <EyeIcon class="h-3 w-3" />
                  <span>{{ article.view_count || 0 }}</span>
                </div>
                <div class="flex items-center space-x-1">
                  <HeartIcon class="h-3 w-3" />
                  <span>{{ article.like_count || 0 }}</span>
                </div>
                <div class="flex items-center space-x-1">
                  <MessageCircleIcon class="h-3 w-3" />
                  <span>{{ article.comment_count || 0 }}</span>
                </div>
              </div>
              <div class="text-blue-600 hover:text-blue-800">
                阅读更多 →
              </div>
            </div>
          </article>
          
          <!-- 分页 -->
          <div v-if="searchResults.total_pages > 1" class="flex justify-center mt-8">
            <nav class="flex items-center space-x-2">
              <button
                @click="goToPage(currentPage - 1)"
                :disabled="currentPage <= 1"
                class="px-3 py-2 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                上一页
              </button>
              <span class="text-sm text-gray-600">
                第 {{ currentPage }} 页，共 {{ searchResults.total_pages }} 页
              </span>
              <button
                @click="goToPage(currentPage + 1)"
                :disabled="currentPage >= searchResults.total_pages"
                class="px-3 py-2 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                下一页
              </button>
            </nav>
          </div>
        </div>

        <!-- 空搜索结果 -->
        <div v-else class="text-center py-12">
          <SearchIcon class="h-16 w-16 text-gray-300 mx-auto mb-4" />
          <h3 class="text-lg font-medium text-gray-900 mb-2">没有找到相关内容</h3>
          <p class="text-gray-600 mb-6">
            尝试使用不同的关键词或检查拼写错误
          </p>
          <div class="space-y-2 text-sm text-gray-500">
            <p>搜索建议：</p>
            <ul class="space-y-1">
              <li>• 使用更通用的关键词</li>
              <li>• 检查拼写和语法</li>
              <li>• 尝试相关的同义词</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- 热门搜索词 -->
      <div v-if="!hasSearched" class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">热门搜索</h3>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="keyword in popularKeywords"
            :key="keyword"
            @click="searchWithKeyword(keyword)"
            class="px-3 py-1 text-sm bg-gray-100 hover:bg-blue-100 hover:text-blue-700 text-gray-700 rounded-full transition-colors"
          >
            {{ keyword }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  SearchIcon,
  LoaderIcon,
  UserIcon,
  EyeIcon,
  HeartIcon,
  MessageCircleIcon
} from 'lucide-vue-next'
import { ArticleApi } from '@/api/article'
import { CategoryApi } from '@/api/category'
import type { Article, Category, PaginatedResponse } from '@/api/types'
import Navbar from '@/components/Navbar.vue'

const route = useRoute()
const router = useRouter()

// 响应式数据
const searchQuery = ref('')
const currentSearchQuery = ref('')
const searchType = ref('articles')
const selectedCategory = ref('')
const sortBy = ref('created_at')
const isSearching = ref(false)
const hasSearched = ref(false)
const searchTime = ref(0)
const currentPage = ref(1)
const searchResults = ref<PaginatedResponse<Article>>({
  items: [],
  total: 0,
  page: 1,
  size: 10,
  total_pages: 0
})
const categories = ref<Category[]>([])

// 热门搜索词
const popularKeywords = ref([
  '育儿知识', '健康饮食', '早教方法', '亲子关系', '学习指导',
  '心理健康', '家庭教育', '成长发育', '安全防护', '兴趣培养'
])

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

// 获取文章摘要
const getArticleSummary = (article: Article) => {
  if (article.summary) {
    return article.summary.substring(0, 150) + (article.summary.length > 150 ? '...' : '')
  }
  if (article.content) {
    // 移除HTML标签
    const plainText = article.content.replace(/<[^>]*>/g, '')
    return plainText.substring(0, 150) + (plainText.length > 150 ? '...' : '')
  }
  return '暂无摘要'
}

// 高亮关键词
const highlightKeyword = (text: string) => {
  if (!currentSearchQuery.value || !text) return text
  const regex = new RegExp(`(${escapeRegExp(currentSearchQuery.value)})`, 'gi')
  return text.replace(regex, '<mark class="bg-yellow-200 px-1 rounded">$1</mark>')
}

// 转义正则表达式特殊字符
const escapeRegExp = (string: string) => {
  return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

// 执行搜索
const performSearch = async (page = 1) => {
  if (!searchQuery.value.trim()) return
  
  try {
    const startTime = performance.now()
    isSearching.value = true
    currentSearchQuery.value = searchQuery.value.trim()
    currentPage.value = page
    
    const params = {
      keyword: currentSearchQuery.value,
      page: page,
      size: 10,
      category_id: selectedCategory.value ? Number(selectedCategory.value) : undefined,
      sort: sortBy.value,
      status: 1 // 只搜索已发布的文章
    }
    
    const response = await ArticleApi.getArticleList(params)
    searchResults.value = response.data
    
    const endTime = performance.now()
    searchTime.value = Math.round(endTime - startTime)
    hasSearched.value = true
    
    // 更新URL参数
    router.replace({
      query: {
        ...route.query,
        q: currentSearchQuery.value,
        page: page.toString(),
        category: selectedCategory.value || undefined,
        sort: sortBy.value
      }
    })
    
  } catch (error) {
    console.error('搜索失败:', error)
  } finally {
    isSearching.value = false
  }
}

// 使用关键词搜索
const searchWithKeyword = (keyword: string) => {
  searchQuery.value = keyword
  performSearch()
}

// 跳转到指定页
const goToPage = (page: number) => {
  if (page >= 1 && page <= searchResults.value.total_pages) {
    performSearch(page)
  }
}

// 跳转到文章详情
const goToArticle = (articleId: number) => {
  router.push(`/articles/${articleId}`)
}

// 加载分类列表
const loadCategories = async () => {
  try {
    const response = await CategoryApi.getCategoryList()
    categories.value = Array.isArray(response.data) ? response.data : []
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

// 监听路由变化
watch(
  () => route.query,
  (newQuery) => {
    if (newQuery.q && typeof newQuery.q === 'string') {
      searchQuery.value = newQuery.q
      currentSearchQuery.value = newQuery.q
      if (newQuery.category) {
        selectedCategory.value = newQuery.category as string
      }
      if (newQuery.sort) {
        sortBy.value = newQuery.sort as string
      }
      if (newQuery.page) {
        currentPage.value = Number(newQuery.page) || 1
      }
      // 如果URL有搜索参数但还没搜索过，执行搜索
      if (!hasSearched.value) {
        performSearch(currentPage.value)
      }
    }
  },
  { immediate: true }
)

// 组件挂载
onMounted(() => {
  loadCategories()
  
  // 如果URL中有搜索参数，自动执行搜索
  if (route.query.q) {
    searchQuery.value = route.query.q as string
    performSearch()
  }
})
</script>

<style scoped>
mark {
  background-color: #fef3c7;
  color: inherit;
  padding: 0 0.25rem;
  border-radius: 0.25rem;
}
</style>