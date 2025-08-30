<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 导航栏 -->
    <Navbar />
    <!-- 加载状态 -->
    <div v-if="isLoading" class="flex justify-center items-center min-h-screen">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-pink-600"></div>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="flex flex-col items-center justify-center min-h-screen">
      <div class="text-red-600 mb-4 text-lg">{{ error }}</div>
      <button
        @click="loadArticle"
        class="bg-pink-600 text-white px-6 py-2 rounded-lg hover:bg-pink-700 transition-colors"
      >
        重试
      </button>
    </div>

    <!-- 文章内容 -->
    <div v-else-if="article" class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
      <!-- 返回按钮 -->
      <div class="mb-6">
        <button
          @click="goBack"
          class="flex items-center text-gray-600 hover:text-blue-600 transition-colors duration-300 bg-white/80 backdrop-blur-sm px-4 py-2 rounded-full shadow-sm hover:shadow-md border border-gray-200 hover:border-blue-300"
        >
          <ArrowLeftIcon class="h-4 w-4 sm:h-5 sm:w-5 mr-2" />
          <span class="text-sm sm:text-base">返回文章列表</span>
        </button>
      </div>

      <!-- 文章头部 -->
      <article class="bg-white rounded-xl shadow-sm overflow-hidden mb-6 sm:mb-8 border border-gray-100">
        <!-- 文章封面 -->
        <div v-if="article.cover_image && !imageLoadError" class="h-48 sm:h-64 md:h-80 relative overflow-hidden">
          <img
            :src="article.cover_image"
            :alt="article.title"
            class="w-full h-full object-cover transition-transform duration-500 hover:scale-105"
            @error="handleImageError"
            @load="handleImageLoad"
          />
          <!-- 渐变遮罩 -->
          <div class="absolute inset-0 bg-gradient-to-t from-black/30 via-transparent to-transparent"></div>
          <div class="absolute bottom-6 left-6 right-6">
            <div class="flex items-center space-x-2 mb-2">
              <span class="bg-gradient-to-r from-blue-500 to-purple-600 text-white px-3 py-1.5 rounded-full text-xs sm:text-sm font-medium shadow-sm">
                {{ getCategoryName(article.category_id) }}
              </span>
              <!-- 精选标签已移除 -->
            </div>
            <h1 class="text-xl sm:text-2xl md:text-3xl lg:text-4xl font-bold text-white mb-2 leading-tight">
              {{ article.title }}
            </h1>
          </div>
        </div>

        <!-- 文章信息 -->
        <div class="p-4 sm:p-6 md:p-8">
          <!-- 无封面时的标题 -->
          <div v-if="!article.cover_image" class="mb-6">
            <div class="flex items-center space-x-2 mb-4">
              <span class="bg-pink-100 text-pink-600 px-3 py-1 rounded-full text-sm font-medium">
                {{ getCategoryName(article.category_id) }}
              </span>
              <!-- 精选标签已移除 -->
            </div>
            <h1 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              {{ article.title }}
            </h1>
          </div>

          <!-- 作者和发布信息 -->
          <div class="flex flex-col md:flex-row md:items-center md:justify-between mb-6 pb-6 border-b border-gray-200">
            <div class="flex items-center space-x-4 mb-4 md:mb-0">
              <div class="w-12 h-12 bg-gradient-to-r from-blue-400 to-purple-500 rounded-full flex items-center justify-center shadow-sm">
                <UserIcon class="h-6 w-6 text-white" />
              </div>
              <div>
                <div class="font-medium text-gray-900">
                  {{ article.author?.username || '匿名用户' }}
                </div>
                <div class="text-sm text-gray-500">
                  发布于 {{ formatDate(article.created_at) }}
                  <span v-if="article.updated_at !== article.created_at">
                    · 更新于 {{ formatDate(article.updated_at) }}
                  </span>
                </div>
              </div>
            </div>
            
            <!-- 文章统计 -->
            <div class="grid grid-cols-3 gap-4 md:flex md:items-center md:space-x-6 md:gap-0">
              <div class="flex items-center justify-center md:justify-start bg-gray-50 rounded-lg p-3 md:bg-transparent md:p-0">
                <EyeIcon class="h-4 w-4 mr-2 text-blue-500" />
                <div class="text-center md:text-left">
                  <div class="text-sm font-semibold text-gray-900">{{ article.views || 0 }}</div>
                  <div class="text-xs text-gray-500">阅读</div>
                </div>
              </div>
              <div class="flex items-center justify-center md:justify-start bg-gray-50 rounded-lg p-3 md:bg-transparent md:p-0">
                <HeartIcon class="h-4 w-4 mr-2 text-red-500" />
                <div class="text-center md:text-left">
                  <div class="text-sm font-semibold text-gray-900">{{ article.likes || 0 }}</div>
                  <div class="text-xs text-gray-500">点赞</div>
                </div>
              </div>
              <div class="flex items-center justify-center md:justify-start bg-gray-50 rounded-lg p-3 md:bg-transparent md:p-0">
                <MessageCircleIcon class="h-4 w-4 mr-2 text-green-500" />
                <div class="text-center md:text-left">
                  <div class="text-sm font-semibold text-gray-900">{{ article.comment_count || 0 }}</div>
                  <div class="text-xs text-gray-500">评论</div>
                </div>
              </div>
            </div>
          </div>

          <!-- 文章摘要 -->
          <div v-if="article.summary" class="mb-8">
            <h2 class="text-lg font-semibold text-gray-800 mb-3 flex items-center">
              <span class="w-1 h-5 bg-gradient-to-b from-blue-500 to-purple-600 rounded-full mr-3"></span>
              文章摘要
            </h2>
            <div class="bg-gradient-to-r from-blue-50 to-purple-50 p-4 sm:p-6 rounded-xl border-l-4 border-blue-500">
              <p class="text-gray-700 leading-relaxed text-sm sm:text-base italic">{{ article.summary }}</p>
            </div>
          </div>

          <!-- 文章内容 -->
          <div class="mb-8">
            <h2 class="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <span class="w-1 h-5 bg-gradient-to-b from-green-500 to-blue-600 rounded-full mr-3"></span>
              正文内容
            </h2>
            <div class="bg-white rounded-xl border border-gray-200 p-6 sm:p-8">
              <div class="prose prose-sm sm:prose-base lg:prose-lg max-w-none leading-relaxed">
                <div v-if="article.content" v-html="formatContent(article.content)" class="article-content"></div>
                <div v-else class="text-gray-500 italic text-center py-8">
                  暂无内容
                </div>
              </div>
            </div>
          </div>

          <!-- 文章标签 -->
          <div v-if="article.tags && getArticleTags(article.tags).length > 0" class="mt-8 pt-6 border-t border-gray-200">
            <h3 class="text-base sm:text-lg font-semibold text-gray-900 mb-3 flex items-center">
              <span class="w-1 h-6 bg-gradient-to-b from-blue-500 to-purple-600 rounded-full mr-3"></span>
              标签
            </h3>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="tag in getArticleTags(article.tags)"
                :key="tag"
                class="bg-gradient-to-r from-gray-100 to-gray-200 hover:from-blue-100 hover:to-purple-100 text-gray-700 hover:text-blue-700 px-3 py-1.5 rounded-full text-xs sm:text-sm transition-all duration-300 cursor-pointer shadow-sm hover:shadow-md"
              >
                #{{ tag }}
              </span>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="mt-8 pt-6 border-t border-gray-200">
            <div class="bg-white rounded-xl shadow-sm p-4 sm:p-6 mb-6 sm:mb-8 border border-gray-100">
              <div class="grid grid-cols-2 sm:flex sm:flex-wrap gap-3 sm:gap-4">
                <button
                  @click="toggleLike"
                  :disabled="isLiking"
                  class="flex items-center justify-center sm:justify-start px-4 py-3 rounded-xl transition-all duration-300 transform hover:scale-105 shadow-sm hover:shadow-md"
                  :class="isLiked ? 'bg-gradient-to-r from-red-500 to-pink-500 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
                >
                  <HeartIcon class="h-4 w-4 sm:h-5 sm:w-5 mr-2" :class="isLiked ? 'fill-current' : ''" />
                  <div class="text-center sm:text-left">
                    <div class="text-xs sm:text-sm font-medium">{{ isLiked ? '已点赞' : '点赞' }}</div>
                    <div class="text-xs opacity-75">({{ article.likes || 0 }})</div>
                  </div>
                </button>
                
                <button
                  @click="shareArticle"
                  class="flex items-center justify-center sm:justify-start px-4 py-3 bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-xl hover:from-blue-600 hover:to-purple-600 transition-all duration-300 transform hover:scale-105 shadow-sm hover:shadow-md"
                >
                  <ShareIcon class="h-4 w-4 sm:h-5 sm:w-5 mr-2" />
                  <span class="text-xs sm:text-sm font-medium">分享</span>
                </button>
                
                <button
                  v-if="canEdit"
                  @click="$router.push(`/articles/${article.id}/edit`)"
                  class="flex items-center justify-center sm:justify-start px-4 py-3 bg-gradient-to-r from-green-500 to-emerald-500 text-white rounded-xl hover:from-green-600 hover:to-emerald-600 transition-all duration-300 transform hover:scale-105 shadow-sm hover:shadow-md"
                >
                  <EditIcon class="h-4 w-4 sm:h-5 sm:w-5 mr-2" />
                  <span class="text-xs sm:text-sm font-medium">编辑</span>
                </button>
                
                <button
                  v-if="canEdit"
                  @click="deleteArticle"
                  class="flex items-center justify-center sm:justify-start px-4 py-3 bg-gradient-to-r from-red-500 to-rose-500 text-white rounded-xl hover:from-red-600 hover:to-rose-600 transition-all duration-300 transform hover:scale-105 shadow-sm hover:shadow-md"
                >
                  <TrashIcon class="h-4 w-4 sm:h-5 sm:w-5 mr-2" />
                  <span class="text-xs sm:text-sm font-medium">删除</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </article>

      <!-- 评论区域 -->
      <div class="mt-6 sm:mt-8">
        <!-- 评论标题栏 - 紧凑化 -->
        <div class="flex items-center justify-between mb-3 sm:mb-4 p-4 bg-gray-50 rounded-lg border border-gray-200">
          <div class="flex items-center">
            <MessageCircleIcon class="w-5 h-5 mr-2 text-blue-500" />
            <h2 class="text-lg font-semibold text-gray-900">评论</h2>
            <span class="ml-2 px-2 py-1 text-xs bg-blue-100 text-blue-600 rounded-full">
              {{ article.comment_count || 0 }}
            </span>
          </div>
          <button
            v-if="!showComments"
            @click="showComments = true"
            class="text-sm text-blue-600 hover:text-blue-800 transition-colors"
          >
            {{ article.comment_count > 0 ? '展开评论' : '发表评论' }}
          </button>
          <button
            v-else
            @click="showComments = false"
            class="text-sm text-gray-500 hover:text-gray-700 transition-colors"
          >
            收起评论
          </button>
        </div>

        <!-- 评论内容区 - 可折叠 -->
        <div 
          v-show="showComments" 
          class="transition-all duration-300 ease-in-out"
          :class="showComments ? 'opacity-100' : 'opacity-0'"
        >
          <CommentSection
            :article-id="article.id"
            @comment-added="handleCommentAdded"
            @comment-deleted="handleCommentDeleted"
          />
        </div>

        <!-- 无评论时的空状态提示 -->
        <div v-if="!showComments && article.comment_count === 0" class="text-center py-6 text-gray-500">
          <MessageCircleIcon class="w-12 h-12 mx-auto mb-3 text-gray-300" />
          <p class="text-sm">还没有人评论，点击"发表评论"来留下第一个评论吧！</p>
        </div>
      </div>

      <!-- 相关文章推荐 -->
      <div v-if="relatedArticles.length > 0" class="mt-8 sm:mt-12">
        <div class="flex items-center mb-4 sm:mb-6">
          <div class="w-1 h-6 bg-gradient-to-b from-pink-500 to-orange-500 rounded-full mr-3"></div>
          <h2 class="text-xl sm:text-2xl font-bold text-gray-900">相关文章</h2>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6">
          <article
            v-for="relatedArticle in relatedArticles"
            :key="relatedArticle.id"
            class="group bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden hover:shadow-lg hover:border-pink-200 transition-all duration-300 cursor-pointer transform hover:scale-105"
            @click="goToArticle(relatedArticle.id)"
          >
            <div class="relative h-32 sm:h-36 bg-gradient-to-br from-pink-200 via-purple-200 to-orange-200 overflow-hidden">
              <img
                v-if="relatedArticle.cover_image"
                :src="relatedArticle.cover_image"
                :alt="relatedArticle.title"
                class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
            </div>
            <div class="p-3 sm:p-4">
              <h3 class="font-semibold text-gray-900 mb-2 line-clamp-2 text-sm sm:text-base group-hover:text-pink-600 transition-colors duration-200">
                {{ relatedArticle.title }}
              </h3>
              <div class="flex items-center justify-between text-xs sm:text-sm text-gray-500">
                <div class="flex items-center space-x-1">
                  <UserIcon class="w-3 h-3 sm:w-4 sm:h-4" />
                  <span class="truncate max-w-20 sm:max-w-none">{{ relatedArticle.author?.username }}</span>
                </div>
                <span class="text-xs">{{ formatDate(relatedArticle.created_at) }}</span>
              </div>
              <div class="flex items-center justify-between mt-2 pt-2 border-t border-gray-100">
                <div class="flex items-center space-x-3 text-xs text-gray-400">
                  <div class="flex items-center space-x-1">
                    <EyeIcon class="w-3 h-3" />
                    <span>{{ relatedArticle.views || 0 }}</span>
                  </div>
                  <div class="flex items-center space-x-1">
                    <HeartIcon class="w-3 h-3" />
                    <span>{{ relatedArticle.likes || 0 }}</span>
                  </div>
                </div>
                <div class="text-xs text-pink-500 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                  阅读更多 →
                </div>
              </div>
            </div>
          </article>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowLeftIcon,
  UserIcon,
  EyeIcon,
  MessageCircleIcon,
  HeartIcon,
  StarIcon,
  ShareIcon,
  EditIcon,
  TrashIcon
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { ArticleApi } from '@/api/article'
import { CategoryApi } from '@/api/category'
import type { Article, Category } from '@/api/types'
import CommentSection from '@/components/CommentSection.vue'
import Navbar from '@/components/Navbar.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// 响应式数据
const article = ref<Article | null>(null)
const categories = ref<Category[]>([])
const relatedArticles = ref<Article[]>([])
const isLoading = ref(false)
const error = ref('')
const isLiking = ref(false)
const isLiked = ref(false)
const imageLoadError = ref(false)
const showComments = ref(false)

