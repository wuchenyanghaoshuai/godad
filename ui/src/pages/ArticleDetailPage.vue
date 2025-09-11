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
          <div class="flex items-center justify-between mb-6 pb-6 border-b border-gray-200">
            <div class="flex items-center space-x-4">
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
            
            <!-- 关注按钮 -->
            <div v-if="showFollowButton" class="flex items-center space-x-2">
              <button
                v-if="!isFollowing"
                @click="followAuthor"
                :disabled="isFollowLoading"
                class="flex items-center px-4 py-2 bg-pink-600 text-white rounded-lg hover:bg-pink-700 transition-colors text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <HeartIcon class="h-4 w-4 mr-1" />
                <span v-if="isFollowLoading">关注中...</span>
                <span v-else>关注</span>
              </button>
              
              <button
                v-else
                @click="unfollowAuthor"
                :disabled="isFollowLoading"
                class="flex items-center px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <HeartIcon class="h-4 w-4 mr-1 fill-current" />
                <span v-if="isFollowLoading">处理中...</span>
                <span v-else>已关注</span>
              </button>
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

          <!-- 互动操作区 -->
          <div class="mt-8 pt-6 border-t border-gray-200">
            <div class="bg-white rounded-lg shadow-sm p-4 border border-gray-200">
              <div class="flex items-center justify-between">
                <!-- 左侧：数据统计 -->
                <div class="flex items-center space-x-6">
                  <div class="flex items-center text-gray-600">
                    <EyeIcon class="h-4 w-4 mr-1.5" />
                    <span class="text-sm font-medium">{{ article.view_count || 0 }} 阅读</span>
                  </div>
                  <div class="flex items-center text-gray-600">
                    <HeartIcon class="h-4 w-4 mr-1.5" />
                    <span class="text-sm font-medium">{{ article.like_count || 0 }} 点赞</span>
                  </div>
                  <div class="flex items-center text-gray-600">
                    <MessageCircleIcon class="h-4 w-4 mr-1.5" />
                    <span class="text-sm font-medium">{{ article.comment_count || 0 }} 评论</span>
                  </div>
                </div>
                
                <!-- 右侧：操作按钮 -->
                <div class="flex items-center space-x-2">
                  <!-- 点赞按钮 -->
                  <button
                    @click="toggleLike"
                    :disabled="isLiking"
                    class="group flex items-center px-3 py-2 rounded-lg font-medium transition-all duration-300 transform hover:scale-105 active:scale-95 shadow-sm hover:shadow-md disabled:opacity-50 disabled:cursor-not-allowed text-sm"
                    :class="isLiked 
                      ? 'bg-red-500 text-white hover:bg-red-600' 
                      : 'bg-white text-gray-700 border border-gray-200 hover:border-red-300 hover:text-red-600 hover:bg-red-50'"
                  >
                    <HeartIcon 
                      class="h-4 w-4 mr-1 transition-transform duration-300 group-hover:scale-110" 
                      :class="isLiked ? 'fill-current' : 'group-hover:text-red-500'" 
                    />
                    <span>{{ isLiked ? '已点赞' : '点赞' }}</span>
                  </button>
                  
                  <!-- 评论按钮 -->
                  <button
                    @click="showComments = !showComments"
                    class="group flex items-center px-3 py-2 bg-white text-gray-700 border border-gray-200 rounded-lg font-medium transition-all duration-300 transform hover:scale-105 active:scale-95 shadow-sm hover:shadow-md hover:border-green-300 hover:text-green-600 hover:bg-green-50 text-sm"
                  >
                    <MessageCircleIcon class="h-4 w-4 mr-1 transition-transform duration-300 group-hover:scale-110 group-hover:text-green-500" />
                    <span>评论</span>
                  </button>
                  
                  <!-- 分享按钮 -->
                  <button
                    @click="shareArticle"
                    class="group flex items-center px-3 py-2 bg-white text-gray-700 border border-gray-200 rounded-lg font-medium transition-all duration-300 transform hover:scale-105 active:scale-95 shadow-sm hover:shadow-md hover:border-blue-300 hover:text-blue-600 hover:bg-blue-50 text-sm"
                  >
                    <ShareIcon class="h-4 w-4 mr-1 transition-transform duration-300 group-hover:scale-110 group-hover:rotate-12 group-hover:text-blue-500" />
                    <span>分享</span>
                  </button>
                  
                  <!-- 编辑按钮（作者可见） -->
                  <button
                    v-if="canEdit"
                    @click="$router.push(`/articles/${article.id}/edit`)"
                    class="group flex items-center px-3 py-2 bg-white text-gray-700 border border-gray-200 rounded-lg font-medium transition-all duration-300 transform hover:scale-105 active:scale-95 shadow-sm hover:shadow-md hover:border-emerald-300 hover:text-emerald-600 hover:bg-emerald-50 text-sm"
                  >
                    <EditIcon class="h-4 w-4 mr-1 transition-transform duration-300 group-hover:scale-110 group-hover:text-emerald-500" />
                    <span>编辑</span>
                  </button>
                  
                  <!-- 删除按钮（作者可见） -->
                  <button
                    v-if="canEdit"
                    @click="deleteArticle"
                    class="group flex items-center px-3 py-2 bg-white text-gray-700 border border-gray-200 rounded-lg font-medium transition-all duration-300 transform hover:scale-105 active:scale-95 shadow-sm hover:shadow-md hover:border-red-300 hover:text-red-600 hover:bg-red-50 text-sm"
                  >
                    <TrashIcon class="h-4 w-4 mr-1 transition-transform duration-300 group-hover:scale-110 group-hover:text-red-500" />
                    <span>删除</span>
                  </button>
                </div>
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
                    <span>{{ relatedArticle.view_count || 0 }}</span>
                  </div>
                  <div class="flex items-center space-x-1">
                    <HeartIcon class="w-3 h-3" />
                    <span>{{ relatedArticle.like_count || 0 }}</span>
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
import { FollowApi } from '@/api/follow'
import type { Article, Category } from '@/api/types'
import CommentSection from '@/components/CommentSection.vue'
import Navbar from '@/components/Navbar.vue'
import { useToast } from '@/composables/useToast'

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
const { toast } = useToast()

