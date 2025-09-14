<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 导航栏 -->
    <Navbar />

    <!-- 搜索和筛选 -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 sm:py-6">
      <div class="bg-white rounded-xl shadow-sm p-4 sm:p-6 border border-gray-100">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <!-- 搜索框 -->
          <div class="relative sm:col-span-2 lg:col-span-1">
            <SearchIcon class="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
            <input
              v-model="searchQuery"
              @input="handleSearch"
              type="text"
              placeholder="搜索文章标题、内容..."
              class="w-full pl-10 pr-4 py-2.5 sm:py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-300 text-sm sm:text-base placeholder-gray-500"
            />
          </div>
          
          <!-- 标签搜索 -->
          <div class="relative">
            <input
              v-model="searchTags"
              @input="handleTagSearch"
              type="text"
              placeholder="搜索标签（用逗号分隔）..."
              class="w-full px-4 py-2.5 sm:py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 transition-all duration-300 text-sm sm:text-base placeholder-gray-500"
            />
          </div>
          
          <!-- 分类筛选 -->
          <div class="relative">
            <select
              v-model="selectedCategory"
              @change="handleCategoryChange"
              class="w-full px-4 py-2.5 sm:py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-300 text-sm sm:text-base bg-white appearance-none cursor-pointer"
            >
              <option value="">全部分类</option>
              <option v-for="category in categories" :key="category.id" :value="category.id">
                {{ category.name }}
              </option>
            </select>
            <ChevronDownIcon class="absolute right-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400 pointer-events-none" />
          </div>
          
          <!-- 排序 -->
          <div class="relative">
            <select
              v-model="sortBy"
              @change="handleSortChange"
              class="w-full px-4 py-2.5 sm:py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-300 text-sm sm:text-base bg-white appearance-none cursor-pointer"
            >
              <option value="created_at">最新发布</option>
              <option value="updated_at">最近更新</option>
              <option value="views">浏览量</option>
              <option value="likes">点赞数</option>
            </select>
            <ChevronDownIcon class="absolute right-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400 pointer-events-none" />
          </div>
        </div>
      </div>
    </div>

    <!-- 文章列表 -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-12">
      <!-- 加载状态 -->
      <div v-if="isLoading" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-pink-600"></div>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error" class="text-center py-12">
        <div class="text-red-600 mb-4">{{ error }}</div>
        <button
          @click="loadArticles"
          class="bg-pink-600 text-white px-4 py-2 rounded-lg hover:bg-pink-700 transition-colors"
        >
          重试
        </button>
      </div>

      <!-- 文章网格 -->
      <div v-else-if="articles && articles.length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6">
        <article
          v-for="article in articles"
          :key="article.id"
          class="bg-white rounded-xl shadow-sm overflow-hidden hover:shadow-lg transition-all duration-300 cursor-pointer border border-gray-100 hover:border-gray-200 transform hover:scale-105"
          @click="goToArticle(article.id)"
        >
          <!-- 文章封面 -->
          <div class="aspect-video bg-gray-200 relative overflow-hidden">
            <img
              v-if="article.cover_image"
              :src="article.cover_image"
              :alt="article.title"
              class="w-full h-full object-cover transition-transform duration-300 hover:scale-110"
              @error="(e) => handleImageError(e, article)"
            />
            <div v-if="!article.cover_image || article._imageError" class="w-full h-full bg-gradient-to-br from-blue-400 via-purple-500 to-pink-500 flex items-center justify-center">
              <span class="text-white text-xl font-bold drop-shadow-lg">{{ article.title.charAt(0) }}</span>
            </div>
            <!-- 分类标签 -->
            <div class="absolute top-3 left-3">
              <span class="bg-white/95 backdrop-blur-sm text-gray-700 px-3 py-1 rounded-full text-xs font-medium shadow-sm">
                {{ getCategoryName(article.category_id) }}
              </span>
            </div>

          </div>
          
          <!-- 文章内容 -->
          <div class="p-4 sm:p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-2 line-clamp-2 hover:text-blue-600 transition-colors duration-300">
              {{ article.title }}
            </h3>
            <p class="text-gray-600 text-sm mb-3 line-clamp-3 leading-relaxed">
              {{ article.summary || article.content?.substring(0, 100) + '...' }}
            </p>
            
            <!-- 标签 -->
            <div v-if="article.tags" class="flex flex-wrap gap-1 mb-3">
              <span
                v-for="tag in getArticleTags(article.tags)"
                :key="tag"
                class="inline-flex items-center px-2 py-1 text-xs font-medium bg-purple-100 text-purple-700 rounded-full hover:bg-purple-200 cursor-pointer transition-colors"
                @click.stop="searchByTag(tag)"
              >
                #{{ tag }}
              </span>
            </div>
            
            <!-- 文章交互区 -->
            <div class="flex items-center justify-between">
              <!-- 互动按钮 -->
              <div class="flex items-center space-x-2">
                <!-- 浏览量 -->
                <div class="flex items-center px-2.5 py-1.5 bg-blue-50 text-blue-700 rounded-full text-xs font-medium hover:bg-blue-100 transition-all duration-300">
                  <EyeIcon class="h-3.5 w-3.5 mr-1" />
                  <span>{{ article.view_count || 0 }}</span>
                </div>
                
                <!-- 点赞 -->
                <button class="flex items-center px-2.5 py-1.5 bg-red-50 text-red-700 rounded-full text-xs font-medium hover:bg-red-100 hover:scale-105 transition-all duration-300 group" @click.stop="quickLike(article)">
                  <HeartIcon class="h-3.5 w-3.5 mr-1 group-hover:scale-110 transition-transform" />
                  <span>{{ article.like_count || 0 }}</span>
                </button>
                
                <!-- 评论 -->
                <div class="flex items-center px-2.5 py-1.5 bg-green-50 text-green-700 rounded-full text-xs font-medium hover:bg-green-100 transition-all duration-300">
                  <MessageCircleIcon class="h-3.5 w-3.5 mr-1" />
                  <span>{{ article.comment_count || 0 }}</span>
                </div>

                <!-- 收藏数量 -->
                <div class="flex items-center px-2.5 py-1.5 bg-amber-50 text-amber-700 rounded-full text-xs font-medium hover:bg-amber-100 transition-all duration-300">
                  <BookmarkIcon class="h-3.5 w-3.5 mr-1" />
                  <span>{{ article.favorite_count || 0 }}</span>
                </div>
              </div>
              
              <!-- 作者信息和日期 -->
              <div class="flex items-center space-x-2">
                <!-- 作者信息 -->
                <div v-if="article.author" class="flex items-center space-x-1">
                  <router-link
                    v-if="article.author.username"
                    :to="`/users/${article.author.username}`"
                    @click.stop
                    class="w-5 h-5 rounded-full overflow-hidden flex-shrink-0 hover:ring-2 hover:ring-pink-300 transition-all"
                  >
                    <img
                      v-if="article.author.avatar"
                      :src="article.author.avatar"
                      :alt="article.author.nickname || article.author.username"
                      class="w-full h-full object-cover"
                    />
                    <div
                      v-else
                      class="w-full h-full bg-gradient-to-r from-blue-400 to-purple-500 flex items-center justify-center"
                    >
                      <span class="text-white font-bold text-xs">
                        {{ (article.author.nickname || article.author.username || 'U').charAt(0).toUpperCase() }}
                      </span>
                    </div>
                  </router-link>
                  <router-link
                    v-if="article.author.username"
                    :to="`/users/${article.author.username}`"
                    @click.stop
                    class="text-xs text-gray-600 hover:text-pink-600 transition-colors font-medium truncate max-w-16"
                  >
                    {{ article.author.nickname || article.author.username }}
                  </router-link>
                  <span v-else class="text-xs text-gray-600">
                    {{ article.author.nickname || article.author.username || '匿名用户' }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </article>
      </div>

      <!-- 空状态 -->
      <div v-else class="text-center py-12">
        <div class="text-gray-500 mb-4">暂无文章</div>
        <router-link
          v-if="authStore.isAuthenticated"
          to="/articles/create"
          class="bg-pink-600 text-white px-6 py-2 rounded-lg font-medium hover:bg-pink-700 transition-colors inline-flex items-center"
        >
          <PlusIcon class="h-5 w-5 mr-2" />
          发布第一篇文章
        </router-link>
      </div>

      <!-- 分页 -->
      <div v-if="totalPages > 1" class="flex justify-center mt-8">
        <nav class="flex items-center space-x-2">
          <button
            :disabled="currentPage === 1"
            @click="goToPage(currentPage - 1)"
            class="px-3 py-2 rounded-lg border border-gray-300 text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <ChevronLeftIcon class="h-5 w-5" />
          </button>
          
          <button
            v-for="page in visiblePages"
            :key="page"
            :class="[
              'px-3 py-2 rounded-lg border',
              page === currentPage
                ? 'bg-pink-600 text-white border-pink-600'
                : 'border-gray-300 text-gray-500 hover:bg-gray-50'
            ]"
            @click="goToPage(page)"
          >
            {{ page }}
          </button>
          
          <button
            :disabled="currentPage === totalPages"
            @click="goToPage(currentPage + 1)"
            class="px-3 py-2 rounded-lg border border-gray-300 text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <ChevronRightIcon class="h-5 w-5" />
          </button>
        </nav>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  SearchIcon,
  FilterIcon,
  EyeIcon,
  HeartIcon,
  UserIcon,
  PlusIcon,
  StarIcon,
  MessageCircleIcon,
  BookmarkIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  ChevronDownIcon,
  ShareIcon
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { ArticleApi } from '@/api/article'
import { CategoryApi } from '@/api/category'
import { FavoriteApi } from '@/api/favorite'
import type { Article, Category } from '@/api/types'
import Navbar from '@/components/Navbar.vue'