// 计算属性
const canEdit = computed(() => {
  return authStore.isAuthenticated && 
         article.value && 
         (authStore.user?.id === article.value.author_id || authStore.isAdmin)
})

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
    month: 'long',
    day: 'numeric'
  })
}

// 获取文章标签数组
const getArticleTags = (tags: string) => {
  return tags ? tags.split(',').map(tag => tag.trim()).filter(tag => tag) : []
}

// 格式化文章内容
const formatContent = (content: string | undefined | null) => {
  // 空值检查
  if (!content) {
    return ''
  }
  
  // 简单的 Markdown 转换（实际项目中建议使用专业的 Markdown 解析器）
  return content
    .replace(/\n/g, '<br>')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
    .replace(/`(.*?)`/g, '<code>$1</code>')
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 跳转到文章
const goToArticle = (articleId: number) => {
  router.push(`/articles/${articleId}`)
}

// 切换点赞状态
const toggleLike = async () => {
  if (!authStore.isAuthenticated || !article.value) return
  
  try {
    isLiking.value = true
    
    if (isLiked.value) {
      await ArticleApi.unlikeArticle(article.value.id)
      article.value.likes = (article.value.likes || 1) - 1
      isLiked.value = false
    } else {
      await ArticleApi.likeArticle(article.value.id)
      article.value.likes = (article.value.likes || 0) + 1
      isLiked.value = true
    }
  } catch (err: any) {
    console.error('点赞操作失败:', err)
    // 这里可以显示错误提示
  } finally {
    isLiking.value = false
  }
}

// 分享文章
const shareArticle = () => {
  if (navigator.share && article.value) {
    navigator.share({
      title: article.value.title,
      text: article.value.summary || article.value.content.substring(0, 100),
      url: window.location.href
    })
  } else {
    // 复制链接到剪贴板
    navigator.clipboard.writeText(window.location.href)
    // 这里可以显示复制成功的提示
    alert('链接已复制到剪贴板')
  }
}

// 删除文章
const deleteArticle = async () => {
  if (!article.value || !confirm('确定要删除这篇文章吗？此操作不可恢复。')) return
  
  try {
    await ArticleApi.deleteArticle(article.value.id)
    router.push('/articles')
  } catch (err: any) {
    console.error('删除文章失败:', err)
    alert('删除失败：' + (err.message || '未知错误'))
  }
}

// 评论事件处理
const handleCommentAdded = () => {
  // 评论添加后可以更新文章的评论数
  if (article.value) {
    article.value.comment_count = (article.value.comment_count || 0) + 1
    // 添加评论后自动展开评论区
    showComments.value = true
  }
}

const handleCommentDeleted = () => {
  // 评论删除后可以更新文章的评论数
  if (article.value && article.value.comment_count > 0) {
    article.value.comment_count = article.value.comment_count - 1
  }
}

// 图片错误处理
const handleImageError = () => {
  imageLoadError.value = true
  console.warn('封面图片加载失败:', article.value?.cover_image)
}

const handleImageLoad = () => {
  imageLoadError.value = false
}

// 加载文章详情
const loadArticle = async () => {
  const articleId = Number(route.params.id)
  if (!articleId) {
    error.value = '无效的文章ID'
    return
  }
  
  try {
    isLoading.value = true
    error.value = ''
    imageLoadError.value = false // 重置图片错误状态
    
    // 加载文章详情
    const response = await ArticleApi.getArticleDetail(articleId)
    article.value = response.data
    
    // 如果有评论，自动展开评论区
    showComments.value = (article.value.comment_count || 0) > 0
    
    // 增加浏览量 (暂时禁用，后端接口未实现)
    // ArticleApi.incrementViewCount(articleId).catch(console.error)
    
    // 检查是否已点赞（如果用户已登录）
    if (authStore.isAuthenticated) {
      // 这里需要后端提供检查点赞状态的接口
      // isLiked.value = await ArticleApi.checkLikeStatus(articleId)
    }
    
    // 加载相关文章
    loadRelatedArticles(articleId, article.value.category_id)
    
  } catch (err: any) {
    error.value = err.message || '加载文章失败'
    console.error('加载文章失败:', err)
  } finally {
    isLoading.value = false
  }
}

// 加载相关文章
const loadRelatedArticles = async (currentArticleId: number, categoryId: number) => {
  try {
    // 获取当前文章的标签
    const currentTags = article.value?.tags ? getArticleTags(article.value.tags) : []
    
    // 1. 先获取更多文章数据用于筛选（获取更多候选文章）
    const response = await ArticleApi.getArticleList({
      page: 1,
      size: 20, // 获取更多文章用于智能筛选
      status: 1 // 1 = published
    })
    
    if (!response.data || !Array.isArray(response.data)) {
      relatedArticles.value = []
      return
    }
    
    // 2. 过滤掉当前文章
    const candidateArticles = response.data.filter(a => a.id !== currentArticleId)
    
    // 3. 计算相关性得分
    const articlesWithScore = candidateArticles.map(article => {
      let score = 0
      
      // 同分类得分更高 (+3分)
      if (article.category_id === categoryId) {
        score += 3
      }
      
      // 标签匹配得分 (每个匹配的标签 +2分)
      if (article.tags && currentTags.length > 0) {
        const articleTags = getArticleTags(article.tags)
        const matchingTags = currentTags.filter(tag => articleTags.includes(tag))
        score += matchingTags.length * 2
      }
      
      // 浏览量和点赞数权重 (热门文章优先)
      score += Math.log(1 + (article.views || 0)) * 0.1
      score += (article.likes || 0) * 0.2
      
      return { ...article, relevanceScore: score }
    })
    
    // 4. 按相关性得分排序，取前3篇
    relatedArticles.value = articlesWithScore
      .filter(article => article.relevanceScore > 0) // 只显示有相关性的文章
      .sort((a, b) => b.relevanceScore - a.relevanceScore)
      .slice(0, 3)
      
  } catch (err) {
    console.error('加载相关文章失败:', err)
    relatedArticles.value = []
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

// 监听路由参数变化
watch(
  () => route.params.id,
  (newId, oldId) => {
    if (newId !== oldId) {
      // 重置状态
      article.value = null
      relatedArticles.value = []
      showComments.value = false
      // 重新加载数据
      loadArticle()
    }
  }
)

// 组件挂载时加载数据
onMounted(() => {
  loadCategories()
  loadArticle()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.prose {
  color: #374151;
  line-height: 1.75;
}

.prose h1,
.prose h2,
.prose h3,
.prose h4,
.prose h5,
.prose h6 {
  color: #111827;
  font-weight: 600;
  margin-top: 2rem;
  margin-bottom: 1rem;
}

.prose p {
  margin-bottom: 1.25rem;
}

.prose code {
  background-color: #f3f4f6;
  padding: 0.125rem 0.25rem;
  border-radius: 0.25rem;
  font-size: 0.875em;
}

.prose strong {
  font-weight: 600;
}

.prose em {
  font-style: italic;
}

/* 文章内容区域特殊样式 */
.article-content {
  line-height: 1.8;
  color: #374151;
  font-size: 16px;
}

.article-content h1,
.article-content h2,
.article-content h3,
.article-content h4,
.article-content h5,
.article-content h6 {
  color: #1f2937;
  font-weight: 600;
  margin-top: 2rem;
  margin-bottom: 1rem;
}

.article-content h1 {
  font-size: 1.875rem;
  border-bottom: 2px solid #e5e7eb;
  padding-bottom: 0.5rem;
}

.article-content h2 {
  font-size: 1.5rem;
}

.article-content h3 {
  font-size: 1.25rem;
}

.article-content p {
  margin-bottom: 1rem;
  text-align: justify;
}

.article-content ul,
.article-content ol {
  margin: 1rem 0;
  padding-left: 2rem;
}

.article-content li {
  margin-bottom: 0.5rem;
}

.article-content blockquote {
  border-left: 4px solid #3b82f6;
  background-color: #f8fafc;
  padding: 1rem 1.5rem;
  margin: 1.5rem 0;
  font-style: italic;
  color: #4b5563;
}

.article-content img {
  max-width: 100%;
  height: auto;
  border-radius: 0.5rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  margin: 1.5rem auto;
  display: block;
}
</style>