// 关注相关状态
const isFollowing = ref(false)
const isFollowLoading = ref(false)

// 计算属性
const canEdit = computed(() => {
  return authStore.isAuthenticated && 
         article.value && 
         (authStore.user?.id === article.value.author_id || authStore.isAdmin)
})

// 显示关注按钮的条件：用户已登录且文章存在且不是自己的文章
const showFollowButton = computed(() => {
  return authStore.isAuthenticated && 
         article.value && 
         authStore.user?.id !== article.value.author_id
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

// 返回文章列表页
const goBack = () => {
  router.push('/articles')
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
    
    const response = await ArticleApi.toggleLike(article.value.id)
    
    // 根据后端返回的数据判断点赞状态
    // 后端逻辑：点赞成功返回点赞对象，取消点赞返回null
    if (response.data !== null && response.data !== undefined) {
      // 返回了点赞对象，说明是点赞成功
      isLiked.value = true
      article.value.like_count = (article.value.like_count || 0) + 1
    } else {
      // 返回null，说明是取消点赞成功
      isLiked.value = false
      article.value.like_count = Math.max(0, (article.value.like_count || 1) - 1)
    }
  } catch (err: any) {
    console.error('点赞操作失败:', err)
    // 显示错误提示但不更改状态
    alert('点赞失败：' + (err.message || '未知错误'))
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

// 关注作者
const followAuthor = async () => {
  if (!article.value?.author_id) return
  
  try {
    isFollowLoading.value = true
    const response = await FollowApi.followUser(article.value.author_id)
    
    if (response.code === 200) {
      isFollowing.value = true
    }
  } catch (error: any) {
    // 如果是已经关注的错误，更新状态但不显示错误
    if (error.message && error.message.includes('already following')) {
      isFollowing.value = true
    } else {
      toast.error(error.message || '关注失败，请稍后重试')
    }
  } finally {
    isFollowLoading.value = false
  }
}

// 取消关注作者
const unfollowAuthor = async () => {
  if (!article.value?.author_id) return
  
  try {
    isFollowLoading.value = true
    const response = await FollowApi.unfollowUser(article.value.author_id)
    
    if (response.code === 200) {
      isFollowing.value = false
    }
  } catch (error: any) {
    toast.error(error.message || '取消关注失败，请稍后重试')
  } finally {
    isFollowLoading.value = false
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
    
    // 获取点赞状态 - 无论用户是否登录都尝试获取，让后端来判断
    try {
      const likeStatusResponse = await ArticleApi.getLikeStatus(articleId)
      isLiked.value = likeStatusResponse.data?.is_liked || false
    } catch (error) {
      // 如果API调用失败，默认设置为未点赞
      isLiked.value = false
    }
    
    // 获取关注状态 - 如果用户已登录且不是自己的文章
    if (authStore.isAuthenticated && article.value.author_id && authStore.user?.id !== article.value.author_id) {
      try {
        const followStatusResponse = await FollowApi.checkFollowStatus(article.value.author_id)
        // 后端直接返回 {"is_following": true}，不包装在data中
        isFollowing.value = followStatusResponse.is_following || false
      } catch (error) {
        // 如果API调用失败，默认设置为未关注
        isFollowing.value = false
      }
    } else {
      isFollowing.value = false
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
      score += Math.log(1 + (article.view_count || 0)) * 0.1
      score += (article.like_count || 0) * 0.2
      
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
      isLiked.value = false
      isFollowing.value = false
      // 重新加载数据
      loadArticle()
    }
  }
)

// 监听认证状态变化，重新加载点赞状态和关注状态
watch(
  () => authStore.isAuthenticated,
  (newVal) => {
    if (article.value) {
      // 无论登录还是登出都重新获取点赞状态，让后端判断
      ArticleApi.getLikeStatus(article.value.id)
        .then(response => {
          isLiked.value = response.data?.is_liked || false
        })
        .catch(error => {
          isLiked.value = false
        })
      
      // 重新获取关注状态
      if (newVal && article.value.author_id && authStore.user?.id !== article.value.author_id) {
        FollowApi.checkFollowStatus(article.value.author_id)
          .then(response => {
            // 后端直接返回 {"is_following": true}，不包装在data中
            isFollowing.value = response.is_following || false
          })
          .catch(error => {
            isFollowing.value = false
          })
      } else {
        isFollowing.value = false
      }
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

/* 文章内容中的链接样式 */
.article-content a {
  color: #2563eb;
  text-decoration: underline;
  text-decoration-color: rgba(37, 99, 235, 0.6);
  text-underline-offset: 2px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-weight: 500;
  position: relative;
}

.article-content a:hover {
  color: #1e40af;
  background-color: rgba(37, 99, 235, 0.1);
  text-decoration-color: #1e40af;
  border-radius: 3px;
  padding: 2px 4px;
  margin: 0 -2px;
}

/* 为链接添加小图标 */
.article-content a.editor-link::after {
  content: '🔗';
  display: inline;
  margin-left: 2px;
  font-size: 0.8em;
  opacity: 0.7;
  transition: opacity 0.2s ease;
}

.article-content a.editor-link:hover::after {
  opacity: 1;
}

/* 外部链接图标 */
.article-content a[target="_blank"]::before {
  content: '';
  display: inline-block;
  width: 0;
  height: 0;
  border-left: 3px solid currentColor;
  border-top: 3px solid transparent;
  border-bottom: 3px solid transparent;
  margin-right: 3px;
  opacity: 0.6;
  vertical-align: middle;
  transition: opacity 0.2s ease;
}

.article-content a[target="_blank"]:hover::before {
  opacity: 1;
}
</style>