const router = useRouter()
const authStore = useAuthStore()

// 响应式数据
const articles = ref<Article[]>([])
const categories = ref<Category[]>([])
const isLoading = ref(false)
const error = ref('')
const bookmarkStatus = ref<Record<number, boolean>>({})

// 搜索和筛选
const searchQuery = ref('')
const searchTags = ref('')
const selectedCategory = ref('')
const sortBy = ref('created_at')

// 分页
const currentPage = ref(1)
const pageSize = ref(12)
const totalCount = ref(0)
const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value))

// 计算可见的页码
const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, currentPage.value + 2)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  
  return pages
})

// 防抖搜索
let searchTimeout: number
const handleSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
    loadArticles()
  }, 500)
}

// 防抖标签搜索
let tagSearchTimeout: number
const handleTagSearch = () => {
  clearTimeout(tagSearchTimeout)
  tagSearchTimeout = setTimeout(() => {
    currentPage.value = 1
    loadArticles()
  }, 500)
}

// 分类变化处理
const handleCategoryChange = () => {
  currentPage.value = 1
  loadArticles()
}

// 排序变化处理
const handleSortChange = () => {
  currentPage.value = 1
  loadArticles()
}

// 获取分类名称
const getCategoryName = (categoryId: number) => {
  const category = categories.value.find(c => c.id === categoryId)
  return category?.name || '未分类'
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

// 跳转到文章详情
const goToArticle = (articleId: number) => {
  router.push(`/articles/${articleId}`)
}

// 获取文章标签数组
const getArticleTags = (tags: string) => {
  return tags ? tags.split(',').map(tag => tag.trim()).filter(tag => tag) : []
}

// 点击标签进行搜索
const searchByTag = (tag: string) => {
  searchTags.value = tag
  currentPage.value = 1
  loadArticles()
}

// 跳转到指定页面
const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    loadArticles()
  }
}

// 加载收藏状态
const loadBookmarkStatus = async (articleIds: number[]) => {
  if (!authStore.isAuthenticated || articleIds.length === 0) {
    return
  }
  
  try {
    const response = await FavoriteApi.batchGetFavoriteStatus({ article_ids: articleIds })
    if (response.success && response.data) {
      bookmarkStatus.value = { ...bookmarkStatus.value, ...response.data }
    }
  } catch (error) {
    console.error('加载收藏状态失败:', error)
  }
}

// 加载文章列表
const loadArticles = async () => {
  try {
    isLoading.value = true
    error.value = ''
    
    const params = {
      page: currentPage.value,
      size: pageSize.value,
      keyword: searchQuery.value || undefined,
      tags: searchTags.value || undefined,
      category_id: selectedCategory.value ? Number(selectedCategory.value) : undefined,
      sort: sortBy.value,
      status: 1 // 只显示已发布的文章（1=已发布）
    }
    
    const response = await ArticleApi.getArticleList(params)
    // Handle nested data structure from API
    if (response.data && response.data.list) {
      articles.value = response.data.list || []
      totalCount.value = response.data.total || 0
    } else {
      articles.value = (response.data && Array.isArray(response.data)) ? response.data : []
      totalCount.value = response.total || 0
    }
    
    // 加载收藏状态
    const articleIds = articles.value.map(article => article.id)
    await loadBookmarkStatus(articleIds)
    
  } catch (err: any) {
    error.value = err.message || '加载文章失败'
    console.error('加载文章失败:', err)
  } finally {
    isLoading.value = false
  }
}

// 加载分类列表
const loadCategories = async () => {
  try {
    const response = await CategoryApi.getCategoryList()
    categories.value = (response.data && Array.isArray(response.data)) ? response.data : []
  } catch (err) {
    console.error('加载分类失败:', err)
  }
}

// 图片错误处理
const handleImageError = (event: Event, article: any) => {
  console.warn('封面图片加载失败:', article.cover_image)
  // 标记这个文章的图片加载失败
  article._imageError = true
}

// 快速点赞
const quickLike = async (article: Article) => {
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }
  
  try {
    // 这里应该调用点赞API，暂时模拟
    // await ArticleApi.likeArticle(article.id)
    // 暂时增加点赞数用于演示
    article.like_count = (article.like_count || 0) + 1
  } catch (error) {
    console.error('点赞失败:', error)
  }
}

// 快速收藏
const quickBookmark = async (article: Article) => {
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }
  
  try {
    const response = await FavoriteApi.toggleFavorite({ article_id: article.id })
    
    if (response.success) {
      // 更新本地状态
      bookmarkStatus.value[article.id] = response.data.is_favorited
      
      // 更新文章收藏数
      const currentCount = article.favorite_count || 0
      if (response.data.is_favorited) {
        article.favorite_count = currentCount + 1
      } else {
        article.favorite_count = Math.max(0, currentCount - 1)
      }
    }
  } catch (error) {
    console.error('收藏操作失败:', error)
  }
}

// 快速分享
const quickShare = (article: Article) => {
  const shareUrl = `${window.location.origin}/articles/${article.id}`
  
  if (navigator.share) {
    navigator.share({
      title: article.title,
      text: article.summary || article.content?.substring(0, 100),
      url: shareUrl
    })
  } else {
    // 复制链接到剪贴板
    navigator.clipboard.writeText(shareUrl).then(() => {
      // 这里可以显示一个提示消息
    })
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadCategories()
  loadArticles()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>